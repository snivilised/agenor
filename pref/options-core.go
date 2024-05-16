package pref

import (
	"github.com/snivilised/traverse/enums"
)

type CoreOptions struct {
	// Subscription defines which node types are visited
	//
	Subscription enums.Subscription

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
	Hibernate HibernateOptions
}

func WithSubscription(subscription enums.Subscription) Option {
	return func(o *Options) error {
		o.Core.Subscription = subscription

		return nil
	}
}
