package pref

import (
	"log/slog"

	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
)

// package: pref contains user option definitions; do not use anything in kernel (cyclic)

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
	}

	// OptionFn functional traverse options
	OptionFn func(o *Options, reg *Registry) error
)

func RequestOptions(reg *Registry, with ...OptionFn) *Options {
	o := defaultOptions()
	o.Events.Bind(&reg.Notification)

	for _, option := range with {
		// TODO: check error
		_ = option(o, reg)
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
	}

	return o
}
