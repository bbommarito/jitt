package jira

import (
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
