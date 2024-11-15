package filter

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/filtering"
	"github.com/snivilised/agenor/pref"
)

type nannyScheme struct {
	common
	filter core.ChildTraverseFilter
}

func (s *nannyScheme) create() error {
	filter, err := filtering.NewChild(s.o.Filter.Child)

	if err != nil {
		return err
	}
	s.filter = filter

	if s.o.Filter.Sink != nil {
		s.o.Filter.Sink(pref.FilterReply{
			Child: s.filter,
		})
	}

	return nil
}

func (s *nannyScheme) init(pi *enclave.PluginInit, crate *enclave.Crate) {
	s.common.init(pi, crate)
}

func (s *nannyScheme) next(_ core.Servant,
	inspection enclave.Inspection,
) (bool, error) {
	files := inspection.Sort(enums.EntryTypeFile)
	matching := s.filter.Matching(files)

	inspection.AssignChildren(matching)
	s.crate.Metrics[enums.MetricNoChildFilesFound].Times(uint(len(matching)))

	filteredOut := len(files) - len(matching)
	s.crate.Metrics[enums.MetricNoChildFilesFilteredOut].Times(uint(filteredOut))

	return true, nil
}
