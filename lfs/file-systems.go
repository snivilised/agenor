package lfs

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/snivilised/traverse/locale"
)

// 🔥 An important note about using standard golang file systems (io.fs/fs.FS)
// as opposed to using the native os calls directly (os.XXX).
// Note that (io.fs/fs.FS) represents a virtual file system where as os.XXX
// represent operations on the local file system. Working with either of
// these is fundamentally different to working with the other; bear this in
// mind to avoid confusion.
//
// virtual file system
// ===================
//
// The client is expected to create a file system rooted at a particular path.
// This path must be absolute. Any function call on the resulting file system
// that requires a path must be relative to this root and therefore must not
// begin or end with a slash.
//
// When composing paths to use with a file system, one might think that using
// filepath.Separator and building paths with filepath.Join is the most
// prudent thing to do to ensure correct functioning on different platforms. When
// it comes to file systems, this is most certainly not the case. The paths
// are virtual and they are mapped into an underlying file system, which typically
// is the local file system. This means that paths used only need to use '/'. And
// the silly thing is, characters like ':', or '\' for windows should not be
// treated as separators by the underlying file system. So really using
// filepath.Separator with a virtual file system is not valid.
//

func sanitise(root string) string {
	return root
}

// 🧩 ---> open

// 🎯 openFS
type openFS struct {
	fS   fs.FS
	root string
}

func (f *openFS) Open(name string) (fs.File, error) {
	return f.fS.Open(name)
}

// 🧩 ---> stat

// 🎯 statFS
type statFS struct {
	*openFS
}

func NewStatFS(at At) fs.StatFS {
	ents := compose(sanitise(at.Root))
	return &ents.stat
}

func (f *statFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(filepath.Join(f.root, name))
}

// 🧩 ---> file system query

// 🎯 readDirFS
type readDirFS struct {
	*openFS
}

// NewReadDirFS creates a file system with read directory capability
func NewReadDirFS(at At) fs.ReadDirFS {
	ents := compose(sanitise(at.Root))

	return &ents.exists
}

// Open opens the named file.
//
// When Open returns an error, it should be of type *PathError
// with the Op field set to "open", the Path field set to name,
// and the Err field describing the problem.
//
// Open should reject attempts to open names that do not satisfy
// ValidPath(name), returning a *PathError with Err set to
// ErrInvalid or ErrNotExist.
func (n *readDirFS) Open(name string) (fs.File, error) {
	return n.fS.Open(name)
}

// ReadDir reads the named directory
// and returns a list of directory entries sorted by filename.
//
// If fs implements [ReadDirFS], ReadDir calls fs.ReadDir.
// Otherwise ReadDir calls fs.Open and uses ReadDir and Close
// on the returned file.
func (n *readDirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(n.fS, name)
}

// 🎯 queryStatusFS
type queryStatusFS struct {
	*statFS
	*readDirFS
}

func (q *queryStatusFS) Open(name string) (fs.File, error) {
	return q.statFS.fS.Open(name)
}

// Stat returns a [FileInfo] describing the named file.
// If there is an error, it will be of type [*PathError].
func (q *queryStatusFS) Stat(name string) (fs.FileInfo, error) {
	return q.statFS.Stat(name)
}

// 🎯 existsInFS
type existsInFS struct {
	*queryStatusFS
}

// ExistsInFS
func NewExistsInFS(at At) ExistsInFS {
	ents := compose(sanitise(at.Root))

	return &ents.exists
}

// FileExists does file exist at the path specified
func (f *existsInFS) FileExists(name string) bool {
	info, err := f.Stat(name)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

// DirectoryExists does directory exist at the path specified
func (f *existsInFS) DirectoryExists(name string) bool {
	info, err := f.Stat(name)
	if err != nil {
		return false
	}

	if !info.IsDir() {
		return false
	}

	return true
}

// 🎯 readFileFS
type readFileFS struct {
	*queryStatusFS
}

func NewReadFileFS(at At) ReadFileFS {
	ents := compose(sanitise(at.Root))

	return &ents.reader
}

// ReadFile reads the named file from the file system fs and returns its contents.
// A successful call returns a nil error, not [io.EOF].
// (Because ReadFile reads the whole file, the expected EOF
// from the final Read is not treated as an error to be reported.)
//
// If fs implements [ReadFileFS], ReadFile calls fs.ReadFile.
// Otherwise ReadFile calls fs.Open and uses Read and Close
// on the returned [File].
func (f *readFileFS) ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(f.queryStatusFS.statFS.fS, name)
}

