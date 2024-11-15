package filtering

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
)

func NewPermissiveTraverseFilter(def *core.FilterDef) core.TraverseFilter {
	return &permissiveTraverseFilter{
		Base: Base{
			name:  def.Description,
			scope: enums.ScopeTop,
		},
	}
}

type permissiveTraverseFilter struct {
	Base
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
		Base: Base{
			name:  def.Description,
			scope: enums.ScopeTop,
		},
	}
}

type prohibitiveTraverseFilter struct {
	Base
	match bool
}

func (f *prohibitiveTraverseFilter) IsApplicable(_ *core.Node) bool {
	return true
}

func (f *prohibitiveTraverseFilter) IsMatch(_ *core.Node) bool {
	return false
}
