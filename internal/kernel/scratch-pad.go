package kernel

import (
	"github.com/snivilised/traverse/pref"
)

// scratchPad contains core data that is derived from the options. Any
// decoration that occurs happens on the scratch pad, this way we can
// leave the options to remain unchanged so it always reflects what the
// client set.
type scratchPad struct {
	o *pref.Options
}

func newScratch(o *pref.Options) *scratchPad {
	return &scratchPad{
		o: o,
	}
}
