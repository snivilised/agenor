package kernel

import (
	"io/fs"

	"github.com/snivilised/traverse/pref"
)

func read(sys fs.FS, o *pref.Options, path string) (*DirectoryContents, error) {
	entries, err := o.Hooks.ReadDirectory.Invoke()(sys, path)

	contents := newDirectoryContents(o, entries)
	return contents, err
}
