package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/geminal/skube/internal/completion"
	"github.com/geminal/skube/internal/config"
	"github.com/geminal/skube/internal/help"
	"github.com/geminal/skube/internal/parser"
)

// execCommand is a variable to allow mocking in tests
var execCommand = exec.Command

func ExecuteCommand(ctx *parser.Context) error {
	// Sanitize inputs
	ctx.SearchTerm = sanitizeInput(ctx.SearchTerm)
	ctx.FilePath = sanitizeInput(ctx.FilePath)
	ctx.SourcePath = sanitizeInput(ctx.SourcePath)
	ctx.DestPath = sanitizeInput(ctx.DestPath)

	if ctx.Command == "" {
		help.PrintHelp()
		return nil
	}

	switch ctx.Command {
	case "logs":
		return handleLogs(ctx)
	case "shell":
		return handleShell(ctx)
	case "restart":
		return handleRestart(ctx)
	case "pods":
		return handlePods(ctx)
	case "scale":
		return handleScale(ctx)
	case "rollback":
		return handleRollback(ctx)
	case "forward":
		return handlePortForward(ctx)
	case "describe":
		return handleDescribe(ctx)
	case "services":
		return handleServices(ctx)
	case "deployments":
		return handleDeployments(ctx)
	case "namespaces":
		return handleNamespaces(ctx)
	case "nodes":
		return handleNodes(ctx)
	case "configmaps":
		return handleConfigMaps(ctx)
	case "secrets":
		return handleSecrets(ctx)
	case "ingresses":
		return handleIngresses(ctx)
	case "pvcs":
		return handlePVCs(ctx)
	case "status":
		return handleStatus(ctx)
	case "events":
		return handleEvents(ctx)
	case "all":
		return handleAll(ctx)
	case "completion":
		return completion.HandleCompletion(ctx)
	case "apply":
		return handleApply(ctx)
	case "delete":
		return handleDelete(ctx)
	case "edit":
		return handleEdit(ctx)
	case "config":
		return handleConfig(ctx)
	case "metrics":
		return handleMetrics(ctx)
	case "copy":
		return handleCopy(ctx)
	case "explain":
		return handleExplain(ctx)
	case "update":
		return handleUpdate()
	case "version":
		help.PrintVersion()
		return nil
	case "help":
		help.PrintHelp(ctx.ResourceType)
		return nil
	default:
		return fmt.Errorf("unknown command: %s\nRun 'skube help' for usage", ctx.Command)
	}
}

func handleLogs(ctx *parser.Context) error {
	kubectlArgs := []string{"logs"}

	if ctx.AppName != "" {
		kubectlArgs = append(kubectlArgs, "-l", "app="+ctx.AppName)
		if ctx.Prefix {
			kubectlArgs = append(kubectlArgs, "--prefix=true")
		}
		fmt.Printf("%s📋 Fetching logs from app: %s%s\n", config.ColorCyan, ctx.AppName, config.ColorReset)
	} else if ctx.PodName != "" {
		kubectlArgs = append(kubectlArgs, ctx.PodName)
		fmt.Printf("%s📋 Fetching logs from pod: %s%s\n", config.ColorCyan, ctx.PodName, config.ColorReset)
	} else {
		return fmt.Errorf("need pod or app\nUsage: skube logs from app <name> in <namespace>\n       skube logs from pod <name> in <namespace>")
	}

	if ctx.Follow {
		kubectlArgs = append(kubectlArgs, "-f")
	}
	if ctx.TailLines > 0 {
		kubectlArgs = append(kubectlArgs, "--tail="+strconv.Itoa(ctx.TailLines))
	}
	if ctx.MaxLogRequests > 0 {
		kubectlArgs = append(kubectlArgs, "--max-log-requests="+strconv.Itoa(ctx.MaxLogRequests))
	}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	if ctx.SearchTerm != "" {
		return runKubectlPiped(kubectlArgs, []string{"--color=always", ctx.SearchTerm}, ctx.DryRun)
	}

	err := runKubectl(kubectlArgs, ctx.DryRun)
	if err != nil && !ctx.DryRun {
		// Check for common errors
		if strings.Contains(err.Error(), "ContainerCreating") || strings.Contains(err.Error(), "CrashLoopBackOff") {
			fmt.Printf("%s💡 Tip: The pod seems to be having trouble starting. Try describing it:%s\n", config.ColorYellow, config.ColorReset)
			fmt.Printf("   skube describe pod %s\n", ctx.PodName)
		}
	}
	return err
}

