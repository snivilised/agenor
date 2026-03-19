package command

import (
	"fmt"

	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"

	"github.com/snivilised/agenor/cmd/ui"
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
		return 0, fmt.Errorf(
			"invalid --subscribe value %q: must be %q, %q or %q",
			flag, SubscribeFlagFiles, SubscribeFlagDirs, SubscribeFlagAll,
		)
	}
}

// ---------------------------------------------------------------------------
// Shared families bundle - embedded in both WalkInputs and RunInputs
// ---------------------------------------------------------------------------

// SharedFamilies groups the mamba param-set pointers registered on the
// root command as persistent flags, inherited by all sub-commands.
// Note: CliInteractionParameterSet is NOT included here because --tui is
// a string flag on RootParameterSet, allowing the user to select a named
// display mode rather than a simple on/off boolean.
type SharedFamilies struct {
	Preview  *assist.ParamSet[store.PreviewParameterSet]
	Cascade  *assist.ParamSet[store.CascadeParameterSet]
	Sampling *assist.ParamSet[store.SamplingParameterSet]
}

// ---------------------------------------------------------------------------
// WalkInputs - consumed by the walk RunE handler
// ---------------------------------------------------------------------------

// WalkInputs collects all flag values needed to build a Walk invocation.
type WalkInputs struct {
	SharedFamilies

	// Tree is the positional directory argument.
	Tree string

	// UI is the display manager selected by --tui. All output to the
	// terminal is routed through this interface.
	UI ui.Manager

	// Jay-specific flags
	ParamSet *assist.ParamSet[WalkParameterSet]

	// Per-command filter family (not inherited from root)
	PolyFam *assist.ParamSet[store.PolyFilterParameterSet]
}

// ---------------------------------------------------------------------------
// RunInputs - consumed by the run RunE handler
// ---------------------------------------------------------------------------

// RunInputs collects all flag values needed to build a Run invocation.
type RunInputs struct {
	SharedFamilies

	// Tree is the positional directory argument.
	Tree string

	// UI is the display manager selected by --tui. All output to the
	// terminal is routed through this interface.
	UI ui.Manager

	// Jay-specific flags
	ParamSet *assist.ParamSet[RunParameterSet]

	// Per-command filter family (not inherited from root)
	PolyFam *assist.ParamSet[store.PolyFilterParameterSet]

	// Run-only family
	WorkerPool *assist.ParamSet[store.WorkerPoolParameterSet]
}
