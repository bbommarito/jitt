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

func TestInit_CreatesJiraFileWhenInGitRepo(t *testing.T) {
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

func TestInit_CreatesJiraFileWithProject(t *testing.T) {
	tmpDir := t.TempDir()

	err := os.Mkdir(filepath.Join(tmpDir, ".git"), 0755)
	assert.NoError(t, err)

	oldCwd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Chdir(oldCwd))
	}()

	err = os.Chdir(tmpDir)
	require.NoError(t, err)

	var exitCode int
	osExit = func(code int) { exitCode = code }
	defer func() { osExit = os.Exit }()

	Handle([]string{"init", "ABC"})

	content, err := os.ReadFile(filepath.Join(tmpDir, ".jira"))
	require.NoError(t, err)
	require.Contains(t, string(content), `project = "ABC"`)
	require.Equal(t, 0, exitCode)
}

func TestInit_FailsOutsideGitRepo(t *testing.T) {
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

func TestInit_FailsIfFileAlreadyExists(t *testing.T) {
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
