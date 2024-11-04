package tv

import (
	"time"

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/filtering"
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
		Extent(bs *Builders) core.Navigator
	}

	director func(bs *Builders) core.Navigator
)

func (fn director) Extent(bs *Builders) core.Navigator {
	return fn(bs)
}

// NavigatorFactory
type NavigatorFactory interface {
	// ConfigureL is a factory function that creates a navigator.
	// We don't return an error here as that would make using the factory
	// awkward. Instead, if there is an error during the build process,
	// we return a fake navigator that when invoked immediately returns
	// a traverse error indicating the build issue.
	//
	Configure() Director
}

type (
	// 🌀 core
	Client  = core.Client
	Node    = core.Node
	Servant = core.Servant

	// 🌀 enums
	Subscription   = enums.Subscription
	ResumeStrategy = enums.ResumeStrategy

	// 🌀 nef
	ExistsInFS  = nef.ExistsInFS
	Rel         = nef.Rel
	RenameFS    = nef.RenameFS
	TraverseFS  = nef.TraverseFS
	WriteFileFS = nef.WriteFileFS
	WriterFS    = nef.WriterFS

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
	SubscribeFiles                = enums.SubscribeFiles
	SubscribeDirectories          = enums.SubscribeDirectories
	SubscribeDirectoriesWithFiles = enums.SubscribeDirectoriesWithFiles
	SubscribeUniversal            = enums.SubscribeUniversal
)

