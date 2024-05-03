# 🪅 Design Notes

<!-- MD013/Line Length -->
<!-- MarkDownLint-disable MD013 -->

<!-- MD014/commands-show-output: Dollar signs used before commands without showing output mark down lint -->
<!-- MarkDownLint-disable MD014 -->

<!-- MD033/no-inline-html: Inline HTML -->
<!-- MarkDownLint-disable MD033 -->

<!-- MD040/fenced-code-language: Fenced code blocks should have a language specified -->
<!-- MarkDownLint-disable MD040 -->

<!-- MD028/no-blanks-blockquote: Blank line inside blockquote -->
<!-- MarkDownLint-disable MD028 -->

<!-- MD010/no-hard-tabs: Hard tabs -->
<!-- MarkDownLint-disable MD010 -->

## Tapable

Given an interface definition:

```go
type (
  // Activity represents an entity that performs an action that is tapable. The type
  // F should be a function signature, defined by the tapable, so F is not really
  // any, but it can't be enforced to be a func, so do not try to instantiate with
  // a type that is not a func.
	Activity[F any] interface {
		On(name string, handler F)
	}

	// Role is similar to Activity, but it allows the tapable entity to request
	// that the client specify which specific behaviour needs to be tapped. This is
	// useful when the entity has multi behaviours that clients may wish to tap.
  // To enable a client to augment the core action, Tap returns a func of type F which
  // represents the core behaviour. If required, the client should store this and then
  // invoke it as appropriate when it is called back.
	Role[F any, E constraints.Integer] interface {
		Tap(name string, role E, fn F) F
	}
)
```

a component can declare that it wants to expose a hook to enable external customisation of a behaviour. That piece of behaviour is the unit which become tab-able. So we have a Component, it's action and a Client that want to tap into this behaviour, by providing a hook.

```go
// Widget is some domain abstraction
type Widget struct {
	Name   string
	Amount int
}

// FuncSimpleWithWidgetAndError is the action func; this function
// can be hooked/is tapable. In fact, this function can either enable core
// functionality to be overridden, decorated with auxiliary behaviour or
// simply be a life-cycle event; ie we have come to a significant point in
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

func (c *Client) WithComponent(component *Component) {
	component.hook.Tap("client", func(name string, amount int) (*Widget, error) {
		widget := &Widget{
			Name:   name,
			Amount: amount,
		}
		return widget, nil
	})
}

```

So we have 3 scenarios to consider: override core behaviour, supplement core behaviour or life cycle.

### Override

In extendio, we have the following hook-able actions:

```go
// TraverseHooks defines the suite of items that can be customised by the client
type TraverseHooks struct {
	QueryStatus   QueryStatusHookFn
	ReadDirectory ReadDirectoryHookFn
	FolderSubPath SubPathHookFn
	FileSubPath   SubPathHookFn
	InitFilters   FilterInitHookFn
	Sort          SortEntriesHookFn
	Extend        ExtendHookFn
}
```

These are all core functions that can be overridden by the client. Furthermore, only a single instance is invoked; ie for each one, just a single action is invoked, its neither a notification or broadcast mechanism.

It may also be pertinent to use the tapable.Role, with the roles being defined as QueryStatus, ReadDirectory, FolderSubPath, ... using a new enum definition. This would serve as an indication to the client that only 1 external entity is able to tap these hook-able actions.

### Supplement

This is a piggy-back on top of the Override scenario, where we would like to augment the core functionality. This then makes the hook, a synchronous notification mechanism, which allows the client to control how core functionality is invoked. The client may choose to invoke custom behaviour before/after the core behaviour or completely override it altogether. The client should be able to invoke the core behaviour, which implies that the Tap should return the default func for this action.

### Life Cycle
