package services

import "github.com/snivilised/extendio/bus"

const (
	format               = "%03d"
	TopicContextExpired  = "context.expired"
	TopicOptionsAnnounce = "options.announce"
	TopicOptionsBefore   = "options.before"
	TopicOptionsComplete = "options.complete"
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
	b, err := bus.New(&bus.Sequential{
		Format: format,
	})

	if err != nil {
		panic(err)
	}

	b.RegisterTopics(topics...)

	Broker = b
}
