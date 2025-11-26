package executor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/geminal/skube/internal/parser"
)

// Mocking exec.Command for testing
// var execCommand = exec.Command // Removed: defined in executor.go

func TestExecuteCommand_Empty(t *testing.T) {
	ctx := &parser.Context{}
	// Capture stdout to avoid cluttering test output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := ExecuteCommand(ctx)

	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	if err != nil {
		t.Errorf("ExecuteCommand(empty) returned error: %v", err)
	}
	// Should print help
	if !strings.Contains(buf.String(), "USAGE:") {
		t.Errorf("Expected help output, got: %s", buf.String())
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
	// Print the command and arguments to stdout so we can verify them
	args := os.Args[3:]
	fmt.Printf("MOCK_EXEC: %s\n", strings.Join(args, " "))
	os.Exit(0)
}

func TestHandleLogs(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	tests := []struct {
		name     string
		ctx      *parser.Context
		expected string
	}{
		{
			name: "Logs for Pod",
			ctx: &parser.Context{
				Command: "logs",
				PodName: "mypod",
			},
			expected: "MOCK_EXEC: kubectl logs mypod",
		},
		{
			name: "Logs for App in Namespace",
			ctx: &parser.Context{
				Command:   "logs",
				AppName:   "myapp",
				Namespace: "qa",
			},
			expected: "MOCK_EXEC: kubectl logs -l app=myapp -n qa",
		},
		{
			name: "Logs with Follow and Tail",
			ctx: &parser.Context{
				Command:   "logs",
				PodName:   "mypod",
				Follow:    true,
				TailLines: 100,
			},
			expected: "MOCK_EXEC: kubectl logs mypod -f --tail=100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := ExecuteCommand(tt.ctx)

			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if err != nil {
				t.Errorf("ExecuteCommand returned error: %v", err)
			}
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output containing %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestHandlePods(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	tests := []struct {
		name     string
		ctx      *parser.Context
		expected string
	}{
		{
			name: "Get Pods",
			ctx: &parser.Context{
				Command: "pods",
			},
			expected: "MOCK_EXEC: kubectl get pods -o wide",
		},
		{
			name: "Get Pods in Namespace",
			ctx: &parser.Context{
				Command:   "pods",
				Namespace: "dev",
			},
			expected: "MOCK_EXEC: kubectl get pods -o wide -n dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := ExecuteCommand(tt.ctx)

			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if err != nil {
				t.Errorf("ExecuteCommand returned error: %v", err)
			}
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output containing %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestHandleScale(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	tests := []struct {
		name     string
		ctx      *parser.Context
		expected string
	}{
		{
			name: "Scale Deployment",
			ctx: &parser.Context{
				Command:        "scale",
				DeploymentName: "web",
				Replicas:       "3",
			},
			expected: "MOCK_EXEC: kubectl scale deployment web --replicas=3",
		},
		{
			name: "Scale Deployment in Namespace",
			ctx: &parser.Context{
				Command:        "scale",
				DeploymentName: "web",
				Replicas:       "5",
				Namespace:      "prod",
			},
			expected: "MOCK_EXEC: kubectl scale deployment web --replicas=5 -n prod",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := ExecuteCommand(tt.ctx)

			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if err != nil {
				t.Errorf("ExecuteCommand returned error: %v", err)
			}
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output containing %q, got %q", tt.expected, output)
			}
		})
	}
}

func TestHandleDelete(t *testing.T) {
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.Command }()

	tests := []struct {
		name     string
		ctx      *parser.Context
		expected string
	}{
		{
			name: "Delete Pod",
			ctx: &parser.Context{
				Command:      "delete",
				ResourceType: "pod",
				ResourceName: "mypod",
			},
			expected: "MOCK_EXEC: kubectl delete pod mypod",
		},
		{
			name: "Delete Deployment",
			ctx: &parser.Context{
				Command:      "delete",
				ResourceType: "deployment",
				ResourceName: "web",
			},
			expected: "MOCK_EXEC: kubectl delete deployment web",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := ExecuteCommand(tt.ctx)

			w.Close()
			os.Stdout = oldStdout
			var buf bytes.Buffer
			io.Copy(&buf, r)
			output := buf.String()

			if err != nil {
				t.Errorf("ExecuteCommand returned error: %v", err)
			}
			if !strings.Contains(output, tt.expected) {
				t.Errorf("Expected output containing %q, got %q", tt.expected, output)
			}
		})
	}
}
