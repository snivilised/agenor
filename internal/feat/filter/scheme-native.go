package filter

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/filtering"
	"github.com/snivilised/traverse/pref"
)

type nativeScheme struct {
	common
	filter core.TraverseFilter
}

func (s *nativeScheme) create() error {
	filter, err := filtering.New(s.o.Filter.Node, &s.o.Filter)
	if err != nil {
		return err
	}

	s.filter = filter

	if s.o.Filter.Sink != nil {
		s.o.Filter.Sink(pref.FilterReply{
			Node: s.filter,
		})
	}

	return nil
}

func (s *nativeScheme) next(servant core.Servant,
	_ enclave.Inspection,
) (bool, error) {
	return matchNext(s.filter, servant.Node(), s.crate)
}
