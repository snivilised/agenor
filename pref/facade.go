package pref

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/tfs"
)

type (
	// Facade
	Facade interface {
		Path() string
		Client() core.Client
		Forest() BuildForest
		Validate() error
		OfExtent() string
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

	// Relic contains information required to instigate a resume
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
	TraverseFileSystemBuilder interface {
		Build(root string) tfs.TraversalFS
	}

	CreateTraverseFS func(root string) tfs.TraversalFS
)

func (fn CreateTraverseFS) Build(root string) tfs.TraversalFS {
	return fn(root)
}

type (
	ResumeFileSystemBuilder interface {
		Build() tfs.TraversalFS
	}

	CreateResumeFS func() tfs.TraversalFS
)

func (fn CreateResumeFS) Build() tfs.TraversalFS {
	return fn()
}

type (
	ForestBuilder interface {
		Build(root string) *core.Forest
	}

	BuildForest func(root string) *core.Forest
)

func (fn BuildForest) Build(root string) *core.Forest {
	return fn(root)
}

func (f *Head) Validate() error {
	if f.Handler == nil {
		return locale.ErrUsageMissingHandler
	}

	return nil
}

func (f *Using) Path() string {
	return f.Tree
}

func (f *Using) Client() core.Client {
	return f.Handler
}

func (f *Using) Forest() BuildForest {
	return f.GetForest
}

func (f *Using) Validate() error {
	if f.Tree == "" {
		return locale.ErrUsageMissingTreePath
	}

	if f.Subscription == enums.SubscribeUndefined {
		return locale.ErrUsageMissingSubscription
	}

	return f.Head.Validate()
}

func (f *Using) OfExtent() string {
	return "prime"
}

func (f *Relic) Path() string {
	return f.From
}

func (f *Relic) Client() core.Client {
	return f.Handler
}

func (f *Relic) Forest() BuildForest {
	return f.GetForest
}

func (f *Relic) Validate() error {
	if f.From == "" {
		return locale.ErrUsageMissingRestorePath
	}

	if f.Strategy == enums.ResumeStrategyUndefined {
		return locale.ErrUsageMissingSubscription // TODO: THIS IS WRONG ERROR
	}

	return f.Head.Validate()
}

func (f *Relic) OfExtent() string {
	return "resume"
}
