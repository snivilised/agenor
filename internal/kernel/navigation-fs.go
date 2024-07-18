package kernel

import (
	"io/fs"
)

func read(sys fs.ReadDirFS, o *readOptions, path string) (*Contents, error) {
	entries, err := o.hooks.read.Invoke()(sys, path)

	contents := newContents(
		o.behaviour, o.hooks.sort, entries,
	)
	return contents, err
}
