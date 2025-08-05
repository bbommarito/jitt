package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

func TestJittPassesArgumentsToGet(t *testing.T) {
	origArgs := os.Args
	defer func() { os.Args = origArgs }()

	cmd := exec.Command("go", "run", "main.go", "status")
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err)
	assert.Contains(t, string(output), "On branch")
}
