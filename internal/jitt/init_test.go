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
