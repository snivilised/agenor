package pref

import (
	"github.com/snivilised/traverse/core"
)

type CoreOptions struct {
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
