package command

import (
	"fmt"
	"strings"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/locale"
)

// createSettings translates shared flag values into agenor option functions.
// Shared between walk and run commands.
func createSettings(families NavFamilies) []pref.Option {
	var opts []pref.Option

	if families.Cascade.Native.NoRecurse {
		opts = append(opts, agenor.WithNoRecurse())
	}

	if d := families.Cascade.Native.Depth; d > 0 {
		opts = append(opts, agenor.WithDepth(d))
	}

	// TODO: implement DryRun on agenor
	// if families.Preview.Native.DryRun {
	//     opts = append(opts, agenor.WithDryRun())
	// }

	if families.Sampling.Native.IsSampling {
		opts = append(opts, agenor.WithSamplingOptions(&pref.SamplingOptions{
			NoOf: pref.EntryQuantities{
				Files:       families.Sampling.Native.NoFiles,
				Directories: families.Sampling.Native.NoDirectories,
			},
		}))
	}

	poly := families.PolyFam.Native
	hasFileGlob := poly.FilesExGlob != ""
	hasFileRegex := poly.FilesRegEx != ""
	hasFolderGlob := poly.DirectoriesGlob != ""
	hasFolderRegex := poly.DirectoriesRegEx != ""

	if hasFileGlob || hasFileRegex || hasFolderGlob || hasFolderRegex {
		fileDef := core.FilterDef{}

		if hasFileGlob {
			fileDef.Type = enums.FilterTypeGlob
			fileDef.Pattern = poly.FilesExGlob
		} else if hasFileRegex {
			fileDef.Type = enums.FilterTypeRegex
			fileDef.Pattern = poly.FilesRegEx
		}

		dirDef := core.FilterDef{}

		if hasFolderGlob {
			dirDef.Type = enums.FilterTypeGlob
			dirDef.Pattern = poly.DirectoriesGlob
		} else if hasFolderRegex {
			dirDef.Type = enums.FilterTypeRegex
			dirDef.Pattern = poly.DirectoriesRegEx
		}

		opts = append(opts, agenor.WithFilter(&pref.FilterOptions{
			Node: &core.FilterDef{
				Poly: &core.PolyFilterDef{
					File:      fileDef,
					Directory: dirDef,
				},
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
