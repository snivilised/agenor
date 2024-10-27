package filtering

import (
	"slices"
	"strings"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
)

func NewChild(def *core.ChildFilterDef) (core.ChildTraverseFilter, error) {
	var (
		filter core.ChildTraverseFilter
	)

	switch def.Type {
	case enums.FilterTypeGlobEx:
		var (
			err                error
			segments, suffixes []string
		)

		if segments, suffixes, err = splitGlobExPattern(def.Pattern); err != nil {
			return nil, locale.NewInvalidIncaseFilterDefError(def.Pattern)
		}

		base, exclusion := splitGlob(segments[0])

		filter = &ChildGlobExFilter{
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

// ChildFilter ================================================================

// Child filter used when subscription is DirectoriesWithFiles
type Child struct {
	Name    string
	Pattern string
	Negate  bool
}

func (f *Child) Description() string {
	return f.Name
}

func (f *Child) Validate() error {
	return nil
}

func (f *Child) Source() string {
	return f.Pattern
}

func (f *Child) invert(result bool) bool {
	return lo.Ternary(f.Negate, !result, result)
}
