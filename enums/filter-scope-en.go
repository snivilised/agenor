package enums

import (
	"math"
)

//go:generate stringer -type=FilterScope -linecomment -trimprefix=Scope -output filter-scope-en-auto.go

// FilterScope allows client to define which node types should be filtered.
// Filters can be applied to multiple node types by bitwise or-ing the XXXNodes
// definitions. A node may have multiple scope designations, eg a node may be top
// level and leaf if the top level directory does not itself contain further
// sub-directories thereby making it also a leaf.
// It should be noted a file is only a leaf node all of its siblings are all files
// only
type FilterScope uint32

const (
	ScopeUndefined FilterScope = 0 // undefined-scope

	// ScopeRoot, the Root scope
	//
	ScopeRoot FilterScope = 1 << (iota - 1) // root-scope

	// ScopeTop, any node that is a direct descendent of the root node
	//
	ScopeTop // top-scope

	// ScopeLeaf, for directories, any node that has no sub folders. For files, any node
	// that appears under a leaf directory node
	//
	ScopeLeaf // leaf-scope

	// ScopeIntermediate, apply filter to nodes which are neither leaf or top nodes
	//
	ScopeIntermediate // intermediate-scope

	// ScopeFile attributed to file nodes
	//
	ScopeFile // file-scope

	// ScopeFolder attributed to directory nodes
	//
	ScopeFolder // folder-scope

	// ScopeCustom, client defined categorisation (yet to be confirmed)
	//
	ScopeCustom // custom-scope

	// ScopeAll represents any node type
	//
	ScopeAll = math.MaxUint32 // all-scopes
)

// Set sets the bit position indicated by mask
func (f *FilterScope) Set(mask FilterScope) {
	*f |= mask
}

// Clear clears the bit position indicated by mask
func (f *FilterScope) Clear(mask FilterScope) {
	*f &^= mask
}
