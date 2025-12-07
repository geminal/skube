package parser

import (
	"strings"

	"github.com/geminal/skube/internal/config"
)

// ResourceResolver helps match user input to actual cluster resources
type ResourceResolver struct {
	patterns *config.ClusterPatterns
}

// NewResourceResolver creates a new resource resolver with cluster patterns
func NewResourceResolver() *ResourceResolver {
	patterns, err := config.LoadClusterPatterns()
	if err != nil {
		// Return resolver with empty patterns if loading fails
		patterns = &config.ClusterPatterns{
			Namespaces:         []string{},
			CommonApps:         []string{},
			Deployments:        []string{},
			Services:           []string{},
			Pods:               []string{},
			Patterns:           []string{},
			MultiWordResources: []string{},
			AppLabels:          make(map[string]string),
		}
	}

	return &ResourceResolver{
		patterns: patterns,
	}
}

// ResolveAppName attempts to match user input to an actual app/deployment name
func (r *ResourceResolver) ResolveAppName(input string, namespace string) string {
	if input == "" {
		return input
	}

	// Try multiple naming conventions
	variants := generateNamingVariants(input)

	// Try exact match first for each variant
	for _, variant := range variants {
		if match := r.findExactMatch(variant, namespace); match != "" {
			return match
		}
	}

	// Try fuzzy matching on deployments with all variants
	for _, variant := range variants {
		if match := r.fuzzyMatchDeployment(variant, namespace); match != "" {
			return match
		}
	}

	// Try fuzzy matching on common apps
	for _, variant := range variants {
		if match := r.fuzzyMatchApp(variant); match != "" {
			return match
		}
	}

	// Try pattern-based matching (e.g., {app}-{namespace})
	if namespace != "" {
		for _, variant := range variants {
			if match := r.patternBasedMatch(variant, namespace); match != "" {
				return match
			}
		}
	}

	// Return the first variant (hyphen-separated) as fallback
	return variants[0]
}

// ResolveServiceName attempts to match user input to an actual service name
func (r *ResourceResolver) ResolveServiceName(input string, namespace string) string {
	if input == "" {
		return input
	}

	variants := generateNamingVariants(input)

	// Try exact match for each variant
	for _, variant := range variants {
		if match := r.findExactServiceMatch(variant, namespace); match != "" {
			return match
		}
	}

	// Try fuzzy matching
	for _, variant := range variants {
		if match := r.fuzzyMatchService(variant, namespace); match != "" {
			return match
		}
	}

	return variants[0]
}

// ResolvePodName attempts to match user input to an actual pod name
func (r *ResourceResolver) ResolvePodName(input string, namespace string) string {
	if input == "" {
		return input
	}

	variants := generateNamingVariants(input)

	// Try exact match
	for _, variant := range variants {
		if match := r.findExactPodMatch(variant, namespace); match != "" {
			return match
		}
	}

	// Try fuzzy matching
	for _, variant := range variants {
		if match := r.fuzzyMatchPod(variant, namespace); match != "" {
			return match
		}
	}

	return variants[0]
}

// ResolveNamespace attempts to match user input to an actual namespace
func (r *ResourceResolver) ResolveNamespace(input string) string {
	if input == "" {
		return input
	}

	inputLower := strings.ToLower(input)

	// Try exact match
	for _, ns := range r.patterns.Namespaces {
		if strings.ToLower(ns) == inputLower {
			return ns
		}
	}

	// Try fuzzy match
	if match, ok := FuzzyMatchWithThreshold(input, r.patterns.Namespaces); ok {
		return match
	}

	return input
}

// findExactMatch looks for exact match in deployments
func (r *ResourceResolver) findExactMatch(name string, namespace string) string {
	nameLower := strings.ToLower(name)

	for _, deployment := range r.patterns.Deployments {
		parts := strings.Split(deployment, "/")
		if len(parts) != 2 {
			continue
		}

		ns := parts[0]
		depName := parts[1]

		// If namespace is specified, only match in that namespace
		if namespace != "" && strings.ToLower(ns) != strings.ToLower(namespace) {
			continue
		}

		if strings.ToLower(depName) == nameLower {
			return depName
		}
	}

	return ""
}

// findExactServiceMatch looks for exact match in services
func (r *ResourceResolver) findExactServiceMatch(name string, namespace string) string {
	nameLower := strings.ToLower(name)

	for _, service := range r.patterns.Services {
		parts := strings.Split(service, "/")
		if len(parts) != 2 {
			continue
		}

		ns := parts[0]
		svcName := parts[1]

		if namespace != "" && strings.ToLower(ns) != strings.ToLower(namespace) {
			continue
		}

		if strings.ToLower(svcName) == nameLower {
			return svcName
		}
	}

	return ""
}

// findExactPodMatch looks for exact match in pods
func (r *ResourceResolver) findExactPodMatch(name string, namespace string) string {
	nameLower := strings.ToLower(name)

	for _, pod := range r.patterns.Pods {
		parts := strings.Split(pod, "/")
		if len(parts) != 2 {
			continue
		}

		ns := parts[0]
		podName := parts[1]

		if namespace != "" && strings.ToLower(ns) != strings.ToLower(namespace) {
			continue
		}

		if strings.ToLower(podName) == nameLower {
			return podName
		}
	}

	return ""
}

// fuzzyMatchDeployment tries fuzzy matching on deployment names
func (r *ResourceResolver) fuzzyMatchDeployment(name string, namespace string) string {
	var candidates []string

	for _, deployment := range r.patterns.Deployments {
		parts := strings.Split(deployment, "/")
		if len(parts) != 2 {
			continue
		}

		ns := parts[0]
		depName := parts[1]

		if namespace != "" && strings.ToLower(ns) != strings.ToLower(namespace) {
			continue
		}

		candidates = append(candidates, depName)
	}

	if match, ok := FuzzyMatchWithThreshold(name, candidates); ok {
		return match
	}

	return ""
}

