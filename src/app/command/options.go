package command

import (
	"fmt"
	"strings"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/locale"
)

// buildOptions translates shared flag values into agenor option functions.
// Shared between walk and run commands.
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

// resolveResumeStrategy maps the --resume flag string to the agenor constant.
func resolveResumeStrategy(resume string) (agenor.ResumeStrategy, error) {
	switch resume {
	case ResumeStrategySpawn:
		return agenor.ResumeStrategySpawn, nil
	case ResumeStrategyFastward:
		return agenor.ResumeStrategyFastward, nil
	default:
		return 0, locale.NewInvalidResumeValueError(
			resume,
			fmt.Sprintf("'%s'", strings.Join([]string{
				ResumeStrategySpawn, ResumeStrategyFastward,
			}, ", ")),
		)
	}
}
