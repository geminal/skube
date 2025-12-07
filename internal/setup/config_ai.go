package setup

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/geminal/skube/internal/config"
)

func ImportAIConfig(jsonFilePath string) error {
	// Read the user's JSON file
	data, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON to validate it
	var cfg config.AIConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Save to skube config location
	if err := config.SaveAIConfig(&cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	configPath := config.GetConfigPath()
	fmt.Printf("✓ AI configuration imported successfully!\n")
	fmt.Printf("  Saved to: %s\n\n", configPath)

	// Show what was configured
	if len(cfg.CommonApps) > 0 {
		fmt.Printf("  Common apps: %d configured\n", len(cfg.CommonApps))
	}
	if len(cfg.Namespaces) > 0 {
		fmt.Printf("  Namespaces: %d configured\n", len(cfg.Namespaces))
	}
	if len(cfg.AppPatterns) > 0 {
		fmt.Printf("  App patterns: %d configured\n", len(cfg.AppPatterns))
	}
	if len(cfg.CustomHints) > 0 {
		fmt.Printf("  Custom hints: %d configured\n", len(cfg.CustomHints))
	}

	fmt.Println("\nYour AI is now customized for your cluster!")
	fmt.Println("Try it: skube --ai \"your natural language command\"")

	return nil
}
