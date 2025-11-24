package parser

import (
	"testing"
)

func TestNaturalLanguageParsing(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Context
	}{
		{
			name: "skube please get the pods",
			args: []string{"please", "get", "the", "pods"},
			expected: Context{
				Command: "pods",
			},
		},
		{
			name: "skube show me logs for myapp",
			args: []string{"show", "me", "logs", "for", "myapp"},
			expected: Context{
				Command: "logs",
				PodName: "myapp", // "myapp" is treated as pod name by default
			},
		},
		{
			name: "skube list all deployments in qa",
			args: []string{"list", "all", "deployments", "in", "qa"},
			expected: Context{
				Command:   "deployments",
				Namespace: "qa",
			},
		},
		{
			name: "skube give me the status",
			args: []string{"give", "me", "the", "status"},
			expected: Context{
				Command: "status",
			},
		},
		{
			name: "skube check usage for pods",
			args: []string{"check", "usage", "for", "pods"},
			expected: Context{
				Command:      "metrics",
				ResourceType: "pods",
			},
		},
		{
			name: "skube show me the config",
			args: []string{"show", "me", "the", "config"},
			expected: Context{
				Command:      "config",
				ResourceType: "view",
			},
		},
		{
			name: "skube restart the backend deployment in staging",
			args: []string{"restart", "the", "backend", "deployment", "in", "staging"},
			expected: Context{
				Command:        "restart",
				DeploymentName: "backend",
				Namespace:      "staging",
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
			if ctx.PodName != tt.expected.PodName {
				t.Errorf("expected pod name %s, got %s", tt.expected.PodName, ctx.PodName)
			}
			if ctx.DeploymentName != tt.expected.DeploymentName {
				t.Errorf("expected deployment name %s, got %s", tt.expected.DeploymentName, ctx.DeploymentName)
			}
		})
	}
}
