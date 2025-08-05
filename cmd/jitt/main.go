package main

import (
	"github.com/bbommarito/jitt/internal/gitpassthrough"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	args := os.Args[1:]

	err := gitpassthrough.RunGit(args)
	if err != nil {
		os.Exit(1)
	}
}
