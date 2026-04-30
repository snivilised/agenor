package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/prism"
)

// ---------------------------------------------------------------------------
// Named display modes - legal values for --tui
// ---------------------------------------------------------------------------

const (
	// ModeLinear is the default stream view. One styled line per node
	// via prism's stream renderer with lipgloss formatting.
	ModeLinear = "linear"

	// ModeDefault is the display used when --tui is not specified.
	ModeDefault = ModeLinear
)

// ---------------------------------------------------------------------------
// Registry
// ---------------------------------------------------------------------------

// Factory is the constructor signature all display mode implementations
// must satisfy. It receives the resolved Palette so that the prism
// renderer is built with the correct colours at construction time.
type Factory func(palette prism.Palette) (report.Presenter, error)

var registry = map[string]Factory{
	ModeLinear: func(palette prism.Palette) (report.Presenter, error) {
		renderer, err := prism.New(prism.StreamView, palette, os.Stdout)
		if err != nil {
			return nil, err
		}

		return &linear{renderer: renderer}, nil
	},
}

// RegisterMode adds a new display mode to the registry. Returns an
// error if the name is already registered. Callers should treat a
// duplicate registration as a programming error and fail at startup.
func RegisterMode(name string, f Factory) error {
	if _, exists := registry[name]; exists {
		return fmt.Errorf("display mode %q is already registered", name)
	}

	registry[name] = f

	return nil
}

// New returns the Presenter for the requested mode, constructed with
// the given palette. Returns an error if the mode is not registered or
// if the palette contains unrecognised colour names.
func New(mode string, palette prism.Palette) (report.Presenter, error) {
	if mode == "" {
		mode = ModeDefault
	}

	f, ok := registry[mode]
	if !ok {
		return nil, fmt.Errorf(
			"unknown display mode %q (valid modes: %s)",
			mode,
			strings.Join(registeredModes(), ", "),
		)
	}

	return f(palette)
}

// registeredModes returns all known mode names, for error messages.
func registeredModes() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}

	return names
}
