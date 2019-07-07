package types

import (
	"fmt"

	"github.com/hashicorp/hcl2/hcl"
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

func (s *Provider) ToHCL(b *hclwrite.Body) {
	block := b.AppendNewBlock("provider", []string{s.name})
	block.Body().SetAttributeValue("alias", cty.StringVal(s.Name()))
	block.Body().SetAttributeValue("version", cty.StringVal(string(s.meta.Version)))
	s.Resource.doToHCLAttributes(block.Body())

	s.dataSources.ToHCL(b)
	s.resources.ToHCL(b)
}

func (t *MapSchema) ToHCL(b *hclwrite.Body) {
	for _, c := range t.collections {
		c.ToHCL(b)
	}
}

func (r *ResourceCollection) ToHCL(b *hclwrite.Body) {
	for i := 0; i < r.Len(); i++ {
		resource := r.Index(i).(*Resource)
		resource.ToHCL(b)
	}
}

func (r *Resource) ToHCL(b *hclwrite.Body) {
	if len(b.Blocks()) != 0 || len(b.Attributes()) != 0 {
		b.AppendNewline()
	}

	var block *hclwrite.Block
	if r.kind != NestedKind {
		block = b.AppendNewBlock(string(r.kind), []string{r.typ, r.Name()})
	} else {
		block = b.AppendNewBlock(r.typ, nil)
	}

	body := block.Body()

	if r.parent.kind == ProviderKind {
		body.SetAttributeTraversal("provider", hcl.Traversal{hcl.TraverseRoot{
			Name: fmt.Sprintf("%s.%s", r.parent.typ, r.parent.Name()),
		}})
	}

	r.doToHCLAttributes(body)
}

func (r *Resource) doToHCLAttributes(body *hclwrite.Body) {
	for k := range r.block.Attributes {
		v, ok := r.values[k]
		if !ok {
			continue
		}

		// TODO(mcuadros): I don't know how to do this properly, meanwhile, this works.
		if c, ok := v.v.(*Computed); ok {
			body.SetAttributeTraversal(k, hcl.Traversal{hcl.TraverseRoot{
				Name: c.String(),
			}})

			continue
		}

		body.SetAttributeValue(k, v.Cty())
	}

	for k := range r.block.BlockTypes {
		v, ok := r.values[k]
		if !ok {
			continue
		}

		if collection, ok := v.Value().(HCLCompatible); ok {
			collection.ToHCL(body)
		}
	}
}
