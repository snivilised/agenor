package pref

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
)

// Using contains essential properties required for a traversal. If
// any of the required properties are missing, then traversal will
// result in an error indicating as such.
type Using struct {
	// Root is the path of the directory tree to be traversed
	Root string

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

	// GetReadDirFS is optional and enables the client to specify how the
	// file system for a path is created
	GetReadDirFS CreateReadDirFS

	// GetQueryStatusFS is optional and enables the client to specify how the
	// file system for a path is created. When specified can probably use
	// the same instance used to create the ReadDi fs, that is because
	// fstest.MapFS implements the required method Stat.
	GetQueryStatusFS CreateQueryStatusFS
}

// Validate checks that the properties on Using are all valid.
func (u Using) Validate() error {
	if u.Root == "" {
		return UsageError{
			message: "missing root path",
		}
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
func (a Was) Validate() error { //nolint:gocritic // heavy, so what, low frequency
	if a.From == "" {
		return UsageError{
			message: "missing restore from path",
		}
	}

	if a.Strategy == enums.ResumeStrategyUndefined {
		return UsageError{
			message: "missing subscription",
		}
	}

	return validate(&a.Using)
}

func validate(using *Using) error {
	if using.Subscription == enums.SubscribeUndefined {
		return UsageError{
			message: "missing subscription",
		}
	}

	if using.Handler == nil {
		return UsageError{
			message: "missing handler",
		}
	}

	return nil
}

type (
	ReadDirFileSystemBuilder interface {
		Build() fs.ReadDirFS
	}

	CreateReadDirFS func() fs.ReadDirFS
)

func (fn CreateReadDirFS) Build() fs.ReadDirFS {
	return fn()
}

type (
	QueryStatusFileSystemBuilder interface {
		Build(fsys fs.FS) fs.StatFS
	}

	CreateQueryStatusFS func(fsys fs.FS) fs.StatFS
)

func (fn CreateQueryStatusFS) Build(fsys fs.FS) fs.StatFS {
	return fn(fsys)
}
