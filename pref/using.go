package pref

import (
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
}

// Validate checks that the properties on Using are all valid.
func (u Using) Validate() error {
	if u.Root == "" {
		return UsageError{
			message: "missing root path",
		}
	}

	if u.Subscription == enums.SubscribeUndefined {
		return UsageError{
			message: "missing subscription",
		}
	}

	if u.Handler == nil {
		return UsageError{
			message: "missing handler",
		}
	}

	return nil
}

// As is similar to Using except that it is required for Resume
// exclusively and contains properties required to support
// restoring a session from a previously terminated run.
type As struct {
	Using

	// From is the path to the resumption file from which a prior
	// traverse session is loaded.
	From string

	// Strategy represent what type of resume is run.
	Strategy enums.ResumeStrategy
}

// Validate checks that the properties on Using and As are all valid.
func (a As) Validate() error {
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

	return a.Using.Validate()
}
