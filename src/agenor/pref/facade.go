package pref

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/locale"
	"github.com/snivilised/jaywalk/src/agenor/tfs"
)

type (
	// Facade is the main interface that clients interact with to perform file
	// system traversal. It provides methods to access the client callback,
	// the file system builder, validate the configuration, and determine the magnitude
	// of the traversal (e.g., "prime" for primary navigation or "resume" for resumption).
	// The Facade interface abstracts away the underlying implementation details and
	// provides a clean API for clients to use when setting up and executing traversals.
	Facade interface {
		// Client is the callback invoked for each file system node found during traversal.
		Client() core.Client

		// Forest returns the client's file system builder, if provided, otherwise it returns
		// the default file system builder.
		Forest() BuildForest

		// Validate checks that the required properties of the configuration are set and returns
		// an error if not.
		Validate() error

		// Magnitude returns a string that represents the magnitude of the traversal (e.g., "prime" for
		// primary navigation or "resume" for resumption).
		Magnitude() string
	}

	// Head contains information common to both primary and resume
	// navigation sessions.
	Head struct {
		// Handler is the callback function invoked for each encountered
		// file system node.
		Handler core.Client

		// GetForest is optional and enables the client to specify how the
		// file systems for a path is created. Typically used by unit tests,
		// but can be used by client to specify different file systems
		// than the default; eg if the client needs to integrate with a
		// file system like afero, they can do so by providing the required
		// adapters.
		GetForest BuildForest
	}

	// Using contains information required to instigate a primary
	// navigation.
	Using struct {
		Head
		// Subscription indicates which file system nodes the client's
		// callback function will be invoked for.
		Subscription enums.Subscription

		// Tree is the root of the directory tree to be traversed and
		// should not be confused with the Root of the file system when
		// the file system in use is relative.
		Tree string

		// O is the optional Options entity. If provided, then these
		// options will be used verbatim, without requiring WithXXX
		// options setters. This is useful if multiple traversals are
		// required, eg a preview traversal followed by a full
		// traversal; in this case the full traversal can reuse the
		// same options that was used in the preview, by setting this
		// property.
		O *Options
	}

	// Relic contains information required to instigate a resume.
	Relic struct {
		Head
		// From is the path to the resumption file from which a prior
		// traverse session is loaded.
		From string

		// Strategy represent what type of resume is run.
		Strategy enums.ResumeStrategy
	}
)

type (
	// TraverseFileSystemBuilder represents the callback function a client can
	// provide to enable them to specify how the file system for a path is
	// created.	Typically used by unit tests, but can be used by client to specify different
	// file systems than the default; eg if the client needs to integrate
	// with a file system like afero, they can do so by providing the required
	// adapters.
	TraverseFileSystemBuilder interface {
		// Build creates a file system for the specified root path.
		Build(root string) tfs.TraversalFS
	}

	// CreateTraverseFS is a function type that implements TraverseFileSystemBuilder.
	CreateTraverseFS func(root string) tfs.TraversalFS
)

// Build creates a file system for the specified root path.
func (fn CreateTraverseFS) Build(root string) tfs.TraversalFS {
	return fn(root)
}

type (
	// ResumeFileSystemBuilder represents the callback function a client can
	// provide to enable them to specify how the file system for a path is
	// created when resuming.	Typically used by unit tests, but can be used by
	// client to specify different file systems than the default; eg if the client
	// needs to integrate with a file system like afero, they can do so by providing
	// the required	adapters.
	ResumeFileSystemBuilder interface {
		// Build creates a file system for resumption.
		Build() tfs.TraversalFS
	}

	// CreateResumeFS is a function type that implements ResumeFileSystemBuilder.
	CreateResumeFS func() tfs.TraversalFS
)

// Build creates a file system for resumption.
func (fn CreateResumeFS) Build() tfs.TraversalFS {
	return fn()
}

type (
	// ForestBuilder represents the callback function a client can provide to
	// enable them to specify how the file system for a path is created.	Typically
	// used by unit tests, but can be used by client to specify different file systems
	// than the default; eg if the client needs to integrate with a file system like
	// afero, they can do so by providing the required adapters.
	ForestBuilder interface {
		// Build creates a file system for the specified root path.
		Build(root string) *core.Forest
	}

	// BuildForest is a function type that implements ForestBuilder.
	BuildForest func(root string) *core.Forest
)

// Build creates a file system for the specified root path.
func (fn BuildForest) Build(root string) *core.Forest {
	return fn(root)
}

// Validate checks that the required properties of the Head are set and returns
// an error if not.
func (f *Head) Validate() error {
	if f.Handler == nil {
		return locale.ErrUsageMissingHandler
	}

	return nil
}

// Client is the callback invoked for each file system node found
// during traversal.
func (f *Head) Client() core.Client {
	return f.Handler
}

// Forest returns the client's file system builder, if provided, otherwise it returns
// the default file system builder.
func (f *Head) Forest() BuildForest {
	return f.GetForest
}

// Validate checks that the required properties of the Using are set and returns
// an error if not.
func (f *Using) Validate() error {
	if f.Tree == "" {
		return locale.ErrUsageMissingTreePath
	}

	if f.Subscription == enums.SubscribeUndefined {
		return locale.ErrUsageMissingSubscription
	}

	return f.Head.Validate()
}

// Magnitude returns a string that represents the magnitude of the Using.
func (f *Using) Magnitude() string {
	return "prime"
}

// Validate checks that the required properties of the Relic are set and returns
// an error if not.
func (f *Relic) Validate() error {
	if f.From == "" {
		return locale.ErrUsageMissingRestorePath
	}

	if f.Strategy == enums.ResumeStrategyUndefined {
		return locale.ErrUsageMissingSubscription // TODO: THIS IS WRONG ERROR
	}

	return f.Head.Validate()
}

// Magnitude returns a string that represents the magnitude of the Relic.
func (f *Relic) Magnitude() string {
	return "resume"
}
