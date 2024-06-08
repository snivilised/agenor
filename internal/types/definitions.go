package types

// package types internal types

type (
	ContextExpiry interface {
		Expired() // ??? ctx context.Context, cancel context.CancelFunc
	}

	// NavigateResult
	NavigateResult struct {
		Err error
	}
)

type Plugin interface {
	Name() string
	Init() error
}

// UsePlugin invoked by the plugin to the navigator
type UsePlugin interface {
	// this interface needs to be exposed internally but not externally
	Register(plugin Plugin) error
	Interceptor() Interception
	Facilitate() Facilities
}

// Facilities is the interface provided to plugins to enable them
// to initialise successfully.
type Facilities interface {
	Foo()
}

type Interception interface {
	Intercept()
}

func (r NavigateResult) Error() error {
	return r.Err
}
