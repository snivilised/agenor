package prism_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPrism(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Prism Suite")
}
