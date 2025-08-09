package jitt

import (
	"fmt"

	"github.com/bbommarito/jitt/internal/config"
)

// HandleDoctor handles the 'jitt doctor' command
func HandleDoctor(args []string) {
	var issues []string
	var warnings []string

	// Check if we're in a Git repository
	if !isGitRepo() {
		issues = append(issues, "❌ Not inside a Git repository")
	} else {
		fmt.Println("✅ Git repository found")
	}

	// Check if .jitt.yaml file exists
	if !HasConfigFile() {
		issues = append(issues, "❌ .jitt.yaml file not found")
	} else {
		fmt.Println("✅ .jitt.yaml file exists")

		// Check if there's a project configured
		cfg, err := config.Load()
		if err != nil {
			issues = append(issues, fmt.Sprintf("❌ Error loading .jitt.yaml: %v", err))
		} else if cfg.Jira.Project == "" {
			warnings = append(warnings, "⚠️  No project configured in .jitt.yaml")
		} else {
			fmt.Printf("✅ Project configured: %s\n", cfg.Jira.Project)
		}
	}

	// Print warnings
	for _, warning := range warnings {
		fmt.Println(warning)
	}

	// Print issues and exit with error code if any found
	if len(issues) > 0 {
		fmt.Println()
		for _, issue := range issues {
			fmt.Println(issue)
		}
		fmt.Println()
		fmt.Println("Run 'jitt init' to set up your project.")
		osExit(1)
		return
	}

	// All good!
	if len(warnings) == 0 {
		fmt.Println()
		fmt.Println("🎉 Everything looks good!")
	} else {
		fmt.Println()
		fmt.Println("✨ Setup is functional but could be improved.")
	}
}
