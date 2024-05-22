package pref

import (
	"github.com/snivilised/traverse/cycle"
)

// this package should be internal

type (
	// Binder contains items derived from Options
	Binder struct {
		Controls cycle.Controls
		Loaded   *LoadInfo
	}
)

func NewBinder() *Binder {
	return &Binder{
		Controls: cycle.NewControls(),
	}
}
