package filtering

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/internal/third/lo"
	"github.com/snivilised/jaywalk/locale"
)

// NewChild creates a new child filter.
func NewChild(def *core.ChildFilterDef) (core.ChildTraverseFilter, error) {
	var (
		filter core.ChildTraverseFilter
	)

	switch def.Type {
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

	case enums.FilterTypeGlobEx:
		return nil, locale.ErrFilterChildGlobExNotSupported

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
	// Name is the child filter name
	Name string
	// Pattern is the child filter pattern
	Pattern string
	// Negate is true if the filter should be negated
	Negate bool
}

// Description returns the description of the filter
func (f *Child) Description() string {
	return f.Name
}

// Validate validates the filter
func (f *Child) Validate() error {
	return nil
}

// Source returns the pattern of the filter
func (f *Child) Source() string {
	return f.Pattern
}

func (f *Child) invert(result bool) bool {
	return lo.Ternary(f.Negate, !result, result)
}
