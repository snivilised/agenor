package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/feat/resume"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type extent interface {
	using() *pref.Using
	was() *pref.Was
	plugin(*kernel.Artefacts) types.Plugin
	options(...pref.Option) (*pref.Options, *opts.Binder, error)
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

func (ex *primeExtent) plugin(*kernel.Artefacts) types.Plugin {
	return nil
}

func (ex *primeExtent) options(settings ...pref.Option) (*pref.Options, *opts.Binder, error) {
	return opts.Get(settings...)
}

func (ex *primeExtent) complete() bool {
	return true
}

type resumeExtent struct {
	baseExtent
	w      *pref.Was
	loaded *opts.LoadInfo
	rp     *resume.Plugin
}

func (ex *resumeExtent) using() *pref.Using {
	return &ex.w.Using
}

func (ex *resumeExtent) was() *pref.Was {
	return ex.w
}

func (ex *resumeExtent) plugin(artefacts *kernel.Artefacts) types.Plugin {
	ex.rp = &resume.Plugin{
		BasePlugin: kernel.BasePlugin{
			Mediator: artefacts.Mediator,
		},
		IfResult: artefacts.IfResult,
	}

	return ex.rp
}

func (ex *resumeExtent) options(settings ...pref.Option) (*pref.Options, *opts.Binder, error) {
	loaded, binder, err := resume.Load(&types.RestoreState{
		Path:   ex.w.From,
		FS:     ex.trees.R,
		Resume: ex.w.Strategy,
	}, settings...)

	ex.loaded = loaded

	// TODO: get the resume point from the resume persistence file
	// then set up hibernation with this defined as a hibernation
	// filter.
	//
	return loaded.O, binder, err
}

func (ex *resumeExtent) complete() bool {
	return ex.rp.IfResult.IsComplete()
}
