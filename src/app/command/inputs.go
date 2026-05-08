package command

import (
	"strings"

	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/locale"
	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
)

// ---------------------------------------------------------------------------
// Activator validation
// ---------------------------------------------------------------------------

// requireActivator returns an error if neither --action nor --pipeline has
// been supplied. This enforces the one-required constraint that cobra cannot
// express across inherited persistent flags.
// See https://github.com/spf13/cobra/issues/921.
func requireActivator(action, pipeline string) error {
	flags := strings.Join([]string{action, pipeline}, ", ")
	if action == "" && pipeline == "" {
		return locale.NewMarkInheritedFlagsOneRequiredError("cmd", flags)
	}

	return nil
}

// ---------------------------------------------------------------------------
// NavFamilies - shared by all navigation commands
// ---------------------------------------------------------------------------

// NavFamilies groups the mamba param-set pointers registered on the ghost
// nav command as persistent flags. All navigation commands (walk, sprint, query)
// inherit these automatically through the cobra parent/child relationship.
type NavFamilies struct {
	Preview  *assist.ParamSet[store.PreviewParameterSet]
	Cascade  *assist.ParamSet[store.CascadeParameterSet]
	Sampling *assist.ParamSet[store.SamplingParameterSet]
	PolyFam  *assist.ParamSet[store.PolyFilterParameterSet]

	// Tree is the positional directory argument.
	Tree string

	// UI is the display manager selected by --tui. All output to the
	// terminal is routed through this interface.
	UI report.Presenter

	// NavPs holds the nav-level flags (subscribe, action, pipeline)
	// inherited from the ghost nav command.
	NavPs *assist.ParamSet[NavParameterSet]

	// ExecPs holds the exec-level flags (resume) inherited from the
	// ghost exec command.
	ExecPs *assist.ParamSet[ExecParameterSet]
}

// ---------------------------------------------------------------------------
// WalkInputs - consumed by the walk RunE handler
// ---------------------------------------------------------------------------

// WalkInputs collects all flag values needed to build a walk invocation.
type WalkInputs struct {
	NavFamilies
}

// ---------------------------------------------------------------------------
// SprintInputs - consumed by the sprint RunE handler
// ---------------------------------------------------------------------------

// SprintInputs collects all flag values needed to build a sprint invocation.
type SprintInputs struct {
	NavFamilies

	// WorkerPool holds the sprint-exclusive concurrency flags (--cpu, --now).
	WorkerPool *assist.ParamSet[store.WorkerPoolParameterSet]
}

// ---------------------------------------------------------------------------
// QueryInputs - consumed by the query RunE handler
// ---------------------------------------------------------------------------

// QueryInputs collects all flag values needed to build a query invocation.
// Query is a read-only traversal - it visits nodes on a single thread and
// displays which action or pipeline would be invoked, without executing them.
type QueryInputs struct {
	NavFamilies

	// Tree is the positional directory argument.
	Tree string

	// UI is the display manager selected by --tui. All output to the
	// terminal is routed through this interface.
	UI report.Presenter

	// NavPs holds the nav-level flags (subscribe, action, pipeline)
	// inherited from the ghost nav command.
	NavPs *assist.ParamSet[NavParameterSet]
}
