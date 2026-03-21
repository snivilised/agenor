package filtering

import (
	"io/fs"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/internal/third/lo"
	"github.com/snivilised/jaywalk/locale"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	nef "github.com/snivilised/nefilim"
)

// BaseSampleFilter ===========================================================

type (
	candidates func(entries []fs.DirEntry) (wanted, others []fs.DirEntry)
)

// Sample is a filter that samples files or directories based on a pattern.
type Sample struct {
	Base
}

// NewSample creates a new sample filter based on the given definition and options.
func NewSample(def *core.SampleFilterDef,
	so *pref.SamplingOptions,
) (core.SampleTraverseFilter, error) {
	var (
		filter core.SampleTraverseFilter
	)

	if def == nil {
		return nil, locale.ErrFilterIsNil
	}

	base := Sample{
		Base: Base{
			name:    def.Description,
			scope:   def.Scope.Scrub(),
			pattern: def.Pattern,
			negate:  def.Negate,
		},
	}

	if base.scope.IsFile() && so.NoOf.Files == 0 {
		return nil, locale.ErrInvalidFileSamplingSpecMissingFiles
	}

	if base.scope.IsDirectory() && so.NoOf.Directories == 0 {
		return nil, locale.ErrInvalidSamplingSpecMissingDirectories
	}

	switch def.Type {
	case enums.FilterTypeGlobEx:
	case enums.FilterTypeRegex:
		filter = &SampleRegex{
			Sample: base,
		}
	case enums.FilterTypeGlob:
		filter = &SampleGlob{
			Sample: base,
		}

	case enums.FilterTypeCustom:
		if def.Custom == nil {
			return nil, locale.ErrFilterIsNil
		}

		filter = def.Custom
	case enums.FilterTypePoly:
	case enums.FilterTypeUndefined:
		return nil, locale.ErrFilterMissingType
	}

	if filter != nil {
		if err := filter.Validate(); err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func (f *Sample) files(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	wanted, others = nef.Separate(entries)
	return wanted, others
}

func (f *Sample) directories(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	others, wanted = nef.Separate(entries)
	return wanted, others
}

func (f *Sample) all(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	return entries, []fs.DirEntry{}
}

func (f *Sample) fn() candidates {
	if f.scope.IsDirectory() {
		return f.directories
	}

	if f.scope.IsFile() {
		return f.files
	}

	return f.all
}

func (f *Sample) fetch(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	return f.fn()(entries)
}

// GetMatching sampler func.
type GetMatching func(entry fs.DirEntry, index int) bool

// Matching returns the collection of files contained within this
// node's directory that matches this filter.
func (f *Sample) Matching(children []fs.DirEntry,
	get GetMatching,
) []fs.DirEntry {
	filterable, bypass := f.fetch(children)
	filtered := lo.Filter(filterable, get)

	return append(filtered, bypass...)
}

// NewCustomSampleFilter only needs to be called explicitly when defining
// a custom sample filter.
func NewCustomSampleFilter(scope enums.FilterScope) Sample {
	return Sample{
		Base: Base{
			scope: scope,
		},
	}
}
