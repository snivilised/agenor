package filter

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/third/lo"
)

type customScheme struct {
	common
	filter core.TraverseFilter
}

func (s *customScheme) create() error {
	s.filter = s.o.Filter.Custom

	return s.filter.Validate()
}

func (s *customScheme) next(servant core.Servant,
	_ enclave.Inspection,
) (bool, error) {
	return matchNext(s.filter, servant.Node(), s.crate)
}

func matchNext(filter core.TraverseFilter,
	node *core.Node, crate *core.Crate,
) (bool, error) {
	matched := filter.IsMatch(node)

	if !matched {
		filteredOutMetric := lo.Ternary(node.IsDirectory(),
			enums.MetricNoDirectoriesFilteredOut,
			enums.MetricNoFilesFilteredOut,
		)
		crate.Metrics[filteredOutMetric].Tick()
	}

	return matched, nil
}
