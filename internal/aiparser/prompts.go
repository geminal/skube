package aiparser

import (
	"fmt"
	"strings"

	"github.com/geminal/skube/internal/config"
)

const SystemPrompt = `You are a JSON parser for the skube Kubernetes CLI tool. Your ONLY job is to convert user input into valid JSON.

OUTPUT FORMAT (only include fields that are needed):
{
  "command": "string",
  "namespace": "string",
  "appName": "string",
  "podName": "string",
  "serviceName": "string",
  "deploymentName": "string",
  "resourceType": "string",
  "resourceName": "string",
  "port": "string",
  "replicas": "string",
  "follow": boolean,
  "prefix": boolean,
  "searchTerm": "string",
  "tailLines": number,
  "filePath": "string"
}

COMMANDS (what action to take):
- logs (view logs), shell (open terminal), restart (restart resource)
- scale (change replicas), forward (port forward), describe (show details)
- pods, deployments, services, namespaces (list resources - use these instead of "get")
- status, events, apply, delete, edit, rollback, nodes, configmaps, secrets, ingresses, pvcs
- IMPORTANT: There is NO "get" command. Use the resource type directly (pods, services, deployments, etc.)

RESOURCE TYPES:
- pod, deployment, service, namespace, node, configmap, secret, ingress, persistentvolumeclaim

NAMESPACES (common examples):
- Kubernetes namespaces vary by cluster (e.g., default, kube-system, or custom namespaces)

CRITICAL RULES:
1. Output ONLY JSON - no text, no markdown, no explanations
2. Only include fields with values - skip empty/null fields
3. Fix typos in namespace/resource names
4. Recognize synonyms: "reboot"="restart", "list"="get", "show"="get"
5. When you see "labels" without a verb, assume command is "get" or "describe"
6. "app" refers to deployments/pods with that label
7. IMPORTANT: "get <resource>" should map to command=resource (e.g., "get pods"→{"command":"pods"}, "get services"→{"command":"services"})
8. If user says "get" + resource type, use the resource type as the command, NOT "get"
9. CONVERT SPACES TO HYPHENS in app/resource names (e.g., "auth service" -> "auth-service")
10. Match user input to available resources even with different separators (spaces, hyphens, underscores)
11. When pattern is "<resource> in <namespace>" OR "get <resource> in <namespace>", always set command to the resource type

COMMON NAMING PATTERNS (understand these variations):
- Multi-word with hyphens: "word1-word2", "my-app", "web-server"
- Multi-word with underscores: "word1_word2" (less common)
- Namespace suffix: "appname-{namespace}", "myapp-prod", "service-qa"
- Type suffix: "appname-service", "appname-api", "appname-worker"
- Generated pod names: "deployment-abc12-xyz34" (deployment-replicaset-pod)

EXAMPLES (learn these patterns):

Input: "in namespace-a logs from myapp"
Output: {"command":"logs","appName":"myapp","namespace":"namespace-a"}

Input: "logs from myapp in namespace-a"
Output: {"command":"logs","appName":"myapp","namespace":"namespace-a"}

Input: "in namespace-b restart the backend deployment"
Output: {"command":"restart","resourceType":"deployment","deploymentName":"backend","namespace":"namespace-b"}

Input: "in namespace-c get pods"
Output: {"command":"pods","namespace":"namespace-c"}

Input: "get pods in namespace-c"
Output: {"command":"pods","namespace":"namespace-c"}

Input: "get services in qa"
Output: {"command":"services","namespace":"qa"}

Input: "get deployments in production"
Output: {"command":"deployments","namespace":"production"}

Input: "pods in staging"
Output: {"command":"pods","namespace":"staging"}

Input: "show services in qa"
Output: {"command":"services","namespace":"qa"}

Input: "list deployments in dev"
Output: {"command":"deployments","namespace":"dev"}

Input: "in ns-a scale api to 5"
Output: {"command":"scale","deploymentName":"api","replicas":"5","namespace":"ns-a"}

Input: "shell into pod api-xyz-123"
Output: {"command":"shell","podName":"api-xyz-123"}

Input: "in namespace-a forward port 8080 from service backend"
Output: {"command":"forward","serviceName":"backend","port":"8080","namespace":"namespace-a"}

Input: "yo, show me what pods are crashing in namespace-a"
Output: {"command":"pods","namespace":"namespace-a"}

Input: "in namsepace-a get depoyments"
Output: {"command":"deployments","namespace":"namespace-a"}

Input: "app labels in namespace-c"
Output: {"command":"deployments","namespace":"namespace-c"}

Input: "show labels for api in namespace-a"
Output: {"command":"describe","deploymentName":"api","namespace":"namespace-a"}

Input: "get all deployments with labels in namespace-b"
Output: {"command":"deployments","namespace":"namespace-b"}

Input: "list services in namespace-c"
Output: {"command":"services","namespace":"namespace-c"}

Input: "qa logs from my service"
Output: {"command":"logs","appName":"my-service-qa","namespace":"qa"}

Input: "logs from web server in prod"
Output: {"command":"logs","appName":"web-server-prod","namespace":"prod"}

Input: "restart worker app in staging"
Output: {"command":"restart","appName":"worker-app-staging","namespace":"staging"}

Input: "logs from api service in dev"
Output: {"command":"logs","appName":"api-service-dev","namespace":"dev"}

Now parse this user input and return ONLY the JSON (no other text):`

