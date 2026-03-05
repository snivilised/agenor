package age

import (
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/pref"
)

type optionHarvest struct {
	o      *pref.Options
	binder *opts.Binder
	loaded *opts.LoadInfo
}

// Options returns the final options after applying any necessary adjustments
// based on the loaded configuration. It ensures that certain concurrency
// settings are set appropriately, especially when cascading behaviour is
// involved. The method takes into account the loaded configuration if available,
// otherwise it uses the initially provided options.
func (a *optionHarvest) Options() *pref.Options {
	return a.afterwards(lo.TernaryF(a.loaded != nil,
		func() *pref.Options {
			return a.loaded.O
		},
		func() *pref.Options {
			return a.o
		},
	))
}

// Binder returns the binder associated with the options. The binder is used to
// manage the relationships between options and their sources, allowing for
// dynamic updates and adjustments based on the configuration loading process.
func (a *optionHarvest) Binder() *opts.Binder {
	return a.binder
}

// Loaded returns the LoadInfo containing details about the loaded configuration,
// including the final options and any relevant metadata. This information can be
// used to understand how the options were derived and to make informed decisions
// about further adjustments or actions based on the loaded configuration.
func (a *optionHarvest) Loaded() *opts.LoadInfo {
	return a.loaded
}

func (a *optionHarvest) afterwards(o *pref.Options) *pref.Options {
	if o.Behaviours.Cascade.NoRecurse {
		o.Behaviours.Cascade.Depth = 1
	}

	if o.Concurrency.Input.Size == 0 {
		o.Concurrency.Input.Size = o.Concurrency.NoW
	}

	if o.Concurrency.Output.On != nil {
		if o.Concurrency.Output.Size == 0 {
			o.Concurrency.Output.Size = o.Concurrency.NoW
		}
	}

	return o
}

type pluginsBuilder interface {
	build(o *pref.Options,
		ext extent,
		artefacts *kernel.Artefacts,
		others ...enclave.Plugin,
	) ([]enclave.Plugin, error)
}

type activated func(*pref.Options,
	extent,
	*kernel.Artefacts,
	...enclave.Plugin,
) ([]enclave.Plugin, error)

func (fn activated) build(o *pref.Options,
	ext extent,
	artefacts *kernel.Artefacts,
	others ...enclave.Plugin,
) ([]enclave.Plugin, error) {
	return fn(o, ext, artefacts, others...)
}

type scaffoldBuilder interface {
	build(addons ...Addon) (scaffold, error)
}

type scaffolding func(addons ...Addon) (scaffold, error)

// build creates a scaffold using the provided addons for configuration.
// The addons parameter allows for customisation of the scaffolding process.
func (fn scaffolding) build(addons ...Addon) (scaffold, error) {
	return fn(addons...)
}
