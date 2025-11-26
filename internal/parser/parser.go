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
	FilePath       string
	SourcePath     string
	DestPath       string
}

func ParseNaturalLanguage(args []string) *Context {
	ctx := &Context{}
	input := strings.Join(args, " ")

	// Early namespace detection (namespace-first syntax)
	// Supports: "skube in production logs from app myapp"
	if len(args) > 1 && strings.ToLower(args[0]) == "in" {
		ctx.Namespace = args[1]
		args = args[2:] // Remove "in <namespace>" from args
		input = strings.Join(args, " ")
	}

	for i := 0; i < len(args); i++ {
		word := strings.ToLower(args[i])

		// Skip stop words
		if stopWords[word] {
			continue
		}

		// Try to parse as a command
		if parseCommand(word, args, &i, ctx) {
			continue
		}

		// Try to parse as a resource type
		if parseResource(word, args, &i, ctx) {
			continue
		}

		// Try to parse flags and modifiers
		if parseFlags(word, args, &i, ctx) {
			continue
		}

		// Try to parse prepositions and context
		if parsePrepositions(word, args, &i, ctx) {
			continue
		}

		// Default fallback logic
		parseDefault(word, input, ctx)
	}

	return ctx
}

var commandAliases = map[string]string{
	"completion": "completion",
	"update":     "update",
	"version":    "version", "-v": "version", "--version": "version",
	"help": "help", "-h": "help", "--help": "help",
	"apply": "apply", "create": "apply",
	"delete": "delete", "remove": "delete", "destroy": "delete",
	"edit": "edit", "change": "edit", "modify": "edit",
	"use": "config", "switch": "config", "config": "config",
	"copy": "copy", "cp": "copy",
	"explain": "explain", "what": "explain",
	"logs": "logs", "log": "logs", "monitor": "logs", "tail": "logs",
	"shell": "shell", "exec": "shell", "ssh": "shell", "connect": "shell",
	"restart": "restart", "reboot": "restart", "bounce": "restart",
	"scale": "scale", "resize": "scale",
	"rollback": "rollback", "undo": "rollback", "revert": "rollback",
	"forward": "forward", "port-forward": "forward", "tunnel": "forward",
	"describe": "describe", "inspect": "describe", "details": "describe",
	"status": "status", "health": "status",
	"events": "events", "history": "events",
	"get": "get", "list": "get", "show": "get", "fetch": "get", "give": "get", "check": "get",
}

var stopWords = map[string]bool{
	"the": true, "a": true, "an": true,
	"my": true, "our": true, "your": true,
	"please": true, "plz": true, "kindly": true,
	"me": true, "us": true,
	"for": true, "target": true,
	"resource": true, "resources": true, "object": true, "objects": true,
	"here": true, "now": true,
}

var resourceAliases = map[string]string{
	"namespaces": "namespace", "ns": "namespace", "namespace": "namespace",
	"pods": "pod", "pod": "pod",
	"deployments": "deployment", "deploy": "deployment", "deployment": "deployment",
	"services": "service", "svc": "service", "service": "service",
	"nodes": "node", "no": "node",
	"configmaps": "configmap", "cm": "configmap",
	"secrets":   "secret",
	"ingresses": "ingress", "ing": "ingress",
	"persistentvolumeclaims": "persistentvolumeclaim", "pvc": "persistentvolumeclaim",
}

var getCommandMap = map[string]string{
	"namespaces": "namespaces", "ns": "namespaces",
	"pods": "pods", "pod": "pods",
	"deployments": "deployments", "deploy": "deployments",
	"services": "services", "svc": "services",
	"nodes": "nodes", "no": "nodes",
	"configmaps": "configmaps", "cm": "configmaps",
	"secrets":   "secrets",
	"ingresses": "ingresses", "ing": "ingresses",
	"persistentvolumeclaims": "pvcs", "pvc": "pvcs",
	"all": "all",
}

