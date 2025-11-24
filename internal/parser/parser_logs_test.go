package parser

import (
	"testing"
)

func TestParseLogsOfApp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Context
	}{
		{
			name: "logs of app myapp",
			args: []string{"logs", "of", "app", "myapp", "in", "qa"},
			expected: Context{
				Command:   "logs",
				AppName:   "myapp",
				Namespace: "qa",
			},
		},
		{
			name: "logs of myapp",
			args: []string{"logs", "of", "myapp", "in", "qa"},
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
			if ctx.AppName != tt.expected.AppName {
				t.Errorf("expected app name %s, got %s", tt.expected.AppName, ctx.AppName)
			}
			if ctx.Namespace != tt.expected.Namespace {
				t.Errorf("expected namespace %s, got %s", tt.expected.Namespace, ctx.Namespace)
			}
		})
	}
}
