package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/feat/resume"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/pref"
)

type extent interface {
	facade() pref.Facade
	plugin(*kernel.Artefacts) enclave.Plugin
	options(...pref.Option) (enclave.OptionHarvest, error)
	forest() *core.Forest
	complete() bool
}

type fileSystems struct {
	fS TraverseFS
}

type baseExtent struct {
	trees *core.Forest
	fac   pref.Facade
}

func (ex *baseExtent) forest() *core.Forest {
	return ex.trees
}

func (ex *baseExtent) facade() pref.Facade {
	return ex.fac
}

type primeExtent struct {
	baseExtent
	using *pref.Using
}

func (ex *primeExtent) plugin(*kernel.Artefacts) enclave.Plugin {
	return nil
}

func (ex *primeExtent) options(settings ...pref.Option) (enclave.OptionHarvest, error) {
	o, binder, err := opts.Get(settings...)

	return &optionHarvest{
		o:      o,
		binder: binder,
	}, err
}

func (ex *primeExtent) complete() bool {
	return true
}

type resumeExtent struct {
	baseExtent
	relic  *pref.Relic
	loaded *opts.LoadInfo
	pin    *resume.Plugin
}

func (ex *resumeExtent) plugin(artefacts *kernel.Artefacts) enclave.Plugin {
	ex.pin = resume.New(&resume.From{
		Active:   ex.loaded.State,
		Mediator: artefacts.Mediator,
		Strategy: ex.relic.Strategy,
		IfResult: artefacts.IfResult,
	})

	return ex.pin
}

func (ex *resumeExtent) options(settings ...pref.Option) (enclave.OptionHarvest, error) {
	loaded, binder, err := resume.Load(&enclave.RestoreState{
		Path:   ex.relic.From,
		FS:     ex.trees.R,
		Resume: ex.relic.Strategy,
	}, settings...)

	ex.loaded = loaded

	if ex.relic.Restorer != nil {
		err = ex.relic.Restorer(ex.loaded.O, ex.loaded.State)
	}

	return &optionHarvest{
		o:      loaded.O,
		binder: binder,
		loaded: loaded,
	}, err
}

func (ex *resumeExtent) complete() bool {
	return ex.pin.IfResult.IsComplete()
}
