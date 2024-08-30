package filtering

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
)

// Filter =====================================================================

// Filter base filter struct.
type Filter struct {
	name            string
	pattern         string
	scope           enums.FilterScope // defines which file system nodes the filter should be applied to
	negate          bool              // select to define a negative match
	ifNotApplicable bool
}

// Description description of the filter
func (f *Filter) Description() string {
	return f.name
}

// Source text defining the filter
func (f *Filter) Source() string {
	return f.pattern
}

func (f *Filter) IsApplicable(node *core.Node) bool {
	return (f.scope & node.Extension.Scope) > 0
}

func (f *Filter) Scope() enums.FilterScope {
	return f.scope
}

func (f *Filter) invert(result bool) bool {
	return lo.Ternary(f.negate, !result, result)
}

func (f *Filter) Validate() error {
	if f.scope == enums.ScopeUndefined {
		f.scope = enums.ScopeAll
	}

	return nil
}

// ChildFilter ================================================================

// ChildFilter filter used when subscription is FoldersWithFiles
type ChildFilter struct {
	Name    string
	Pattern string
	Negate  bool
}

func (f *ChildFilter) Description() string {
	return f.Name
}

func (f *ChildFilter) Validate() error {
	return nil
}

func (f *ChildFilter) Source() string {
	return f.Pattern
}

func (f *ChildFilter) invert(result bool) bool {
	return lo.Ternary(f.Negate, !result, result)
}
