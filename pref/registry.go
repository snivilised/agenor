package pref

import (
	"github.com/snivilised/traverse/cycle"
)

// this package should be internal

type (
	// Registry contains items derived from Options
	Registry struct {
		O            *Options
		Notification cycle.Controls
	}
)

func NewRegistry() *Registry {
	return &Registry{
		Notification: cycle.Controls{
			Ascend: cycle.NotificationCtrl[cycle.NodeHandler]{
				Dispatch: cycle.DescendDispatcher,
			},
			Begin: cycle.NotificationCtrl[cycle.BeginHandler]{
				Dispatch: cycle.BeginDispatcher,
			},
			Descend: cycle.NotificationCtrl[cycle.NodeHandler]{
				Dispatch: cycle.DescendDispatcher,
			},
			End: cycle.NotificationCtrl[cycle.EndHandler]{
				Dispatch: cycle.EndDispatcher,
			},
			Start: cycle.NotificationCtrl[cycle.HibernateHandler]{
				Dispatch: cycle.StartDispatcher,
			},
			Stop: cycle.NotificationCtrl[cycle.HibernateHandler]{
				Dispatch: cycle.StopDispatcher,
			},
		},
	}
}
