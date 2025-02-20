package level

import (
	"path/filepath"
	"strings"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
)

// 📦 pkg: level - contains functionality concerned only with depth
// management.

// Periscope: depth and scope manager
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

func (p *Periscope) Depth() int {
	return p.depth - 1
}

func (p *Periscope) Delta(tree, current string) (err error) {
	rootSize := len(strings.Split(tree, string(filepath.Separator)))
	currentSize := len(strings.Split(current, string(filepath.Separator)))

	if rootSize > currentSize {
		return core.NewInvalidPeriscopeRootPathError(tree, current)
	}

	return nil
}

func (p *Periscope) Descend(maximum uint) bool {
	if maximum > 0 && p.depth > int(maximum) { //nolint:gosec // ok
		return false
	}

	p.depth++

	return true
}

func (p *Periscope) Ascend() {
	p.depth--
}
