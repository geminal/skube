#!/bin/bash
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
            # Suggest namespaces with simple caching
            local cache_file="${TMPDIR:-/tmp}/skube-cache/namespaces"
            local now=$(date +%s)
            
            if [[ -f "$cache_file" ]]; then
                local mtime
                if [[ "$OSTYPE" == "darwin"* ]]; then
                    mtime=$(stat -f %m "$cache_file" 2>/dev/null)
                else
                    mtime=$(stat -c %Y "$cache_file" 2>/dev/null)
                fi
                
                if [[ -n "$mtime" && $((now - mtime)) -lt 5 ]]; then
                    local namespaces=$(cat "$cache_file")
                    COMPREPLY=( $(compgen -W "${namespaces}" -- ${cur}) )
                    return 0
                fi
            fi
            
            # Fetch and cache
            mkdir -p "${TMPDIR:-/tmp}/skube-cache"
            local namespaces=$(kubectl get namespaces -o jsonpath='{.items[*].metadata.name}' 2>/dev/null)
            echo "$namespaces" > "$cache_file"
            COMPREPLY=( $(compgen -W "${namespaces}" -- ${cur}) )
            return 0
            ;;
        *)
            ;;
    esac
}

complete -F _skube_completions skube
