package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/geminal/skube/internal/config"
	"github.com/geminal/skube/internal/executor"
	"github.com/geminal/skube/internal/parser"
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
		executor.PrintHelp()
		os.Exit(0)
	}

	args := os.Args[1:]
	ctx := parser.ParseNaturalLanguage(args)

	if err := executor.ExecuteCommand(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%sError: %v%s\n", config.ColorRed, err, config.ColorReset)
		os.Exit(1)
	}
}
