package pref

import (
	"github.com/snivilised/traverse/core"
)

// 📦 pkg: pref - contains user option definitions; do not use anything
// in kernel (cyclic).

const (
	badge = "badge: option-requester"
)

type (
	// Restorer function defined by client invoked as part of the resume
	// process
	Restorer func(o *Options, active *core.ActiveState) error
)
