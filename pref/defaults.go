package pref

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/lo"
)

// DefaultReadEntriesHook reads the contents of a directory. The resulting
// slice is left un-sorted
func DefaultReadEntriesHook(sys fs.FS, dirname string) ([]fs.DirEntry, error) {
	const all = -1

	contents, err := fs.ReadDir(sys, dirname)
	if err != nil {
		return nil, err
	}

	return lo.Filter(contents, func(item fs.DirEntry, _ int) bool {
		return item.Name() != ".DS_Store"
	}), nil
}

// CaseSensitiveSortHook hook function for case sensitive directory traversal. A
// directory of "a" will be visited after a sibling directory "B".
func CaseSensitiveSortHook(entries []fs.DirEntry, _ ...any) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
}

// CaseInSensitiveSortHook hook function for case insensitive directory traversal. A
// directory of "a" will be visited before a sibling directory "B".
func CaseInSensitiveSortHook(entries []fs.DirEntry, _ ...any) {
	sort.Slice(entries, func(i, j int) bool {
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})
}

// tail extracts the end part of a string, starting from the offset
func tail(input string, offset int) string {
	asRunes := []rune(input)

	if offset >= len(asRunes) {
		return ""
	}

	return string(asRunes[offset:])
}

// difference returns the difference between a child path and a parent path
// Designed to be used with paths created from the file system rather than
// custom created or user provided input. For this reason, if there is no
// relationship between the parent and child paths provided then a panic
// may occur.
func difference(parent, child string) string {
	return tail(child, len(parent))
}

// RootItemSubPathHook
func RootItemSubPathHook(info *core.SubPathInfo) string {
	return difference(info.Root, info.Node.Path)
}

// RootParentSubPathHook
func RootParentSubPathHook(info *core.SubPathInfo) string {
	if info.Node.Extension.Scope == enums.ScopeTop {
		return lo.Ternary(info.KeepTrailingSep, string(filepath.Separator), "")
	}

	return difference(info.Root, info.Node.Extension.Parent)
}

func DefaultFaultHandler(*NavigationFault) error {
	return nil
}

func DefaultPanicHandler() {
	// may this should invoke save
}

func DefaultSkipHandler(*core.Node, core.DirectoryContents, error) (enums.SkipTraversal, error) {
	return enums.SkipNoneTraversal, nil
}
