package filtering

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
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

func (f *Base) IsApplicable(node *core.Node) bool {
	return (f.scope & node.Extension.Scope) > 0
}

func (f *Base) Scope() enums.FilterScope {
	return f.scope
}

func (f *Base) invert(result bool) bool {
	return lo.Ternary(f.negate, !result, result)
}

func (f *Base) Validate() error {
	if f.scope == enums.ScopeUndefined {
		f.scope = enums.ScopeAll
	}

	return nil
}
