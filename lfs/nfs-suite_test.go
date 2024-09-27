package lfs_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
)

func TestNfs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nfs Suite")
}

type (
	ensureTE struct {
		given     string
		should    string
		relative  string
		expected  string
		directory bool
	}

	RPEntry struct {
		given  string
		should string
		path   string
		expect string
	}
)

const (
	perm = 0o766
)

var (
	fakeHome      = filepath.Join(string(filepath.Separator), "home", "rabbitweed")
	fakeAbsCwd    = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music", "xpander")
	fakeAbsParent = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music")
)

func fakeHomeResolver() (string, error) {
	return fakeHome, nil
}

func fakeAbsResolver(path string) (string, error) {
	if strings.HasPrefix(path, "..") {
		return filepath.Join(fakeAbsParent, path[2:]), nil
	}

	if strings.HasPrefix(path, ".") {
		return filepath.Join(fakeAbsCwd, path[1:]), nil
	}

	return path, nil
}

type (
	mkDirAllMapFS struct {
		mapFS fstest.MapFS
	}
)

func (f *mkDirAllMapFS) FileExists(path string) bool {
	fi, err := f.mapFS.Stat(path)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return false
	}

	return true
}

func (f *mkDirAllMapFS) DirectoryExists(path string) bool {
	if strings.HasPrefix(path, string(filepath.Separator)) {
		path = path[1:]
	}

	fileInfo, err := f.mapFS.Stat(path)
	if err != nil {
		return false
	}

	if !fileInfo.IsDir() {
		return false
	}

	return true
}

func (f *mkDirAllMapFS) MkDirAll(path string, perm os.FileMode) error {
	var current string
	segments := filepath.SplitList(path)

	for _, part := range segments {
		if current == "" {
			current = part
		} else {
			current += string(filepath.Separator) + part
		}

		if exists := f.DirectoryExists(current); !exists {
			f.mapFS[current] = &fstest.MapFile{
				Mode: fs.ModeDir | perm,
			}
		}
	}

	return nil
}
