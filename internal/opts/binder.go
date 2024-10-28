package opts

import (
	"github.com/snivilised/traverse/life"
)

type (
	// Binder contains items derived from Options
	Binder struct {
		Controls life.Controls
	}
)

func NewBinder() *Binder {
	return &Binder{
		Controls: life.NewControls(),
	}
}
