package lfs

import (
	"io/fs"
)

type syntheticEntry struct {
	fi fs.FileInfo
}

func (e *syntheticEntry) Name() string {
	return e.fi.Name()
}

func (e *syntheticEntry) IsDir() bool {
	return e.fi.IsDir()
}

func (e *syntheticEntry) Type() fs.FileMode {
	return e.fi.Mode()
}

func (e *syntheticEntry) Info() (fs.FileInfo, error) {
	return e.fi, nil
}

// FromFileInfo synthesises a DirEntry instance from a FileInfo
func FromFileInfo(fi fs.FileInfo) fs.DirEntry {
	return &syntheticEntry{
		fi: fi,
	}
}
