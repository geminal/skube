package cluster

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/geminal/skube/internal/config"
)

// ResourceNames holds lists of common resource names found in the cluster
type ResourceNames struct {
	KubeContext  string    `json:"kubeContext"`  // Kubernetes context this cache is for
	Namespaces   []string  `json:"namespaces"`
	Deployments  []string  `json:"deployments"`
	StatefulSets []string  `json:"statefulsets"`
	DaemonSets   []string  `json:"daemonsets"`
	Services     []string  `json:"services"`
	LastUpdated  time.Time `json:"last_updated"`
}

const (
	cacheSubDir   = "resource-cache"
	cacheDuration = 10 * time.Minute
)

// GetCommonResourceNames fetches common resource names from the cluster with a timeout
// It tries to load from cache first, and if expired or missing, fetches from K8s
func GetCommonResourceNames(timeout time.Duration) (*ResourceNames, error) {
	// Get current context
	currentContext, err := config.GetCurrentKubeContext()
	if err != nil {
		// If we can't get context, return empty resources
		return &ResourceNames{
			Namespaces:   []string{},
			Deployments:  []string{},
			StatefulSets: []string{},
			DaemonSets:   []string{},
			Services:     []string{},
		}, nil
	}

	// Try to load from cache
	if cached, err := loadCache(currentContext); err == nil {
		// Verify context matches
		if cached.KubeContext == currentContext && time.Since(cached.LastUpdated) < cacheDuration {
			return cached, nil
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resources := &ResourceNames{
		KubeContext: currentContext,
		LastUpdated: time.Now(),
	}
	errChan := make(chan error, 5)

	go func() {
		ns, err := getNamespaces(ctx)
		if err == nil {
			resources.Namespaces = ns
		}
		errChan <- err
	}()

	go func() {
		deps, err := getDeployments(ctx)
		if err == nil {
			resources.Deployments = deps
		}
		errChan <- err
	}()

	go func() {
		sts, err := getStatefulSets(ctx)
		if err == nil {
			resources.StatefulSets = sts
		}
		errChan <- err
	}()

	go func() {
		ds, err := getDaemonSets(ctx)
		if err == nil {
			resources.DaemonSets = ds
		}
		errChan <- err
	}()

	go func() {
		svcs, err := getServices(ctx)
		if err == nil {
			resources.Services = svcs
		}
		errChan <- err
	}()

	// Wait for all goroutines to finish or timeout
	for i := 0; i < 5; i++ {
		select {
		case <-ctx.Done():
			// If we have partial results, return them (and maybe cache them?)
			// For now, let's just return what we have if context expires
			return resources, ctx.Err()
		case <-errChan:
			// We ignore individual errors to return partial results
			continue
		}
	}

	// Save to cache
	_ = saveCache(resources)

	return resources, nil
}

func getCachePath(context string) (string, error) {
	// Use the same getConfigDir function as cluster patterns for consistency
	configDir, err := getSkubeConfigDir()
	if err != nil {
		return "", err
	}

	// Store in resource-cache subdirectory
	cacheDir := filepath.Join(configDir, cacheSubDir)

	// Sanitize context name for filename
	safeContext := sanitizeContextName(context)
	return filepath.Join(cacheDir, safeContext+".json"), nil
}

// getSkubeConfigDir returns the skube configuration directory path
// This duplicates the logic from config package to avoid circular import
func getSkubeConfigDir() (string, error) {
	var configDir string

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Always use ~/.config/skube for consistency across all cache types
	configDir = filepath.Join(home, ".config", "skube")

	return configDir, nil
}

func loadCache(context string) (*ResourceNames, error) {
	path, err := getCachePath(context)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var resources ResourceNames
	if err := json.Unmarshal(data, &resources); err != nil {
		return nil, err
	}

	return &resources, nil
}

func saveCache(resources *ResourceNames) error {
	if resources.KubeContext == "" {
		// No context set, skip caching
		return nil
	}

	path, err := getCachePath(resources.KubeContext)
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.Marshal(resources)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
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

func getNamespaces(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Fields(string(out)), nil
}

func getDeployments(ctx context.Context) ([]string, error) {
	return getNamespacedResources(ctx, "deployments")
}

func getStatefulSets(ctx context.Context) ([]string, error) {
	return getNamespacedResources(ctx, "statefulsets")
}

func getDaemonSets(ctx context.Context) ([]string, error) {
	return getNamespacedResources(ctx, "daemonsets")
}

func getNamespacedResources(ctx context.Context, resourceType string) ([]string, error) {
	// Get resources with namespace context: namespace/resource-name
	cmd := exec.CommandContext(ctx, "kubectl", "get", resourceType, "--all-namespaces", "-o", "jsonpath={range .items[*]}{.metadata.namespace}/{.metadata.name}{\"\\n\"}{end}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	uniqueNames := make(map[string]bool)
	var result []string

	// Add both full paths and just names for flexibility
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Add full path (e.g., "namespace-a/app-name")
		if !uniqueNames[line] {
			uniqueNames[line] = true
			result = append(result, line)
		}

		// Also add just the resource name
		parts := strings.Split(line, "/")
		if len(parts) == 2 {
			name := parts[1]
			if !uniqueNames[name] {
				uniqueNames[name] = true
				result = append(result, name)
			}
		}
	}
	return result, nil
}

func getServices(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "services", "--all-namespaces", "-o", "jsonpath={.items[*].metadata.name}")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	// Deduplicate names
	names := strings.Fields(string(out))
	uniqueNames := make(map[string]bool)
	var result []string
	for _, name := range names {
		if !uniqueNames[name] {
			uniqueNames[name] = true
			result = append(result, name)
		}
	}
	return result, nil
}
