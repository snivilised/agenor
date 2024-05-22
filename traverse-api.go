package tv

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/pref"
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
	WithCancel               = pref.WithCancel
	WithCPU                  = pref.WithCPU
	WithContext              = pref.WithContext
	WithDepth                = pref.WithDepth
	WithFilter               = pref.WithFilter
	WithHibernation          = pref.WithHibernation
	WithHibernationBehaviour = pref.WithHibernationBehaviour
	WithNavigationBehaviours = pref.WithNavigationBehaviours
	WithNoRecurse            = pref.WithNoRecurse
	WithNoW                  = pref.WithNoW
	WithSamplerOptions       = pref.WithSamplerOptions
	WithSampling             = pref.WithSampling
	WithSamplingInReverse    = pref.WithSamplingInReverse
	WithSamplingNoOf         = pref.WithSamplingNoOf
	WithSamplingOptions      = pref.WithSamplingOptions
	WithSamplingType         = pref.WithSamplingType
	WithSortBehaviour        = pref.WithSortBehaviour
	WithSubPathBehaviour     = pref.WithSubPathBehaviour
	WithSubscription         = pref.WithSubscription
)

type Using = pref.Using

// sub package description:
//

// This high level list assumes everything can use core and enums; dependencies
// can only point downwards.
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
