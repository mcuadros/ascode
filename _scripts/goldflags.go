package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rogpeppe/go-internal/modfile"
)

func main() {
	pkg := os.Args[1]
	path := os.Args[2]

	f, err := readGoMod(path)
	if err != nil {
		log.Fatal(err)
	}

	flags, err := getFlags(f, path, os.Args[3:])
	if err != nil {
		log.Fatal(err)
	}

	if pkg != "main" {
		pkg = filepath.Join(f.Module.Mod.Path, pkg)
	}

	fmt.Printf(renderLDFLAGS(pkg, flags))
}

func getFlags(f *modfile.File, path string, pkgs []string) (map[string]string, error) {
	var err error
	flags := make(map[string]string, 0)
	flags["build"] = time.Now().Format(time.RFC3339)
	flags["version"], flags["commit"], err = readVersion(path)

	if err != nil {
		return nil, err
	}

	for _, v := range pkgs {
		parts := strings.SplitN(v, "=", 2)
		key := parts[0]
		pkg := parts[1]

		flags[key] = getPackageVersion(f, pkg)
	}

	return flags, nil
}

func readVersion(path string) (string, string, error) {
	r, err := git.PlainOpen(path)
	if err != nil {
		return "", "", err
	}

	ref, err := r.Head()
	if err != nil {
		return "", "", err
	}

	if !ref.Name().IsBranch() {
		ref, err = findTag(r, ref.Hash())
		if err != nil {
			return "", "", err
		}
	}

	return ref.Name().Short(), ref.Hash().String()[:7], nil
}

func findTag(r *git.Repository, h plumbing.Hash) (*plumbing.Reference, error) {
	tagrefs, err := r.Tags()
	if err != nil {
		return nil, err
	}

	var match *plumbing.Reference
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		if t.Hash() == h {
			match = t
		}

		return nil
	})

	return match, err
}

func getVersionFromBranch(ref *plumbing.Reference) string {
	name := ref.Name().Short()
	pattern := "dev-%s"
	if name != "master" {
		pattern = fmt.Sprintf("dev-%s-%%s", name)
	}

	hash := ref.Hash().String()[:7]
	return fmt.Sprintf(pattern, hash)
}

func readGoMod(path string) (*modfile.File, error) {
	content, err := ioutil.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		return nil, err
	}

	return modfile.ParseLax("", content, nil)
}

func getPackageVersion(f *modfile.File, pkg string) string {
	for _, r := range f.Require {
		if r.Mod.Path == pkg {
			return r.Mod.Version
		}
	}
	return ""
}

func renderLDFLAGS(pkg string, flags map[string]string) string {
	output := make([]string, 0)
	for k, v := range flags {
		output = append(output, fmt.Sprintf("-X %s.%s=%s", pkg, k, v))
	}

	return strings.Join(output, " ")
}
