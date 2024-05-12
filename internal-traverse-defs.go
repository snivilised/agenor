package traverse

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/pref"
)

// extentFunc
type extentFunc func(ob OptionsBuilder) *pref.Options

type OptionsBuilder interface {
	get() *pref.Options
}

type optionals func() *pref.Options

func (fn optionals) get() *pref.Options {
	return fn()
}

// TODO: do we pass in another func to the directory that represents the sync?
type direct func(ob OptionsBuilder) core.Navigator

func (fn direct) Extent(ob OptionsBuilder) core.Navigator {
	return fn(ob)
}

type syncBuilder interface {
	wake(at string) error // we might need to pass in options
}

type sync func(at string) error

func (fn sync) wake(at string) error {
	return fn(at)
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
