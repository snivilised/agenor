package ui

import (
	"fmt"

	"github.com/snivilised/jaywalk/src/app/report"
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
// It is purely reactive - all methods are event notifications. The UI
// decides how to render each event; no formatting logic lives outside
// the UI layer.
type Manager interface {
	// OnNodeEvent is called per node visit when no action or pipeline
	// is configured.
	OnNodeEvent(e *report.NeutralEvent)

	// OnActionEvent is called when a configured action has been executed
	// against a node.
	OnActionEvent(e *report.ActionEvent)

	// OnPipelineEvent is called when a configured pipeline has been
	// executed against a node.
	OnPipelineEvent(e *report.PipelineEvent)

	// OnSkipEvent is called when an action is skipped for a node because
	// a placeholder in the action's cmd string resolved to a path at or
	// above the traversal root.
	OnSkipEvent(e *report.SkipEvent)

	// OnComplete is called once at the end of a traversal with the full
	// structured outcome.
	OnComplete(t *report.Traversal)
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
