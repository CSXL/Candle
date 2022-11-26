package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Function to capture standard output and run the main function
func CaptureStandardOutput() []byte {
	r, w, _ := os.Pipe()
	os.Stdout = w
	main()
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = os.Stderr
	return out
}

var original_arguments []string

func TestMain(m *testing.M) {
	original_arguments = os.Args // Some tests need to access the original arguments to re-run the test command.
	os.Exit(m.Run())
}

// Test the get-info command
func TestGetInfo(t *testing.T) {
	os.Args = []string{"candle", "get-info", "IBM"}
	var out = CaptureStandardOutput()
	assert.NotEmpty(t, out)
}

// Test the about command
func TestAbout(t *testing.T) {
	os.Args = []string{"candle", "about"}
	var out = CaptureStandardOutput()
	assert.Equal(t, "--> Candle CLI\nCSX Labs\nMade w/ <3 by @absozero and @ecsbeats\n", string(out))
}

// Test the exit command
func TestExit(t *testing.T) {
	// Save current function and restore at the end:
	oldOsExit := osExit
	defer func() { osExit = oldOsExit }()

	var exit_code int
	mockExit := func(code int) {
		exit_code = code
	}

	osExit = mockExit
	main()
	assert.Equal(t, 0, exit_code, "Exit code should be 0, but was %d", exit_code)
}

// Test an invalid command
func TestInvalidCommand(t *testing.T) {
	os.Args = []string{"candle", "veryinvalidcommand"}
	var out = CaptureStandardOutput()
	assert.Equal(t, "--> Invalid command: veryinvalidcommand\n", string(out))
}
