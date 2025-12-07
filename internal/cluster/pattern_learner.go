package cluster

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/geminal/skube/internal/config"
)

// LearnClusterPatterns queries the cluster and learns naming patterns
func LearnClusterPatterns(showProgress bool) (*config.ClusterPatterns, error) {
	// Get current kubectl context
	currentContext, err := config.GetCurrentKubeContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubectl context: %w", err)
	}

	clusterName := config.GetCurrentClusterName()

	if showProgress {
		fmt.Printf("Analyzing cluster patterns for context: %s%s%s\n",
			config.ColorCyan, currentContext, config.ColorReset)
		if clusterName != "" && clusterName != currentContext {
			fmt.Printf("  Cluster: %s\n", clusterName)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	patterns := &config.ClusterPatterns{
		KubeContext:        currentContext,
		ClusterName:        clusterName,
		LastUpdated:        time.Now(),
		Namespaces:         []string{},
		CommonApps:         []string{},
		Deployments:        []string{},
		Services:           []string{},
		Pods:               []string{},
		Patterns:           []string{},
		MultiWordResources: []string{},
		AppLabels:          make(map[string]string),
	}

	// Fetch namespaces
	namespaces, err := getNamespacesWithContext(ctx)
	if err == nil {
		patterns.Namespaces = namespaces
	}

	// Fetch deployments from all namespaces
	deployments, err := getDeploymentsAllNamespaces(ctx)
	if err == nil {
		patterns.Deployments = deployments
		patterns.MultiWordResources = append(patterns.MultiWordResources, extractMultiWordResources(deployments)...)
	}

	// Fetch services from all namespaces
	services, err := getServicesAllNamespaces(ctx)
	if err == nil {
		patterns.Services = services
		patterns.MultiWordResources = append(patterns.MultiWordResources, extractMultiWordResources(services)...)
	}

	// Fetch pods and their app labels
	pods, appLabels, err := getPodsWithAppLabels(ctx)
	if err == nil {
		patterns.Pods = pods
		patterns.AppLabels = appLabels
		patterns.CommonApps = extractCommonApps(appLabels)
	}

	// Detect naming patterns
	patterns.Patterns = detectNamingPatterns(patterns)

	// Detect the dominant naming convention
	patterns.NamingConvention = detectNamingConvention(patterns)

	// Remove duplicates from MultiWordResources
	patterns.MultiWordResources = uniqueStrings(patterns.MultiWordResources)

	if showProgress {
		fmt.Printf("Done! Found %d namespaces, %d deployments, %d services, %d apps.\n",
			len(patterns.Namespaces),
			len(patterns.Deployments),
			len(patterns.Services),
			len(patterns.CommonApps))

		// Show detected naming convention
		if patterns.NamingConvention != "" {
			fmt.Printf("Detected naming convention: %s\n", patterns.NamingConvention)
		}

		// Show where the patterns are cached
		safeContext := sanitizeContextForDisplay(currentContext)
		fmt.Printf("Patterns cached for context '%s' in ~/.config/skube/patterns/ (not committed)\n", safeContext)
	}

	return patterns, nil
}

// getNamespacesWithContext fetches all namespaces
func getNamespacesWithContext(ctx context.Context) ([]string, error) {
	return getNamespaces(ctx)
}

// getDeploymentsAllNamespaces fetches all deployments from all namespaces
func getDeploymentsAllNamespaces(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "deployments", "--all-namespaces", "-o", "jsonpath={range .items[*]}{.metadata.namespace}/{.metadata.name}{'\\n'}{end}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var deployments []string
	for _, line := range lines {
		if line != "" {
			deployments = append(deployments, line)
		}
	}
	return deployments, nil
}

// getServicesAllNamespaces fetches all services from all namespaces
func getServicesAllNamespaces(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "services", "--all-namespaces", "-o", "jsonpath={range .items[*]}{.metadata.namespace}/{.metadata.name}{'\\n'}{end}")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var services []string
	for _, line := range lines {
		if line != "" {
			services = append(services, line)
		}
	}
	return services, nil
}

// getPodsWithAppLabels fetches all pods and their app labels
func getPodsWithAppLabels(ctx context.Context) ([]string, map[string]string, error) {
	cmd := exec.CommandContext(ctx, "kubectl", "get", "pods", "--all-namespaces", "-o", "jsonpath={range .items[*]}{.metadata.namespace}/{.metadata.name}{'|'}{.metadata.labels.app}{'\\n'}{end}")
	output, err := cmd.Output()
	if err != nil {
		return nil, nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var pods []string
	appLabels := make(map[string]string)

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) >= 1 {
			podName := parts[0]
			pods = append(pods, podName)

			// Store app label if it exists
			if len(parts) == 2 && parts[1] != "" && parts[1] != "<no value>" {
				appLabels[podName] = parts[1]
			}
		}
	}

	return pods, appLabels, nil
}

