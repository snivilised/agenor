package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type navigatorBuilder interface {
	build(o *pref.Options) (core.Navigator, error)
}

type factory func(o *pref.Options) (core.Navigator, error)

func (fn factory) build(o *pref.Options) (core.Navigator, error) {
	return fn(o)
}

type optionsBuilder interface {
	build() (*pref.Options, error)
}

type optionals func() (*pref.Options, error)

func (fn optionals) build() (*pref.Options, error) {
	return fn()
}

type pluginsBuilder interface {
	build(*pref.Options) ([]types.Plugin, error)
}

type features func(*pref.Options) ([]types.Plugin, error)

func (fn features) build(o *pref.Options) ([]types.Plugin, error) {
	return fn(o)
}

// TODO: do we pass in another func to the directory that represents the sync?
type director func(bs *Builders) core.Navigator

func (fn director) Extent(bs *Builders) core.Navigator {
	return fn(bs)
}

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
// navigationImpl: navigatorBase // implFiles, implFolders, implUniversal // agent

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
