package completion

import (
	"fmt"

	"github.com/geminal/skube/internal/parser"
)

func HandleCompletion(ctx *parser.Context) error {
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

    # Dynamic completion helpers - query actual Kubernetes cluster
    _skube_get_namespaces() {
        kubectl get namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
    }

    _skube_get_apps() {
        local namespace="$1"
        if [[ -n "$namespace" ]]; then
            kubectl get pods -n "$namespace" -o jsonpath='{.items[*].metadata.labels.app}' 2>/dev/null | tr ' ' '\n' | sort -u | grep -v '^$'
        else
            # Try all namespaces first
            if ! kubectl get pods --all-namespaces -o jsonpath='{.items[*].metadata.labels.app}' 2>/dev/null | tr ' ' '\n' | sort -u | grep -v '^$'; then
                # Fallback to current namespace if all-namespaces is forbidden
                kubectl get pods -o jsonpath='{.items[*].metadata.labels.app}' 2>/dev/null | tr ' ' '\n' | sort -u | grep -v '^$'
            fi
        fi
    }

    _skube_get_pods() {
        local namespace="$1"
        if [[ -n "$namespace" ]]; then
            kubectl get pods -n "$namespace" -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        else
            kubectl get pods --all-namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        fi
    }

    _skube_get_deployments() {
        local namespace="$1"
        if [[ -n "$namespace" ]]; then
            kubectl get deployments -n "$namespace" -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        else
            kubectl get deployments --all-namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        fi
    }

    _skube_get_services() {
        local namespace="$1"
        if [[ -n "$namespace" ]]; then
            kubectl get services -n "$namespace" -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        else
            kubectl get services --all-namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        fi
    }

    # Extract namespace from previous words if present
    _skube_extract_namespace() {
        # Check namespace-first syntax (position 2-3)
        if [[ "${words[2]}" == "in" && -n "${words[3]}" ]]; then
            echo "${words[3]}"
            return
        fi
        # Fallback: search for "in" keyword in command
        local found_in=0
        for word in "${words[@]}"; do
            if [[ "$word" == "in" ]]; then
                found_in=1
                continue
            fi
            if [[ $found_in -eq 1 ]]; then
                echo "$word"
                return
            fi
        done
    }

    _arguments -C \
        '1: :_describe "command" commands' \
        '*:: :->args'

    case $state in
        args)
            case $words[1] in
                in)
                    # "skube in <namespace>"
                    if [[ $CURRENT -eq 2 ]]; then
                        local -a namespaces
                        namespaces=(${(f)"$(_skube_get_namespaces)"})
                        compadd "${namespaces[@]}"
                    elif [[ $CURRENT -eq 3 ]]; then
                        _describe "command" commands
                    else
                        # After "skube in <namespace> <command>", delegate to normal flow
                        # This is a simplification, ideally we'd recurse but zsh is tricky
                        # For now, just suggest resources if command is get
                        if [[ "${words[3]}" == "get" ]]; then
                             _describe "resource" resources
                        fi
                    fi
                    ;;
                get)
                    if [[ $CURRENT -eq 2 ]]; then
                        _describe "resource" resources
                    elif [[ "${words[2]}" == "pods" || "${words[2]}" == "pod" ]]; then
                        if [[ "${words[3]}" == "from" || "${words[3]}" == "in" ]]; then
                             # suggest namespaces
                             local -a namespaces
                             namespaces=(${(f)"$(_skube_get_namespaces)"})
                             compadd "${namespaces[@]}"
                        fi
                    fi
                    ;;
                logs)
                    # "skube logs ..."
                    if [[ $CURRENT -eq 2 ]]; then
                        # suggest "from", "of", "in" or app names directly?
                        # Let's suggest keywords
                        local -a keywords
                        keywords=('from:Source' 'of:App source' 'in:Namespace context')
                        _describe "keyword" keywords
                    else
                        case "${words[CURRENT-1]}" in
                            from)
                                local -a types
                                types=('app:Application' 'pod:Specific Pod')
                                _describe "type" types
                                ;;
                            of)
                                # "logs of <app>"
                                local namespace=$(_skube_extract_namespace)
                                local -a apps
                                apps=(${(f)"$(_skube_get_apps "$namespace")"})
                                compadd "${apps[@]}"
                                ;;
                            app)
                                # After "app", suggest actual app names from cluster
                                local namespace=$(_skube_extract_namespace)
                                local -a apps
                                apps=(${(f)"$(_skube_get_apps "$namespace")"})
                                compadd "${apps[@]}"
                                ;;
                            in)
                                # suggest namespaces
                                local -a namespaces
                                namespaces=(${(f)"$(_skube_get_namespaces)"})
                                compadd "${namespaces[@]}"
                                ;;
                            pod)
                                # After "pod" keyword, suggest actual pod names from namespace
                                local namespace=$(_skube_extract_namespace)
                                if [[ -n "$namespace" ]]; then
                                    local -a pods
                                    pods=(${(f)"$(_skube_get_pods "$namespace")"})
                                    compadd "${pods[@]}"
                                fi
                                ;;
                        esac
                    fi
                    ;;
                restart)
                     if [[ $CURRENT -eq 2 ]]; then
                        local -a types
                        types=('deployment:Restart Deployment' 'pod:Restart Pod')
                        _describe "type" types
                     else
                        case "${words[CURRENT-1]}" in
                            deployment)
                                local namespace=$(_skube_extract_namespace)
                                local -a deploys
                                deploys=(${(f)"$(_skube_get_deployments "$namespace")"})
                                compadd "${deploys[@]}"
                                ;;
                            pod)
                                local namespace=$(_skube_extract_namespace)
                                local -a pods
                                pods=(${(f)"$(_skube_get_pods "$namespace")"})
                                compadd "${pods[@]}"
                                ;;
                             in)
                                local -a namespaces
                                namespaces=(${(f)"$(_skube_get_namespaces)"})
                                compadd "${namespaces[@]}"
                                ;;
                        esac
                     fi
                    ;;
                 scale)
                     if [[ $CURRENT -eq 2 ]]; then
                        compadd "deployment"
                     elif [[ "${words[2]}" == "deployment" && $CURRENT -eq 3 ]]; then
                        local namespace=$(_skube_extract_namespace)
                        local -a deploys
                        deploys=(${(f)"$(_skube_get_deployments "$namespace")"})
                        compadd "${deploys[@]}"
                     fi
                    ;;
            esac
            ;;
    esac
}

# Register the completion function
compdef _skube skube
`
}

func getBashCompletion() string {
	return `#!/bin/bash
# Bash completion for skube

_skube_completions()
{
    local cur prev
    cur=${COMP_WORDS[COMP_CWORD]}
    prev=${COMP_WORDS[COMP_CWORD-1]}

    # Basic commands
    local commands="get logs shell restart scale rollback forward describe show apply delete edit config copy explain completion update help"
    local resources="namespaces pods deployments services nodes configmaps secrets ingresses pvcs"

    case "${prev}" in
        skube)
            COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
            return 0
            ;;
        get|describe|delete|edit)
            COMPREPLY=( $(compgen -W "${resources}" -- ${cur}) )
            return 0
            ;;
        in|from)
            # Suggest namespaces
            local namespaces=$(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null)
            COMPREPLY=( $(compgen -W "${namespaces}" -- ${cur}) )
            return 0
            ;;
        *)
            ;;
    esac
}

complete -F _skube_completions skube
`
}
