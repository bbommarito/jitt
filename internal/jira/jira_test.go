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
			tmpDir   string
			oldCwd   string
		)

		BeforeEach(func() {
			tmpDir = GinkgoT().TempDir()
			var err error
			oldCwd, err = os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			err = os.Chdir(tmpDir)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			err := os.Chdir(oldCwd)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when no .jira file exists", func() {
			It("should return false", func() {
				Expect(jira.HasJiraFile()).To(BeFalse())
			})
		})

		Context("when .jira file exists", func() {
			BeforeEach(func() {
				err := os.WriteFile(".jira", []byte("test"), 0644)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return true", func() {
				Expect(jira.HasJiraFile()).To(BeTrue())
			})
		})
	})
})
