package command

import (
	"sync"

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/app/controller"
	"github.com/snivilised/jaywalk/src/locale"
)

func (b *Bootstrap) buildRunCommand(container *assist.CobraContainer) {
	runCmd := &cobra.Command{
		Use:   "run <directory>",
		Short: li18ngo.Text(locale.RunCmdShortDescTemplData{}),
		Long:  li18ngo.Text(locale.RunCmdLongDescTemplData{}),
		Args:  cobra.ExactArgs(1),
		RunE:  b.runRun,
	}

	// family: worker-pool [--cpu, --now]
	b.workerPoolFam = assist.NewParamSet[store.WorkerPoolParameterSet](runCmd)
	b.workerPoolFam.Native.BindAll(b.workerPoolFam, runCmd.Flags())
	container.MustRegisterParamSet(WorkerPoolFamName, b.workerPoolFam)

	container.MustRegisterCommand("exec", runCmd)
}

// runRun is the RunE handler for the run command. It reads flags from
// the nav and exec param-sets (all inherited) plus the run-exclusive
// worker-pool family, constructs the agenor.Hare scenario, and delegates
// to the coordinator. The WaitGroup is owned here - the adapter created
// it and waits on it after the coordinator returns.
func (b *Bootstrap) runRun(cmd *cobra.Command, args []string) error {
	if err := requireActivator(b.navPs.Native.Action, b.navPs.Native.Pipeline); err != nil {
		return err
	}

	subscription, err := ResolveSubscription(b.navPs.Native.Subscribe)
	if err != nil {
		return err
	}

	settings := createSettings(b.navFamilies())
	settings = append(settings, pref.WithTraversalConfigurer(b.UI))

	if b.workerPoolFam.Native.CPU {
		settings = append(settings, agenor.WithCPU())
	} else if n := b.workerPoolFam.Native.NoWorkers; n > 0 {
		settings = append(settings, agenor.WithNoW(uint(n)))
	}

	isPrime := b.execPs.Native.Resume == ""
	wg := &sync.WaitGroup{}

	base := controller.Request{
		Subscription: subscription,
		Settings:     settings,
		ActionName:   b.navPs.Native.Action,
		PipelineName: b.navPs.Native.Pipeline,
		Scenario:     agenor.Hare(isPrime, wg),
		UI:           b.UI,
	}

	var execErr error

	if isPrime {
		execErr = b.coord.ExecutePrime(cmd.Context(), &controller.PrimeRequest{
			Request: base,
			Tree:    args[0],
		})
	} else {
		strategy, e := resolveResumeStrategy(b.execPs.Native.Resume)
		if e != nil {
			return e
		}

		execErr = b.coord.ExecuteResume(cmd.Context(), &controller.ResumeRequest{
			Request:  base,
			Strategy: strategy,
		})
	}

	wg.Wait()

	return execErr
}
