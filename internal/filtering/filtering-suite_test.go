package filtering_test

import (
	"path/filepath"
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
)

func TestFiltering(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Filtering Suite")
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