func parseCommand(word string, args []string, index *int, ctx *Context) bool {
	i := *index

	// Special case for "get"
	if word == "get" {
		if i+1 < len(args) {
			nextWord := args[i+1]
			if cmd, ok := getCommandMap[nextWord]; ok {
				ctx.Command = cmd
				return true
			}
			if nextWord == "last" && i+2 < len(args) {
				if lines, err := strconv.Atoi(args[i+2]); err == nil {
					ctx.TailLines = lines
				}
			}
		}
		return true
	}

	// Special case for "check usage" -> metrics
	if word == "check" {
		if i+1 < len(args) && args[i+1] == "usage" {
			ctx.Command = "metrics"
			if i+2 < len(args) {
				ctx.ResourceType = args[i+2]
				*index += 2
			}
			return true
		}
	}

	// Special case for "show"
	if word == "show" {
		if i+1 < len(args) {
			sub := args[i+1]
			if sub == "status" {
				ctx.Command = "status"
				return true
			} else if sub == "events" {
				ctx.Command = "events"
				return true
			} else if sub == "config" {
				ctx.Command = "config"
				ctx.ResourceType = "view"
				return true
			} else if sub == "metrics" {
				ctx.Command = "metrics"
				if i+2 < len(args) {
					ctx.ResourceType = args[i+2]
					*index += 2
				}
				return true
			}
		}
		// If not a special show command, fall through to alias lookup (show -> get)
	}

	// Lookup in alias map
	if cmd, ok := commandAliases[word]; ok {
		ctx.Command = cmd

		// Handle command-specific arguments
		switch cmd {
		case "completion":
			if i+1 < len(args) {
				ctx.ResourceType = args[i+1]
			}
		case "apply":
			if i+1 < len(args) && (args[i+1] == "file" || args[i+1] == "-f") {
				if i+2 < len(args) {
					ctx.FilePath = args[i+2]
					*index += 2
				}
			}
		case "config":
			// "use context" or "use namespace"
			if i+1 < len(args) {
				if args[i+1] == "context" {
					ctx.ResourceType = "context"
					if i+2 < len(args) {
						ctx.ResourceName = args[i+2]
						*index += 2
					}
				} else if args[i+1] == "namespace" || args[i+1] == "ns" {
					ctx.ResourceType = "namespace"
					if i+2 < len(args) {
						ctx.ResourceName = args[i+2]
						*index += 2
					}
				}
			}
		case "copy":
			if i+1 < len(args) && args[i+1] == "file" {
				*index++
			}
		case "explain":
			if word == "what" && i+1 < len(args) && args[i+1] == "is" {
				*index++
			}
		}
		return true
	}

	return false
}

func parseResource(word string, args []string, index *int, ctx *Context) bool {
	i := *index

	// Check resource aliases
	if resType, ok := resourceAliases[word]; ok {
		// If command is empty OR command is generic "get", upgrade to specific command
		if cmd, isCmd := getCommandMap[word]; isCmd {
			if ctx.Command == "" || ctx.Command == "get" {
				ctx.Command = cmd
				return true
			}
		}

		if isResourceCommand(ctx.Command) {
			ctx.ResourceType = resType
			return true
		}

		// Type Correction: If we see "deployment" but have a PodName (likely from parseDefault),
		// and no DeploymentName, assume the PodName was actually the DeploymentName.
		if resType == "deployment" && ctx.PodName != "" && ctx.DeploymentName == "" {
			ctx.DeploymentName = ctx.PodName
			ctx.PodName = ""
			// Don't return true yet, we might still want to consume next word if it's a name?
			// But usually "backend deployment" -> backend is the name.
			// If "deployment backend" -> backend is next.
			// Let's continue to check next word just in case?
			// No, if we corrected the name, we are good for this word.
			return true
		}

		// Context setting (e.g. "deployment api")
		if resType == "deployment" && i+1 < len(args) && ctx.DeploymentName == "" {
			// Check if next word is a stop word or preposition, if so, don't consume it
			nextWord := strings.ToLower(args[i+1])
			if !stopWords[nextWord] && nextWord != "in" && nextWord != "from" && nextWord != "to" {
				ctx.DeploymentName = args[i+1]
				*index++
			}
			return true
		}
		if resType == "service" && i+1 < len(args) && ctx.ServiceName == "" {
			nextWord := strings.ToLower(args[i+1])
			if !stopWords[nextWord] && nextWord != "in" && nextWord != "from" && nextWord != "to" {
				ctx.ServiceName = args[i+1]
				*index++
			}
			return true
		}
		if resType == "namespace" && i+1 < len(args) && ctx.Namespace == "" {
			nextWord := strings.ToLower(args[i+1])
			if !stopWords[nextWord] && nextWord != "in" && nextWord != "from" && nextWord != "to" {
				ctx.Namespace = args[i+1]
				*index++
			}
			return true
		}

		return true
	}

	if word == "app" {
		if i+1 < len(args) && ctx.AppName == "" {
			ctx.AppName = args[i+1]
			*index++
		}
		return true
	}

	return false
}

func isResourceCommand(cmd string) bool {
	return cmd == "delete" || cmd == "edit" || cmd == "explain" || cmd == "describe"
}

func parseFlags(word string, args []string, index *int, ctx *Context) bool {
	i := *index
	switch word {
	case "--dry-run":
		ctx.DryRun = true
		return true

	case "to":
		if i+1 < len(args) {
			if ctx.Command == "copy" {
				ctx.DestPath = args[i+1]
				*index++
			} else {
				ctx.Replicas = args[i+1]
				*index++
			}
		}
		return true

	case "port":
		if i+1 < len(args) {
			ctx.Port = args[i+1]
			*index++
		}
		return true

	case "follow", "-f":
		ctx.Follow = true
		return true

	case "prefix", "prefixes", "with":
		if i+1 < len(args) && args[i+1] == "prefix" {
			ctx.Prefix = true
			*index++
		} else if word == "prefix" || word == "prefixes" {
			ctx.Prefix = true
		}
		return true

	case "search", "find", "filter", "grep":
		if i+1 < len(args) {
			ctx.SearchTerm = strings.Trim(args[i+1], `"'`)
			*index++
		}
		return true

	case "max":
		// "max 30" or "max log requests 30" for --max-log-requests
		if i+1 < len(args) {
			nextWord := args[i+1]
			if nextWord == "log" && i+3 < len(args) && args[i+2] == "requests" {
				if maxReqs, err := strconv.Atoi(args[i+3]); err == nil {
					ctx.MaxLogRequests = maxReqs
					*index += 3
				}
			} else if maxReqs, err := strconv.Atoi(nextWord); err == nil {
				ctx.MaxLogRequests = maxReqs
				*index++
			}
		}
		return true

	case "-n", "--namespace":
		if i+1 < len(args) {
			ctx.Namespace = args[i+1]
			*index++
		}
		return true
	}
	return false
}

