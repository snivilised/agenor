package tv

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

type optionHarvest struct {
	o      *pref.Options
	binder *opts.Binder
	loaded *opts.LoadInfo
}

func (a *optionHarvest) Options() *pref.Options {
	return lo.TernaryF(a.loaded != nil,
		func() *pref.Options {
			return a.loaded.O
		},
		func() *pref.Options {
			return a.o
		},
	)
}

func (a *optionHarvest) Binder() *opts.Binder {
	return a.binder
}

func (a *optionHarvest) Loaded() *opts.LoadInfo {
	return a.loaded
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

// We need an entity that manages the decoration of the client handler. The
// following scenarios need to be catered for:
//
// [prime]
// - raw: handler not decorated
// - filtered: decorated by filter
// - sample: decorated by filter
// - hiber: decorated by hiber wake/sleep filter
// [resume-fastward]
// - raw: handler decorated by hiber wake/filter >>> fast forwarding
// [resume-spawn]
// - raw: handler not decorated
//
// We need something that manages the client callback, let's call this
// the guardian; it kind of replaces the agent, but has a more limited scope.
// The guardian simply manages decorations for all the scenarios we have listed
// above. And to ensure that it can't be abused, we hide it behind an interface
// and let's call that the sentinel; so the guardian is an implementation of
// the sentinel.
//
// So we have to make sure that we have a persistence object that records
// everything about the internal state

// ✨ =========================================================================
//																					extent
//													primary												resume
// ============================================================================
// sync
//			sequential
//
//			reactive
//

// [options builder.GetOptions]
//																					extent
//													primary												resume
// ============================================================================
// sync
//			sequential					from params										from file
//			reactive						from params										from file
//

// ✨ ==========================================================================
// KERNEL: mediator

// KERNEL: navigationController
// navigationImpl: navigatorBase // implFiles, implDirectories, implUniversal // agent

// FEATURE: resume hydrate (depends on FEATURE: hibernation)
// resumeController
// resumeStrategy

// FEATURE: sampling
// samplingController

// FEATURE: filter

// FEATURE: hibernation

// ⏰ TIMELINE
// !!! keep in mind that the bootstrap must fall away after initialisation.
// Any orchestration required during navigation time is the responsibility
// of the mediator
//
// --> pre init
// * features register handler with message bus; if they can
//
// --> client invoke new
// * create driver/session
//
// --> acquire options
// * get default options
// * primary: user applies choices to the defaults
// * resume: load options from file and apply to defaults
// * announce options.available (might not be required, let's see)
//
// --> configure features (filter,sampler,hibernation,resume)
//
