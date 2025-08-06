package gitpassthrough

import (
	"os"
	"os/exec"
)

var Command = exec.Command

func RunGit(args []string) error {
	cmd := Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
