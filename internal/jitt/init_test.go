package jitt

import (
	"os"
	"os/exec"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("jitt init command", func() {
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
		It("should refuse to create config file with helpful error", func() {
			command := exec.Command(pathToJittBinary, "init")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("Not inside a Git repo"))
			Expect(string(session.Err.Contents())).To(ContainSubstring("Config not created"))

			// Verify no config file was created
			Expect(".jitt.yaml").NotTo(BeAnExistingFile())
		})

		It("should refuse even when a project name is provided", func() {
			command := exec.Command(pathToJittBinary, "init", "MYPROJECT")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("Not inside a Git repo"))
			Expect(".jitt.yaml").NotTo(BeAnExistingFile())
		})

		It("should show helpful help message", func() {
			command := exec.Command(pathToJittBinary, "help")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring("jitt - Jira + Git + Tiny Tooling"))
			Expect(output).To(ContainSubstring("init [project]"))
			Expect(output).To(ContainSubstring("Initialize .jitt.yaml configuration file"))
		})
	})

	Context("inside a Git repository", func() {
		BeforeEach(func() {
			Expect(os.Mkdir(filepath.Join(tmpDir, ".git"), 0o755)).To(Succeed())
		})

		Context("with no existing config file", func() {
			It("should create .jitt.yaml file with empty project and show success message", func() {
				command := exec.Command(pathToJittBinary, "init")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				Expect(string(session.Out.Contents())).To(ContainSubstring(".jitt.yaml created"))

				// Verify file was created with YAML structure
				Expect(".jitt.yaml").To(BeAnExistingFile())
				content, err := os.ReadFile(".jitt.yaml")
				Expect(err).To(Succeed())
				Expect(string(content)).To(ContainSubstring("jira:"))
				Expect(string(content)).To(ContainSubstring("project: \"\""))
			})

			It("should create .jitt.yaml file with project when provided", func() {
				command := exec.Command(pathToJittBinary, "init", "TESTPROJ")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				Expect(string(session.Out.Contents())).To(ContainSubstring(".jitt.yaml created"))

				// Verify file has project configuration in YAML format
				content, err := os.ReadFile(".jitt.yaml")
				Expect(err).To(Succeed())
				Expect(string(content)).To(ContainSubstring("jira:"))
				Expect(string(content)).To(ContainSubstring("project: TESTPROJ"))
			})
		})

		Context("with existing .jitt.yaml file", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: existing"), 0o600)).To(Succeed())
			})

			It("should refuse to overwrite and show helpful error", func() {
				command := exec.Command(pathToJittBinary, "init")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(string(session.Err.Contents())).To(ContainSubstring(".jitt.yaml already exists"))
				Expect(string(session.Err.Contents())).To(ContainSubstring("not overwriting"))

				// Verify original content is preserved
				content, err := os.ReadFile(".jitt.yaml")
				Expect(err).To(Succeed())
				Expect(string(content)).To(Equal("jira:\n  project: existing"))
			})
		})
	})
})
