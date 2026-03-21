package command

import (
	"context"
	"fmt"

	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
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
	//
	walkPs.BindString(
		assist.NewFlagInfo(
			"subscribe node types to visit: \"files\", \"dirs\" or \"all\" (default)",
			"s",
			defaultWalkSubscribe,
		),
		&walkPs.Native.Subscribe,
	)

	// --action(-a): name of a config-defined action to run per node
	//
	walkPs.BindString(
		assist.NewFlagInfo(
			"action name of the config-defined action to invoke for each matched node",
			"a",
			defaultWalkAction,
		),
		&walkPs.Native.Action,
	)

	// --pipeline(-p): name of a config-defined pipeline to execute
	//
	walkPs.BindString(
		assist.NewFlagInfo(
			"pipeline name of the config-defined pipeline to execute",
			"p",
			defaultWalkPipeline,
		),
		&walkPs.Native.Pipeline,
	)

	// --resume(-r): resume strategy for interrupted traversals
	//
	walkPs.BindString(
		assist.NewFlagInfo(
			`resume strategy for an interrupted traversal: "spawn" or "fastward"`,
			"r",
			defaultWalkResume,
		),
		&walkPs.Native.Resume,
	)

	// poly-filter family [--files-glob(b), --files-rx, --folders-glob, --folders-rx]
	// Not passed to PersistentFlags so it is local to walk only.
	//
	polyFam := assist.NewParamSet[store.PolyFilterParameterSet](walkCmd)
	polyFam.Native.BindAll(polyFam)

	container.MustRegisterRootedCommand(walkCmd)
	container.MustRegisterParamSet(WalkPsName, walkPs)
	container.MustRegisterParamSet(PolyFamName+"-walk", polyFam)

	b.walkPs = walkPs
	b.walkPolyFam = polyFam
}

// runWalk is the RunE handler for the walk command.
func (b *Bootstrap) runWalk(cmd *cobra.Command, args []string) error {
	inputs := &WalkInputs{
		Tree:           args[0],
		UI:             b.UI,
		ParamSet:       b.walkPs,
		PolyFam:        b.walkPolyFam,
		SharedFamilies: b.sharedFamilies(),
	}

	return executeWalk(cmd.Context(), inputs)
}

// executeWalk builds and runs an agenor walk traversal using agenor.Tortoise,
// which supports both prime and resume modes for synchronous walking.
func executeWalk(ctx context.Context, inputs *WalkInputs) error {
	subscription, err := ResolveSubscription(inputs.ParamSet.Native.Subscribe)
	if err != nil {
		return err
	}

	isPrime := inputs.ParamSet.Native.Resume == ""
	opts := buildOptions(inputs.SharedFamilies)

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

	result, err := agenor.Tortoise(isPrime)(facade, opts...).Navigate(ctx)
	if err != nil {
		inputs.UI.Error(fmt.Sprintf("walk failed: %v", err))
		return err
	}

	inputs.UI.Info(fmt.Sprintf(
		"walk complete: %d files, %d dirs visited",
		result.Metrics().Count(enums.MetricNoFilesInvoked),
		result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
	))

	return nil
}

// buildOptions translates shared flag values into agenor option functions.
// Shared between walk and run to avoid duplication.
func buildOptions(families SharedFamilies) []pref.Option {
	var opts []pref.Option

	if families.Cascade.Native.NoRecurse {
		opts = append(opts, agenor.WithNoRecurse())
	}

	if d := families.Cascade.Native.Depth; d > 0 {
		opts = append(opts, agenor.WithDepth(d))
	}

	// TODO: implement DryRun on agenor
	// if families.Preview.Native.DryRun {
	// 	opts = append(opts, agenor.WithDryRun())
	// }

	if families.Sampling.Native.IsSampling {
		opts = append(opts, agenor.WithSamplingOptions(&pref.SamplingOptions{
			NoOf: pref.EntryQuantities{
				Files:       families.Sampling.Native.NoFiles,
				Directories: families.Sampling.Native.NoFolders,
			},
		}))
	}

	return opts
}

// resolveResumeStrategy maps the --resume string to the agenor constant.
func resolveResumeStrategy(resume string) (agenor.ResumeStrategy, error) {
	switch resume {
	case ResumeStrategySpawn:
		return agenor.ResumeStrategySpawn, nil
	case ResumeStrategyFastward:
		return agenor.ResumeStrategyFastward, nil
	default:
		return 0, fmt.Errorf(
			"invalid --resume value %q: must be %q or %q",
			resume, ResumeStrategySpawn, ResumeStrategyFastward,
		)
	}
}
