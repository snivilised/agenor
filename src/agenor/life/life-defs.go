package life

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
)

// beforeX
// afterX

// eg beforeOptions
// afterOptions

type (
	// Event is the interface for life cycle events, which can be subscribed
	// to by handlers. The type parameter F represents the type of the
	// handler function that can be subscribed to the event.
	Event[F any] interface {
		// On subscribes to a life cycle event
		On(handler F)
	}

	// SimpleHandler is a function that takes no extra custom parameters and can
	// be used by any notification with this signature.
	SimpleHandler func()

	// BeginState represents the state at the beginning of traversal, which can be
	// used by the BeginHandler to provide context about the traversal.
	BeginState struct {
		// Tree represents the tree being traversed. This can be used by the BeginHandler
		// to provide context about the traversal.
		Tree string
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
