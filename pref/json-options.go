package pref

import (
	"github.com/snivilised/agenor/core"
)

// JSONOptions defines the JSON persist format for options.
type JSONOptions struct {
	// all fields should be flattened out here when implemented

	// Behaviours collection of behaviours that adjust the way navigation occurs,
	// that can be tweaked by the client.
	//
	Behaviours NavigationBehaviours

	// Sampling options
	//
	Sampling SamplingOptions

	// Filter
	//
	Filter FilterOptions

	// Hibernation
	//
	Hibernate core.HibernateOptions

	// Concurrency contains options relating concurrency
	//
	Concurrency ConcurrencyOptions
}
