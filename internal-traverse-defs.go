package tv

import (
	"io/fs"

	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

// optionsBuilder
type optionsBuilder interface {
	build(ext extent) (*pref.Options, error)
}

type optionals func(ext extent) (*pref.Options, error)

func (fn optionals) build(ext extent) (*pref.Options, error) {
	return fn(ext)
}

// pluginsBuilder
type pluginsBuilder interface {
	build(*pref.Options, types.Mediator, ...types.Plugin) ([]types.Plugin, error)
}

type features func(*pref.Options, types.Mediator, ...types.Plugin) ([]types.Plugin, error)

func (fn features) build(o *pref.Options, mediator types.Mediator, others ...types.Plugin) ([]types.Plugin, error) {
	return fn(o, mediator, others...)
}

type fsBuilder interface {
	build(path string) fs.FS
}

type filesystem func(path string) fs.FS

func (fn filesystem) build(path string) fs.FS {
	return fn(path)
}

type extentBuilder interface {
	build(fsys fs.FS) extent
}

type extension func(fs.FS) extent

func (fn extension) build(fsys fs.FS) extent {
	return fn(fsys)
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
