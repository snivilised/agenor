package command

import (
	"fmt"

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/locale"
)

const (
	defaultNavSubscribe = SubscribeFlagDefault
	defaultNavAction    = ""
	defaultNavPipeline  = ""
	defaultNavResume    = ""
)

// buildNavCommand constructs the ghost nav command. It is never intended
// to be invoked directly by the user - its sole purpose is to own the
// persistent flag families shared by all navigation commands (walk, run,
// query). By placing these families here rather than on root, utility
// commands (verify, theme) are kept clean of navigation flags.
//
// nav is registered as a rooted command (child of root) so that cobra's
// persistent flag propagation flows: root -> nav -> walk/run/query.
func (b *Bootstrap) buildNavCommand(container *assist.CobraContainer) {
	root := container.Root()

	navCmd := &cobra.Command{
		Use:    "nav",
		Hidden: true,
		Short:  "nav is a ghost command",
		Long:   "nav is a ghost command",

		// nav has no RunE. If the user accidentally types "jay nav", Run
		// prints a friendly explanation then shows root help so they can
		// find what they actually need.
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println(
				li18ngo.Text(locale.NewCommandIsNotUserInvocablePromptTemplData("nav")),
			)
			_ = root.Help()
		},
	}

	// nav-ps: subscribe, action, pipeline, resume.
	b.navPs = assist.NewParamSet[NavParameterSet](navCmd)

	// --subscribe(-s): which node types to visit.
	b.navPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			li18ngo.Text(locale.SubscribeFlagDescTemplData{}),
			"s",
			defaultNavSubscribe,
			navCmd.PersistentFlags(),
		),
		&b.navPs.Native.Subscribe,
	)

	// --action(-a): name of a config-defined action to run per node.
	b.navPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			li18ngo.Text(locale.ActionFlagDescTemplData{}),
			"a",
			defaultNavAction,
			navCmd.PersistentFlags(),
		),
		&b.navPs.Native.Action,
	)

	// --pipeline(-p): name of a config-defined pipeline to execute.
	b.navPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			li18ngo.Text(locale.PipelineFlagDescTemplData{}),
			"p",
			defaultNavPipeline,
			navCmd.PersistentFlags(),
		),
		&b.navPs.Native.Pipeline,
	)

	// --resume(-r): resume strategy for interrupted traversals.
	b.navPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			li18ngo.Text(locale.ResumeFlagDescTemplData{}),
			"r",
			defaultNavResume,
			navCmd.PersistentFlags(),
		),
		&b.navPs.Native.Resume,
	)

	// family: preview [--dry-run]
	b.previewFam = assist.NewParamSet[store.PreviewParameterSet](navCmd)
	b.previewFam.Native.BindAll(b.previewFam, navCmd.PersistentFlags())

	// family: cascade [--depth, --no-recurse]
	b.cascadeFam = assist.NewParamSet[store.CascadeParameterSet](navCmd)
	b.cascadeFam.Native.BindAll(b.cascadeFam, navCmd.PersistentFlags())

	// family: sampling [--sample, --num-files, --num-folders, --last]
	b.samplingFam = assist.NewParamSet[store.SamplingParameterSet](navCmd)
	b.samplingFam.Native.BindAll(b.samplingFam, navCmd.PersistentFlags())

	// family: poly-filter [--files-glob, --file-regex, --folders-glob, --folders-regex]
	b.polyFam = assist.NewParamSet[store.PolyFilterParameterSet](navCmd)
	b.polyFam.Native.BindAll(b.polyFam, navCmd.PersistentFlags())

	container.MustRegisterRootedCommand(navCmd)
	container.MustRegisterParamSet(NavPsName, b.navPs)
	container.MustRegisterParamSet(PreviewFamName, b.previewFam)
	container.MustRegisterParamSet(CascadeFamName, b.cascadeFam)
	container.MustRegisterParamSet(SamplingFamName, b.samplingFam)
	container.MustRegisterParamSet(PolyFamName, b.polyFam)
}
