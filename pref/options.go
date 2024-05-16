package pref

import (
	"context"
	"log/slog"
	"runtime"

	"github.com/snivilised/extendio/bus"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/services"
)

// package: pref contains user option definitions; do not use anything in kernel (cyclic)

const (
	badge = "badge: option-requester"
)

func initTbd() {
	h := bus.Handler{
		Handle: func(_ context.Context, m bus.Message) {
			_ = m.Data
		},
		Matcher: services.TopicOptionsAnnounce,
	}

	services.Broker.RegisterHandler(badge, h)
}

type (
	Options struct {
		Core CoreOptions

		// Persist contains options for persisting traverse options
		//
		Persist PersistOptions

		// Sampler defines options for sampling directory entries. There are
		// multiple ways of performing sampling. The client can either:
		// A) Use one of the four predefined functions see (SamplerOptions.Fn)
		// B) Use a Custom iterator. When setting the Custom iterator properties
		//
		Sampler SamplerOptions

		// Events provides the ability to tap into life cycle events
		//
		Events cycle.Events

		// Monitor contains externally provided logger
		//
		Monitor MonitorOptions

		// Acceleration contains options relating concurrency
		//
		Acceleration AccelerationOptions

		Binder *Binder
	}

	// Option functional traverse options
	Option func(o *Options) error
)

func Get(settings ...Option) (*Options, error) {
	o := DefaultOptions()
	binder := NewBinder()
	o.Events.Bind(&binder.Notification)

	apply(o, settings...)

	if o.Acceleration.ctx == nil {
		o.Acceleration.ctx = context.Background()
	}

	o.Binder = binder

	return o, nil
}

type LoadInfo struct {
	// O      *Options
	WakeAt string
}

func Load(from string, settings ...Option) (*Options, error) {
	o := DefaultOptions()
	// do load
	_ = from
	binder := NewBinder()
	o.Events.Bind(&binder.Notification)
	o.Binder = binder

	// TODO: save any active state on the binder, eg the wake point

	apply(o, settings...)

	if o.Acceleration.ctx == nil {
		o.Acceleration.ctx = context.Background()
	}

	o.Binder.Loaded = &LoadInfo{
		// O:      o,
		WakeAt: "tbd",
	}

	return o, nil
}

func apply(o *Options, settings ...Option) {
	for _, option := range settings {
		// TODO: check error
		_ = option(o)
	}
}

func DefaultOptions() *Options {
	nopLogger := &slog.Logger{}

	o := &Options{
		Core: CoreOptions{
			Subscription: enums.SubscribeUniversal,
			Behaviours: NavigationBehaviours{
				SubPath: SubPathBehaviour{
					KeepTrailingSep: true,
				},
				Sort: SortBehaviour{
					IsCaseSensitive:     false,
					DirectoryEntryOrder: enums.DirectoryContentsOrderFoldersFirst,
				},
				Hibernation: HibernationBehaviour{
					InclusiveStart: true,
					InclusiveStop:  false,
				},
			},
		},
		Persist: PersistOptions{
			Format: enums.PersistJSON,
		},
		Monitor: MonitorOptions{
			Log: nopLogger,
		},
		Acceleration: AccelerationOptions{
			now: runtime.NumCPU(),
		},
	}

	return o
}
