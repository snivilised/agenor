package kernel_test

import (
	"fmt"

	. "github.com/onsi/gomega/types" //nolint:revive // ok
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/filtering"
)

// === MatchCurrentRegexFilter ===
//

type IsCurrentRegexMatchMatcher struct {
	rxFilter interface{}
}

func MatchCurrentRegexFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentRegexMatchMatcher{
		rxFilter: expected,
	}
}

func (m *IsCurrentRegexMatchMatcher) Match(actual interface{}) (bool, error) {
	node, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", node)
	}

	rxFilter, filterOk := m.rxFilter.(*filtering.RegExpr)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *RegexFilter (%T)", rxFilter)
	}

	return rxFilter.IsMatch(node), nil
}

func (m *IsCurrentRegexMatchMatcher) FailureMessage(actual interface{}) string {
	node, _ := actual.(*core.Node)
	rxFilter, _ := m.rxFilter.(*filtering.RegExpr)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match regex\n\t%v\n",
		node.Extension.Name, rxFilter.Source(),
	)
}

func (m *IsCurrentRegexMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	node, _ := actual.(*core.Node)
	rxFilter, _ := m.rxFilter.(*filtering.RegExpr)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match regex\n\t%v\n",
		node.Extension.Name, rxFilter.Source(),
	)
}

// === MatchCurrentGlobFilter ===
//

type IsCurrentGlobMatchMatcher struct {
	gbFilter interface{}
}

func MatchCurrentGlobFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentGlobMatchMatcher{
		gbFilter: expected,
	}
}

func (m *IsCurrentGlobMatchMatcher) Match(actual interface{}) (bool, error) {
	node, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", node)
	}

	gbFilter, filterOk := m.gbFilter.(*filtering.Glob)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *GlobFilter (%T)", gbFilter)
	}

	return gbFilter.IsMatch(node), nil
}

func (m *IsCurrentGlobMatchMatcher) FailureMessage(actual interface{}) string {
	node, _ := actual.(*core.Node)
	gbFilter, _ := m.gbFilter.(*filtering.Glob)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match glob\n\t%v\n",
		node.Extension.Name, gbFilter.Source(),
	)
}

func (m *IsCurrentGlobMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	node, _ := actual.(*core.Node)
	gbFilter, _ := m.gbFilter.(*filtering.Glob)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match glob\n\t%v\n",
		node.Extension.Name, gbFilter.Source(),
	)
}

// === MatchCurrentGlobExFilter ===
//

type IsCurrentGlobExMatchMatcher struct {
	egbFilter interface{}
}

func MatchCurrentGlobExFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentGlobExMatchMatcher{
		egbFilter: expected,
	}
}

func (m *IsCurrentGlobExMatchMatcher) Match(actual interface{}) (bool, error) {
	node, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", node)
	}

	egbFilter, filterOk := m.egbFilter.(*filtering.GlobEx)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *IncaseFilter (%T)", egbFilter)
	}

	return egbFilter.IsMatch(node), nil
}

func (m *IsCurrentGlobExMatchMatcher) FailureMessage(actual interface{}) string {
	node, _ := actual.(*core.Node)
	egbFilter, _ := m.egbFilter.(*filtering.GlobEx)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match incase\n\t%v\n",
		node.Extension.Name, egbFilter.Source(),
	)
}

func (m *IsCurrentGlobExMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	node, _ := actual.(*core.Node)
	egbFilter, _ := m.egbFilter.(*filtering.GlobEx)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match incase\n\t%v\n",
		node.Extension.Name, egbFilter.Source(),
	)
}

// === MatchCurrentCustomFilter ===
//

type IsCurrentCustomMatchMatcher struct {
	tvFilter interface{}
}

func MatchCurrentCustomFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentCustomMatchMatcher{
		tvFilter: expected,
	}
}

func (m *IsCurrentCustomMatchMatcher) Match(actual interface{}) (bool, error) {
	node, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", node)
	}

	tvFilter, filterOk := m.tvFilter.(core.TraverseFilter)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a core.TraverseFilter (%T)", tvFilter)
	}

	return tvFilter.IsMatch(node), nil
}

func (m *IsCurrentCustomMatchMatcher) FailureMessage(actual interface{}) string {
	node, _ := actual.(*core.Node)
	tvFilter, _ := m.tvFilter.(core.TraverseFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match custom filter\n\t%v\n",
		node.Extension.Name, tvFilter.Source(),
	)
}

func (m *IsCurrentCustomMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	node, _ := actual.(*core.Node)
	tvFilter, _ := m.tvFilter.(core.TraverseFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match custom filter\n\t%v\n",
		node.Extension.Name, tvFilter.Source(),
	)
}
