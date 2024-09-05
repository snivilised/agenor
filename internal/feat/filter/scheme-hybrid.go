package filter

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/types"
)

// hybridScheme required because node based filtering can be active at
// the same time as child filtering.
//
// The hybridScheme is related to the nanny plugin. If no filtering is active,
// then the nanny will become active and be responsible for handling the
// children. If filtering is active, then the nanny will remain dormant as
// the filter plugin becomes active. The nannyScheme will take over handling
// the children, using the child filter to do so. The primary scheme
// performs node based filtering.

type hybridScheme struct {
	common
	primary scheme
	nanny   scheme
}

func (s *hybridScheme) create() error {
	if s.primary == nil && s.nanny == nil {
		return ErrNoSubordinateHybridSchemesDefined
	}

	if s.primary != nil {
		if err := s.primary.create(); err != nil {
			return err
		}
	}

	if s.nanny != nil {
		return s.nanny.create()
	}

	return nil
}

func (s *hybridScheme) init(pi *types.PluginInit, crate *measure.Crate) {
	s.common.init(pi, crate)

	if s.primary != nil {
		s.primary.init(pi, crate)
	}

	if s.nanny != nil {
		s.nanny.init(pi, crate)
	}
}

func (s *hybridScheme) next(node *core.Node,
	inspection types.Inspection,
) (bool, error) {
	if s.primary != nil {
		invokeNext, err := s.primary.next(node, inspection)

		if invokeNext && s.nanny != nil {
			// The nanny has no say in wether the next link is invoked,
			// therefore we ignore its next result
			_, err := s.nanny.next(node, inspection)

			return invokeNext, err
		}

		return invokeNext, err
	}

	return true, nil
}
