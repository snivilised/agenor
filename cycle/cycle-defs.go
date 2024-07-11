package cycle

import (
	"github.com/snivilised/traverse/core"
)

// cycle represents life cycle events; can't use prefs

// beforeX
// afterX

// eg beforeOptions
// afterOptions

type (
	Event[F any] interface {
		// On subscribes to a life cycle event
		On(handler F)
	}

	// SimpleHandler is a function that takes no extra custom parameters and can
	// be used by any notification with this signature.
	SimpleHandler func()

	BeginState struct {
		Root string
	}

	// BeginHandler invoked before traversal begins
	BeginHandler func(state *BeginState)

	// EndHandler invoked at the end of traversal
	EndHandler func(result core.TraverseResult)

	// HibernateHandler is a generic handler that is used by hibernation
	// to indicate wake or sleep.
	HibernateHandler func(description string)

	// NodeHandler is a generic handler that is for any notification that contains
	// the traversal node, such as directory ascend or descend.
	NodeHandler func(node *core.Node)
)
