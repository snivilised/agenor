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
	badge = "option-requester"
)

func init() {
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

		binder *Binder
	}

	// Option functional traverse options
	Option func(o *Options) error
)

func Request(binder *Binder, opts ...Option) *Options {
	o := DefaultOptions()
	o.Events.Bind(&binder.Notification)

	apply(o, opts...)

	o.binder = binder

	return o
}

type LoadInfo struct {
	O      *Options
	WakeAt string
}

func Load(binder *Binder, from string, opts ...Option) (*LoadInfo, error) {
	o := DefaultOptions()
	// do load
	_ = from
	o.binder = binder

	apply(o, opts...)

	return &LoadInfo{
		O:      o,
		WakeAt: "tbd",
	}, nil
}

func apply(o *Options, opts ...Option) {
	for _, option := range opts {
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
