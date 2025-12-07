package aiparser

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/geminal/skube/internal/cluster"
	"github.com/geminal/skube/internal/config"
	"github.com/geminal/skube/internal/parser"
	"github.com/ollama/ollama/api"
)

type AIParser struct {
	manager      *ModelManager
	openAIClient *OpenAIClient
	provider     string
}

func NewAIParser() (*AIParser, error) {
	cfg, err := config.LoadAIConfig()
	if err != nil {
		return nil, err
	}

	// Default to ollama if provider not specified
	provider := cfg.Provider
	if provider == "" {
		if cfg.OpenAIAPIKey != "" {
			provider = "openai"
		} else {
			provider = "ollama"
		}
	}

	parser := &AIParser{
		provider: provider,
	}

	if provider == "openai" {
		if cfg.OpenAIAPIKey == "" {
			return nil, errors.New("OpenAI API key not configured. Please run 'skube setup-ai' or add 'openai_api_key' to your config")
		}
		model := cfg.OpenAIModel
		if model == "" {
			model = "gpt-4o-mini"
		}
		parser.openAIClient = NewOpenAIClient(cfg.OpenAIAPIKey, model)
	} else {
		manager, err := NewModelManager()
		if err != nil {
			return nil, err
		}

		if err := manager.EnsureModelAvailable(); err != nil {
			return nil, err
		}

		parser.manager = manager
	}

	return parser, nil
}

func (p *AIParser) Parse(args []string) (*parser.Context, error) {
	userInput := strings.Join(args, " ")

	// Fetch cluster resources for context (with reasonable timeout)
	// Increased to 5s to get more complete data for better matching
	resources, _ := cluster.GetCommonResourceNames(5 * time.Second)

	// Flatten resources into a single list for the prompt
	var resourceList []string
	if resources != nil {
		resourceList = append(resourceList, resources.Namespaces...)
		resourceList = append(resourceList, resources.Deployments...)
		resourceList = append(resourceList, resources.StatefulSets...)
		resourceList = append(resourceList, resources.DaemonSets...)
		resourceList = append(resourceList, resources.Services...)
	}

	prompt := FormatPrompt(userInput, resourceList)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var response string
	var err error

	if p.provider == "openai" {
		response, err = p.openAIClient.Generate(ctx, prompt)
		if err != nil {
			return nil, err
		}
	} else {
		// Use Ollama
		req := &api.GenerateRequest{
			Model:  p.manager.GetModelName(),
			Prompt: prompt,
			Stream: new(bool), // Disable streaming
			Options: map[string]interface{}{
				"temperature": 0.0,  // Lower = more deterministic (0.0 to 1.0)
				"top_p":       0.95, // Higher = more diverse vocabulary
				"top_k":       40,   // Limit token choices for consistency
				"num_predict": 200,  // Limit response length to just the JSON
			},
		}

		respFunc := func(resp api.GenerateResponse) error {
			response += resp.Response
			return nil
		}

		err = p.manager.GetClient().Generate(ctx, req, respFunc)
		if err != nil {
			return nil, err
		}
	}

	response = strings.TrimSpace(response)
	response = extractJSON(response)

	parsedCtx, err := parseJSONToContext(response)
	if err != nil {
		return nil, err
	}

	return parsedCtx, nil
}

func extractJSON(text string) string {
	text = strings.TrimSpace(text)

	if start := strings.Index(text, "{"); start != -1 {
		if end := strings.LastIndex(text, "}"); end != -1 && end > start {
			return text[start : end+1]
		}
	}

	return text
}

func parseJSONToContext(jsonStr string) (*parser.Context, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return nil, errors.New("failed to parse AI response as JSON")
	}

	ctx := &parser.Context{}

	if v, ok := raw["command"].(string); ok {
		ctx.Command = v
	}
	if v, ok := raw["namespace"].(string); ok {
		ctx.Namespace = v
	}
	if v, ok := raw["appName"].(string); ok {
		ctx.AppName = v
	}
	if v, ok := raw["podName"].(string); ok {
		ctx.PodName = v
	}
	if v, ok := raw["serviceName"].(string); ok {
		ctx.ServiceName = v
	}
	if v, ok := raw["deploymentName"].(string); ok {
		ctx.DeploymentName = v
	}
	if v, ok := raw["resourceType"].(string); ok {
		ctx.ResourceType = v
	}
	if v, ok := raw["resourceName"].(string); ok {
		ctx.ResourceName = v
	}
	if v, ok := raw["port"].(string); ok {
		ctx.Port = v
	}
	if v, ok := raw["replicas"].(string); ok {
		ctx.Replicas = v
	}
	if v, ok := raw["follow"].(bool); ok {
		ctx.Follow = v
	}
	if v, ok := raw["prefix"].(bool); ok {
		ctx.Prefix = v
	}
	if v, ok := raw["searchTerm"].(string); ok {
		ctx.SearchTerm = v
	}
	if v, ok := raw["tailLines"].(float64); ok {
		ctx.TailLines = int(v)
	}
	if v, ok := raw["filePath"].(string); ok {
		ctx.FilePath = v
	}
	if v, ok := raw["sourcePath"].(string); ok {
		ctx.SourcePath = v
	}
	if v, ok := raw["destPath"].(string); ok {
		ctx.DestPath = v
	}

	return ctx, nil
}
