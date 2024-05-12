package services

import "github.com/snivilised/extendio/bus"

const (
	format               = "%03d"
	TopicContextExpired  = "context.expired"
	TopicOptionsAnnounce = "options.announce"
	TopicOptionsBefore   = "options.before"
	TopicOptionsComplete = "options.complete"
	TopicTraverseResult  = "traverse.result"
)

var (
	Broker *bus.Broker
	topics = []string{
		TopicContextExpired,
		TopicOptionsAnnounce,
		TopicOptionsBefore,
		TopicOptionsComplete,
	}
)

func init() {
	Reset()
}

func Reset() {
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
}
