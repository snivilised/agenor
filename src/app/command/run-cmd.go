package command

import (
	"context"
	"fmt"
	"sync"

	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
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
	//
	runPs.BindString(
		assist.NewFlagInfo(
			"subscribe node types to visit: \"files\", \"dirs\" or \"all\" (default)",
			"s",
			defaultRunSubscribe,
		),
		&runPs.Native.Subscribe,
	)

	// --action(-a)
	//
	runPs.BindString(
		assist.NewFlagInfo(
			"action name of the config-defined action to invoke for each matched node",
			"a",
			defaultRunAction,
		),
		&runPs.Native.Action,
	)

	// --pipeline(-p)
	//
	runPs.BindString(
		assist.NewFlagInfo(
			"pipeline name of the config-defined pipeline to execute",
			"p",
			defaultRunPipeline,
		),
		&runPs.Native.Pipeline,
	)

	// --resume(-r)
	//
	runPs.BindString(
		assist.NewFlagInfo(
			`resume strategy for an interrupted traversal: "spawn" or "fastward"`,
			"r",
			defaultRunResume,
		),
		&runPs.Native.Resume,
	)

	// family: worker-pool [--cpu(C), --now]
	// run-only, registered on the run command's local flags.
	//
	workerPoolFam := assist.NewParamSet[store.WorkerPoolParameterSet](runCmd)
	workerPoolFam.Native.BindAll(workerPoolFam, runCmd.Flags())
	container.MustRegisterParamSet(WorkerPoolFamName, workerPoolFam)

	// poly-filter family - local to run, not inherited.
	//
	polyFam := assist.NewParamSet[store.PolyFilterParameterSet](runCmd)
	polyFam.Native.BindAll(polyFam)

	container.MustRegisterRootedCommand(runCmd)
	container.MustRegisterParamSet(RunPsName, runPs)
	container.MustRegisterParamSet(PolyFamName+"-run", polyFam)

	b.runPs = runPs
	b.runPolyFam = polyFam
	b.workerPoolFam = workerPoolFam
}

// runRun is the RunE handler for the run command.
func (b *Bootstrap) runRun(cmd *cobra.Command, args []string) error {
	inputs := &RunInputs{
		Tree:           args[0],
		UI:             b.UI,
		ParamSet:       b.runPs,
		PolyFam:        b.runPolyFam,
		SharedFamilies: b.sharedFamilies(),
		WorkerPool:     b.workerPoolFam,
	}

	return executeRun(cmd.Context(), inputs)
}

// executeRun builds and runs an agenor concurrent traversal using agenor.Hare,
// which supports both prime and resume modes with a worker pool.
func executeRun(ctx context.Context, inputs *RunInputs) error {
	subscription, err := ResolveSubscription(inputs.ParamSet.Native.Subscribe)
	if err != nil {
		return err
	}

	isPrime := inputs.ParamSet.Native.Resume == ""
	opts := buildOptions(inputs.SharedFamilies)

	if inputs.WorkerPool.Native.CPU {
		opts = append(opts, agenor.WithCPU())
	} else if n := inputs.WorkerPool.Native.NoWorkers; n > 0 {
		opts = append(opts, agenor.WithNoW(uint(n)))
	}

	var facade pref.Facade

	if isPrime {
		facade = &pref.Using{
			Subscription: subscription,
			Head: pref.Head{
				Handler: func(servant agenor.Servant) error {
					return inputs.UI.OnNode(servant.Node())
				},
			},
			Tree: inputs.Tree,
		}
	} else {
		strategy, e := resolveResumeStrategy(inputs.ParamSet.Native.Resume)
		if e != nil {
			return e
		}

		facade = &pref.Relic{
			Head: pref.Head{
				Handler: func(servant agenor.Servant) error {
					return inputs.UI.OnNode(servant.Node())
				},
			},
			Strategy: strategy,
		}
	}

	wg := sync.WaitGroup{}

	result, err := agenor.Hare(isPrime, &wg)(facade, opts...).Navigate(ctx)
	wg.Wait()

	if err != nil {
		inputs.UI.Error(fmt.Sprintf("run failed: %v", err))
		return err
	}

	inputs.UI.Info(fmt.Sprintf(
		"run complete: %d files, %d dirs visited",
		result.Metrics().Count(enums.MetricNoFilesInvoked),
		result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
	))

	return nil
}
