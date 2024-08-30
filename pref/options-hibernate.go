package pref

import (
	"github.com/snivilised/traverse/core"
)

func WithHibernationBehaviourExclusiveWake() Option {
	return func(o *Options) error {
		o.Hibernate.Behaviour.InclusiveWake = false

		return nil
	}
}

func WithHibernationBehaviourInclusiveSleep() Option {
	return func(o *Options) error {
		o.Hibernate.Behaviour.InclusiveSleep = true
		return nil
	}
}

func WithHibernationFilterWake(wake *core.FilterDef) Option {
	return func(o *Options) error {
		o.Hibernate.WakeAt = wake

		return nil
	}
}

func WithHibernationFilterSleep(sleep *core.FilterDef) Option {
	return func(o *Options) error {
		o.Hibernate.SleepAt = sleep

		return nil
	}
}

func WithHibernationOptions(ho *core.HibernateOptions) Option {
	return func(o *Options) error {
		o.Hibernate = *ho

		return nil
	}
}
