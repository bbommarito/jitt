package jira_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/bbommarito/jitt/internal/jira"
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
				Expect(jira.HasJiraFile()).To(BeFalse())
			})
			
			It("should indicate .jira file is absent", func() {
				Expect(".jira").NotTo(BeAnExistingFile())
			})
		})

		Context("when .jira file exists", func() {
			BeforeEach(func() {
				// Write file and expect it to succeed  
				Expect(os.WriteFile(".jira", []byte("test"), 0644)).To(Succeed())
			})

			It("should return true", func() {
				Expect(jira.HasJiraFile()).To(BeTrue())
			})
			
			It("should have the .jira file present", func() {
				Expect(".jira").To(BeAnExistingFile())
			})
			
			It("should contain the expected content", func() {
				Expect(".jira").To(BeAnExistingFile())
				Expect(os.ReadFile(".jira")).To(Equal([]byte("test")))
			})
		})
	})
})
