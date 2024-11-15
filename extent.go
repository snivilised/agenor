package age

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/feat/resume"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/agenor/tfs"
)

type (
	extent interface {
		facade() pref.Facade
		subscription() enums.Subscription
		plugin(*kernel.Artefacts) enclave.Plugin
		options([]Addon, ...pref.Option) (enclave.OptionHarvest, error)
		forest() *core.Forest
		complete() bool
	}
)

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

func (x *primeExtent) subscription() enums.Subscription {
	return x.using.Subscription
}

func (x *primeExtent) plugin(*kernel.Artefacts) enclave.Plugin {
	return nil
}

func (x *primeExtent) options(_ []Addon, settings ...pref.Option) (enclave.OptionHarvest, error) {
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

func (x *resumeExtent) subscription() enums.Subscription {
	return x.loaded.State.Subscription
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

func (x *resumeExtent) options(addons []Addon,
	settings ...pref.Option,
) (enclave.OptionHarvest, error) {
	loaded, binder, err := resume.Load(&enclave.RestoreState{
		Path:   x.relic.From,
		FS:     x.trees.R,
		Resume: x.relic.Strategy,
	}, settings...)

	x.loaded = loaded

	if handler := x.seek(addons); handler != nil {
		handler.OnLoad(loaded.State)
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

func (x *resumeExtent) seek(addons []Addon) enclave.StateHandler {
	if addon, found := lo.Find(addons, func(item Addon) bool {
		_, ok := item.(enclave.StateHandler)
		return ok
	}); found {
		result, _ := addon.(enclave.StateHandler)
		return result
	}

	return nil
}
