package pref

import (
	"log/slog"
	"runtime"

	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
)

// package: pref contains user option definitions; do not use anything in kernel (cyclic)

const (
	badge = "badge: option-requester"
)

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

		// Concurrency contains options relating concurrency
		//
		Concurrency ConcurrencyOptions

		Binder *Binder
	}

	// Option functional traverse options
	Option func(o *Options) error
)

func Get(settings ...Option) (o *Options, err error) {
	o = DefaultOptions()
	binder := NewBinder()
	o.Events.Bind(&binder.Controls)

	err = apply(o, settings...)
	o.Binder = binder

	return
}

type LoadInfo struct {
	// O      *Options
	WakeAt string
}

func Load(from string, settings ...Option) (o *Options, err error) {
	o = DefaultOptions()
	// do load
	_ = from
	binder := NewBinder()
	o.Events.Bind(&binder.Controls)
	o.Binder = binder

	// TODO: save any active state on the binder, eg the wake point

	err = apply(o, settings...)
	o.Binder.Loaded = &LoadInfo{
		// O:      o,
		WakeAt: "tbd",
	}

	return
}

func apply(o *Options, settings ...Option) (err error) {
	for _, option := range settings {
		err = option(o)

		if err != nil {
			return err
		}
	}

	return
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
		Concurrency: ConcurrencyOptions{
			NoW: uint(runtime.NumCPU()),
		},
	}

	return o
}
