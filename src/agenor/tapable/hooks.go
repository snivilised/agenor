package tapable

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
)

type (
	// Hooks contains the controllers for all the hooks that can be registered to
	// during traversal. Each hook is defined by its default function, and the
	// broadcaster function that is used to chain additional functionality to the
	// default when clients call the Chain method on the hook controller.
	Hooks struct {
		// FileSubPath is a hook that allow clients to modify the sub-path
		// that is used for a file.
		FileSubPath Hook[core.SubPathHook, core.ChainSubPathHook]

		// DirectorySubPath is a hook that allow clients to modify the sub-path
		// that is used for a directory.
		DirectorySubPath Hook[core.SubPathHook, core.ChainSubPathHook]

		// ReadDirectory is a hook that allows clients to override the reading of
		// the contents directory.
		ReadDirectory Hook[core.ReadDirectoryHook, core.ChainReadDirectoryHook]

		// QueryStatus is a hook that allows clients to override the querying of
		// a node's status. QueryStatus is only used for the top node of a traversal,
		// and is not used for any other nodes. This is because the top node is the
		// only node that is not guaranteed to exist in the forest, and therefore
		// may require special handling.
		QueryStatus Hook[core.QueryStatusHook, core.ChainQueryStatusHook]

		// Sort is a hook that allows clients to override the sorting of a directory's
		// contents. This is used in conjunction with the SortBehaviour to determine
		// how the sorting should be applied.
		Sort Hook[core.SortHook, core.ChainSortHook]
	}

	// HookCtrl contains the handler function to be invoked.
	// Each HookCtrl is defined by its default function. If other parties wish
	// to chain functionality onto the result of invoking the default, they can
	// invoke the Chain method. When this happens, the registered handler is
	// converted into a broadcaster. The broadcaster then becomes responsible for
	// invoking the default function, then chaining the result of this to the new
	// chain that is formed by the Chain method.
	// Note that Chain is distinct from Tap which is used to replace the default
	// functionality entirely. It does not make sense to Chain and Tap the same
	// hook.
	//
	// F: core hook function
	// C: chained client hook, ie the hook the client provides when they call Chain
	// B: pre-defined broadcaster function
	HookCtrl[F, C, B any] struct {
		handler     F
		def         F
		broadcaster B
		adapter     attacher[F, C, B]
		listeners   []C
	}
)

type (
	// listenerProvider
	// C: chained client hook, ie the hook the client provides when they call Chain
	listenerProvider[C any] interface {
		// get returns the collection of interested listeners
		get() []C
	}
)

type (
	attacher[F, C, B any] func(def F, provider listenerProvider[C], broadcaster B) F
)

func (fn attacher[F, C, B]) attach(def F, provider listenerProvider[C], broadcaster B) F {
	return fn(def, provider, broadcaster)
}

// NewHookCtrl creates a new hook controller
func NewHookCtrl[F, C, B any](
	def F,
	broadcaster B,
	adapter attacher[F, C, B],
) *HookCtrl[F, C, B] {
	// The control is agnostic to the handler's signature and therefore can not
	// invoke it; this is the reason why there is delegation to hook specific
	// functions which are signature aware, in particular, the broadcaster and
	// the adapter, eg: GetSubPathBroadcaster/SubPathAttacher,
	//
	return &HookCtrl[F, C, B]{
		handler:     def,
		def:         def,
		broadcaster: broadcaster,
		adapter:     adapter,
	}
}

// Tap registers the handler as the function to be invoked when the hook is invoked.
// This replaces any previously registered handler, and does not chain to the
// default functionality. If the client wishes to chain to the default functionality,
// they should use the Chain method instead.
func (c *HookCtrl[F, C, B]) Tap(handler F) {
	c.handler = handler
}

// Chain registers the handler as the function to be invoked when the hook is invoked.
// This chains to the default functionality, and any previously registered handlers,
// by converting the handler into a broadcaster, and invoking the adapter to
// attach this broadcaster to the default function. If the client wishes to
// replace the default functionality entirely, they should use the Tap method instead.
func (c *HookCtrl[F, C, B]) Chain(handler C) {
	if c.listeners == nil {
		c.listeners = []C{handler}
		c.handler = c.adapter.attach(c.handler, c, c.broadcaster)

		return
	}

	c.listeners = append(c.listeners, handler)
}

// Default returns the default function for the hook, which is the function that is
// invoked if no handlers are registered, or if the handlers chain to the default
// functionality.
func (c *HookCtrl[F, C, B]) Default() F {
	return c.def
}

// Invoke returns the current handler for the hook, which is the function that is
// invoked when the hook is invoked. This may be the default function if no handlers
// have been registered, or it may be a chained function if handlers have been
// registered using the Chain method, or it may be a completely replaced function
// if a handler has been registered using the Tap method.
func (c *HookCtrl[F, C, B]) Invoke() F {
	return c.handler
}

func (c *HookCtrl[F, C, B]) get() []C { //nolint:unused // ok
	return c.listeners
}
