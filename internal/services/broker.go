package services

import (
	"github.com/snivilised/agenor/internal/third/bus"
)

const (
	format = "%03d"
	// TopicInitPlugins topic used to indicate initialisation of plugins.
	TopicInitPlugins = "topic:init.plugins"

	// TopicInterceptNavigator topic used to indicate navigator can
	// creation can be intercepted.
	TopicInterceptNavigator = "topic:intercept.navigator"

	// TopicNavigationComplete topic used to indicate navigation has completed.
	TopicNavigationComplete = "topic:navigation.complete"

	// TopicOptionsAnnounce topic used to indicate options have been processed.
	TopicOptionsAnnounce = "topic:options.announce"

	// TopicOptionsBefore topic used to indicate options are about to be processed.
	TopicOptionsBefore = "topic:options.before"

	// TopicOptionsComplete topic used to indicate options have been processed.
	TopicOptionsComplete = "topic:options.complete"
)

var (
	// Broker is the broker instance
	Broker *bus.Broker
	topics = []string{
		TopicInitPlugins,
		TopicInterceptNavigator,
		TopicNavigationComplete,
		TopicOptionsAnnounce,
		TopicOptionsBefore,
		TopicOptionsComplete,
	}
)

func init() {
	Reset()
}

type (
	// InitBroker is the broker initialisation interface
	InitBroker interface {
		// Available indicates availability of the broker
		Available(b *bus.Broker)
	}
)

// Reset creates a new Broker
func Reset() *bus.Broker {
	b, err := bus.New(&bus.Sequential{
		Format: format,
	})
	if err != nil {
		panic(err)
	}

	b.RegisterTopics(topics...)

	// Access to the broker is currently un-synchronised; that is because interaction
	// with the broker is only expected to come from a single thread. However, if we
	// really wanted to grant access to it to other threads, we can define a wrapper
	// function/object around it that implements synchronisation using locks.
	//
	Broker = b

	return b
}
