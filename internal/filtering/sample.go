package filtering

import (
	"io/fs"

	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/nfs"
)

// BaseSampleFilter ===========================================================

type (
	candidates func(entries []fs.DirEntry) (wanted, others []fs.DirEntry)
)

type SampleFilter struct {
	Filter
}

// NewSampleFilter only needs to be called explicitly when defining
// a custom sample filter.
func NewSampleFilter(scope enums.FilterScope) SampleFilter {
	return SampleFilter{
		Filter: Filter{
			scope: scope,
		},
	}
}

func (f *SampleFilter) files(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	wanted, others = nfs.Separate(entries)
	return wanted, others
}

func (f *SampleFilter) folders(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	others, wanted = nfs.Separate(entries)
	return wanted, others
}

func (f *SampleFilter) all(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	return entries, []fs.DirEntry{}
}

func (f *SampleFilter) fn() candidates {
	if f.scope.IsFolder() {
		return f.folders
	}

	if f.scope.IsFile() {
		return f.files
	}

	return f.all
}

func (f *SampleFilter) fetch(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	return f.fn()(entries)
}

// GetMatching sampler func.
type GetMatching func(entry fs.DirEntry, index int) bool

func (f *SampleFilter) Matching(children []fs.DirEntry,
	get GetMatching,
) []fs.DirEntry {
	filterable, bypass := f.fetch(children)
	filtered := lo.Filter(filterable, get)

	return append(filtered, bypass...)
}
