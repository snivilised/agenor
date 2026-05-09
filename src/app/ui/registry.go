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
// View factory functions - on-demand creation
// ---------------------------------------------------------------------------

// newLinearPresenter creates a linear stream view presenter with the
// given palette. The presenter wraps a prism stream renderer and
// applies the palette's theme settings (colors, icons, styles). Custom
// tree icons from the palette are explicitly applied via WithIcons to
// ensure they override the defaults.
func newLinearPresenter(palette prism.Palette) (report.Presenter, error) {
	renderer, err := prism.New(
		prism.StreamView,
		palette,
		os.Stdout,
		prism.WithIcons(palette.TreeIcons),
	)
	if err != nil {
		return nil, err
	}

	return &linear{renderer: renderer}, nil
}

// ---------------------------------------------------------------------------
// Public API
// ---------------------------------------------------------------------------

// New returns the Presenter for the requested mode, constructed with
// the given palette. Only the selected view is instantiated; other views
// are not created. Returns an error if the mode is unknown or if the
// palette contains unrecognised colour names.
func New(mode string, palette prism.Palette) (report.Presenter, error) {
	if mode == "" {
		mode = ModeDefault
	}

	switch mode {
	case ModeLinear:
		return newLinearPresenter(palette)
	default:
		return nil, fmt.Errorf(
			"unknown display mode %q (valid modes: %s)",
			mode,
			strings.Join(availableModes(), ", "),
		)
	}
}

// availableModes returns all known mode names, for error messages.
func availableModes() []string {
	return []string{ModeLinear}
}