// fuzzyMatchService tries fuzzy matching on service names
func (r *ResourceResolver) fuzzyMatchService(name string, namespace string) string {
	var candidates []string

	for _, service := range r.patterns.Services {
		parts := strings.Split(service, "/")
		if len(parts) != 2 {
			continue
		}

		ns := parts[0]
		svcName := parts[1]

		if namespace != "" && strings.ToLower(ns) != strings.ToLower(namespace) {
			continue
		}

		candidates = append(candidates, svcName)
	}

	if match, ok := FuzzyMatchWithThreshold(name, candidates); ok {
		return match
	}

	return ""
}

// fuzzyMatchPod tries fuzzy matching on pod names
func (r *ResourceResolver) fuzzyMatchPod(name string, namespace string) string {
	var candidates []string

	for _, pod := range r.patterns.Pods {
		parts := strings.Split(pod, "/")
		if len(parts) != 2 {
			continue
		}

		ns := parts[0]
		podName := parts[1]

		if namespace != "" && strings.ToLower(ns) != strings.ToLower(namespace) {
			continue
		}

		candidates = append(candidates, podName)
	}

	if match, ok := FuzzyMatchWithThreshold(name, candidates); ok {
		return match
	}

	return ""
}

// fuzzyMatchApp tries fuzzy matching on common app names
func (r *ResourceResolver) fuzzyMatchApp(name string) string {
	if match, ok := FuzzyMatchWithThreshold(name, r.patterns.CommonApps); ok {
		return match
	}
	return ""
}

// patternBasedMatch tries to construct resource names based on detected patterns
func (r *ResourceResolver) patternBasedMatch(name string, namespace string) string {
	// Try {app}-{namespace} pattern
	for _, pattern := range r.patterns.Patterns {
		if pattern == "{app}-{namespace}" {
			candidate := name + "-" + namespace
			// Check if this exists in deployments
			for _, deployment := range r.patterns.Deployments {
				parts := strings.Split(deployment, "/")
				if len(parts) == 2 {
					if strings.ToLower(parts[1]) == strings.ToLower(candidate) {
						return parts[1]
					}
				}
			}
		}
	}

	return ""
}

// HasPatterns returns true if the resolver has learned cluster patterns
func (r *ResourceResolver) HasPatterns() bool {
	return len(r.patterns.Deployments) > 0 || len(r.patterns.Namespaces) > 0
}

// GetAllDeploymentNames returns all deployment names (without namespace prefix)
func (r *ResourceResolver) GetAllDeploymentNames() []string {
	var names []string
	for _, deployment := range r.patterns.Deployments {
		parts := strings.Split(deployment, "/")
		if len(parts) == 2 {
			names = append(names, parts[1])
		}
	}
	return names
}

// IsValidNamespace checks if a namespace exists in the cluster
func (r *ResourceResolver) IsValidNamespace(namespace string) bool {
	namespaceLower := strings.ToLower(namespace)
	for _, ns := range r.patterns.Namespaces {
		if strings.ToLower(ns) == namespaceLower {
			return true
		}
	}
	return false
}

// generateNamingVariants creates multiple naming convention variants from user input
// The order is based on the detected cluster naming convention (if available)
// Examples:
//   "my app" -> ["my-app", "myApp", "my_app", "MyApp", "myapp", "my app"]
//   "web server" -> ["web-server", "webServer", "web_server", "WebServer", "webserver", "web server"]
func generateNamingVariants(input string) []string {
	if !strings.Contains(input, " ") {
		// No spaces, return as-is (might already be hyphenated, camelCase, etc.)
		return []string{input}
	}

	words := strings.Fields(input)
	if len(words) == 0 {
		return []string{input}
	}

	// Generate all possible variants
	variantMap := make(map[string]string)

	// Hyphen-separated (most common in Kubernetes)
	variantMap["hyphen"] = strings.ToLower(strings.Join(words, "-"))

	// camelCase (first word lowercase, rest capitalized)
	if len(words) > 1 {
		camelCase := strings.ToLower(words[0])
		for i := 1; i < len(words); i++ {
			if len(words[i]) > 0 {
				camelCase += strings.ToUpper(words[i][:1]) + strings.ToLower(words[i][1:])
			}
		}
		variantMap["camelCase"] = camelCase
	}

	// Underscore-separated
	variantMap["underscore"] = strings.ToLower(strings.Join(words, "_"))

	// PascalCase (all words capitalized)
	pascalCase := ""
	for _, word := range words {
		if len(word) > 0 {
			pascalCase += strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	variantMap["PascalCase"] = pascalCase

	// All lowercase, no separators
	variantMap["nospace"] = strings.ToLower(strings.Join(words, ""))

	// Original input
	variantMap["original"] = input

	// Load cluster patterns to check detected naming convention
	patterns, _ := config.LoadClusterPatterns()
	detectedConvention := ""
	if patterns != nil && patterns.NamingConvention != "" {
		detectedConvention = patterns.NamingConvention
	}

	// Order variants with detected convention first
	variants := make([]string, 0, 6)

	// Add detected convention first (if exists)
	if detectedConvention != "" && variantMap[detectedConvention] != "" {
		variants = append(variants, variantMap[detectedConvention])
		delete(variantMap, detectedConvention)
	}

	// Add remaining variants in priority order
	priority := []string{"hyphen", "camelCase", "underscore", "PascalCase", "nospace", "original"}
	for _, key := range priority {
		if val, exists := variantMap[key]; exists {
			variants = append(variants, val)
		}
	}

	return variants
}
