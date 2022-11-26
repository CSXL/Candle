package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewScanner(t *testing.T) {
	var scanner = NewScanner()
	assert.NotNil(t, scanner)
}

func TestScanner_Scan(t *testing.T) {
	r, w, _ := os.Pipe()
	os.Stdin = r
	var scanner = NewScanner()
	w.WriteString("CSX Labs Rules!")
	w.Close()
	assert.True(t, scanner.Scan())
	assert.Equal(t, "CSX Labs Rules!", scanner.Scanner.Text())
}

func TestNewMockScanner(t *testing.T) {
	var scanner = NewMockScanner()
	assert.NotNil(t, scanner)
}

func TestMockScanner_Scan(t *testing.T) {
	var scanner = NewMockScanner()
	scanner.WriteString("CSX Labs Rules!")
	assert.True(t, scanner.Scan())
	assert.Equal(t, "CSX Labs Rules!", scanner.Scanner.Text())
}
