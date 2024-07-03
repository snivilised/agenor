package kernel

import (
	"io/fs"

	"github.com/snivilised/traverse/pref"
)

func read(sys fs.FS, o *pref.Options, path string) (*Contents, error) {
	entries, err := o.Hooks.ReadDirectory.Invoke()(sys, path)

	contents := newContents(o, entries)
	return contents, err
}
