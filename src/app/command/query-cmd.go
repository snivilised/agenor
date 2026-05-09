package command

import (
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/agenor"
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

	b.bindNavFlags(queryCmd, &b.query.navState)

	// query intentionally omits bindExecFlags: it is a read-only traversal
	// and --resume has no meaning here.

	container.MustRegisterParamSet(QueryNavPsName, b.query.navPs)
	container.MustRegisterParamSet(QueryPreviewFamName, b.query.previewFam)
	container.MustRegisterParamSet(QueryCascadeFamName, b.query.cascadeFam)
	container.MustRegisterParamSet(QuerySamplingFamName, b.query.samplingFam)
	container.MustRegisterParamSet(QueryPolyFamName, b.query.polyFam)

	container.MustRegisterRootedCommand(queryCmd)
}

// runQuery is the RunE handler for the query command. It reads flags from
// the query-local param-sets, constructs the agenor.Goldfish scenario,
// and delegates to the coordinator.
//
// When --action or --pipeline is supplied, query displays which activator
// would be invoked per node without executing it. When neither is supplied,
// query displays the nodes that would be visited.
func (b *Bootstrap) runQuery(cmd *cobra.Command, args []string) error {
	subscription, err := controller.ResolveSubscription(b.query.navPs.Native.Subscribe)
	if err != nil {
		return err
	}

	settings := controller.BuildTraversalSettings(
		createTraversalSettingsIntent(navFamilies(&b.query.navState)),
		b.UI,
	)

	base := controller.Request{
		Subscription: subscription,
		Settings:     settings,
		ActionName:   b.query.navPs.Native.Action,
		PipelineName: b.query.navPs.Native.Pipeline,
		Scenario:     agenor.SlowPrime,
		UI:           b.UI,
	}

	return b.coord.ExecutePrime(cmd.Context(), &controller.PrimeRequest{
		Request: base,
		Tree:    args[0],
	})
}
