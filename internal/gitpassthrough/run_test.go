package gitpassthrough

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunGit_CallsExecCommandWithCorrectArgs(t *testing.T) {
	var captured []string

	original := Command
	defer func() { Command = original }()

	Command = func(name string, args ...string) *exec.Cmd {
		captured = append([]string{name}, args...)
		return exec.Command("echo")
	}

	_ = RunGit([]string{"log", "--oneline"})

	expected := []string{"git", "log", "--oneline"}
	assert.Equal(t, expected, captured)
}
