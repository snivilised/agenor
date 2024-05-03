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
}

func WithSubscription(subscription enums.Subscription) OptionFn {
	return func(o *Options, _ *Registry) error {
		o.Core.Subscription = subscription

		return nil
	}
}
