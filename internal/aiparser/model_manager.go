package aiparser

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/geminal/skube/internal/config"
	"github.com/ollama/ollama/api"
)

const (
	DefaultModel   = "qwen2.5-coder:3b"
	OllamaEndpoint = "http://localhost:11434"
)

type ModelManager struct {
	client        *api.Client
	dockerManager *DockerManager
	useDocker     bool
	modelName     string
}

func NewModelManager() (*ModelManager, error) {
	dockerManager := NewDockerManager()
	useDocker := false
	modelName := DefaultModel

	// Load config to check user preference
	cfg, err := config.LoadAIConfig()
	if err == nil && cfg.Enabled {
		useDocker = cfg.UseDocker
		if cfg.Model != "" {
			modelName = cfg.Model
		}
	} else if dockerManager.IsDockerAvailable() {
		// Fallback: Try Docker if available (backward compatibility)
		useDocker = true
	}

	var client *api.Client
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	if useDocker {
		// Use Docker endpoint
		endpoint, _ := url.Parse(dockerManager.GetDockerEndpoint())
		client = api.NewClient(endpoint, httpClient)
	} else {
		// Use local Ollama
		client, err = api.ClientFromEnvironment()
		if err != nil {
			endpoint, _ := url.Parse(OllamaEndpoint)
			client = api.NewClient(endpoint, httpClient)
		}
	}

	return &ModelManager{
		client:        client,
		dockerManager: dockerManager,
		useDocker:     useDocker,
		modelName:     modelName,
	}, nil
}

func (m *ModelManager) GetModelName() string {
	return m.modelName
}

func (m *ModelManager) IsOllamaRunning() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := m.client.Heartbeat(ctx)
	return err == nil
}

func (m *ModelManager) IsModelAvailable(modelName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	listResp, err := m.client.List(ctx)
	if err != nil {
		return false, err
	}

	for _, model := range listResp.Models {
		if model.Name == modelName {
			return true, nil
		}
	}

	return false, nil
}

func (m *ModelManager) PullModel(modelName string) error {
	fmt.Printf("\nPulling AI model '%s' from Ollama...\n", modelName)
	fmt.Println("This will take a few minutes depending on your connection speed.")

	ctx := context.Background()

	req := &api.PullRequest{
		Model: modelName,
	}

	progressFunc := func(resp api.ProgressResponse) error {
		if resp.Status != "" {
			if resp.Total > 0 {
				percent := float64(resp.Completed) / float64(resp.Total) * 100
				fmt.Printf("\r%s: %.1f%% (%s / %s)",
					resp.Status,
					percent,
					formatBytes(resp.Completed),
					formatBytes(resp.Total))
			} else {
				fmt.Printf("\r%s", resp.Status)
			}
		}
		return nil
	}

	err := m.client.Pull(ctx, req, progressFunc)
	if err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}

	fmt.Println("\n✓ Model pulled successfully!")
	return nil
}

func (m *ModelManager) EnsureModelAvailable() error {
	// If using Docker, start container if needed
	if m.useDocker {
		if !m.dockerManager.IsOllamaContainerRunning() {
			if err := m.dockerManager.StartOllamaContainer(); err != nil {
				return fmt.Errorf("failed to start Ollama container: %w\n"+
					"Make sure Docker is running and you have permissions to use it", err)
			}
		}

		// Check if model is already pulled in container
		modelPulled, err := m.dockerManager.IsModelPulledInContainer(m.modelName)
		if err == nil && !modelPulled {
			fmt.Printf("\nFirst time setup: Pulling AI model '%s'...\n", m.modelName)
			if err := m.dockerManager.PullModelInContainer(m.modelName); err != nil {
				return err
			}
			fmt.Println("✓ Model ready!")
		}

		return nil
	}

	// Fallback to local Ollama
	if !m.IsOllamaRunning() {
		return errors.New("Ollama is not running. Please install and start Ollama:\n" +
			"  macOS/Linux: curl -fsSL https://ollama.com/install.sh | sh\n" +
			"  Windows: https://ollama.com/download\n" +
			"Then run: ollama serve\n\n" +
			"Or install Docker to use containerized Ollama automatically.")
	}

	available, err := m.IsModelAvailable(m.modelName)
	if err != nil {
		return fmt.Errorf("failed to check model availability: %w", err)
	}

	if !available {
		if err := m.PullModel(m.modelName); err != nil {
			return err
		}
	}

	return nil
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (m *ModelManager) GetClient() *api.Client {
	return m.client
}
