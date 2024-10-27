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
	return &RegExpr{
		Base: Base{
			name:            def.Description,
			scope:           def.Scope,
			pattern:         def.Pattern,
			negate:          def.Negate,
			ifNotApplicable: ifNotApplicable,
		},
	}
}

// RegexFilter ================================================================

// RegExpr regex filter.
type RegExpr struct {
	Base
	rex *regexp.Regexp
}

// Validate ensures the filter definition is valid, panics when invalid
func (f *RegExpr) Validate() error {
	if err := f.Base.Validate(); err != nil {
		return err
	}

	var (
		err error
	)
	f.rex, err = regexp.Compile(f.pattern)

	return err
}

// IsMatch
func (f *RegExpr) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		return f.invert(f.rex.MatchString(node.Extension.Name))
	}

	return f.ifNotApplicable
}

// ChildRegexFilter ===========================================================

type ChildRegex struct {
	Child
	rex *regexp.Regexp
}

func (f *ChildRegex) Validate() error {
	var (
		err error
	)
	f.rex, err = regexp.Compile(f.Pattern)

	return err
}

func (f *ChildRegex) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, _ int) bool {
		return f.invert(f.rex.MatchString(entry.Name()))
	})
}

// SampleRegexFilter ==========================================================

// SampleRegex is a hybrid between a child filter and a node filter. It
// is used to filter on a compound basis but has some differences to ChildRegexFilter
// that necessitates its use. The biggest difference is that ChildRegexFilter is
// designed to only be applied to file directory entries, where as SampleRegex
// can be applied to files or directories. It also possesses a scope field used to
// distinguish only between files and directories.
type SampleRegex struct {
	Sample
	rex *regexp.Regexp
}

// Validate ensures the filter definition is valid, panics when invalid
func (f *SampleRegex) Validate() error {
	if err := f.Base.Validate(); err != nil {
		return err
	}

	var (
		err error
	)
	f.rex, err = regexp.Compile(f.pattern)

	return err
}

func (f *SampleRegex) Matching(entries []fs.DirEntry) []fs.DirEntry {
	filterable, bypass := f.fetch(entries)
	filtered := lo.Filter(filterable, func(entry fs.DirEntry, _ int) bool {
		return f.invert(f.rex.MatchString(entry.Name()))
	})

	filtered = append(filtered, bypass...)

	return filtered
}
