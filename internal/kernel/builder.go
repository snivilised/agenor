package kernel

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/pref"
)

type (
	// Artefacts is the result of building a navigator.
	Artefacts struct {
		// Kontroller is the kernel controller.
		Kontroller enclave.KernelController
		// Mediator is the mediator.
		Mediator enclave.Mediator
		// Resources is the resources.
		Resources *enclave.Resources
		// IfResult is the result of the navigation.
		IfResult core.Completion
		// Error is the error.
		Error error
	}

	// Inception is the input to the builder.
	Inception struct {
		// Facade is the facade.
		Facade pref.Facade
		// Subscription is the subscription.
		Subscription enums.Subscription
		// Harvest is the harvest.
		Harvest enclave.OptionHarvest
		// Resources is the resources.
		Resources *enclave.Resources
	}

	// NavigatorBuilder is the interface for building a navigator.
	NavigatorBuilder interface {
		Build(inception *Inception) *Artefacts
	}

	// Builder is a function that builds a navigator.
	Builder func(inception *Inception) *Artefacts
)

// Build builds a navigator.
func (fn Builder) Build(inception *Inception) *Artefacts {
	return fn(inception)
}

// NavigationTree returns the navigation tree.
func (i *Inception) NavigationTree() string {
	if using, ok := i.Facade.(*pref.Using); ok {
		return using.Tree
	}

	return i.Harvest.Loaded().State.Tree
}
