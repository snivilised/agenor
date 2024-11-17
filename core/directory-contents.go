package core

import (
	"io/fs"
)

type (
	// DirectoryContents represents the contents of a directory's contents and
	// handles sorting order which by default is different between various
	// operating systems. This abstraction removes the differences in sorting
	// behaviour on different platforms.
	DirectoryContents interface {
		All() []fs.DirEntry
		Directories() []fs.DirEntry
		Files() []fs.DirEntry
	}
)
