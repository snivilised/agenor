package tapable

// 📚 package: tapable enables entities to expose hooks

type (
	Invokable[F any] interface {
		// Invoke returns the hook function for execution.
		Invoke() F
	}

	Hook[F any] interface {
		Invokable[F]
		// Tap overrides the default tap-able core function
		Tap(handler F)

		// Default returns the default function for this hook
		Default() F
	}
)
