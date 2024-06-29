package types

import (
	"context"
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/pref"
)

// package types defines internal types

type (
	// Link represents a single decorator in the chain
	Link interface {
		// Next invokes this decorator which returns true if
		// next link in the chain can be run or false to stop
		// execution of subsequent links.
		Next(node *core.Node) (bool, error)

		// Role indicates the identity of the link
		Role() enums.Role
	}

	// GuardianSealer protects against invalid decorations. There can only
	// be 1 sealer (the master) and currently that only comes into play
	// for fastward resume. An ordinary filter is decorate-able, so it
	// can't be the sealer. It is not mandatory for a master to be registered.
	// When no master is registered, Benign will be used.
	GuardianSealer interface {
		Seal(link Link) error
		IsSealed(top Link) bool
	}

	// Guardian is the gateway to accessing the invocation chain.
	Guardian interface {
		Decorate(link Link) error
		Unwind(role enums.Role) error
	}

	Mediator interface {
		Guardian
		Navigate(ctx context.Context) (core.TraverseResult, error)
		Spawn(ctx context.Context, root string) (core.TraverseResult, error)
	}

	FileSystems struct {
		N fs.FS
		R fs.FS
	}
	// Resources are dependencies required for navigation
	Resources struct {
		FS FileSystems
	}

	// Plugin used to define interaction with supplementary features
	Plugin interface {
		Name() string
		Register() error
		Init() error
	}

	// Restoration; tbd...
	Restoration interface {
		Inject(state pref.ActiveState)
	}

	// Facilities is the interface provided to plugins to enable them
	// to initialise successfully.
	Facilities interface {
		Restoration
	}

	// TraverseController
	TraverseController interface {
		core.Navigator
	}
)

type KernelResult struct {
	Err error
}

func (r KernelResult) Error() error {
	return r.Err
}
