package shell

import (
	"github.com/snivilised/jaywalk/src/agenor/enums"
)

// LocateFunc reports whether a named executable, builtin, or shell
// function is invokable in the current environment. It returns the
// resolved path or description on success, or an error if the token
// cannot be found. The implementation is platform and environment
// specific - callers should always obtain a LocateFunc from
// Environment.Locate rather than constructing one directly.
type LocateFunc func(name string) (string, error)

// Environment is the result of Detect(). It identifies the shell hosting
// jay and carries a LocateFunc appropriate for that environment.
type Environment struct {
	// Kind is the detected shell environment.
	Kind enums.ShellKind

	// Locate is the platform and environment appropriate function for
	// validating whether a named token is invokable. It is initialised
	// by Detect() and is ready to use without further configuration.
	Locate LocateFunc
}

// Detect inspects the process environment to determine the shell context
// in which jay is running and returns an Environment whose Locate field
// is configured appropriately. An error is returned if detection itself
// fails in an unrecoverable way (e.g. a required subprocess cannot be
// spawned to probe the environment). Callers should treat a non-nil error
// as a fatal startup condition.
func Detect() (Environment, error) {
	return detect()
}
