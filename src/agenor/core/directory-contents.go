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
		// All returns all directory entries in a consistent order across platforms.
		All() []fs.DirEntry

		// Directories returns only the directory entries that are directories.
		Directories() []fs.DirEntry

		// Files returns only the directory entries that are files.
		Files() []fs.DirEntry
	}
)
