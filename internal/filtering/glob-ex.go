package filtering

import (
	"cmp"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/third/lo"
)

const (
	excludeExtSymbol = '!'
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

	// ensure that the scope is appropriate for glob-ex
	//
	def.Scope.Clear(enums.ScopeDirectory)
	def.Scope.Set(enums.ScopeFile)

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

type GlobEx struct {
	Base
	spec      *globSpec
	fileGlobs []string
}

func (f *GlobEx) IsApplicable(node *core.Node) bool {
	if f.Base.IsApplicable(node) {
		return f.isDirectoryMatch(node)
	}

	return true
}

func (f *GlobEx) isDirectoryMatch(node *core.Node) bool {
	if node.Parent != nil {
		name := strings.ToLower(node.Parent.Extension.Name)

		if result, _ := filepath.Match(
			f.spec.directoryGlob, name,
		); !result {
			return false
		}

		if excluded, _ := filepath.Match(
			f.spec.directoryExclusion, name,
		); excluded {
			return false
		}
	}

	return true
}

// IsMatch does this node match the filter
func (f *GlobEx) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		return f.invert(f.spec.IsMatch(node))
	}

	return f.ifNotApplicable
}

// patternSpec represents a file pattern specification, which consists of
// the base and extension parts
type (
	patternSpec struct {
		pattern    string
		ext        string
		grossExt   string
		excludeExt bool
	}

	globSpec struct {
		specs              []*patternSpec // */<excl>|*.o*.flac,f*.jpg
		directoryGlob      string         // *
		directoryExclusion string         // <excl>
	}
)

func newSpec(directoryBase, directoryExclusion string, patterns []string) *globSpec {
	specs := lo.Map(patterns, func(s string, _ int) *patternSpec {
		s = strings.ToLower(strings.TrimSpace(s))
		dot := strings.LastIndex(s, ".")

		if dot == -1 { // pattern denotes an extension only; dot not specified: !jpg/jpg
			const minWithoutDot = 1

			if len(s) == minWithoutDot { // single letter extension
				return &patternSpec{
					pattern:  s,
					ext:      s,
					grossExt: "." + s,
				}
			}

			excludedExt := s[0] == excludeExtSymbol
			return &patternSpec{
				ext:        lo.Ternary(excludedExt, s[1:], s),
				excludeExt: excludedExt,
			}
		}

		if dot == 0 { // pattern denotes an extension only; dot specified: .!jpg/.jpg
			const minWithDot = 2

			if len(s) <= minWithDot { // single letter extension
				if s[1] == excludeExtSymbol {
					// .! this is invalid, so ignore
					return nil
				}

				pattern := s[1:]
				return &patternSpec{
					pattern:  pattern,
					ext:      pattern,
					grossExt: s,
				}
			}

			excludedExt := s[1] == excludeExtSymbol
			grossExt := lo.Ternary(excludedExt, "."+s[2:], s)
			ext := grossExt[1:]
			return &patternSpec{
				pattern:    grossExt,
				ext:        ext,
				grossExt:   grossExt,
				excludeExt: excludedExt,
			}
		}

		// *.o*.!flac/*.o*.flac
		excludedExt := s[dot+1] == excludeExtSymbol
		base := s[:dot]
		grossExt := lo.Ternary(excludedExt, "."+s[dot+2:], s[dot+1:])
		ext := grossExt[1:]
		pattern := base + ext
		return &patternSpec{ // pattern denotes both parts
			pattern:    pattern,
			ext:        ext,
			grossExt:   grossExt,
			excludeExt: excludedExt,
		}
	})

	specs = lo.Reject(specs, func(item *patternSpec, _ int) bool {
		return item == nil
	})

	return &globSpec{
		specs:              specs,
		directoryGlob:      directoryBase,
		directoryExclusion: directoryExclusion,
	}
}

func (s *patternSpec) excluded(name string) bool {
	if !s.excludeExt {
		return false
	}

	m, _ := filepath.Match(name, s.pattern)
	return m
}

func (s *patternSpec) match(name string) bool {
	m, _ := filepath.Match(name, s.pattern)
	return m
}

func (s globSpec) IsMatch(node *core.Node) bool { //nolint: gocritic
	name := node.Extension.Name

	for _, spec := range s.specs {
		if spec.excluded(name) {
			return false
		}

		if !spec.match(name) {
			return false
		}
	}

	return true
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
