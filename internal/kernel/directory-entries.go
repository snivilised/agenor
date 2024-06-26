package kernel

import (
	"io/fs"

	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/pref"
)

func newDirectoryContents(o *pref.Options, entries []fs.DirEntry) *DirectoryContents {
	contents := DirectoryContents{
		o: o,
	}

	contents.arrange(entries)

	return &contents
}

// DirectoryContents represents the contents of a directory's contents and
// handles sorting order which by default is different between various
// operating systems. This abstraction removes the differences in sorting
// behaviour on different platforms.
type DirectoryContents struct {
	Folders []fs.DirEntry
	Files   []fs.DirEntry
	o       *pref.Options
}

// All returns the contents of a directory respecting the directory sorting
// order defined in the traversal options.
func (e *DirectoryContents) All() []fs.DirEntry {
	result := make([]fs.DirEntry, 0, len(e.Files)+len(e.Folders))

	switch e.o.Core.Behaviours.Sort.DirectoryEntryOrder {
	case enums.DirectoryContentsOrderFoldersFirst:
		result = e.Folders
		result = append(result, e.Files...)

	case enums.DirectoryContentsOrderFilesFirst:
		result = e.Files
		result = append(result, e.Folders...)
	}

	return result
}

func (e *DirectoryContents) arrange(entries []fs.DirEntry) {
	grouped := lo.GroupBy(entries, func(entry fs.DirEntry) bool {
		return entry.IsDir()
	})

	e.Folders = grouped[true]
	e.Files = grouped[false]

	if e.Folders == nil {
		e.Folders = []fs.DirEntry{}
	}

	if e.Files == nil {
		e.Files = []fs.DirEntry{}
	}
}