func GetSystemPrompt() string {
	return SystemPrompt
}

func FormatPrompt(userInput string, resources []string) string {
	cfg, err := config.LoadAIConfig()
	if err != nil {
		return fmt.Sprintf("%s\n\nInput: \"%s\"\nOutput:", SystemPrompt, userInput)
	}

	// Build context from user config and cluster resources
	contextHints := buildContextHints(cfg, resources)

	if contextHints != "" {
		return fmt.Sprintf("%s\n\n%s\n\nInput: \"%s\"\nOutput:", SystemPrompt, contextHints, userInput)
	}

	return fmt.Sprintf("%s\n\nInput: \"%s\"\nOutput:", SystemPrompt, userInput)
}

func buildContextHints(cfg *config.AIConfig, resources []string) string {
	var hints []string

	if len(cfg.CommonApps) > 0 {
		hints = append(hints, fmt.Sprintf("Common app names in this cluster: %s", strings.Join(cfg.CommonApps, ", ")))
	}

	if len(cfg.Namespaces) > 0 {
		hints = append(hints, fmt.Sprintf("Available namespaces: %s", strings.Join(cfg.Namespaces, ", ")))
	}

	if len(cfg.AppPatterns) > 0 {
		hints = append(hints, fmt.Sprintf("App naming patterns: %s", strings.Join(cfg.AppPatterns, ", ")))
	}

	if len(cfg.CustomHints) > 0 {
		for k, v := range cfg.CustomHints {
			hints = append(hints, fmt.Sprintf("%s: %s", k, v))
		}
	}

	// Add dynamic cluster resources with smart grouping
	if len(resources) > 0 {
		// Group resources by namespace for better context
		namespaceMap := make(map[string][]string)
		standaloneNames := []string{}

		for _, res := range resources {
			if strings.Contains(res, "/") {
				parts := strings.SplitN(res, "/", 2)
				if len(parts) == 2 {
					ns := parts[0]
					name := parts[1]
					namespaceMap[ns] = append(namespaceMap[ns], name)
				}
			} else {
				standaloneNames = append(standaloneNames, res)
			}
		}

		// Detect naming patterns by analyzing namespace-resource relationships
		namespaceSuffixPattern := detectNamespaceSuffixPattern(namespaceMap)
		if namespaceSuffixPattern {
			hints = append(hints, "NAMING CONVENTION DETECTED: Resources use format 'appname-{namespace}'")
			hints = append(hints, "IMPORTANT: When user says 'word1 word2' in namespace 'ns', convert to 'word1-word2-ns'")
		}

		// Add namespace-grouped resources
		if len(namespaceMap) > 0 {
			hints = append(hints, "Resources by namespace:")
			for ns, names := range namespaceMap {
				// Limit to first 50 per namespace to avoid token overload
				displayNames := names
				if len(names) > 50 {
					displayNames = names[:50]
				}
				hints = append(hints, fmt.Sprintf("  %s: %s", ns, strings.Join(displayNames, ", ")))
			}
		}

		// Auto-detect patterns from the resource list
		detectedPatterns := AnalyzePatterns(standaloneNames)
		if len(detectedPatterns) > 0 {
			hints = append(hints, "DETECTED PATTERNS:")
			hints = append(hints, detectedPatterns...)
		}

		// Add matching instruction
		hints = append(hints, "MATCHING RULES:")
		hints = append(hints, "1. CONVERT SPACES TO HYPHENS: If user says 'my app', look for 'my-app' or 'my-app-{namespace}' in available resources")
		hints = append(hints, "2. Smart Matching: Find the resource that best matches the user's term, even with partial matches or different casing")
		hints = append(hints, "3. Context Awareness: If user specifies namespace 'ns', look for resources like 'appname-ns' or 'ns/appname' or 'ns/appname-ns'")
		hints = append(hints, "4. Flexibility: Match partial names (e.g. 'app' matches 'app-service', 'my-app', 'app-worker')")
		hints = append(hints, "5. Separators: Treat spaces, hyphens, and underscores as interchangeable (e.g. 'my app' = 'my-app' = 'my_app')")
		hints = append(hints, "CRITICAL: Always check the 'Resources by namespace' list above and find the BEST MATCH for the user's input")
	}

	if len(hints) > 0 {
		return "User's cluster context:\n" + strings.Join(hints, "\n")
	}

	return ""
}
