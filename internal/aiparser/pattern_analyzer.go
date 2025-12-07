package aiparser

import (
	"fmt"
	"strings"
)

// AnalyzePatterns examines a list of resource names to find common naming conventions
// It returns a list of human-readable hints about detected patterns
func AnalyzePatterns(resources []string) []string {
	if len(resources) < 3 {
		return nil
	}

	var hints []string
	suffixes := make(map[string]int)
	prefixes := make(map[string]int)
	total := len(resources)

	// Analyze suffixes (e.g., -qa, -prod, -svc)
	for _, res := range resources {
		// Strip namespace if present (ns/name)
		name := res
		if parts := strings.Split(res, "/"); len(parts) == 2 {
			name = parts[1]
		}

		parts := strings.Split(name, "-")
		if len(parts) > 1 {
			// Check last part as suffix
			suffix := "-" + parts[len(parts)-1]
			suffixes[suffix]++

			// Check first part as prefix
			prefix := parts[0] + "-"
			prefixes[prefix]++
		}
	}

	// Threshold: if > 30% of resources share a pattern (and at least 3 items)
	threshold := total / 3
	if threshold < 3 {
		threshold = 3
	}

	for suffix, count := range suffixes {
		if count >= threshold {
			hints = append(hints, fmt.Sprintf("Detected pattern: Many resources end with '%s' (e.g. app%s)", suffix, suffix))
		}
	}

	for prefix, count := range prefixes {
		if count >= threshold {
			hints = append(hints, fmt.Sprintf("Detected pattern: Many resources start with '%s' (e.g. %sapp)", prefix, prefix))
		}
	}

	return hints
}

// detectNamespaceSuffixPattern checks if resources follow the pattern: app-name-{namespace}
// Returns true if a significant portion of resources end with their namespace name
func detectNamespaceSuffixPattern(namespaceMap map[string][]string) bool {
	totalResources := 0
	matchingResources := 0

	for ns, resources := range namespaceMap {
		for _, res := range resources {
			totalResources++
			// Check if resource ends with "-{namespace}"
			if strings.HasSuffix(res, "-"+ns) {
				matchingResources++
			}
		}
	}

	// If at least 40% of resources follow this pattern, consider it detected
	if totalResources > 0 && float64(matchingResources)/float64(totalResources) >= 0.4 {
		return true
	}

	return false
}
