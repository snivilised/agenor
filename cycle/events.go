package cycle

import (
	"github.com/snivilised/traverse/core"
)

type (
	broadcasterFunc[F any] func(listeners []F) F

	Dispatch[F any] struct {
		Invoke      F
		broadcaster broadcasterFunc[F]
	}

	// NotificationCtrl contains the handler function to be invoked. The control
	// is agnostic to the handler's signature and therefore can not invoke it.
	NotificationCtrl[F any] struct {
		Dispatch   Dispatch[F]
		subscribed bool
		listeners  []F
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

var (
	AscendDispatcher  Dispatch[NodeHandler]
	BeginDispatcher   Dispatch[BeginHandler]
	DescendDispatcher Dispatch[NodeHandler]
	EndDispatcher     Dispatch[EndHandler]
	StartDispatcher   Dispatch[HibernateHandler]
	StopDispatcher    Dispatch[HibernateHandler]
)

func init() {
	AscendDispatcher = Dispatch[NodeHandler]{
		Invoke:      func(_ *core.Node) {},
		broadcaster: BroadcastNode,
	}

	BeginDispatcher = Dispatch[BeginHandler]{
		Invoke:      func(_ string) {},
		broadcaster: BroadcastBegin,
	}

	DescendDispatcher = Dispatch[NodeHandler]{
		Invoke:      func(_ *core.Node) {},
		broadcaster: BroadcastNode,
	}

	EndDispatcher = Dispatch[EndHandler]{
		Invoke:      func(_ core.TraverseResult) {},
		broadcaster: BroadcastEnd,
	}

	StartDispatcher = Dispatch[HibernateHandler]{
		Invoke:      func(_ string) {},
		broadcaster: BroadcastHibernate,
	}

	StopDispatcher = Dispatch[HibernateHandler]{
		Invoke:      func(_ string) {},
		broadcaster: BroadcastHibernate,
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
		c.Dispatch.Invoke = handler
		c.subscribed = true

		return
	}

	if c.listeners == nil {
		const size = 2
		c.listeners = make([]F, 0, size)
		c.listeners = append(c.listeners, c.Dispatch.Invoke)
	}

	c.listeners = append(c.listeners, handler)
	c.Dispatch.Invoke = c.broadcaster(c.listeners)
}

func (c *NotificationCtrl[F]) broadcaster(listeners []F) F {
	return c.Dispatch.broadcaster(listeners)
}

func BroadcastBegin(listeners []BeginHandler) BeginHandler {
	return func(root string) {
		for _, listener := range listeners {
			listener(root)
		}
	}
}

func BroadcastEnd(listeners []EndHandler) EndHandler {
	return func(result core.TraverseResult) {
		for _, listener := range listeners {
			listener(result)
		}
	}
}

func BroadcastNode(listeners []NodeHandler) NodeHandler {
	return func(node *core.Node) {
		for _, listener := range listeners {
			listener(node)
		}
	}
}

func BroadcastHibernate(listeners []HibernateHandler) HibernateHandler {
	return func(description string) {
		for _, listener := range listeners {
			listener(description)
		}
	}
}

func BroadcastSimple(listeners []SimpleHandler) SimpleHandler {
	return func() {
		for _, listener := range listeners {
			listener()
		}
	}
}
