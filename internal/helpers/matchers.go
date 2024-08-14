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

	return fmt.Sprintf("🔥 Expected\n\t%v\nto match contents\n\t%v\n",
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

	return fmt.Sprintf("🔥 Expected\n\t%v\nNOT to match contents\n\t%v\n",
		strings.Join(names, ", "), strings.Join(expected, ", "),
	)
}
