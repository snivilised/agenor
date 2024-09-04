package hiber

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/override"
)

const (
	badge = "badge: hibernator"
)

// 📦 pkg: hiber - represents the facility to be able to start navigation in
// hibernated state, ie we navigate but dont invoke a client action,
// until a certain condition occurs, specified by a node matching a
// filter. This is what used to be called listening in extendio. We
// could call these conditions, waypoints. We could wake or sleep type
// waypoints.
//
// Hibernation depends on filtering.
//

// subscribe to options.before
func RestoreOptions() {
	// called by resume to load options from json file and
	// setup binder to reflect this
}

type (
	nextFn func(node *core.Node, inspection override.Inspection) (bool, error)

	state struct {
		next nextFn
	}
	hibernateStates map[enums.Hibernation]state

	triggers struct {
		wake  core.TraverseFilter
		sleep core.TraverseFilter
	}

	profile interface {
		init(controls *cycle.Controls) error
		next(node *core.Node, inspection override.Inspection) (bool, error)
	}

	common struct {
		ho       *core.HibernateOptions
		triggers triggers
		controls *cycle.Controls
	}
)

func launch(ho *core.HibernateOptions) enums.Hibernation {
	if ho.WakeAt != nil {
		return enums.HibernationPending
	}

	return enums.HibernationActive
}