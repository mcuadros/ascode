package main

import (
	"fmt"

	"github.com/hashicorp/terraform/configs/configschema"

	"go.starlark.net/starlark"
)

type NestedBlock struct {
	typ   string
	block *configschema.Block
	*PointerList
}

func NewNestedBlock(typ string, block *configschema.Block, refs *PointerList) *NestedBlock {
	return &NestedBlock{typ: typ, block: block, PointerList: refs}
}

// String honors the starlark.Value interface.
func (r *NestedBlock) String() string {
	return fmt.Sprintf("%s", r.typ)
}

// Type honors the starlark.Value interface.
func (r *NestedBlock) Type() string {
	return fmt.Sprintf("%s_collection", r.typ)
}

// Truth honors the starlark.Value interface.
func (r *NestedBlock) Truth() starlark.Bool {
	return true // even when empty
}

// Freeze honors the starlark.Value interface.
func (r *NestedBlock) Freeze() {}

// Hash honors the starlark.Value interface.
func (r *NestedBlock) Hash() (uint32, error) { return 42, nil }

// Name honors the starlark.Callable interface.
func (b *NestedBlock) Name() string {
	return b.typ
}

// CallInternal honors the starlark.Callable interface.
func (b *NestedBlock) CallInternal(thread *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String
	if len(args) != 0 {
		name = args.Index(0).(starlark.String)
	}

	resource, err := MakeResourceInstance(string(name), b.typ, b.block, kwargs)
	if err != nil {
		return nil, err
	}

	if err := b.PointerList.Append(resource); err != nil {
		return nil, err
	}

	return resource, nil
}
