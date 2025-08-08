package gitpassthrough

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RunGit", func() {
	var (
		originalCommand func(string, ...string) *exec.Cmd
		captured        []string
	)

	BeforeEach(func() {
		// Save original Command function
		originalCommand = Command
		captured = []string{}

		// Mock Command to capture arguments
		Command = func(name string, args ...string) *exec.Cmd {
			captured = append([]string{name}, args...)
			return exec.Command("echo") // Use echo as a safe command that won't fail
		}
	})

	AfterEach(func() {
		// Restore original Command function
		Command = originalCommand
	})

	Context("when calling RunGit with arguments", func() {
		It("should call exec.Command with 'git' and the provided arguments", func() {
			_ = RunGit([]string{"log", "--oneline"})

			expected := []string{"git", "log", "--oneline"}
			Expect(captured).To(Equal(expected))
		})

		It("should handle empty arguments", func() {
			_ = RunGit([]string{})

			expected := []string{"git"}
			Expect(captured).To(Equal(expected))
		})

		It("should handle multiple arguments", func() {
			_ = RunGit([]string{"status", "--porcelain", "--untracked-files=no"})

			expected := []string{"git", "status", "--porcelain", "--untracked-files=no"}
			Expect(captured).To(Equal(expected))
		})
	})
})
