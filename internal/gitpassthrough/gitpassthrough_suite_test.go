package gitpassthrough

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGitpassthrough(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gitpassthrough Suite")
}
