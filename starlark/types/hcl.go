package types

import (
	"fmt"
	"math/big"
	"regexp"
	"sort"
	"unicode"
	"unicode/utf8"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"go.starlark.net/starlark"
)

// HCLCompatible defines if the struct is suitable of by encoded in HCL.
type HCLCompatible interface {
	ToHCL(b *hclwrite.Body)
}

// BuiltinHCL returns a starlak.Builtin function to generate HCL from objects
// implementing the HCLCompatible interface.
//
//   outline: types
//     functions:
//       hcl(resource) string
//         Returns the HCL encoding of the given resource.
//         params:
//           resource <resource>
//             resource to be encoded.
//
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

// ToHCL honors the HCLCompatible interface.
func (s *Terraform) ToHCL(b *hclwrite.Body) {
	if s.b != nil {
		s.b.ToHCL(b)
	}

	s.p.ToHCL(b)
}

// ToHCL honors the HCLCompatible interface.
func (s *Dict) ToHCL(b *hclwrite.Body) {
	for _, v := range s.Keys() {
		p, _, _ := s.Get(v)
		hcl, ok := p.(HCLCompatible)
		if !ok {
			continue
		}

		hcl.ToHCL(b)
	}
}

// ToHCL honors the HCLCompatible interface.
func (s *Provider) ToHCL(b *hclwrite.Body) {
	block := b.AppendNewBlock("provider", []string{s.typ})

	block.Body().SetAttributeValue("alias", cty.StringVal(s.name))
	block.Body().SetAttributeValue("version", cty.StringVal(string(s.meta.Version)))
	s.Resource.doToHCLAttributes(block.Body())

	s.dataSources.ToHCL(b)
	s.resources.ToHCL(b)
	b.AppendNewline()
}

// ToHCL honors the HCLCompatible interface.
func (s *Provisioner) ToHCL(b *hclwrite.Body) {
	block := b.AppendNewBlock("provisioner", []string{s.typ})
	s.Resource.doToHCLAttributes(block.Body())
}

// ToHCL honors the HCLCompatible interface.
func (s *Backend) ToHCL(b *hclwrite.Body) {
	parent := b.AppendNewBlock("terraform", nil)

	block := parent.Body().AppendNewBlock("backend", []string{s.typ})
	s.Resource.doToHCLAttributes(block.Body())
	b.AppendNewline()
}

// ToHCL honors the HCLCompatible interface.
func (t *ResourceCollectionGroup) ToHCL(b *hclwrite.Body) {
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

// ToHCL honors the HCLCompatible interface.
func (c *ResourceCollection) ToHCL(b *hclwrite.Body) {
	for i := 0; i < c.Len(); i++ {
		c.Index(i).(*Resource).ToHCL(b)
	}
}

// ToHCL honors the HCLCompatible interface.
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

		tokens := appendTokensForValue(v.v, nil)
		body.SetAttributeRaw(v.Name, tokens)
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
	if len(r.dependencies) == 0 {
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

	l := len(r.dependencies)
	for i, dep := range r.dependencies {
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

var containsInterpolation = regexp.MustCompile(`(?mU)\$\{.*\}`)

func appendTokensForValue(val starlark.Value, toks hclwrite.Tokens) hclwrite.Tokens {
	switch v := val.(type) {
	case starlark.NoneType:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(`null`),
		})
	case starlark.Bool:
		var src []byte
		if v {
			src = []byte(`true`)
		} else {
			src = []byte(`false`)
		}
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: src,
		})
	case starlark.Float:
		bf := big.NewFloat(float64(v))
		srcStr := bf.Text('f', -1)
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenNumberLit,
			Bytes: []byte(srcStr),
		})
	case starlark.Int:
		srcStr := fmt.Sprintf("%d", v)
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenNumberLit,
			Bytes: []byte(srcStr),
		})
	case starlark.String:
		src := []byte(v.GoString())
		if !containsInterpolation.Match(src) {
			src = escapeQuotedStringLit(v.GoString())
		}

		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenOQuote,
			Bytes: []byte{'"'},
		})
		if len(src) > 0 {
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenQuotedLit,
				Bytes: []byte(src),
			})
		}
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenCQuote,
			Bytes: []byte{'"'},
		})
	case *starlark.List:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenOBrack,
			Bytes: []byte{'['},
		})

		for i := 0; i < v.Len(); i++ {
			if i > 0 {
				toks = append(toks, &hclwrite.Token{
					Type:  hclsyntax.TokenComma,
					Bytes: []byte{','},
				})
			}

			toks = appendTokensForValue(v.Index(i), toks)
		}

		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenCBrack,
			Bytes: []byte{']'},
		})

	case *starlark.Dict:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenOBrace,
			Bytes: []byte{'{'},
		})

		i := 0
		for _, eKey := range v.Keys() {
			if i > 0 {
				toks = append(toks, &hclwrite.Token{
					Type:  hclsyntax.TokenComma,
					Bytes: []byte{','},
				})
			}

			eVal, _, _ := v.Get(eKey)
			if hclsyntax.ValidIdentifier(eKey.(starlark.String).GoString()) {
				toks = append(toks, &hclwrite.Token{
					Type:  hclsyntax.TokenIdent,
					Bytes: []byte(eKey.(starlark.String).GoString()),
				})
			} else {
				toks = appendTokensForValue(eKey, toks)
			}
			toks = append(toks, &hclwrite.Token{
				Type:  hclsyntax.TokenEqual,
				Bytes: []byte{'='},
			})
			toks = appendTokensForValue(eVal, toks)
			i++
		}

		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenCBrace,
			Bytes: []byte{'}'},
		})
	case *Attribute:
		toks = append(toks, &hclwrite.Token{
			Type:  hclsyntax.TokenIdent,
			Bytes: []byte(v.sString.String()),
		})
	default:
		panic(fmt.Sprintf("cannot produce tokens for %#v", val))
	}

	return toks
}

func escapeQuotedStringLit(s string) []byte {
	if len(s) == 0 {
		return nil
	}
	buf := make([]byte, 0, len(s))
	for i, r := range s {
		switch r {
		case '\n':
			buf = append(buf, '\\', 'n')
		case '\r':
			buf = append(buf, '\\', 'r')
		case '\t':
			buf = append(buf, '\\', 't')
		case '"':
			buf = append(buf, '\\', '"')
		case '\\':
			buf = append(buf, '\\', '\\')
		case '$', '%':
			buf = appendRune(buf, r)
			remain := s[i+1:]
			if len(remain) > 0 && remain[0] == '{' {
				// Double up our template introducer symbol to escape it.
				buf = appendRune(buf, r)
			}
		default:
			if !unicode.IsPrint(r) {
				var fmted string
				if r < 65536 {
					fmted = fmt.Sprintf("\\u%04x", r)
				} else {
					fmted = fmt.Sprintf("\\U%08x", r)
				}
				buf = append(buf, fmted...)
			} else {
				buf = appendRune(buf, r)
			}
		}
	}
	return buf
}

func appendRune(b []byte, r rune) []byte {
	l := utf8.RuneLen(r)
	for i := 0; i < l; i++ {
		b = append(b, 0) // make room at the end of our buffer
	}
	ch := b[len(b)-l:]
	utf8.EncodeRune(ch, r)
	return b
}
