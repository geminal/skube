# skube (Simple Kube)

Talk to Kubernetes in plain English. No more remembering complex kubectl syntax - just say what you want!

## Installation

### Option 1: Go Install (Recommended)

If you have Go installed:

```bash
go install github.com/geminal/skube/cmd/skube@latest
```

Make sure `~/go/bin` is in your PATH. Add this to your `~/.zshrc` or `~/.bashrc`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

Then reload your shell:
```bash
source ~/.zshrc  # or source ~/.bashrc
```

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

`skube` features **smart autocomplete** that queries your actual Kubernetes cluster!

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

### What Makes It Smart?

Unlike traditional autocomplete, `skube` **dynamically queries your cluster** to suggest:

- ✅ **Real namespaces** from your cluster (not hardcoded `prod`, `staging`)
- ✅ **Real pods** when you type `skube logs <TAB>`
- ✅ **Real deployments** when you type `skube restart deployment <TAB>`
- ✅ **Context-aware suggestions** - if you specify a namespace, it only shows resources from that namespace

**Example:**
```bash
$ skube logs pod <TAB>
# Shows: nginx-7d8b49557c-abc12  redis-6b8f9c-def34  ...
# (your actual pods!)

$ skube pods in <TAB>
# Shows: default  kube-system  production  my-app-namespace  ...
# (your actual namespaces!)
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

### Additional Resources

```bash
# Nodes
skube get nodes
skube describe node worker-1

# ConfigMaps
skube get configmaps in production
skube get cm in qa  # shorthand

# Secrets
skube get secrets in staging

# Ingress
skube get ingresses in production
skube get ing in qa  # shorthand

# PersistentVolumeClaims
skube get pvcs in production
skube get pvc in staging  # shorthand
```

### Configuration & Management

```bash
# Apply configuration
skube apply file deployment.yaml
skube create from file config.yaml

# Edit resources
skube edit deployment api in production
skube edit service backend in staging

# Delete resources
skube delete pod mypod in qa
skube delete deployment old-app in staging

# Context/Namespace management
skube use context production-cluster
skube use namespace staging
skube show config

# Copy files to/from pods
skube copy file local.txt to /tmp/remote.txt in qa
skube cp /tmp/remote.txt to local.txt in production

# Resource documentation
skube explain pod
skube what is service
skube what is ingress
```

### Metrics & Monitoring

```bash
# Resource metrics
skube show metrics pods in production
skube check usage nodes
skube check usage pods in qa
```

### Cluster Info

```bash
# Show status
skube show status in production
skube get all in qa

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

## Natural Language Features

### Talk Naturally

`skube` understands conversational English! You can use:

**Filler words** (automatically ignored):
- `the`, `a`, `an`, `please`, `me`, `for`, `my`, `our`

**Command synonyms**:
- `list`, `show`, `give`, `fetch` → `get`
- `tail`, `monitor` → `logs`
- `ssh`, `connect` → `shell`
- `change`, `modify` → `edit`
- `remove`, `destroy` → `delete`
- `reboot`, `bounce` → `restart`

**Examples:**
```bash
# All of these work!
skube please get the pods
skube show me logs for myapp
skube list all deployments in qa
skube give me the status
skube restart the backend deployment in staging
skube ssh into pod mypod
```

### Simplified Syntax

Prepositions are optional! These are equivalent:

```bash
# Traditional (still works)
skube get pods in qa
skube logs of app myapp in production

# Simplified (new!)
skube pods qa
skube logs myapp production
skube logs app myapp qa
```

### Keywords Reference

- **Actions**: `get`, `logs`, `shell`, `restart`, `scale`, `forward`, `describe`, `show`, `apply`, `delete`, `edit`, `copy`, `explain`
- **Prepositions**: `of`, `from`, `in`, `into`, `with`, `to`
- **Resources**: `pod`, `deployment`, `service`, `namespace`, `node`, `configmap`, `secret`, `ingress`, `pvc`
- **Modifiers**: `follow`, `prefix`, `search`, `find`, `last`, `max`

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

## Advanced Features

### Smart Autocomplete

Press TAB to see suggestions from your **actual cluster**:
- **Commands**: `get`, `logs`, `shell`, `restart`, `scale`, etc.
- **Keywords**: `of`, `from`, `in`, `pod`, `deployment`, `service`, etc.
- **Resources**: `namespaces`, `pods`, `deployments`, `services`
- **Real namespaces**: From your cluster (not hardcoded!)
- **Real pods**: From your cluster
- **Real deployments**: From your cluster

**Try it:**
```bash
skube <TAB><TAB>
skube get <TAB><TAB>
skube logs <TAB><TAB>
skube logs pod <TAB><TAB>  # Shows YOUR actual pods!
skube pods in <TAB><TAB>   # Shows YOUR actual namespaces!
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on:

- Development setup
- Coding standards
- Testing requirements
- Pull request process
- How to add new commands and resources

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

