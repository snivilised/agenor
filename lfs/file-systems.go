package lfs

import (
	"io/fs"
	"os"
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

// 🎯 openFS
type openFS struct {
	fsys fs.FS
}

func (f *openFS) Open(path string) (fs.File, error) {
	return f.fsys.Open(path)
}

// 🎯 statFS
type statFS struct {
	openFS
}

func NewStatFS(path string) fs.StatFS {
	ents := compose(path)
	return &ents.stat
}

func (f *statFS) Stat(path string) (fs.FileInfo, error) {
	return fs.Stat(f.fsys, path)
}

// 🧩 ---> file system query

// 🎯 readDirFS
type readDirFS struct {
	openFS
}

// NewReadDirFS creates a native file system.
func NewReadDirFS(path string) fs.ReadDirFS {
	ents := compose(path)

	return &readDirFS{
		openFS: ents.open,
	}
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
func (n *readDirFS) Open(path string) (fs.File, error) {
	return n.fsys.Open(path)
}

// ReadDir reads the named directory
// and returns a list of directory entries sorted by filename.
//
// If fs implements [ReadDirFS], ReadDir calls fs.ReadDir.
// Otherwise ReadDir calls fs.Open and uses ReadDir and Close
// on the returned file.
func (n *readDirFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(n.fsys, name)
}

// 🎯 queryStatusFS
type queryStatusFS struct {
	statFS
	readDirFS
}

func NewQueryStatusFS(path string) fs.StatFS {
	ents := compose(path)

	return &queryStatusFS{
		readDirFS: ents.read,
	}
}

// QueryStatusFromFS defines a file system that has a Stat
// method to determine the existence of a path.
func QueryStatusFromFS(fsys fs.FS) fs.StatFS {
	return &queryStatusFS{
		readDirFS: readDirFS{
			openFS: openFS{
				fsys: fsys,
			},
		},
	}
}

func (q *queryStatusFS) Open(name string) (fs.File, error) {
	return q.statFS.fsys.Open(name)
}

// Stat returns a [FileInfo] describing the named file.
// If there is an error, it will be of type [*PathError].
func (q *queryStatusFS) Stat(name string) (fs.FileInfo, error) {
	return q.statFS.Stat(name)
}

// 🎯 existsInFS
type existsInFS struct {
	queryStatusFS
}

// ExistsInFS
func NewExistsInFS(path string) ExistsInFS {
	ents := compose(path)

	return &ents.exists
}

// FileExists does file exist at the path specified
func (f *existsInFS) FileExists(path string) bool {
	info, err := f.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

// DirectoryExists does directory exist at the path specified
func (f *existsInFS) DirectoryExists(path string) bool {
	info, err := f.Stat(path)
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
	queryStatusFS
}

func NewReadFileFS(path string) ReadFileFS {
	ents := compose(path)

	return &readFileFS{
		queryStatusFS: ents.query,
	}
}

// ReadFile reads the named file from the file system fs and returns its contents.
// A successful call returns a nil error, not [io.EOF].
// (Because ReadFile reads the whole file, the expected EOF
// from the final Read is not treated as an error to be reported.)
//
// If fs implements [ReadFileFS], ReadFile calls fs.ReadFile.
// Otherwise ReadFile calls fs.Open and uses Read and Close
// on the returned [File].
func (f *readFileFS) ReadFile(path string) ([]byte, error) {
	return fs.ReadFile(f.queryStatusFS.statFS.fsys, path)
}

// 🧩 ---> file system mutation

// 🎯 baseWriterFS
type baseWriterFS struct {
	openFS
	overwrite bool
}

// 🎯 MkDirAllFS
type mkDirAllFS struct {
	existsInFS
}

// NewMkDirAllFS
func NewMkDirAllFS() MkDirAllFS {
	return &mkDirAllFS{}
}

// MkdirAll creates a directory named path,
// along with any necessary parents, and returns nil,
// or else returns an error.
// The permission bits perm (before umask) are used for all
// directories that MkdirAll creates.
// If path is already a directory, MkdirAll does nothing
// and returns nil.
func (*mkDirAllFS) MkDirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm) // !!! WRONG use the correct fs
}

// 🎯 writeFileFS
type writeFileFS struct {
	baseWriterFS
}

func NewWriteFileFS(path string, overwrite bool) WriteFileFS {
	ents := compose(path)

	return &writeFileFS{
		baseWriterFS: baseWriterFS{
			openFS:    ents.open,
			overwrite: overwrite,
		},
	}
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
	return os.Create(name)
}

// WriteFile writes data to the named file, creating it if necessary.
// If the file does not exist, WriteFile creates it with permissions perm (before umask);
// otherwise WriteFile truncates it before writing, without changing permissions.
// Since WriteFile requires multiple system calls to complete, a failure mid-operation
// can leave the file in a partially written state.
func (f *writeFileFS) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// 🧩 ---> file system aggregators

// 🎯 readerFS
type readerFS struct {
	readDirFS
	readFileFS
	existsInFS
	statFS
}

// NewReaderFS
func NewReaderFS(path string) ReaderFS {
	ents := compose(path)

	return &readerFS{
		readDirFS: ents.read,
		readFileFS: readFileFS{
			queryStatusFS: ents.query,
		},
		existsInFS: ents.exists,
		statFS:     ents.stat,
	}
}

// 🎯 writerFS
type writerFS struct {
	mkDirAllFS
	writeFileFS
}

func NewWriterFS(path string, overwrite bool) WriterFS {
	ents := compose(path)

	return &writerFS{
		writeFileFS: writeFileFS{
			baseWriterFS: baseWriterFS{
				openFS:    ents.open,
				overwrite: overwrite,
			},
		},
	}
}

// 🎯 traverseFS
type traverseFS struct {
	readerFS
	writerFS
}

func NewTraverseFS(path string, overwrite bool) TraverseFS {
	ents := compose(path)

	return &traverseFS{
		readerFS: readerFS{
			readDirFS: ents.read,
			readFileFS: readFileFS{
				queryStatusFS: ents.query,
			},
			existsInFS: ents.exists,
			statFS:     ents.stat,
		},
		writerFS: writerFS{
			mkDirAllFS: mkDirAllFS{
				existsInFS: ents.exists,
			},
			writeFileFS: writeFileFS{
				baseWriterFS: baseWriterFS{
					openFS:    ents.open,
					overwrite: overwrite,
				},
			},
		},
	}
}

type entities struct {
	open   openFS
	read   readDirFS
	stat   statFS
	query  queryStatusFS
	exists existsInFS
}

func compose(root string) *entities {
	open := openFS{
		fsys: os.DirFS(root),
	}
	read := readDirFS{
		openFS: open,
	}
	stat := statFS{
		openFS: open,
	}
	query := queryStatusFS{
		statFS: statFS{
			openFS: open,
		},
		readDirFS: read,
	}
	exists := existsInFS{
		queryStatusFS: query,
	}

	return &entities{
		open:   open,
		read:   read,
		stat:   stat,
		query:  query,
		exists: exists,
	}
}
