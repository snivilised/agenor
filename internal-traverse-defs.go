package age

import (
	"io/fs"

	"github.com/snivilised/agenor/core"
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

func (a *optionHarvest) Binder() *opts.Binder {
	return a.binder
}

func (a *optionHarvest) Loaded() *opts.LoadInfo {
	return a.loaded
}

func (a *optionHarvest) afterwards(o *pref.Options) *pref.Options {
	if o.Behaviours.Cascade.NoRecurse {
		o.Behaviours.Cascade.Depth = 1
	}

	return o
}

// optionsBuilder
type optionsBuilder interface {
	build(ext extent) (enclave.OptionHarvest, error)
}

type optionBuilder func(ext extent) (enclave.OptionHarvest, error)

func (fn optionBuilder) build(ext extent) (enclave.OptionHarvest, error) {
	return fn(ext)
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

type fsBuilder interface {
	build(path string) fs.FS
}

type filesystem func(path string) fs.FS

func (fn filesystem) build(path string) fs.FS {
	return fn(path)
}

type extentBuilder interface {
	build(forest *core.Forest) extent
}

type extension func(forest *core.Forest) extent

func (fn extension) build(forest *core.Forest) extent {
	return fn(forest)
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
