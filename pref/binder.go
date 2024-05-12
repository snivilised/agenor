package pref

import (
	"github.com/snivilised/traverse/cycle"
)

// this package should be internal

type (
	// Binder contains items derived from Options
	Binder struct {
		Notification cycle.Controls
	}
)

func NewBinder() *Binder {
	return &Binder{
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
