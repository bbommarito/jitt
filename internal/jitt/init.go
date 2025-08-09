package jitt

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bbommarito/jitt/internal/config"
)

var osExit = os.Exit

func HasConfigFile() bool {
	return config.Exists()
}

func isGitRepo() bool {
	dir, err := os.Getwd()
	if err != nil {
		return false
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}

		dir = parent
	}
}

// HandleInit handles the 'jitt init' command
func HandleInit(args []string) {
	if !isGitRepo() {
		fmt.Fprintln(os.Stderr, "Not inside a Git repo. Config not created")
		osExit(1)
		return
	}

	if HasConfigFile() {
		fmt.Fprintln(os.Stderr, ".jitt.yaml already exists — not overwriting.")
		osExit(1)
		return
	}

	var project string
	if len(args) >= 1 {
		project = args[0]
	}

	err := config.Create(project)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating .jitt.yaml: %v\n", err)
		osExit(1)
		return
	}

	fmt.Println(".jitt.yaml created")
}

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

// HandleConfig handles the 'jitt config' command
func HandleConfig(args []string) {
	// Check if we're in a Git repository
	if !isGitRepo() {
		fmt.Fprintln(os.Stderr, "Not inside a Git repo.")
		osExit(1)
		return
	}

	// Check if config file exists
	if !HasConfigFile() {
		fmt.Fprintln(os.Stderr, ".jitt.yaml file not found - run 'jitt init' first")
		osExit(1)
		return
	}

	// If no args, show all config
	if len(args) == 0 {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
			osExit(1)
			return
		}

		fmt.Println("Current configuration:")
		fmt.Printf("  jira.project = %s\n", cfg.Jira.Project)
		return
	}

	// Handle specific config keys
	switch args[0] {
	case "project":
		if len(args) == 1 {
			// Show current project
			cfg, err := config.Load()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
				osExit(1)
				return
			}
			if cfg.Jira.Project == "" {
				fmt.Println("No project configured")
			} else {
				fmt.Printf("jira.project = %s\n", cfg.Jira.Project)
			}
		} else {
			// Set project
			newProject := args[1]
			err := config.Update("jira.project", newProject)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error updating config: %v\n", err)
				osExit(1)
				return
			}
			fmt.Printf("Set jira.project = %s\n", newProject)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown config key: %s\n", args[0])
		fmt.Fprintln(os.Stderr, "Available keys: project")
		osExit(1)
	}
}
