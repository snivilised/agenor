package command

import (
	"sync"

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

	// query is a read-only traversal - it visits nodes on a single thread
	// but executes no actions or pipelines. The --action and --pipeline
	// flags are inherited from nav but ignored at the controller level.

	container.MustRegisterCommand("nav", queryCmd)
}

// runQuery is the RunE handler for the query command. It reads flags from
// the nav-level param-sets (all inherited), constructs the agenor.Goldfish
// scenario, and delegates to the coordinator.
func (b *Bootstrap) runQuery(cmd *cobra.Command, args []string) error {
	subscription, err := ResolveSubscription(b.navPs.Native.Subscribe)
	if err != nil {
		return err
	}

	opts := buildOptions(b.navFamilies())

	const isWalk = true
	wg := &sync.WaitGroup{}

	base := controller.Request{
		Subscription: subscription,
		Options:      opts,
		Scenario:     agenor.Goldfish(isWalk, wg),
		UI:           b.UI,
	}

	return b.coord.ExecutePrime(cmd.Context(), &controller.PrimeRequest{
		Request: base,
		Tree:    args[0],
	})
}
