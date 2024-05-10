package core

import (
	"github.com/snivilised/extendio/bus"
)

const (
	format               = "%03d"
	TopicOptionsAnnounce = "options.announce"
	TopicOptionsBefore   = "options.before"
	TopicOptionsComplete = "options.complete"
	TopicContextExpired  = "context.expired"
)

var (
	Broker *bus.Broker
	topics = []string{
		TopicOptionsAnnounce,
		TopicOptionsBefore,
		TopicOptionsComplete,
	}
)

func init() {
	b, err := bus.New(&bus.Sequential{
		Format: format,
	})

	if err != nil {
		panic(err)
	}

	b.RegisterTopics(topics...)

	Broker = b
}
