package command

import (
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/locale"
	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
)

// ---------------------------------------------------------------------------
// Subscription resolution
// ---------------------------------------------------------------------------

// ResolveSubscription maps the user-supplied --subscribe string to the
// corresponding agenor enums.Subscription value. Returns an error if the
// value is not one of the three legal strings.
func ResolveSubscription(flag string) (enums.Subscription, error) {
	switch flag {
	case SubscribeFlagFiles, "":
		return enums.SubscribeFiles, nil
	case SubscribeFlagDirs:
		return enums.SubscribeDirectories, nil
	case SubscribeFlagAll:
		return enums.SubscribeUniversal, nil
	default:
		return 0, locale.ErrInvalidSubscription
	}
}

// ---------------------------------------------------------------------------
// NavFamilies - shared by all navigation commands
// ---------------------------------------------------------------------------

// NavFamilies groups the mamba param-set pointers registered on the ghost
// nav command as persistent flags. All navigation commands (walk, run, query)
// inherit these automatically through the cobra parent/child relationship.
type NavFamilies struct {
	Preview  *assist.ParamSet[store.PreviewParameterSet]
	Cascade  *assist.ParamSet[store.CascadeParameterSet]
	Sampling *assist.ParamSet[store.SamplingParameterSet]
	PolyFam  *assist.ParamSet[store.PolyFilterParameterSet]
}

// ---------------------------------------------------------------------------
// WalkInputs - consumed by the walk RunE handler
// ---------------------------------------------------------------------------

// WalkInputs collects all flag values needed to build a walk invocation.
type WalkInputs struct {
	NavFamilies

	// Tree is the positional directory argument.
	Tree string

	// UI is the display manager selected by --tui. All output to the
	// terminal is routed through this interface.
	UI report.Presenter

	// ParamSet holds the nav-level flags (subscribe, action, pipeline,
	// resume) inherited from the ghost nav command.
	ParamSet *assist.ParamSet[NavParameterSet]
}

// ---------------------------------------------------------------------------
// RunInputs - consumed by the run RunE handler
// ---------------------------------------------------------------------------

// RunInputs collects all flag values needed to build a run invocation.
type RunInputs struct {
	NavFamilies

	// Tree is the positional directory argument.
	Tree string

	// UI is the display manager selected by --tui. All output to the
	// terminal is routed through this interface.
	UI report.Presenter

	// ParamSet holds the nav-level flags (subscribe, action, pipeline,
	// resume) inherited from the ghost nav command.
	ParamSet *assist.ParamSet[NavParameterSet]

	// WorkerPool holds the run-exclusive concurrency flags (--cpu, --now).
	WorkerPool *assist.ParamSet[store.WorkerPoolParameterSet]
}

// ---------------------------------------------------------------------------
// QueryInputs - consumed by the query RunE handler
// ---------------------------------------------------------------------------

// QueryInputs collects all flag values needed to build a query invocation.
// Query is a single-threaded, read-only traversal - it visits nodes but
// executes no actions or pipelines.
type QueryInputs struct {
	NavFamilies

	// Tree is the positional directory argument.
	Tree string

	// UI is the display manager selected by --tui. All output to the
	// terminal is routed through this interface.
	UI report.Presenter

	// ParamSet holds the nav-level flags (subscribe, action, pipeline,
	// resume) inherited from the ghost nav command.
	ParamSet *assist.ParamSet[NavParameterSet]
}
