package command

import (
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/app/controller"
	"github.com/snivilised/jaywalk/src/locale"
)

func (b *Bootstrap) buildWalkCommand(container *assist.CobraContainer) {
	walkCmd := &cobra.Command{
		Use:   "walk <directory>",
		Short: li18ngo.Text(locale.WalkCmdShortDescTemplData{}),
		Long:  li18ngo.Text(locale.WalkCmdLongDescTemplData{}),
		Args:  cobra.ExactArgs(1),
		RunE:  b.runWalk,
	}

	container.MustRegisterCommand("exec", walkCmd)
}

// runWalk is the RunE handler for the walk command. It reads flags from
// the nav and exec param-sets (all inherited), constructs the
// agenor.Tortoise scenario, and delegates to the coordinator.
func (b *Bootstrap) runWalk(cmd *cobra.Command, args []string) error {
	if err := requireActivator(b.navPs.Native.Action, b.navPs.Native.Pipeline); err != nil {
		return err
	}

	subscription, err := ResolveSubscription(b.navPs.Native.Subscribe)
	if err != nil {
		return err
	}

	settings := createSettings(b.navFamilies())
	settings = append(settings, pref.WithTraversalConfigurer(b.UI))
	isPrime := b.execPs.Native.Resume == ""

	base := controller.Request{
		Subscription: subscription,
		Settings:     settings,
		ActionName:   b.navPs.Native.Action,
		PipelineName: b.navPs.Native.Pipeline,
		Scenario:     agenor.Tortoise(isPrime),
		UI:           b.UI,
		GetForest:    b.options.GetForest,
	}

	if isPrime {
		return b.coord.ExecutePrime(cmd.Context(), &controller.PrimeRequest{
			Request: base,
			Tree:    args[0],
		})
	}

	strategy, err := resolveResumeStrategy(b.execPs.Native.Resume)
	if err != nil {
		return err
	}

	return b.coord.ExecuteResume(cmd.Context(), &controller.ResumeRequest{
		Request:  base,
		Strategy: strategy,
	})
}
