package scan_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/jncornett/scan"
)

const ExampleText = `
package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, playground")
}
`

func TestFilterScanner(t *testing.T) {
	tests := []struct {
		name   string
		text   string
		filter scan.Filter
		lines  []string
	}{
		{
			name: "Example",
			text: ExampleText,
			filter: func(v scan.View) bool {
				text := v.Text()
				return len(text) != 0 && strings.IndexAny(text, " \t") != 0
			},
			lines: []string{
				"package main",
				"import (",
				")",
				"func main() {",
				"}",
			},
		},
		{
			name:   "Empty",
			text:   "",
			filter: func(scan.View) bool { return true },
			lines:  []string{},
		},
		{
			name:   "NilFilter",
			text:   ExampleText,
			filter: nil,
			lines: []string{
				"",
				"package main",
				"",
				"import (",
				"	\"fmt\"",
				")",
				"",
				"func main() {",
				"	fmt.Println(\"Hello, playground\")",
				"}",
			},
		},
		// TODO add more test cases
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := scan.FilterScanner{
				Scanner: bufio.NewScanner(strings.NewReader(test.text)),
				Filter:  test.filter,
			}
			for i, line := range test.lines {
				if !s.Scan() {
					t.Fatal("Scan returned false too early (at line ", i, ")")
					return
				}
				result := s.Text()
				if line != result {
					t.Errorf("Expected %v, got %v (line %v)", line, result, i)
				}
			}
			if s.Scan() {
				t.Fatal("Scanner has line(s) left over")
				return
			}
			err := s.Err()
			if err != nil {
				t.Error(err)
			}
		})
	}
}
