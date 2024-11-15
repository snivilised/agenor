package filter_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/third/lo"
)

func TestFilter(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filter Suite")
}

type FilterTE struct {
	lab.NaviTE
	Description     string
	Pattern         string
	Scope           enums.FilterScope
	Negate          bool
	ErrorContains   string
	IfNotApplicable enums.TriStateBool
	Custom          core.TraverseFilter
	Type            enums.FilterType
	Sample          core.SampleTraverseFilter
}

type PolyTE struct {
	lab.NaviTE
	File      core.FilterDef
	Directory core.FilterDef
}

type customFilter struct {
	name            string
	pattern         string
	scope           enums.FilterScope
	negate          bool
	ifNotApplicable bool
}

// Description describes filter
func (f *customFilter) Description() string {
	return f.name
}

func (f *customFilter) Validate() error {
	if f.scope == enums.ScopeUndefined {
		f.scope = enums.ScopeAll
	}

	return nil
}

func (f *customFilter) Source() string {
	return f.pattern
}

func (f *customFilter) invert(result bool) bool {
	return lo.Ternary(f.negate, !result, result)
}

func (f *customFilter) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		matched, _ := filepath.Match(f.pattern, node.Extension.Name)
		return f.invert(matched)
	}

	return f.ifNotApplicable
}

func (f *customFilter) IsApplicable(node *core.Node) bool {
	return (f.scope & node.Extension.Scope) > 0
}

func (f *customFilter) Scope() enums.FilterScope {
	return f.scope
}
