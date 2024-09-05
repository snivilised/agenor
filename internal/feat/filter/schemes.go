package filter

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type (
	scheme interface {
		create() error
		init(pi *types.PluginInit, crate *measure.Crate)
		next(node *core.Node, inspection types.Inspection) (bool, error)
	}
)

type common struct {
	o     *pref.Options
	crate *measure.Crate
}

func (f *common) init(_ *types.PluginInit, crate *measure.Crate) {
	f.crate = crate
}

func newScheme(o *pref.Options) scheme {
	c := common{o: o}
	primary, nanny := binary(o)

	if primary != nil && nanny != nil {
		return &hybridScheme{
			common:  c,
			primary: primary,
			nanny:   nanny,
		}
	}

	if primary != nil {
		return primary
	}

	return nanny
}

func binary(o *pref.Options) (primary, nanny scheme) {
	c := common{o: o}

	primary = unary(c)
	nanny = lo.TernaryF(o.Filter.IsChildFilteringActive(),
		func() scheme {
			return &nannyScheme{
				common: c,
			}
		},
		func() scheme {
			return nil
		},
	)

	if nanny == nil {
		return primary, nil
	}

	return primary, nanny
}

func unary(c common) scheme {
	if c.o.Filter.IsCustomFilteringActive() {
		return &customScheme{
			common: c,
		}
	}

	if c.o.Filter.IsNodeFilteringActive() {
		return &nativeScheme{
			common: c,
		}
	}

	if c.o.Filter.IsSampleFilteringActive() {
		return &samplerScheme{
			common: c,
		}
	}

	return nil // only nanny is active
}
