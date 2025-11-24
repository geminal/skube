package executor

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/geminal/skube/internal/config"
	"github.com/geminal/skube/internal/parser"
)

func ExecuteCommand(ctx *parser.Context) error {
	if ctx.Command == "" {
		PrintHelp()
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
	case "status":
		return handleStatus(ctx)
	case "events":
		return handleEvents(ctx)
	case "all":
		return handleAll(ctx)
	case "completion":
		return handleCompletion(ctx)
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
	case "help":
		PrintHelp()
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

	return runKubectl(kubectlArgs, ctx.DryRun)
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
	return runKubectl(kubectlArgs, ctx.DryRun)
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
	return runKubectl(kubectlArgs, ctx.DryRun)
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

func handleCompletion(ctx *parser.Context) error {
	shell := ctx.ResourceType
	if shell == "" {
		return fmt.Errorf("please specify shell type\nUsage: skube completion <zsh|bash>")
	}

	switch shell {
	case "zsh":
		fmt.Print(getZshCompletion())
		return nil
	case "bash":
		fmt.Print(getBashCompletion())
		return nil
	default:
		return fmt.Errorf("unsupported shell: %s\nSupported shells: zsh, bash", shell)
	}
}

func handleUpdate() error {
	fmt.Printf("%s🔄 Updating skube...%s\n", config.ColorCyan, config.ColorReset)

	// Check if go is installed
	_, err := exec.LookPath("go")
	if err != nil {
		return fmt.Errorf("go not found. Please download the latest release from:\n  https://github.com/geminal/skube/releases")
	}

	fmt.Printf("%sRunning: go install github.com/geminal/skube/cmd/skube@latest%s\n", config.ColorYellow, config.ColorReset)

	cmd := exec.Command("go", "install", "github.com/geminal/skube/cmd/skube@latest")
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
	cmd := exec.Command("kubectl", args...)

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

	if isInteractive {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	// Capture output for analysis
	output, err := cmd.CombinedOutput()
	fmt.Print(string(output))

	if err != nil {
		return err
	}

	// Check for empty results on success
	outStr := string(output)
	if strings.TrimSpace(outStr) == "" || strings.Contains(outStr, "No resources found") {
		fmt.Printf("\n%s⚠️  Tip: No resources found. Please check if the app label or namespace is correct.%s\n", config.ColorYellow, config.ColorReset)
	}

	return nil
}

func runKubectlPiped(kubectlArgs []string, grepArgs []string, dryRun bool) error {
	if dryRun {
		fmt.Printf("%s📋 DRY RUN: Would execute:%s\n", config.ColorYellow, config.ColorReset)
		fmt.Printf("kubectl %s | grep %s\n", strings.Join(kubectlArgs, " "), strings.Join(grepArgs, " "))
		return nil
	}
	kubectlCmd := exec.Command("kubectl", kubectlArgs...)
	grepCmd := exec.Command("grep", grepArgs...)

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

func getZshCompletion() string {
	return `#compdef skube

# Initialize completion system if not already done
if [[ -z "$_comps" ]]; then
    autoload -Uz compinit && compinit
fi

_skube() {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    local -a commands
    commands=(
        'get:Get resources (namespaces, pods, deployments, services)'
        'logs:View logs from pods or apps'
        'shell:Open shell in a pod'
        'restart:Restart a pod or deployment'
        'scale:Scale a deployment'
        'rollback:Rollback a deployment'
        'forward:Port forward to a service'
        'describe:Describe a resource'
        'show:Show status, events, or metrics'
        'apply:Apply configuration from file'
        'delete:Delete resources'
        'edit:Edit resources'
        'config:Manage configuration'
        'copy:Copy files'
        'explain:Resource documentation'
        'completion:Generate completion script'
        'update:Update skube to latest version'
        'help:Show help message'
    )

    local -a resources
    resources=(
        'namespaces:List all namespaces (start here!)'
        'pods:List all pods'
        'deployments:List deployments'
        'services:List services'
        'ns:Shorthand for namespaces'
        'pod:Shorthand for pods'
        'deploy:Shorthand for deployments'
        'svc:Shorthand for services'
    )

    local -a common_namespaces
    common_namespaces=(
        'production'
        'staging'
        'qa'
        'dev'
        'prod'
        'test'
    )

    local -a log_options
    log_options=(
        'follow:Follow logs in real-time'
        'prefix:Show pod names in logs'
        'max:Set max concurrent log streams'
        'search:Search for term in logs'
        'find:Find term in logs'
    )

    # Helper function to check if word exists in previous words
    _contains_word() {
        local word="$1"
        shift
        for w in "$@"; do
            [[ "$w" == "$word" ]] && return 0
        done
        return 1
    }

    # Get previous word (word before cursor)
    local prev="${words[CURRENT-1]}"

    # First argument - show commands
    if [[ $CURRENT -eq 2 ]]; then
        _describe 'command' commands
        return
    fi

    # Context-aware completion based on command and previous words
    case "${words[2]}" in
        get)
            case "$prev" in
                get)
                    _describe 'resource' resources
                    ;;
                pods|pod|deployments|deploy|services|svc)
                    compadd of in
                    _values 'namespace' "${common_namespaces[@]}"
                    ;;
                of)
                    # After "of", suggest app names (user types them)
                    ;;
                in)
                    _values 'namespace' "${common_namespaces[@]}"
                    ;;
                namespaces|ns)
                    ;;
                *)
                    if (_contains_word "pods" "${words[@]}" || _contains_word "pod" "${words[@]}" || \
                        _contains_word "deployments" "${words[@]}" || _contains_word "deploy" "${words[@]}" || \
                        _contains_word "services" "${words[@]}" || _contains_word "svc" "${words[@]}") && \
                       _contains_word "of" "${words[@]}"; then
                        compadd in
                        _values 'namespace' "${common_namespaces[@]}"
                    fi
                    ;;
            esac
            ;;

        logs)
            case "$prev" in
                logs)
                    compadd of from
                    ;;
                of)
                    # After "of", user types app name
                    ;;
                from)
                    compadd pod
                    ;;
                pod)
                    ;;
                *)
                    if _contains_word "of" "${words[@]}" || _contains_word "pod" "${words[@]}"; then
                        if ! _contains_word "in" "${words[@]}"; then
                            compadd in
                        fi
                        _describe 'log options' log_options
                        _values 'namespace' "${common_namespaces[@]}"
                    else
                        compadd of from
                    fi
                    ;;
            esac
            ;;

        shell)
            case "$prev" in
                shell)
                    compadd into
                    ;;
                into)
                    compadd pod
                    ;;
                pod)
                    ;;
                *)
                    if _contains_word "pod" "${words[@]}"; then
                        compadd in
                        _values 'namespace' "${common_namespaces[@]}"
                    fi
                    ;;
            esac
            ;;

        restart)
            case "$prev" in
                restart)
                    compadd deployment pod
                    ;;
                deployment|pod)
                    ;;
                *)
                    if _contains_word "deployment" "${words[@]}" || _contains_word "pod" "${words[@]}"; then
                        compadd in
                        _values 'namespace' "${common_namespaces[@]}"
                    fi
                    ;;
            esac
            ;;

        scale)
            case "$prev" in
                scale)
                    compadd deployment
                    ;;
                deployment)
                    ;;
                to)
                    ;;
                *)
                    if _contains_word "deployment" "${words[@]}"; then
                        if ! _contains_word "to" "${words[@]}"; then
                            compadd to
                        elif ! _contains_word "in" "${words[@]}"; then
                            compadd in
                        else
                            _values 'namespace' "${common_namespaces[@]}"
                        fi
                    fi
                    ;;
            esac
            ;;

        rollback)
            case "$prev" in
                rollback)
                    compadd deployment
                    ;;
                deployment)
                    ;;
                *)
                    if _contains_word "deployment" "${words[@]}"; then
                        compadd in
                        _values 'namespace' "${common_namespaces[@]}"
                    fi
                    ;;
            esac
            ;;

        forward)
            case "$prev" in
                forward)
                    compadd service
                    ;;
                service)
                    ;;
                port)
                    ;;
                *)
                    if _contains_word "service" "${words[@]}"; then
                        if ! _contains_word "port" "${words[@]}"; then
                            compadd port
                        elif ! _contains_word "in" "${words[@]}"; then
                            compadd in
                        else
                            _values 'namespace' "${common_namespaces[@]}"
                        fi
                    fi
                    ;;
            esac
            ;;

        describe)
            if [[ $CURRENT -eq 3 ]]; then
                compadd pod deployment service namespace
            elif [[ $CURRENT -eq 4 ]]; then
                :
            else
                if ! _contains_word "in" "${words[@]}"; then
                    compadd in
                fi
                _values 'namespace' "${common_namespaces[@]}"
            fi
            ;;

        show)
            if [[ $CURRENT -eq 3 ]]; then
                compadd status events
            else
                compadd in
                _values 'namespace' "${common_namespaces[@]}"
            fi
            ;;

        completion)
            if [[ $CURRENT -eq 3 ]]; then
                compadd zsh bash
            fi
            ;;

        *)
            ;;
    esac
}

# Register the completion function
compdef _skube skube
`
}

func getBashCompletion() string {
	return `#!/bin/bash

_skube_completions() {
    local cur prev words cword
    _init_completion || return

    local commands="get logs shell restart scale rollback forward describe show completion update help apply delete edit config metrics copy explain"
    local keywords="of from in into pod deployment service namespace to port follow prefix with search find get last file context"
    local resources="namespaces pods deployments services status events nodes"
    local namespaces="production staging qa dev prod test"

    case "${prev}" in
        skube)
            COMPREPLY=($(compgen -W "${commands}" -- "${cur}"))
            return 0
            ;;
        get|delete|edit|explain)
            COMPREPLY=($(compgen -W "${resources}" -- "${cur}"))
            return 0
            ;;
        logs|shell|restart|describe|scale|rollback|forward|apply|copy|metrics)
            COMPREPLY=($(compgen -W "${keywords}" -- "${cur}"))
            return 0
            ;;
        show)
            COMPREPLY=($(compgen -W "status events metrics config" -- "${cur}"))
            return 0
            ;;
        config)
             COMPREPLY=($(compgen -W "use view" -- "${cur}"))
             return 0
             ;;
        use)
             COMPREPLY=($(compgen -W "context namespace" -- "${cur}"))
             return 0
             ;;
        of|from|in|into)
            COMPREPLY=($(compgen -W "pod deployment service namespace ${namespaces}" -- "${cur}"))
            return 0
            ;;
        pod|deployment|service)
            return 0
            ;;
        namespace)
            COMPREPLY=($(compgen -W "${namespaces}" -- "${cur}"))
            return 0
            ;;
        completion)
            COMPREPLY=($(compgen -W "zsh bash" -- "${cur}"))
            return 0
            ;;
        *)
            COMPREPLY=($(compgen -W "${keywords} ${namespaces}" -- "${cur}"))
            return 0
            ;;
    esac
}

complete -F _skube_completions skube
`
}

func PrintHelp() {
	help := fmt.Sprintf(`%sskube%s - Talk to Kubernetes in plain English

%sUSAGE:%s
  skube %s<command>%s %s<resource>%s %sfrom|in%s %s<name>%s %s<namespace>%s

%sCOMMANDS:%s
  %sget%s         List resources (namespaces, pods, deployments, services)
  %slogs%s        View and search logs from pods or apps
  %sshell%s       Open interactive shell in a pod
  %srestart%s     Restart pods or deployments
  %sscale%s       Scale deployment replicas
  %srollback%s    Rollback deployment to previous version
  %sforward%s     Port forward to a service
  %sdescribe%s    Show detailed resource information
  %sshow%s        Display cluster status, events, or metrics
  %sapply%s       Apply configuration from file
  %sdelete%s      Delete resources
  %sedit%s        Edit resources
  %sconfig%s      Manage configuration (context/namespace)
  %scopy%s        Copy files to/from pods
  %sexplain%s     Documentation for resources
  %scompletion%s  Generate shell completion script (zsh, bash)
  %supdate%s      Update skube to latest version

%sRESOURCES:%s
  %snamespaces%s    Kubernetes namespaces (environments)
  %spods%s          Running pod instances
  %sdeployments%s   Deployment configurations
  %sservices%s      Service endpoints
  %sapp%s           Filter by application label
  %spod%s           Specific pod name

%sOPTIONS:%s
  %sfollow%s        Tail logs in real-time
  %sprefix%s        Show pod names in multi-pod logs
  %ssearch%s        Filter logs by keyword
  %sfind%s          Same as search
  %sget last N%s    Show last N lines of logs
  %s--dry-run%s     Show kubectl command without executing

%sEXAMPLES:%s
  %s# Investigation%s
  skube get namespaces
  skube get pods from %s<namespace>%s
  skube get pods from app %s<app-name>%s in %s<namespace>%s
  skube logs from app %s<app-name>%s in %s<namespace>%s follow with prefix
  skube logs from pod %s<pod-name>%s get last 100 in %s<namespace>%s
  skube logs from pod %s<pod-name>%s search "%serror%s" in %s<namespace>%s
  skube show metrics pods in %s<namespace>%s
  skube explain pod

  %s# Operations%s
  skube shell into pod %s<pod-name>%s in %s<namespace>%s
  skube restart deployment %s<name>%s in %s<namespace>%s
  skube scale deployment %s<name>%s to %s<N>%s in %s<namespace>%s
  skube forward service %s<name>%s port %s<port>%s in %s<namespace>%s
  skube apply file %s<filename>%s
  skube delete pod %s<name>%s in %s<namespace>%s
  skube copy file %s<src>%s to %s<dest>%s in %s<namespace>%s
  skube use context %s<name>%s
  skube use namespace %s<name>%s
`,
		config.ColorGreen, config.ColorReset,
		config.ColorYellow, config.ColorReset,
		config.ColorCyan, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorYellow, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorYellow, config.ColorReset,
		config.ColorYellow, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
		config.ColorYellow, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorYellow, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorYellow, config.ColorReset,
		config.ColorGreen, config.ColorReset,

		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorGreen, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset,
	)

	fmt.Print(help)
}
