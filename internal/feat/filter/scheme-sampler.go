package filter

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/filtering"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/nfs"
)

type samplerScheme struct {
	common
	filter core.SampleTraverseFilter
}

func (s *samplerScheme) create() error {
	filter, err := filtering.NewSample(s.o.Filter.Sample, &s.o.Sampling)

	if err != nil {
		return err
	}

	s.filter = filter

	// the filter plugin performs premature filtering (with fs.DirEntry as opposed
	// to core.Node) on behalf of the sampler.
	s.o.Hooks.ReadDirectory.Chain(
		func(result []fs.DirEntry, err error,
			_ fs.ReadDirFS, _ string,
		) ([]fs.DirEntry, error) {
			return s.filter.Matching(result), err
		},
	)

	return filter.Validate()
}

func (s *samplerScheme) next(node *core.Node,
	inspection types.Inspection,
) (bool, error) {
	if node.Extension.Scope.IsRoot() {
		matching := s.filter.Matching(
			[]fs.DirEntry{nfs.FromFileInfo(node.Info)},
		)
		result := len(matching) > 0

		lo.Ternary(result,
			s.crate.Mums[enums.MetricNoChildFilesFound],
			s.crate.Mums[enums.MetricNoChildFilesFilteredOut],
		).Times(uint(len(inspection.Contents().Files())))

		return result, nil
	}

	return true, nil
}
