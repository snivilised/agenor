package flow

import (
	"fmt"
	"io"

	"github.com/snivilised/jaywalk/src/prism"
)

// New constructs a linear renderer for linear-style output.
func New(palette prism.Palette, writer io.Writer, opts ...LinearOption) (prism.Renderer, error) {
	theme, err := prism.NewTheme(palette, writer)
	if err != nil {
		return nil, fmt.Errorf("flow.New: %w", err)
	}

	r := &renderer{
		theme:     theme,
		writer:    writer,
		treeIcons: theme.TreeIcons,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r, nil
}

// Register installs the linear view factory into prism's shared factory map.
// Call this explicitly during application bootstrap before invoking prism.New.
func Register() {
	prism.RegisterFactory(prism.LinearView, func(palette prism.Palette, writer io.Writer) (prism.Renderer, error) {
		return New(palette, writer)
	})
}
