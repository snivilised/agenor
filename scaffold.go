package age

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/tfs"
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
				R: tfs.New(),
			},
		},
		relic: relic,
	}
	harvest, err := ext.options(addons, settings...)
	if err != nil {
		return resume.base.pacify(facade, settings...),
			err
	}

	ext.trees = resume.buildForest(relic, harvest.Loaded().State)
	if ext.trees == nil {
		return resume.base.pacify(facade, settings...), core.ErrNilForest
	}

	if ext.trees.T == nil {
		err = locale.NewTraverseFsMismatchError()
	} else if ext.trees.R == nil {
		err = locale.NewResumeFsMismatchError()
	}

	return &platform{
		fac: facade,
		ext: ext,
		oh:  harvest,
	}, err
}

type basePlatform struct {
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
	fn := using.Forest()

	if fn != nil {
		return fn(using.Tree)
	}
	// Create an absolute file system for both navigation and resume. We
	// can share the same instance because absolute fs have no state, as
	// opposed to a relative fs, which needs use the root path as state
	// which would be different for navigation and resume purposes.
	fS := tfs.New()

	return &core.Forest{
		T: fS,
		R: fS,
	}
}

func (p *primaryPlatform) buildOptions(using *pref.Using,
	addons []Addon,
	ext *primeExtent,
	settings ...pref.Option,
) (enclave.OptionHarvest, error) {
	return func(ve error) (enclave.OptionHarvest, error) {
		if using.O != nil {
			return &optionHarvest{
				o:      using.O,
				binder: opts.Push(using.O),
			}, ve
		}

		harvest, err := ext.options(addons, settings...)

		return harvest, lo.Ternary(ve != nil, ve, err)
	}(using.Validate())
}

type resumePlatform struct {
	base basePlatform
}

func (p *resumePlatform) buildForest(relic *pref.Relic, active *core.ActiveState) *core.Forest {
	fn := relic.Forest()
	forest := lo.TernaryF(fn != nil,
		func() *core.Forest {
			return fn(active.Tree)
		},
		func() *core.Forest {
			// Create an absolute file system for both navigation and resume. We
			// can share the same instance because absolute fs have no state, as
			// opposed to a relative fs, which needs to use the root path as state
			// tha is different for navigation and resume purposes.
			fS := tfs.New()

			return &core.Forest{
				T: fS,
				R: fS,
			}
		},
	)

	if forest == nil {
		return nil
	}

	// The client may not have provided a function to create the forest,
	// instead relying on the default, but as we create an absolute fs by
	// default, this may not match the type of traverse fs created by
	// the original traversal session as per the active state loaded.
	// If this is the case, then the client should provide a file system
	// that matches the loaded active state from the previous session.
	//

	if forest.T.IsRelative() != active.TraverseDescription.IsRelative {
		forest.T = nil
	}

	if forest.R.IsRelative() != active.ResumeDescription.IsRelative {
		forest.R = nil
	}

	return forest
}

func (p *platform) extent() extent {
	return p.ext
}

func (p *platform) harvest() enclave.OptionHarvest {
	return p.oh
}
