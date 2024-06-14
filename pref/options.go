package pref

import (
	"io/fs"
	"log/slog"
	"os"
	"runtime"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/tapable"
)

// package: pref contains user option definitions; do not use anything in kernel (cyclic)

const (
	badge = "badge: option-requester"
)

type (
	Options struct {
		Core CoreOptions

		// Sampler defines options for sampling directory entries. There are
		// multiple ways of performing sampling. The client can either:
		// A) Use one of the four predefined functions see (SamplerOptions.Fn)
		// B) Use a Custom iterator. When setting the Custom iterator properties
		//
		Sampler SamplerOptions

		// Events provides the ability to tap into life cycle events
		//
		Events cycle.Events

		// Hooks contains client hook-able functions
		//
		Hooks tapable.Hooks

		// Monitor contains externally provided logger
		//
		Monitor MonitorOptions

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

type ActiveState struct {
}

type LoadInfo struct {
	O      *Options
	State  *ActiveState
	WakeAt string
}

func Load(_ fs.FS, from string, settings ...Option) (*LoadInfo, error) {
	o := DefaultOptions()
	// do load
	_ = from
	binder := NewBinder()
	o.Events.Bind(&binder.Controls)
	o.Binder = binder

	// TODO: save any active state on the binder, eg the wake point

	err := apply(o, settings...)
	o.Binder.Loaded = &LoadInfo{
		// O:      o,
		WakeAt: "tbd",
	}

	return &LoadInfo{
		O:      o,
		WakeAt: "tbd",
	}, err
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
			Concurrency: ConcurrencyOptions{
				NoW: uint(runtime.NumCPU()),
			},
			Persist: PersistOptions{
				Format: enums.PersistJSON,
			},
		},

		Hooks: tapable.Hooks{
			FileSubPath:   tapable.NewHookCtrl[core.SubPathHook](RootParentSubPathHook),
			FolderSubPath: tapable.NewHookCtrl[core.SubPathHook](RootParentSubPathHook),
			ReadDirectory: tapable.NewHookCtrl[core.ReadDirectoryHook](DefaultReadEntriesHook),
			QueryStatus:   tapable.NewHookCtrl[core.QueryStatusHook](os.Lstat),
			Sort:          tapable.NewHookCtrl[core.SortHook](CaseInSensitiveSortHook),
		},

		Monitor: MonitorOptions{
			Log: nopLogger,
		},
	}

	return o
}