// 🧩 ---> file system mutation

// 🎯 copyFS

type copyFS struct {
	*openFS
}

func (f *copyFS) Copy(from, to string) error {
	return fmt.Errorf("copy not implemented yet (from: %q, to: %q)", from, to)
}

// CopyFS copies the file system fsys into the directory dir,
// creating dir if necessary.
func (f *copyFS) CopyFS(dir string, fsys fs.FS) error {
	_ = fsys
	return fmt.Errorf("copyFS not implemented yet (dir: %q)", dir)
}

// 🎯 baseWriterFS
type baseWriterFS struct {
	*openFS
	*existsInFS
	overwrite bool
}

// 🎯 MakeDirFS
type makeDirAllFS struct {
	*existsInFS
}

// NewMakeDirFS
func NewMakeDirFS(at At) MakeDirFS {
	ents := compose(sanitise(at.Root)).mutator(at.Overwrite)

	return &ents.writer
}

// Mkdir creates a new directory with the specified name and permission
// bits (before umask).
// If there is an error, it will be of type *PathError.
func (f *makeDirAllFS) MakeDir(name string, perm os.FileMode) error {
	if !fs.ValidPath(name) {
		return locale.NewInvalidPathError(name)
	}

	if f.DirectoryExists(name) {
		return nil
	}

	path := filepath.Join(f.statFS.root, name)
	return os.Mkdir(path, perm)
}

// MakeDirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates.
// If path is already a directory, MakeDirAll does nothing
// and returns nil.
func (f *makeDirAllFS) MakeDirAll(name string, perm os.FileMode) error {
	if !fs.ValidPath(name) {
		return locale.NewInvalidPathError(name)
	}

	if f.DirectoryExists(name) {
		return nil
	}
	path := filepath.Join(f.statFS.root, name)
	return os.MkdirAll(path, perm)
}

// 🎯 removeFS

type removeFS struct {
	*openFS
}

func (f *removeFS) Remove(name string) error {
	if !fs.ValidPath(name) {
		return locale.NewInvalidPathError(name)
	}

	path := filepath.Join(f.root, filepath.Clean(name))
	return os.Remove(path)
}

func (f *removeFS) RemoveAll(path string) error {
	if !fs.ValidPath(path) {
		return locale.NewInvalidPathError(path)
	}

	return os.RemoveAll(filepath.Join(f.root, filepath.Clean(path)))
}

// 🎯 renameFS

type renameFS struct {
	*openFS
}

// Rename delegates to the Rename functionality implemented in the standard
// library.
func (s *renameFS) Rename(from, to string) error {
	return os.Rename(
		filepath.Join(s.root, from),
		filepath.Join(s.root, to),
	)
}

// 🎯 writeFileFS
type writeFileFS struct {
	*baseWriterFS
}

func NewWriteFileFS(at At) WriteFileFS {
	ents := compose(sanitise(at.Root)).mutator(at.Overwrite)

	return &ents.writer
}

// Create creates or truncates the named file. If the file already exists,
// it is truncated. If the file does not exist, it is created with mode 0o666
// (before umask). If successful, methods on the returned File can
// be used for I/O; the associated file descriptor has mode O_RDWR.
// If there is an error, it will be of type *PathError.
//
// We need to maintain conformity with apis in the standard library. Ideally,
// this Create method would have the overwrite bool passed in as an argument,
// but doing so would break standard lib compatibility. Instead, the underlying
// implementation has to decide wether to Create on an override basis itself.
// The disadvantage of this approach is that the client can not decide on
// the fly wether a call to Create is on a override basis or not. This decision
// has to be made at the point of creating the file system. This is less
// flexible and just results in friction, but this is out of our power.
func (f *writeFileFS) Create(name string) (*os.File, error) {
	if !fs.ValidPath(name) {
		return nil, locale.NewInvalidPathError(name)
	}

	if !f.overwrite && f.FileExists(name) {
		return nil, os.ErrExist
	}

	path := filepath.Join(f.root, name)
	return os.Create(path)
}

