package tapable

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// tapable: enables entities to expose hooks

type (
	// ActivityL represents an entity that performs an action that is tapable. The type
	// F should be a function signature, defined by the tapable, so F is not really
	// any, but it can't be enforced to be a func, so do not try to instantiate with
	// a type that is not a func.
	// To enable a client to augment the core action, Tap returns a func of type F which
	// represents the core behaviour. If required, the client should store this and then
	// invoke it as appropriate when it is called back.
	ActivityL[F any] interface {
		Tap(name string, fn F) (F, error)
	}

	// RoleL is similar to Activity, but it allows the tapable entity to request
	// that the client specify which specific behaviour needs to be tapped. This is
	// useful when the entity has multiple behaviours that clients may wish to tap.
	RoleL[F any, E constraints.Integer] interface {
		Tap(name string, role E, fn F) error
	}

	// NotifyL used for life cycle events. There can be multiple subscribers to
	// life cycle events
	NotifyL[F any, E constraints.Integer] interface {
		On(name string, event E, handler F) error
	}
)

// we should have MonoHook (1 client allowed);
// or MultiHook (multiple clients allowed)

// scratch ...

// type FuncSimpleWithError func() error

// type SimpleHookFunc Activity[FuncSimpleWithError]

// Widget is some domain abstraction
type Widget struct {
	Name   string
	Amount int
}

// FuncSimpleWithWidgetAndError is the action func; invoke this function
// that can be hooked/tapable. In fact, this function can either enable core
// functionality to be overridden, decorated with auxiliary behaviour or
// simply as a life-cycle event; ie we have come to a significant point in
// the component's workflow and the some external entity's function needs
// to be invoked.
type FuncSimpleWithWidgetAndError func(name string, amount int) (*Widget, error)

type ActionHook struct {
	name   string
	action FuncSimpleWithWidgetAndError
}

// Tap invoked by client to enable registration of the hook
func (h *ActionHook) Tap(_ string, fn FuncSimpleWithWidgetAndError) {
	h.action = fn
}

type NotificationHook struct {
	name   string
	action FuncSimpleWithWidgetAndError
}

// Tap invoked by client to enable registration of the hook
func (h *NotificationHook) On(_ string, fn FuncSimpleWithWidgetAndError) {
	h.action = fn
}

// Component contains tapable behaviour
type Component struct {
	hook ActionHook
}

func (c *Component) DoWork() {
	if _, err := c.hook.action("work", 0); err != nil {
		panic(fmt.Errorf("work failed: '%v'", err))
	}
}

type Client struct {
}

func (c *Client) WithComponent(from *Component) {
	from.hook.Tap("client", func(name string, amount int) (*Widget, error) {
		widget := &Widget{
			Name:   name,
			Amount: amount,
		}
		return widget, nil
	})
}

// Receiver contains tapable behaviour received as a notification
type Receiver struct {
	hook NotificationHook
}

func (r *Receiver) When(from *Component) {
	from.hook.Tap("client", func(name string, amount int) (*Widget, error) {
		widget := &Widget{
			Name:   name,
			Amount: amount,
		}
		return widget, nil
	})
}

// The above is still confused. Let's start again we have these scenarios:
//
// --> internal (during bootstrap): ==> actually, this is di, not tap
// * 1 component needs to custom another
//
// --> external (via options):
// * notification of life-cycle events (broadcast) | [On/Notify] (eg, OnBegin/On(enums.cycle.begin))
// * customise core behaviour, by role (targeted) | [Tap/Role] (eg, ReadDirectory, role=directory-reader)
// *

// Since we now have finer grain control; ie there are more but smaller packages
// organised as features, each feature can expose its own set of hooks. Having
// said this, I can still only think of nav as needing to expose hooks, but others
// may emerge.
//
// For a component, we have the following situations
// - broadcast, multiple callbacks
// - targeted, single callback
//
// - may expose multiple hooks with different signatures
// the problem this poses is that we can't have a collection of different
// items. This means we need to define a hook container struct that contains the hooks.
// The component aggregates this hook container with a member called hooks.
// For example, in extendio, the options contains a hooks struct TraverseHooks:
//
// type TraverseHooks struct {
// 	QueryStatus   QueryStatusHookFn
// 	ReadDirectory ReadDirectoryHookFn
// 	FolderSubPath SubPathHookFn
// 	FileSubPath   SubPathHookFn
// 	InitFilters   FilterInitHookFn
// 	Sort          SortEntriesHookFn
// 	Extend        ExtendHookFn
// }
//
// But we need to ba able to tap these,
//
// if hooks is of type TraverseHooks, in object Component
// component.hooks.ReadDirectory.tap("name", hookFn)
// therefore ReadDirectoryHookFn, can't be the function, there must be a
// level of indirection in-between,
//
// in TraverseHooks, ReadDirectory must be of type ReadDirectoryHook,
// which is an instantiation of a generic type:
// HookFunc[F any], where HookFunc contains the Tap function
//
// type ReadDirectoryHook tapable.HookFunc[ReadDirectoryHookFn]
//
// type TraverseHooks struct {
// 	ReadDirectory ReadDirectoryHook
// }
