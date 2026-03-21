package filtering

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
)

// NewPermissiveTraverseFilter creates a new permissive traverse filter.
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
}

// IsApplicable is this filter applicable to this node's scope
func (f *permissiveTraverseFilter) IsApplicable(_ *core.Node) bool {
	return true
}

// IsMatch determines if this node matches the filter.
func (f *permissiveTraverseFilter) IsMatch(_ *core.Node) bool {
	return true
}

// NewProhibitiveTraverseFilter creates a new prohibitive traverse filter.
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
}

// IsApplicable is this filter applicable to this node's scope
func (f *prohibitiveTraverseFilter) IsApplicable(_ *core.Node) bool {
	return true
}

// IsMatch determines if this node matches the filter.
func (f *prohibitiveTraverseFilter) IsMatch(_ *core.Node) bool {
	return false
}
