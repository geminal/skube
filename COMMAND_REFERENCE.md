# skube Conversational Command Reference

Complete mapping of skube conversational commands to kubectl equivalents.

---

> [!TIP]
> You can also use the `--dry-run` flag with any command to see the exact `kubectl` command it would execute, without actually running it.
> Example: `skube logs of api in prod --dry-run`

---

## Investigation Commands

### List Namespaces

| skube | kubectl equivalent |
|----------|-------------------|
| `skube get namespaces` | `kubectl get namespaces` |
| `skube namespaces` | `kubectl get namespaces` |

### List Pods

| skube | kubectl equivalent |
|----------|-------------------|
| `skube get pods` | `kubectl get pods -o wide` |
| `skube get pods from production namespace` | `kubectl get pods -n production -o wide` |
| `skube get pods in qa` | `kubectl get pods -n qa -o wide` |
| `skube get pods of myapp` | `kubectl get pods -l app=myapp -o wide` |
| `skube get pods of myapp in qa` | `kubectl get pods -l app=myapp -n qa -o wide` |
| `skube pods in staging` | `kubectl get pods -n staging -o wide` |

### View Logs

| skube | kubectl equivalent |
|----------|-------------------|
| `skube logs from pod api-abc123 in qa` | `kubectl logs api-abc123 -n qa` |
| `skube logs from pod api-abc123 in qa follow` | `kubectl logs api-abc123 -f -n qa` |
| `skube logs from pod api-abc123 get last 100 in qa` | `kubectl logs api-abc123 --tail=100 -n qa` |

### Logs from All App Pods

| skube | kubectl equivalent |
|----------|-------------------|
| `skube logs of myapp in prod follow` | `kubectl logs -l app=myapp -f -n prod` |
| `skube logs of api in prod follow with prefix` | `kubectl logs -l app=api -f --prefix=true -n prod` |
| `skube logs of backend in staging with prefix` | `kubectl logs -l app=backend --prefix=true -n staging` |
| `skube logs of webapp in production follow max 30` | `kubectl logs -l app=webapp -f --max-log-requests=30 -n namespace` |

### Search Logs

| skube | kubectl equivalent |
|----------|-------------------|
| `skube logs from pod api-abc123 search "error" in qa` | `kubectl logs api-abc123 -n qa \| grep --color=always "error"` |
| `skube logs of myapp find timeout in prod` | `kubectl logs -l app=myapp -n prod \| grep --color=always timeout` |
| `skube logs from pod xyz search "500" in staging` | `kubectl logs xyz -n staging \| grep --color=always "500"` |

---

## Pod Operations

### Shell into Pod

| skube | kubectl equivalent |
|----------|-------------------|
| `skube shell into pod api-abc123 in staging` | `kubectl exec -it api-abc123 -n staging -- sh` |
| `skube shell pod backend-xyz in qa` | `kubectl exec -it backend-xyz -n qa -- sh` |

### Restart Pod

| skube | kubectl equivalent |
|----------|-------------------|
| `skube restart pod api-abc123 in production` | `kubectl delete pod api-abc123 -n production` |
| `skube restart pod backend-xyz in staging` | `kubectl delete pod backend-xyz -n staging` |

### Describe Pod

| skube | kubectl equivalent |
|----------|-------------------|
| `skube describe pod api-abc123 in production` | `kubectl describe pod api-abc123 -n production` |
| `skube describe pod backend-xyz in qa` | `kubectl describe pod backend-xyz -n qa` |

---

## Deployment Operations

### List Deployments

| skube | kubectl equivalent |
|----------|-------------------|
| `skube get deployments` | `kubectl get deployments -o wide` |
| `skube get deployments in staging` | `kubectl get deployments -n staging -o wide` |
| `skube get deployments from prod` | `kubectl get deployments -n prod -o wide` |
| `skube deployments in prod` | `kubectl get deployments -n prod -o wide` |

### Restart Deployment

| skube | kubectl equivalent |
|----------|-------------------|
| `skube restart deployment backend in prod` | `kubectl rollout restart deployment backend -n prod` |
| `skube restart deployment api in staging` | `kubectl rollout restart deployment api -n staging` |

### Scale Deployment

| skube | kubectl equivalent |
|----------|-------------------|
| `skube scale deployment api to 5 in production` | `kubectl scale deployment api --replicas=5 -n production` |
| `skube scale deployment backend to 3 in staging` | `kubectl scale deployment backend --replicas=3 -n staging` |

### Rollback Deployment

| skube | kubectl equivalent |
|----------|-------------------|
| `skube rollback deployment api in staging` | `kubectl rollout undo deployment api -n staging` |
| `skube rollback deployment backend in prod` | `kubectl rollout undo deployment backend -n prod` |

---

## Service Operations

### List Services

| skube | kubectl equivalent |
|----------|-------------------|
| `skube get services` | `kubectl get services -o wide` |
| `skube get services in production` | `kubectl get services -n production -o wide` |
| `skube get services from qa` | `kubectl get services -n qa -o wide` |
| `skube services in qa` | `kubectl get services -n qa -o wide` |

### Port Forward

