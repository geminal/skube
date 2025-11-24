package parser

import (
	"testing"
)

func TestParseResources(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Context
	}{
		{
			name: "get nodes",
			args: []string{"get", "nodes"},
			expected: Context{
				Command: "nodes",
			},
		},
		{
			name: "get nodes shorthand",
			args: []string{"get", "no"},
			expected: Context{
				Command: "nodes",
			},
		},
		{
			name: "get configmaps in ns",
			args: []string{"get", "configmaps", "in", "prod"},
			expected: Context{
				Command:   "configmaps",
				Namespace: "prod",
			},
		},
		{
			name: "get cm shorthand",
			args: []string{"get", "cm", "in", "dev"},
			expected: Context{
				Command:   "configmaps",
				Namespace: "dev",
			},
		},
		{
			name: "get secrets",
			args: []string{"get", "secrets"},
			expected: Context{
				Command: "secrets",
			},
		},
		{
			name: "get ingress",
			args: []string{"get", "ingresses", "in", "staging"},
			expected: Context{
				Command:   "ingresses",
				Namespace: "staging",
			},
		},
		{
			name: "get ing shorthand",
			args: []string{"get", "ing", "in", "qa"},
			expected: Context{
				Command:   "ingresses",
				Namespace: "qa",
			},
		},
		{
			name: "get pvc",
			args: []string{"get", "pvc", "in", "default"},
			expected: Context{
				Command:   "pvcs",
				Namespace: "default",
			},
		},
		{
			name: "describe node",
			args: []string{"describe", "node", "node-1"},
			expected: Context{
				Command:      "describe",
				ResourceType: "node",
				ResourceName: "node-1",
			},
		},
		{
			name: "delete secret",
			args: []string{"delete", "secret", "my-secret", "in", "prod"},
			expected: Context{
				Command:      "delete",
				ResourceType: "secret",
				ResourceName: "my-secret",
				Namespace:    "prod",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := ParseNaturalLanguage(tt.args)

			if ctx.Command != tt.expected.Command {
				t.Errorf("expected command %s, got %s", tt.expected.Command, ctx.Command)
			}
			if ctx.ResourceType != tt.expected.ResourceType {
				t.Errorf("expected resource type %s, got %s", tt.expected.ResourceType, ctx.ResourceType)
			}
			if ctx.ResourceName != tt.expected.ResourceName {
				t.Errorf("expected resource name %s, got %s", tt.expected.ResourceName, ctx.ResourceName)
			}
			if ctx.Namespace != tt.expected.Namespace {
				t.Errorf("expected namespace %s, got %s", tt.expected.Namespace, ctx.Namespace)
			}
		})
	}
}
