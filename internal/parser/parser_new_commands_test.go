package parser

import (
	"testing"
)

func TestParseNewCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Context
	}{
		{
			name: "apply file",
			args: []string{"apply", "file", "pod.yaml"},
			expected: Context{
				Command:  "apply",
				FilePath: "pod.yaml",
			},
		},
		{
			name: "create from file",
			args: []string{"create", "from", "file", "deployment.yaml"},
			expected: Context{
				Command:  "apply",
				FilePath: "deployment.yaml",
			},
		},
		{
			name: "delete pod",
			args: []string{"delete", "pod", "mypod", "in", "prod"},
			expected: Context{
				Command:      "delete",
				ResourceType: "pod",
				ResourceName: "mypod",
				Namespace:    "prod",
			},
		},
		{
			name: "edit deployment",
			args: []string{"edit", "deployment", "mydeploy", "in", "staging"},
			expected: Context{
				Command:      "edit",
				ResourceType: "deployment",
				ResourceName: "mydeploy",
				Namespace:    "staging",
			},
		},
		{
			name: "use context",
			args: []string{"use", "context", "minikube"},
			expected: Context{
				Command:      "config",
				ResourceType: "context",
				ResourceName: "minikube",
			},
		},
		{
			name: "use namespace",
			args: []string{"use", "namespace", "dev"},
			expected: Context{
				Command:      "config",
				ResourceType: "namespace",
				ResourceName: "dev",
			},
		},
		{
			name: "show config",
			args: []string{"show", "config"},
			expected: Context{
				Command:      "config",
				ResourceType: "view",
			},
		},
		{
			name: "show metrics pods",
			args: []string{"show", "metrics", "pods", "in", "kube-system"},
			expected: Context{
				Command:      "metrics",
				ResourceType: "pods",
				Namespace:    "kube-system",
			},
		},
		{
			name: "check usage nodes",
			args: []string{"check", "usage", "nodes"},
			expected: Context{
				Command:      "metrics",
				ResourceType: "nodes",
			},
		},
		{
			name: "copy file to pod",
			args: []string{"copy", "file", "./local.txt", "to", "/tmp/remote.txt", "in", "pod-123", "in", "qa"},
			expected: Context{
				Command:    "copy",
				SourcePath: "./local.txt",
				DestPath:   "/tmp/remote.txt",
				PodName:    "pod-123",
				Namespace:  "qa",
			},
		},
		{
			name: "explain pod",
			args: []string{"explain", "pod"},
			expected: Context{
				Command:      "explain",
				ResourceType: "pod",
			},
		},
		{
			name: "what is service",
			args: []string{"what", "is", "service"},
			expected: Context{
				Command:      "explain",
				ResourceType: "service",
			},
		},
		{
			name: "what is ingress",
			args: []string{"what", "is", "ingress"},
			expected: Context{
				Command:      "explain",
				ResourceType: "ingress",
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
			if ctx.FilePath != tt.expected.FilePath {
				t.Errorf("expected file path %s, got %s", tt.expected.FilePath, ctx.FilePath)
			}
			if ctx.SourcePath != tt.expected.SourcePath {
				t.Errorf("expected source path %s, got %s", tt.expected.SourcePath, ctx.SourcePath)
			}
			if ctx.DestPath != tt.expected.DestPath {
				t.Errorf("expected dest path %s, got %s", tt.expected.DestPath, ctx.DestPath)
			}
		})
	}
}
