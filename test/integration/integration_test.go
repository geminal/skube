package integration

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	// Skip if SKUBE_INTEGRATION is not set
	if os.Getenv("SKUBE_INTEGRATION") == "" {
		t.Skip("Skipping integration tests (SKUBE_INTEGRATION not set)")
	}

	// Build the binary
	tempDir := t.TempDir()
	binaryPath := filepath.Join(tempDir, "skube")
	cmd := exec.Command("go", "build", "-o", binaryPath, "../../cmd/skube")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build skube: %v", err)
	}

	// Helper to run skube
	runSkube := func(args ...string) (string, error) {
		cmd := exec.Command(binaryPath, args...)
		// Mock kubectl if needed, or use real one if configured
		// For this test, we assume a real cluster or a mocked environment via PATH
		// But to be safe and self-contained, we can mock kubectl by adding a script to PATH

		// Create a mock kubectl
		mockKubectl := filepath.Join(tempDir, "kubectl")
		script := `#!/bin/sh
echo "mock-kubectl: $@"
`
		if err := os.WriteFile(mockKubectl, []byte(script), 0755); err != nil {
			return "", err
		}

		// Prepend tempDir to PATH
		cmd.Env = append(os.Environ(), "PATH="+tempDir+string(os.PathListSeparator)+os.Getenv("PATH"))

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		return out.String(), err
	}

	tests := []struct {
		name    string
		args    []string
		wantOut string
		wantErr bool
	}{
		{
			name:    "Help",
			args:    []string{"help"},
			wantOut: "skube - Talk to Kubernetes in plain English",
		},
		{
			name:    "Get Pods",
			args:    []string{"get", "pods"},
			wantOut: "mock-kubectl: get pods -o wide",
		},
		{
			name:    "Logs from App",
			args:    []string{"logs", "from", "app", "myapp"},
			wantOut: "mock-kubectl: logs -l app=myapp",
		},
	}

	// Helper to strip ANSI codes
	stripANSI := func(str string) string {
		// Simple regex to strip ANSI codes
		// \x1b\[[0-9;]*m
		// Since we don't want to import regexp just for this if possible, let's use a simple string replacement loop
		// or just use regexp.
		// Let's use a simple replacement for common codes if we know them, or just ignore color in verification.
		// Better: verify the content exists regardless of color.
		// But "skube - Talk" might have color codes in between.
		// Let's just use a simple loop to remove anything starting with ESC[ and ending with m
		var ret strings.Builder
		inCode := false
		for _, r := range str {
			if r == '\x1b' {
				inCode = true
				continue
			}
			if inCode {
				if r == 'm' {
					inCode = false
				}
				continue
			}
			ret.WriteRune(r)
		}
		return ret.String()
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runSkube(tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("runSkube() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			cleanGot := stripANSI(got)
			if !strings.Contains(cleanGot, tt.wantOut) {
				t.Errorf("runSkube() got = %v (clean: %v), want substring %v", got, cleanGot, tt.wantOut)
			}
		})
	}
}
