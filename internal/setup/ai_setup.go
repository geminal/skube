package setup

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/geminal/skube/internal/aiparser"
	"github.com/geminal/skube/internal/config"
)

func RunAISetup() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🤖 AI Features Setup for skube")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	fmt.Println("AI features allow you to use natural language with the --ai flag:")
	fmt.Println("  Example: skube --ai \"show me all crashing pods in production\"")
	fmt.Println()

	// Ask if user wants AI
	fmt.Print("Enable AI features? [y/N]: ")
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response != "y" && response != "yes" {
		cfg := &config.AIConfig{
			Enabled:   false,
			UseDocker: false,
			Model:     "qwen2.5-coder:3b",
		}
		if err := config.SaveAIConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("\n✓ AI features disabled. You can enable them later with: skube setup-ai")
		return nil
	}

	// Check Docker availability
	dockerMgr := aiparser.NewDockerManager()
	hasDocker := dockerMgr.IsDockerAvailable()

	fmt.Println()
	fmt.Println("Choose AI provider:")
	fmt.Println()

	if hasDocker {
		fmt.Println("  1) Docker (Recommended - Free local AI)")
		fmt.Println("     ✓ Automatic setup")
		fmt.Println("     ✓ Downloads Ollama image (~700MB)")
		fmt.Println("     ✓ Pulls AI model (~2GB)")
		fmt.Println("     ✓ No system installation needed")
		fmt.Println()
		fmt.Println("  2) Local Ollama (Free local AI)")
		fmt.Println("     • Requires manual installation")
		fmt.Println("     • Run: curl -fsSL https://ollama.com/install.sh | sh")
		fmt.Println()
		fmt.Println("  3) OpenAI (Requires API key)")
		fmt.Println("     • Uses OpenAI's GPT models")
		fmt.Println("     • Requires internet connection")
		fmt.Println("     • Pay per use")
		fmt.Println()
		fmt.Print("Choose [1/2/3]: ")
	} else {
		fmt.Println("  Docker not detected.")
		fmt.Println()
		fmt.Println("  1) Local Ollama (Free local AI)")
		fmt.Println("     • Install with: curl -fsSL https://ollama.com/install.sh | sh")
		fmt.Println("     • Windows: https://ollama.com/download")
		fmt.Println()
		fmt.Println("  2) OpenAI (Requires API key)")
		fmt.Println("     • Uses OpenAI's GPT models")
		fmt.Println("     • Requires internet connection")
		fmt.Println("     • Pay per use")
		fmt.Println()
		fmt.Print("Choose [1/2]: ")
	}

	response, _ = reader.ReadString('\n')
	response = strings.TrimSpace(response)

	useDocker := false
	provider := "ollama"
	var openAIKey string
	var openAIModel string

	if hasDocker {
		if response == "3" {
			provider = "openai"
		} else if response == "1" || response == "" {
			useDocker = true
		} else if response != "2" {
			fmt.Println("\n✗ Invalid choice")
			return nil
		}
	} else {
		if response == "2" {
			provider = "openai"
		} else if response != "1" {
			fmt.Println("\n✗ Setup cancelled")
			return nil
		}
	}

	// If OpenAI, ask for API key
	if provider == "openai" {
		// Load existing config to check for existing API key
		existingCfg, err := config.LoadAIConfig()
		if err == nil && existingCfg.OpenAIAPIKey != "" {
			fmt.Println()
			fmt.Print("Use existing OpenAI API key? [Y/n]: ")
			response, _ := reader.ReadString('\n')
			response = strings.ToLower(strings.TrimSpace(response))

			if response == "" || response == "y" || response == "yes" {
				openAIKey = existingCfg.OpenAIAPIKey
			}
		}

		if openAIKey == "" {
			fmt.Println()
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println("OpenAI Configuration")
			fmt.Println()
			fmt.Println("Get your API key from: https://platform.openai.com/api-keys")
			fmt.Println()
			fmt.Print("Enter your OpenAI API key: ")
			openAIKey, _ = reader.ReadString('\n')
			openAIKey = strings.TrimSpace(openAIKey)

			if openAIKey == "" {
				fmt.Println("\n✗ API key is required for OpenAI")
				return nil
			}
		}

		// Check for existing model preference
		if err == nil && existingCfg.OpenAIModel != "" && openAIKey == existingCfg.OpenAIAPIKey {
			fmt.Println()
			fmt.Printf("Current OpenAI model: %s\n", existingCfg.OpenAIModel)
			fmt.Print("Keep current model? [Y/n]: ")
			response, _ := reader.ReadString('\n')
			response = strings.ToLower(strings.TrimSpace(response))

			if response == "" || response == "y" || response == "yes" {
				openAIModel = existingCfg.OpenAIModel
			}
		}

		if openAIModel == "" {
			fmt.Println()
			fmt.Println("Choose OpenAI model:")
			fmt.Println("  1) gpt-4o-mini (Recommended - Fast and affordable)")
			fmt.Println("  2) gpt-4o (More accurate, higher cost)")
			fmt.Println("  3) gpt-3.5-turbo (Legacy, cheapest)")
			fmt.Print("\nChoose [1/2/3]: ")
			modelChoice, _ := reader.ReadString('\n')
			modelChoice = strings.TrimSpace(modelChoice)

			switch modelChoice {
			case "2":
				openAIModel = "gpt-4o"
			case "3":
				openAIModel = "gpt-3.5-turbo"
			default:
				openAIModel = "gpt-4o-mini"
			}
		}
	}

	cfg := &config.AIConfig{
		Enabled:      true,
		Provider:     provider,
		UseDocker:    useDocker,
		Model:        "qwen2.5-coder:3b",
		OpenAIAPIKey: openAIKey,
		OpenAIModel:  openAIModel,
	}

	if err := config.SaveAIConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Check if user has skube-config.json in current directory
	localConfigPath := "skube-config.json"
	if _, err := os.Stat(localConfigPath); err == nil {
		fmt.Println("\n📋 Found skube-config.json in current directory!")
		fmt.Print("   Import this config? [Y/n]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		if response == "" || response == "y" || response == "yes" {
			// Read and merge the local config
			localData, err := os.ReadFile(localConfigPath)
			if err == nil {
				var localCfg config.AIConfig
				if json.Unmarshal(localData, &localCfg) == nil {
					// Preserve enabled and useDocker from setup, but import everything else
					localCfg.Enabled = cfg.Enabled
					localCfg.UseDocker = cfg.UseDocker
					localCfg.Model = cfg.Model

					// Preserve provider if not set in local config
					if localCfg.Provider == "" {
						localCfg.Provider = cfg.Provider
					}
					// Preserve OpenAI settings if not set in local config
					if localCfg.OpenAIAPIKey == "" {
						localCfg.OpenAIAPIKey = cfg.OpenAIAPIKey
					}
					if localCfg.OpenAIModel == "" {
						localCfg.OpenAIModel = cfg.OpenAIModel
					}

					cfg = &localCfg

					if err := config.SaveAIConfig(cfg); err == nil {
						fmt.Println("   ✓ Config imported successfully!")
					}
				}
			}
		}
	}

	// Show config location
	configPath := config.GetConfigPath()
	fmt.Printf("\n✓ Configuration saved to: %s\n", configPath)

	// If OpenAI, we're done
	if provider == "openai" {
		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("✓ OpenAI setup complete!")
		fmt.Println()
		fmt.Printf("Using model: %s\n", openAIModel)
		fmt.Println()
		fmt.Println("Try it out:")
		fmt.Println("  skube --ai \"in production get pods\"")
		fmt.Println("  skube --ai \"show me all crashing deployments\"")
		fmt.Println()
		return nil
	}

	// If Docker, set up now
	if useDocker {
		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("Setting up Ollama in Docker...")
		fmt.Println()

		// Start container
		if err := dockerMgr.StartOllamaContainer(); err != nil {
			return fmt.Errorf("failed to start Ollama container: %w", err)
		}

		// Pull model
		fmt.Println()
		fmt.Printf("Pulling AI model '%s' (this may take a few minutes)...\n", cfg.Model)
		if err := dockerMgr.PullModelInContainer(cfg.Model); err != nil {
			return fmt.Errorf("failed to pull model: %w", err)
		}

		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("✓ Setup complete!")
		fmt.Println()
		fmt.Println("Try it out:")
		fmt.Println("  skube --ai \"in production get pods\"")
		fmt.Println("  skube --ai \"show me all crashing deployments\"")
		fmt.Println()
	} else {
		fmt.Println()
		fmt.Println("✓ Configuration saved!")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("  1. Install Ollama: curl -fsSL https://ollama.com/install.sh | sh")
		fmt.Println("  2. Start Ollama: ollama serve")
		fmt.Println("  3. Pull model: ollama pull qwen2.5-coder:3b")
		fmt.Println("  4. Use skube: skube --ai \"in production get pods\"")
		fmt.Println()
	}

	return nil
}

