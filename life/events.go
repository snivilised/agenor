package life

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/tapable"
)

type (
	// NotificationCtrl contains the handler function to be invoked. The control
	// is agnostic to the handler's signature and therefore can not invoke it.
	NotificationCtrl[F any] struct {
		dispatcher tapable.Dispatcher[F]
		nop        F
		subscribed bool
		listeners  []F
		muted      bool
	}

	// Events contains the life cycle events that can be subscribed to by handlers.
	// The events are agnostic to the handler's signature and therefore can not invoke
	// the handlers. The events are intended to be used by the binder to bind the
	// underlying notification controllers to the events.
	Events struct { // --> options
		// Ascend is invoked before ascending a directory. The handler function takes a
		// Node as a parameter, which represents the directory being ascended.
		Ascend Event[NodeHandler]

		// Begin is invoked before traversal begins. The handler function takes a BeginState
		// as a parameter, which represents the state at the beginning of traversal. This can
		// be used by the handler to provide context about the traversal.
		Begin Event[BeginHandler]

		// Descend is invoked after descending a directory. The handler function takes a Node
		// as a parameter, which represents the directory being descended.
		Descend Event[NodeHandler]

		// End is invoked at the end of traversal. The handler function takes a TraverseResult
		// as a parameter, which represents the result of the traversal. This can be used by
		// the handler to provide context about the traversal.
		End Event[EndHandler]

		// Wake is invoked when hibernation wakes. The handler function takes a string as a parameter,
		// which represents a description of the wake event. This can be used by the handler to
		// provide context about the wake event.
		Wake Event[HibernateHandler]

		// Sleep is invoked when hibernation sleeps. The handler function takes a string as a parameter,
		// which represents a description of the sleep event. This can be used by the handler to
		// provide context about the sleep event.
		Sleep Event[HibernateHandler]
	}

	// Controls contain notification controls
	// (Since the Controls are only required internally as they are used
	// by binder, they should be moved to an internal package. this
	// would also necessitate moving the handler definitions to core
	// so that they can be shared.)
	Controls struct { // --> binder
		// Ascend is the notification controller for the Ascend event.
		Ascend NotificationCtrl[NodeHandler]

		// Begin is the notification controller for the Begin event.
		Begin NotificationCtrl[BeginHandler]

		// Descend is the notification controller for the Descend event.
		Descend NotificationCtrl[NodeHandler]

		// End is the notification controller for the End event.
		End NotificationCtrl[EndHandler]

		// Wake is the notification controller for the Wake event.
		Wake NotificationCtrl[HibernateHandler]

		// Sleep is the notification controller for the Sleep event.
		Sleep NotificationCtrl[HibernateHandler]
	}
)

// NewNotificationCtrl creates a new NotificationCtrl with the provided
// nop and broadcaster.
func NewNotificationCtrl[F any](nop F,
	broadcaster tapable.Announce[F],
) *NotificationCtrl[F] {
	return &NotificationCtrl[F]{
		dispatcher: tapable.Dispatcher[F]{
			Invoke:      nop,
			Broadcaster: broadcaster,
		},
		nop: nop,
	}
}

// NewControls creates a new Controls with the provided nop and broadcaster functions.
func NewControls() *Controls {
	return &Controls{
		Ascend:  *NewNotificationCtrl[NodeHandler](nopNode, broadcastNode),
		Begin:   *NewNotificationCtrl[BeginHandler](nopBegin, broadcastBegin),
		Descend: *NewNotificationCtrl[NodeHandler](nopNode, broadcastNode),
		End:     *NewNotificationCtrl[EndHandler](nopEnd, broadcastEnd),
		Wake:    *NewNotificationCtrl[HibernateHandler](nopHibernate, broadcastHibernate),
		Sleep:   *NewNotificationCtrl[HibernateHandler](nopHibernate, broadcastHibernate),
	}
}

// MuteAll mutes all the notification controllers in the Controls. This is
// useful when the handlers need to be temporarily disabled without unsubscribing them.
func (c *Controls) MuteAll() {
	c.Ascend.Mute()
	c.Begin.Mute()
	c.Descend.Mute()
	c.End.Mute()
	c.Wake.Mute()
	c.Sleep.Mute()
}

// UnmuteAll unmutes all the notification controllers in the Controls. This is useful
// when the handlers need to be re-enabled after being temporarily disabled by MuteAll.
// Unmuting allows the handlers to be invoked again when the corresponding events
// are dispatched.
func (c *Controls) UnmuteAll() {
	c.Ascend.Unmute()
	c.Begin.Unmute()
	c.Descend.Unmute()
	c.End.Unmute()
	c.Sleep.Unmute()
}

// Bind attaches the underlying notification controllers to the
// Events.
func (e *Events) Bind(cs *Controls) {
	e.Ascend = &cs.Ascend
	e.Begin = &cs.Begin
	e.Descend = &cs.Descend
	e.End = &cs.End
	e.Wake = &cs.Wake
	e.Sleep = &cs.Sleep
}

// On subscribes to a life cycle event
func (c *NotificationCtrl[F]) On(handler F) {
	if !c.subscribed {
		c.dispatcher.Invoke = handler
		c.subscribed = true

		return
	}

	if c.listeners == nil {
		const size = 2

		c.listeners = make([]F, 0, size)
		c.listeners = append(c.listeners, c.dispatcher.Invoke)
	}

	c.listeners = append(c.listeners, handler)
	c.dispatcher.Invoke = c.dispatcher.Broadcaster(c.listeners)
}

// Dispatch returns the handler function to be invoked. If the notification
// controller is muted, it returns the nop function instead. This allows the
// caller to invoke the returned function without worrying about whether the
// notification controller is muted or not.
func (c *NotificationCtrl[F]) Dispatch() F {
	if c.muted {
		return c.nop
	}

	return c.dispatcher.Invoke
}

// Off unsubscribes from a life cycle event. This is achieved by resetting the
// notification controller to its initial state, which is effectively the same
// as unsubscribing all handlers.
func (c *NotificationCtrl[F]) Off() {
	c.dispatcher.Invoke = c.nop
	c.subscribed = false
	c.listeners = nil
}

// Mute temporarily disables the handlers by setting the muted flag to true. When
// a notification controller is muted, the Dispatch method will return the nop
// function instead of the actual handler function, effectively preventing any
// handlers from being invoked when the corresponding event is dispatched.
func (c *NotificationCtrl[F]) Mute() {
	c.muted = true
}

// Unmute re-enables the handlers by setting the muted flag to false. When a
// notification controller is unmuted, the Dispatch method will return the
// actual handler function again, allowing the handlers to be invoked when the
// corresponding event is dispatched.
func (c *NotificationCtrl[F]) Unmute() {
	c.muted = false
}

func broadcastBegin(listeners []BeginHandler) BeginHandler {
	return func(state *BeginState) {
		for _, listener := range listeners {
			listener(state)
		}
	}
}

func nopBegin(*BeginState) {}

func broadcastEnd(listeners []EndHandler) EndHandler {
	return func(result core.TraverseResult) {
		for _, listener := range listeners {
			listener(result)
		}
	}
}

func nopEnd(_ core.TraverseResult) {}

func broadcastNode(listeners []NodeHandler) NodeHandler {
	return func(node *core.Node) {
		for _, listener := range listeners {
			listener(node)
		}
	}
}

func nopNode(_ *core.Node) {}

func broadcastHibernate(listeners []HibernateHandler) HibernateHandler {
	return func(description string) {
		for _, listener := range listeners {
			listener(description)
		}
	}
}

func nopHibernate(_ string) {}
