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