func SwitchProvider() error {
	reader := bufio.NewReader(os.Stdin)

	// Load current config
	cfg, err := config.LoadAIConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if !cfg.Enabled {
		fmt.Println("❌ AI features are currently disabled.")
		fmt.Print("\nEnable AI features? [y/N]: ")
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))

		if response != "y" && response != "yes" {
			fmt.Println("\n✓ AI remains disabled. Run 'skube setup-ai' to enable.")
			return nil
		}
		cfg.Enabled = true
	}

	// Show current provider
	currentProvider := cfg.Provider
	if currentProvider == "" {
		if cfg.OpenAIAPIKey != "" {
			currentProvider = "openai"
		} else {
			currentProvider = "ollama"
		}
	}

	fmt.Println("🔄 Switch AI Provider")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("\nCurrent provider: %s\n", currentProvider)
	fmt.Println()

	// Check Docker availability
	dockerMgr := aiparser.NewDockerManager()
	hasDocker := dockerMgr.IsDockerAvailable()

	fmt.Println("Available providers:")
	fmt.Println()

	optionMap := make(map[string]string)
	option := 1

	if hasDocker {
		fmt.Printf("  %d) Docker (Free local AI)\n", option)
		fmt.Println("     ✓ Automatic setup")
		fmt.Println("     ✓ No system installation needed")
		optionMap[fmt.Sprintf("%d", option)] = "docker"
		option++
		fmt.Println()
	}

	fmt.Printf("  %d) Local Ollama (Free local AI)\n", option)
	fmt.Println("     • Requires manual installation")
	optionMap[fmt.Sprintf("%d", option)] = "ollama"
	option++
	fmt.Println()

	fmt.Printf("  %d) OpenAI (Requires API key)\n", option)
	fmt.Println("     • Uses OpenAI's GPT models")
	fmt.Println("     • Pay per use")
	optionMap[fmt.Sprintf("%d", option)] = "openai"
	fmt.Println()

	fmt.Printf("Choose [1-%d]: ", option-1)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	selectedProvider, ok := optionMap[response]
	if !ok {
		fmt.Println("\n✗ Invalid choice")
		return nil
	}

	// Handle provider switch
	switch selectedProvider {
	case "docker":
		cfg.Provider = "ollama"
		cfg.UseDocker = true

		if err := config.SaveAIConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("Setting up Ollama in Docker...")
		fmt.Println()

		if err := dockerMgr.StartOllamaContainer(); err != nil {
			return fmt.Errorf("failed to start Ollama container: %w", err)
		}

		fmt.Println()
		fmt.Printf("Pulling AI model '%s' (this may take a few minutes)...\n", cfg.Model)
		if err := dockerMgr.PullModelInContainer(cfg.Model); err != nil {
			return fmt.Errorf("failed to pull model: %w", err)
		}

		fmt.Println()
		fmt.Println("✓ Switched to Docker-based Ollama!")

	case "ollama":
		cfg.Provider = "ollama"
		cfg.UseDocker = false

		if err := config.SaveAIConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println()
		fmt.Println("✓ Switched to local Ollama!")
		fmt.Println()
		fmt.Println("Make sure Ollama is installed and running:")
		fmt.Println("  1. Install: curl -fsSL https://ollama.com/install.sh | sh")
		fmt.Println("  2. Start: ollama serve")
		fmt.Println("  3. Pull model: ollama pull qwen2.5-coder:3b")

	case "openai":
		var openAIKey string
		var openAIModel string

		// Check if API key already exists
		if cfg.OpenAIAPIKey != "" {
			fmt.Println()
			fmt.Print("Use existing OpenAI API key? [Y/n]: ")
			response, _ := reader.ReadString('\n')
			response = strings.ToLower(strings.TrimSpace(response))

			if response == "" || response == "y" || response == "yes" {
				openAIKey = cfg.OpenAIAPIKey
			}
		}

		if openAIKey == "" {
			fmt.Println()
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println("OpenAI Configuration")
			fmt.Println()
			fmt.Println("Get your API key from: https://platform.openai.com/api-keys")
			fmt.Println()
			fmt.Print("Enter your OpenAI API key: ")
			openAIKey, _ = reader.ReadString('\n')
			openAIKey = strings.TrimSpace(openAIKey)

			if openAIKey == "" {
				fmt.Println("\n✗ API key is required for OpenAI")
				return nil
			}
		}

		// Check if model already exists
		if cfg.OpenAIModel != "" {
			fmt.Println()
			fmt.Printf("Current OpenAI model: %s\n", cfg.OpenAIModel)
			fmt.Print("Keep current model? [Y/n]: ")
			response, _ := reader.ReadString('\n')
			response = strings.ToLower(strings.TrimSpace(response))

			if response == "" || response == "y" || response == "yes" {
				openAIModel = cfg.OpenAIModel
			}
		}

		if openAIModel == "" {
			fmt.Println()
			fmt.Println("Choose OpenAI model:")
			fmt.Println("  1) gpt-4o-mini (Recommended - Fast and affordable)")
			fmt.Println("  2) gpt-4o (More accurate, higher cost)")
			fmt.Println("  3) gpt-3.5-turbo (Legacy, cheapest)")
			fmt.Print("\nChoose [1/2/3]: ")
			modelChoice, _ := reader.ReadString('\n')
			modelChoice = strings.TrimSpace(modelChoice)

			switch modelChoice {
			case "2":
				openAIModel = "gpt-4o"
			case "3":
				openAIModel = "gpt-3.5-turbo"
			default:
				openAIModel = "gpt-4o-mini"
			}
		}

		cfg.Provider = "openai"
		cfg.UseDocker = false
		cfg.OpenAIAPIKey = openAIKey
		cfg.OpenAIModel = openAIModel

		if err := config.SaveAIConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		fmt.Println()
		fmt.Println("✓ Switched to OpenAI!")
		fmt.Printf("  Using model: %s\n", openAIModel)
	}

	fmt.Println()
	fmt.Println("Try it out:")
	fmt.Println("  skube --ai \"in production get pods\"")
	fmt.Println()

	return nil
}
