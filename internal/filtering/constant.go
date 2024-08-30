package filtering

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
)

func NewPermissiveTraverseFilter(def *core.FilterDef) core.TraverseFilter {
	return &permissiveTraverseFilter{
		Filter: Filter{
			name:  def.Description,
			scope: enums.ScopeTop,
		},
	}
}

type permissiveTraverseFilter struct {
	Filter
	match bool
}

func (f *permissiveTraverseFilter) IsApplicable(_ *core.Node) bool {
	return true
}

func (f *permissiveTraverseFilter) IsMatch(_ *core.Node) bool {
	return true
}

func NewProhibitiveTraverseFilter(def *core.FilterDef) core.TraverseFilter {
	return &prohibitiveTraverseFilter{
		Filter: Filter{
			name:  def.Description,
			scope: enums.ScopeTop,
		},
	}
}

type prohibitiveTraverseFilter struct {
	Filter
	match bool
}

func (f *prohibitiveTraverseFilter) IsApplicable(_ *core.Node) bool {
	return true
}

func (f *prohibitiveTraverseFilter) IsMatch(_ *core.Node) bool {
	return false
}
