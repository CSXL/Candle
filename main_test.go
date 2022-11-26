package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Function to capture standard output from the main function with the given os.Args
func CaptureStandardOutput(os_args ...string) []byte {
	os.Args = os_args
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
	var out = CaptureStandardOutput("candle", "get-info", "IBM")
	assert.NotEmpty(t, out)
}

// Test the about command
func TestAbout(t *testing.T) {
	var out = CaptureStandardOutput("candle", "about")
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
	os.Args = []string{"candle", "exit"}
	main()
	assert.Equal(t, 0, exit_code, "Exit code should be 0, but was %d", exit_code)
}

// Test an invalid command
func TestInvalidCommand(t *testing.T) {
	var out = CaptureStandardOutput("candle", "veryinvalidcommand")
	assert.Equal(t, "--> Invalid command: veryinvalidcommand\n", string(out))
}
