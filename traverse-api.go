package tv

import (
	"time"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/pref"
)

// Director
type Director interface {
	// Extent represents the magnitude of the traversal; ie we can
	// perform a full Prime run, or Resume from a previously
	// cancelled run.
	//
	// I wonder if this could be replaced by Prime/Resume, to simplify
	//
	Extent(bs *Builders) core.Navigator
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

const (
	OutputChSize       = 10
	CheckCloseInterval = time.Second / 10
	TimeoutOnSend      = time.Second * 2
)

// traverse is the front line user facing interface to this module. It sits
// on the top of the code stack and is allowed to use anything, but nothing
// else can depend on definitions here, except unit tests.

type Client = core.Client
type Node = core.Node
type Option = pref.Option
type Options = pref.Options
type Subscription = enums.Subscription

const (
	SubscribeFiles            = enums.SubscribeFiles
	SubscribeFolders          = enums.SubscribeFolders
	SubscribeFoldersWithFiles = enums.SubscribeFoldersWithFiles
	SubscribeUniversal        = enums.SubscribeUniversal
)

type ResumeStrategy = enums.ResumeStrategy

const (
	ResumeStrategySpawn    = enums.ResumeStrategySpawn
	ResumeStrategyFastward = enums.ResumeStrategyFastward
)

type Was = pref.Was

var (
	WithCPU                  = pref.WithCPU
	WithDepth                = pref.WithDepth
	WithFilter               = pref.WithFilter
	WithHibernationWake      = pref.WithHibernationWake
	WithHibernationSleep     = pref.WithHibernationSleep
	WithHibernationBehaviour = pref.WithHibernationBehaviour
	WithHookQueryStatus      = pref.WithHookQueryStatus
	WithHookReadDirectory    = pref.WithHookReadDirectory
	WithHookSort             = pref.WithHookSort
	WithHookFileSubPath      = pref.WithHookFileSubPath
	WithHookFolderSubPath    = pref.WithHookFolderSubPath
	WithNavigationBehaviours = pref.WithNavigationBehaviours
	WithOnAscend             = pref.WithOnAscend
	WithOnBegin              = pref.WithOnBegin
	WithOnDescend            = pref.WithOnDescend
	WithOnEnd                = pref.WithOnEnd
	WithOnStart              = pref.WithOnStart
	WithOnStop               = pref.WithOnStop
	WithNoRecurse            = pref.WithNoRecurse
	WithNoW                  = pref.WithNoW
	WithSampler              = pref.WithSampler
	WithSampling             = pref.WithSampling
	WithSamplingInReverse    = pref.WithSamplingInReverse
	WithSamplingNoOf         = pref.WithSamplingNoOf
	WithSamplingOptions      = pref.WithSamplingOptions
	WithSamplingType         = pref.WithSamplingType
	WithSortBehaviour        = pref.WithSortBehaviour
	WithSubPathBehaviour     = pref.WithSubPathBehaviour
)

type Using = pref.Using

// sub package description:
//

// This high level list assumes everything can use core and enums; dependencies
// can only point downwards. NB: These restrictions do not apply to the unit tests;
// eg, "cycle_test" defines tests that are dependent on "pref", but "cycle" is prohibited
// from using "cycle".
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
// ---
//
// 🔆 support layer
// tapable: [core, pref]
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
// core: []
// enums: [none]
// ---
// ============================================================================
//
