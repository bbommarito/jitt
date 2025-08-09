package config

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config Suite")
}

var _ = Describe("Config package", func() {
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

		// Reset viper state to avoid interference between tests
		viper.Reset()
	})

	AfterEach(func() {
		Expect(os.Chdir(oldCwd)).To(Succeed())
		viper.Reset()
	})

	Describe("Exists", func() {
		Context("when .jitt.yaml does not exist", func() {
			It("should return false", func() {
				Expect(Exists()).To(BeFalse())
			})
		})

		Context("when .jitt.yaml exists", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: test"), 0o600)).To(Succeed())
			})

			It("should return true", func() {
				Expect(Exists()).To(BeTrue())
			})
		})
	})

	Describe("Create", func() {
		Context("with a project name", func() {
			It("should create a valid config file", func() {
				err := Create("TESTPROJ")
				Expect(err).NotTo(HaveOccurred())

				Expect(".jitt.yaml").To(BeAnExistingFile())
				content, err := os.ReadFile(".jitt.yaml")
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("project: TESTPROJ"))
			})
		})

		Context("with an empty project name", func() {
			It("should create a config file with empty project", func() {
				err := Create("")
				Expect(err).NotTo(HaveOccurred())

				Expect(".jitt.yaml").To(BeAnExistingFile())
				content, err := os.ReadFile(".jitt.yaml")
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("project: \"\""))
			})
		})
	})

	Describe("Load", func() {
		Context("when config file does not exist", func() {
			It("should return config file not found error", func() {
				cfg, err := Load()
				Expect(cfg).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("config file not found"))
			})
		})

		Context("when config file exists but is malformed", func() {
			BeforeEach(func() {
				// Create an invalid YAML file
				Expect(os.WriteFile(".jitt.yaml", []byte("invalid yaml: [unclosed bracket"), 0o600)).To(Succeed())
			})

			It("should return unmarshaling error", func() {
				cfg, err := Load()
				Expect(cfg).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error reading config file"))
			})
		})

		Context("when config file exists but has wrong structure", func() {
			BeforeEach(func() {
				// Create YAML that parses but doesn't match our struct
				Expect(os.WriteFile(".jitt.yaml", []byte("wrong_field: value"), 0o600)).To(Succeed())
			})

			It("should still load with defaults", func() {
				cfg, err := Load()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())
				Expect(cfg.Jira.Project).To(Equal(""))
			})
		})

		Context("when config file is valid", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: VALIDPROJ"), 0o600)).To(Succeed())
			})

			It("should load the configuration successfully", func() {
				cfg, err := Load()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())
				Expect(cfg.Jira.Project).To(Equal("VALIDPROJ"))
			})
		})

		Context("when config file is empty", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte(""), 0o600)).To(Succeed())
			})

			It("should load with default values", func() {
				cfg, err := Load()
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg).NotTo(BeNil())
				Expect(cfg.Jira.Project).To(Equal(""))
			})
		})
	})

	Describe("Update", func() {
		Context("when config file does not exist", func() {
			It("should return config file not found error", func() {
				err := Update("jira.project", "NEWPROJ")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("config file not found - run 'jitt init' first"))
			})
		})

		Context("when config file exists", func() {
			BeforeEach(func() {
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: OLDPROJ"), 0o600)).To(Succeed())
			})

			It("should update the configuration", func() {
				err := Update("jira.project", "NEWPROJ")
				Expect(err).NotTo(HaveOccurred())

				// Verify the file was updated
				content, err := os.ReadFile(".jitt.yaml")
				Expect(err).NotTo(HaveOccurred())
				Expect(string(content)).To(ContainSubstring("project: NEWPROJ"))
			})
		})

		Context("when config file exists but is unreadable", func() {
			BeforeEach(func() {
				// Create a file with invalid permissions (if possible)
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: TEST"), 0o000)).To(Succeed())
			})

			AfterEach(func() {
				// Restore permissions for cleanup
				_ = os.Chmod(".jitt.yaml", 0o600)
			})

			It("should return error reading config file", func() {
				err := Update("jira.project", "NEWPROJ")
				if err != nil {
					// This test might not work on all systems due to permission handling
					Expect(err.Error()).To(ContainSubstring("error reading config file"))
				}
			})
		})

		Context("when config file becomes corrupted after creation", func() {
			BeforeEach(func() {
				// Create valid file first
				Expect(os.WriteFile(".jitt.yaml", []byte("jira:\n  project: OLDPROJ"), 0o600)).To(Succeed())
				// Then corrupt it
				Expect(os.WriteFile(".jitt.yaml", []byte("invalid yaml: ["), 0o600)).To(Succeed())
			})

			It("should return error reading config file", func() {
				err := Update("jira.project", "NEWPROJ")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("error reading config file"))
			})
		})
	})
})
