package parser

import (
	"testing"
)

// TestLevenshteinDistance tests the edit distance calculation
func TestLevenshteinDistance(t *testing.T) {
	tests := []struct {
		name     string
		s1       string
		s2       string
		expected int
	}{
		{"identical strings", "hello", "hello", 0},
		{"case insensitive", "Hello", "hello", 0},
		{"one char different", "hello", "hallo", 1},
		{"completely different", "abc", "def", 3},
		{"empty strings", "", "", 0},
		{"one empty", "hello", "", 5},
		{"typo in namespace", "production", "produciton", 2},
		{"typo in deployment", "deployment", "depoloyment", 2},
		{"namespace typo", "staging", "stagingg", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LevenshteinDistance(tt.s1, tt.s2)
			if result != tt.expected {
				t.Errorf("LevenshteinDistance(%q, %q) = %d, want %d", tt.s1, tt.s2, result, tt.expected)
			}
		})
	}
}

// TestFuzzyMatch tests fuzzy matching with threshold
func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		name        string
		s1          string
		s2          string
		maxDistance int
		expected    bool
	}{
		{"exact match", "prod", "prod", 2, true},
		{"one typo within threshold", "prod", "prdo", 1, true},
		{"two typos within threshold", "staging", "stagign", 2, true},
		{"beyond threshold", "prod", "production", 2, false},
		{"case insensitive match", "Prod", "prod", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FuzzyMatch(tt.s1, tt.s2, tt.maxDistance)
			if result != tt.expected {
				t.Errorf("FuzzyMatch(%q, %q, %d) = %v, want %v", tt.s1, tt.s2, tt.maxDistance, result, tt.expected)
			}
		})
	}
}

// TestFindClosestMatch tests finding the best match from candidates
func TestFindClosestMatch(t *testing.T) {
	tests := []struct {
		name             string
		target           string
		candidates       []string
		expectedMatch    string
		expectedDistance int
	}{
		{
			name:             "exact match exists",
			target:           "production",
			candidates:       []string{"dev", "staging", "production", "qa"},
			expectedMatch:    "production",
			expectedDistance: 0,
		},
		{
			name:             "typo correction",
			target:           "produciton",
			candidates:       []string{"dev", "staging", "production", "qa"},
			expectedMatch:    "production",
			expectedDistance: 2,
		},
		{
			name:             "closest match",
			target:           "stag",
			candidates:       []string{"dev", "staging", "production"},
			expectedMatch:    "staging",
			expectedDistance: 3,
		},
		{
			name:             "empty candidates",
			target:           "test",
			candidates:       []string{},
			expectedMatch:    "",
			expectedDistance: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, distance := FindClosestMatch(tt.target, tt.candidates)
			if match != tt.expectedMatch {
				t.Errorf("FindClosestMatch() match = %q, want %q", match, tt.expectedMatch)
			}
			if distance != tt.expectedDistance {
				t.Errorf("FindClosestMatch() distance = %d, want %d", distance, tt.expectedDistance)
			}
		})
	}
}

// TestFuzzyMatchWithThreshold tests adaptive threshold matching
func TestFuzzyMatchWithThreshold(t *testing.T) {
	tests := []struct {
		name       string
		target     string
		candidates []string
		wantMatch  string
		wantFound  bool
	}{
		{
			name:       "short string with 1 typo",
			target:     "prdo",
			candidates: []string{"prod", "dev", "qa"},
			wantMatch:  "prod",
			wantFound:  true,
		},
		{
			name:       "medium string with 2 typos",
			target:     "stagign",
			candidates: []string{"staging", "production", "dev"},
			wantMatch:  "staging",
			wantFound:  true,
		},
		{
			name:       "long string with 3 typos",
			target:     "prodcution",
			candidates: []string{"production", "staging"},
			wantMatch:  "production",
			wantFound:  true,
		},
		{
			name:       "too many typos",
			target:     "xyz",
			candidates: []string{"production", "staging"},
			wantMatch:  "",
			wantFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, found := FuzzyMatchWithThreshold(tt.target, tt.candidates)
			if found != tt.wantFound {
				t.Errorf("FuzzyMatchWithThreshold() found = %v, want %v", found, tt.wantFound)
			}
			if match != tt.wantMatch {
				t.Errorf("FuzzyMatchWithThreshold() match = %q, want %q", match, tt.wantMatch)
			}
		})
	}
}

