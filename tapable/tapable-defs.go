package tapable

type (
	Hook[F any] interface {
		// Tap overrides the default tap-able core function
		Tap(handler F)

		// Default returns the default function for this hook
		Default() F

		// Invoke returns the hook function for execution.
		Invoke() F
	}
)
