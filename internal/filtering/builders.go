package filtering

import (
	"strings"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
)

func buildNativeNodeFilter(definition *core.FilterDef) (core.TraverseFilter, error) {
	var (
		filter          core.TraverseFilter
		ifNotApplicable = applicable(definition.IfNotApplicable)
		err             error
	)

	switch definition.Type {
	case enums.FilterTypeGlobEx:
		filter, err = createGlobExFilter(definition, ifNotApplicable)

	case enums.FilterTypeRegex:
		filter = createRegexFilter(definition, ifNotApplicable)

	case enums.FilterTypeGlob:
		filter = createGlobFilter(definition, ifNotApplicable)

	case enums.FilterTypeCustom, enums.FilterTypePoly:
		return nil, nil

	case enums.FilterTypeUndefined:
		return nil, locale.ErrFilterMissingType
	}

	if err != nil {
		return nil, err
	}

	return filter, filter.Validate()
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

func buildPolyNodeFilter(definition *core.FilterDef,
	fo *pref.FilterOptions,
	nativeFn filterNativeFunc,
	customFn filterUsingOptionsFunc,
) (core.TraverseFilter, error) {
	if definition.Type != enums.FilterTypePoly {
		return nil, nil
	}
	polyDef := fo.Node.Poly

	// enforce the correct filter scopes
	//
	polyDef.File.Scope.Set(enums.ScopeFile)
	polyDef.File.Scope.Clear(enums.ScopeDirectory)

	polyDef.Directory.Scope.Set(enums.ScopeDirectory)
	polyDef.Directory.Scope.Clear(enums.ScopeFile)

	if polyDef.File.Type == enums.FilterTypePoly {
		return nil, locale.ErrPolyFilterIsInvalid
	}

	file, err := buildConstituent(&polyDef.File, fo, nativeFn, customFn)

	if err != nil {
		return nil, err
	}

	if polyDef.Directory.Type == enums.FilterTypePoly {
		return nil, locale.ErrPolyFilterIsInvalid
	}

	directory, err := buildConstituent(&polyDef.Directory, fo, nativeFn, customFn)

	if err != nil {
		return nil, err
	}

	filter := &Poly{
		File:      file,
		Directory: directory,
	}

	return filter, nil
}

func buildConstituent(definition *core.FilterDef,
	fo *pref.FilterOptions,
	nativeFn filterNativeFunc,
	customFn filterUsingOptionsFunc,
) (core.TraverseFilter, error) {
	filter, err := OrFuncE(
		func() (core.TraverseFilter, error) {
			return customFn(definition, fo)
		},
		func() (core.TraverseFilter, error) {
			return nativeFn(definition)
		},
	)

	if err != nil {
		return nil, err
	}

	if err := filter.Validate(); err != nil {
		return nil, err
	}

	return filter, err
}

func getCustomFilter(_ *core.FilterDef,
	fo *pref.FilterOptions,
) (core.TraverseFilter, error) {
	return fo.Custom, nil
}

func splitGlobExPattern(pattern string) (segments, suffixes []string, err error) {
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