| skube | kubectl equivalent |
|----------|-------------------|
| `skube forward service my-service port 8080 in prod` | `kubectl port-forward service/my-service 8080:8080 -n prod` |
| `skube forward service backend port 3000 in staging` | `kubectl port-forward service/backend 3000:3000 -n staging` |

### Describe Service

| skube | kubectl equivalent |
|----------|-------------------|
| `skube describe service api in qa` | `kubectl describe service api -n qa` |
| `skube describe service backend in prod` | `kubectl describe service backend -n prod` |

### Finding Pods Behind a Service

Services don't "contain" pods - they use label selectors to route traffic to matching pods.

**Workflow to find pods for a service:**
```bash
# 1. Describe service to see its selector labels
skube describe service myservice in qa
# Look for "Selector:" in output, e.g., "app=myservice"

# 2. Get pods matching those labels
skube get pods of myservice in qa
```

**Equivalent kubectl commands:**
```bash
# 1. Get service selector
kubectl describe service myservice -n qa | grep Selector

# 2. Get matching pods
kubectl get pods -l app=myservice -n qa -o wide
```

---

## Cluster Info

### Show Status

| skube | kubectl equivalent |
|----------|-------------------|
| `skube show status in production` | `kubectl get all -n production` |
| `skube status in qa` | `kubectl get all -n qa` |

### Show Events

| skube | kubectl equivalent |
|----------|-------------------|
| `skube show events in production` | `kubectl get events --sort-by=.lastTimestamp -n production` |
| `skube events in qa` | `kubectl get events --sort-by=.lastTimestamp -n qa` |

### Get All Resources

| skube | kubectl equivalent |
|----------|-------------------|
| `skube get all` | `kubectl get all -o wide` |
| `skube get all in production` | `kubectl get all -n production -o wide` |

### More Resources

| skube | kubectl equivalent |
|----------|-------------------|
| `skube get nodes` | `kubectl get nodes -o wide` |
| `skube get configmaps in prod` | `kubectl get configmaps -n prod` |
| `skube get secrets in qa` | `kubectl get secrets -n qa` |
| `skube get ingress in staging` | `kubectl get ingress -n staging` |
| `skube get pvc in dev` | `kubectl get pvc -n dev` |

---

## Utility Commands

### Update skube

| skube | Description |
|----------|-------------------|
| `skube update` | Updates skube to the latest version via `go install` |

### Generate Completion

| skube | Description |
|----------|-------------------|
| `skube completion zsh` | Generates Zsh completion script |
| `skube completion bash` | Generates Bash completion script |

---

## Natural Language Patterns

skube understands natural language patterns. These are all equivalent:

### Namespace Patterns
- `in production`
- `from production`
- `in production namespace`
- `from production namespace`
- `-n production` (traditional flag style)

### Resource Shorthands
- `namespaces` = `ns`
- `pods` = `pod`
- `deployments` = `deploy`
- `services` = `svc`

Examples:
- `skube get ns` = `skube get namespaces`
- `skube get deploy in qa` = `skube get deployments in qa`
- `skube get svc from prod` = `skube get services from prod`

### App Selection
- `of myapp` (preferred - cleaner syntax)
- Using full label selector with kubectl: `-l app=myapp`

### Pod Selection
- `from pod api-abc123`
- `pod api-abc123`

### Alternative Command Syntax
- `shell` = `exec`
- `show status` = `status`
- `show events` = `events`

### Log Modifiers
- `follow` = `-f`
- `with prefix` or just `prefix` = `--prefix=true`
- `search "term"` or `find "term"` = `| grep term`
- `get last 100` = `--tail=100`
- `max 30` = `--max-log-requests=30` (for following logs from many pods)

---

## Daily Workflow Examples

### Morning Health Check

```bash
# Check all environments
skube get namespaces

# Check production status
skube show status in production
skube show events in production

# Check specific apps
skube get pods of api in production
skube get pods of backend in production
```

### Investigating Issues

```bash
# Find the app
skube get pods of myapp in qa

# Tail logs from all pods
skube logs of myapp in qa follow with prefix

# Search for specific errors
skube logs of myapp in qa find "connection refused"

# Get recent logs from specific pod
skube logs from pod myapp-abc123 get last 100 in qa

# Shell in if needed
skube shell into pod myapp-abc123 in qa
```

### Deployment Tasks

```bash
# Check current state
skube get deployments in staging

# Scale up
skube scale deployment api to 10 in production

# Restart after config change
skube restart deployment backend in prod

# Rollback if needed
skube rollback deployment backend in prod
```

### Local Development

```bash
# Port forward to access service
skube forward service my-service port 8080 in staging

# Watch logs
skube logs of my-service in staging follow

# Debug specific pod
skube shell into pod my-service-abc123 in staging
```

---

## Tips

1. **Flexible syntax** - Word order doesn't matter much: `in qa` = `from qa`
2. **Natural prepositions** - Use `from`, `in`, `into`, `with` to make it readable
3. **Both styles work** - Traditional flags (`-n`, `-l`) still work alongside natural language
4. **App = all pods** - `of <name>` targets all pods with that app label
5. **Combine modifiers** - `follow with prefix` works together
6. **Quote search terms** - Use quotes for multi-word searches: `search "connection refused"`
