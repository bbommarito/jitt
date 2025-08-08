package jira

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Jira", func() {
	Describe("HasJiraFile", func() {
		var (
			tmpDir string
			oldCwd string
		)

		BeforeEach(func() {
			tmpDir = GinkgoT().TempDir()

			// Get current directory and expect it to succeed
			var err error
			oldCwd, err = os.Getwd()
			Expect(err).To(Succeed())

			// Change to temp directory
			Expect(os.Chdir(tmpDir)).To(Succeed())
		})

		AfterEach(func() {
			// Restore original directory
			Expect(os.Chdir(oldCwd)).To(Succeed())
		})

		Context("when no .jira file exists", func() {
			It("should return false", func() {
				Expect(HasJiraFile()).To(BeFalse())
			})

			It("should indicate .jira file is absent", func() {
				Expect(".jira").NotTo(BeAnExistingFile())
			})
		})

		Context("when .jira file exists", func() {
			BeforeEach(func() {
				// Write file and expect it to succeed
				Expect(os.WriteFile(".jira", []byte("test"), 0o600)).To(Succeed())
			})

			It("should return true", func() {
				Expect(HasJiraFile()).To(BeTrue())
			})

			It("should have the .jira file present", func() {
				Expect(".jira").To(BeAnExistingFile())
			})

			It("should contain the expected content", func() {
				Expect(".jira").To(BeAnExistingFile())

				content, err := os.ReadFile(".jira")
				Expect(err).To(Succeed())
				Expect(content).To(Equal([]byte("test")))
			})
		})
	})

	Describe("HandleInit command", func() {
		var (
			tmpDir         string
			oldCwd         string
			originalOsExit func(int)
			exitCode       int
		)

		BeforeEach(func() {
			tmpDir = GinkgoT().TempDir()

			// Get current directory and expect it to succeed
			var err error
			oldCwd, err = os.Getwd()
			Expect(err).To(Succeed())

			// Change to temp directory
			Expect(os.Chdir(tmpDir)).To(Succeed())

			// Mock os.Exit to capture exit code
			originalOsExit = osExit
			exitCode = -1
			osExit = func(code int) {
				exitCode = code
			}
		})

		AfterEach(func() {
			// Restore original directory and os.Exit
			Expect(os.Chdir(oldCwd)).To(Succeed())
			osExit = originalOsExit
		})

		Context("when inside a Git repository", func() {
			BeforeEach(func() {
				// Create .git directory
				Expect(os.Mkdir(filepath.Join(tmpDir, ".git"), 0o755)).To(Succeed())
			})

			Context("when no .jira file exists", func() {
				It("should create .jira file and not call exit (success)", func() {
					HandleInit([]string{})

					Expect(exitCode).To(Equal(-1)) // osExit was never called
					Expect(HasJiraFile()).To(BeTrue())
				})

				It("should create .jira file with default content", func() {
					HandleInit([]string{})

					content, err := os.ReadFile(".jira")
					Expect(err).To(Succeed())
					Expect(string(content)).To(Equal("# jitt config\n"))
				})
			})

			Context("when project name is provided", func() {
				It("should create .jira file with project configuration", func() {
					HandleInit([]string{"ABC"})

					content, err := os.ReadFile(".jira")
					Expect(err).To(Succeed())
					Expect(string(content)).To(ContainSubstring(`project = "ABC"`))
					Expect(exitCode).To(Equal(-1)) // osExit was never called (success)
				})
			})

			Context("when .jira file already exists", func() {
				BeforeEach(func() {
					Expect(os.WriteFile(".jira", []byte("existing"), 0o600)).To(Succeed())
				})

				It("should fail and not overwrite existing file", func() {
					HandleInit([]string{})

					Expect(exitCode).NotTo(Equal(0))

					data, err := os.ReadFile(".jira")
					Expect(err).To(Succeed())
					Expect(data).To(Equal([]byte("existing")))
				})
			})
		})

		Context("when outside a Git repository", func() {
			It("should fail and not create .jira file", func() {
				HandleInit([]string{})

				Expect(exitCode).NotTo(Equal(0))
				Expect(HasJiraFile()).To(BeFalse())
			})
		})
	})
})
