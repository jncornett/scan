package scan_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/jncornett/scan"
)

// FIXME add test cases for error conditions
func TestBufferedScanner(t *testing.T) {
	text := "foo\n"
	s := scan.NewBufferedScanner(bufio.NewScanner(strings.NewReader(text)))
	if !s.Scan() {
		t.Fatal("Scan returned false on first call")
	}
	line1 := s.Text()
	s.Unscan()
	if !s.Scan() {
		t.Fatal("Scan returned false after Unscan")
	}
	line2 := s.Text()
	if line1 != line2 {
		t.Errorf("Expected %q to equal %q", line1, line2)
	}
}
