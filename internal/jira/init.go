package jira

import (
	"fmt"
	"os"
	"path/filepath"
)

const jiraFilePerms = 0o600

var osExit = os.Exit

func HasJiraFile() bool {
	_, err := os.Stat(".jira")
	return err == nil
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
		fmt.Fprintln(os.Stderr, "Not inside a Git repo. .jira not created")
		osExit(1)
		return
	}

	if HasJiraFile() {
		fmt.Fprintln(os.Stderr, ".jira already exists — not overwriting.")
		osExit(1)
		return
	}

	var content string

	if len(args) >= 1 {
		project := args[0]
		content = fmt.Sprintf("project = %q\n", project)
	} else {
		content = "# jitt config\n"
	}

	err := os.WriteFile(".jira", []byte(content), jiraFilePerms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating .jira: %v\n", err)
		osExit(1)
		return
	}

	fmt.Println(".jira created")
}
