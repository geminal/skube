#compdef skube

# Initialize completion system if not already done
if [[ -z "$_comps" ]]; then
    autoload -Uz compinit && compinit
fi

_skube() {
    local curcontext="$curcontext" state line
    typeset -A opt_args

    local -a _skube_cmds
    _skube_cmds=(
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

    local -a _skube_resources
    _skube_resources=(
        'namespaces:List all namespaces (start here!)'
        'pods:List all pods'
        'deployments:List deployments'
        'services:List services'
        'ns:Shorthand for namespaces'
        'pod:Shorthand for pods'
        'deploy:Shorthand for deployments'
        'svc:Shorthand for services'
    )

    # Caching helper
    _skube_cache_get() {
        local key="$1"
        local func="$2"
        shift 2
        local cache_dir="${TMPDIR:-/tmp}/skube-cache"
        local cache_file="$cache_dir/$key"
        local now=$(date +%s)

        mkdir -p "$cache_dir"

        if [[ -f "$cache_file" ]]; then
            local mtime
            if [[ "$OSTYPE" == "darwin"* ]]; then
                mtime=$(stat -f %m "$cache_file")
            else
                mtime=$(stat -c %Y "$cache_file" 2>/dev/null)
            fi
            
            if [[ -n "$mtime" && $((now - mtime)) -lt 5 ]]; then
                cat "$cache_file"
                return
            fi
        fi

        "$func" "$@" > "$cache_file" 2>/dev/null
        cat "$cache_file"
    }

    # Fetch functions
    _skube_fetch_namespaces() {
        kubectl get namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
    }

    _skube_fetch_apps() {
        local namespace="$1"
        if [[ -n "$namespace" ]]; then
            kubectl get pods -n "$namespace" -o jsonpath='{.items[*].metadata.labels.app}' 2>/dev/null | tr ' ' '\n' | sort -u | grep -v '^$'
        else
            if ! kubectl get pods --all-namespaces -o jsonpath='{.items[*].metadata.labels.app}' 2>/dev/null | tr ' ' '\n' | sort -u | grep -v '^$'; then
                kubectl get pods -o jsonpath='{.items[*].metadata.labels.app}' 2>/dev/null | tr ' ' '\n' | sort -u | grep -v '^$'
            fi
        fi
    }

    _skube_fetch_pods() {
        local namespace="$1"
        if [[ -n "$namespace" ]]; then
            kubectl get pods -n "$namespace" -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        else
            kubectl get pods --all-namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        fi
    }

    _skube_fetch_deployments() {
        local namespace="$1"
        if [[ -n "$namespace" ]]; then
            kubectl get deployments -n "$namespace" -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        else
            kubectl get deployments --all-namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        fi
    }

    _skube_fetch_services() {
        local namespace="$1"
        if [[ -n "$namespace" ]]; then
            kubectl get services -n "$namespace" -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        else
            kubectl get services --all-namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null | tr ' ' '\n'
        fi
    }

    # Dynamic completion helpers - query actual Kubernetes cluster
    _skube_get_namespaces() {
        _skube_cache_get "namespaces" _skube_fetch_namespaces
    }

    _skube_get_apps() {
        local namespace="$1"
        local cache_key="apps_${namespace:-all}"
        _skube_cache_get "$cache_key" _skube_fetch_apps "$namespace"
    }

    _skube_get_pods() {
        local namespace="$1"
        local cache_key="pods_${namespace:-all}"
        _skube_cache_get "$cache_key" _skube_fetch_pods "$namespace"
    }

    _skube_get_deployments() {
        local namespace="$1"
        local cache_key="deployments_${namespace:-all}"
        _skube_cache_get "$cache_key" _skube_fetch_deployments "$namespace"
    }

    _skube_get_services() {
        local namespace="$1"
        local cache_key="services_${namespace:-all}"
        _skube_cache_get "$cache_key" _skube_fetch_services "$namespace"
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

    _arguments \
        '1:command:->command' \
        '*::arg:->args'

    case $state in
        command)
            _describe 'command' _skube_cmds
            ;;
        args)
            case $words[1] in
                in)
                    # "skube in <namespace>"
                    if [[ $CURRENT -eq 2 ]]; then
                        local -a namespaces
                        namespaces=(${(f)"$(_skube_get_namespaces)"})
                        compadd "${namespaces[@]}"
                    elif [[ $CURRENT -eq 3 ]]; then
                        _describe "command" _skube_cmds
                    else
                        # After "skube in <namespace> <command>", delegate to normal flow
                        # This is a simplification, ideally we'd recurse but zsh is tricky
                        # For now, just suggest resources if command is get
                        if [[ "${words[3]}" == "get" ]]; then
                             _describe "resource" _skube_resources
                        fi
                    fi
                    ;;
                get)
                    if [[ $CURRENT -eq 2 ]]; then
                        _describe "resource" _skube_resources
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
                    # "skube logs in <namespace> ..." - namespace-first syntax
                    if [[ $CURRENT -eq 2 ]]; then
                        # Only suggest "in" to enforce namespace-first
                        compadd "in"
                    elif [[ $CURRENT -eq 3 && "${words[2]}" == "in" ]]; then
                        # After "in", suggest namespaces
                        local -a namespaces
                        namespaces=(${(f)"$(_skube_get_namespaces)"})
                        compadd "${namespaces[@]}"
                    elif [[ $CURRENT -eq 4 && "${words[2]}" == "in" ]]; then
                        # After namespace, suggest "from"
                        compadd "from"
                    elif [[ $CURRENT -eq 5 && "${words[4]}" == "from" ]]; then
                        # After "from", suggest app or pod
                        local -a types
                        types=('app:Application' 'pod:Specific Pod')
                        _describe "type" types
                    elif [[ $CURRENT -eq 6 ]]; then
                        # After type (app/pod), suggest actual resources from namespace
                        local namespace="${words[3]}"
                        case "${words[5]}" in
                            app)
                                local -a apps
                                apps=(${(f)"$(_skube_get_apps "$namespace")"})
                                compadd "${apps[@]}"
                                ;;
                            pod)
                                local -a pods
                                pods=(${(f)"$(_skube_get_pods "$namespace")"})
                                compadd "${pods[@]}"
                                ;;
                        esac
                    elif [[ $CURRENT -gt 6 ]]; then
                        # After app/pod name, suggest additional options
                        local -a log_options
                        log_options=(
                            'follow:Follow log output (-f)'
                            '-f:Follow log output'
                            'prefix:Show pod name prefix'
                            'search:Search/filter logs'
                            'tail:Show last N lines'
                        )
                        _describe "log options" log_options
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
