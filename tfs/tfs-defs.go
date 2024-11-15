package tfs

import (
	nef "github.com/snivilised/nefilim"
)

// 📦 tfs: nef - contains definitions for traversal file system. Should
// not depend on anything else in agenor.

type (
	// TraversalFS non streaming file system with reader and some
	// writer capabilities
	TraversalFS interface {
		nef.MakeDirFS
		nef.ReaderFS
		nef.WriteFileFS
	}
)
