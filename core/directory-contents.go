package core

import (
	"io/fs"
)

// DirectoryContents represents the contents of a directory's contents and
// handles sorting order which by default is different between various
// operating systems. This abstraction removes the differences in sorting
// behaviour on different platforms.
type (
	DirectoryContents interface {
		All() []fs.DirEntry
		Directories() []fs.DirEntry
		Files() []fs.DirEntry
	}
)
