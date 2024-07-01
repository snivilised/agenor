package core

import (
	"time"

	"github.com/snivilised/traverse/measure"
)

// core contains universal definitions and handles user facing cross
// cutting concerns try to keep to a minimum to reduce rippling changes.

type (
	// ResultCompletion used to determine if the result really represents
	// final navigation completion.
	ResultCompletion interface {
		IsComplete() bool
	}

	Completion func() bool

	// Session represents a traversal session and keeps tracks of
	// timing.
	Session interface {
		ResultCompletion
		StartedAt() time.Time
		Elapsed() time.Duration
	}

	// TraverseResult
	TraverseResult interface {
		Metrics() measure.Reporter
		Session() Session
		Error() error
	}

	// Client is the callback invoked for each file system node found
	// during traversal.
	Client func(node *Node) error

	// SimpleHandler is a function that takes no parameters and can
	// be used by any notification with this signature.
	SimpleHandler func()

	// BeginHandler invoked before traversal begins
	BeginHandler func(root string)

	// EndHandler invoked at the end of traversal
	EndHandler func(result TraverseResult)

	// HibernateHandler is a generic handler that is used by hibernation
	// to indicate wake or sleep.
	HibernateHandler func(description string)

	// NodeHandler is a generic handler that is for any notification that contains
	// the traversal node, such as directory ascend or descend.
	NodeHandler func(node *Node)
)

func (fn Completion) IsComplete() bool {
	return fn()
}
