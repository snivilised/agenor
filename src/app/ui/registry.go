package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/prism"
)

// ---------------------------------------------------------------------------
// Named display modes - these are the legal values for --tui
// ---------------------------------------------------------------------------

const (
	// ModeLinear is the default stream view. Writes one styled line per
	// node using prism's stream renderer with lipgloss formatting.
	ModeLinear = "linear"

	// ModeDefault is the display used when --tui is not specified.
	ModeDefault = ModeLinear
)

// ---------------------------------------------------------------------------
// Registry
// ---------------------------------------------------------------------------

// Factory is the constructor signature all display mode implementations
// must satisfy. It is exported so that external packages can reference
// the type when calling RegisterMode.
type Factory func() report.Presenter

var registry = map[string]Factory{
	ModeLinear: func() report.Presenter {
		return &linear{
			renderer: prism.New(prism.StreamView, os.Stdout),
		}
	},
}

// RegisterMode adds a new display mode to the registry. Returns an error
// if the name is already registered. Callers should treat a duplicate
// registration as a programming error and fail startup explicitly.
func RegisterMode(name string, f Factory) error {
	if _, exists := registry[name]; exists {
		return fmt.Errorf("display mode '%s' is already registered", name)
	}

	registry[name] = f

	return nil
}

// New returns the Presenter for the requested mode. Returns an error if
// the mode has not been registered.
func New(mode string) (report.Presenter, error) {
	if mode == "" {
		mode = ModeDefault
	}

	f, ok := registry[mode]
	if !ok {
		return nil, fmt.Errorf(
			"unknown display mode '%s' (valid modes: '%s')",
			mode,
			strings.Join(registeredModes(), ", "),
		)
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
