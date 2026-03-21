package level

import (
	"path/filepath"
	"strings"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
)

// Periscope depth and scope manager
type Periscope struct {
	depth int
}

// New creates a new periscope instance for a fresh session
func New() *Periscope {
	return &Periscope{}
}

// Offset should be invoked with the Depth loaded into the ActiveState
// during bridging from the previous resume session and the current
// resume session.
func (p *Periscope) Offset(by int) {
	p.depth += by
}

// Scope returns the scope of the current node based on whether it is a leaf
// and the current depth. The scope is determined using bitwise operations on
// the enums.FilterScope values, allowing for combinations of scopes to be
// represented efficiently.
func (p *Periscope) Scope(isLeaf bool) enums.FilterScope {
	result := enums.ScopeIntermediate

	// Tree=0
	// Top=1
	//
	depth := p.Depth()

	switch {
	case isLeaf && depth == 0:
		result = enums.ScopeTree | enums.ScopeLeaf
	case depth == 0:
		result = enums.ScopeTree
	case isLeaf && depth == 1:
		result = enums.ScopeTop | enums.ScopeLeaf
	case depth == 1:
		result = enums.ScopeTop
	case isLeaf:
		result = enums.ScopeLeaf
	}

	return result
}

// Depth returns the current depth of the periscope, adjusted by a decrement
// to account for the initial state where depth is zero before any descent
// has occurred. This allows the depth to reflect the actual level of traversal
// in the directory structure, where a depth of zero corresponds to the
// root level.
func (p *Periscope) Depth() int {
	return p.depth - 1
}

// Delta checks the tree path, with the path of the current node
func (p *Periscope) Delta(tree, current string) (err error) {
	rootSize := len(strings.Split(tree, string(filepath.Separator)))
	currentSize := len(strings.Split(current, string(filepath.Separator)))

	if rootSize > currentSize {
		return core.NewInvalidPeriscopeRootPathError(tree, current)
	}

	return nil
}

// Descend increments the depth of the periscope and checks if it exceeds
// the specified maximum depth. If the maximum depth is exceeded, it returns
// false, indicating that the traversal should not continue deeper. Otherwise,
// it returns true, allowing the traversal to proceed to the next level. The use
// of a maximum depth helps to prevent infinite recursion and allows for
// controlled traversal of directory structures.
func (p *Periscope) Descend(maximum uint) bool {
	if maximum > 0 && p.depth > int(maximum) { //nolint:gosec // ok
		return false
	}

	p.depth++

	return true
}

// Ascend decrements the depth of the periscope, allowing the traversal to
// move back up the directory structure. This is typically called after
// processing a directory's contents, signaling that the traversal is moving
// back up to the parent level. By managing the depth in this way, the periscope
// can accurately track the current level of traversal and ensure that operations
// are performed at the correct depth in the directory hierarchy.
func (p *Periscope) Ascend() {
	p.depth--
}
