package tv

import (
	"time"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/feat/refine"
	"github.com/snivilised/traverse/nfs"
	"github.com/snivilised/traverse/pref"
)

// 📦 pkg: traverse - is the front line user facing interface to this module. It sits
// on the top of the code stack and is allowed to use anything, but nothing
// else can depend on definitions here, except unit tests.

type (
	// Director
	Director interface {
		// Extent represents the magnitude of the traversal; ie we can
		// perform a full Prime run, or Resume from a previously
		// cancelled run.
		//
		// I wonder if this could be replaced by Prime/Resume, to simplify
		//
		Extent(bs *Builders) core.Navigator
	}

	director func(bs *Builders) core.Navigator
)

func (fn director) Extent(bs *Builders) core.Navigator {
	return fn(bs)
}

// NavigatorFactory
type NavigatorFactory interface {
	// Configure is a factory function that creates a navigator.
	// We don't return an error here as that would make using the factory
	// awkward. Instead, if there is an error during the build process,
	// we return a fake navigator that when invoked immediately returns
	// a traverse error indicating the build issue.
	//
	Configure() Director
}

type (
	// 🌀 core
	Node   = core.Node
	Client = core.Client

	// 🌀 enums
	Subscription   = enums.Subscription
	ResumeStrategy = enums.ResumeStrategy

	// 🌀 refine
	BaseSampleFilter = refine.SampleFilter

	// 🌀 nfs
	FileSystems = nfs.FileSystems

	// 🌀 pref
	Option  = pref.Option
	Options = pref.Options
	Using   = pref.Using
	Was     = pref.Was
)

const (
	OutputChSize       = 10
	CheckCloseInterval = time.Second / 10
	TimeoutOnSend      = time.Second * 2

	// 🌀 enum: ResumeStrategy
	ResumeStrategySpawn    = enums.ResumeStrategySpawn
	ResumeStrategyFastward = enums.ResumeStrategyFastward

	// 🌀 enum:Subscribe
	SubscribeFiles            = enums.SubscribeFiles
	SubscribeFolders          = enums.SubscribeFolders
	SubscribeFoldersWithFiles = enums.SubscribeFoldersWithFiles
	SubscribeUniversal        = enums.SubscribeUniversal
)

var (
	// 🌀 nfs
	NewNativeFS      = nfs.NewNativeFS
	NewQueryStatusFS = nfs.NewQueryStatusFS

	// 🌀 refine
	NewSampleFilter = refine.NewSampleFilter

	// 🌀 pref
	IfOption                               = pref.IfOption
	IfOptionF                              = pref.IfOptionF
	WithCPU                                = pref.WithCPU
	WithDepth                              = pref.WithDepth
	WithFaultHandler                       = pref.WithFaultHandler
	WithFilter                             = pref.WithFilter
	WithHibernationBehaviourExclusiveWake  = pref.WithHibernationBehaviourExclusiveWake
	WithHibernationBehaviourInclusiveSleep = pref.WithHibernationBehaviourInclusiveSleep
	WithHibernationFilterSleep             = pref.WithHibernationFilterSleep
	WithHibernationFilterWake              = pref.WithHibernationFilterWake
	WithHibernationOptions                 = pref.WithHibernationOptions
	WithHibernationBehaviour               = pref.WithHibernationOptions
	WithHookQueryStatus                    = pref.WithHookQueryStatus
	WithHookReadDirectory                  = pref.WithHookReadDirectory
	WithHookSort                           = pref.WithHookSort
	WithHookCaseSensitiveSort              = pref.WithHookCaseSensitiveSort
	WithHookFileSubPath                    = pref.WithHookFileSubPath
	WithHookFolderSubPath                  = pref.WithHookFolderSubPath
	WithNavigationBehaviours               = pref.WithNavigationBehaviours
	WithOnAscend                           = pref.WithOnAscend
	WithOnBegin                            = pref.WithOnBegin
	WithOnDescend                          = pref.WithOnDescend
	WithOnEnd                              = pref.WithOnEnd
	WithOnStart                            = pref.WithOnStart
	WithOnStop                             = pref.WithOnStop
	WithPanicHandler                       = pref.WithPanicHandler
	WithNoRecurse                          = pref.WithNoRecurse
	WithNoW                                = pref.WithNoW
	WithSampling                           = pref.WithSampling
	WithSkipHandler                        = pref.WithSkipHandler
	WithSortBehaviour                      = pref.WithSortBehaviour
	WithSubPathBehaviour                   = pref.WithSubPathBehaviour
)

// sub package description:
//

// This high level list assumes everything can use core and enums; dependencies
// can only point downwards. NB: These restrictions do not apply to the unit tests;
// eg, "cycle_test" defines tests that are dependent on "pref", but "cycle" is prohibited
// from using "pref".
// ============================================================================
// 🔆 user interface layer
// traverse: [everything]
// ---
//
// 🔆 feature layer
// resume: ["pref"]
// sampling: ["refine"]
// hiber: ["refine", "services"]
// refine: []
//
// 🔆 central layer
// kernel: []
// types: [measure, pref, override]
// override: [tapable], !("types")
// ---
//
// 🔆 support layer
// pref: ["cycle", "services", "persist(to-be-confirmed)"] actually, persist should be part of pref
// persist: []
// services: []
// ---
//
// 🔆 intermediary layer
// cycle: [], !("pref")
// ---
//
// 🔆 platform layer
// tapable: [core]
// core: []
// enums: [none]
// measure: []
// nfs:
// ---
// ============================================================================
//
