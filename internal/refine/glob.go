package refine

import (
	"io/fs"
	"path/filepath"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/lo"
)

// GlobFilter wildcard filter.
type GlobFilter struct {
	Filter
}

// IsMatch does this node match the filter
func (f *GlobFilter) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		matched, _ := filepath.Match(f.pattern, node.Extension.Name)
		return f.invert(matched)
	}

	return f.ifNotApplicable
}

// ChildGlobFilter ============================================================

type ChildGlobFilter struct {
	ChildFilter
}

// Matching returns the collection of files contained within this
// node's folder that matches this filter.
func (f *ChildGlobFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, _ int) bool {
		matched, _ := filepath.Match(f.Pattern, entry.Name())
		return f.invert(matched)
	})
}