// WriteFile writes data to the named file, creating it if necessary.
// If the file does not exist, WriteFile creates it with permissions perm (before umask);
// otherwise WriteFile truncates it before writing, without changing permissions.
// Since WriteFile requires multiple system calls to complete, a failure mid-operation
// can leave the file in a partially written state.
func (f *writeFileFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	if !fs.ValidPath(name) {
		return locale.NewInvalidPathError(name)
	}

	path := filepath.Join(f.root, name)
	return os.WriteFile(path, data, perm)
}

// 🧩 ---> file system aggregators

// 🎯 readerFS
type readerFS struct {
	*existsInFS
	*readDirFS
	*readFileFS
	*statFS
}

// NewReaderFS
func NewReaderFS(at At) ReaderFS {
	ents := compose(sanitise(at.Root))

	return &ents.reader
}

// 🎯 moverFS
type moverFS struct {
	mover lazyMover
}

// 🎯 writerFS
type writerFS struct {
	*copyFS
	*makeDirAllFS
	*moverFS
	*removeFS
	*renameFS
	*writeFileFS
}

func NewWriterFS(at At) WriterFS {
	ents := compose(sanitise(at.Root)).mutator(at.Overwrite)

	return &ents.writer
}

// 🎯 mutatorFS
type mutatorFS struct {
	*readerFS
	*writerFS
}

// Move is similar to rename but it has distinctly different semantics, which
// also varies depending on whether the file system was created with overwrite
// enabled or not.
// When overwrite is enabled, a move with overwrite the destination. If not
// enabled, Move will return as error (os.ErrExist).
// The paths denoted by from and to must be in different locations, otherwise
// the move amounts to a rename and the client should use Rename instead of
// move. When this scenario is detected, an error is returned.
func (f *mutatorFS) Move(from, to string) error {
	return f.mover.instance(f.baseWriterFS.root, f.overwrite, f).move(from, to)
}

func NewTraverseFS(at At) TraverseFS {
	ents := compose(sanitise(at.Root)).mutator(at.Overwrite)

	return &mutatorFS{
		readerFS: &ents.reader,
		writerFS: &ents.writer,
	}
}

func NewUniversalFS(at At) UniversalFS {
	ents := compose(sanitise(at.Root)).mutator(at.Overwrite)

	return &mutatorFS{
		readerFS: &ents.reader,
		writerFS: &ents.writer,
	}
}

// 🧩 ---> construction

type (
	entities struct {
		open   openFS
		read   readDirFS
		stat   statFS
		query  queryStatusFS
		exists existsInFS
		reader readerFS
		copy   copyFS
		remove removeFS
		rename renameFS
		writer writerFS
	}
)

func (e *entities) mutator(overwrite bool) *entities {
	e.writer = writerFS{
		copyFS: &copyFS{
			openFS: &e.open,
		},
		makeDirAllFS: &makeDirAllFS{
			existsInFS: &e.exists,
		},
		moverFS: &moverFS{},
		removeFS: &removeFS{
			openFS: &e.open,
		},
		renameFS: &renameFS{
			openFS: &e.open,
		},
		writeFileFS: &writeFileFS{
			baseWriterFS: &baseWriterFS{
				existsInFS: &e.exists,
				openFS:     &e.open,
				overwrite:  overwrite,
			},
		},
	}

	return e
}

func compose(root string) *entities {
	open := openFS{
		fS:   os.DirFS(root),
		root: root,
	}
	read := readDirFS{
		openFS: &open,
	}
	stat := statFS{
		openFS: &open,
	}
	query := queryStatusFS{
		statFS: &statFS{
			openFS: &open,
		},
		readDirFS: &read,
	}
	exists := existsInFS{
		queryStatusFS: &query,
	}

	reader := readerFS{
		readDirFS: &read,
		readFileFS: &readFileFS{
			queryStatusFS: &query,
		},
		existsInFS: &exists,
		statFS:     &stat,
	}

	return &entities{
		open:   open,
		read:   read,
		stat:   stat,
		query:  query,
		exists: exists,
		reader: reader,
	}
}
