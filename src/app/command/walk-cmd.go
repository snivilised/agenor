package command

import (
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/app/controller"
	"github.com/snivilised/jaywalk/src/locale"
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
		Short: li18ngo.Text(locale.WalkCmdShortDescTemplData{}),
		Long:  li18ngo.Text(locale.WalkCmdLongDescTemplData{}),
		Args:  cobra.ExactArgs(1),
		RunE:  b.runWalk,
	}

	walkPs := assist.NewParamSet[WalkParameterSet](walkCmd)

	// --subscribe(-s): which node types to visit
	walkPs.BindString(
		assist.NewFlagInfo(
			li18ngo.Text(locale.SubscribeFlagDescTemplData{}),
			"s",
			defaultWalkSubscribe,
		),
		&walkPs.Native.Subscribe,
	)

	// --action(-a): name of a config-defined action to run per node
	walkPs.BindString(
		assist.NewFlagInfo(
			li18ngo.Text(locale.ActionFlagDescTemplData{}),
			"a",
			defaultWalkAction,
		),
		&walkPs.Native.Action,
	)

	// --pipeline(-p): name of a config-defined pipeline to execute
	walkPs.BindString(
		assist.NewFlagInfo(
			li18ngo.Text(locale.PipelineFlagDescTemplData{}),
			"p",
			defaultWalkPipeline,
		),
		&walkPs.Native.Pipeline,
	)

	// --resume(-r): resume strategy for interrupted traversals
	walkPs.BindString(
		assist.NewFlagInfo(
			li18ngo.Text(locale.ResumeFlagDescTemplData{}),
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
