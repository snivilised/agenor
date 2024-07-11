package tv_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	tv "github.com/snivilised/traverse"
)

func TestTraverse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Traverse Suite")
}

var (
	errBuildOptions = errors.New("options build error")
)

const (
	RootPath    = "traversal-root-path"
	RestorePath = "/from-restore-path"
	files       = 3
	folders     = 2
)

var noOpHandler = func(_ *tv.Node) error {
	return nil
}
