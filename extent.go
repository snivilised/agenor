package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/feat/resume"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/tfs"
)

type extent interface {
	facade() pref.Facade
	plugin(*kernel.Artefacts) enclave.Plugin
	options(...pref.Option) (enclave.OptionHarvest, error)
	forest() *core.Forest
	complete() bool
}

type fileSystems struct {
	fS tfs.TraversalFS
}

type baseExtent struct {
	trees *core.Forest
	fac   pref.Facade
}

func (x *baseExtent) forest() *core.Forest {
	return x.trees
}

func (x *baseExtent) facade() pref.Facade {
	return x.fac
}

type primeExtent struct {
	baseExtent
	using *pref.Using
}

func (x *primeExtent) plugin(*kernel.Artefacts) enclave.Plugin {
	return nil
}

func (x *primeExtent) options(settings ...pref.Option) (enclave.OptionHarvest, error) {
	o, binder, err := opts.Get(settings...)

	return &optionHarvest{
		o:      o,
		binder: binder,
	}, err
}

func (x *primeExtent) complete() bool {
	return true
}

type resumeExtent struct {
	baseExtent
	relic  *pref.Relic
	loaded *opts.LoadInfo
	pin    *resume.Plugin
}

func (x *resumeExtent) plugin(artefacts *kernel.Artefacts) enclave.Plugin {
	x.pin = resume.New(&resume.From{
		Active:   x.loaded.State,
		Mediator: artefacts.Mediator,
		Strategy: x.relic.Strategy,
		IfResult: artefacts.IfResult,
	})

	return x.pin
}

func (x *resumeExtent) options(settings ...pref.Option) (enclave.OptionHarvest, error) {
	loaded, binder, err := resume.Load(&enclave.RestoreState{
		Path:   x.relic.From,
		FS:     x.trees.R,
		Resume: x.relic.Strategy,
	}, settings...)

	x.loaded = loaded

	if x.relic.Restorer != nil {
		err = x.relic.Restorer(x.loaded.O, x.loaded.State)
	}

	return &optionHarvest{
		o:      loaded.O,
		binder: binder,
		loaded: loaded,
	}, err
}

func (x *resumeExtent) complete() bool {
	return x.pin.IfResult.IsComplete()
}
