package kernel_test

import (
	"fmt"

	. "github.com/onsi/gomega/types" //nolint:revive // ok
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/refine"
)

// === MatchCurrentRegexFilter ===
//

type IsCurrentRegexMatchMatcher struct {
	filter interface{}
}

func MatchCurrentRegexFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentRegexMatchMatcher{
		filter: expected,
	}
}

func (m *IsCurrentRegexMatchMatcher) Match(actual interface{}) (bool, error) {
	item, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	filter, filterOk := m.filter.(*refine.RegexFilter)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *RegexFilter (%T)", filter)
	}

	return filter.IsMatch(item), nil
}

func (m *IsCurrentRegexMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	filter, _ := m.filter.(*refine.RegexFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match regex\n\t%v\n", item.Extension.Name, filter.Source())
}

func (m *IsCurrentRegexMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	filter, _ := m.filter.(*refine.RegexFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match regex\n\t%v\n", item.Extension.Name, filter.Source())
}

// === MatchCurrentGlobFilter ===
//

type IsCurrentGlobMatchMatcher struct {
	filter interface{}
}

func MatchCurrentGlobFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentGlobMatchMatcher{
		filter: expected,
	}
}

func (m *IsCurrentGlobMatchMatcher) Match(actual interface{}) (bool, error) {
	item, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	filter, filterOk := m.filter.(*refine.GlobFilter)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *GlobFilter (%T)", filter)
	}

	return filter.IsMatch(item), nil
}

func (m *IsCurrentGlobMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	filter, _ := m.filter.(*refine.GlobFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match glob\n\t%v\n", item.Extension.Name, filter.Source())
}

func (m *IsCurrentGlobMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	filter, _ := m.filter.(*refine.GlobFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match glob\n\t%v\n", item.Extension.Name, filter.Source())
}

// === MatchCurrentExtendedGlobFilter ===
//

type IsCurrentExtendedGlobMatchMatcher struct {
	filter interface{}
}

func MatchCurrentExtendedFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentExtendedGlobMatchMatcher{
		filter: expected,
	}
}

func (m *IsCurrentExtendedGlobMatchMatcher) Match(actual interface{}) (bool, error) {
	item, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	filter, filterOk := m.filter.(*refine.ExtendedGlobFilter)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a *IncaseFilter (%T)", filter)
	}

	return filter.IsMatch(item), nil
}

func (m *IsCurrentExtendedGlobMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	filter, _ := m.filter.(*refine.ExtendedGlobFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match incase\n\t%v\n",
		item.Extension.Name, filter.Source(),
	)
}

func (m *IsCurrentExtendedGlobMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	filter, _ := m.filter.(*refine.ExtendedGlobFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match incase\n\t%v\n",
		item.Extension.Name, filter.Source(),
	)
}

// === MatchCurrentCustomFilter ===
//

type IsCurrentCustomMatchMatcher struct {
	filter interface{}
}

func MatchCurrentCustomFilter(expected interface{}) GomegaMatcher {
	return &IsCurrentCustomMatchMatcher{
		filter: expected,
	}
}

func (m *IsCurrentCustomMatchMatcher) Match(actual interface{}) (bool, error) {
	item, itemOk := actual.(*core.Node)
	if !itemOk {
		return false, fmt.Errorf("matcher expected a *TraverseItem (%T)", item)
	}

	filter, filterOk := m.filter.(core.TraverseFilter)
	if !filterOk {
		return false, fmt.Errorf("matcher expected a core.TraverseFilter (%T)", filter)
	}

	return filter.IsMatch(item), nil
}

func (m *IsCurrentCustomMatchMatcher) FailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	filter, _ := m.filter.(core.TraverseFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match custom filter\n\t%v\n",
		item.Extension.Name, filter.Source(),
	)
}

func (m *IsCurrentCustomMatchMatcher) NegatedFailureMessage(actual interface{}) string {
	item, _ := actual.(*core.Node)
	filter, _ := m.filter.(core.TraverseFilter)

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match custom filter\n\t%v\n",
		item.Extension.Name, filter.Source(),
	)
}
