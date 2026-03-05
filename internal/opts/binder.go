package opts

import (
	"github.com/snivilised/agenor/life"
)

type (
	// Binder contains items derived from Options
	Binder struct {
		// Controls contains the controls for the current traversal session
		Controls *life.Controls
	}
)

// NewBinder creates a new Binder instance with default values
func NewBinder() *Binder {
	return &Binder{
		Controls: life.NewControls(),
	}
}
