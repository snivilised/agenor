package pref

import (
	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/locale"
)

// Using contains essential properties required for a traversal. If
// any of the required properties are missing, then traversal will
// result in an error indicating as such.
type Using struct {
	// Tree is the root of the directory tree to be traversed and
	// should not be confused with the Root of the file system when
	// the file system in use is relative.
	Tree string

	// Subscription indicates which file system nodes the client's
	// callback function will be invoked for.
	Subscription enums.Subscription

	// Handler is the callback function invoked for each encountered
	// file system node.
	Handler core.Client

	// O is the optional Options entity. If provided, then these
	// options will be used verbatim, without requiring WithXXX
	// options setters. This is useful if multiple traversals are
	// required, eg a preview traversal followed by a full
	// traversal; in this case the full traversal can reuse the
	// same options that was used in the preview, by setting this
	// property.
	O *Options

	// GetTraverseFS is optional and enables the client to specify how the
	// file system for a path is created
	GetTraverseFS CreateTraverseFS
}

// Validate checks that the properties on Using are all valid.
func (u Using) Validate() error {
	if u.Tree == "" {
		return locale.ErrUsageMissingTreePath
	}

	return validate(&u)
}

// Was is similar to Using except that it is required for Resume
// exclusively and contains properties required to support
// restoring a session from a previously terminated run.
type Was struct {
	Using

	// From is the path to the resumption file from which a prior
	// traverse session is loaded.
	From string

	// Strategy represent what type of resume is run.
	Strategy enums.ResumeStrategy
}

// Validate checks that the properties on Using and Was are all valid.
func (a Was) Validate() error {
	if a.From == "" {
		return locale.ErrUsageMissingRestorePath
	}

	if a.Strategy == enums.ResumeStrategyUndefined {
		return locale.ErrUsageMissingSubscription
	}

	return validate(&a.Using)
}

func validate(using *Using) error {
	if using.Subscription == enums.SubscribeUndefined {
		return locale.ErrUsageMissingSubscription
	}

	if using.Handler == nil {
		return locale.ErrUsageMissingHandler
	}

	return nil
}

type (
	TraverseFileSystemBuilder interface {
		Build(root string) nef.TraverseFS
	}

	CreateTraverseFS func(root string) nef.TraverseFS
)

func (fn CreateTraverseFS) Build(root string) nef.TraverseFS {
	return fn(root)
}
