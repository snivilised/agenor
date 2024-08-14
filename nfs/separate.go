package nfs

import (
	"io/fs"

	"github.com/snivilised/traverse/internal/third/lo"
)

func Separate(entries []fs.DirEntry) (files, folders []fs.DirEntry) {
	grouped := lo.GroupBy(entries, func(entry fs.DirEntry) bool {
		return entry.IsDir()
	})

	const (
		asFile   = false
		asFolder = true
	)

	// incase lo.GroupBy has returned a nil for a particular grouping,
	// we make sure we at least have an empty slice instead of allowing
	// nil to be returned to represent an empty result set.
	//
	files = lo.Ternary(grouped[asFile] == nil,
		[]fs.DirEntry{}, grouped[asFile],
	)
	folders = lo.Ternary(grouped[asFolder] == nil,
		[]fs.DirEntry{}, grouped[asFolder],
	)

	return files, folders
}
