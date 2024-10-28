package life

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/tapable"
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

	Events struct { // --> options
		Ascend  Event[NodeHandler]
		Begin   Event[BeginHandler]
		Descend Event[NodeHandler]
		End     Event[EndHandler]
		Wake    Event[HibernateHandler]
		Sleep   Event[HibernateHandler]
	}

	// since the Controls are only required internally as they are used
	// by binder, they should be moved to an internal package. this
	// would also necessitate moving the handler definitions to core
	// so that they can be shared.

	Controls struct { // --> binder
		Ascend  NotificationCtrl[NodeHandler]
		Begin   NotificationCtrl[BeginHandler]
		Descend NotificationCtrl[NodeHandler]
		End     NotificationCtrl[EndHandler]
		Wake    NotificationCtrl[HibernateHandler]
		Sleep   NotificationCtrl[HibernateHandler]
	}
)

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

func NewControls() Controls {
	return Controls{
		Ascend:  *NewNotificationCtrl[NodeHandler](nopNode, broadcastNode),
		Begin:   *NewNotificationCtrl[BeginHandler](nopBegin, broadcastBegin),
		Descend: *NewNotificationCtrl[NodeHandler](nopNode, broadcastNode),
		End:     *NewNotificationCtrl[EndHandler](nopEnd, broadcastEnd),
		Wake:    *NewNotificationCtrl[HibernateHandler](nopHibernate, broadcastHibernate),
		Sleep:   *NewNotificationCtrl[HibernateHandler](nopHibernate, broadcastHibernate),
	}
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

func (c *NotificationCtrl[F]) Dispatch() F {
	if c.muted {
		return c.nop
	}

	return c.dispatcher.Invoke
}

func (c *NotificationCtrl[F]) Mute() {
	c.muted = true
}

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
