package lab

import (
	"os"
	"testing/fstest"
)

type TestTraverseFS struct {
	fstest.MapFS
}

func (f *TestTraverseFS) FileExists(path string) bool {
	if mapFile, found := f.MapFS[path]; found && !mapFile.Mode.IsDir() {
		return true
	}

	return false
}

func (f *TestTraverseFS) DirectoryExists(path string) bool {
	if mapFile, found := f.MapFS[path]; found && mapFile.Mode.IsDir() {
		return true
	}

	return false
}

func (f *TestTraverseFS) Create(name string) (*os.File, error) {
	_ = name
	panic("NOT-IMPL: TestTraverseFS.Create")
}

func (f *TestTraverseFS) MkDirAll(path string, perm os.FileMode) error {
	_ = path
	_ = perm
	panic("NOT-IMPL: TestTraverseFS.MkDirAll")
}

func (f *TestTraverseFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	_ = name
	_ = data
	_ = perm

	panic("NOT-IMPL: TestTraverseFS.WriteFile")
}
