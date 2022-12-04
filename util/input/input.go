package util

import (
	"bufio"
	"io"
	"os"
	"sync"
)

type ScannerInterface interface {
	Scan() bool
	Text() string
	WriteString(str string)
}

type Scanner struct {
	*bufio.Scanner
}

func NewScanner() *Scanner {
	return &Scanner{bufio.NewScanner(os.Stdin)}
}

func (s *Scanner) Scan() bool {
	return s.Scanner.Scan()
}

func (s *Scanner) Text() string {
	return s.Scanner.Text()
}

func (s *Scanner) WriteString(str string) {
	// Do nothing
}

type MockScanner struct {
	*bufio.Scanner
	w      *io.PipeWriter
	buffer []byte
}

func NewMockScanner() *MockScanner {
	r, w := io.Pipe()
	return &MockScanner{bufio.NewScanner(r), w, []byte{}}
}

func (s *MockScanner) Scan() bool {
	var wg sync.WaitGroup
	wg.Add(2)
	var ret bool
	go func() {
		defer wg.Done()
		ret = s.Scanner.Scan()
	}()
	go func() {
		defer wg.Done()
		s.w.Write(s.buffer)
		s.w.Close()
		s.buffer = []byte{} // Flush the buffer
	}()
	wg.Wait()
	return ret
}

func (s *MockScanner) WriteString(str string) {
	s.buffer = append(s.buffer, []byte(str)...)
}