var (
	// 🌀 nef

	// NewReadDirFS creates a file system with read directory capability
	NewReadDirFS = nef.NewReadDirFS

	// NewReaderFS creates a read only file system
	NewReaderFS = nef.NewReaderFS

	// NewReadFileFS  creates a file system with read file capability
	NewReadFileFS = nef.NewReadFileFS

	// NewStatFS creates a file system with Stat method
	NewStatFS = nef.NewStatFS

	// NewTraverseFS creates a file system that implements functionality required
	// by the traverse
	NewTraverseFS = nef.NewTraverseFS

	// NewUniversalFS creates a file system that's got it all
	NewUniversalFS = nef.NewUniversalFS

	// NewWriteFileFS creates a file system with write file capability
	NewWriteFileFS = nef.NewWriteFileFS

	// NewWriterFS creates a file system with writer capabilities
	NewWriterFS = nef.NewWriterFS

	// 🌀 filtering

	// NewCustomSampleFilter only needs to be called explicitly when defining
	// a custom sample filter.
	NewCustomSampleFilter = filtering.NewCustomSampleFilter

	// 🌀 pref

	// IfOption enables options to be conditional. IfOption condition evaluates to true
	// then the option is returned, otherwise nil.
	IfOption = pref.IfOption

	// IfOptionF allows the delaying of creation of the option until the condition
	// is known to be true. This is in contrast to IfOption where the Option is
	// pre-created, regardless of the condition.
	IfOptionF = pref.IfOptionF

	// WithCPU configures the worker pool used for concurrent traversal sessions
	// in the Run function to utilise a number of go-routines equal to the available
	// CPU count, optimising performance based on the system's processing capabilities.
	WithCPU = pref.WithCPU

	// WithDepth sets the maximum number of directories deep the navigator
	// will traverse to.
	WithDepth = pref.WithDepth

	// WithFaultHandler defines a custom handler to handle an error that occurs
	// when 'Stat'ing the tree root directory. When an error occurs, traversal terminates
	// immediately. The handler specified allows custom functionality to be invoked
	// when an error occurs here.
	WithFaultHandler = pref.WithFaultHandler

	// WithFilter used to determine which file system nodes (files or directories)
	// the client defined handler is invoked for. Note that the filter does not
	// determine navigation, it only determines wether the callback is invoked.
	WithFilter = pref.WithFilter

	// WithHibernationBehaviourExclusiveWake activates hibernation
	// with a wake condition. The wake condition should be defined
	// using WithHibernationFilterWake.
	WithHibernationBehaviourExclusiveWake = pref.WithHibernationBehaviourExclusiveWake

	// WithHibernationBehaviourInclusiveSleep activates hibernation
	// with a sleep condition. The sleep condition should be defined
	// using WithHibernationFilterSleep.
	WithHibernationBehaviourInclusiveSleep = pref.WithHibernationBehaviourInclusiveSleep

	// WithHibernationFilterSleep defines the sleep condition
	// for hibernation based traversal sessions.
	WithHibernationFilterSleep = pref.WithHibernationFilterSleep

	// WithHibernationFilterWake defines the wake condition
	// for hibernation based traversal sessions.
	WithHibernationFilterWake = pref.WithHibernationFilterWake

	// WithHibernationOptions defines options for a hibernation traversal
	// session.
	WithHibernationOptions = pref.WithHibernationOptions

	// WithHookCaseSensitiveSort specifies that a directory's contents
	// should be sorted with case sensitivity.
	WithHookCaseSensitiveSort = pref.WithHookCaseSensitiveSort

	// WithHookDirectorySubPath defines an custom hook to override the
	// default behaviour for obtaining the sub-path of a directory.
	WithHookDirectorySubPath = pref.WithHookDirectorySubPath

	// WithHookFileSubPath defines an custom hook to override the
	// default behaviour for obtaining the sub-path of a file.
	WithHookFileSubPath = pref.WithHookFileSubPath

	// WithHookQueryStatus defines an custom hook to override the
	// default behaviour for Stating a directory.
	WithHookQueryStatus = pref.WithHookQueryStatus

	// WithHookReadDirectory defines an custom hook to override the
	// default behaviour for reading a directory's contents.
	WithHookReadDirectory = pref.WithHookReadDirectory

	// WithHookSort defines an custom hook to override the
	// default behaviour for sorting a directory's contents.
	WithHookSort = pref.WithHookSort

	// WithLogger defines a structure logger
	WithLogger = pref.WithLogger

	// WithNavigationBehaviours defines all navigation behaviours
	WithNavigationBehaviours = pref.WithNavigationBehaviours

	// WithOnAscend sets ascend handler, invoked when navigator
	// traverses up a directory, ie after all children have been
	// visited.
	WithOnAscend = pref.WithOnAscend

	// WithOnBegin sets the begin handler, invoked before the start
	// of a traversal session.
	WithOnBegin = pref.WithOnBegin

	// WithOnDescend sets the descend handler, invoked when navigator
	// traverses down into a child directory.
	WithOnDescend = pref.WithOnDescend

	// WithOnEnd sets the end handler, invoked at the end of a traversal
	// session.
	WithOnEnd = pref.WithOnEnd

	// WithOnSleep sets the sleep handler, when hibernation is active
	// and the sleep condition has occurred, ie when a file system
	// node is encountered that matches the hibernation's sleep filter.
	WithOnSleep = pref.WithOnSleep

	// WithOnWake sets the wake handler, when hibernation is active
	// and the wake condition has occurred, ie when a file system
	// node is encountered that matches the hibernation's wake filter.
	WithOnWake = pref.WithOnWake

	// WithPanicHandler defines a custom handler to handle a panic.
	WithPanicHandler = pref.WithPanicHandler

	// WithNoRecurse sets the navigator to not descend sub-directories.
	WithNoRecurse = pref.WithNoRecurse

	// WithNoW sets the number of go-routines to use in the worker
	// pool used for concurrent traversal sessions requested by using
	// the Run function.
	WithNoW = pref.WithNoW

	// WithSamplingOptions specifies the sampling options.
	// SampleType: the type of sampling to use
	// SampleInReverse: determines the direction of iteration for the sampling
	// operation
	// NoOf: specifies number of items required in each sample (only applies
	// when not using Custom iterator options)
	// Iteration: allows the client to customise how a directory's contents are sampled.
	// The default way to sample is either by slicing the directory's contents or
	// by using the filter to select either the first/last n entries (using the
	// SamplingOptions). If the client requires an alternative way of creating a
	// sample, eg to take all files greater than a certain size, then this
	// can be achieved by specifying Each and While inside Iteration.
	WithSamplingOptions = pref.WithSamplingOptions

	// WithSkipHandler defines a handler that will be invoked if the
	// client callback returns an error during traversal. The client
	// can control if traversal is either terminated early (fs.SkipAll)
	// or the remaining items in a directory are skipped (fs.SkipDir).
	WithSkipHandler = pref.WithSkipHandler

	// WithSortBehaviour enabling setting of all sorting behaviours.
	WithSortBehaviour = pref.WithSortBehaviour

	// WithSubPathBehaviour defines all sub-path behaviours.
	WithSubPathBehaviour = pref.WithSubPathBehaviour
)

// sub package description:
//

// This high level list assumes everything can use core and enums; dependencies
// can only point downwards. NB: These restrictions do not apply to the unit tests;
// eg, "life_test" defines tests that are dependent on "pref", but "life" is prohibited
// from using "pref".
// ============================================================================
// 🔆 user interface layer
// traverse: [everything]
// ---
//
// 🔆 feature layer
// resume: ["pref", "opts", "kernel"]
// sampling: ["filter"]
// hiber: ["filter", "services"]
// filter: []
//
// 🔆 central layer
// kernel: []
// enclave: [pref, override]
// opts: [pref]
// override: [tapable], !("enclave")
// ---
//
// 🔆 support layer
// pref: ["life", "services", "persist(to-be-confirmed)"] actually, persist should be part of pref
// persist: []
// services: []
// ---
//
// 🔆 intermediary layer
// life: [], !("pref")
// ---
//
// 🔆 platform layer
// tapable: [core]
// core: []
// enums: [none]
// lfs:
// ---
// ============================================================================
//
