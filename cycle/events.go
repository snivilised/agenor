package cycle

import (
	"github.com/snivilised/traverse/core"
)

type (
	announce[F any] func(listeners []F) F

	dispatcher[F any] struct {
		invoke      F
		broadcaster announce[F]
	}

	// NotificationCtrl contains the handler function to be invoked. The control
	// is agnostic to the handler's signature and therefore can not invoke it.
	NotificationCtrl[F any] struct {
		dispatcher dispatcher[F]
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
		Start   Event[HibernateHandler]
		Stop    Event[HibernateHandler]
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
		Start   NotificationCtrl[HibernateHandler]
		Stop    NotificationCtrl[HibernateHandler]
	}
)

func NewControls() Controls {
	return Controls{
		Ascend: NotificationCtrl[NodeHandler]{
			dispatcher: dispatcher[NodeHandler]{
				invoke:      nopNode,
				broadcaster: broadcastNode,
			},
			nop: nopNode,
		},
		Begin: NotificationCtrl[BeginHandler]{
			dispatcher: dispatcher[BeginHandler]{
				invoke:      nopBegin,
				broadcaster: broadcastBegin,
			},
			nop: nopBegin,
		},
		Descend: NotificationCtrl[NodeHandler]{
			dispatcher: dispatcher[NodeHandler]{
				invoke:      nopNode,
				broadcaster: broadcastNode,
			},
			nop: nopNode,
		},
		End: NotificationCtrl[EndHandler]{
			dispatcher: dispatcher[EndHandler]{
				invoke:      nopEnd,
				broadcaster: broadcastEnd,
			},
			nop: nopEnd,
		},
		Start: NotificationCtrl[HibernateHandler]{
			dispatcher: dispatcher[HibernateHandler]{
				invoke:      nopHibernate,
				broadcaster: broadcastHibernate,
			},
			nop: nopHibernate,
		},
		Stop: NotificationCtrl[HibernateHandler]{
			dispatcher: dispatcher[HibernateHandler]{
				invoke:      nopHibernate,
				broadcaster: broadcastHibernate,
			},
			nop: nopHibernate,
		},
	}
}

// Bind attaches the underlying notification controllers to the
// Events.
func (e *Events) Bind(cs *Controls) {
	e.Ascend = &cs.Ascend
	e.Begin = &cs.Begin
	e.Descend = &cs.Descend
	e.End = &cs.End
	e.Start = &cs.Start
	e.Stop = &cs.Stop
}

// On subscribes to a life cycle event
func (c *NotificationCtrl[F]) On(handler F) {
	if !c.subscribed {
		c.dispatcher.invoke = handler
		c.subscribed = true

		return
	}

	if c.listeners == nil {
		const size = 2
		c.listeners = make([]F, 0, size)
		c.listeners = append(c.listeners, c.dispatcher.invoke)
	}

	c.listeners = append(c.listeners, handler)
	c.dispatcher.invoke = c.dispatcher.broadcaster(c.listeners)
}

func (c *NotificationCtrl[F]) Dispatch() F {
	if c.muted {
		return c.nop
	}

	return c.dispatcher.invoke
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

func broadcastSimple(listeners []SimpleHandler) SimpleHandler {
	return func() {
		for _, listener := range listeners {
			listener()
		}
	}
}
