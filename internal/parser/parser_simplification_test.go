package parser

import (
	"testing"
)

func TestPrepositionlessParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Context
	}{
		{
			name: "pods qa (implicit namespace)",
			args: []string{"pods", "qa"},
			expected: Context{
				Command:   "pods",
				Namespace: "qa", // Currently this will likely fail/be empty
			},
		},
		{
			name: "logs myapp (implicit pod)",
			args: []string{"logs", "myapp"},
			expected: Context{
				Command: "logs",
				PodName: "myapp",
			},
		},
		{
			name: "logs app myapp (implicit app)",
			args: []string{"logs", "app", "myapp"},
			expected: Context{
				Command: "logs",
				AppName: "myapp",
			},
		},
		{
			name: "logs myapp qa (implicit pod and namespace)",
			args: []string{"logs", "myapp", "qa"},
			expected: Context{
				Command:   "logs",
				PodName:   "myapp",
				Namespace: "qa",
			},
		},
		{
			name: "logs app myapp qa (implicit app and namespace)",
			args: []string{"logs", "app", "myapp", "qa"},
			expected: Context{
				Command:   "logs",
				AppName:   "myapp",
				Namespace: "qa",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ParseNaturalLanguage(tt.args)

			if ctx.Command != tt.expected.Command {
				t.Errorf("expected command %s, got %s", tt.expected.Command, ctx.Command)
			}
			if ctx.Namespace != tt.expected.Namespace {
				t.Errorf("expected namespace %s, got %s", tt.expected.Namespace, ctx.Namespace)
			}
			if ctx.AppName != tt.expected.AppName {
				t.Errorf("expected app name %s, got %s", tt.expected.AppName, ctx.AppName)
			}
		})
	}
}
