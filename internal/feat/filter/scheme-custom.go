package filter

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type customScheme struct {
	common
	filter core.TraverseFilter
}

func (s *customScheme) create() error {
	s.filter = s.o.Filter.Custom

	if s.o.Filter.Sink != nil {
		s.o.Filter.Sink(pref.FilterReply{
			Node: s.filter,
		})
	}

	return s.filter.Validate()
}

func (s *customScheme) next(node *core.Node,
	_ types.Inspection,
) (bool, error) {
	return matchNext(s.filter, node, s.crate)
}

func matchNext(filter core.TraverseFilter,
	node *core.Node, crate *measure.Crate,
) (bool, error) {
	matched := filter.IsMatch(node)

	if !matched {
		filteredOutMetric := lo.Ternary(node.IsFolder(),
			enums.MetricNoFoldersFilteredOut,
			enums.MetricNoFilesFilteredOut,
		)
		crate.Mums[filteredOutMetric].Tick()
	}

	return matched, nil
}
