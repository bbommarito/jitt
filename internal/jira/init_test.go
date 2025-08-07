package jira

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestHasJiraFile(t *testing.T) {
	tmpDir := t.TempDir()

	oldCwd, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() {
		require.NoError(t, os.Chdir(oldCwd))
	}()

	assert.False(t, HasJiraFile())

	err := os.WriteFile(".jira", []byte("test"), 0644)
	assert.NoError(t, err)

	assert.True(t, HasJiraFile())
}

func TestInitCreatesJiraFile_WhenInGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.Mkdir(filepath.Join(tmpDir, ".git"), 0755)
	assert.NoError(t, err)

	oldCwd, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() {
		require.NoError(t, os.Chdir(oldCwd))
	}()

	code := runInitAndCaptureExitCode()

	assert.Equal(t, 0, code)

	assert.True(t, HasJiraFile())
}

func TestInitFailsOutsideGitRepo(t *testing.T) {
	tmpDir := t.TempDir()
	oldCwd, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() {
		require.NoError(t, os.Chdir(oldCwd))
	}()

	code := runInitAndCaptureExitCode()
	assert.NotEqual(t, 0, code)
	assert.False(t, HasJiraFile())
}

func TestInitFailsIfFileAlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()

	require.NoError(t, os.Mkdir(filepath.Join(tmpDir, ".git"), 0755))

	oldCwd, _ := os.Getwd()
	require.NoError(t, os.Chdir(tmpDir))
	defer func() {
		require.NoError(t, os.Chdir(oldCwd))
	}()

	err := os.WriteFile(".jira", []byte("existing"), 0644)
	assert.NoError(t, err)

	code := runInitAndCaptureExitCode()
	assert.NotEqual(t, 0, code)

	data, err := os.ReadFile(".jira")
	assert.NoError(t, err)
	assert.Equal(t, []byte("existing"), data)
}

func runInitAndCaptureExitCode() int {
	oldExit := osExit
	defer func() { osExit = oldExit }()

	var code int
	osExit = func(c int) { code = c }

	Handle([]string{"init"})
	return code
}
