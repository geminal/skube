package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

type AIConfig struct {
	Enabled      bool              `json:"enabled"`
	Provider     string            `json:"provider,omitempty"` // "ollama" or "openai"
	UseDocker    bool              `json:"use_docker"`
	Model        string            `json:"model"`
	OpenAIAPIKey string            `json:"openai_api_key,omitempty"`
	OpenAIModel  string            `json:"openai_model,omitempty"`
	AppPatterns  []string          `json:"app_patterns,omitempty"`
	CommonApps   []string          `json:"common_apps,omitempty"`
	Namespaces   []string          `json:"namespaces,omitempty"`
	CustomHints  map[string]string `json:"custom_hints,omitempty"`
}

func GetConfigPath() string {
	var configDir string

	switch runtime.GOOS {
	case "darwin", "linux":
		home, err := os.UserHomeDir()
		if err != nil {
			configDir = "/tmp/skube"
		} else {
			configDir = filepath.Join(home, ".config", "skube")
		}
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

	return filepath.Join(configDir, "config.json")
}

func LoadAIConfig() (*AIConfig, error) {
	configPath := GetConfigPath()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &AIConfig{
				Enabled:      false,
				Provider:     "ollama",
				UseDocker:    false,
				Model:        "qwen2.5:3b",
				OpenAIModel:  "gpt-4o-mini",
				AppPatterns:  []string{},
				CommonApps:   []string{},
				Namespaces:   []string{},
				CustomHints:  map[string]string{},
			}, nil
		}
		return nil, err
	}

	// Strip comments (lines starting with //) before parsing JSON
	cleanData := stripJSONComments(data)

	var cfg AIConfig
	if err := json.Unmarshal(cleanData, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func stripJSONComments(data []byte) []byte {
	lines := []byte{}
	inString := false
	escape := false

	for i := 0; i < len(data); i++ {
		char := data[i]

		// Track if we're inside a string
		if char == '"' && !escape {
			inString = !inString
		}
		escape = char == '\\' && !escape

		// Skip comment lines (only outside strings)
		if !inString && i < len(data)-1 && char == '/' && data[i+1] == '/' {
			// Skip until end of line
			for i < len(data) && data[i] != '\n' {
				i++
			}
			continue
		}

		lines = append(lines, char)
	}

	return lines
}

func SaveAIConfig(cfg *AIConfig) error {
	configPath := GetConfigPath()
	configDir := filepath.Dir(configPath)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	// Add helpful comment at the top
	configWithComment := []byte("// Customize AI behavior for your cluster\n" +
		"// Add your app names, namespaces, and patterns below\n" +
		"// See: https://github.com/geminal/skube#customizing-ai-for-your-cluster\n" +
		string(data))

	return os.WriteFile(configPath, configWithComment, 0644)
}
