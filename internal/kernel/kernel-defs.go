package kernel

import (
	"context"
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
)

// 📚 package: kernel contains the core traversal functionality. Kernel
// is concerned only with the core task of navigation. Supplementary
// functionality is implemented externally in plugins.

type (
	// NavigatorImpl
	NavigatorImpl interface {
		// Ignite
		Ignite(ignition *types.Ignition)

		// Top
		Top(ctx context.Context,
			ns *navigationStatic,
		) (*types.KernelResult, error)

		// Traverse
		Traverse(ctx context.Context,
			ns *navigationStatic,
			current *core.Node,
		) (bool, error)

		// Result
		Result(ctx context.Context, err error) *types.KernelResult
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
		Decorate(link types.Link) error

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
		Invoke(node *core.Node) error
	}

	// Mutant represents the mutable interface to the Guardian
	Mutant interface {
		Gateway
		Invokable
	}

	// navigationStatic contains static info, ie info that is established during
	// bootstrap and doesn't change after navigation begins. Used to help
	// minimise allocations.
	navigationStatic struct {
		mediator *mediator
		root     string
	}

	// navigationVapour represents short-lived navigation data whose state relates
	// only to the current Node. (equivalent to inspection in extendio)
	navigationVapour struct { // after content has been read
		ns      *navigationStatic
		present *core.Node
		cargo   *Contents
		ents    []fs.DirEntry
	}

	navigationInfo struct { // pre content read
	}

	inspection interface { // after content has been read
		core.Inspection
		static() *navigationStatic
		clear()
	}

	navigationAssets struct {
		ns     navigationStatic
		vapour *navigationVapour
	}
)

func (v *navigationVapour) static() *navigationStatic {
	return v.ns
}

func (v *navigationVapour) Current() *core.Node {
	return v.present
}

func (v *navigationVapour) Contents() core.DirectoryContents {
	return v.cargo
}

func (v *navigationVapour) Entries() []fs.DirEntry {
	return v.ents
}

func (v *navigationVapour) clear() {
	if v.cargo != nil {
		v.cargo.clear()
	} else {
		newEmptyContents()
	}
}

func (v *navigationVapour) Sort(et enums.EntryType) []fs.DirEntry {
	v.cargo.Sort(et)

	// change SortHook to return entries so we don't have to do this switch?
	switch et {
	case enums.EntryTypeAll:
		return v.cargo.All()
	case enums.EntryTypeFolder:
		return v.cargo.folders
	case enums.EntryTypeFile:
		return v.cargo.files
	}

	return nil
}

func (v *navigationVapour) Pick(et enums.EntryType) {
	switch et {
	case enums.EntryTypeAll:
		v.ents = v.cargo.All()
	case enums.EntryTypeFolder:
		v.ents = v.cargo.folders
	case enums.EntryTypeFile:
		v.ents = v.cargo.files
	}
}

func (v *navigationVapour) AssignChildren(children []fs.DirEntry) {
	v.present.Children = children
}

type NodeInvoker func(node *core.Node) error

func (fn NodeInvoker) Invoke(node *core.Node) error {
	return fn(node)
}
