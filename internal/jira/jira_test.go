package jira

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var pathToJittBinary string

var _ = BeforeSuite(func() {
	var err error
	// Build the jitt binary for testing
	pathToJittBinary, err = gexec.Build("github.com/bbommarito/jitt/cmd/jitt")
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

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
		It("should refuse to create .jira file with helpful error", func() {
			command := exec.Command(pathToJittBinary, "init")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("Not inside a Git repo"))
			Expect(string(session.Err.Contents())).To(ContainSubstring(".jira not created"))

			// Verify no .jira file was created
			Expect(".jira").NotTo(BeAnExistingFile())
		})

		It("should refuse even when a project name is provided", func() {
			command := exec.Command(pathToJittBinary, "init", "MYPROJECT")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("Not inside a Git repo"))
			Expect(".jira").NotTo(BeAnExistingFile())
		})

		It("should show helpful help message", func() {
			command := exec.Command(pathToJittBinary, "help")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			output := string(session.Out.Contents())
			Expect(output).To(ContainSubstring("jitt - Jira + Git + Tiny Tooling"))
			Expect(output).To(ContainSubstring("init [project]"))
			Expect(output).To(ContainSubstring("Initialize .jira configuration file"))
		})
	})

	Context("inside a Git repository", func() {
		BeforeEach(func() {
			Expect(os.Mkdir(filepath.Join(tmpDir, ".git"), 0o755)).To(Succeed())
		})

		Context("with no existing .jira file", func() {
			It("should create .jira file with default content and show success message", func() {
				command := exec.Command(pathToJittBinary, "init")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				Expect(string(session.Out.Contents())).To(ContainSubstring(".jira created"))

				// Verify file was created with correct content
				Expect(".jira").To(BeAnExistingFile())
				content, err := os.ReadFile(".jira")
				Expect(err).To(Succeed())
				Expect(string(content)).To(Equal("# jitt config\n"))
			})

			It("should create .jira file with project when provided", func() {
				command := exec.Command(pathToJittBinary, "init", "TESTPROJ")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(0))
				Expect(string(session.Out.Contents())).To(ContainSubstring(".jira created"))

				// Verify file has project configuration
				content, err := os.ReadFile(".jira")
				Expect(err).To(Succeed())
				contentStr := strings.TrimSpace(string(content))
				Expect(contentStr).To(Equal(`project = "TESTPROJ"`))
			})
		})

		Context("with existing .jira file", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jira", []byte("existing config"), 0o600)).To(Succeed())
			})

			It("should refuse to overwrite and show helpful error", func() {
				command := exec.Command(pathToJittBinary, "init")
				session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
				Expect(err).NotTo(HaveOccurred())

				Eventually(session).Should(gexec.Exit(1))
				Expect(string(session.Err.Contents())).To(ContainSubstring(".jira already exists"))
				Expect(string(session.Err.Contents())).To(ContainSubstring("not overwriting"))

				// Verify original content is preserved
				content, err := os.ReadFile(".jira")
				Expect(err).To(Succeed())
				Expect(string(content)).To(Equal("existing config"))
			})
		})
	})
})
