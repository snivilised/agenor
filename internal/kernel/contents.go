package kernel

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/tapable"
)

func newContents(behaviour *pref.SortBehaviour,
	hook tapable.Hook[core.SortHook],
	entries []fs.DirEntry,
) *Contents {
	contents := Contents{
		behaviour: behaviour,
		hook:      hook,
	}

	contents.arrange(entries)

	return &contents
}

// Contents represents the contents of a directory and handles sorting
// order which by default is different between various operating systems.
// This abstraction removes the differences in sorting behaviour on
// different platforms.
type Contents struct {
	folders   []fs.DirEntry
	files     []fs.DirEntry
	behaviour *pref.SortBehaviour
	hook      tapable.Hook[core.SortHook]
}

func (c *Contents) Folders() []fs.DirEntry {
	return c.folders
}

func (c *Contents) Files() []fs.DirEntry {
	return c.files
}

// All returns the contents of a directory respecting the directory sorting
// order defined in the traversal options.
func (c *Contents) All() []fs.DirEntry {
	result := make([]fs.DirEntry, 0, len(c.files)+len(c.folders))

	switch c.behaviour.DirectoryEntryOrder {
	case enums.DirectoryContentsOrderFoldersFirst:
		result = c.folders
		result = append(result, c.files...)

	case enums.DirectoryContentsOrderFilesFirst:
		result = c.files
		result = append(result, c.folders...)
	}

	return result
}

// Sort will sort either the folders or files or both.
func (c *Contents) Sort(et enums.EntryType) {
	type selectors map[enums.EntryType]func() [][]fs.DirEntry

	sortables := selectors{
		enums.EntryTypeAll: func() [][]fs.DirEntry {
			return [][]fs.DirEntry{
				c.folders, c.files,
			}
		},
		enums.EntryTypeFolder: func() [][]fs.DirEntry {
			return [][]fs.DirEntry{c.folders}
		},
		enums.EntryTypeFile: func() [][]fs.DirEntry {
			return [][]fs.DirEntry{c.files}
		},
	}

	for _, entries := range sortables[et]() {
		c.hook.Invoke()(entries)
	}
}

func (c *Contents) clear() {
	c.files = []fs.DirEntry{}
	c.folders = []fs.DirEntry{}
}

func (c *Contents) arrange(entries []fs.DirEntry) {
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

func newEmptyContents(prealloc ...*pref.EntryQuantities) *Contents {
	return lo.TernaryF(len(prealloc) == 0,
		func() *Contents {
			return &Contents{
				files:   []fs.DirEntry{},
				folders: []fs.DirEntry{},
			}
		},
		func() *Contents {
			return &Contents{
				files:   make([]fs.DirEntry, 0, prealloc[0].Files),
				folders: make([]fs.DirEntry, 0, prealloc[0].Folders),
			}
		},
	)
}
