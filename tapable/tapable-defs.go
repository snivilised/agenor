package tapable

type (
	// Invokable represents an entity that can be invoked to return a function of type F.
	// F: core hook function
	Invokable[F any] interface {
		// Invoke returns the hook function for execution.
		Invoke() F
	}

	// Hook represents core functionality that can be hooked by multiple
	// entities via a chain starting off with the default hook.
	// F: core hook function
	// C: chained client hook, ie the hook the client provides when they call Chain
	// B: pre-defined broadcaster function
	Hook[F, C any] interface {
		Invokable[F]
		// Tap overrides the default tap-able core function
		Tap(handler F)

		// Chain augments previously registered behaviour. The default
		// behaviour will be invoked first, followed by any other
		// handlers registered in the order of registration.
		Chain(handler C)

		// Default returns the default function for this hook
		Default() F
	}

	// Announce is the function type that defines what to broadcast to multiple listeners.
	Announce[F any] func(listeners []F) F

	// Dispatcher is the struct that holds the core hook function and the broadcaster function.
	Dispatcher[F any] struct {
		// Invoke is the core hook function that will be executed when the hook is called.
		Invoke F

		// Broadcaster is the function that defines what to broadcast to multiple listeners.
		Broadcaster Announce[F]
	}
)
