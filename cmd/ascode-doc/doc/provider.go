package doc

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/ascode-dev/ascode/starlark/types"
	"github.com/b5/outline/lib"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type ResourceDocumentation struct {
	Name    string
	Type    string
	Attribs map[string]string
	Blocks  map[string]map[string]string
}

func NewResourceDocumentation(typ, name string) *ResourceDocumentation {
	return &ResourceDocumentation{
		Name:    name,
		Type:    typ,
		Attribs: make(map[string]string, 0),
		Blocks:  make(map[string]map[string]string, 0),
	}
}

var re = regexp.MustCompile(`\*[. ][\x60](.*)[\x60].*\) (.*)`)

// https://regex101.com/r/hINfBI/2
var blockRe = regexp.MustCompile(`^[^\*].*\x60(.*)\x60 block(?:s?)`)

func (r *ResourceDocumentation) Decode(doc io.Reader) error {
	var block string

	buf := bufio.NewReader(doc)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		parts := re.FindStringSubmatch(line)
		if len(parts) == 3 {
			if block == "" {
				r.AddAttrib(parts[1], parts[2])
				continue
			}

			r.AddBlockAttrib(block, parts[1], parts[2])
			continue
		}

		parts = blockRe.FindStringSubmatch(line)
		if len(parts) == 2 {
			block = parts[1]
		}
	}

	return nil
}

func (r *ResourceDocumentation) AddAttrib(name, desc string) {
	r.Attribs[name] = desc
}

func (r *ResourceDocumentation) AddBlockAttrib(block, name, desc string) {
	if _, ok := r.Blocks[block]; !ok {
		r.Blocks[block] = make(map[string]string, 0)
	}

	r.Blocks[block][name] = desc
}

type Documentation struct {
	name       string
	repository *git.Repository
	head       *object.Commit
	resources  map[string]map[string]string
}

func NewDocumentation(name string) (*Documentation, error) {
	d := &Documentation{
		name:      name,
		resources: make(map[string]map[string]string, 0),
	}

	return d, d.initRepository()
}

func (d *Documentation) initRepository() error {
	storer := memory.NewStorage()

	var err error
	d.repository, err = git.Clone(storer, nil, &git.CloneOptions{
		URL:   fmt.Sprintf("https://github.com/terraform-providers/terraform-provider-%s.git", d.name),
		Depth: 1,

		// as git does, when you make a clone, pull or some other operations the
		// server sends information via the sideband, this information can being
		// collected provinding a io.Writer to the CloneOptions options
		Progress: os.Stdout,
	})

	h, err := d.repository.Head()
	if err != nil {
		return err
	}

	d.head, err = d.repository.CommitObject(h.Hash())
	return err
}

func (d *Documentation) Resource(typ, name string) (*ResourceDocumentation, error) {
	parts := strings.SplitN(name, "_", 2)
	name = parts[1]

	filename := fmt.Sprintf("website/docs/%s/%s.html.md", typ, name)

	file, err := d.head.File(filename)
	if err != nil {
		return nil, err
	}

	r, err := file.Reader()
	if err != nil {
		return nil, err
	}

	resource := NewResourceDocumentation(typ, name)
	return resource, resource.Decode(r)
}

func (d *Documentation) Do(name string, schema providers.GetSchemaResponse) *lib.Doc {
	doc := &lib.Doc{}
	doc.Name = name
	doc.Path = name

	for name, schema := range schema.DataSources {
		doc.Types = append(doc.Types, d.schemaToDoc(name, &schema)...)
	}

	return doc
}

func (d *Documentation) schemaToDoc(resource string, s *providers.Schema) []*lib.Type {
	rd, err := d.Resource("d", resource)
	if err != nil {
		panic(err)
	}

	fmt.Println(resource, rd)
	typ := &lib.Type{}
	typ.Name = resource

	for name, attr := range s.Block.Attributes {
		typ.Fields = append(typ.Fields, d.attributeToField(rd.Attribs, name, attr))
	}

	types := []*lib.Type{typ}

	for name, block := range s.Block.BlockTypes {
		types = append(types, d.blockToType(rd, resource, name, block))
		typ.Fields = append(typ.Fields, d.blockToField(rd, resource, name, block))
	}

	return types
}

func (d *Documentation) blockToType(rd *ResourceDocumentation, resource, name string, block *configschema.NestedBlock) *lib.Type {
	typ := &lib.Type{}
	typ.Name = fmt.Sprintf("%s.%s", resource, name)

	for n, attr := range block.Attributes {
		//fmt.Println(rd.Blocks[name])

		typ.Fields = append(typ.Fields, d.attributeToField(rd.Blocks[name], n, attr))
	}

	return typ
}

func (d *Documentation) blockToField(rd *ResourceDocumentation, resource, name string, block *configschema.NestedBlock) *lib.Field {
	field := &lib.Field{}
	nested := fmt.Sprintf("%s.%s", resource, name)

	field.Name = name
	if block.MaxItems != 1 {
		field.Type = fmt.Sprintf("collection<%s>", nested)
	} else {
		field.Type = nested
	}

	return field
}

func (d *Documentation) attributeToField(doc map[string]string, name string, attr *configschema.Attribute) *lib.Field {
	field := &lib.Field{}
	field.Name = name
	field.Description, _ = doc[name]
	if attr.Computed && !attr.Optional {
		field.Type = "computed"
	} else {
		field.Type = types.MustTypeFromCty(attr.Type).Starlark()
	}

	return field
}