func parsePrepositions(word string, args []string, index *int, ctx *Context) bool {
	i := *index
	switch word {
	case "of":
		// "of" keyword implies app context: "logs of myapp in qa"
		if i+1 < len(args) {
			if args[i+1] == "app" {
				if i+2 < len(args) {
					ctx.AppName = args[i+2]
					*index += 2
				}
			} else if args[i+1] == "pod" {
				if i+2 < len(args) {
					ctx.PodName = args[i+2]
					*index += 2
				}
			} else {
				ctx.AppName = args[i+1]
				*index++
			}
		}
		return true

	case "from", "in", "into":
		if i+1 < len(args) {
			nextWord := args[i+1]
			if nextWord == "pod" && i+2 < len(args) {
				ctx.PodName = args[i+2]
				*index += 2
			} else if nextWord == "deployment" && i+2 < len(args) {
				ctx.DeploymentName = args[i+2]
				*index += 2
			} else if nextWord == "service" && i+2 < len(args) {
				ctx.ServiceName = args[i+2]
				*index += 2
			} else if nextWord == "namespace" && i+2 < len(args) {
				ctx.Namespace = args[i+2]
				*index += 2
			} else if nextWord == "app" && i+2 < len(args) {
				ctx.AppName = args[i+2]
				*index += 2
			} else if nextWord == "file" && i+2 < len(args) {
				// for apply or copy
				if ctx.Command == "apply" {
					ctx.FilePath = args[i+2]
				} else if ctx.Command == "copy" {
					ctx.SourcePath = args[i+2]
				}
				*index += 2
			} else if nextWord != "pod" && nextWord != "deployment" && nextWord != "service" && nextWord != "file" && nextWord != "app" {
				ctx.Namespace = nextWord
				*index++
			}
		}
		return true
	}
	return false
}

func parseDefault(word string, input string, ctx *Context) {
	// If we have a command that lists resources, and namespace is empty, assume this word is the namespace
	// e.g. "skube pods qa" -> Command="pods", Namespace="qa"
	if ctx.Namespace == "" && (ctx.Command == "pods" || ctx.Command == "deployments" ||
		ctx.Command == "services" || ctx.Command == "nodes" || ctx.Command == "configmaps" ||
		ctx.Command == "secrets" || ctx.Command == "ingresses" || ctx.Command == "pvcs" ||
		ctx.Command == "events" || ctx.Command == "status" || ctx.Command == "all") {
		// Ensure it's not a flag or modifier we missed (though parseFlags should have caught it)
		if !strings.HasPrefix(word, "-") {
			ctx.Namespace = word
			return
		}
	}

	// If we have a specific resource selected (Pod, App, Deployment, Service), and namespace is empty,
	// assume this word is the namespace.
	// e.g. "skube logs myapp qa" -> Command="logs", AppName="myapp", Namespace="qa"
	if ctx.Namespace == "" && (ctx.PodName != "" || ctx.AppName != "" || ctx.DeploymentName != "" || ctx.ServiceName != "" || ctx.ResourceName != "") {
		if !strings.HasPrefix(word, "-") {
			ctx.Namespace = word
			return
		}
	}

	// Default resource name inference
	if ctx.PodName == "" && ctx.DeploymentName == "" && ctx.ServiceName == "" && ctx.AppName == "" && ctx.ResourceName == "" {
		if ctx.Command == "logs" || ctx.Command == "shell" || ctx.Command == "restart" {
			ctx.PodName = word
		} else if ctx.Command == "scale" || ctx.Command == "rollback" {
			ctx.DeploymentName = word
		} else if ctx.Command == "forward" {
			ctx.ServiceName = word
		} else if isResourceCommand(ctx.Command) {
			if ctx.ResourceType == "" {
				ctx.ResourceType = word
			} else if ctx.ResourceName == "" {
				ctx.ResourceName = word
			}
		} else if ctx.Command == "copy" {
			if ctx.SourcePath == "" {
				ctx.SourcePath = word
			} else if ctx.DestPath == "" {
				ctx.DestPath = word
			}
		}
	}

	if strings.Contains(word, ":") && ctx.Port == "" && ctx.Command != "copy" {
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
