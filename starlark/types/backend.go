package types

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/addrs"
	"github.com/hashicorp/terraform/backend"
	binit "github.com/hashicorp/terraform/backend/init"
	"github.com/hashicorp/terraform/providers"
	"github.com/hashicorp/terraform/states"
	"github.com/hashicorp/terraform/states/statemgr"
	"github.com/mcuadros/ascode/terraform"
	"github.com/qri-io/starlib/util"
	"go.starlark.net/starlark"
)

func init() {
	binit.Init(nil)
}

// BuiltinBackend returns a starlak.Builtin function capable of instantiate
// new Backend instances.
//
//   outline: types
//     functions:
//       backend(type) Backend
//         Instantiates a new [`Backend`](#backend)
//
//         params:
//           type string
//             [Backend type](https://www.terraform.io/docs/backends/types/index.html).
//
func BuiltinBackend() starlark.Value {
	return starlark.NewBuiltin("backend", MakeBackend)
}

// MakeBackend defines the Backend constructor.
func MakeBackend(
	t *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var name starlark.String
	switch len(args) {
	case 1:
		var ok bool
		name, ok = args.Index(0).(starlark.String)
		if !ok {
			return nil, fmt.Errorf("expected string, got %s", args.Index(0).Type())
		}
	default:
		return nil, fmt.Errorf("unexpected positional arguments count")
	}

	pm := t.Local(PluginManagerLocal).(*terraform.PluginManager)
	p, err := NewBackend(pm, name.GoString())
	if err != nil {
		return nil, err
	}

	return p, p.loadKeywordArgs(kwargs)
}

// Backend represent a Terraform Backend.
//
//   outline: types
//     types:
//       Backend
//         The part of Terraform's core that determines how Terraform stores
//         state and performs operations (like plan, apply, import, etc.).
//         Terraform has multiple [backends](https://www.terraform.io/docs/backends/index.html)
//         to choose from, which can be configured in a variety of ways.
//
//         fields:
//           __kind__ string
//             Kind of the backend. Fixed value `backend`.
//           __type__ string
//             Type of the backend. Eg.: `local`.
//           __dict__ Dict
//             A dictionary containing all the config values of the backend.
//           <argument> <scalar>
//             Arguments defined by the backend schema, thus can be of any
//             scalar type.
//
//         methods:
//           state(module="", workspace="default") State
//             Loads the latest state for a given module or workspace.
//             params:
//               module string
//                 name of the module, empty equals to root.
//               workspace string
//                 backend workspace
type Backend struct {
	pm *terraform.PluginManager
	b  backend.Backend
	*Resource
}

var _ starlark.Value = &Backend{}
var _ starlark.HasAttrs = &Backend{}
var _ starlark.Comparable = &Backend{}

// NewBackend returns a new Backend instance based on given arguments,
func NewBackend(pm *terraform.PluginManager, typ string) (*Backend, error) {
	fn := binit.Backend(typ)
	if fn == nil {
		return nil, fmt.Errorf("unable to find backend %q", typ)
	}

	b := fn()

	return &Backend{
		pm:       pm,
		b:        b,
		Resource: NewResource("", typ, BackendKind, b.ConfigSchema(), nil, nil),
	}, nil
}

// Attr honors the starlark.HasAttrs interface.
func (b *Backend) Attr(name string) (starlark.Value, error) {
	switch name {
	case "state":
		return starlark.NewBuiltin("state", b.state), nil
	}

	return b.Resource.Attr(name)
}

// AttrNames honors the starlark.HasAttrs interface.
func (b *Backend) AttrNames() []string {
	return append(b.Resource.AttrNames(), "state")
}

func (b *Backend) getStateMgr(workspace string) (statemgr.Full, error) {
	values, diag := b.b.PrepareConfig(b.values.Cty(b.b.ConfigSchema()))
	if err := diag.Err(); err != nil {
		return nil, err
	}

	diag = b.b.Configure(values)
	if err := diag.Err(); err != nil {
		return nil, err
	}

	workspaces, err := b.b.Workspaces()
	if err != nil {
		return nil, err
	}

	var found bool
	for _, w := range workspaces {
		if w == workspace {
			found = true
			break
		}
	}

	if !found {
		return nil, fmt.Errorf("unable to find %q workspace", workspace)
	}

	return b.b.StateMgr(workspace)
}

func (b *Backend) state(
	_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple,
) (starlark.Value, error) {

	workspace := "default"
	module := ""

	err := starlark.UnpackArgs("state", args, kwargs, "module?", &module, "workspace?", &workspace)
	if err != nil {
		return nil, err
	}

	sm, err := b.getStateMgr(workspace)
	if err != nil {
		return nil, err
	}

	if err := sm.RefreshState(); err != nil {
		return nil, err
	}

	state := sm.State()
	if state == nil {
		return starlark.None, nil
	}

	return NewState(b.pm, module, state)

}

