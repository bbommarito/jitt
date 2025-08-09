package jitt

import (
	"fmt"
	"os"

	"github.com/bbommarito/jitt/internal/config"
)

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
