package pref

import (
	"context"
	"log/slog"
	"runtime"

	"github.com/snivilised/extendio/bus"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
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
		Matcher: core.TopicOptionsAnnounce,
	}

	core.Broker.RegisterHandler(badge, h)
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
	}

	// Option functional traverse options
	Option func(o *Options) error
)

func RequestOptions(reg *Registry, with ...Option) *Options {
	o := defaultOptions()
	o.Events.Bind(&reg.Notification)

	for _, option := range with {
		// TODO: check error
		_ = option(o)
	}

	reg.O = o

	return o
}

func defaultOptions() *Options {
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
