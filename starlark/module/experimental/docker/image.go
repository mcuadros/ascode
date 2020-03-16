package docker

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/Masterminds/semver/v3"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/docker/reference"
	"github.com/containers/image/v5/types"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

const (
	// ModuleName defines the expected name for this Module when used
	// in starlark's load() function, eg: load('experimental/docker', 'docker')
	ModuleName = "experimental/docker"

	ImageFuncName = "image"

	latestTag = "lastest"
)

var (
	once         sync.Once
	dockerModule starlark.StringDict
)

// LoadModule loads the os module.
// It is concurrency-safe and idempotent.
//
//   outline: docker
//     path: experimental/docker
func LoadModule() (starlark.StringDict, error) {
	once.Do(func() {
		dockerModule = starlark.StringDict{
			"docker": &starlarkstruct.Module{
				Name: "docker",
				Members: starlark.StringDict{
					ImageFuncName: starlark.NewBuiltin(ImageFuncName, Image),
				},
			},
		}
	})

	return dockerModule, nil
}

type sString = starlark.String
type image struct {
	tags       []string
	ref        types.ImageReference
	constraint string
	sString
}

func Image(
	thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var image, constraint string
	err := starlark.UnpackArgs(ImageFuncName, args, kwargs, "image", &image, "constraint", &constraint)
	if err != nil {
		return nil, err
	}

	return newImage(image, constraint)
}

func newImage(name, constraint string) (*image, error) {
	ref, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return nil, err
	}

	if !reference.IsNameOnly(ref) {
		return nil, errors.New("no tag or digest allowed in reference")
	}

	dref, err := docker.NewReference(reference.TagNameOnly(ref))
	if err != nil {
		return nil, err
	}

	return &image{
		ref:        dref,
		constraint: constraint,
		sString:    starlark.String(ref.Name()),
	}, nil
}

func (i *image) Attr(name string) (starlark.Value, error) {
	switch name {
	case "name":
		return starlark.String(i.ref.DockerReference().Name()), nil
	case "tags":
		return i.getTags()
	case "version":
		return starlark.NewBuiltin("version", i.builtinVersionFunc), nil
	}

	return nil, nil
}

func (i *image) builtinVersionFunc(
	_ *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple,
) (starlark.Value, error) {

	var full bool
	starlark.UnpackArgs(ImageFuncName, args, kwargs, "full", &full)

	v, err := i.getVersion()
	if err != nil {
		return starlark.None, err
	}

	if full {
		v = fmt.Sprintf("%s:%s", i.ref.DockerReference().Name(), v)
	}

	return starlark.String(v), nil
}

func (i *image) getTags() (*starlark.List, error) {
	if len(i.tags) != 0 {
		return listToStarlark(i.tags), nil
	}

	var err error
	i.tags, err = docker.GetRepositoryTags(context.TODO(), nil, i.ref)
	if err != nil {
		return nil, fmt.Errorf("error listing repository tags: %v", err)
	}

	i.tags = sortTags(i.tags)
	return listToStarlark(i.tags), nil
}

func (i *image) getVersion() (string, error) {
	if i.constraint == latestTag {
		return latestTag, nil
	}

	_, err := i.getTags()
	if err != nil {
		return "", err
	}

	if len(i.tags) == 0 {
		return "", fmt.Errorf("no tags form this image")
	}

	c, err := semver.NewConstraint(i.constraint)
	if err != nil {
		return i.doGetVersionExactTag(i.constraint)
	}

	return i.doGetVersionWithConstraint(c)
}

func (i *image) doGetVersionWithConstraint(c *semver.Constraints) (string, error) {
	// it assumes tags are always sorted from higher to lower
	for _, tag := range i.tags {
		v, err := semver.NewVersion(tag)
		if err == nil {
			if c.Check(v) {
				return tag, nil
			}
		}
	}

	return "", nil
}

func (i *image) doGetVersionExactTag(expected string) (string, error) {
	for _, tag := range i.tags {
		if tag == expected {
			return tag, nil
		}
	}

	return "", fmt.Errorf("tag %q not found in repository", expected)
}

func sortTags(tags []string) []string {
	versions, others := listToVersion(tags)
	sort.Sort(sort.Reverse(semver.Collection(versions)))
	return versionToList(versions, others)
}

func listToStarlark(input []string) *starlark.List {
	output := make([]starlark.Value, len(input))
	for i, v := range input {
		output[i] = starlark.String(v)
	}

	return starlark.NewList(output)
}

func listToVersion(input []string) ([]*semver.Version, []string) {
	versions := make([]*semver.Version, 0)
	other := make([]string, 0)

	for _, text := range input {
		v, err := semver.NewVersion(text)
		if err == nil && v.Prerelease() == "" {
			versions = append(versions, v)
			continue
		}

		other = append(other, text)
	}

	return versions, other
}

func versionToList(versions []*semver.Version, other []string) []string {
	output := make([]string, 0)
	for _, v := range versions {
		output = append(output, v.Original())
	}

	return append(output, other...)
}
