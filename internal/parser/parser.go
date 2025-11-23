package parser

import (
	"strconv"
	"strings"
)

type Context struct {
	Command        string
	Namespace      string
	AppName        string
	PodName        string
	ServiceName    string
	DeploymentName string
	ResourceType   string
	ResourceName   string
	Port           string
	Replicas       string
	Follow         bool
	Prefix         bool
	DryRun         bool
	SearchTerm     string
	TailLines      int
	MaxLogRequests int
}

func ParseNaturalLanguage(args []string) *Context {
	ctx := &Context{}
	input := strings.Join(args, " ")
	words := args

	for i := 0; i < len(words); i++ {
		word := words[i]

		switch word {
		case "completion":
			ctx.Command = "completion"
			if i+1 < len(words) {
				ctx.ResourceType = words[i+1]
			}
			return ctx

		case "update":
			ctx.Command = "update"
			return ctx

		case "help", "-h", "--help":
			ctx.Command = "help"
			return ctx

		case "--dry-run":
			ctx.DryRun = true

		case "get":
			if i+1 < len(words) {
				nextWord := words[i+1]
				if nextWord == "namespaces" || nextWord == "ns" {
					ctx.Command = "namespaces"
				} else if nextWord == "pods" || nextWord == "pod" {
					ctx.Command = "pods"
				} else if nextWord == "deployments" || nextWord == "deploy" {
					ctx.Command = "deployments"
				} else if nextWord == "services" || nextWord == "svc" {
					ctx.Command = "services"
				} else if nextWord == "all" {
					ctx.Command = "all"
				} else if nextWord == "last" && i+2 < len(words) {
					if lines, err := strconv.Atoi(words[i+2]); err == nil {
						ctx.TailLines = lines
					}
				}
			}

		case "logs":
			ctx.Command = "logs"

		case "shell", "exec":
			ctx.Command = "shell"

		case "restart":
			ctx.Command = "restart"

		case "scale":
			ctx.Command = "scale"

		case "rollback":
			ctx.Command = "rollback"

		case "forward":
			ctx.Command = "forward"

		case "describe":
			ctx.Command = "describe"

		case "show":
			if i+1 < len(words) && words[i+1] == "status" {
				ctx.Command = "status"
			} else if i+1 < len(words) && words[i+1] == "events" {
				ctx.Command = "events"
			}

		case "status":
			if ctx.Command == "" {
				ctx.Command = "status"
			}

		case "events":
			if ctx.Command == "" {
				ctx.Command = "events"
			}

		case "namespaces", "ns":
			if ctx.Command == "" {
				ctx.Command = "namespaces"
			}

		case "pods", "pod":
			if ctx.Command == "" {
				ctx.Command = "pods"
			}

		case "deployments", "deploy":
			if ctx.Command == "" {
				ctx.Command = "deployments"
			}

		case "services", "svc":
			if ctx.Command == "" {
				ctx.Command = "services"
			}

		case "of":
			// "of" keyword implies app context: "logs of myapp in qa"
			if i+1 < len(words) {
				ctx.AppName = words[i+1]
				i++
			}

		case "from", "in", "into":
			if i+1 < len(words) {
				nextWord := words[i+1]
				if nextWord == "pod" && i+2 < len(words) {
					ctx.PodName = words[i+2]
					i += 2
				} else if nextWord == "deployment" && i+2 < len(words) {
					ctx.DeploymentName = words[i+2]
					i += 2
				} else if nextWord == "service" && i+2 < len(words) {
					ctx.ServiceName = words[i+2]
					i += 2
				} else if nextWord == "namespace" && i+2 < len(words) {
					ctx.Namespace = words[i+2]
					i += 2
				} else if nextWord != "pod" && nextWord != "deployment" && nextWord != "service" {
					ctx.Namespace = nextWord
					i++
				}
			}

		case "app":
			if i+1 < len(words) && ctx.AppName == "" {
				ctx.AppName = words[i+1]
				i++
			}

		case "deployment":
			if i+1 < len(words) && ctx.DeploymentName == "" {
				ctx.DeploymentName = words[i+1]
				i++
			}

		case "service":
			if i+1 < len(words) && ctx.ServiceName == "" {
				ctx.ServiceName = words[i+1]
				i++
			}

		case "namespace":
			if i+1 < len(words) && ctx.Namespace == "" {
				ctx.Namespace = words[i+1]
				i++
			}

		case "to":
			if i+1 < len(words) {
				ctx.Replicas = words[i+1]
				i++
			}

		case "port":
			if i+1 < len(words) {
				ctx.Port = words[i+1]
				i++
			}

		case "follow", "-f":
			ctx.Follow = true

		case "prefix", "prefixes", "with":
			if i+1 < len(words) && words[i+1] == "prefix" {
				ctx.Prefix = true
				i++
			} else if word == "prefix" || word == "prefixes" {
				ctx.Prefix = true
			}

		case "search", "find", "filter", "grep":
			if i+1 < len(words) {
				ctx.SearchTerm = strings.Trim(words[i+1], `"'`)
				i++
			}

		case "max":
			// "max 30" or "max log requests 30" for --max-log-requests
			if i+1 < len(words) {
				nextWord := words[i+1]
				if nextWord == "log" && i+3 < len(words) && words[i+2] == "requests" {
					if maxReqs, err := strconv.Atoi(words[i+3]); err == nil {
						ctx.MaxLogRequests = maxReqs
						i += 3
					}
				} else if maxReqs, err := strconv.Atoi(nextWord); err == nil {
					ctx.MaxLogRequests = maxReqs
					i++
				}
			}

		case "-n", "--namespace":
			if i+1 < len(words) {
				ctx.Namespace = words[i+1]
				i++
			}

		default:
			if ctx.PodName == "" && ctx.DeploymentName == "" && ctx.ServiceName == "" && ctx.AppName == "" {
				if ctx.Command == "logs" || ctx.Command == "shell" || ctx.Command == "restart" {
					ctx.PodName = word
				} else if ctx.Command == "scale" || ctx.Command == "rollback" {
					ctx.DeploymentName = word
				} else if ctx.Command == "forward" {
					ctx.ServiceName = word
				} else if ctx.Command == "describe" {
					if ctx.ResourceType == "" {
						ctx.ResourceType = word
					} else if ctx.ResourceName == "" {
						ctx.ResourceName = word
					}
				}
			}
			if strings.Contains(word, ":") && ctx.Port == "" {
				ctx.Port = word
			}
			if _, err := strconv.Atoi(word); err == nil && ctx.Port == "" && ctx.Replicas == "" {
				if strings.Contains(input, "scale") {
					ctx.Replicas = word
				} else if strings.Contains(input, "forward") || strings.Contains(input, "port") {
					ctx.Port = word
				}
			}
		}
	}

	return ctx
}
