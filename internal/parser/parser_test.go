package parser

import (
	"testing"
)

func TestParseNaturalLanguage(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Context
	}{
		{
			name: "get pods in default namespace",
			args: []string{"get", "pods"},
			expected: Context{
				Command: "pods",
			},
		},
		{
			name: "get pods in specific namespace",
			args: []string{"get", "pods", "in", "production"},
			expected: Context{
				Command:   "pods",
				Namespace: "production",
			},
		},
		{
			name: "logs of app in namespace",
			args: []string{"logs", "of", "myapp", "in", "staging"},
			expected: Context{
				Command:   "logs",
				AppName:   "myapp",
				Namespace: "staging",
			},
		},
		{
			name: "logs with follow and prefix",
			args: []string{"logs", "of", "myapp", "follow", "with", "prefix"},
			expected: Context{
				Command: "logs",
				AppName: "myapp",
				Follow:  true,
				Prefix:  true,
			},
		},
		{
			name: "scale deployment",
			args: []string{"scale", "deployment", "backend", "to", "5", "in", "prod"},
			expected: Context{
				Command:        "scale",
				DeploymentName: "backend",
				Replicas:       "5",
				Namespace:      "prod",
			},
		},
		{
			name: "restart pod",
			args: []string{"restart", "pod", "mypod", "in", "dev"},
			expected: Context{
				Command:   "restart",
				PodName:   "mypod",
				Namespace: "dev",
			},
		},
		{
			name: "port forward service",
			args: []string{"forward", "service", "web", "port", "8080", "in", "qa"},
			expected: Context{
				Command:     "forward",
				ServiceName: "web",
				Port:        "8080",
				Namespace:   "qa",
			},
		},
		{
			name: "search logs",
			args: []string{"logs", "of", "api", "search", "error"},
			expected: Context{
				Command:    "logs",
				AppName:    "api",
				SearchTerm: "error",
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
			if ctx.ServiceName != tt.expected.ServiceName {
				t.Errorf("expected service name %s, got %s", tt.expected.ServiceName, ctx.ServiceName)
			}
			if ctx.Replicas != tt.expected.Replicas {
				t.Errorf("expected replicas %s, got %s", tt.expected.Replicas, ctx.Replicas)
			}
			if ctx.Port != tt.expected.Port {
				t.Errorf("expected port %s, got %s", tt.expected.Port, ctx.Port)
			}
			if ctx.Follow != tt.expected.Follow {
				t.Errorf("expected follow %v, got %v", tt.expected.Follow, ctx.Follow)
			}
			if ctx.Prefix != tt.expected.Prefix {
				t.Errorf("expected prefix %v, got %v", tt.expected.Prefix, ctx.Prefix)
			}
			if ctx.SearchTerm != tt.expected.SearchTerm {
				t.Errorf("expected search term %s, got %s", tt.expected.SearchTerm, ctx.SearchTerm)
			}
		})
	}
}
