package filtering

import (
	"cmp"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/third/lo"
)

func createGlobExFilter(def *core.FilterDef,
	ifNotApplicable bool,
) (core.TraverseFilter, error) {
	var (
		err                error
		segments, suffixes []string
	)
	if segments, suffixes, err = splitGlobExPattern(def.Pattern); err != nil {
		return nil, err
	}

	base, exclusion := splitGlob(segments[0])

	filter := &GlobEx{
		Base: Base{
			name:            def.Description,
			scope:           def.Scope,
			pattern:         def.Pattern,
			negate:          def.Negate,
			ifNotApplicable: ifNotApplicable,
		},
		baseGlob: base,
		suffixes: lo.Map(suffixes, func(s string, _ int) string {
			return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(s), "."))
		}),
		anyExtension: slices.Contains(suffixes, "*"),
		exclusion:    exclusion,
	}

	return filter, nil
}

type GlobEx struct {
	Base
	baseGlob     string
	suffixes     []string
	anyExtension bool
	exclusion    string
}

// IsMatch does this node match the filter
func (f *GlobEx) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		result := lo.TernaryF(node.IsDirectory(),
			func() bool {
				result, _ := filepath.Match(f.baseGlob, strings.ToLower(node.Extension.Name))

				return result
			},
			func() bool {
				return filterFileByGlobEx(
					node.Extension.Name, f.baseGlob, f.exclusion, f.suffixes, f.anyExtension,
				)
			},
		)

		return f.invert(result)
	}

	return f.ifNotApplicable
}

// ChildGlobExFilter ================================================================

type ChildGlobExFilter struct {
	Child
	baseGlob     string
	exclusion    string
	suffixes     []string
	anyExtension bool
}

func (f *ChildGlobExFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, _ int) bool {
		name := entry.Name()

		return f.invert(filterFileByGlobEx(
			name, f.baseGlob, f.exclusion, f.suffixes, f.anyExtension,
		))
	})
}

func filterFileByGlobEx(name, base, exclusion string,
	suffixes []string, anyExtension bool,
) bool {
	extension := filepath.Ext(name)
	baseName := strings.ToLower(strings.TrimSuffix(name, extension))

	if baseMatch, _ := filepath.Match(base, baseName); !baseMatch {
		return false
	}

	if excluded, _ := filepath.Match(exclusion, baseName); excluded {
		return false
	}

	return cmp.Or(
		func() bool {
			return anyExtension
		}(),
		func() bool {
			return extension == "" && len(suffixes) == 0
		}(),
		func() bool {
			return lo.Contains(
				suffixes, strings.ToLower(strings.TrimPrefix(extension, ".")),
			)
		}(),
	)
}
