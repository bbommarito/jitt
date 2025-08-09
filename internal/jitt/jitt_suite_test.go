package jitt

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJitt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Jitt Suite")
}
