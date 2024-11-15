package pref

import (
	"github.com/snivilised/agenor/life"
)

// WithOnAscend sets ascend handler, invoked when navigator
// traverses up a directory, ie after all children have been
// visited.
func WithOnAscend(handler life.NodeHandler) Option {
	return func(o *Options) error {
		o.Events.Ascend.On(handler)

		return nil
	}
}

// WithOnBegin sets the begin handler, invoked before the start
// of a traversal session.
func WithOnBegin(handler life.BeginHandler) Option {
	return func(o *Options) error {
		o.Events.Begin.On(handler)

		return nil
	}
}

// WithOnDescend sets the descend handler, invoked when navigator
// traverses down into a child directory.
func WithOnDescend(handler life.NodeHandler) Option {
	return func(o *Options) error {
		o.Events.Descend.On(handler)

		return nil
	}
}

// WithOnEnd sets the end handler, invoked at the end of a traversal
// session.
func WithOnEnd(handler life.EndHandler) Option {
	return func(o *Options) error {
		o.Events.End.On(handler)

		return nil
	}
}

// WithOnWake sets the wake handler, when hibernation is active
// and the wake condition has occurred, ie when a file system
// node is encountered that matches the hibernation's wake filter.
func WithOnWake(handler life.HibernateHandler) Option {
	return func(o *Options) error {
		o.Events.Wake.On(handler)

		return nil
	}
}

// WithOnSleep sets the sleep handler, when hibernation is active
// and the sleep condition has occurred, ie when a file system
// node is encountered that matches the hibernation's sleep filter.
func WithOnSleep(handler life.HibernateHandler) Option {
	return func(o *Options) error {
		o.Events.Sleep.On(handler)

		return nil
	}
}
