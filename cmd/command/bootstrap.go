package command

import (
	"fmt"
	"os"

	"github.com/cubiest/jibberjabber"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	macfg "github.com/snivilised/mamba/assist/cfg"
	si18n "github.com/snivilised/mamba/locale"
	"github.com/snivilised/mamba/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/language"

	"github.com/snivilised/agenor/cmd/internal/cfg"
	"github.com/snivilised/agenor/cmd/ui"
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

// ConfigInfo describes the configuration file that should be loaded,
// including its name, type, path and the viper instance to use.
type ConfigInfo struct {
	Name       string
	ConfigType string
	ConfigPath string
	Viper      macfg.ViperConfig
}

// ConfigureOptions groups configuration options that influence how
// Bootstrap initialises localisation and configuration.
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

// Bootstrap constructs the full cobra command tree and owns all
// mamba param-set registrations. It is the single entry point for
// application startup wiring.
type Bootstrap struct {
	container *assist.CobraContainer
	options   ConfigureOptions

	// Cfg is populated after configure() reads viper.
	Cfg *cfg.Config

	// UI is constructed from the --tui flag value in PersistentPreRunE
	// and injected into every command's Inputs struct.
	UI ui.Manager

	// root param-set - stashed so PersistentPreRunE and RunE can read it.
	rootPs *assist.ParamSet[RootParameterSet]

	// shared family param-sets (persistent, inherited by all sub-commands)
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
// to be executed.
func (b *Bootstrap) Root(options ...ConfigureOptionFn) *cobra.Command {
	b.prepare()

	for _, fo := range options {
		fo(&b.options)
	}

	b.configure()

	b.container = assist.NewCobraContainer(
		&cobra.Command{
			Use:     ApplicationName,
			Short:   li18ngo.Text(locale.RootCmdShortDescTemplData{}),
			Long:    li18ngo.Text(locale.RootCmdLongDescTemplData{}),
			Version: fmt.Sprintf("'%v'", Version),

			// PersistentPreRunE runs after flag parsing but before any
			// RunE handler. It resolves --tui into a UI manager so all
			// sub-commands receive a fully constructed b.UI.
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
// configure - viper + i18n (your existing logic, unchanged)
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

	// Load jay's typed config on top of viper now that the file is read.
	// GlobalViperConfig delegates to viper's global instance, so we use
	// viper.GetViper() to obtain the underlying *viper.Viper directly,
	// which allows cfg.Load to skip ReadInConfig on an already-read instance.
	b.Cfg, err = cfg.Load(cfg.LoadOptions{
		ViperInstance: viper.GetViper(),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "jay: config error: %v\n", err)
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
// buildRootCommand - root param-set and shared persistent families
// ---------------------------------------------------------------------------

func (b *Bootstrap) buildRootCommand(container *assist.CobraContainer) {
	root := container.Root()

	// root param-set: --tui (and any future root-level flags)
	b.rootPs = assist.NewParamSet[RootParameterSet](root)

	// --tui(-t) <mode>: selects the display renderer; defaults to "linear".
	// Validated eagerly so a bad value is rejected before traversal starts.
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

	// family: preview [--dry-run(D)]
	b.previewFam = assist.NewParamSet[store.PreviewParameterSet](root)
	b.previewFam.Native.BindAll(b.previewFam, root.PersistentFlags())
	container.MustRegisterParamSet(PreviewFamName, b.previewFam)

	// family: cascade [--depth, --no-recurse]
	b.cascadeFam = assist.NewParamSet[store.CascadeParameterSet](root)
	b.cascadeFam.Native.BindAll(b.cascadeFam, root.PersistentFlags())
	container.MustRegisterParamSet(CascadeFamName, b.cascadeFam)

	// family: sampling [--sample, --no-files, --no-folders, --last]
	b.samplingFam = assist.NewParamSet[store.SamplingParameterSet](root)
	b.samplingFam.Native.BindAll(b.samplingFam, root.PersistentFlags())
	container.MustRegisterParamSet(SamplingFamName, b.samplingFam)
}

// sharedFamilies is a convenience accessor for RunE closures.
func (b *Bootstrap) sharedFamilies() SharedFamilies {
	return SharedFamilies{
		Preview:  b.previewFam,
		Cascade:  b.cascadeFam,
		Sampling: b.samplingFam,
	}
}
