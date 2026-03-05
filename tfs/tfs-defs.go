package tfs

import (
	nef "github.com/snivilised/nefilim"
)

type (
	// TraversalFS non streaming file system with reader and some
	// writer capabilities
	TraversalFS interface {
		nef.MakeDirFS
		nef.ReaderFS
		nef.WriteFileFS
	}
)
