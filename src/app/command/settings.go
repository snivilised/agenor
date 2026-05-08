package command

import (
	"fmt"
	"strings"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/app/controller"
	"github.com/snivilised/jaywalk/src/locale"
)

func createTraversalSettingsIntent(families NavFamilies) controller.TraversalSettingsIntent {
	return controller.TraversalSettingsIntent{
		NoRecurse:     families.Cascade.Native.NoRecurse,
		Depth:         families.Cascade.Native.Depth,
		IsSampling:    families.Sampling.Native.IsSampling,
		NoFiles:       families.Sampling.Native.NoFiles,
		NoDirectories: families.Sampling.Native.NoDirectories,
		Filter: controller.FilterIntent{
			FilesExGlob:      families.PolyFam.Native.FilesExGlob,
			FilesRegEx:       families.PolyFam.Native.FilesRegEx,
			DirectoriesGlob:  families.PolyFam.Native.DirectoriesGlob,
			DirectoriesRegEx: families.PolyFam.Native.DirectoriesRegEx,
		},
	}
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
