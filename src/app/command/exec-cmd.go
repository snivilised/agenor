package command

import (
	"fmt"

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/locale"
)

const (
	defaultExecResume = ""
)

// buildExecCommand constructs the ghost exec command. It is never intended
// to be invoked directly by the user - its sole purpose is to own the
// persistent flags that apply only to execution commands (walk, run) and
// not to query.
//
// Specifically exec owns:
//   - --resume: resume strategy, meaningless for a read-only query
//
// The constraint that at least one of --action or --pipeline must be
// supplied is enforced in runWalk and runRun rather than here, because
// cobra's MarkFlagsOneRequired cannot resolve flags inherited from parent
// commands. See https://github.com/spf13/cobra/issues/921.
//
// exec is registered as a child of nav so that cobra's persistent flag
// propagation flows: root -> nav -> exec -> walk/run
func (b *Bootstrap) buildExecCommand(container *assist.CobraContainer) {
	root := container.Root()

	execCmd := &cobra.Command{
		Use:    "exec",
		Hidden: true,
		Short:  li18ngo.Text(locale.NewCommandIsAGhostTemplData("exec")),
		Long:   li18ngo.Text(locale.NewCommandIsAGhostTemplData("exec")),

		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(
				li18ngo.Text(locale.NewCommandIsNotUserInvocablePromptTemplData("exec")),
			)
			_ = root.Help()
		},
	}

	b.execPs = assist.NewParamSet[ExecParameterSet](execCmd)

	// --resume(-r): resume strategy for interrupted traversals.
	b.execPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			li18ngo.Text(locale.ResumeFlagDescTemplData{}),
			"r",
			defaultExecResume,
			execCmd.PersistentFlags(),
		),
		&b.execPs.Native.Resume,
	)

	container.MustRegisterCommand("nav", execCmd)
	container.MustRegisterParamSet(ExecPsName, b.execPs)
}
