package filter

import (
	"io/fs"

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/filtering"
	"github.com/snivilised/traverse/internal/third/lo"
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

func (s *samplerScheme) next(servant core.Servant,
	inspection enclave.Inspection,
) (bool, error) {
	node := servant.Node()

	if node.Extension.Scope.IsTree() {
		matching := s.filter.Matching(
			[]fs.DirEntry{nef.FromFileInfo(node.Info)},
		)
		result := len(matching) > 0

		lo.Ternary(result,
			s.crate.Metrics[enums.MetricNoChildFilesFound],
			s.crate.Metrics[enums.MetricNoChildFilesFilteredOut],
		).Times(uint(len(inspection.Contents().Files())))

		return result, nil
	}

	return true, nil
}
