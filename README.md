# skube

Talk to Kubernetes in plain English. No more remembering complex kubectl syntax - just say what you want!

## Installation

### Option 1: Go Install (Recommended for Go developers)

If you have Go installed:

```bash
go install github.com/geminal/skube/cmd/skube@latest
```

Ensure your `$GOPATH/bin` is in your `$PATH`.

### Option 2: Pre-built Binaries

Download the latest release for your OS/Arch from the [Releases](https://github.com/geminal/skube/releases) page.

1. Unzip the archive
2. Move `skube` to a directory in your `$PATH` (e.g., `/usr/local/bin`)

### Option 3: Build from Source

```bash
git clone https://github.com/geminal/skube.git
cd skube
go install ./cmd/skube
```




### Verify Installation

```bash
skube help
skube get namespaces
```

## Autocomplete

### Zsh

Add this to your `~/.zshrc`:

```bash
source <(skube completion zsh)
```

### Bash

Add this to your `~/.bashrc`:

```bash
source <(skube completion bash)
```


### Update

To get the latest version:

```bash
skube update
```

This runs `go install github.com/geminal/skube/cmd/skube@latest` under the hood.

Alternatively, you can run the `go install` command manually or download the latest binary release.




## Quick Start

Talk to Kubernetes naturally:

```bash
# Instead of: kubectl get namespaces
skube get namespaces

# Instead of: kubectl get pods -n production
skube get pods from production namespace

# Instead of: kubectl logs -f -l app=myapp --prefix=true -n prod
skube logs of myapp in prod follow with prefix

# Instead of: kubectl logs my-service -n staging | grep ERROR
skube logs of my-service in staging search "ERROR"

# Instead of: kubectl logs pod-abc123 --tail=100 -n qa
skube logs from pod pod-abc123 get last 100 in qa
```

## Conversational Commands

### Quick Investigation

Perfect for your daily workflow:

```bash
# List environments
skube get namespaces

# Check what's running in an environment
skube get pods from production namespace
skube get pods from qa namespace

# Check specific app pods
skube get pods of myapp in qa

# Tail logs from all pods of an app (with pod names shown)
skube logs of myapp in prod follow with prefix

# Tail logs from many pods (increase concurrent stream limit)
skube logs of webapp in production follow max 30

# Search for errors in logs
skube logs of myapp in qa search "error"
skube logs of myapp in qa find "timeout"

# Get last N lines from logs
skube logs from pod api-abc123 get last 100 in staging
```

### Pod Operations

```bash
# List pods
skube get pods in production
skube get pods of myapp in qa

# View logs
skube logs from pod api-abc123 in staging
skube logs from pod api-abc123 in staging follow

# Shell into pod
skube shell into pod api-abc123 in qa

# Restart pod
skube restart pod api-abc123 in production

# Describe pod
skube describe pod api-abc123 in staging
```

### Deployment Operations

```bash
# List deployments
skube get deployments in staging

# Restart deployment
skube restart deployment backend in prod

# Scale deployment
skube scale deployment api to 5 in production
skube scale deployment backend to 3 in staging

# Rollback deployment
skube rollback deployment api in staging
```

### Service Operations

```bash
# List services
skube get services in production

# Port forward
skube forward service my-service port 8080 in prod
skube forward service backend port 3000 in staging

# Describe service
skube describe service api in qa
```

### Cluster Info

```bash
# Show status
skube show status in production

# Show events
skube show events in qa
```

## Real-World Examples

### Daily Investigation Workflow

```bash
# 1. Check available environments
skube get namespaces

# 2. See what's running in QA
skube get pods from qa namespace

# 3. Check specific app
skube get pods of myapp in qa

# 4. Tail logs from all pods (with pod name prefixes)
skube logs of myapp in qa follow with prefix

# 5. Search for specific errors
skube logs of myapp in qa find "connection refused"

# 6. Get last 100 lines from a specific pod
skube logs from pod myapp-abc123 get last 100 in qa

# 7. If needed, shell into a pod
skube shell into pod myapp-abc123 in qa
```

### Quick Operations

```bash
# Restart a deployment
skube restart deployment backend in prod

# Scale up for traffic
skube scale deployment api to 10 in production

# Port forward for local testing
skube forward service my-service port 8080 in staging

# Check events for debugging
skube show events in production
```

## Requirements

- kubectl must be installed and configured
- Active Kubernetes context

## Natural Language Keywords

skube understands these natural language patterns:

- **Actions**: `get`, `logs`, `shell`, `restart`, `scale`, `forward`, `describe`, `show`
- **Prepositions**: `of`, `from`, `in`, `into`, `with`
- **Resources**: `pod`, `deployment`, `service`, `namespace`
- **Modifiers**: `follow`, `prefix`, `search`, `find`, `last`, `max`

Mix and match them naturally:
- `skube get pods of myapp in qa`
- `skube logs of api in prod follow with prefix`
- `skube logs from pod xyz get last 100 in staging`

## Tips

- **Use TAB autocomplete** - Press TAB to see available commands, keywords, and namespaces
- **Talk naturally** - Use `of`, `from`, `in`, `to` to make commands readable
- **Both syntaxes work** - Old flag style (`-n namespace`) still works
- **Flexible word order** - `in qa` and `from qa` both work for namespaces
- **Log all pods** - Use `of <appname>` to get logs from all pods of an app
- **Prefixes help** - Add `with prefix` to see which pod each log line comes from
- **Search logs** - Use `search "term"` or `find "term"` to filter logs
- **Last N lines** - Use `get last 100` to tail specific number of lines
- **Many pods** - Use `max 30` to increase concurrent log stream limit (default is 5)
- **Dry Run** - Use `--dry-run` to see the kubectl command without executing it

## Autocomplete

Autocomplete is enabled automatically during installation. It suggests:
- **Commands**: `get`, `logs`, `shell`, `restart`, `scale`, etc.
- **Keywords**: `of`, `from`, `in`, `pod`, `deployment`, `service`, etc.
- **Resources**: `namespaces`, `pods`, `deployments`, `services`
- **Common namespaces**: `production`, `staging`, `qa`, `dev`, `prod`

**Try it:**
```bash
skube <TAB><TAB>
skube get <TAB><TAB>
skube logs <TAB><TAB>
skube logs of myapp in <TAB><TAB>
```

## Contributing

Contributions are welcome! To add new commands:

1. Add handler function in `main.go`
2. Add case in `executeCommand()` switch statement
3. Update help text in `printHelp()`
4. Update this README

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

