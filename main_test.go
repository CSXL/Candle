package main

import (
	"io"
	"os"
	"os/exec"
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
	os.Args = []string{"candle", "exit"}
	if os.Getenv("BE_CRASHER") == "1" {
		main()
		return
	}
	cmd := exec.Command(original_arguments[0], "-test.run=TestCrasher", "-args", os.Args[0], os.Args[1])
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("Process ran with err %v, want exit status 0", err)
}

// Test an invalid command
func TestInvalidCommand(t *testing.T) {
	os.Args = []string{"candle", "veryinvalidcommand"}
	var out = CaptureStandardOutput()
	assert.Equal(t, "--> Invalid command: veryinvalidcommand\n", string(out))
}
