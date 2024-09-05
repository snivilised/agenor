package filtering

import (
	"cmp"
	"io/fs"
	"path/filepath"
	"slices"
	"strings"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/third/lo"
)

func createExtendedGlobFilter(def *core.FilterDef,
	ifNotApplicable bool,
) (core.TraverseFilter, error) {
	var (
		err                error
		segments, suffixes []string
	)
	if segments, suffixes, err = splitExtendedGlobPattern(def.Pattern); err != nil {
		return nil, err
	}

	base, exclusion := splitGlob(segments[0])

	filter := &ExtendedGlob{
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

type ExtendedGlob struct {
	Base
	baseGlob     string
	suffixes     []string
	anyExtension bool
	exclusion    string
}

// IsMatch does this node match the filter
func (f *ExtendedGlob) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		result := lo.TernaryF(node.IsFolder(),
			func() bool {
				result, _ := filepath.Match(f.baseGlob, strings.ToLower(node.Extension.Name))

				return result
			},
			func() bool {
				return filterFileByExtendedGlob(
					node.Extension.Name, f.baseGlob, f.exclusion, f.suffixes, f.anyExtension,
				)
			},
		)

		return f.invert(result)
	}

	return f.ifNotApplicable
}

// ChildExtendedGlobFilter ==========================================================

type ChildExtendedGlobFilter struct {
	Child
	baseGlob     string
	exclusion    string
	suffixes     []string
	anyExtension bool
}

func (f *ChildExtendedGlobFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, _ int) bool {
		name := entry.Name()

		return f.invert(filterFileByExtendedGlob(
			name, f.baseGlob, f.exclusion, f.suffixes, f.anyExtension,
		))
	})
}

func filterFileByExtendedGlob(name, base, exclusion string,
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
