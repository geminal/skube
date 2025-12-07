package aiparser

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	OllamaDockerImage     = "ollama/ollama:latest"
	OllamaContainerName   = "skube-ollama"
	OllamaDockerPort      = "11434"
	OllamaDockerEndpoint  = "http://localhost:11434"
)

type DockerManager struct{}

func NewDockerManager() *DockerManager {
	return &DockerManager{}
}

func (d *DockerManager) IsDockerAvailable() bool {
	cmd := exec.Command("docker", "--version")
	return cmd.Run() == nil
}

func (d *DockerManager) IsOllamaContainerRunning() bool {
	cmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", OllamaContainerName), "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == OllamaContainerName
}

func (d *DockerManager) StartOllamaContainer() error {
	if d.IsOllamaContainerRunning() {
		return nil
	}

	// Check if container exists but is stopped
	checkCmd := exec.Command("docker", "ps", "-a", "--filter", fmt.Sprintf("name=%s", OllamaContainerName), "--format", "{{.Names}}")
	output, _ := checkCmd.Output()
	containerExists := strings.TrimSpace(string(output)) == OllamaContainerName

	if containerExists {
		// Start existing container
		fmt.Println("Starting Ollama container...")
		cmd := exec.Command("docker", "start", OllamaContainerName)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to start Ollama container: %w", err)
		}
	} else {
		// Check if image exists locally
		if !d.hasOllamaImage() {
			fmt.Println("Downloading Ollama image (~700MB)...")
			fmt.Println("This is a one-time download and may take a few minutes.")
			if err := d.pullOllamaImage(); err != nil {
				return fmt.Errorf("failed to pull Ollama image: %w", err)
			}
			fmt.Println("✓ Image downloaded successfully!")
		}

		// Create and start new container
		fmt.Println("Creating Ollama container...")

		cmd := exec.Command("docker", "run", "-d",
			"--name", OllamaContainerName,
			"-p", fmt.Sprintf("%s:11434", OllamaDockerPort),
			OllamaDockerImage,
		)

		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to create Ollama container: %w\nOutput: %s", err, string(output))
		}
		fmt.Printf("✓ Container created!\n")
	}

	// Wait for container to be ready
	fmt.Print("Waiting for Ollama to be ready")
	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print(".")

		cmd := exec.Command("docker", "exec", OllamaContainerName, "ollama", "list")
		if cmd.Run() == nil {
			fmt.Println(" ✓")
			return nil
		}
	}

	return fmt.Errorf("timeout waiting for Ollama container to be ready")
}

func (d *DockerManager) PullModelInContainer(modelName string) error {
	if !d.IsOllamaContainerRunning() {
		return fmt.Errorf("Ollama container is not running")
	}

	fmt.Printf("Pulling model '%s' in Docker container...\n", modelName)

	cmd := exec.Command("docker", "exec", OllamaContainerName, "ollama", "pull", modelName)
	cmd.Stdout = nil // We'll show our own progress
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull model: %w", err)
	}

	return nil
}

func (d *DockerManager) IsModelPulledInContainer(modelName string) (bool, error) {
	if !d.IsOllamaContainerRunning() {
		return false, nil
	}

	cmd := exec.Command("docker", "exec", OllamaContainerName, "ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return false, err
	}

	return strings.Contains(string(output), modelName), nil
}

func (d *DockerManager) StopOllamaContainer() error {
	if !d.IsOllamaContainerRunning() {
		return nil
	}

	cmd := exec.Command("docker", "stop", OllamaContainerName)
	return cmd.Run()
}

func (d *DockerManager) GetDockerEndpoint() string {
	return OllamaDockerEndpoint
}

func (d *DockerManager) hasOllamaImage() bool {
	cmd := exec.Command("docker", "images", "-q", OllamaDockerImage)
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(strings.TrimSpace(string(output))) > 0
}

func (d *DockerManager) pullOllamaImage() error {
	cmd := exec.Command("docker", "pull", OllamaDockerImage)
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Show progress by running command and checking image size periodically
	done := make(chan error)
	go func() {
		done <- cmd.Run()
	}()

	// Show progress dots while pulling
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case err := <-done:
			fmt.Println() // New line after progress dots
			return err
		case <-ticker.C:
			fmt.Print(".")
		}
	}
}
