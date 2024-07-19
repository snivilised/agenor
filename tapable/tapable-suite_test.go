package tapable_test

import (
	"io/fs"
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
)

func TestTapable(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tapable Suite")
}

type fakeDirEntry struct {
	name string
	dir  bool
	mode fs.FileMode
	info fs.FileInfo
}

func (e *fakeDirEntry) Name() string {
	return e.name
}

func (e *fakeDirEntry) IsDir() bool {
	return e.dir
}

func (e *fakeDirEntry) Type() fs.FileMode {
	return e.mode
}

func (e *fakeDirEntry) Info() (fs.FileInfo, error) {
	return e.info, nil
}
