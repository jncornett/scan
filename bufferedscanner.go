package scan

type buffer struct {
	data  []byte
	empty bool
}

func (b *buffer) pop() []byte {
	b.empty = true
	return b.data
}

func (b *buffer) push(p []byte) {
	b.empty = false
	b.data = p
}

type bufferedScanner struct {
	s           Scanner
	back, front buffer
}

// NewBufferedScanner returns a BufferedScanner that is backed by s.
// Calls to BufferedScanner.Unscan that are unmatched by a preceding call to
// BufferedScanner.Scan will panic.
func NewBufferedScanner(s Scanner) BufferedScanner {
	return &bufferedScanner{s: s}
}

func (s bufferedScanner) Err() error {
	return s.s.Err()
}

func (s bufferedScanner) Bytes() []byte {
	return s.front.data
}

func (s bufferedScanner) Text() string {
	return string(s.Bytes())
}

func (s *bufferedScanner) Scan() bool {
	if s.back.empty {
		// need to fetch data from backing scanner
		if s.s.Scan() {
			s.back.push(copySlice(s.s.Bytes()))
		}
	}
	// now shift results to front buffer
	if s.back.empty {
		return false
	}
	s.front.push(s.back.pop())
	return true
}

func (s *bufferedScanner) Unscan() {
	if s.front.empty {
		panic("Unmatched call to Unscan or call to Unscan after Scan returned false")
	}
	s.back.push(s.front.pop())
}

func copySlice(p []byte) []byte {
	newSlice := make([]byte, len(p))
	copy(newSlice, p)
	return newSlice
}
