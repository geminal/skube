package aiparser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatPrompt_GenericService(t *testing.T) {
	// This test simulates the user's failing scenario
	userInput := "qa logs from my service"
	resources := []string{
		"qa/my-service-qa",
		"qa/other-service",
	}

	prompt := FormatPrompt(userInput, resources)

	// We want to ensure the prompt contains instructions that would help the AI
	// connect "my service" to "my-service-qa"

	// Check for generic rules that should cover this
	assert.Contains(t, prompt, "Smart Matching")
	assert.Contains(t, prompt, "Flexibility")
	assert.Contains(t, prompt, "Separators")

	// We might need to add a specific rule about spaces vs hyphens if this isn't enough
	// For now, let's just see what the prompt looks like
	t.Logf("Generated Prompt:\n%s", prompt)
}
