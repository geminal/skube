package parser

import (
	"reflect"
	"testing"
)

func TestParseNaturalLanguage(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected *Context
	}{
		{
			name: "Basic get pods",
			args: []string{"get", "pods"},
			expected: &Context{
				Command: "pods",
			},
		},
		{
			name: "Get pods in namespace",
			args: []string{"get", "pods", "in", "qa"},
			expected: &Context{
				Command:   "pods",
				Namespace: "qa",
			},
		},
		{
			name: "Namespace first syntax",
			args: []string{"in", "qa", "get", "pods"},
			expected: &Context{
				Command:   "pods",
				Namespace: "qa",
			},
		},
		{
			name: "Logs from app",
			args: []string{"logs", "from", "app", "myapp", "in", "qa"},
			expected: &Context{
				Command:   "logs",
				AppName:   "myapp",
				Namespace: "qa",
			},
		},
		{
			name: "Logs from app (namespace first)",
			args: []string{"in", "qa", "logs", "from", "app", "myapp"},
			expected: &Context{
				Command:   "logs",
				AppName:   "myapp",
				Namespace: "qa",
			},
		},
		{
			name: "Logs of pod",
			args: []string{"logs", "of", "pod", "mypod"},
			expected: &Context{
				Command: "logs",
				PodName: "mypod",
			},
		},
		{
			name: "Logs of pod in namespace",
			args: []string{"in", "qa", "logs", "of", "pod", "mypod"},
			expected: &Context{
				Command:   "logs",
				PodName:   "mypod",
				Namespace: "qa",
			},
		},
		{
			name: "Scale deployment",
			args: []string{"scale", "deployment", "backend", "to", "3", "in", "prod"},
			expected: &Context{
				Command:        "scale",
				DeploymentName: "backend",
				Replicas:       "3",
				Namespace:      "prod",
			},
		},
		{
			name: "Port forward",
			args: []string{"forward", "service", "web", "port", "8080", "in", "dev"},
			expected: &Context{
				Command:     "forward",
				ServiceName: "web",
				Port:        "8080",
				Namespace:   "dev",
			},
		},
		{
			name: "Delete pod",
			args: []string{"delete", "pod", "mypod", "in", "qa"},
			expected: &Context{
				Command:      "delete",
				ResourceType: "pod",
				ResourceName: "mypod",
				Namespace:    "qa",
			},
		},
		{
			name: "Show metrics",
			args: []string{"show", "metrics", "pods", "in", "qa"},
			expected: &Context{
				Command:      "metrics",
				ResourceType: "pods",
				Namespace:    "qa",
			},
		},
		{
			name: "Apply file",
			args: []string{"apply", "file", "deploy.yaml"},
			expected: &Context{
				Command:  "apply",
				FilePath: "deploy.yaml",
			},
		},
		{
			name: "Copy file",
			args: []string{"copy", "file", "local.txt", "to", "/tmp/remote.txt", "in", "qa"},
			expected: &Context{
				Command:    "copy",
				SourcePath: "local.txt",
				DestPath:   "/tmp/remote.txt",
				Namespace:  "qa",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseNaturalLanguage(tt.args)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseNaturalLanguage() = %+v, want %+v", got, tt.expected)
			}
		})
	}
}
