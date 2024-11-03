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
	using() *pref.Using
	was() *pref.Was
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
}

func (ex *baseExtent) forest() *core.Forest {
	return ex.trees
}

type primeExtent struct {
	baseExtent
	u *pref.Using
}

func (ex *primeExtent) using() *pref.Using {
	return ex.u
}

func (ex *primeExtent) was() *pref.Was {
	return nil
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
	w      *pref.Was
	loaded *opts.LoadInfo
	pin    *resume.Plugin
}

func (ex *resumeExtent) using() *pref.Using {
	return &ex.w.Using
}

func (ex *resumeExtent) was() *pref.Was {
	return ex.w
}

func (ex *resumeExtent) plugin(artefacts *kernel.Artefacts) enclave.Plugin {
	ex.pin = resume.New(&resume.From{
		Active:   ex.loaded.State,
		Mediator: artefacts.Mediator,
		Strategy: ex.w.Strategy,
		IfResult: artefacts.IfResult,
	})

	return ex.pin
}

func (ex *resumeExtent) options(settings ...pref.Option) (enclave.OptionHarvest, error) {
	loaded, binder, err := resume.Load(&enclave.RestoreState{
		Path:   ex.w.From,
		FS:     ex.trees.R,
		Resume: ex.w.Strategy,
	}, settings...)

	ex.loaded = loaded

	if ex.w.Restorer != nil {
		err = ex.w.Restorer(ex.loaded.O, ex.loaded.State)
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
