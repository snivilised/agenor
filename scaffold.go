package tv

import (
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

type (
	scaffold interface {
		extent() extent
		harvest() enclave.OptionHarvest
	}

	platform struct {
		fac pref.Facade
		ext extent
		oh  enclave.OptionHarvest
	}
)

func newPrimaryPlatform(facade pref.Facade,
	addons []Addon,
	settings ...pref.Option,
) (*platform, error) {
	primary := primaryPlatform{}
	using, ok := facade.(*pref.Using)

	if !ok {
		return primary.base.pacify(facade, settings...),
			core.ErrWrongPrimaryFacade
	}

	ext := &primeExtent{
		baseExtent: baseExtent{
			fac:   using,
			trees: primary.buildForest(using),
		},
		using: using,
	}
	harvest, err := primary.buildOptions(using, addons, ext, settings...)

	return &platform{
		fac: facade,
		ext: ext,
		oh:  harvest,
	}, err
}

func newResumePlatform(facade pref.Facade,
	addons []Addon,
	settings ...pref.Option,
) (*platform, error) {
	resume := resumePlatform{}
	relic, ok := facade.(*pref.Relic)

	if !ok {
		return resume.base.pacify(facade, settings...),
			core.ErrWrongResumeFacade
	}

	ext := &resumeExtent{
		baseExtent: baseExtent{
			fac: facade,
			trees: &core.Forest{
				R: nef.NewTraverseABS(),
			},
		},
		relic: relic,
	}
	harvest, err := ext.options(addons, settings...)
	if err != nil {
		return resume.base.pacify(facade, settings...),
			err
	}

	ext.trees = resume.buildForest(relic, harvest.Loaded().State.Tree)

	return &platform{
		fac: facade,
		ext: ext,
		oh:  harvest,
	}, err
}

type basePlatform struct {
}

func (p *basePlatform) buildForest(facade pref.Facade, tree string) *core.Forest {
	fn := facade.Forest()

	return lo.TernaryF(fn != nil,
		func() *core.Forest {
			return fn(tree)
		},
		func() *core.Forest {
			// Create an absolute file system for both navigation and resume. We
			// can share the same instance because absolute fs have no state, as
			// opposed to a relative fs, which needs use the root path as state
			// which would be different for navigation and resume purposes.
			fS := nef.NewTraverseABS()

			return &core.Forest{
				T: fS,
				R: fS,
			}
		},
	)
}

func (p *basePlatform) pacify(facade pref.Facade,
	settings ...pref.Option,
) *platform {
	// this error doesn't matter because pacify is being called
	// in the presence of a prior error
	o, binder, _ := opts.Get(settings...)

	return &platform{
		fac: facade,
		ext: &primeExtent{
			baseExtent: baseExtent{
				fac: facade,
			},
		},
		oh: &optionHarvest{
			o:      o,
			binder: binder,
		},
	}
}

type primaryPlatform struct {
	base basePlatform
}

func (p *primaryPlatform) buildForest(using *pref.Using) *core.Forest {
	return p.base.buildForest(using, using.Path())
}

func (p *primaryPlatform) buildOptions(using *pref.Using,
	addons []Addon,
	ext *primeExtent,
	settings ...pref.Option,
) (oh enclave.OptionHarvest, err error) {
	type baggage struct {
		harvest enclave.OptionHarvest
		err     error
	}

	b := func(ve error) *baggage {
		return lo.TernaryF(using.O != nil,
			func() *baggage {
				return &baggage{
					harvest: &optionHarvest{
						o:      using.O,
						binder: opts.Push(using.O),
					},
					err: ve,
				}
			},
			func() *baggage {
				harvest, err := ext.options(addons, settings...)

				return &baggage{
					harvest: harvest,
					err:     lo.Ternary(ve != nil, ve, err),
				}
			},
		)
	}(using.Validate())

	return b.harvest, b.err
}

type resumePlatform struct {
	base basePlatform
}

func (p *resumePlatform) buildForest(relic *pref.Relic, tree string) *core.Forest {
	return p.base.buildForest(relic, tree)
}

func (p *platform) extent() extent {
	return p.ext
}

func (p *platform) harvest() enclave.OptionHarvest {
	return p.oh
}
