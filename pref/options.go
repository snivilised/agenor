package pref

import (
	"io"
	"log/slog"
	"runtime"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/tapable"
)

// 📦 pkg: pref - contains user option definitions; do not use anything
// in kernel (cyclic).

const (
	badge = "badge: option-requester"
)

type (
	Options struct {
		// Behaviours collection of behaviours that adjust the way navigation occurs,
		// that can be tweaked by the client.
		//
		Behaviours NavigationBehaviours

		// Sampling options
		// There are multiple ways of performing sampling. The client can either:
		// A) Use one of the four predefined functions see (SamplerOptions.Fn)
		// B) Use a Custom iterator. When setting the Custom iterator properties
		//
		Sampling SamplingOptions

		// Filter
		//
		Filter FilterOptions

		// Hibernation
		//
		Hibernate core.HibernateOptions

		// Concurrency contains options relating concurrency
		//
		Concurrency ConcurrencyOptions

		// Events provides the ability to tap into life cycle events
		//
		Events cycle.Events

		// Hooks contains client hook-able functions
		//
		Hooks tapable.Hooks

		// Monitor contains externally provided logger
		//
		Monitor MonitorOptions

		// Defects contains error handling options
		//
		Defects DefectOptions
	}

	// Option functional traverse options
	Option func(o *Options) error
)

// IfOption enables options to be conditional. IfOption condition evaluates to true
// then the option is returned, otherwise nil.
func IfOption(condition bool, option Option) Option {
	if condition {
		return option
	}

	return nil
}

// ConditionalOption allows the delaying of creation of the option until
// the condition is known to be true. This is in contrast to If where the
// Option is pre-created, regardless of the condition.
type ConditionalOption func() Option

// IfOptionF
func IfOptionF(condition bool, option ConditionalOption) Option {
	if condition {
		return option()
	}

	return nil
}

// DefaultOptions
func DefaultOptions() *Options {
	// If a change is made to the default, make sure to make the corresponding
	// same change in internal/persist/data/test-restore.DEFAULT.json
	nopLogger := slog.New(slog.NewTextHandler(io.Discard, nil))

	o := &Options{
		Hibernate: core.HibernateOptions{
			Behaviour: core.HibernationBehaviour{
				InclusiveWake:  true,
				InclusiveSleep: false,
			},
		},
		Behaviours: NavigationBehaviours{
			SubPath: SubPathBehaviour{
				KeepTrailingSep: true,
			},
			Sort: SortBehaviour{
				IsCaseSensitive: false,
				SortFilesFirst:  false,
			},
		},
		Concurrency: ConcurrencyOptions{
			NoW: uint(runtime.NumCPU()),
		},
		Hooks: newHooks(),
		Monitor: MonitorOptions{
			Log: nopLogger,
		},

		Defects: DefectOptions{
			Fault: Accepter(DefaultFaultHandler),
			Panic: Rescuer(DefaultPanicHandler),
			Skip:  Asker(DefaultSkipHandler),
		},
	}

	return o
}

func newHooks() tapable.Hooks {
	return tapable.Hooks{
		FileSubPath: tapable.NewHookCtrl[
			core.SubPathHook, core.ChainSubPathHook, tapable.SubPathBroadcaster,
		](
			DefaultSubPathHook,
			tapable.GetSubPathBroadcaster,
			tapable.SubPathAttacher,
		),

		FolderSubPath: tapable.NewHookCtrl[
			core.SubPathHook, core.ChainSubPathHook, tapable.SubPathBroadcaster,
		](
			DefaultSubPathHook,
			tapable.GetSubPathBroadcaster,
			tapable.SubPathAttacher,
		),

		ReadDirectory: tapable.NewHookCtrl[
			core.ReadDirectoryHook, core.ChainReadDirectoryHook, tapable.ReadDirectoryBroadcaster,
		](
			DefaultReadEntriesHook,
			tapable.GetReadDirectoryBroadcaster,
			tapable.ReadDirectoryAttacher,
		),

		QueryStatus: tapable.NewHookCtrl[
			core.QueryStatusHook, core.ChainQueryStatusHook, tapable.QueryStatusBroadcaster,
		](
			DefaultQueryStatusHook,
			tapable.GetQueryStatusBroadcaster,
			tapable.QueryStatusAttacher,
		),

		Sort: tapable.NewHookCtrl[
			core.SortHook, core.ChainSortHook, tapable.SortBroadcaster,
		](
			DefaultCaseInSensitiveSortHook,
			tapable.GetSortBroadcaster,
			tapable.SortAttacher,
		),
	}
}
