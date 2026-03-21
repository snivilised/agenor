package filtering

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/internal/third/lo"
)

// Filter =====================================================================

// Base base filter struct.
type Base struct {
	name            string
	pattern         string
	scope           enums.FilterScope // defines which file system nodes the filter should be applied to
	negate          bool              // select to define a negative match
	ifNotApplicable bool
}

// Description description of the filter
func (f *Base) Description() string {
	return f.name
}

// Source text defining the filter
func (f *Base) Source() string {
	return f.pattern
}

// IsApplicable is this filter applicable to this node's scope
func (f *Base) IsApplicable(node *core.Node) bool {
	return (f.scope & node.Extension.Scope) > 0
}

// Scope what items this filter applies to
func (f *Base) Scope() enums.FilterScope {
	return f.scope
}

// invert inverts the result of the filter
func (f *Base) invert(result bool) bool {
	return lo.Ternary(f.negate, !result, result)
}

// Validate validates the filter
func (f *Base) Validate() error {
	if f.scope == enums.ScopeUndefined {
		f.scope = enums.ScopeAll
	}

	return nil
}
