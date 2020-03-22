package types

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcl/hclsyntax"
	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"go.starlark.net/starlark"
)

type HCLCompatible interface {
	ToHCL(b *hclwrite.Body)
}

func BuiltinHCL() starlark.Value {
	return starlark.NewBuiltin("hcl", func(_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
		if args.Len() != 1 {
			return nil, fmt.Errorf("exactly one argument is required")
		}

		value := args.Index(0)
		hcl, ok := value.(HCLCompatible)
		if !ok {
			return nil, fmt.Errorf("value type %s doesn't support HCL conversion", value.Type())
		}

		f := hclwrite.NewEmptyFile()
		hcl.ToHCL(f.Body())
		return starlark.String(string(f.Bytes())), nil
	})
}

func (s *Terraform) ToHCL(b *hclwrite.Body) {
	if s.b != nil {
		s.b.ToHCL(b)
	}

	s.p.ToHCL(b)
}

func (s *AttrDict) ToHCL(b *hclwrite.Body) {
	for _, v := range s.Keys() {
		p, _, _ := s.Get(v)
		hcl, ok := p.(HCLCompatible)
		if !ok {
			continue
		}

		hcl.ToHCL(b)
	}
}

func (s *Provider) ToHCL(b *hclwrite.Body) {
	block := b.AppendNewBlock("provider", []string{s.typ})

	block.Body().SetAttributeValue("alias", cty.StringVal(s.name))
	block.Body().SetAttributeValue("version", cty.StringVal(string(s.meta.Version)))
	s.Resource.doToHCLAttributes(block.Body())

	s.dataSources.ToHCL(b)
	s.resources.ToHCL(b)
	b.AppendNewline()
}

func (s *Provisioner) ToHCL(b *hclwrite.Body) {
	block := b.AppendNewBlock("provisioner", []string{s.typ})
	s.Resource.doToHCLAttributes(block.Body())
}

func (s *Backend) ToHCL(b *hclwrite.Body) {
	parent := b.AppendNewBlock("terraform", nil)

	block := parent.Body().AppendNewBlock("backend", []string{s.typ})
	s.Resource.doToHCLAttributes(block.Body())
	b.AppendNewline()
}

func (t *MapSchema) ToHCL(b *hclwrite.Body) {
	names := make(sort.StringSlice, len(t.collections))
	var i int
	for name := range t.collections {
		names[i] = name
		i++
	}

	sort.Sort(names)
	for _, name := range names {
		t.collections[name].ToHCL(b)
	}
}

func (c *ResourceCollection) ToHCL(b *hclwrite.Body) {
	for i := 0; i < c.Len(); i++ {
		c.Index(i).(*Resource).ToHCL(b)
	}
}

func (r *Resource) ToHCL(b *hclwrite.Body) {
	if len(b.Blocks()) != 0 || len(b.Attributes()) != 0 {
		b.AppendNewline()
	}

	var block *hclwrite.Block
	if r.kind != NestedKind {
		labels := []string{r.typ, r.Name()}
		block = b.AppendNewBlock(string(r.kind), labels)
	} else {
		block = b.AppendNewBlock(r.typ, nil)
	}

	body := block.Body()

	if r.kind != NestedKind && r.parent != nil && r.parent.kind == ProviderKind {
		body.SetAttributeTraversal("provider", hcl.Traversal{
			hcl.TraverseRoot{Name: r.parent.typ},
			hcl.TraverseAttr{Name: r.parent.Name()},
		})
	}

	r.doToHCLAttributes(body)
	r.doToHCLDependencies(body)
	r.doToHCLProvisioner(body)
}

func (r *Resource) doToHCLAttributes(body *hclwrite.Body) {
	r.values.ForEach(func(v *NamedValue) error {
		if _, ok := r.block.Attributes[v.Name]; !ok {
			return nil
		}

		if c, ok := v.v.(*Computed); ok {
			body.SetAttributeTraversal(v.Name, hcl.Traversal{
				hcl.TraverseRoot{Name: c.String()},
			})

			return nil
		}

		body.SetAttributeValue(v.Name, v.Cty())
		return nil
	})

	r.values.ForEach(func(v *NamedValue) error {
		if _, ok := r.block.BlockTypes[v.Name]; !ok {
			return nil
		}

		if collection, ok := v.Starlark().(HCLCompatible); ok {
			collection.ToHCL(body)
		}

		return nil
	})
}

func (r *Resource) doToHCLDependencies(body *hclwrite.Body) {
	if len(r.dependenies) == 0 {
		return
	}

	toks := []*hclwrite.Token{}
	toks = append(toks, &hclwrite.Token{
		Type:  hclsyntax.TokenIdent,
		Bytes: []byte("depends_on"),
	})

	toks = append(toks, &hclwrite.Token{
		Type: hclsyntax.TokenEqual, Bytes: []byte{'='},
	}, &hclwrite.Token{
		Type: hclsyntax.TokenOBrack, Bytes: []byte{'['},
	})

	l := len(r.dependenies)
	for i, dep := range r.dependenies {
		name := fmt.Sprintf("%s.%s", dep.typ, dep.Name())
		toks = append(toks, &hclwrite.Token{
			Type: hclsyntax.TokenIdent, Bytes: []byte(name),
		})

		if i+1 == l {
			break
		}

		toks = append(toks, &hclwrite.Token{
			Type: hclsyntax.TokenComma, Bytes: []byte{','},
		})
	}

	toks = append(toks, &hclwrite.Token{
		Type: hclsyntax.TokenCBrack, Bytes: []byte{']'},
	})

	body.AppendUnstructuredTokens(toks)
	body.AppendNewline()
}

func (r *Resource) doToHCLProvisioner(body *hclwrite.Body) {
	if len(r.provisioners) == 0 {
		return
	}

	for _, p := range r.provisioners {
		if len(body.Blocks()) != 0 || len(body.Attributes()) != 0 {
			body.AppendNewline()
		}

		p.ToHCL(body)
	}
}
