package command

import (
	"fmt"
	"os"

	"github.com/cubiest/jibberjabber"
	"github.com/snivilised/jaywalk/src/internal/third/lo"
	"github.com/snivilised/jaywalk/src/locale"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	macfg "github.com/snivilised/mamba/assist/cfg"
	si18n "github.com/snivilised/mamba/locale"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/language"

	bedrock "github.com/snivilised/jaywalk/src/app/bedrock"
	jac "github.com/snivilised/jaywalk/src/app/controller"
	"github.com/snivilised/jaywalk/src/app/ui"
)

// ---------------------------------------------------------------------------
// Locale detection
// ---------------------------------------------------------------------------

// LocaleDetector abstracts the detection of the user's preferred
// language as a BCP 47 language tag.
type LocaleDetector interface {
	Scan() language.Tag
}

// Jabber is a LocaleDetector implemented using jibberjabber.
type Jabber struct{}

// Scan returns the detected language tag.
func (j *Jabber) Scan() language.Tag {
	lang, _ := jibberjabber.DetectIETF()
	return language.MustParse(lang)
}

// ---------------------------------------------------------------------------
// ConfigureOptions
// ---------------------------------------------------------------------------

// ConfigInfo describes the configuration file that should be loaded.
type ConfigInfo struct {
	Name       string
	ConfigType string
	ConfigPath string
	Viper      macfg.ViperConfig
}

// ConfigureOptions groups options that influence how Bootstrap
// initialises localisation and configuration.
type ConfigureOptions struct {
	Detector LocaleDetector
	Config   ConfigInfo
}

// ConfigureOptionFn is a functional option used to modify
// ConfigureOptions before Bootstrap performs its setup.
type ConfigureOptionFn func(*ConfigureOptions)

// ---------------------------------------------------------------------------
// Bootstrap
// ---------------------------------------------------------------------------

// Bootstrap is a pure composition root. Its sole responsibility is
// wiring: it constructs the cobra command tree, registers param-sets,
// creates the ApplicationController, and connects everything together.
// It contains no business logic and no traversal decisions.
//
// Code smell checklist (should remain clean):
//   - No direct calls to agenor from Bootstrap
//   - No UI rendering logic in Bootstrap
//   - No flag interpretation beyond resolving --tui into a ui.Manager
type Bootstrap struct {
	container *assist.CobraContainer
	options   ConfigureOptions

	// Cfg is populated after configure() reads viper.
	Cfg *bedrock.Config

	// UI is resolved from --tui in PersistentPreRunE and passed
	// into requests; Bootstrap does not use it directly.
	UI ui.Manager

	// coord is the single ApplicationController instance wired at
	// startup and shared by all command handlers.
	coord *jac.Coordinator

	// root param-set
	rootPs *assist.ParamSet[RootParameterSet]

	// shared persistent families
	previewFam  *assist.ParamSet[store.PreviewParameterSet]
	cascadeFam  *assist.ParamSet[store.CascadeParameterSet]
	samplingFam *assist.ParamSet[store.SamplingParameterSet]

	// walk command
	walkPs      *assist.ParamSet[WalkParameterSet]
	walkPolyFam *assist.ParamSet[store.PolyFilterParameterSet]

	// run command
	runPs         *assist.ParamSet[RunParameterSet]
	runPolyFam    *assist.ParamSet[store.PolyFilterParameterSet]
	workerPoolFam *assist.ParamSet[store.WorkerPoolParameterSet]
}

// ---------------------------------------------------------------------------
// Root
// ---------------------------------------------------------------------------

func (b *Bootstrap) prepare() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	b.options = ConfigureOptions{
		Detector: &Jabber{},
		Config: ConfigInfo{
			Name:       ApplicationName,
			ConfigType: "yaml",
			ConfigPath: home,
			Viper:      &macfg.GlobalViperConfig{},
		},
	}
}

