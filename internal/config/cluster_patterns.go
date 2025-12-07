package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ClusterPatterns holds learned patterns from the Kubernetes cluster
type ClusterPatterns struct {
	KubeContext        string            `json:"kubeContext"`      // Kubernetes context this cache is for
	ClusterName        string            `json:"clusterName"`      // Cluster name (optional, for display)
	LastUpdated        time.Time         `json:"lastUpdated"`
	Namespaces         []string          `json:"namespaces"`
	CommonApps         []string          `json:"commonApps"`
	Deployments        []string          `json:"deployments"`
	Services           []string          `json:"services"`
	Pods               []string          `json:"pods"`
	Patterns           []string          `json:"patterns"`
	MultiWordResources []string          `json:"multiWordResources"`
	AppLabels          map[string]string `json:"appLabels"`        // pod name -> app label
	NamingConvention   string            `json:"namingConvention"` // detected naming style: "hyphen", "camelCase", "underscore", "PascalCase", "mixed"
}

const (
	patternsSubDir   = "patterns"
	patternsCacheTTL = 24 * time.Hour
)

// GetCurrentKubeContext returns the current kubectl context
func GetCurrentKubeContext() (string, error) {
	cmd := exec.Command("kubectl", "config", "current-context")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current kubectl context: %w", err)
	}

	context := strings.TrimSpace(string(output))
	if context == "" {
		return "", fmt.Errorf("no kubectl context is currently set")
	}

	return context, nil
}

// GetCurrentClusterName returns the cluster name for the current context (optional)
func GetCurrentClusterName() string {
	cmd := exec.Command("kubectl", "config", "view", "--minify", "-o", "jsonpath={.clusters[0].name}")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// sanitizeContextName converts a kubectl context name to a safe filename
func sanitizeContextName(context string) string {
	// Replace characters that are problematic in filenames
	safe := strings.ReplaceAll(context, "/", "_")
	safe = strings.ReplaceAll(safe, ":", "_")
	safe = strings.ReplaceAll(safe, "\\", "_")
	safe = strings.ReplaceAll(safe, " ", "_")
	safe = strings.ReplaceAll(safe, "*", "_")
	safe = strings.ReplaceAll(safe, "?", "_")
	safe = strings.ReplaceAll(safe, "\"", "_")
	safe = strings.ReplaceAll(safe, "<", "_")
	safe = strings.ReplaceAll(safe, ">", "_")
	safe = strings.ReplaceAll(safe, "|", "_")
	return safe
}

// getConfigDir returns the skube configuration directory path
func getConfigDir() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "darwin", "linux":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(home, ".config", "skube")
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			configDir = filepath.Join(os.TempDir(), "skube")
		} else {
			configDir = filepath.Join(localAppData, "skube")
		}
	default:
		configDir = filepath.Join(os.TempDir(), "skube")
	}

	return configDir, nil
}

// LoadClusterPatterns loads the cluster patterns cache from disk for the current context
func LoadClusterPatterns() (*ClusterPatterns, error) {
	currentContext, err := GetCurrentKubeContext()
	if err != nil {
		// If we can't get the current context, return empty patterns
		return &ClusterPatterns{
			Namespaces:         []string{},
			CommonApps:         []string{},
			Deployments:        []string{},
			Services:           []string{},
			Pods:               []string{},
			Patterns:           []string{},
			MultiWordResources: []string{},
			AppLabels:          make(map[string]string),
		}, nil
	}

	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	// Use context-specific file path
	patternsDir := filepath.Join(configDir, patternsSubDir)
	safeContext := sanitizeContextName(currentContext)
	filePath := filepath.Join(patternsDir, safeContext+".json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty patterns if file doesn't exist (first run for this context)
			return &ClusterPatterns{
				KubeContext:        currentContext,
				Namespaces:         []string{},
				CommonApps:         []string{},
				Deployments:        []string{},
				Services:           []string{},
				Pods:               []string{},
				Patterns:           []string{},
				MultiWordResources: []string{},
				AppLabels:          make(map[string]string),
			}, nil
		}
		return nil, err
	}

	var patterns ClusterPatterns
	if err := json.Unmarshal(data, &patterns); err != nil {
		return nil, err
	}

	// Verify the loaded patterns match the current context
	if patterns.KubeContext != currentContext {
		// Context mismatch - return empty patterns and let it refresh
		return &ClusterPatterns{
			KubeContext:        currentContext,
			Namespaces:         []string{},
			CommonApps:         []string{},
			Deployments:        []string{},
			Services:           []string{},
			Pods:               []string{},
			Patterns:           []string{},
			MultiWordResources: []string{},
			AppLabels:          make(map[string]string),
		}, nil
	}

	return &patterns, nil
}

// SaveClusterPatterns saves the cluster patterns cache to disk
func SaveClusterPatterns(patterns *ClusterPatterns) error {
	// Ensure context is set
	if patterns.KubeContext == "" {
		currentContext, err := GetCurrentKubeContext()
		if err != nil {
			return fmt.Errorf("cannot save patterns: %w", err)
		}
		patterns.KubeContext = currentContext
	}

	// Set cluster name if not already set
	if patterns.ClusterName == "" {
		patterns.ClusterName = GetCurrentClusterName()
	}

	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	// Create patterns directory if it doesn't exist
	patternsDir := filepath.Join(configDir, patternsSubDir)
	if err := os.MkdirAll(patternsDir, 0755); err != nil {
		return err
	}

	patterns.LastUpdated = time.Now()

	data, err := json.MarshalIndent(patterns, "", "  ")
	if err != nil {
		return err
	}

	// Save to context-specific file
	safeContext := sanitizeContextName(patterns.KubeContext)
	filePath := filepath.Join(patternsDir, safeContext+".json")
	return os.WriteFile(filePath, data, 0644)
}

// IsClusterPatternsCacheStale checks if the patterns cache needs refresh
func IsClusterPatternsCacheStale() bool {
	patterns, err := LoadClusterPatterns()
	if err != nil {
		// If we can't load, consider it stale
		return true
	}

	// Check if cache is empty (first run)
	if len(patterns.Namespaces) == 0 && len(patterns.Deployments) == 0 {
		return true
	}

	// Check if cache has expired
	return time.Since(patterns.LastUpdated) > patternsCacheTTL
}

// GetClusterPatternsPath returns the full path to the patterns cache file for the current context
func GetClusterPatternsPath() (string, error) {
	currentContext, err := GetCurrentKubeContext()
	if err != nil {
		return "", err
	}

	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	patternsDir := filepath.Join(configDir, patternsSubDir)
	safeContext := sanitizeContextName(currentContext)
	return filepath.Join(patternsDir, safeContext+".json"), nil
}

// DeleteClusterPatterns removes the patterns cache file
func DeleteClusterPatterns() error {
	filePath, err := GetClusterPatternsPath()
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}
