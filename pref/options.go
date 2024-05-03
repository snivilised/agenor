package pref

import (
	"log/slog"

	"github.com/snivilised/traverse/enums"
)

// package: pref contains user option definitions; do not use anything in nav (cyclic)

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

		// Monitor contains externally provided logger
		//
		Monitor MonitorOptions
	}

	// OptionFn functional traverse options
	OptionFn func(o *Options, reg *Registry) error
)

func requestOptions(with ...OptionFn) *Options {
	o := getDetDefaultOptions()
	reg := &Registry{}

	for _, functionalOption := range with {
		// TODO: check error
		_ = functionalOption(o, reg)
	}

	return o
}

func getDetDefaultOptions() *Options {
	nopLogger := &slog.Logger{}

	return &Options{
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
}
