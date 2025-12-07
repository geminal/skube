package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/geminal/skube/internal/aiparser"
	"github.com/geminal/skube/internal/cluster"
	"github.com/geminal/skube/internal/config"
	"github.com/geminal/skube/internal/executor"
	"github.com/geminal/skube/internal/help"
	"github.com/geminal/skube/internal/parser"
	"github.com/geminal/skube/internal/setup"
)

func main() {
	// Check if kubectl is installed
	if _, err := exec.LookPath("kubectl"); err != nil {
		fmt.Fprintf(os.Stderr, "%s⚠️  kubectl not found in PATH%s\n\n", config.ColorYellow, config.ColorReset)
		fmt.Println("skube requires kubectl to interact with Kubernetes clusters.")
		fmt.Println("\nTo install kubectl, choose your platform:")

		// Detect OS and provide specific instructions
		switch runtime.GOOS {
		case "darwin":
			fmt.Println("  macOS:")
			fmt.Println("    brew install kubectl")
			fmt.Println("    # or")
			fmt.Println("    curl -LO \"https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/darwin/amd64/kubectl\"")
		case "linux":
			fmt.Println("  Linux:")
			fmt.Println("    # Debian/Ubuntu")
			fmt.Println("    sudo apt-get update && sudo apt-get install -y kubectl")
			fmt.Println("    # or")
			fmt.Println("    curl -LO \"https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl\"")
		case "windows":
			fmt.Println("  Windows:")
			fmt.Println("    choco install kubernetes-cli")
			fmt.Println("    # or download from:")
			fmt.Println("    # https://kubernetes.io/docs/tasks/tools/install-kubectl-windows/")
		default:
			fmt.Println("  Visit: https://kubernetes.io/docs/tasks/tools/")
		}

		fmt.Printf("\n%sAfter installing kubectl, run skube again.%s\n", config.ColorCyan, config.ColorReset)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		help.PrintHelp()
		os.Exit(0)
	}

	// Check for version flag
	if os.Args[1] == "--version" || os.Args[1] == "-v" || os.Args[1] == "version" {
		help.PrintVersion()
		os.Exit(0)
	}

	// Check for setup-ai command
	if os.Args[1] == "setup-ai" {
		if err := setup.RunAISetup(); err != nil {
			fmt.Fprintf(os.Stderr, "%sError: %v%s\n", config.ColorRed, err, config.ColorReset)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check for switch-ai command
	if os.Args[1] == "switch-ai" {
		if err := setup.SwitchProvider(); err != nil {
			fmt.Fprintf(os.Stderr, "%sError: %v%s\n", config.ColorRed, err, config.ColorReset)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check for config-ai command
	if os.Args[1] == "config-ai" {
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "%sUsage: skube config-ai <config-file.json>%s\n\n", config.ColorYellow, config.ColorReset)
			fmt.Println("Steps:")
			fmt.Println("  1. Copy the example: cp ai-config.example.json my-config.json")
			fmt.Println("  2. Edit my-config.json with your apps, namespaces, and patterns")
			fmt.Println("  3. Import it: skube config-ai my-config.json")
			fmt.Println("\nSee: https://github.com/geminal/skube#customizing-ai-for-your-cluster")
			os.Exit(1)
		}

		if err := setup.ImportAIConfig(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "%sError: %v%s\n", config.ColorRed, err, config.ColorReset)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Check for init command (learn cluster patterns)
	if os.Args[1] == "init" || os.Args[1] == "refresh-patterns" {
		patterns, err := cluster.LearnClusterPatterns(true)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError learning cluster patterns: %v%s\n", config.ColorRed, err, config.ColorReset)
			os.Exit(1)
		}

		if err := config.SaveClusterPatterns(patterns); err != nil {
			fmt.Fprintf(os.Stderr, "%sError saving cluster patterns: %v%s\n", config.ColorRed, err, config.ColorReset)
			os.Exit(1)
		}

		os.Exit(0)
	}

	// Check if patterns cache is stale and auto-refresh
	if config.IsClusterPatternsCacheStale() {
		// Get current context for display
		currentContext, _ := config.GetCurrentKubeContext()
		if currentContext != "" {
			fmt.Printf("%sRefreshing cluster patterns for %s%s%s...%s ",
				config.ColorYellow, config.ColorCyan, currentContext, config.ColorYellow, config.ColorReset)
		} else {
			fmt.Printf("%sRefreshing cluster patterns...%s ", config.ColorYellow, config.ColorReset)
		}

		patterns, err := cluster.LearnClusterPatterns(false)
		if err == nil {
			config.SaveClusterPatterns(patterns)
			fmt.Printf("%s✓%s\n", config.ColorGreen, config.ColorReset)
		} else {
			fmt.Printf("%s(skipped)%s\n", config.ColorYellow, config.ColorReset)
		}
	}

	// Check for model command
	if os.Args[1] == "model" {
		cfg, err := config.LoadAIConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%sError loading AI config: %v%s\n", config.ColorRed, err, config.ColorReset)
			os.Exit(1)
		}

		fmt.Printf("%sAI Configuration:%s\n", config.ColorGreen, config.ColorReset)
		fmt.Printf("  Status:   %s", config.ColorCyan)
		if cfg.Enabled {
			fmt.Printf("Enabled%s\n", config.ColorReset)
		} else {
			fmt.Printf("Disabled%s\n", config.ColorReset)
		}

		fmt.Printf("  Provider: %s%s%s\n", config.ColorCyan, cfg.Provider, config.ColorReset)

		if cfg.Provider == "ollama" {
			fmt.Printf("  Model:    %s%s%s\n", config.ColorCyan, cfg.Model, config.ColorReset)
			fmt.Printf("  Docker:   %s", config.ColorCyan)
			if cfg.UseDocker {
				fmt.Printf("Yes%s\n", config.ColorReset)
			} else {
				fmt.Printf("No%s\n", config.ColorReset)
			}
		} else if cfg.Provider == "openai" {
			fmt.Printf("  Model:    %s%s%s\n", config.ColorCyan, cfg.OpenAIModel, config.ColorReset)
			if cfg.OpenAIAPIKey != "" {
				fmt.Printf("  API Key:  %s[configured]%s\n", config.ColorCyan, config.ColorReset)
			} else {
				fmt.Printf("  API Key:  %s[not configured]%s\n", config.ColorYellow, config.ColorReset)
			}
		}

		configPath := config.GetConfigPath()
		fmt.Printf("\n%sConfig file: %s%s\n", config.ColorYellow, configPath, config.ColorReset)

		if !cfg.Enabled {
			fmt.Printf("\n%sTo enable AI features, run: %sskube setup-ai%s\n", config.ColorYellow, config.ColorCyan, config.ColorReset)
		}

		os.Exit(0)
	}

	args := os.Args[1:]

	var ctx *parser.Context
	if aiparser.HasAIFlag(args) {
		args = aiparser.StripAIFlag(args)
		aiParser, err := aiparser.NewAIParser()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s⚠️  AI Setup Error: %v%s\n", config.ColorYellow, err, config.ColorReset)
			fmt.Fprintf(os.Stderr, "%sFalling back to regular parser...%s\n\n", config.ColorYellow, config.ColorReset)
			ctx = parser.ParseNaturalLanguage(args)
		} else {
			ctx, err = aiParser.Parse(args)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s⚠️  AI Parse Error: %v%s\n", config.ColorYellow, err, config.ColorReset)
				fmt.Fprintf(os.Stderr, "%sFalling back to regular parser...%s\n\n", config.ColorYellow, config.ColorReset)
				ctx = parser.ParseNaturalLanguage(args)
			}
		}
	} else {
		ctx = parser.ParseNaturalLanguage(args)
	}

	if err := executor.ExecuteCommand(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", config.ColorRed, err, config.ColorReset)
		os.Exit(1)
	}
}
