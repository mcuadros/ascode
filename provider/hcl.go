package provider

import (
	"fmt"

	"github.com/hashicorp/hcl2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	"go.starlark.net/starlark"
)

type HCLCompatible interface {
	ToHCL(b *hclwrite.Body)
}

func BuiltinToHCL(hcl HCLCompatible, f *hclwrite.File) starlark.Value {
	return starlark.NewBuiltin("to_hcl", func(_ *starlark.Thread, _ *starlark.Builtin, _ starlark.Tuple, _ []starlark.Tuple) (starlark.Value, error) {
		hcl.ToHCL(f.Body())
		return starlark.String(string(f.Bytes())), nil
	})
}

func (s *ProviderInstance) ToHCL(b *hclwrite.Body) {
	s.dataSources.ToHCL(b)
	s.resources.ToHCL(b)
}

func (t *MapSchemaIntance) ToHCL(b *hclwrite.Body) {
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
	if !r.nested {
		block = b.AppendNewBlock("resource", []string{r.typ, r.name})
	} else {
		block = b.AppendNewBlock(r.typ, nil)
	}

	body := block.Body()
	for k, attr := range r.block.Attributes {
		v, ok := r.values[k]
		if !ok {
			continue
		}

		body.SetAttributeValue(k, EncodeToCty(attr.Type, ValueToNative(v)))
	}

	for k := range r.block.BlockTypes {
		v, ok := r.values[k]
		if !ok {
			continue
		}

		if collection, ok := v.(HCLCompatible); ok {
			collection.ToHCL(block.Body())
		}
	}
}

func EncodeToCty(t cty.Type, v interface{}) cty.Value {
	switch value := v.(type) {
	case string:
		return cty.StringVal(value)
	case int64:
		return cty.NumberIntVal(value)
	case bool:
		return cty.BoolVal(value)
	case []interface{}:
		if len(value) == 0 {
			return cty.ListValEmpty(t)
		}

		values := make([]cty.Value, len(value))
		for i, v := range value {
			values[i] = EncodeToCty(t, v)
		}

		return cty.ListVal(values)
	default:
		return cty.StringVal(fmt.Sprintf("unhandled: %T", v))
	}
}
