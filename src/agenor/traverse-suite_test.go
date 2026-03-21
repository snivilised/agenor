package agenor_test

import (
	"errors"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/snivilised/jaywalk/src/agenor"
)

func TestTraverse(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Traverse Suite")
}

var (
	errBuildOptions = errors.New("options build error")
)

const (
	TreePath    = "traversal-tree-path"
	files       = 3
	directories = 2
)

var noOpHandler = func(_ agenor.Servant) error {
	return nil
}

type TestWriter struct {
	assertFn func()
}

func (tw *TestWriter) Write([]byte) (int, error) {
	tw.assertFn()
	return 0, nil
}
