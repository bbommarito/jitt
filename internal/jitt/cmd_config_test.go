package jitt

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("jitt config command", func() {
	var (
		tmpDir string
		oldCwd string
	)

	BeforeEach(func() {
		tmpDir = GinkgoT().TempDir()

		var err error
		oldCwd, err = os.Getwd()
		Expect(err).To(Succeed())
		Expect(os.Chdir(tmpDir)).To(Succeed())
	})

	AfterEach(func() {
		Expect(os.Chdir(oldCwd)).To(Succeed())
	})

	Context("outside a Git repository", func() {
		It("should refuse to run config command", func() {
			command := exec.Command(pathToJittBinary, "config")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("Not inside a Git repo"))
		})
	})

	Context("inside a Git repository", func() {
		BeforeEach(func() {
			Expect(os.Mkdir(filepath.Join(tmpDir, ".git"), 0o755)).To(Succeed())
		})

		Context("with no .jitt.yaml file", func() {
			It("should report missing config file", func() {
				command := exec.Command(pathToJittBinary, "config")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(string(session.Err.Contents())).To(ContainSubstring(".jitt.yaml file not found"))
				Expect(string(session.Err.Contents())).To(ContainSubstring("run 'jitt init' first"))
			})
		})

		Context("with existing .jitt.yaml file", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: TESTPROJ"), 0o600)).To(Succeed())
			})

			It("should show all config when no arguments provided", func() {
				command := exec.Command(pathToJittBinary, "config")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("Current configuration:"))
				Expect(output).To(ContainSubstring("jira.project = TESTPROJ"))
			})

			It("should show current project when 'config project' is called", func() {
				command := exec.Command(pathToJittBinary, "config", "project")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("jira.project = TESTPROJ"))
			})

			It("should set project when 'config project VALUE' is called", func() {
				command := exec.Command(pathToJittBinary, "config", "project", "NEWPROJ")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("Set jira.project = NEWPROJ"))

				// Verify the file was actually updated
				content, err := os.ReadFile(".jitt.yaml")
				Expect(err).To(Succeed())
				Expect(string(content)).To(ContainSubstring("project: NEWPROJ"))
			})

			It("should handle unknown config keys", func() {
				command := exec.Command(pathToJittBinary, "config", "unknown")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				output := string(session.Err.Contents())
				Expect(output).To(ContainSubstring("Unknown config key: unknown"))
				Expect(output).To(ContainSubstring("Available keys: project"))
			})
		})

		Context("with .jitt.yaml file but no project configured", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: \"\""), 0o600)).To(Succeed())
			})

			It("should show empty project when getting project", func() {
				command := exec.Command(pathToJittBinary, "config", "project")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("No project configured"))
			})

			It("should show empty project in all config", func() {
				command := exec.Command(pathToJittBinary, "config")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("Current configuration:"))
				Expect(output).To(ContainSubstring("jira.project = "))
			})

			It("should be able to set project from empty", func() {
				command := exec.Command(pathToJittBinary, "config", "project", "FIRSTPROJ")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("Set jira.project = FIRSTPROJ"))

				// Verify the file was actually updated
				content, err := os.ReadFile(".jitt.yaml")
				Expect(err).To(Succeed())
				Expect(string(content)).To(ContainSubstring("project: FIRSTPROJ"))
			})
		})
	})

	Context("help message", func() {
		It("should include config command in help", func() {
			command := exec.Command(pathToJittBinary, "help")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring("config [key] [value]  Get or set configuration values"))
			Expect(output).To(ContainSubstring("jitt config       # Show all configuration"))
			Expect(output).To(ContainSubstring("jitt config project       # Show current project"))
			Expect(output).To(ContainSubstring("jitt config project XYZ   # Set project to XYZ"))
		})
	})
})
