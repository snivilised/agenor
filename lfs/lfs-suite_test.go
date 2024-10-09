package lfs_test

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/locale"
)

func TestLfs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Lfs Suite")
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

	funcFS[T any] func(entry fsTE[T], fS T)

	fsTE[T any] struct {
		given     string
		should    string
		note      string
		op        string
		overwrite bool
		directory bool
		require   string
		target    string
		from      string
		to        string
		arrange   funcFS[T]
		action    funcFS[T]
	}
)

func (t *fsTE[T]) run(fS T) {
	if t.arrange != nil {
		t.arrange(*t, fS)
	}
	t.action(*t, fS)
}

var (
	fakeHome      = filepath.Join(string(filepath.Separator), "home", "rabbitweed")
	fakeAbsCwd    = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music", "xpander")
	fakeAbsParent = filepath.Join(string(filepath.Separator), "home", "rabbitweed", "music")
)

// require ensures that a path exists. If files are also provided,
// it will create these files too. The files are relative to the root
// and should be prefixed by parent; that is to say, when a test needs
// scratch/foo.txt, parent = 'scratch' and file = 'scratch/foo.txt';
// ie te file still needs to be relative to root, not parent.
func require(root, parent string, files ...string) error {
	if err := os.MkdirAll(filepath.Join(root, parent), lab.Perms.Dir.Perm()); err != nil {
		return fmt.Errorf("failed to create directory: %q (%w)", parent, err)
	}

	for _, name := range files {
		handle, err := os.Create(filepath.Join(root, name))
		if err != nil {
			return fmt.Errorf("failed to create file: %q (%w)", name, err)
		}

		handle.Close()
	}

	return nil
}

func scratch(root string) {
	scratchPath := filepath.Join(root, lab.Static.FS.Scratch)

	if _, err := os.Stat(scratchPath); err == nil {
		Expect(os.RemoveAll(scratchPath)).To(Succeed(),
			fmt.Sprintf("failed to delete existing directory %q", scratchPath),
		)
	}
}

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

func IsLinkError(err error, reason string) {
	var linkErr *os.LinkError
	Expect(errors.As(err, &linkErr)).To(BeTrue(), fmt.Sprintf("not LinkError, %q", reason))
}

func IsSameDirMoveRejectionError(err error, reason string) {
	Expect(locale.IsRejectSameDirMoveError(err)).To(BeTrue(),
		fmt.Sprintf("not SameDirMoveRejectionError, %q", reason),
	)
}

type (
	makeDirMapFS struct {
		mapFS fstest.MapFS
	}
)

func (f *makeDirMapFS) FileExists(path string) bool {
	fi, err := f.mapFS.Stat(path)
	if err != nil {
		return false
	}

	if fi.IsDir() {
		return false
	}

	return true
}

func (f *makeDirMapFS) DirectoryExists(path string) bool {
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

func (f *makeDirMapFS) MakeDir(path string, perm os.FileMode) error {
	if exists := f.DirectoryExists(path); !exists {
		f.mapFS[path] = &fstest.MapFile{
			Mode: fs.ModeDir | perm,
		}
	}

	return nil
}

func (f *makeDirMapFS) MakeDirAll(path string, perm os.FileMode) error {
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
