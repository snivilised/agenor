package ui

import (
	"fmt"

	"github.com/snivilised/jaywalk/src/agenor/core"
)

// ---------------------------------------------------------------------------
// Named display modes - these are the legal values for --tui
// ---------------------------------------------------------------------------

const (
	// ModeLinear is the default plain-text display. Writes one line per node
	// using fmt.Println. No external dependencies.
	ModeLinear = "linear"

	// ModeDefault is the display used when --tui is not specified.
	ModeDefault = ModeLinear
)

// ---------------------------------------------------------------------------
// Manager interface
// ---------------------------------------------------------------------------

// Manager is the single interface all UI implementations satisfy.
// Command handlers and agenor node callbacks interact only with this
// interface, never with a concrete type.
type Manager interface {
	// OnNode is called for every node the traversal visits. The node is
	// obtained by the command layer via servant.Node() before being passed
	// here. It is safe to call from multiple goroutines (implementations
	// must ensure this).
	OnNode(node *core.Node) error

	// Info writes a general informational message to the display.
	Info(msg string)

	// Warn writes a warning message to the display.
	Warn(msg string)

	// Error writes an error message to the display.
	Error(msg string)
}

// ---------------------------------------------------------------------------
// Factory
// ---------------------------------------------------------------------------

// ErrUnknownMode is returned by New when the requested mode is not registered.
type ErrUnknownMode struct {
	Mode string
}

func (e *ErrUnknownMode) Error() string {
	return fmt.Sprintf("ui: unknown display mode %q (valid modes: %v)", e.Mode, registeredModes())
}

// factory maps a mode name to its constructor.
type factory func() Manager

var registry = map[string]factory{
	ModeLinear: func() Manager { return &linear{} },
}

// RegisterMode adds a new display mode to the registry. Call this from
// an init() function in the package that provides the implementation.
// Panics if the name is already registered.
func RegisterMode(name string, f factory) {
	if _, exists := registry[name]; exists {
		panic(fmt.Sprintf("ui: display mode %q already registered", name))
	}
	registry[name] = f
}

// New returns the Manager for the requested mode. Returns an error if the
// mode is not registered.
func New(mode string) (Manager, error) {
	if mode == "" {
		mode = ModeDefault
	}
	f, ok := registry[mode]
	if !ok {
		return nil, &ErrUnknownMode{Mode: mode}
	}
	return f(), nil
}

// registeredModes returns a slice of all known mode names, for error messages.
func registeredModes() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
