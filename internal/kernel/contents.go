package kernel

import (
	"io/fs"

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/tapable"
)

func NewContents(behaviour *pref.SortBehaviour,
	hook tapable.Hook[core.SortHook, core.ChainSortHook],
	entries []fs.DirEntry,
) *Contents {
	contents := Contents{
		behaviour: behaviour,
		hook:      hook,
	}

	contents.files, contents.directories = nef.Separate(entries)

	return &contents
}

// Contents represents the contents of a directory and handles sorting
// order which by default is different between various operating systems.
// This abstraction removes the differences in sorting behaviour on
// different platforms.
type Contents struct {
	directories []fs.DirEntry
	files       []fs.DirEntry
	behaviour   *pref.SortBehaviour
	hook        tapable.Hook[core.SortHook, core.ChainSortHook]
}

func (c *Contents) Directories() []fs.DirEntry {
	return c.directories
}

func (c *Contents) Files() []fs.DirEntry {
	return c.files
}

// All returns the contents of a directory respecting the directory sorting
// order defined in the traversal options.
func (c *Contents) All() []fs.DirEntry {
	//nolint:ineffassign,staticcheck // prealloc
	result := make([]fs.DirEntry, 0, len(c.files)+len(c.directories))

	if c.behaviour.SortFilesFirst {
		result = c.files
		result = append(result, c.directories...)
	} else {
		result = c.directories
		result = append(result, c.files...)
	}

	return result
}

// Sort will sort either the directories or files or both.
func (c *Contents) Sort(et enums.EntryType) {
	type selectors map[enums.EntryType]func() [][]fs.DirEntry

	sortables := selectors{
		enums.EntryTypeAll: func() [][]fs.DirEntry {
			return [][]fs.DirEntry{
				c.directories, c.files,
			}
		},
		enums.EntryTypeDirectory: func() [][]fs.DirEntry {
			return [][]fs.DirEntry{c.directories}
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
	c.directories = []fs.DirEntry{}
}

func newEmptyContents(prealloc ...*pref.EntryQuantities) *Contents {
	return lo.TernaryF(len(prealloc) == 0,
		func() *Contents {
			return &Contents{
				files:       []fs.DirEntry{},
				directories: []fs.DirEntry{},
			}
		},
		func() *Contents {
			return &Contents{
				files:       make([]fs.DirEntry, 0, prealloc[0].Files),
				directories: make([]fs.DirEntry, 0, prealloc[0].Directories),
			}
		},
	)
}
