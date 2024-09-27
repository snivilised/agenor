package lfs

import (
	"io/fs"
	"os"
)

type localFS struct {
	fsys fs.FS
}

// NewLocalFS creates a native file system.
func NewLocalFS(path string) fs.ReadDirFS {
	return &localFS{
		fsys: os.DirFS(path),
	}
}

func (n *localFS) Open(path string) (fs.File, error) {
	return n.fsys.Open(path)
}

func (n *localFS) ReadDir(name string) ([]fs.DirEntry, error) {
	return fs.ReadDir(n.fsys, name)
}

type queryStatusFS struct {
	fsys fs.FS
}

// NewQueryStatusFS defines a file system that has a Stat
// method to determine the existence of a path.
func NewQueryStatusFS(fsys fs.FS) fs.StatFS {
	return &queryStatusFS{
		fsys: fsys,
	}
}

func (q *queryStatusFS) Open(name string) (fs.File, error) {
	return q.fsys.Open(name)
}

func (q *queryStatusFS) Stat(name string) (fs.FileInfo, error) {
	return os.Stat(name)
}
