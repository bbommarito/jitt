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

var _ = Describe("Integration: jitt binary behavior", func() {
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

		// Change to temp directory (no git repo)
		Expect(os.Chdir(tmpDir)).To(Succeed())
	})

	AfterEach(func() {
		// Restore original directory
		Expect(os.Chdir(oldCwd)).To(Succeed())
	})

	Describe("when running outside a Git repository", func() {
		It("should fail to create .jira file with proper error message", func() {
			command := exec.Command(pathToJittBinary, "init")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("Not inside a Git repo"))
			Expect(string(session.Err.Contents())).To(ContainSubstring(".jira not created"))

			// Verify no .jira file was created
			Expect(".jira").NotTo(BeAnExistingFile())
		})

		It("should fail even with project name provided", func() {
			command := exec.Command(pathToJittBinary, "init", "TESTPROJ")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("Not inside a Git repo"))
			Expect(string(session.Err.Contents())).To(ContainSubstring(".jira not created"))

			// Verify no .jira file was created
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

	Describe("when running inside a Git repository", func() {
		BeforeEach(func() {
			// Create .git directory to simulate being in a git repo
			Expect(os.Mkdir(filepath.Join(tmpDir, ".git"), 0o755)).To(Succeed())
		})

		It("should successfully create .jira file", func() {
			command := exec.Command(pathToJittBinary, "init")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring(".jira created"))

			// Verify .jira file was created with correct content
			Expect(".jira").To(BeAnExistingFile())
			content, err := os.ReadFile(".jira")
			Expect(err).To(Succeed())
			Expect(string(content)).To(Equal("# jitt config\n"))
		})

		It("should create .jira file with project configuration", func() {
			command := exec.Command(pathToJittBinary, "init", "MYPROJ")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(0))
			Expect(string(session.Out.Contents())).To(ContainSubstring(".jira created"))

			// Verify .jira file was created with project content
			Expect(".jira").To(BeAnExistingFile())
			content, err := os.ReadFile(".jira")
			Expect(err).To(Succeed())
			contentStr := strings.TrimSpace(string(content))
			Expect(contentStr).To(Equal(`project = "MYPROJ"`))
		})

		It("should refuse to overwrite existing .jira file", func() {
			// Create existing .jira file
			Expect(os.WriteFile(".jira", []byte("existing content"), 0o600)).To(Succeed())

			command := exec.Command(pathToJittBinary, "init")
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			Eventually(session).Should(gexec.Exit(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring(".jira already exists"))
			Expect(string(session.Err.Contents())).To(ContainSubstring("not overwriting"))

			// Verify original content preserved
			content, err := os.ReadFile(".jira")
			Expect(err).To(Succeed())
			Expect(string(content)).To(Equal("existing content"))
		})
	})
})
