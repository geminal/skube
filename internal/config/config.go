package config

import (
	"os"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

// FindSkubeRepo finds the skube repository path
func FindSkubeRepo() string {
	// 1. Check SKUBE_REPO environment variable
	if repoPath := os.Getenv("SKUBE_REPO"); repoPath != "" {
		if isValidSkubeRepo(repoPath) {
			return repoPath
		}
	}

	// 2. Search for git repos in common locations
	homeDir := os.Getenv("HOME")
	searchPaths := []string{
		homeDir + "/Projects",
		homeDir + "/Code",
		homeDir + "/Workspace",
		homeDir + "/Work",
		homeDir + "/Git",
		homeDir + "/repos",
		homeDir + "/src",
		homeDir + "/go/src",
	}

	for _, basePath := range searchPaths {
		if found := searchForSkubeRepo(basePath); found != "" {
			return found
		}
	}

	return ""
}

func isValidSkubeRepo(path string) bool {
	// Check if directory exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	// Check if it's a git repository
	gitDir := path + "/.git"
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return false
	}

	// Check if main.go exists (checking for go.mod might be better in the future)
	goMod := path + "/go.mod"
	if _, err := os.Stat(goMod); os.IsNotExist(err) {
		return false
	}

	return true
}

func searchForSkubeRepo(basePath string) string {
	// Check if base path exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return ""
	}

	// Search for skube at base level (generic pattern)
	fullPath := basePath + "/skube"
	if isValidSkubeRepo(fullPath) {
		return fullPath
	}

	return ""
}