// Root builds the command tree and returns the root command, ready
// to be executed. This is the composition root: all wiring happens
// here and only here.
func (b *Bootstrap) Root(options ...ConfigureOptionFn) *cobra.Command {
	b.prepare()

	for _, fo := range options {
		fo(&b.options)
	}

	b.configure()

	// Construct the ApplicationController once. Command handlers receive
	// it via b.ctrl - they never construct it themselves.
	b.coord = jac.New()

	b.container = assist.NewCobraContainer(
		&cobra.Command{
			Use:     ApplicationName,
			Short:   li18ngo.Text(locale.RootCmdShortDescTemplData{}),
			Long:    li18ngo.Text(locale.RootCmdLongDescTemplData{}),
			Version: fmt.Sprintf("'%v'", Version),

			// PersistentPreRunE resolves --tui into a ui.Manager and
			// stores it on Bootstrap so command handlers can include it
			// in their requests. This is the only UI concern Bootstrap
			// is permitted to touch.
			PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
				mgr, err := ui.New(b.rootPs.Native.TUI)
				if err != nil {
					return err
				}
				b.UI = mgr
				return nil
			},

			Run: func(_ *cobra.Command, _ []string) {
				fmt.Println("=== jay ===")
			},
		},
	)

	b.buildRootCommand(b.container)
	b.buildWalkCommand(b.container)
	b.buildRunCommand(b.container)

	return b.container.Root()
}

// ---------------------------------------------------------------------------
// configure
// ---------------------------------------------------------------------------

func (b *Bootstrap) configure() {
	vc := b.options.Config.Viper
	ci := b.options.Config

	vc.SetConfigName(ci.Name)
	vc.SetConfigType(ci.ConfigType)
	vc.AddConfigPath(ci.ConfigPath)
	vc.AutomaticEnv()

	err := vc.ReadInConfig()

	handleLangSetting()

	if err != nil {
		msg := li18ngo.Text(locale.UsingConfigFileTemplData{
			ConfigFileName: viper.ConfigFileUsed(),
		})
		fmt.Fprintln(os.Stderr, msg)
	}

	b.Cfg, err = bedrock.Load(bedrock.LoadOptions{
		ViperInstance: viper.GetViper(),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func handleLangSetting() {
	tag := lo.TernaryF(viper.InConfig("lang"),
		func() language.Tag {
			lang := viper.GetString("lang")
			parsedTag, err := language.Parse(lang)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return parsedTag
		},
		func() language.Tag {
			return li18ngo.DefaultLanguage
		},
	)

	err := li18ngo.Use(func(uo *li18ngo.UseOptions) {
		uo.Tag = tag
		uo.From = li18ngo.LoadFrom{
			Sources: li18ngo.TranslationFiles{
				SourceID:            li18ngo.TranslationSource{Name: ApplicationName},
				si18n.MambaSourceID: li18ngo.TranslationSource{Name: "mamba"},
			},
		}
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// ---------------------------------------------------------------------------
// buildRootCommand
// ---------------------------------------------------------------------------

func (b *Bootstrap) buildRootCommand(container *assist.CobraContainer) {
	root := container.Root()

	b.rootPs = assist.NewParamSet[RootParameterSet](root)
	b.rootPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			`tui display mode: "linear" (default) or a named Charm-based renderer`,
			"t",
			ui.ModeDefault,
			root.PersistentFlags(),
		),
		&b.rootPs.Native.TUI,
	)
	container.MustRegisterParamSet(RootPsName, b.rootPs)

	b.previewFam = assist.NewParamSet[store.PreviewParameterSet](root)
	b.previewFam.Native.BindAll(b.previewFam, root.PersistentFlags())
	container.MustRegisterParamSet(PreviewFamName, b.previewFam)

	b.cascadeFam = assist.NewParamSet[store.CascadeParameterSet](root)
	b.cascadeFam.Native.BindAll(b.cascadeFam, root.PersistentFlags())
	container.MustRegisterParamSet(CascadeFamName, b.cascadeFam)

	b.samplingFam = assist.NewParamSet[store.SamplingParameterSet](root)
	b.samplingFam.Native.BindAll(b.samplingFam, root.PersistentFlags())
	container.MustRegisterParamSet(SamplingFamName, b.samplingFam)
}

// sharedFamilies is a convenience accessor used by runWalk and runRun.
func (b *Bootstrap) sharedFamilies() SharedFamilies {
	return SharedFamilies{
		Preview:  b.previewFam,
		Cascade:  b.cascadeFam,
		Sampling: b.samplingFam,
	}
}
