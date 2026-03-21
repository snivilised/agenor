package command

import (
	"sync"

	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/app/controller"
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
		Short: "run a concurrent directory tree traversal using a worker pool",
		Long: `Run traverses a directory tree concurrently via an agenor worker pool.
Flags are identical to walk, plus --cpu and --now for worker pool control.
Use --action or --pipeline to name a config-defined operation.`,
		Args: cobra.ExactArgs(1),
		RunE: b.runRun,
	}

	runPs := assist.NewParamSet[RunParameterSet](runCmd)

	// --subscribe(-s): which node types to visit
	runPs.BindString(
		assist.NewFlagInfo(
			"subscribe node types to visit: \"files\", \"dirs\" or \"all\" (default)",
			"s",
			defaultRunSubscribe,
		),
		&runPs.Native.Subscribe,
	)

	// --action(-a)
	runPs.BindString(
		assist.NewFlagInfo(
			"action name of the config-defined action to invoke for each matched node",
			"a",
			defaultRunAction,
		),
		&runPs.Native.Action,
	)

	// --pipeline(-p)
	runPs.BindString(
		assist.NewFlagInfo(
			"pipeline name of the config-defined pipeline to execute",
			"p",
			defaultRunPipeline,
		),
		&runPs.Native.Pipeline,
	)

	// --resume(-r)
	runPs.BindString(
		assist.NewFlagInfo(
			`resume strategy for an interrupted traversal: "spawn" or "fastward"`,
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
