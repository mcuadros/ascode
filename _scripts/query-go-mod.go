package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rogpeppe/go-internal/modfile"
)

func main() {
	path := os.Args[1]
	pkg := os.Args[2]

	f, err := readGoMod(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range f.Require {
		if r.Mod.Path != pkg {
			continue
		}

		parts := strings.Split(r.Mod.Version, "-")
		if len(parts) > 1 {
			fmt.Println(parts[len(parts)-1])
		}
	}
}

func readGoMod(path string) (*modfile.File, error) {
	content, err := ioutil.ReadFile(filepath.Join(path, "go.mod"))
	if err != nil {
		return nil, err
	}

	return modfile.ParseLax("", content, nil)
}
