package lab

import (
	"io/fs"
	"os"
	"strings"
	"testing/fstest"

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/internal/third/lo"
)

type testMapFile struct {
	f fstest.MapFile
}

type TestTraverseFS struct {
	fstest.MapFS
}

func (f *TestTraverseFS) FileExists(name string) bool {
	if mapFile, found := f.MapFS[name]; found && !mapFile.Mode.IsDir() {
		return true
	}

	return false
}

func (f *TestTraverseFS) DirectoryExists(name string) bool {
	if mapFile, found := f.MapFS[name]; found && mapFile.Mode.IsDir() {
		return true
	}

	return false
}

func (f *TestTraverseFS) Create(name string) (*os.File, error) {
	if _, err := f.Stat(name); err == nil {
		return nil, fs.ErrExist
	}

	file := &fstest.MapFile{
		Mode: Perms.File,
	}

	f.MapFS[name] = file
	// TODO: this needs a resolution using a file interface
	// rather than using os.File which is a struct not an
	// interface
	dummy := &os.File{}

	return dummy, nil
}

func (f *TestTraverseFS) MakeDir(name string, perm os.FileMode) error {
	if !fs.ValidPath(name) {
		return nef.NewInvalidPathError("MakeDir", name)
	}

	if _, found := f.MapFS[name]; !found {
		f.MapFS[name] = &fstest.MapFile{
			Mode: perm | os.ModeDir,
		}
	}

	return nil
}

func (f *TestTraverseFS) MakeDirAll(name string, perm os.FileMode) error {
	if !fs.ValidPath(name) {
		return nef.NewInvalidPathError("MakeDir", name)
	}

	segments := strings.Split(name, "/")

	_ = lo.Reduce(segments,
		func(acc []string, s string, _ int) []string {
			acc = append(acc, s)
			path := strings.Join(acc, "/")

			if _, found := f.MapFS[path]; !found {
				f.MapFS[path] = &fstest.MapFile{
					Mode: perm | os.ModeDir,
				}
			}

			return acc
		}, []string{},
	)

	return nil
}

func (f *TestTraverseFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	if _, err := f.Stat(name); err == nil {
		return fs.ErrExist
	}

	f.MapFS[name] = &fstest.MapFile{
		Data: data,
		Mode: perm,
	}

	return nil
}

func (f *TestTraverseFS) Ensure(_ nef.PathAs) (string, error) {
	panic("NOT-IMPL: TestTraverseFS.Ensure")
}
