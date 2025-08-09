package main

import (
	"fmt"
	"os"

	"github.com/bbommarito/jitt/internal/jitt"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	switch args[0] {
	case "init":
		jitt.HandleInit(args[1:])
	case "doctor":
		jitt.HandleDoctor(args[1:])
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "jitt: unknown command %q\n\n", args[0])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("jitt - Jira + Git + Tiny Tooling")
	fmt.Println()
	fmt.Println("Usage: jitt <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init [project]    Initialize .jitt.yaml configuration file")
	fmt.Println("  doctor            Check project setup and configuration")
	fmt.Println("  help              Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  jitt init         # Create .jitt.yaml file with empty project")
	fmt.Println("  jitt init ABC     # Create .jitt.yaml file with project=ABC")
	fmt.Println("  jitt doctor       # Check if setup is correct")
}