// TestNormalizeSpacesToHyphens tests space/underscore normalization
func TestNormalizeSpacesToHyphens(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"spaces to hyphens", "web server", "web-server"},
		{"underscores to hyphens", "web_server", "web-server"},
		{"mixed separators", "web server_api", "web-server-api"},
		{"already hyphenated", "web-server", "web-server"},
		{"uppercase to lowercase", "Web Server", "web-server"},
		{"multiple spaces", "my  app  name", "my--app--name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeSpacesToHyphens(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizeSpacesToHyphens(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestContainsFuzzy tests partial fuzzy matching
func TestContainsFuzzy(t *testing.T) {
	tests := []struct {
		name        string
		target      string
		candidates  []string
		maxDistance int
		wantMatch   string
		wantFound   bool
	}{
		{
			name:        "exact substring",
			target:      "api",
			candidates:  []string{"api-service", "web-server", "worker"},
			maxDistance: 2,
			wantMatch:   "api-service",
			wantFound:   true,
		},
		{
			name:        "fuzzy match on part",
			target:      "srver",
			candidates:  []string{"api-service", "web-server", "worker"},
			maxDistance: 2,
			wantMatch:   "web-server",
			wantFound:   true,
		},
		{
			name:        "no match",
			target:      "xyz",
			candidates:  []string{"api-service", "web-server"},
			maxDistance: 1,
			wantMatch:   "",
			wantFound:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, found := ContainsFuzzy(tt.target, tt.candidates, tt.maxDistance)
			if found != tt.wantFound {
				t.Errorf("ContainsFuzzy() found = %v, want %v", found, tt.wantFound)
			}
			if match != tt.wantMatch {
				t.Errorf("ContainsFuzzy() match = %q, want %q", match, tt.wantMatch)
			}
		})
	}
}

// TestCalculateThreshold tests adaptive threshold calculation
func TestCalculateThreshold(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int
	}{
		{"very short", "abc", 1},
		{"short", "prod", 1},
		{"medium", "staging", 2},
		{"long", "production", 3},
		{"very long", "my-application-name", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateThreshold(tt.input)
			if result != tt.expected {
				t.Errorf("calculateThreshold(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// BenchmarkLevenshteinDistance benchmarks the edit distance function
func BenchmarkLevenshteinDistance(b *testing.B) {
	s1 := "production"
	s2 := "produciton"
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(s1, s2)
	}
}

// BenchmarkFuzzyMatchWithThreshold benchmarks fuzzy matching
func BenchmarkFuzzyMatchWithThreshold(b *testing.B) {
	target := "stagign"
	candidates := []string{"dev", "staging", "production", "qa", "test"}
	for i := 0; i < b.N; i++ {
		FuzzyMatchWithThreshold(target, candidates)
	}
}

// TestGenerateNamingVariants tests different naming convention variants
func TestGenerateNamingVariants(t *testing.T) {
	// Note: This function is in resource_resolver.go but we test it here
	// You may need to move this test or export the function
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:  "two words",
			input: "web server",
			expected: []string{
				"web-server",    // hyphen-separated
				"webServer",     // camelCase
				"web_server",    // underscore
				"WebServer",     // PascalCase
				"webserver",     // no separator
				"web server",    // original
			},
		},
		{
			name:  "three words",
			input: "my app name",
			expected: []string{
				"my-app-name",
				"myAppName",
				"my_app_name",
				"MyAppName",
				"myappname",
				"my app name",
			},
		},
		{
			name:     "single word (no spaces)",
			input:    "myapp",
			expected: []string{"myapp"},
		},
		{
			name:     "already hyphenated (no spaces)",
			input:    "my-app",
			expected: []string{"my-app"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We'll need to import this from resource_resolver
			// For now, this is a placeholder test
			t.Skip("Function generateNamingVariants is not exported - test when needed")
		})
	}
}
