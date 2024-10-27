package filtering

import (
	"io/fs"
	"path/filepath"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/third/lo"
)

func createGlobFilter(def *core.FilterDef,
	ifNotApplicable bool,
) core.TraverseFilter {
	return &Glob{
		Base: Base{
			name:            def.Description,
			scope:           def.Scope,
			pattern:         def.Pattern,
			negate:          def.Negate,
			ifNotApplicable: ifNotApplicable,
		},
	}
}

// Glob wildcard filter.
type Glob struct {
	Base
}

// IsMatch does this node match the filter
func (f *Glob) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		matched, _ := filepath.Match(f.pattern, node.Extension.Name)
		return f.invert(matched)
	}

	return f.ifNotApplicable
}

// ChildGlobFilter ============================================================

type ChildGlob struct {
	Child
}

// Matching returns the collection of files contained within this
// node's directory that matches this filter.
func (f *ChildGlob) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children,
		func(entry fs.DirEntry, _ int) bool {
			matched, _ := filepath.Match(f.Pattern, entry.Name())
			return f.invert(matched)
		},
	)
}

// SampleGlobFilter ===========================================================

// SampleGlob is a hybrid between a child filter and a node filter. It
// is used to filter on a compound basis but has some differences to ChildGlobFilter
// that necessitates its use. The biggest difference is that ChildGlobFilter is
// designed to only be applied to file directory entries, where as SampleGlob
// can be applied to files or directories. It also possesses a scope field used to
// distinguish only between files and directories.
type SampleGlob struct {
	Sample
}

func (f *SampleGlob) Matching(entries []fs.DirEntry) []fs.DirEntry {
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
