package level

import (
	"path/filepath"
	"strings"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
)

// 📚 package: level contains functionality concerned only with depth
// management.

// Periscope: depth and scope manager
type Periscope struct {
	offset int
	depth  int
}

// New creates a new periscope instance for a fresh session
func New() *Periscope {
	return &Periscope{}
}

// Restore creates a new periscope instance required for resume sessions
func Restore(offset, depth int) *Periscope {
	return &Periscope{
		offset: offset,
		depth:  depth,
	}
}

func (p *Periscope) Scope(isLeaf bool) enums.FilterScope {
	result := enums.ScopeIntermediate

	// Root=0
	// Top=1
	//
	depth := p.Depth()

	switch {
	case isLeaf && depth == 0:
		result = enums.ScopeRoot | enums.ScopeLeaf
	case depth == 0:
		result = enums.ScopeRoot
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
	return p.offset + p.depth - 1
}

func (p *Periscope) Delta(root, current string) (err error) {
	rootSize := len(strings.Split(root, string(filepath.Separator)))
	currentSize := len(strings.Split(current, string(filepath.Separator)))

	if rootSize > currentSize {
		return core.NewInvalidPeriscopeRootPathNativeError(root, current)
	}

	p.offset = currentSize - rootSize

	return nil
}

func (p *Periscope) Descend(maximum uint) bool {
	if maximum > 0 && p.depth > int(maximum) {
		return false
	}

	p.depth++

	return true
}

func (p *Periscope) Ascend() {
	p.depth--
}
