package kernel

import (
	"io/fs"

	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/pref"
)

func newDirectoryContents(o *pref.Options,
	entries []fs.DirEntry,
) *DirectoryContents {
	contents := DirectoryContents{
		o: o,
	}

	contents.Arrange(entries)

	return &contents
}

// DirectoryContents represents the contents of a directory's contents and
// handles sorting order which by default is different between various
// operating systems. This abstraction removes the differences in sorting
// behaviour on different platforms.
type DirectoryContents struct {
	folders []fs.DirEntry
	files   []fs.DirEntry
	o       *pref.Options
}

func (c *DirectoryContents) Folders() []fs.DirEntry {
	return c.folders
}

func (c *DirectoryContents) Files() []fs.DirEntry {
	return c.files
}

// All returns the contents of a directory respecting the directory sorting
// order defined in the traversal options.
func (c *DirectoryContents) All() []fs.DirEntry {
	result := make([]fs.DirEntry, 0, len(c.files)+len(c.folders))

	switch c.o.Core.Behaviours.Sort.DirectoryEntryOrder {
	case enums.DirectoryContentsOrderFoldersFirst:
		result = c.folders
		result = append(result, c.files...)

	case enums.DirectoryContentsOrderFilesFirst:
		result = c.files
		result = append(result, c.folders...)
	}

	return result
}

// Sort will sort either the folders or files or if no
// entry type is specified, sort both.
func (c *DirectoryContents) Sort(ents ...enums.EntryType) {
	// This looks complicated, but it really isn't. The reason is
	// we only want to sort only what's really required and the sorting
	// of entries must be separated by type, ie we dont want all the
	// files and folders sorted into a single collection and the
	// requested entry types need to be mapped into the corresponding
	// internal directory entries.
	//
	for _, entries := range lo.TernaryF(len(ents) == 0,
		func() [][]fs.DirEntry {
			return [][]fs.DirEntry{
				c.folders, c.files,
			}
		},
		func() [][]fs.DirEntry {
			if ents[0] == enums.EntryTypeFolder {
				return [][]fs.DirEntry{c.folders}
			}

			return [][]fs.DirEntry{c.files}
		},
	) {
		c.o.Hooks.Sort.Invoke()(entries)
	}
}

func (c *DirectoryContents) Clear() {
	c.files = []fs.DirEntry{}
	c.folders = []fs.DirEntry{}
}

func (c *DirectoryContents) Arrange(entries []fs.DirEntry) {
	grouped := lo.GroupBy(entries, func(entry fs.DirEntry) bool {
		return entry.IsDir()
	})

	c.folders = grouped[true]
	c.files = grouped[false]

	if c.folders == nil {
		c.folders = []fs.DirEntry{}
	}

	if c.files == nil {
		c.files = []fs.DirEntry{}
	}
}

func newEmptyDirectoryEntries(o *pref.Options,
	prealloc ...*pref.EntryQuantities,
) *DirectoryContents {
	return lo.TernaryF(len(prealloc) == 0,
		func() *DirectoryContents {
			return &DirectoryContents{
				o:       o,
				files:   []fs.DirEntry{},
				folders: []fs.DirEntry{},
			}
		},
		func() *DirectoryContents {
			return &DirectoryContents{
				o:       o,
				files:   make([]fs.DirEntry, 0, prealloc[0].Files),
				folders: make([]fs.DirEntry, 0, prealloc[0].Folders),
			}
		},
	)
}
