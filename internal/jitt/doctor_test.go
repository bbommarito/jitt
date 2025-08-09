package jitt

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("jitt doctor command", func() {
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
		It("should report missing Git repository and exit with error", func() {
			command := exec.Command(pathToJittBinary, "doctor")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring("❌ Not inside a Git repository"))
			Expect(output).To(ContainSubstring("Run 'jitt init' to set up your project"))
		})
	})

	Context("inside a Git repository", func() {
		BeforeEach(func() {
			Expect(os.Mkdir(filepath.Join(tmpDir, ".git"), 0o755)).To(Succeed())
		})

		Context("with no .jitt.yaml file", func() {
			It("should report missing config file and exit with error", func() {
				command := exec.Command(pathToJittBinary, "doctor")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("✅ Git repository found"))
				Expect(output).To(ContainSubstring("❌ .jitt.yaml file not found"))
				Expect(output).To(ContainSubstring("Run 'jitt init' to set up your project"))
			})
		})

		Context("with .jitt.yaml file but no project configured", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: \"\""), 0o600)).To(Succeed())
			})

			It("should show warning about missing project but exit successfully", func() {
				command := exec.Command(pathToJittBinary, "doctor")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("✅ Git repository found"))
				Expect(output).To(ContainSubstring("✅ .jitt.yaml file exists"))
				Expect(output).To(ContainSubstring("⚠️  No project configured in .jitt.yaml"))
				Expect(output).To(ContainSubstring("✨ Setup is functional but could be improved"))
			})
		})

		Context("with .jitt.yaml file and project configured", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: TESTPROJ"), 0o600)).To(Succeed())
			})

			It("should report everything is good", func() {
				command := exec.Command(pathToJittBinary, "doctor")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("✅ Git repository found"))
				Expect(output).To(ContainSubstring("✅ .jitt.yaml file exists"))
				Expect(output).To(ContainSubstring("✅ Project configured: TESTPROJ"))
				Expect(output).To(ContainSubstring("🎉 Everything looks good!"))
			})
		})

		Context("with malformed .jitt.yaml file", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte("invalid yaml content ["), 0o600)).To(Succeed())
			})

			It("should report config loading error", func() {
				command := exec.Command(pathToJittBinary, "doctor")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				output := string(session.Out.Contents())
				Expect(output).To(ContainSubstring("✅ Git repository found"))
				Expect(output).To(ContainSubstring("✅ .jitt.yaml file exists"))
				Expect(output).To(ContainSubstring("❌ Error loading .jitt.yaml"))
			})
		})
	})

	Context("help message", func() {
		It("should include doctor command in help", func() {
			command := exec.Command(pathToJittBinary, "help")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring("doctor            Check project setup and configuration"))
			Expect(output).To(ContainSubstring("jitt doctor       # Check if setup is correct"))
		})
	})
})
