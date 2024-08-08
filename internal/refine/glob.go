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
	return lo.Filter(children,
		func(entry fs.DirEntry, _ int) bool {
			matched, _ := filepath.Match(f.Pattern, entry.Name())
			return f.invert(matched)
		},
	)
}

// SampleGlobFilter ===========================================================

// SampleGlobFilter is a hybrid between a child filter and a node filter. It
// is used to filter on a compound basis but has some differences to ChildGlobFilter
// that necessitates its use. The biggest difference is that ChildGlobFilter is
// designed to only be applied to file directory entries, where as SampleGlobFilter
// can be applied to files or folders. It also possesses a scope field used to
// distinguish only between files and folders.
type SampleGlobFilter struct {
	SampleFilter
}

func (f *SampleGlobFilter) Matching(entries []fs.DirEntry) []fs.DirEntry {
	filterable, bypass := f.fetch(entries)

	filtered := lo.Filter(filterable,
		func(entry fs.DirEntry, _ int) bool {
			matched, _ := filepath.Match(f.pattern, entry.Name())
			return f.invert(matched)
		},
	)

	filtered = append(filtered, bypass...)

	return filtered
}
