package core

import (
	"io/fs"

	"github.com/snivilised/traverse/enums"
)

// DirectoryContents represents the contents of a directory's contents and
// handles sorting order which by default is different between various
// operating systems. This abstraction removes the differences in sorting
// behaviour on different platforms.
type (
	DirectoryContents interface {
		All() []fs.DirEntry
		Folders() []fs.DirEntry
		Files() []fs.DirEntry
	}

	// Inspection
	Inspection interface {
		Current() *Node
		Contents() DirectoryContents
		Entries() []fs.DirEntry
		Sort(et enums.EntryType) []fs.DirEntry
		Pick(et enums.EntryType)
		AssignChildren(children []fs.DirEntry)
	}
)
