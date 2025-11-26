#!/bin/bash
# Bash completion for skube

_skube_cache_get() {
    local key="$1"
    local cmd="$2"
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

    eval "$cmd" > "$cache_file" 2>/dev/null
    cat "$cache_file"
}

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
            local namespaces=$(_skube_cache_get "namespaces" "kubectl get namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null")
            COMPREPLY=( $(compgen -W "${namespaces}" -- ${cur}) )
            return 0
            ;;
        *)
            ;;
    esac
}

complete -F _skube_completions skube
