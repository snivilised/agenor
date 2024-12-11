package filtering

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/third/lo"
)

const (
	excludeExtSymbol = '!'
	wildcard         = "*"
	GlobExPattern    = `^(?P<base>.*?)\.(?P<neg>!)?(?P<ext>[^.]+)$`
)

func createGlobExFilter(def *core.FilterDef,
	ifNotApplicable bool,
) (filter core.TraverseFilter, err error) {
	var (
		segments, patterns []string
		gs                 *globSpec
	)
	if segments, patterns, err = splitGlobExPattern(def.Pattern); err != nil {
		return nil, err
	}

	// we ignore the scope on the filter because use of this type
	// of filter implies the correct scope, ie file scope.
	def.Scope = enums.ScopeFile

	directoryBase, directoryExclusion := splitGlob(segments[0])

	gs, err = newSpec(directoryBase, directoryExclusion, patterns)
	filter = &GlobEx{
		Base: Base{
			name:            def.Description,
			scope:           def.Scope,
			pattern:         def.Pattern,
			negate:          def.Negate,
			ifNotApplicable: ifNotApplicable,
		},
		spec: gs,
	}

	return filter, err
}

type GlobEx struct {
	Base
	spec *globSpec
}

func (f *GlobEx) IsApplicable(node *core.Node) bool {
	if f.Base.IsApplicable(node) {
		return f.isDirectoryMatch(node)
	}

	return false
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

// IsMatch does this node match the filter; It's important to note that
// IfNotApplicable on the filter definition should be set correctly to avoid
// an unexpected result in certain situations. By default, IfNotApplicable = true,
// but this behaviour could cause confusion if not well understood. If for example
// using a filter pattern like "c*|*.*", (which means any file whose parent
// directory starts with "c"), files whose parent directory doesn't match
// this condition still pass through. This is because, when this filter is
// not applicable to non conforming files, returning IfNotApplicable's default
// value of true, will permit the non conforming file to be invoked. To resolve
// this situation, the client should also set the filter's IfNotApplicable to false.
func (f *GlobEx) IsMatch(node *core.Node) bool {
	if f.IsApplicable(node) {
		return f.invert(f.spec.IsMatch(node.Extension.Name))
	}

	return f.ifNotApplicable
}

// patternSpec represents a file pattern specification, which consists of
// the base and extension parts
type (
	patternSpec struct {
		base       string
		ext        string
		grossExt   string
		matcher    string
		excludeExt bool
	}

	globSpec struct {
		specs              []*patternSpec // */<excl>|*.o*.flac,f*.jpg
		directoryGlob      string         // *
		directoryExclusion string         // <excl>
	}

	indexCaptures struct {
		dot, star int
		base, ext string
	}

	extractor  func(pattern string) (ext, gross, match string)
	indexFuncs map[enums.GlobExtraction]extractor
)

func parse(pattern string, re *regexp.Regexp) (spec *patternSpec, err error) {
	subMatches := re.FindStringSubmatch(pattern)

	if subMatches == nil {
		return nil, fmt.Errorf("invalid glob-ex sub-pattern '%v' (match failed)", pattern)
	}

	result := make(map[string]string)

	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = subMatches[i]
		}
	}

	ext := result["ext"]
	base := result["base"]
	grossExt := "." + ext

	return &patternSpec{
		base:       base,
		ext:        ext,
		grossExt:   grossExt,
		excludeExt: result["neg"] == string(excludeExtSymbol),
		matcher: lo.Ternary(base == "",
			wildcard+grossExt,
			base+grossExt,
		),
	}, nil
}

func newSpec(directoryBase, directoryExclusion string, patterns []string) (*globSpec, error) {
	re := regexp.MustCompile(GlobExPattern)

	specs, err := lo.MapE(patterns, func(pattern string, _ int) (*patternSpec, error) {
		pattern = strings.ToLower(strings.TrimSpace(pattern))

		return parse(pattern, re)
	})

	if err != nil {
		return nil, err
	}

	return &globSpec{
		specs:              specs,
		directoryGlob:      directoryBase,
		directoryExclusion: directoryExclusion,
	}, nil
}

func (s *patternSpec) excluded(name string) bool {
	if !s.excludeExt {
		return false
	}

	m, _ := filepath.Match(s.matcher, name)
	return m
}

func (s *patternSpec) match(name string) bool {
	m, _ := filepath.Match(s.matcher, name)
	return m
}

func (s globSpec) IsMatch(name string) bool { //nolint: gocritic
	for _, spec := range s.specs {
		if spec.excluded(name) {
			return false
		}

		if spec.excludeExt || spec.match(name) {
			return true
		}
	}

	return false
}
