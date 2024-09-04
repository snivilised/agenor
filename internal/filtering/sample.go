package filtering

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/nfs"
	"github.com/snivilised/traverse/pref"
)

// BaseSampleFilter ===========================================================

type (
	candidates func(entries []fs.DirEntry) (wanted, others []fs.DirEntry)
)

type Sample struct {
	Base
}

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

	if base.scope.IsFolder() && so.NoOf.Folders == 0 {
		return nil, locale.ErrInvalidFolderSamplingSpecMissingFolders
	}

	switch def.Type {
	case enums.FilterTypeExtendedGlob:
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
	wanted, others = nfs.Separate(entries)
	return wanted, others
}

func (f *Sample) folders(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	others, wanted = nfs.Separate(entries)
	return wanted, others
}

func (f *Sample) all(entries []fs.DirEntry) (wanted, others []fs.DirEntry) {
	return entries, []fs.DirEntry{}
}

func (f *Sample) fn() candidates {
	if f.scope.IsFolder() {
		return f.folders
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
