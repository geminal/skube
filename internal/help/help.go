package help

import (
	"fmt"
	"runtime/debug"

	"github.com/geminal/skube/internal/config"
)

// Version is set at build time using -ldflags
var Version = "dev"

func PrintVersion() {
	version := Version

	// Try to get version from build info (works with go install)
	if info, ok := debug.ReadBuildInfo(); ok && version == "dev" {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
	}

	fmt.Printf("skube version %s\n", version)
}

var commandHelp = map[string]string{
	"get": `Usage: skube get <resource> [in <namespace>]

List resources in the cluster.

Examples:
  skube get pods
  skube get services in production
  skube get deployments`,

	"logs": `Usage: skube logs from <pod|app> <name> [in <namespace>] [follow] [search "term"]

View logs from a pod or application.

Options:
  follow        Stream logs in real-time
  prefix        Add pod name prefix to lines
  search "term" Filter logs by keyword
  get last N    Show only the last N lines

Examples:
  skube logs from app myapp
  skube logs from pod backend-123 in prod follow`,

	"shell": `Usage: skube shell into pod <name> [in <namespace>]

Open an interactive shell (/bin/sh) in a running pod.

Examples:
  skube shell into pod backend-123
  skube in production shell into pod database-0`,

	"restart": `Usage: skube restart <deployment|pod> <name> [in <namespace>]

Restart a resource. For deployments, it performs a rollout restart. For pods, it deletes the pod.

Examples:
  skube restart deployment backend
  skube restart pod worker-123`,

	"scale": `Usage: skube scale deployment <name> to <N> [in <namespace>]

Scale a deployment to a specific number of replicas.

Examples:
  skube scale deployment backend to 5
  skube scale deployment worker to 0 in staging`,

	"forward": `Usage: skube forward service <name> port <port> [in <namespace>]

Forward a local port to a service in the cluster.

Examples:
  skube forward service web port 8080
  skube forward service db port 5432:5432 in prod`,
}

func PrintHelp(args ...string) {
	if len(args) > 0 && args[0] != "" {
		cmd := args[0]
		if helpText, ok := commandHelp[cmd]; ok {
			fmt.Println(helpText)
			return
		}
		fmt.Printf("No specific help for command: %s\n\n", cmd)
	}

	help := fmt.Sprintf(`%sskube%s - Talk to Kubernetes in plain English

%sUSAGE:%s
  skube %s<command>%s %s<resource>%s %sfrom|in%s %s<name>%s %s<namespace>%s
  skube %sin%s %s<namespace>%s %s<command>%s ...

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
  %sconfig%s      Manage configuration (contexts/namespaces)
  %scopy%s        Copy files to/from pods
  %sexplain%s     Documentation for resources
  %scompletion%s  Generate shell completion script (zsh, bash)
  %supdate%s      Update skube to latest version
  %ssetup-ai%s    Configure AI features (requires Docker or Ollama)
  %sswitch-ai%s   Switch between AI providers (Ollama/OpenAI)
  %sconfig-ai%s   Import AI config from JSON file
  %smodel%s       Show current AI model and provider configuration
  %shelp%s        Show help message (try: skube help logs)

%sRESOURCES:%s
  %snamespaces%s    Kubernetes namespaces (environments)
  %spods%s          Running pod instances
  %sdeployments%s   Deployment configurations
  %sservices%s      Service endpoints
  %snodes%s         Cluster nodes
  %sconfigmaps%s    Configuration data (cm)
  %ssecrets%s       Sensitive data
  %singresses%s     Ingress resources (ing)
  %spvcs%s          PersistentVolumeClaims (pvc)
  %sapp%s           Filter by application label
  %spod%s           Specific pod name

%sOPTIONS:%s
  %sfollow%s        Tail logs in real-time
  %sprefix%s        Show pod names in multi-pod logs
  %ssearch%s        Filter logs by keyword
  %sfind%s          Same as search
  %sget last N%s    Show last N lines of logs
  %s--dry-run%s     Show kubectl command without executing
  %s--ai%s          Use AI to parse natural language (run 'setup-ai' first)

%sEXAMPLES:%s
  %s# Investigation%s
  skube get namespaces
  skube in %s<namespace>%s get pods
  skube in %s<namespace>%s logs from app %s<app-name>%s
  skube logs from app %s<app-name>%s in %s<namespace>%s follow with prefix
  skube logs from pod %s<pod-name>%s get last 100 in %s<namespace>%s
  skube logs from pod %s<pod-name>%s search "%serror%s" in %s<namespace>%s
  skube show metrics pods in %s<namespace>%s
  skube explain pod

  %s# Operations%s
  skube in %s<namespace>%s shell into pod %s<pod-name>%s
  skube in %s<namespace>%s restart deployment %s<name>%s
  skube scale deployment %s<name>%s to %s<N>%s in %s<namespace>%s
  skube forward service %s<name>%s port %s<port>%s in %s<namespace>%s
  skube apply file %s<filename>%s
  skube delete pod %s<name>%s in %s<namespace>%s
  skube copy file %s<src>%s to %s<dest>%s in %s<namespace>%s

  %s# Context Management%s
  skube show context
  skube list contexts
  skube use context %s<name>%s
  skube switch context %s<name>%s
  skube use namespace %s<name>%s
`,
		config.ColorGreen, config.ColorReset,
		config.ColorYellow, config.ColorReset,
		config.ColorCyan, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorYellow, config.ColorReset, config.ColorBlue, config.ColorReset,
		config.ColorYellow, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorYellow, config.ColorReset,
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
		config.ColorCyan, config.ColorReset, // completion
		config.ColorCyan, config.ColorReset, // update
		config.ColorCyan, config.ColorReset, // setup-ai
		config.ColorCyan, config.ColorReset, // switch-ai
		config.ColorCyan, config.ColorReset, // config-ai
		config.ColorCyan, config.ColorReset, // model
		config.ColorCyan, config.ColorReset, // help
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
		config.ColorYellow, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset,
		config.ColorBlue, config.ColorReset, // --ai

		config.ColorYellow, config.ColorReset, // EXAMPLES header
		config.ColorYellow, config.ColorReset, // Investigation header

		config.ColorBlue, config.ColorReset, // in <namespace> get pods
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // in <namespace> logs from app <app-name>
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // logs from app <app-name> in <namespace>
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // logs from pod <pod-name> ... in <namespace>
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // logs from pod ... search "error" in <namespace>
		config.ColorBlue, config.ColorReset, // Extra args needed for alignment
		config.ColorBlue, config.ColorReset, // show metrics ... in <namespace>

		config.ColorYellow, config.ColorReset, // Operations header

		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // shell
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // restart
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // scale
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // forward
		config.ColorBlue, config.ColorReset, // apply
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // delete
		config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, config.ColorBlue, config.ColorReset, // copy

		config.ColorYellow, config.ColorReset, // Context Management header
		// show context (no params)
		// list contexts (no params)
		config.ColorBlue, config.ColorReset, // use context <name>
		config.ColorBlue, config.ColorReset, // switch context <name>
		config.ColorBlue, config.ColorReset, // use namespace <name>
	)

	fmt.Print(help)
}
