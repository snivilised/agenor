package filtering

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/locale"
)

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
