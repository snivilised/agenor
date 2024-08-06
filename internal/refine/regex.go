package refine

import (
	"io/fs"
	"regexp"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/lo"
)

// RegexFilter ================================================================

// RegexFilter regex filter.
type RegexFilter struct {
	Filter
	rex *regexp.Regexp
}

// Validate ensures the filter definition is valid, panics when invalid
func (f *RegexFilter) Validate() {
	f.Filter.Validate()
	f.rex = regexp.MustCompile(f.pattern)
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

func (f *ChildRegexFilter) Validate() {
	f.rex = regexp.MustCompile(f.Pattern)
}

func (f *ChildRegexFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, _ int) bool {
		return f.invert(f.rex.MatchString(entry.Name()))
	})
}

// SampleGlobFilter ===========================================================

// SampleGlobFilter is a hybrid between a child filter and a node filter. It
// is used to filter on a compound basis but has some differences to ChildGlobFilter
// that necessitates its use. The biggest difference is that ChildGlobFilter is
// designed to only be applied to file directory entries, where as SampleGlobFilter
// can be applied to files or folders. It also possesses a scope field used to
// distinguish only between files and folders.
type SampleRegexFilter struct {
	SampleFilter
	rex *regexp.Regexp
}

// Validate ensures the filter definition is valid, panics when invalid
func (f *SampleRegexFilter) Validate() {
	f.Filter.Validate()
	f.rex = regexp.MustCompile(f.pattern)
}

func (f *SampleRegexFilter) Matching(entries []fs.DirEntry) []fs.DirEntry {
	filterable, bypass := f.fetch(entries)
	filtered := lo.Filter(filterable, func(entry fs.DirEntry, _ int) bool {
		return f.invert(f.rex.MatchString(entry.Name()))
	})

	filtered = append(filtered, bypass...)

	return filtered
}
