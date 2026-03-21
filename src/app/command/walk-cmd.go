package command

import (
	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/app/controller"
)

const (
	defaultWalkSubscribe = SubscribeFlagDefault
	defaultWalkAction    = ""
	defaultWalkPipeline  = ""
	defaultWalkResume    = ""
)

func (b *Bootstrap) buildWalkCommand(container *assist.CobraContainer) {
	walkCmd := &cobra.Command{
		Use:   "walk <directory>",
		Short: "walk a directory tree, invoking an action or pipeline per node",
		Long: `Walk traverses a directory tree synchronously using the agenor library.
For each node that matches the active filter, the handler prints the node path.
Use --action or --pipeline to name a config-defined operation.
Use --resume to re-enter a previously interrupted traversal.`,
		Args: cobra.ExactArgs(1),
		RunE: b.runWalk,
	}

	walkPs := assist.NewParamSet[WalkParameterSet](walkCmd)

	// --subscribe(-s): which node types to visit
	walkPs.BindString(
		assist.NewFlagInfo(
			"subscribe node types to visit: \"files\", \"dirs\" or \"all\" (default)",
			"s",
			defaultWalkSubscribe,
		),
		&walkPs.Native.Subscribe,
	)

	// --action(-a): name of a config-defined action to run per node
	walkPs.BindString(
		assist.NewFlagInfo(
			"action name of the config-defined action to invoke for each matched node",
			"a",
			defaultWalkAction,
		),
		&walkPs.Native.Action,
	)

	// --pipeline(-p): name of a config-defined pipeline to execute
	walkPs.BindString(
		assist.NewFlagInfo(
			"pipeline name of the config-defined pipeline to execute",
			"p",
			defaultWalkPipeline,
		),
		&walkPs.Native.Pipeline,
	)

	// --resume(-r): resume strategy for interrupted traversals
	walkPs.BindString(
		assist.NewFlagInfo(
			`resume strategy for an interrupted traversal: "spawn" or "fastward"`,
			"r",
			defaultWalkResume,
		),
		&walkPs.Native.Resume,
	)

	// poly-filter family [--files-glob, --files-rx, --folders-glob, --folders-rx]
	// Local to walk only, not inherited by sub-commands.
	polyFam := assist.NewParamSet[store.PolyFilterParameterSet](walkCmd)
	polyFam.Native.BindAll(polyFam)

	container.MustRegisterRootedCommand(walkCmd)
	container.MustRegisterParamSet(WalkPsName, walkPs)
	container.MustRegisterParamSet(PolyFamName+"-walk", polyFam)

	b.walkPs = walkPs
	b.walkPolyFam = polyFam
}

// runWalk is the RunE handler for the walk command. It parses flags,
// constructs the agenor.Tortoise scenario, and delegates to the
// coordinator. No agenor traversal logic lives here.
func (b *Bootstrap) runWalk(cmd *cobra.Command, args []string) error {
	subscription, err := ResolveSubscription(b.walkPs.Native.Subscribe)
	if err != nil {
		return err
	}

	opts := buildOptions(b.sharedFamilies())
	isPrime := b.walkPs.Native.Resume == ""

	base := controller.Request{
		Subscription: subscription,
		Options:      opts,
		ActionName:   b.walkPs.Native.Action,
		PipelineName: b.walkPs.Native.Pipeline,
		Scenario:     agenor.Tortoise(isPrime),
		UI:           b.UI,
	}

	if isPrime {
		return b.coord.ExecutePrime(cmd.Context(), &controller.PrimeRequest{
			Request: base,
			Tree:    args[0],
		})
	}

	strategy, err := resolveResumeStrategy(b.walkPs.Native.Resume)
	if err != nil {
		return err
	}

	return b.coord.ExecuteResume(cmd.Context(), &controller.ResumeRequest{
		Request:  base,
		Strategy: strategy,
	})
}
