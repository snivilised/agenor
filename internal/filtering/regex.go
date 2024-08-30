package filtering

import (
	"io/fs"
	"regexp"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/third/lo"
)

func createRegexFilter(def *core.FilterDef,
	ifNotApplicable bool,
) core.TraverseFilter {
	return &RegexFilter{
		Filter: Filter{
			name:            def.Description,
			scope:           def.Scope,
			pattern:         def.Pattern,
			negate:          def.Negate,
			ifNotApplicable: ifNotApplicable,
		},
	}
}

// RegexFilter ================================================================

// RegexFilter regex filter.
type RegexFilter struct {
	Filter
	rex *regexp.Regexp
}

// Validate ensures the filter definition is valid, panics when invalid
func (f *RegexFilter) Validate() error {
	if err := f.Filter.Validate(); err != nil {
		return err
	}

	var (
		err error
	)
	f.rex, err = regexp.Compile(f.pattern)

	return err
}

// IsMatch
func (f *RegexFilter) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		return f.invert(f.rex.MatchString(node.Extension.Name))
	}

	return f.ifNotApplicable
}

// ChildRegexFilter ===========================================================

type ChildRegexFilter struct {
	ChildFilter
	rex *regexp.Regexp
}

func (f *ChildRegexFilter) Validate() error {
	var (
		err error
	)
	f.rex, err = regexp.Compile(f.Pattern)

	return err
}

func (f *ChildRegexFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, _ int) bool {
		return f.invert(f.rex.MatchString(entry.Name()))
	})
}

// SampleRegexFilter ==========================================================

// SampleRegexFilter is a hybrid between a child filter and a node filter. It
// is used to filter on a compound basis but has some differences to ChildRegexFilter
// that necessitates its use. The biggest difference is that ChildRegexFilter is
// designed to only be applied to file directory entries, where as SampleRegexFilter
// can be applied to files or folders. It also possesses a scope field used to
// distinguish only between files and folders.
type SampleRegexFilter struct {
	SampleFilter
	rex *regexp.Regexp
}

// Validate ensures the filter definition is valid, panics when invalid
func (f *SampleRegexFilter) Validate() error {
	if err := f.Filter.Validate(); err != nil {
		return err
	}

	var (
		err error
	)
	f.rex, err = regexp.Compile(f.pattern)

	return err
}

func (f *SampleRegexFilter) Matching(entries []fs.DirEntry) []fs.DirEntry {
	filterable, bypass := f.fetch(entries)
	filtered := lo.Filter(filterable, func(entry fs.DirEntry, _ int) bool {
		return f.invert(f.rex.MatchString(entry.Name()))
	})

	filtered = append(filtered, bypass...)

	return filtered
}
