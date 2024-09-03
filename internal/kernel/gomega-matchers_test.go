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
	item, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	rxFilter, filterOk := m.rxFilter.(*filtering.RegExpr)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *RegexFilter (%T)", rxFilter)
	}

	return rxFilter.IsMatch(item), nil
}

func (m *IsCurrentRegexMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	rxFilter, _ := m.rxFilter.(*filtering.RegExpr)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match regex\n\t%v\n",
		item.Extension.Name, rxFilter.Source(),
	)
}

func (m *IsCurrentRegexMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	rxFilter, _ := m.rxFilter.(*filtering.RegExpr)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match regex\n\t%v\n",
		item.Extension.Name, rxFilter.Source(),
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
	item, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	gbFilter, filterOk := m.gbFilter.(*filtering.Glob)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *GlobFilter (%T)", gbFilter)
	}

	return gbFilter.IsMatch(item), nil
}

func (m *IsCurrentGlobMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	gbFilter, _ := m.gbFilter.(*filtering.Glob)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match glob\n\t%v\n",
		item.Extension.Name, gbFilter.Source(),
	)
}

func (m *IsCurrentGlobMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	gbFilter, _ := m.gbFilter.(*filtering.Glob)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match glob\n\t%v\n",
		item.Extension.Name, gbFilter.Source(),
	)
}

// === MatchCurrentExtendedGlobFilter ===
//

type IsCurrentExtendedGlobMatchMatcher struct {
	egbFilter interface{}
}

func MatchCurrentExtendedFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentExtendedGlobMatchMatcher{
		egbFilter: expected,
	}
}

func (m *IsCurrentExtendedGlobMatchMatcher) Match(actual interface{}) (bool, error) {
	item, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	egbFilter, filterOk := m.egbFilter.(*filtering.ExtendedGlob)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *IncaseFilter (%T)", egbFilter)
	}

	return egbFilter.IsMatch(item), nil
}

func (m *IsCurrentExtendedGlobMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	egbFilter, _ := m.egbFilter.(*filtering.ExtendedGlob)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match incase\n\t%v\n",
		item.Extension.Name, egbFilter.Source(),
	)
}

func (m *IsCurrentExtendedGlobMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	egbFilter, _ := m.egbFilter.(*filtering.ExtendedGlob)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match incase\n\t%v\n",
		item.Extension.Name, egbFilter.Source(),
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
	item, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	tvFilter, filterOk := m.tvFilter.(core.TraverseFilter)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a core.TraverseFilter (%T)", tvFilter)
	}

	return tvFilter.IsMatch(item), nil
}

func (m *IsCurrentCustomMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	tvFilter, _ := m.tvFilter.(core.TraverseFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match custom filter\n\t%v\n",
		item.Extension.Name, tvFilter.Source(),
	)
}

func (m *IsCurrentCustomMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	tvFilter, _ := m.tvFilter.(core.TraverseFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match custom filter\n\t%v\n",
		item.Extension.Name, tvFilter.Source(),
	)
}
