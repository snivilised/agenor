package command

import (
	"fmt"
	"os"

	"github.com/cubiest/jibberjabber"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/locale"
	"github.com/snivilised/jaywalk/src/third/lo"
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
	Detector  LocaleDetector
	Config    ConfigInfo
	GetForest pref.BuildForest
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
//	root                     (global persistent flags: --tui, --theme, --language)
//	  nav (ghost)            (nav persistent flags: --subscribe, --action, --pipeline
//	  │                       + cascade, sampling, preview, poly-filter families)
//	  │  exec (ghost)        (exec persistent flags: --resume
//	  │  │                    + MarkFlagsOneRequired("action", "pipeline"))
//	  │  │  walk
//	  │  │  sprint
//	  │  query
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

	// nav ghost param-set and persistent families, inherited by walk, sprint, and query.
	navPs       *assist.ParamSet[NavParameterSet]
	previewFam  *assist.ParamSet[store.PreviewParameterSet]
	cascadeFam  *assist.ParamSet[store.CascadeParameterSet]
	samplingFam *assist.ParamSet[store.SamplingParameterSet]
	polyFam     *assist.ParamSet[store.PolyFilterParameterSet]

	// exec ghost param-set, inherited by walk and sprint only.
	execPs *assist.ParamSet[ExecParameterSet]

	// sprint-exclusive family
	workerPoolFam *assist.ParamSet[store.WorkerPoolParameterSet]
}

// ---------------------------------------------------------------------------
// Root
// ---------------------------------------------------------------------------

func (b *Bootstrap) prepare() {
	home, err := core.Home()
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

	env, err := shell.Detect()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	b.themeLoader = bedrock.NewThemeLoader()
	b.coord = jac.New(b.Cfg,
		jac.WithLocate(env.Locate),
		jac.WithForest(b.options.GetForest),
	)

	b.container = assist.NewCobraContainer(
		&cobra.Command{
			Use:     ApplicationName,
			Short:   li18ngo.Text(locale.RootCmdShortDescTemplData{}),
			Long:    li18ngo.Text(locale.RootCmdLongDescTemplData{}),
			Version: fmt.Sprintf("'%v'", Version),

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

	// Registration order matters: parent must be registered before children.
	b.buildRootCommand(b.container)
	b.buildNavCommand(b.container)
	b.buildExecCommand(b.container)
	b.buildWalkCommand(b.container)
	b.buildSprintCommand(b.container)
	b.buildQueryCommand(b.container)

	root := b.container.Root()

	// Inject ghost ancestors into os.Args so that the user can type
	// 'jay walk' instead of 'jay nav exec walk'. SetArgs is used rather
	// than mutating os.Args directly so that cobra's test harness (which
	// calls SetArgs itself) works correctly: the harness sets args after
	// Root() returns, so tests must pipe their args through
	// InjectGhostAncestors themselves (see walk-cmd_test.go).
	root.SetArgs(InjectGhostAncestors(os.Args[1:]))

	return root
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

	b.rootPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			li18ngo.Text(locale.TuiFlagDescTemplData{}),
			"t",
			ui.ModeDefault,
			root.PersistentFlags(),
		),
		&b.rootPs.Native.TUI,
	)

	b.rootPs.BindString(
		assist.NewFlagInfoOnFlagSet(
			li18ngo.Text(locale.NewThemeFlagDescTemplData(b.themeLoader.ThemesDir())),
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

// navFamilies is a convenience accessor used by runWalk, runSprint, and
// runQuery to obtain the nav-level flag families in a single bundle.
func (b *Bootstrap) navFamilies() NavFamilies {
	return NavFamilies{
		Preview:  b.previewFam,
		Cascade:  b.cascadeFam,
		Sampling: b.samplingFam,
		PolyFam:  b.polyFam,
	}
}

// ---------------------------------------------------------------------------
// InjectGhostAncestors
// ---------------------------------------------------------------------------

// ghostPrefixes maps each user-visible leaf command to the ghost ancestors
// that cobra requires for correct flag inheritance. The slice is ordered
// from outermost to innermost ghost, matching the command tree:
//
//	root → nav → exec → walk/sprint
//	root → nav → query
//
// Adding a new leaf under ghost parents requires only a new entry here.
var ghostPrefixes = map[string][]string{
	"walk":   {"nav", "exec"},
	"sprint": {"nav", "exec"},
	"query":  {"nav"},
}

// InjectGhostAncestors takes a cobra args slice (i.e. os.Args[1:] or the
// slice passed to cmd.SetArgs) and, when the first non-flag token is a
// known leaf command, splices in the ghost ancestor tokens that cobra needs
// for correct flag inheritance. All other invocations are returned unchanged.
//
// It is exported so that tests can pipe their args through it before passing
// them to CommandTester, mirroring exactly what Root() does for production.
//
// Example:
//
//	in:  ["walk", "RETRO-WAVE", "--action", "echo"]
//	out: ["nav", "exec", "walk", "RETRO-WAVE", "--action", "echo"]
func InjectGhostAncestors(args []string) []string {
	// Scan for the first non-flag token (the sub-command).
	// Flags before the sub-command (e.g. --tui porthole walk) are preserved.
	insertAt := -1
	leaf := ""

	for i, arg := range args {
		if len(arg) > 0 && arg[0] != '-' {
			insertAt = i
			leaf = arg
			break
		}
	}

	if insertAt == -1 {
		return args // no sub-command token; nothing to rewrite
	}

	prefixes, ok := ghostPrefixes[leaf]
	if !ok {
		return args // not a leaf that needs ghost injection
	}

	// Build the rewritten slice:
	//   args[:insertAt]  — any flags preceding the sub-command
	//   prefixes         — injected ghost tokens
	//   args[insertAt:]  — leaf command + everything that follows
	rewritten := make([]string, 0, len(args)+len(prefixes))
	rewritten = append(rewritten, args[:insertAt]...)
	rewritten = append(rewritten, prefixes...)
	rewritten = append(rewritten, args[insertAt:]...)

	return rewritten
}
