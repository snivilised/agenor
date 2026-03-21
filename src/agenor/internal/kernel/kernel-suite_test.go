package kernel_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestKernel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kernel Suite")
}

const (
	RootPath = "traversal-tree-path"
)
