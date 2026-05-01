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

	"github.com/snivilised/jaywalk/src/app/bedrock"
	jac "github.com/snivilised/jaywalk/src/app/controller"
	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/app/shell"
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
// creates the Coordinator, and connects everything together.
// It contains no business logic and no traversal decisions.
//
// Command hierarchy:
//
//	root  (global persistent flags: --tui, --theme, --language)
//	  nav (ghost; nav persistent flags: all navigation families)
//	    walk
//	    run
//	    query
//	  verify
//	  theme
//
// Code smell checklist (should remain clean):
//   - No direct calls to agenor from Bootstrap
//   - No UI rendering logic in Bootstrap
//   - No flag interpretation beyond resolving --tui and --theme
type Bootstrap struct {
	container *assist.CobraContainer
	options   ConfigureOptions

	// Cfg is populated after configure() reads viper.
	Cfg *bedrock.Config

	// UI is resolved from --tui and --theme in PersistentPreRunE and
	// passed into requests. Bootstrap does not use it directly.
	UI report.Presenter

	// themeLoader resolves named themes from the themes directory.
	// Constructed once in Root() and reused in PersistentPreRunE.
	themeLoader *bedrock.ThemeLoader

	// coord is the single Coordinator instance wired at startup and
	// shared by all command handlers.
	coord *jac.Coordinator

	// root param-set
	rootPs *assist.ParamSet[RootParameterSet]

	// nav ghost param-set and persistent families, inherited by walk, run, and query.
	navPs       *assist.ParamSet[NavParameterSet]
	previewFam  *assist.ParamSet[store.PreviewParameterSet]
	cascadeFam  *assist.ParamSet[store.CascadeParameterSet]
	samplingFam *assist.ParamSet[store.SamplingParameterSet]
	polyFam     *assist.ParamSet[store.PolyFilterParameterSet]

	// run-exclusive family
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

// Root builds the command tree and returns the root command, ready to
// be executed. This is the composition root: all wiring happens here
// and only here.
func (b *Bootstrap) Root(options ...ConfigureOptionFn) *cobra.Command {
	b.prepare()

	for _, fo := range options {
		fo(&b.options)
	}

	b.configure()

	// Detect the shell environment once at startup.
	env, err := shell.Detect()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Construct the ThemeLoader once - it resolves JAY_THEMES_DIR or
	// falls back to the XDG default (~/.config/jay/themes/).
	b.themeLoader = bedrock.NewThemeLoader()

	// Construct the Coordinator once with the detected locate function.
	b.coord = jac.New(b.Cfg, jac.WithLocate(env.Locate))

	b.container = assist.NewCobraContainer(
		&cobra.Command{
			Use:     ApplicationName,
			Short:   li18ngo.Text(locale.RootCmdShortDescTemplData{}),
			Long:    li18ngo.Text(locale.RootCmdLongDescTemplData{}),
			Version: fmt.Sprintf("'%v'", Version),

			// PersistentPreRunE resolves --tui and --theme into a
			// ui.Presenter backed by the appropriate prism.Renderer.
			// This is the only UI concern Bootstrap is permitted to touch.
			PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
				palette, err := b.themeLoader.Load(b.rootPs.Native.Theme)
				if err != nil {
					return err
				}

				mgr, err := ui.New(b.rootPs.Native.TUI, palette)
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

	// Registration order matters - parent must be registered before children.
	b.buildRootCommand(b.container)
	b.buildNavCommand(b.container) // ghost; must precede walk/run/query
	b.buildWalkCommand(b.container)
	b.buildRunCommand(b.container)
	b.buildQueryCommand(b.container)

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
			ConfigFileName: b.options.Config.Name,
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

	err := li18ngo.Register(func(uo *li18ngo.UseOptions) {
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

// buildRootCommand registers the truly global persistent flags on root:
// --tui, --theme, and --language. Navigation families live on the ghost
// nav command so that utility commands (verify, theme) do not inherit them.
func (b *Bootstrap) buildRootCommand(container *assist.CobraContainer) {
	root := container.Root()

	b.rootPs = assist.NewParamSet[RootParameterSet](root)

	// --tui(-t): display mode, inherited by all sub-commands.
	b.rootPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			`tui display mode: "linear" (default), "porthole", "lanes"`,
			"t",
			ui.ModeDefault,
			root.PersistentFlags(),
		),
		&b.rootPs.Native.TUI,
	)

	// --theme: colour theme name, inherited by all sub-commands.
	// "system" (default) uses ANSI-16 colours from the terminal theme.
	// Any other value loads <name>.yaml from the themes directory.
	b.rootPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			fmt.Sprintf(
				`colour theme name (default "system" uses terminal theme colours; `+
					`custom themes loaded from %s)`,
				b.themeLoader.ThemesDir(),
			),
			"",
			bedrock.ThemeSystemName,
			root.PersistentFlags(),
		),
		&b.rootPs.Native.Theme,
	)

	container.MustRegisterParamSet(RootPsName, b.rootPs)
}

// ---------------------------------------------------------------------------
// navFamilies
// ---------------------------------------------------------------------------

// navFamilies is a convenience accessor used by runWalk, runRun, and
// runQuery to obtain the nav-level flag families in a single bundle.
func (b *Bootstrap) navFamilies() NavFamilies {
	return NavFamilies{
		Preview:  b.previewFam,
		Cascade:  b.cascadeFam,
		Sampling: b.samplingFam,
		PolyFam:  b.polyFam,
	}
}
