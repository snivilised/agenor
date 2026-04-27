package command

import (
	"sync"

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/app/controller"
	"github.com/snivilised/jaywalk/src/locale"
)

const (
	defaultRunSubscribe = SubscribeFlagDefault
	defaultRunAction    = ""
	defaultRunPipeline  = ""
	defaultRunResume    = ""
)

func (b *Bootstrap) buildRunCommand(container *assist.CobraContainer) {
	runCmd := &cobra.Command{
		Use:   "run <directory>",
		Short: li18ngo.Text(locale.RunCmdShortDescTemplData{}),
		Long:  li18ngo.Text(locale.RunCmdLongDescTemplData{}),
		Args:  cobra.ExactArgs(1),
		RunE:  b.runRun,
	}

	runPs := assist.NewParamSet[RunParameterSet](runCmd)

	// --subscribe(-s): which node types to visit
	runPs.BindString(
		assist.NewFlagInfo(
			li18ngo.Text(locale.SubscribeFlagDescTemplData{}),
			"s",
			defaultRunSubscribe,
		),
		&runPs.Native.Subscribe,
	)

	// --action(-a)
	runPs.BindString(
		assist.NewFlagInfo(
			li18ngo.Text(locale.ActionFlagDescTemplData{}),
			"a",
			defaultRunAction,
		),
		&runPs.Native.Action,
	)

	// --pipeline(-p)
	runPs.BindString(
		assist.NewFlagInfo(
			li18ngo.Text(locale.PipelineFlagDescTemplData{}),
			"p",
			defaultRunPipeline,
		),
		&runPs.Native.Pipeline,
	)

	// --resume(-r)
	runPs.BindString(
		assist.NewFlagInfo(
			li18ngo.Text(locale.ResumeFlagDescTemplData{}),
			"r",
			defaultRunResume,
		),
		&runPs.Native.Resume,
	)

	// family: worker-pool [--cpu, --now]
	// Local to run, registered on the run command's own flags.
	workerPoolFam := assist.NewParamSet[store.WorkerPoolParameterSet](runCmd)
	workerPoolFam.Native.BindAll(workerPoolFam, runCmd.Flags())
	container.MustRegisterParamSet(WorkerPoolFamName, workerPoolFam)

	// poly-filter family - local to run, not inherited.
	polyFam := assist.NewParamSet[store.PolyFilterParameterSet](runCmd)
	polyFam.Native.BindAll(polyFam)

	container.MustRegisterRootedCommand(runCmd)
	container.MustRegisterParamSet(RunPsName, runPs)
	container.MustRegisterParamSet(PolyFamName+"-run", polyFam)

	b.runPs = runPs
	b.runPolyFam = polyFam
	b.workerPoolFam = workerPoolFam
}

// runRun is the RunE handler for the run command. It parses flags,
// constructs the agenor.Hare scenario, and delegates to the coordinator.
// The WaitGroup is owned here — the adapter created it and waits on it
// after the coordinator returns.
func (b *Bootstrap) runRun(cmd *cobra.Command, args []string) error {
	subscription, err := ResolveSubscription(b.runPs.Native.Subscribe)
	if err != nil {
		return err
	}

	opts := buildOptions(b.sharedFamilies())

	// Worker pool options are appended here — they come from run-specific
	// flags and the coordinator has no knowledge of them.
	if b.workerPoolFam.Native.CPU {
		opts = append(opts, agenor.WithCPU())
	} else if n := b.workerPoolFam.Native.NoWorkers; n > 0 {
		opts = append(opts, agenor.WithNoW(uint(n)))
	}

	isPrime := b.runPs.Native.Resume == ""
	wg := &sync.WaitGroup{}

	base := controller.Request{
		Subscription: subscription,
		Options:      opts,
		ActionName:   b.runPs.Native.Action,
		PipelineName: b.runPs.Native.Pipeline,
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
		strategy, e := resolveResumeStrategy(b.runPs.Native.Resume)
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