// extractMultiWordResources identifies resources with hyphens (multi-word names)
func extractMultiWordResources(resources []string) []string {
	var multiWord []string
	for _, resource := range resources {
		// Extract just the name part (after namespace/)
		parts := strings.Split(resource, "/")
		name := resource
		if len(parts) == 2 {
			name = parts[1]
		}

		// Check if it contains hyphens (indicating multi-word)
		if strings.Contains(name, "-") {
			multiWord = append(multiWord, name)
		}
	}
	return multiWord
}

// extractCommonApps extracts unique app labels
func extractCommonApps(appLabels map[string]string) []string {
	appSet := make(map[string]bool)
	for _, app := range appLabels {
		if app != "" {
			appSet[app] = true
		}
	}

	var apps []string
	for app := range appSet {
		apps = append(apps, app)
	}
	return apps
}

// detectNamingPatterns analyzes resources to detect common naming conventions
func detectNamingPatterns(patterns *config.ClusterPatterns) []string {
	detectedPatterns := make(map[string]bool)

	// Check for {app}-{namespace} pattern
	for _, deployment := range patterns.Deployments {
		parts := strings.Split(deployment, "/")
		if len(parts) != 2 {
			continue
		}
		namespace := parts[0]
		name := parts[1]

		// Check if name ends with -namespace
		if strings.HasSuffix(name, "-"+namespace) {
			detectedPatterns["{app}-{namespace}"] = true
		}

		// Check for {app}-service pattern
		if strings.HasSuffix(name, "-service") {
			detectedPatterns["{app}-service"] = true
		}

		// Check for {app}-api pattern
		if strings.HasSuffix(name, "-api") {
			detectedPatterns["{app}-api"] = true
		}

		// Check for {app}-worker pattern
		if strings.HasSuffix(name, "-worker") {
			detectedPatterns["{app}-worker"] = true
		}
	}

	var result []string
	for pattern := range detectedPatterns {
		result = append(result, pattern)
	}
	return result
}

// uniqueStrings removes duplicates from a string slice
func uniqueStrings(input []string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, str := range input {
		if !seen[str] {
			seen[str] = true
			result = append(result, str)
		}
	}

	return result
}

// detectNamingConvention analyzes resource names to determine the dominant naming style
func detectNamingConvention(patterns *config.ClusterPatterns) string {
	counts := map[string]int{
		"hyphen":     0,
		"camelCase":  0,
		"underscore": 0,
		"PascalCase": 0,
	}

	// Analyze all resource names
	allResources := []string{}
	allResources = append(allResources, patterns.CommonApps...)
	allResources = append(allResources, patterns.MultiWordResources...)

	// Extract just the names from namespace/name format
	for _, deployment := range patterns.Deployments {
		parts := strings.Split(deployment, "/")
		if len(parts) == 2 {
			allResources = append(allResources, parts[1])
		}
	}

	for _, service := range patterns.Services {
		parts := strings.Split(service, "/")
		if len(parts) == 2 {
			allResources = append(allResources, parts[1])
		}
	}

	// Count naming conventions
	for _, name := range allResources {
		convention := identifyNamingStyle(name)
		if convention != "" {
			counts[convention]++
		}
	}

	// Find the dominant convention
	maxCount := 0
	dominant := ""
	totalMultiWord := 0

	for style, count := range counts {
		totalMultiWord += count
		if count > maxCount {
			maxCount = count
			dominant = style
		}
	}

	// If we have a clear winner (>40% of multi-word resources)
	if totalMultiWord > 0 && float64(maxCount)/float64(totalMultiWord) > 0.4 {
		return dominant
	}

	// If no clear pattern or mixed conventions
	if totalMultiWord > 0 {
		return "mixed (hyphen, camelCase, underscore)"
	}

	return "hyphen (kubernetes default)"
}

// identifyNamingStyle determines the naming convention of a single resource name
func identifyNamingStyle(name string) string {
	// Skip if single word or very short
	if len(name) < 3 {
		return ""
	}

	// Check for hyphens
	if strings.Contains(name, "-") {
		return "hyphen"
	}

	// Check for underscores
	if strings.Contains(name, "_") {
		return "underscore"
	}

	// Check for camelCase or PascalCase
	hasUpper := false
	hasLower := false
	startsWithUpper := false

	if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
		startsWithUpper = true
	}

	for _, ch := range name {
		if ch >= 'A' && ch <= 'Z' {
			hasUpper = true
		}
		if ch >= 'a' && ch <= 'z' {
			hasLower = true
		}
	}

	// PascalCase: starts with uppercase and has both cases
	if startsWithUpper && hasUpper && hasLower {
		return "PascalCase"
	}

	// camelCase: starts with lowercase and has uppercase letters
	if !startsWithUpper && hasUpper && hasLower {
		return "camelCase"
	}

	return ""
}

// sanitizeContextForDisplay returns a cleaned up context name for display
func sanitizeContextForDisplay(context string) string {
	// Just return the context as-is for display purposes
	return context
}
