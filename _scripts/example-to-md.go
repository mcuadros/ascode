package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	md, err := exampleToMD(os.Args[1], os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(md)
}

func exampleToMD(filename string, weight string) (string, error) {
	b := bytes.NewBuffer(nil)
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}

	var isComment, isCodeBlock, isPrint bool
	var preLine string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		line := scanner.Text()
		curIsComment := len(line) > 0 && line[0] == '#'
		if isCodeBlock && len(preLine) == 0 && curIsComment {
			isPrint = false
		}

		if len(preLine) >= 3 && preLine[:3] == `"""` {
			isPrint = false
		}

		if isPrint {

			preLine = strings.Trim(preLine, "#")
			fmt.Fprintln(b, preLine)
			isPrint = false
		}

		preLine = line
		if curIsComment {
			line = strings.TrimSpace(line[1:])

			isComment = true
			if isCodeBlock {
				fmt.Fprintln(b, "```\n\n")
				isCodeBlock = false
			}

			if b.Len() == 0 {
				fmt.Fprintf(b,
					"---\ntitle: '%s'\nweight: %s\n---\n\n",
					line, weight,
				)
				continue
			}

			isPrint = true
			continue
		}

		if len(line) == 0 && isComment {
			isPrint = true
			continue
		}

		if isComment {
			isComment = false
			isCodeBlock = true

			if len(line) >= 3 && line[:3] == `"""` {
				fmt.Fprintln(b, "```"+line[3:])
				continue
			}

			fmt.Fprintln(b, "\n\n```python")

		}

		isPrint = true
	}

	if isCodeBlock {
		fmt.Fprintln(b, "```")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	return b.String(), nil
}
