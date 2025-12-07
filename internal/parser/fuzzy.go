package parser

import (
	"strings"
)

// LevenshteinDistance calculates the edit distance between two strings.
// This is a hand-rolled implementation to avoid external dependencies.
// Time complexity: O(len(s1) * len(s2))
// Space complexity: O(len(s2))
func LevenshteinDistance(s1, s2 string) int {
	// Convert to lowercase for case-insensitive comparison
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	// Quick equality check
	if s1 == s2 {
		return 0
	}

	// If either string is empty, distance is the length of the other
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Use runes to handle multi-byte characters correctly
	r1 := []rune(s1)
	r2 := []rune(s2)

	// Optimize memory by using only two rows
	prevRow := make([]int, len(r2)+1)
	currRow := make([]int, len(r2)+1)

	// Initialize first row (distance from empty string)
	for j := 0; j <= len(r2); j++ {
		prevRow[j] = j
	}

	// Calculate edit distance using dynamic programming
	for i := 1; i <= len(r1); i++ {
		currRow[0] = i

		for j := 1; j <= len(r2); j++ {
			cost := 1
			if r1[i-1] == r2[j-1] {
				cost = 0
			}

			currRow[j] = min(
				currRow[j-1]+1,    // insertion
				prevRow[j]+1,      // deletion
				prevRow[j-1]+cost, // substitution
			)
		}

		// Swap rows
		prevRow, currRow = currRow, prevRow
	}

	return prevRow[len(r2)]
}

// min returns the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

// FuzzyMatch returns true if s1 and s2 are similar enough.
// maxDistance defines the maximum allowed edit distance.
// For strings longer than 10 chars, we allow more tolerance.
func FuzzyMatch(s1, s2 string, maxDistance int) bool {
	distance := LevenshteinDistance(s1, s2)
	return distance <= maxDistance
}

// FindClosestMatch finds the closest matching string from a list of candidates.
// Returns the best match and its edit distance, or ("", -1) if no candidates.
func FindClosestMatch(target string, candidates []string) (string, int) {
	if len(candidates) == 0 {
		return "", -1
	}

	bestMatch := candidates[0]
	bestDistance := LevenshteinDistance(target, candidates[0])

	for i := 1; i < len(candidates); i++ {
		distance := LevenshteinDistance(target, candidates[i])
		if distance < bestDistance {
			bestMatch = candidates[i]
			bestDistance = distance
		}
	}

	return bestMatch, bestDistance
}

// FuzzyMatchWithThreshold finds the closest match if it's within the threshold.
// Returns (match, true) if a good match is found, ("", false) otherwise.
// The threshold is adaptive based on string length.
func FuzzyMatchWithThreshold(target string, candidates []string) (string, bool) {
	if len(candidates) == 0 {
		return "", false
	}

	match, distance := FindClosestMatch(target, candidates)

	// Adaptive threshold based on string length
	threshold := calculateThreshold(target)

	if distance <= threshold {
		return match, true
	}

	return "", false
}

// calculateThreshold returns an adaptive threshold based on string length.
// Short strings (< 5 chars): allow 1 typo
// Medium strings (5-10 chars): allow 2 typos
// Long strings (> 10 chars): allow 3 typos
func calculateThreshold(s string) int {
	length := len(s)
	if length <= 4 {
		return 1
	}
	if length <= 10 {
		return 2
	}
	return 3
}

// NormalizeSpacesToHyphens converts spaces to hyphens for resource name matching.
// Also handles underscores by converting them to hyphens.
func NormalizeSpacesToHyphens(s string) string {
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	return strings.ToLower(s)
}

// ContainsFuzzy checks if target appears within any candidate with fuzzy matching.
// Useful for partial name matching (e.g., "web" matches "web-server").
func ContainsFuzzy(target string, candidates []string, maxDistance int) (string, bool) {
	targetLower := strings.ToLower(target)

	for _, candidate := range candidates {
		candidateLower := strings.ToLower(candidate)

		// Exact substring match
		if strings.Contains(candidateLower, targetLower) {
			return candidate, true
		}

		// Fuzzy match on the whole string
		if FuzzyMatch(target, candidate, maxDistance) {
			return candidate, true
		}

		// Check if target fuzzy-matches any part of the candidate
		// Split by common separators
		parts := strings.FieldsFunc(candidateLower, func(r rune) bool {
			return r == '-' || r == '_' || r == '.'
		})

		for _, part := range parts {
			if FuzzyMatch(target, part, maxDistance) {
				return candidate, true
			}
		}
	}

	return "", false
}
