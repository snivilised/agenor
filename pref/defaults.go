package pref

import (
	"io/fs"
	"path/filepath"
	"sort"
	"strings"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
)

// DefaultReadEntriesHook reads the contents of a directory. The resulting
// slice is left un-sorted
func DefaultReadEntriesHook(sys fs.ReadDirFS,
	dirname string,
) ([]fs.DirEntry, error) {
	contents, err := fs.ReadDir(sys, dirname)
	if err != nil {
		return nil, err
	}

	return lo.Filter(contents, func(item fs.DirEntry, _ int) bool {
		return item.Name() != ".DS_Store"
	}), nil
}

// DefaultQueryStatusHook query the status of the path on the file system
// provided.
func DefaultQueryStatusHook(qsys fs.StatFS, path string) (fs.FileInfo, error) {
	return qsys.Stat(path)
}

// CaseSensitiveSortHook hook function for case sensitive directory traversal. A
// directory of "a" will be visited after a sibling directory "B".
func CaseSensitiveSortHook(entries []fs.DirEntry, _ ...any) {
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})
}

// DefaultCaseInSensitiveSortHook hook function for case insensitive directory traversal. A
// directory of "a" will be visited before a sibling directory "B".
func DefaultCaseInSensitiveSortHook(entries []fs.DirEntry, _ ...any) {
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
	return difference(info.Tree, info.Node.Path)
}

// DefaultSubPathHook
func DefaultSubPathHook(info *core.SubPathInfo) string {
	if info.Node.Extension.Scope == enums.ScopeTop {
		return lo.Ternary(info.KeepTrailingSep, string(filepath.Separator), "")
	}

	return difference(info.Tree, info.Node.Extension.Parent)
}

func DefaultFaultHandler(fault *NavigationFault) error {
	return fault.Err
}

func DefaultPanicHandler(recovery Recovery, data RescueData) (string, error) {
	return recovery.Save(data)
}

func DefaultSkipHandler(*core.Node,
	core.DirectoryContents, error,
) (enums.SkipTraversal, error) {
	return enums.SkipNoneTraversal, nil
}