// String honors the starlark.Value interface.
func (b *Backend) String() string {
	return fmt.Sprintf("Backend<%s>", b.typ)
}

// Type honors the starlark.Value interface.
func (b *Backend) Type() string {
	return "Backend"
}

// State represents a Terraform state read by a backend.
// https://www.terraform.io/docs/state/index.html
//
//   outline: types
//     types:
//       State
//         Terraform's cached information about your managed infrastructure
//         and configuration. This [state](https://www.terraform.io/docs/state/index.html)
//         is used to persistently map the same real world resources to your
//         configuration from run-to-run, keep track of metadata, and improve
//         performance for large infrastructures.
//
//         State implements a Dict, where the first level are the providers
//         containing the keys `data` with the data sources and `resources` with
//         the resources.
//
//         examples:
//           backend_local.star
//             An example of how to print a resume of providers and resources
//             count from the state.
//
//         fields:
//           <provider> AttrDict
//             provider and all the resources state.
type State struct {
	*starlark.Dict
	pm *terraform.PluginManager
}

var _ starlark.Value = &State{}
var _ starlark.HasAttrs = &State{}
var _ starlark.Comparable = &State{}

// NewState returns a new instance of State based on the given arguments,
func NewState(pm *terraform.PluginManager, module string, state *states.State) (*State, error) {
	var mod *states.Module
	for _, m := range state.Modules {
		if m.Addr.String() == module {
			mod = m
		}
	}

	if mod == nil {
		return nil, fmt.Errorf("unable to find module with addr %q", module)
	}

	s := &State{
		Dict: starlark.NewDict(0),
		pm:   pm,
	}

	return s, s.initialize(state, mod)
}

func (s *State) initialize(state *states.State, mod *states.Module) error {
	providers := make(map[string]*Provider, 0)
	addrs := state.ProviderAddrs()
	for _, addr := range addrs {
		typ := addr.ProviderConfig.Type.Type
		p, err := NewProvider(s.pm, typ, "", addr.ProviderConfig.Alias)
		if err != nil {
			return err
		}

		providers[addr.ProviderConfig.String()] = p
	}

	for _, r := range mod.Resources {
		provider := r.ProviderConfig.String()
		if err := s.initializeResource(providers[provider], r); err != nil {
			return err
		}
	}

	return nil
}

func (s *State) initializeResource(p *Provider, r *states.Resource) error {
	typ := r.Addr.Type
	name := r.Addr.Name

	mode := addrsResourceModeString(r.Addr.Mode)

	var schema providers.Schema
	switch r.Addr.Mode {
	case addrs.DataResourceMode:
		schema = p.dataSources.schemas[typ]
	case addrs.ManagedResourceMode:
		schema = p.resources.schemas[typ]
	default:
		return fmt.Errorf("invalid resource type")
	}

	multi := r.EachMode != states.NoEach
	for _, instance := range r.Instances {
		r := NewResource(name, typ, ResourceKind, schema.Block, p, p.Resource)

		var val interface{}
		if err := json.Unmarshal(instance.Current.AttrsJSON, &val); err != nil {
			return err
		}

		values, _ := util.Marshal(val)
		if err := r.loadDict(values.(*starlark.Dict)); err != nil {
			return err
		}

		if err := s.set(mode, typ, name, r, multi); err != nil {
			return err
		}
	}

	return nil
}

func addrsResourceModeString(m addrs.ResourceMode) string {
	switch m {
	case addrs.ManagedResourceMode:
		return "resource"
	case addrs.DataResourceMode:
		return "data"
	}

	return ""
}
func (s *State) set(mode, typ, name string, r *Resource, multi bool) error {
	p := starlark.String(r.provider.typ)
	m := starlark.String(mode)
	t := starlark.String(typ[len(r.provider.typ)+1:])
	n := starlark.String(name)

	if _, ok, _ := s.Get(p); !ok {
		s.SetKey(p, starlark.NewDict(0))
	}

	providers, _, _ := s.Get(p)
	if _, ok, _ := providers.(*starlark.Dict).Get(m); !ok {
		providers.(*starlark.Dict).SetKey(m, starlark.NewDict(0))
	}

	modes, _, _ := providers.(*starlark.Dict).Get(m)
	if _, ok, _ := modes.(*starlark.Dict).Get(t); !ok {
		modes.(*starlark.Dict).SetKey(t, starlark.NewDict(0))
	}

	resources, _, _ := modes.(*starlark.Dict).Get(t)

	if !multi {
		return resources.(*starlark.Dict).SetKey(n, r)
	}

	if _, ok, _ := resources.(*starlark.Dict).Get(n); !ok {
		resources.(*starlark.Dict).SetKey(n, starlark.NewList(nil))
	}

	instances, _, _ := resources.(*starlark.Dict).Get(n)
	if err := instances.(*starlark.List).Append(r); err != nil {
		return err
	}

	return nil
}
