package filtering

import (
	"slices"
	"strings"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

func NewNodeFilter(def *core.FilterDef,
	fo *pref.FilterOptions,
) (core.TraverseFilter, error) {
	var (
		filter          core.TraverseFilter
		ifNotApplicable = applicable(def.IfNotApplicable)
		err             error
	)

	switch def.Type {
	case enums.FilterTypeExtendedGlob:
		filter, err = createExtendedGlobFilter(def, ifNotApplicable)

	case enums.FilterTypeRegex:
		filter = createRegexFilter(def, ifNotApplicable)

	// TODO: createGlobFilter
	case enums.FilterTypeGlob:
		filter = createGlobFilter(def, ifNotApplicable)

	// TODO: issue #156: the way we access FilterOptions needs to be
	// tidied up/improved. Just feels way too fishy for custom and poly...
	// perhaps we have extra NewCustomFilter/NewPolyFilter funcs
	//
	// TODO: createCustomFilter
	case enums.FilterTypeCustom:
		if fo != nil {
			if fo.Custom == nil {
				return nil, locale.ErrMissingCustomFilterDefinition
			}
			filter = fo.Custom
		}

	// TODO: createPolyFilter
	case enums.FilterTypePoly:
		if fo != nil {
			var polyE error

			if filter, polyE = createPolyFilter(fo.Node.Poly); polyE != nil {
				return nil, polyE
			}
		}

	case enums.FilterTypeUndefined:
		return nil, locale.ErrFilterMissingType
	}

	if filter != nil {
		err = filter.Validate()
	}

	return filter, err
}

func applicable(ifNotApplicable enums.TriStateBool) bool {
	switch ifNotApplicable {
	case enums.TriStateBoolTrue:
		return true

	case enums.TriStateBoolFalse:
		return false

	case enums.TriStateBoolUndefined:
	}

	return true
}

const (
	exclusionDelim = "/"
)

func splitGlob(baseGlob string) (base, exclusion string) {
	base = strings.ToLower(baseGlob)

	if strings.Contains(base, exclusionDelim) {
		constituents := strings.Split(base, exclusionDelim)
		base = constituents[0]
		exclusion = constituents[1]
	}

	return base, exclusion
}

func NewChildFilter(def *core.ChildFilterDef) (core.ChildTraverseFilter, error) {
	var (
		filter core.ChildTraverseFilter
	)

	if def == nil {
		return nil, locale.ErrFilterIsNil
	}

	switch def.Type {
	case enums.FilterTypeExtendedGlob:
		var (
			err                error
			segments, suffixes []string
		)

		if segments, suffixes, err = splitExtendedGlobPattern(def.Pattern); err != nil {
			return nil, locale.NewInvalidIncaseFilterDefError(def.Pattern)
		}

		base, exclusion := splitGlob(segments[0])

		filter = &ChildExtendedGlobFilter{
			Child: Child{
				Name:    def.Description,
				Pattern: def.Pattern,
				Negate:  def.Negate,
			},
			baseGlob: base,
			suffixes: lo.Map(suffixes, func(s string, _ int) string {
				return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(s), "."))
			}),
			anyExtension: slices.Contains(suffixes, "*"),
			exclusion:    exclusion,
		}

	case enums.FilterTypeRegex:
		filter = &ChildRegex{
			Child: Child{
				Name:    def.Description,
				Pattern: def.Pattern,
				Negate:  def.Negate,
			},
		}

	case enums.FilterTypeGlob:
		filter = &ChildGlob{
			Child: Child{
				Name:    def.Description,
				Pattern: def.Pattern,
				Negate:  def.Negate,
			},
		}

	case enums.FilterTypeCustom:
		return nil, locale.ErrFilterCustomNotSupported

	case enums.FilterTypeUndefined:
		return nil, locale.ErrFilterUndefined

	case enums.FilterTypePoly:
	}

	if filter != nil {
		if err := filter.Validate(); err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func splitExtendedGlobPattern(pattern string) (segments, suffixes []string, err error) {
	if !strings.Contains(pattern, "|") {
		return []string{}, []string{},
			locale.NewInvalidExtGlobFilterMissingSeparatorError(pattern)
	}

	segments = strings.Split(pattern, "|")
	suffixes = strings.Split(segments[1], ",")

	suffixes = lo.Reject(suffixes, func(item string, _ int) bool {
		return item == ""
	})

	return segments, suffixes, nil
}

func newSampleFilter(def *core.SampleFilterDef,
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
