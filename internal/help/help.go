package help

import (
	"fmt"

	"github.com/geminal/skube/internal/config"
)

func PrintVersion() {
	fmt.Printf("skube version %s\n", "v1.0.0")
}

func PrintHelp() {
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
  skube use context %s<name>%s
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
		config.ColorCyan, config.ColorReset,
		config.ColorCyan, config.ColorReset,
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
		config.ColorBlue, config.ColorReset, // use context
		config.ColorBlue, config.ColorReset, // use namespace
	)

	fmt.Print(help)
}
