package core

import "github.com/snivilised/traverse/enums"

// core contains universal definitions and handles cross cutting concerns
// try to keep to a minimum to reduce rippling changes

type (
	// TraverseResult
	TraverseResult interface {
		Error() error
	}

	// Client is the callback invoked for each file system node found
	// during traversal.
	Client func(node *Node) error
)

type (
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

type Using struct {
	Root         string
	Subscription enums.Subscription
	Handler      Client
}

func (u Using) Validate() error {
	if u.Root == "" {
		return UsingError{
			message: "missing root path",
		}
	}

	if u.Subscription == enums.SubscribeUndefined {
		return UsingError{
			message: "missing subscription",
		}
	}

	if u.Handler == nil {
		return UsingError{
			message: "missing handler",
		}
	}

	return nil
}

type As struct {
	Using
	From     string
	Strategy enums.ResumeStrategy
}

func (a As) Validate() error {
	if a.From == "" {
		return UsingError{
			message: "missing restore from path",
		}
	}

	if a.Strategy == enums.ResumeStrategyUndefined {
		return UsingError{
			message: "missing subscription",
		}
	}

	return a.Using.Validate()
}
