package helpers

import (
	"fmt"
	"io/fs"
	"slices"
	"strings"

	. "github.com/onsi/gomega/types" //nolint:stylecheck,revive // ok
	"github.com/snivilised/traverse/internal/third/lo"
)

type DirectoryContentsMatcher struct {
	expected interface{}
}

func HaveDirectoryContents(expected interface{}) GomegaMatcher {
	return &DirectoryContentsMatcher{
		expected: expected,
	}
}

func (m *DirectoryContentsMatcher) Match(actual interface{}) (bool, error) {
	entries, entriesOk := actual.([]fs.DirEntry)
	if !entriesOk {
		return false, fmt.Errorf("matcher expected []fs.DirEntry (%T)", entries)
	}

	names := lo.Map(entries, func(entry fs.DirEntry, _ int) string {
		return entry.Name()
	})

	expected, expectedOk := m.expected.([]string)
	if !expectedOk {
		return false, fmt.Errorf("matcher expected []string (%T)", expected)
	}

	return slices.Compare(names, expected) == 0, nil
}

func (m *DirectoryContentsMatcher) FailureMessage(actual interface{}) string {
	entries, _ := actual.([]fs.DirEntry)
	names := lo.Map(entries, func(entry fs.DirEntry, _ int) string {
		return entry.Name()
	})
	slices.Sort(names)

	expected, _ := m.expected.([]string)
	slices.Sort(expected)

	return fmt.Sprintf("❌ Expected\n\t%v\nto match contents\n\t%v\n",
		strings.Join(names, ", "), strings.Join(expected, ", "),
	)
}

func (m *DirectoryContentsMatcher) NegatedFailureMessage(actual interface{}) string {
	entries, _ := actual.([]fs.DirEntry)
	names := lo.Map(entries, func(entry fs.DirEntry, _ int) string {
		return entry.Name()
	})
	slices.Sort(names)

	expected, _ := m.expected.([]string)
	slices.Sort(expected)

	return fmt.Sprintf("❌ Expected\n\t%v\nNOT to match contents\n\t%v\n",
		strings.Join(names, ", "), strings.Join(expected, ", "),
	)
}

type InvokeNodeMatcher struct {
	expected interface{}
}

func HaveInvokedNode(expected interface{}) GomegaMatcher {
	return &InvokeNodeMatcher{
		expected: expected,
	}
}

func (m *InvokeNodeMatcher) Match(actual interface{}) (bool, error) {
	recording, ok := actual.(RecordingMap)
	if !ok {
		return false, fmt.Errorf("matcher expected actual to be a RecordingMap (%T)", actual)
	}

	mandatory, ok := m.expected.(string)
	if !ok {
		return false, fmt.Errorf("matcher expected string (%T)", actual)
	}

	_, found := recording[mandatory]

	if !found {
		return false, fmt.Errorf("❌ missing invoke for node: '%v'", mandatory)
	}

	return true, nil
}

func (m *InvokeNodeMatcher) FailureMessage(_ interface{}) string {
	mandatory, ok := m.expected.(string)
	if !ok {
		return fmt.Sprintf("🔥 matcher expected string (%T)", m.expected)
	}

	return fmt.Sprintf("❌ Expected\n\t%v\nnode to be invoked\n",
		mandatory,
	)
}

func (m *InvokeNodeMatcher) NegatedFailureMessage(_ interface{}) string {
	mandatory, ok := m.expected.(string)
	if !ok {
		return fmt.Sprintf("🔥 matcher expected string (%T)", m.expected)
	}

	return fmt.Sprintf("❌ Expected\n\t%v\nnode NOT to be invoked\n",
		mandatory,
	)
}

type NotInvokeNodeMatcher struct {
	expected interface{}
}

func HaveNotInvokedNode(expected interface{}) GomegaMatcher {
	return &NotInvokeNodeMatcher{
		expected: expected,
	}
}

func (m *NotInvokeNodeMatcher) Match(actual interface{}) (bool, error) {
	recording, ok := actual.(RecordingMap)
	if !ok {
		return false, fmt.Errorf("matcher expected actual to be a RecordingMap (%T)", actual)
	}

	mandatory, ok := m.expected.(string)
	if !ok {
		return false, fmt.Errorf("matcher expected string (%T)", actual)
	}

	_, found := recording[mandatory]

	if found {
		return false, fmt.Errorf("❌ prohibited invoke occurred for node: '%v'", mandatory)
	}

	return true, nil
}

func (m *NotInvokeNodeMatcher) FailureMessage(_ interface{}) string {
	mandatory, ok := m.expected.(string)
	if !ok {
		return fmt.Sprintf("🔥 matcher expected string (%T)", m.expected)
	}

	return fmt.Sprintf("❌ Expected\n\t%v\nnode to NOT be invoked\n",
		mandatory,
	)
}

func (m *NotInvokeNodeMatcher) NegatedFailureMessage(_ interface{}) string {
	mandatory, ok := m.expected.(string)
	if !ok {
		return fmt.Sprintf("🔥 matcher expected string (%T)", m.expected)
	}

	return fmt.Sprintf("❌ Expected\n\t%v\nnode to be invoked\n",
		mandatory,
	)
}

type ExpectedCount struct {
	Name  string
	Count int
}

type ChildCountMatcher struct {
	expected interface{}
}

func HaveChildCountOf(expected interface{}) GomegaMatcher {
	return &ChildCountMatcher{
		expected: expected,
	}
}

func (m *ChildCountMatcher) Match(actual interface{}) (bool, error) {
	recording, ok := actual.(RecordingMap)
	if !ok {
		return false, fmt.Errorf("matcher expected actual to be a RecordingMap (%T)", actual)
	}

	expected, ok := m.expected.(ExpectedCount)
	if !ok {
		return false, fmt.Errorf("matcher expected ExpectedCount (%T)", actual)
	}

	count, ok := recording[expected.Name]
	if !ok {
		return false, fmt.Errorf("🔥 not found: '%v'", expected.Name)
	}

	if count != expected.Count {
		return false, fmt.Errorf(
			"❌ incorrect child count for: '%v', actual: '%v', expected: '%v'",
			expected.Name,
			count, expected.Count,
		)
	}

	return true, nil
}

func (m *ChildCountMatcher) FailureMessage(_ interface{}) string {
	expected, ok := m.expected.(ExpectedCount)
	if !ok {
		return fmt.Sprintf("🔥 matcher expected ExpectedCount (%T)", m.expected)
	}

	return fmt.Sprintf("❌ Expected\n\t%v\nnode to be invoked\n",
		expected,
	)
}

func (m *ChildCountMatcher) NegatedFailureMessage(_ interface{}) string {
	expected, ok := m.expected.(ExpectedCount)
	if !ok {
		return fmt.Sprintf("🔥 matcher expected ExpectedCount (%T)", m.expected)
	}

	return fmt.Sprintf("❌ Expected\n\t%v\nnode NOT to be invoked\n",
		expected,
	)
}
