package tapable

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

type (
	HookName                               = string
	HooksCollection[R constraints.Integer] map[R]HookName

	// Container manages a collection of hooks defined for different roles. Since
	// the hook is role specific, there is no way
	Container[R constraints.Integer] struct {
		hooks HooksCollection[R]
	}

	TapContainer[R constraints.Integer] interface {
		Tap(name string, role R) error
	}

	// Activity represents an entity that performs an action that is tapable. The type
	// F should be a function signature, defined by the tapable, so F is not really
	// any, but it can't be enforced to be a func, so do not try to instantiate with
	// a type that is not a func.
	// To enable a client to augment the core action, Tap returns a func of type F which
	// represents the core behaviour. If required, the client should store this and then
	// invoke it as appropriate when it is called back.
	Activity[F any] interface {
		Tap(name string, fn F) (F, error)
	}

	// Role is similar to Activity, but it allows the tapable entity to request
	// that the client specify which specific behaviour needs to be tapped. This is
	// useful when the entity has multiple behaviours that clients may wish to tap.
	// To enable a client to augment the core action, Tap returns a func of type F which
	// represents the core behaviour. If required, the client should store this and then
	// invoke it as appropriate when it is called back.
	Role[F any, R constraints.Integer] interface {
		Tap(name string, role R, fn F) (F, error)
	}

	// WithDefault is a helper that binds together a Hook and its associated
	// default action. The default is the pure, non-tapable underlying function.
	WithDefault[F any, R constraints.Integer] struct {
		Name      HookName
		Role      R
		Action    F
		Default   F
		Container *Container[R]
	}

	AlreadyTappedError[R constraints.Integer] struct {
		role     R
		name     string
		existing string
	}
)

func (e AlreadyTappedError[R]) Error() string {
	return fmt.Sprintf("role '%v', already tapped as '%v'", e.role, e.existing)
}

// NewContainer creates a new instance of a hook container
func NewContainer[R constraints.Integer]() *Container[R] {
	return &Container[R]{
		hooks: make(HooksCollection[R]),
	}
}

// Query prevents client from trying to register multiple hooks
// for the same role.
func (c *Container[R]) Query(name string, role R) error {
	if existing, found := c.hooks[role]; found {
		return &AlreadyTappedError[R]{
			role:     role,
			name:     name,
			existing: existing,
		}
	}

	return nil
}

// Tap taps into the hook and captures the default action.
func (d *WithDefault[F, R]) Tap(name string, role R, fn F) (F, error) {
	if err := d.Container.Query(name, role); err != nil {
		return d.Default, err
	}

	d.Name = name
	d.Role = role
	d.Action = fn

	return d.Default, nil
}
