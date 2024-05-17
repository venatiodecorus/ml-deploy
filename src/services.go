package main

import (
	"log"
	"strings"
)

// Handle the entire deployment process, generate Terraform config from instructions, validate, then deploy.
// Only focusing on Hetzner for now.
func hetznerDeploy(instructions string) bool {
	// Generate Terraform config from instructions
	config, err := terraformRequest(instructions)
	if err != nil {
		log.Printf("failed to generate Terraform config: %s", err)
		return false
	}

	// Clean up the returned config, OpenAI always wraps it in a JSON multi-line string
	trimmed := strings.TrimPrefix(config, "```hcl")
	trimmed = strings.TrimSuffix(trimmed, "```")
	trimmed = strings.TrimSpace(trimmed)

	// Validate the Terraform config
	if !validate(trimmed) {
		return false
	}

	// Deploy the Terraform config
	return deploy(trimmed)
}