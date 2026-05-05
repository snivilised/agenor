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

func (b *Bootstrap) buildQueryCommand(container *assist.CobraContainer) {
	queryCmd := &cobra.Command{
		Use:   "query <directory>",
		Short: li18ngo.Text(locale.QueryCmdShortDescTemplData{}),
		Long:  li18ngo.Text(locale.QueryCmdLongDescTemplData{}),
		Args:  cobra.ExactArgs(1),
		RunE:  b.runQuery,
	}

	container.MustRegisterCommand("nav", queryCmd)
}

// runQuery is the RunE handler for the query command. It reads flags from
// the nav-level param-sets (all inherited), constructs the agenor.Goldfish
// scenario, and delegates to the coordinator.
//
// When --action or --pipeline is supplied, query displays which activator
// would be invoked per node without executing it. When neither is supplied,
// query displays the nodes that would be visited.
func (b *Bootstrap) runQuery(cmd *cobra.Command, args []string) error {
	subscription, err := ResolveSubscription(b.navPs.Native.Subscribe)
	if err != nil {
		return err
	}

	settings := createSettings(b.navFamilies())
	settings = append(settings, pref.WithTraversalConfigurer(b.UI))

	base := controller.Request{
		Subscription: subscription,
		Settings:     settings,
		ActionName:   b.navPs.Native.Action,
		PipelineName: b.navPs.Native.Pipeline,
		Scenario:     agenor.SlowPrime,
		UI:           b.UI,
	}

	return b.coord.ExecutePrime(cmd.Context(), &controller.PrimeRequest{
		Request: base,
		Tree:    args[0],
	})
}
