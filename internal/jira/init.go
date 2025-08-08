package jira

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	minArgsForProject = 2
	jiraFilePerms     = 0o600
)

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

func Handle(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "jitt jira: no command given")
		osExit(1)
		return
	}

	switch args[0] {
	case "init":
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

		if len(args) >= minArgsForProject {
			project := args[1]
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
	default:
		fmt.Fprintf(os.Stderr, "jitt jira: %q not recognized\n", args[0])
		osExit(1)
	}
}
