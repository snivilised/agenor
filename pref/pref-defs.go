package pref

import "github.com/snivilised/traverse/enums"

// 📦 pkg: pref - contains user option definitions; do not use anything
// in kernel (cyclic).

const (
	badge = "badge: option-requester"
)

type (
	TraversalState struct {
		Tree        string
		Hibernation enums.Hibernation
		CurrentPath string
		Depth       int
	}

	// Restorer function defined by client invoked as part of the resume
	// process
	Restorer func(o *Options, ts *TraversalState) error
)