func handleShell(ctx *parser.Context) error {
	if ctx.PodName == "" {
		return fmt.Errorf("need pod name\nUsage: skube shell into pod <name> in <namespace>")
	}

	kubectlArgs := []string{"exec", "-it", ctx.PodName}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}
	kubectlArgs = append(kubectlArgs, "--", "sh")

	fmt.Printf("%s🐚 Opening shell in pod: %s%s\n", config.ColorCyan, ctx.PodName, config.ColorReset)
	err := runKubectl(kubectlArgs, ctx.DryRun)
	if err != nil && !ctx.DryRun {
		if strings.Contains(err.Error(), "not found") {
			fmt.Printf("%s💡 Tip: Double check the pod name. List pods with:%s\n", config.ColorYellow, config.ColorReset)
			fmt.Printf("   skube get pods\n")
		}
	}
	return err
}

func handleRestart(ctx *parser.Context) error {
	if ctx.PodName != "" {
		kubectlArgs := []string{"delete", "pod", ctx.PodName}
		if ctx.Namespace != "" {
			kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
		}
		fmt.Printf("%s🔄 Restarting pod: %s%s\n", config.ColorYellow, ctx.PodName, config.ColorReset)
		return runKubectl(kubectlArgs, ctx.DryRun)
	}

	if ctx.DeploymentName == "" {
		return fmt.Errorf("need deployment or pod name\nUsage: skube restart deployment <name> in <namespace>")
	}

	kubectlArgs := []string{"rollout", "restart", "deployment", ctx.DeploymentName}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s🔄 Restarting deployment: %s%s\n", config.ColorYellow, ctx.DeploymentName, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handlePods(ctx *parser.Context) error {
	kubectlArgs := []string{"get", "pods", "-o", "wide"}

	if ctx.AppName != "" {
		kubectlArgs = append(kubectlArgs, "-l", "app="+ctx.AppName)
		fmt.Printf("%s📦 Listing pods from app: %s%s\n", config.ColorCyan, ctx.AppName, config.ColorReset)
	} else {
		fmt.Printf("%s📦 Listing pods%s\n", config.ColorCyan, config.ColorReset)
	}

	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleScale(ctx *parser.Context) error {
	if ctx.DeploymentName == "" || ctx.Replicas == "" {
		return fmt.Errorf("need deployment and replicas\nUsage: skube scale deployment <name> to <N> in <namespace>")
	}

	kubectlArgs := []string{"scale", "deployment", ctx.DeploymentName, "--replicas=" + ctx.Replicas}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s⚖️  Scaling deployment %s to %s replicas%s\n", config.ColorYellow, ctx.DeploymentName, ctx.Replicas, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleRollback(ctx *parser.Context) error {
	if ctx.DeploymentName == "" {
		return fmt.Errorf("need deployment name\nUsage: skube rollback deployment <name> in <namespace>")
	}

	kubectlArgs := []string{"rollout", "undo", "deployment", ctx.DeploymentName}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s⏪ Rolling back deployment: %s%s\n", config.ColorYellow, ctx.DeploymentName, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handlePortForward(ctx *parser.Context) error {
	if ctx.ServiceName == "" || ctx.Port == "" {
		return fmt.Errorf("need service and port\nUsage: skube forward service <name> port <port> in <namespace>")
	}

	port := ctx.Port
	if !strings.Contains(port, ":") {
		port = port + ":" + port
	}

	kubectlArgs := []string{"port-forward", "service/" + ctx.ServiceName, port}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s🔌 Port forwarding service %s on %s%s\n", config.ColorCyan, ctx.ServiceName, port, config.ColorReset)
	err := runKubectl(kubectlArgs, ctx.DryRun)
	if err != nil && !ctx.DryRun {
		fmt.Printf("%s💡 Tip: Check if the service exists and exposes port %s:%s\n", config.ColorYellow, port, config.ColorReset)
		fmt.Printf("   skube get services\n")
	}
	return err
}

func handleDescribe(ctx *parser.Context) error {
	if ctx.ResourceType == "" || ctx.ResourceName == "" {
		return fmt.Errorf("need resource type and name\nUsage: skube describe pod <name> in <namespace>")
	}

	kubectlArgs := []string{"describe", ctx.ResourceType, ctx.ResourceName}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s🔍 Describing %s: %s%s\n", config.ColorCyan, ctx.ResourceType, ctx.ResourceName, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleServices(ctx *parser.Context) error {
	kubectlArgs := []string{"get", "services", "-o", "wide"}

	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s🌐 Listing services%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleDeployments(ctx *parser.Context) error {
	kubectlArgs := []string{"get", "deployments", "-o", "wide"}

	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s🚀 Listing deployments%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleStatus(ctx *parser.Context) error {
	fmt.Printf("%s📊 Cluster Status%s\n\n", config.ColorGreen, config.ColorReset)

	kubectlArgs := []string{"get", "all"}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleEvents(ctx *parser.Context) error {
	kubectlArgs := []string{"get", "events", "--sort-by=.lastTimestamp"}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s📅 Cluster Events%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleAll(ctx *parser.Context) error {
	kubectlArgs := []string{"get", "all", "-o", "wide"}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s📋 All Resources%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleNamespaces(ctx *parser.Context) error {
	fmt.Printf("%s📂 Listing namespaces%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl([]string{"get", "namespaces"}, ctx.DryRun)
}

func handleNodes(ctx *parser.Context) error {
	fmt.Printf("%s🖥️  Listing nodes%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl([]string{"get", "nodes", "-o", "wide"}, ctx.DryRun)
}

func handleConfigMaps(ctx *parser.Context) error {
	kubectlArgs := []string{"get", "configmaps"}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}
	fmt.Printf("%s📄 Listing configmaps%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleSecrets(ctx *parser.Context) error {
	kubectlArgs := []string{"get", "secrets"}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}
	fmt.Printf("%s🔒 Listing secrets%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleIngresses(ctx *parser.Context) error {
	kubectlArgs := []string{"get", "ingress"}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}
	fmt.Printf("%s🌐 Listing ingresses%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handlePVCs(ctx *parser.Context) error {
	kubectlArgs := []string{"get", "pvc"}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}
	fmt.Printf("%s💾 Listing persistent volume claims%s\n", config.ColorCyan, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleApply(ctx *parser.Context) error {
	if ctx.FilePath == "" {
		return fmt.Errorf("need file path\nUsage: skube apply file <filename>")
	}

	kubectlArgs := []string{"apply", "-f", ctx.FilePath}
	fmt.Printf("%s📝 Applying configuration from: %s%s\n", config.ColorYellow, ctx.FilePath, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleDelete(ctx *parser.Context) error {
	if ctx.ResourceType == "" || ctx.ResourceName == "" {
		return fmt.Errorf("need resource type and name\nUsage: skube delete <resource> <name> in <namespace>")
	}

	kubectlArgs := []string{"delete", ctx.ResourceType, ctx.ResourceName}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s🗑️  Deleting %s: %s%s\n", config.ColorRed, ctx.ResourceType, ctx.ResourceName, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleEdit(ctx *parser.Context) error {
	if ctx.ResourceType == "" || ctx.ResourceName == "" {
		return fmt.Errorf("need resource type and name\nUsage: skube edit <resource> <name> in <namespace>")
	}

	kubectlArgs := []string{"edit", ctx.ResourceType, ctx.ResourceName}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s✏️  Editing %s: %s%s\n", config.ColorYellow, ctx.ResourceType, ctx.ResourceName, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleConfig(ctx *parser.Context) error {
	if ctx.ResourceType == "view" {
		fmt.Printf("%s⚙️  Current Configuration%s\n", config.ColorCyan, config.ColorReset)
		return runKubectl([]string{"config", "view", "--minify"}, ctx.DryRun)
	} else if ctx.ResourceType == "context" {
		if ctx.ResourceName == "" {
			return fmt.Errorf("need context name\nUsage: skube use context <name>")
		}
		fmt.Printf("%s🔄 Switching to context: %s%s\n", config.ColorYellow, ctx.ResourceName, config.ColorReset)
		return runKubectl([]string{"config", "use-context", ctx.ResourceName}, ctx.DryRun)
	} else if ctx.ResourceType == "namespace" {
		if ctx.ResourceName == "" {
			return fmt.Errorf("need namespace name\nUsage: skube use namespace <name>")
		}
		fmt.Printf("%s🔄 Switching default namespace to: %s%s\n", config.ColorYellow, ctx.ResourceName, config.ColorReset)
		return runKubectl([]string{"config", "set-context", "--current", "--namespace=" + ctx.ResourceName}, ctx.DryRun)
	}
	return fmt.Errorf("unknown config command")
}

func handleMetrics(ctx *parser.Context) error {
	kubectlArgs := []string{"top"}

	if ctx.ResourceType == "nodes" {
		kubectlArgs = append(kubectlArgs, "nodes")
		fmt.Printf("%s📊 Node Metrics%s\n", config.ColorCyan, config.ColorReset)
	} else {
		// Default to pods
		kubectlArgs = append(kubectlArgs, "pods")
		if ctx.Namespace != "" {
			kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
		}
		fmt.Printf("%s📊 Pod Metrics%s\n", config.ColorCyan, config.ColorReset)
	}

	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleCopy(ctx *parser.Context) error {
	if ctx.SourcePath == "" || ctx.DestPath == "" {
		return fmt.Errorf("need source and destination\nUsage: skube copy file <src> to <dest>")
	}

	kubectlArgs := []string{"cp", ctx.SourcePath, ctx.DestPath}
	if ctx.Namespace != "" {
		kubectlArgs = append(kubectlArgs, "-n", ctx.Namespace)
	}

	fmt.Printf("%s📂 Copying %s to %s%s\n", config.ColorYellow, ctx.SourcePath, ctx.DestPath, config.ColorReset)
	return runKubectl(kubectlArgs, ctx.DryRun)
}

func handleExplain(ctx *parser.Context) error {
	if ctx.ResourceType == "" {
		return fmt.Errorf("need resource type\nUsage: skube explain <resource>")
	}

	fmt.Printf("%s📖 Explaining %s%s\n", config.ColorCyan, ctx.ResourceType, config.ColorReset)
	return runKubectl([]string{"explain", ctx.ResourceType}, ctx.DryRun)
}

func handleUpdate() error {
	fmt.Printf("%s🔄 Updating skube...%s\n", config.ColorCyan, config.ColorReset)

	// Check if go is installed
	_, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go not found. Please download the latest release from:\n  https://github.com/geminal/skube/releases")
	}

	fmt.Printf("%sRunning: go install github.com/geminal/skube/cmd/skube@latest%s\n", config.ColorYellow, config.ColorReset)

	cmd := execCommand("go", "install", "github.com/geminal/skube/cmd/skube@latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("update failed: %v", err)
	}

	fmt.Printf("\n%s✅ Update complete!%s\n", config.ColorGreen, config.ColorReset)
	return nil
}

func runKubectl(args []string, dryRun bool) error {
	if dryRun {
		fmt.Printf("%s📋 DRY RUN: Would execute:%s\n", config.ColorYellow, config.ColorReset)
		fmt.Printf("kubectl %s\n", strings.Join(args, " "))
		return nil
	}

	// Interactive or streaming commands need direct IO
	isInteractive := false
	if len(args) > 0 {
		if args[0] == "exec" || args[0] == "edit" || args[0] == "run" || args[0] == "attach" || args[0] == "port-forward" {
			isInteractive = true
		}
	}
	for _, arg := range args {
		if arg == "-f" || arg == "--follow" || arg == "-w" || arg == "--watch" {
			isInteractive = true
		}
	}

	// For interactive commands or commands that might trigger OIDC auth, use direct IO
	if isInteractive {
		cmd := execCommand("kubectl", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	// For non-interactive commands, try with a timeout first to detect hanging
	// If it might need auth, we'll switch to interactive mode
	cmd := execCommand("kubectl", args...)
	cmd.Stdin = os.Stdin // Always pass stdin for potential OIDC auth

	// Use pipes for stdout/stderr to capture and display
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Read and print output as it comes
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				fmt.Print(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				fmt.Fprint(os.Stderr, string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	err = cmd.Wait()

	if err != nil {
		return err
	}

	return nil
}

func runKubectlPiped(kubectlArgs []string, grepArgs []string, dryRun bool) error {
	if dryRun {
		fmt.Printf("%s📋 DRY RUN: Would execute:%s\n", config.ColorYellow, config.ColorReset)
		fmt.Printf("kubectl %s | grep %s\n", strings.Join(kubectlArgs, " "), strings.Join(grepArgs, " "))
		return nil
	}
	kubectlCmd := execCommand("kubectl", kubectlArgs...)
	grepCmd := execCommand("grep", grepArgs...)

	pipe, err := kubectlCmd.StdoutPipe()
	if err != nil {
		return err
	}

	grepCmd.Stdin = pipe
	grepCmd.Stdout = os.Stdout
	grepCmd.Stderr = os.Stderr
	kubectlCmd.Stderr = os.Stderr

	if err := grepCmd.Start(); err != nil {
		return err
	}

	if err := kubectlCmd.Start(); err != nil {
		return err
	}

	if err := kubectlCmd.Wait(); err != nil {
		return err
	}

	return grepCmd.Wait()
}

func sanitizeInput(input string) string {
	if input == "" {
		return ""
	}

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Block dangerous shell characters that could be used for injection
	// Even though exec.Command passes args safely, we want defense in depth
	dangerousChars := []string{
		";",    // Command separator
		"&",    // Background/AND operator
		"|",    // Pipe operator
		"$",    // Variable expansion
		"`",    // Command substitution
		"\\",   // Escape character (in some contexts)
		"\n",   // Newline
		"\r",   // Carriage return
		"\x00", // Null byte
	}

	for _, char := range dangerousChars {
		if strings.Contains(input, char) {
			// Return empty string if dangerous characters detected
			// This is safer than trying to escape, as it prevents any potential bypass
			return ""
		}
	}

	// Check for command substitution patterns
	if strings.Contains(input, "$(") || strings.Contains(input, "${") {
		return ""
	}

	// Prevent flag injection: if input starts with -, prepend with ./
	// This ensures it's treated as a path/argument, not a flag
	if strings.HasPrefix(input, "-") {
		// For file paths, make it explicit
		// For search terms, the caller should handle this with -- separator
		input = "./" + input
	}

	// Limit length to prevent DoS
	maxLength := 1024
	if len(input) > maxLength {
		input = input[:maxLength]
	}

	return input
}
