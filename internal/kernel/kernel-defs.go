package kernel

import (
	"context"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	nef "github.com/snivilised/nefilim"
)

// 📦 pkg: kernel - contains the core traversal functionality. Kernel
// is concerned only with the core task of navigation. Supplementary
// functionality is implemented externally in plugins.

type (
	// NavigatorImpl
	NavigatorImpl interface {
		// Ignite
		Ignite(ignition *enclave.Ignition)

		// Top
		Top(ctx context.Context,
			ns *navigationStatic,
		) (*enclave.KernelResult, error)

		// Traverse
		Traverse(ctx context.Context,
			ns *navigationStatic,
			servant core.Servant,
		) (bool, error)

		// Result
		Result(ctx context.Context) *enclave.KernelResult
	}

	// NavigatorDriver
	NavigatorDriver interface {
		Impl() NavigatorImpl
	}

	// Gateway provides a barrier around the guardian to prevent accidental
	// misuse.
	Gateway interface {
		// Decorate is used to wrap the existing client. The decoration will
		// result in the decorator being called first before the existing the
		// client. This needs to behave like chain of responsibility, because
		// the decorator has the choice of wether to pass on the call further
		// down the chain and ultimately the client's callback.
		// The returned bool indicates wether the next in the chain is invoked
		// or not; true, means pass on, false means absorb.
		//
		// A filter decorator will return true if the node matches filter, false
		// otherwise.
		// A fastward resume will act likewise. If we have fastward active, but
		// there is also an underlying filter we have a chain that looks like
		// this:
		// fastward-filter => underlying-filter => callback
		//
		// With this in place, we have to think very carefully as to whether we
		// really need a state machine, because the chain is able to fulfill this
		// purpose.
		//
		// but wait, let's think about wake and sleep. In the normal scenario,
		// fastward will start off in sleeping mode and the filter will be a
		// wake condition. Once we encounter the wake condition, the hibernation
		// decorator needs to be removed. But how do we know that the top of the
		// chain is the hibernate decorator? It must be that the resume hibernate
		// can't be decorated, this suggests some kind of priority/authorisation
		// is required. Because of this, we need a role and we also need to make
		// sure the features are initialised in the correct order to make this
		// happen correctly => sequence manifest
		//
		// role indicates the guise under which the decorator is being applied.
		// Not all roles can be decorated. The fastward-resume decorator can
		// not be decorated. If an attempt is made to Decorate a sealed decorator,
		// an error is returned.
		Decorate(link enclave.Link) error

		// Invoke executes the chain which may or may not end up resulting in
		// the invocation of the client's callback, depending on the contents
		// of the chain.
		// Invoke(node *core.Node) error

		// Unwind removes last link in the chain which is expected to be of
		// role specified.
		Unwind(role enums.Role) error
	}

	// Invokable
	Invokable interface {
		Invoke(servant core.Servant, inspection enclave.Inspection) error
	}

	// navigationStatic contains static info, ie info that is established during
	// bootstrap and doesn't change after navigation begins. Used to help
	// minimise allocations.
	navigationStatic struct {
		mediator     *mediator
		tree         string
		calc         nef.PathCalc
		ofExtent     string
		subscription enums.Subscription
	}

	inspection interface { // after content has been read
		enclave.Inspection
		static() *navigationStatic
		active(tree string,
			forest *core.Forest,
			depth int,
			metrics core.Metrics) *core.ActiveState
		clear()
	}
)
