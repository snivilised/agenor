package kernel_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
)

func TestKernel(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Kernel Suite")
}

const (
	RootPath    = "traversal-tree-path"
	RestorePath = "/from-restore-path"
)
