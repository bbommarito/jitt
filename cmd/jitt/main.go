package main

import (
	"os"

	"github.com/bbommarito/jitt/internal/gitpassthrough"
	"github.com/bbommarito/jitt/internal/jira"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		os.Exit(1)
	}

	switch args[0] {
	case "jira":
		jira.Handle(args[1:])
	default:
		err := gitpassthrough.RunGit(args)
		if err != nil {
			os.Exit(1)
		}
	}
}
