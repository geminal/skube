package executor

import (
	"os"
	"os/exec"
	"testing"

	"github.com/geminal/skube/internal/parser"
)

// Mocking exec.Command for testing
var execCommand = exec.Command

func TestExecuteCommand_Empty(t *testing.T) {
	ctx := &parser.Context{}
	err := ExecuteCommand(ctx)
	if err != nil {
		t.Errorf("ExecuteCommand(empty) returned error: %v", err)
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"", ""},
		{"-flag", "./-flag"}, // Prepended with ./ to prevent flag injection
		{"normal-file.txt", "normal-file.txt"},
		{"path/to/file", "path/to/file"},
		{"; rm -rf /", ""},      // Dangerous characters blocked
		{"$(whoami)", ""},       // Command substitution blocked
		{"test&background", ""}, // Background operator blocked
		{"test|pipe", ""},       // Pipe blocked
	}

	for _, tt := range tests {
		got := sanitizeInput(tt.input)
		if got != tt.expected {
			t.Errorf("sanitizeInput(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

// Helper to mock exec.Command
func fakeExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	os.Exit(0)
}
