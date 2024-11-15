package tapable

import (
	"github.com/snivilised/agenor/core"
)

type (
	Hooks struct {
		FileSubPath      Hook[core.SubPathHook, core.ChainSubPathHook]
		DirectorySubPath Hook[core.SubPathHook, core.ChainSubPathHook]
		ReadDirectory    Hook[core.ReadDirectoryHook, core.ChainReadDirectoryHook]
		QueryStatus      Hook[core.QueryStatusHook, core.ChainQueryStatusHook]
		Sort             Hook[core.SortHook, core.ChainSortHook]
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

	// ProvideListeners represents an entity that contains the list of listeners
	// that needs to be provided to the hook broadcaster.
	ProvideListeners[C any] func() []C
)

func (fn ProvideListeners[C]) get() []C {
	return fn()
}

type (
	// broadcastAdapter adapts a default hook so that it can be broadcasted to
	// all members in the chain.
	// F: core hook function
	// C: chained client hook, ie the hook the client provides when they call Chain
	// B: pre-defined broadcaster function
	broadcastAdapter[F, C, B any] interface {
		// attach effectively adds a new listener to the broadcast chain
		// which in itself is attached to the default hook. That is to say,
		// initially, each hook is defined to run default functionality. If
		// an entity registers interest in augmenting the default functionality
		// by invoking Chain on the HookCtrl, then the broadcaster is employed
		// to invoke the default hook, the result (if a result is generated) of
		// which is passed down the invocation chain. This allows subsequent
		// parties to modify the ultimate result.
		attach(def F, provider listenerProvider[C], broadcaster B) F
	}

	// attacher
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

func (c *HookCtrl[F, C, B]) Tap(handler F) {
	c.handler = handler
}

func (c *HookCtrl[F, C, B]) Chain(handler C) {
	if c.listeners == nil {
		c.listeners = []C{handler}
		c.handler = c.adapter.attach(c.handler, c, c.broadcaster)

		return
	}

	c.listeners = append(c.listeners, handler)
}

func (c *HookCtrl[F, C, B]) Default() F {
	return c.def
}

func (c *HookCtrl[F, C, B]) Invoke() F {
	return c.handler
}

func (c *HookCtrl[F, C, B]) get() []C {
	return c.listeners
}
