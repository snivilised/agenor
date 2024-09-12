package pref

import (
	"github.com/snivilised/traverse/core"
)

// WithHibernationBehaviourExclusiveWake activates hibernation
// with a wake condition. The wake condition should be defined
// using WithHibernationFilterWake.
func WithHibernationBehaviourExclusiveWake() Option {
	return func(o *Options) error {
		o.Hibernate.Behaviour.InclusiveWake = false

		return nil
	}
}

// WithHibernationBehaviourInclusiveSleep activates hibernation
// with a sleep condition. The sleep condition should be defined
// using WithHibernationFilterSleep.
func WithHibernationBehaviourInclusiveSleep() Option {
	return func(o *Options) error {
		o.Hibernate.Behaviour.InclusiveSleep = true
		return nil
	}
}

// WithHibernationFilterWake defines the wake condition
// for hibernation based traversal sessions.
func WithHibernationFilterWake(wake *core.FilterDef) Option {
	return func(o *Options) error {
		o.Hibernate.WakeAt = wake

		return nil
	}
}

// WithHibernationFilterSleep defines the sleep condition
// for hibernation based traversal sessions.
func WithHibernationFilterSleep(sleep *core.FilterDef) Option {
	return func(o *Options) error {
		o.Hibernate.SleepAt = sleep

		return nil
	}
}

// WithHibernationOptions defines options for a hibernation traversal
// session.
func WithHibernationOptions(ho *core.HibernateOptions) Option {
	return func(o *Options) error {
		o.Hibernate = *ho

		return nil
	}
}
