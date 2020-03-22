package types

import (
	"bufio"
	"os"
	stdos "os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func TestExamples(t *testing.T) {
	pwd, _ := stdos.Getwd()
	defer func() {
		stdos.Chdir(pwd)
	}()

	stdos.Chdir(filepath.Join(pwd, "testdata/examples"))

	tests, err := filepath.Glob("*.star")
	assert.NoError(t, err)

	for _, test := range tests {
		doTestExample(t, test)
	}
}

func doTestExample(t *testing.T, filename string) {
	var output string
	printer := func(_ *starlark.Thread, msg string) {
		if output != "" {
			output += "\n"
		}

		output += msg
	}

	doTestPrint(t, filename, printer)
	expected := strings.TrimSpace(getExpectedFromExample(t, filename))

	assert.Equal(t, strings.TrimSpace(output), expected)
}

func getExpectedFromExample(t *testing.T, filename string) string {
	f, err := os.Open(filename)
	assert.NoError(t, err)
	defer f.Close()

	var expected []string
	scanner := bufio.NewScanner(f)
	var capture bool
	for scanner.Scan() {
		line := scanner.Text()
		if line == "# Output:" {
			capture = true
			continue
		}

		if !capture {
			continue
		}

		if len(line) >= 2 {
			line = line[2:]
		} else {
			line = ""
		}

		expected = append(expected, line)

	}

	return strings.Join(expected, "\n")
}
