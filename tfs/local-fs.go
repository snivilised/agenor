package tfs

import (
	nef "github.com/snivilised/nefilim"
)

// NewFS creates a relative local file system required for traversal.
func NewFS(rel nef.Rel) TraversalFS {
	return nef.NewUniversalFS(rel)
}

// New creates an absolute local file system required for traversal.
func New() TraversalFS {
	return nef.NewTraverseABS()
}
