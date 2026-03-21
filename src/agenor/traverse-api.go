// Package agenor is the front line user facing interface to this module.
// It sits on the top of the code stack and is allowed to use anything, but
// nothing else can depend on definitions here, except unit tests.
package agenor

import (
	"time"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/internal/filtering"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/agenor/tfs"
	nef "github.com/snivilised/nefilim"
)

type (
	// Director is an interface that represents a director for creating Navigator instances. It
	// defines a single method, Extent, which takes a Builders instance and returns a Navigator.
	// The Director is responsible for orchestrating the construction of a Navigator based on
	// the provided Builders and any configuration specified in the Addons.
	Director interface {
		// Extent represents the magnitude of the traversal; ie we can
		// perform a full 'Prime' traversal, or 'Resume' from a previously
		// cancelled run.
		//
		Extent(bs *Builders) Navigator
	}

	director func(bs *Builders) Navigator
)

// Extent is a method of the Director interface that represents the magnitude of the traversal.
func (fn director) Extent(bs *Builders) Navigator {
	return fn(bs)
}

// NavigatorFactory is a factory interface for creating Navigator instances. It defines a
// single method, Configure, which takes a variable number of Addon arguments and returns
// a Director. The Configure method is responsible for setting up the necessary configuration
// and dependencies based on the provided Addons and returning a Director that can be used
// to create Navigator instances with the specified configuration.
type NavigatorFactory interface {
	// Configure is a factory function that creates a navigator.
	// We don't return an error here as that would make using the factory
	// awkward. Instead, if there is an error during the build process,
	// we return a fake navigator that when invoked immediately returns
	// a traverse error indicating the build issue.
	//
	Configure(addons ...Addon) Director
}

type (
	// 🌀 core

	// Client is the callback invoked for each file system node found
	// during traversal.
	Client = core.Client

	// Navigator represents the core navigation interface. It is the main interface
	// that users interact with to perform file system traversal. The Navigator interface
	// defines a single method, Navigate, which takes a context and returns a TraverseResult
	// and an error. The Navigate method is responsible for executing the traversal logic
	// based on the provided context and returning the results of the traversal.
	Navigator = core.Navigator

	// Node represents a file system node, which can be either a file or a directory. It
	// contains information about the node's path, type, and any extensions that may be associated with
	// it. The Node struct is used to represent the individual elements of the file system that are
	// encountered during traversal.
	Node = core.Node

	// Servant provides the client with facility to request properties
	Servant = core.Servant

	// 🌀 enums

	// Subscription represents the types of file system nodes that can be subscribed to
	// during traversal. It is used to specify whether the client wants to receive callbacks
	// for files, directories, or both.
	Subscription = enums.Subscription

	// ResumeStrategy represents the strategies for resuming a traversal session. It is used to
	// specify how the traversal should continue from a previously cancelled session, such as
	// whether to spawn new sessions or fast-forward to the last known state.
	ResumeStrategy = enums.ResumeStrategy

	// 🌀 nef

	// ExistsInFS contains methods that check the existence of file system items.
	ExistsInFS = nef.ExistsInFS

	// Rel represents generic info required to create a relative file system.
	// Relative just means that a file system is created with a root path and
	// the operations on the file system are invoked with paths that must be
	// relative to the root.
	Rel = nef.Rel

	// RenameFS contains methods for renaming files and directories in a file system.
	RenameFS = nef.RenameFS

	// WriteFileFS contains methods for writing files in a file system.
	WriteFileFS = nef.WriteFileFS

	// WriterFS contains methods for writing files and directories in a file system.
	WriterFS = nef.WriterFS

	// 🌀 pref

	// Accepter is the function signature for functions that can be accepted as options
	// in the Addon interface.
	Accepter = pref.Accepter

	// Head represents the initial configuration of a traversal session. It is used
	// to define the starting point and any initial settings for the traversal.
	Head = pref.Head

	// Option represents a configuration option that can be applied to a traversal
	// session. It is used to specify various settings and behaviors for the traversal,
	// such as filters, handlers, and other preferences.
	Option = pref.Option

	// Options represents a collection of configuration options that can be applied
	// to a traversal session. It is used to group multiple options together for easier
	// management and application to the traversal.
	Options = pref.Options

	// Relic represents a saved state of a traversal session that can be used to resume.
	Relic = pref.Relic

	// Using represents the dependencies required by an Addon to be applied to a
	// traversal session.
	Using = pref.Using

	// 🌀 tfs

	// TraversalFS represents the file system interface used for traversal. It defines
	// the methods that must be implemented by any file system that is to be traversed
	// using the Navigator. The TraversalFS interface includes methods for reading
	// directories, reading files, and obtaining file information, among others. It
	// serves as the abstraction layer between the traversal logic and the underlying
	// file system, allowing for flexibility in the types of file systems that can be traversed.
	TraversalFS = tfs.TraversalFS
)

const (
	// OutputChSize defines the size of the output channel used when WithOutput is specified.
	OutputChSize = 10

	// CheckCloseInterval defines the interval at which the output channel is checked for closure.
	CheckCloseInterval = time.Second / 10

	// TimeoutOnSend defines the duration to wait when sending output before timing out.
	TimeoutOnSend = time.Second * 2

	// 🌀 enum: ResumeStrategy

	// ResumeStrategySpawn indicates that when resuming a traversal session, new sessions
	// should be spawned to continue the traversal from the last known state.
	ResumeStrategySpawn = enums.ResumeStrategySpawn

	// ResumeStrategyFastward indicates that when resuming a traversal session, the traversal
	// should fast-forward to the last known state and continue from there, rather than
	// spawning new sessions.
	ResumeStrategyFastward = enums.ResumeStrategyFastward

	// 🌀 enum:Subscribe

	// SubscribeFiles indicates that the client wants to receive callbacks for file nodes
	// during traversal.
	SubscribeFiles = enums.SubscribeFiles

	// SubscribeDirectories indicates that the client wants to receive callbacks for
	// directory nodes during traversal.
	SubscribeDirectories = enums.SubscribeDirectories

	// SubscribeDirectoriesWithFiles indicates that the client wants to receive callbacks for
	// both file and directory nodes during traversal.
	SubscribeDirectoriesWithFiles = enums.SubscribeDirectoriesWithFiles

	// SubscribeUniversal indicates that the client wants to receive callbacks for all file
	// system nodes	during traversal, regardless of whether they are files or directories.
	SubscribeUniversal = enums.SubscribeUniversal
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

	// IfOptionF allows the delaying of inception of the option until the condition
	// is known to be true. This is in contrast to IfOption where the Option is
	// pre-created, regardless of the condition.
	IfOptionF = pref.IfOptionF

	// IfElseOptionF is similar to IfOptionF except that it accepts 2 options, the
	// first represents the returned option if the condition true and the second
	// if false.
	// IfElseOptionF provides conditional option selection similar to IfOptionF but
	// handles both true and false cases. It accepts a condition and two
	// ConditionalOption functions:
	// tOption (executed when condition is true) and
	// fOption (executed when condition is false).
	IfElseOptionF = pref.IfElseOptionF

	// WithAdminPath defines the path for admin related files
	WithAdminPath = pref.WithAdminPath

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

	// WithOutput requests that the worker pool emits outputs
	WithOutput = pref.WithOutput

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
