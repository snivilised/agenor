package hiber

import (
	"context"

	"github.com/snivilised/extendio/bus"
	"github.com/snivilised/traverse/core"
)

const (
	badge = "hibernator"
)

// hiber represents the facility to be able to start navigation in hibernated state,
// ie we navigate but dont invoke a client action, until a certain condition occurs,
// specified by a node matching a filter. This is what used to be called listening
// in extendio. We could call these conditions, waypoints. We could wake or sleep
// type waypoints
//
// Hibernation depends on filtering.
//

func init() {
	h := bus.Handler{
		Handle: func(_ context.Context, m bus.Message) {
			// The data field will contain the appropriate
			// object (represented behind an interface of some kind) that is related
			// to the topic.
			//
			_ = m.Data
		},
		Matcher: core.TopicOptionsAnnounce,
	}

	core.Broker.RegisterHandler(badge, h)
}

// subscribe to options.before
func RestoreOptions() {
	// called bny resume to load options from json file and
	// setup registry to reflect this
}
