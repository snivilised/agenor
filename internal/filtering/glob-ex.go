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
		segments, patterns []string
	)
	if segments, patterns, err = splitGlobExPattern(def.Pattern); err != nil {
		return nil, err
	}

	// *|*.o*.flac,f*.jpg
	directoryBase, directoryExclusion := splitGlob(segments[0])
	filter := &GlobEx{
		Base: Base{
			name:            def.Description,
			scope:           def.Scope,
			pattern:         def.Pattern,
			negate:          def.Negate,
			ifNotApplicable: ifNotApplicable,
		},
		spec: newSpec(directoryBase, directoryExclusion, patterns),
		fileGlobs: lo.Map(patterns, func(s string, _ int) string {
			return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(s), "."))
		}),
	}

	return filter, nil
}

// patternSpec represents a file pattern specification, which consists of
// the base and extension parts
type (
	patternSpec struct {
		base string
		ext  string
	}

	globSpec struct {
		specs              []*patternSpec // *|*.o*.flac,f*.jpg
		directoryGlob      string         // *
		basePatterns       []string       // *.o*,f*
		extPatterns        []string       // flac,jpg
		fileGlobs          []string       // tbd !!! needs attention
		directoryExclusion string
		anyExtension       bool
	}
)

func newSpec(directoryBase, directoryExclusion string, patterns []string) *globSpec {
	return &globSpec{
		directoryGlob:      directoryBase,
		directoryExclusion: directoryExclusion,
		fileGlobs: lo.Map(patterns, func(s string, _ int) string {
			return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(s), "."))
		}),
		anyExtension: slices.Contains(patterns, "*"),
	}
}

func (s *globSpec) IsMatch(node *core.Node) bool {
	return lo.TernaryF(node.IsDirectory(),
		func() bool {
			result, _ := filepath.Match(
				s.directoryGlob,
				strings.ToLower(node.Extension.Name),
			)

			return result
		},
		func() bool {
			return s.filter(node.Extension.Name)
		},
	)
}

func (s *globSpec) filter(name string) bool {
	extension := filepath.Ext(name)
	baseName := strings.ToLower(strings.TrimSuffix(name, extension))

	if baseMatch, _ := filepath.Match(s.directoryGlob, baseName); !baseMatch {
		return false
	}

	if excluded, _ := filepath.Match(s.directoryExclusion, baseName); excluded {
		return false
	}

	return cmp.Or(
		func() bool {
			return s.anyExtension
		}(),
		func() bool {
			return extension == "" && len(s.fileGlobs) == 0
		}(),
		func() bool {
			return lo.Contains(
				s.fileGlobs, strings.ToLower(strings.TrimPrefix(extension, ".")),
			)
		}(),
	)
}

type GlobEx struct {
	Base
	spec      *globSpec
	fileGlobs []string
}

// IsMatch does this node match the filter
func (f *GlobEx) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		result := f.spec.IsMatch(node)

		return f.invert(result)
	}

	return f.ifNotApplicable
}

// ChildGlobExFilter ================================================================

type ChildGlobExFilter struct {
	Child
	directoryGlob string
	fileGlobs     []string
	anyExtension  bool
	exclusion     string
}

func (f *ChildGlobExFilter) Matching(children []fs.DirEntry) []fs.DirEntry {
	return lo.Filter(children, func(entry fs.DirEntry, _ int) bool {
		name := entry.Name()

		return f.invert(filterFileByGlobExL(
			name, f.directoryGlob, f.exclusion, f.fileGlobs, f.anyExtension,
		))
	})
}

func filterFileByGlobExL(name, base, exclusion string,
	patterns []string, anyExtension bool,
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
			return extension == "" && len(patterns) == 0
		}(),
		func() bool {
			return lo.Contains(
				patterns, strings.ToLower(strings.TrimPrefix(extension, ".")),
			)
		}(),
	)
}
