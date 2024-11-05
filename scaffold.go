package tv

import (
	"errors"

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
		forest() *core.Forest
	}

	platform struct {
		ext   extent
		oh    enclave.OptionHarvest
		trees *core.Forest
	}
)

func newPrimaryPlatform(facade pref.Facade, settings ...pref.Option) (*platform, error) {
	using, ok := facade.(*pref.Using)
	if !ok {
		// TODO: Create a test that accidentally sets relic facade
		//
		return nil, errors.New("incorrect facade") // TODO: create a proper error
	}

	primary := primaryPlatform{}
	ext := &primeExtent{
		baseExtent: baseExtent{
			fac:   using,
			trees: primary.buildForest(using),
		},
	}
	harvest, err := primary.buildOptions(using, ext, settings...)

	return &platform{
		ext:   ext,
		oh:    harvest,
		trees: ext.trees,
	}, err
}

func newResumePlatform(facade pref.Facade, settings ...pref.Option) (*platform, error) {
	relic, ok := facade.(*pref.Relic)
	if !ok {
		// TODO: Create a test that accidentally sets using facade
		//
		return nil, errors.New("incorrect facade") // TODO: create a proper error
	}

	resume := resumePlatform{}
	ext := &resumeExtent{
		baseExtent: baseExtent{
			fac: facade,
			trees: &core.Forest{
				R: nef.NewTraverseABS(),
			},
		},
		relic: relic,
	}
	harvest, err := ext.options(settings...)
	ext.trees = resume.buildForest(relic, harvest.Loaded().State.Tree)

	return &platform{
		ext:   ext,
		oh:    harvest,
		trees: ext.trees,
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
			return &core.Forest{
				T: nef.NewTraverseFS(Rel{
					Root:      tree,
					Overwrite: noOverwrite,
				}),
				R: nef.NewTraverseABS(),
			}
		},
	)
}

type primaryPlatform struct {
	base basePlatform
}

func (p *primaryPlatform) buildForest(using *pref.Using) *core.Forest {
	return p.base.buildForest(using, using.Path())
}

func (p *primaryPlatform) buildOptions(using *pref.Using,
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
				harvest, err := ext.options(settings...)

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

func (p *platform) forest() *core.Forest {
	return p.trees
}
