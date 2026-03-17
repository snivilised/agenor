package command

import (
	"fmt"
	"os"

	"github.com/cubiest/jibberjabber"
	"github.com/snivilised/agenor/internal/third/lo"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/li18ngo"
	"github.com/snivilised/mamba/assist"
	"github.com/snivilised/mamba/assist/cfg"
	si18n "github.com/snivilised/mamba/locale"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

// LocaleDetector abstracts the detection of the user's preferred
// language as a BCP 47 language tag.
type LocaleDetector interface {
	Scan() language.Tag
}

// Jabber is a LocaleDetector implemented using jibberjabber.
type Jabber struct {
}

// Scan returns the detected language tag.
func (j *Jabber) Scan() language.Tag {
	lang, _ := jibberjabber.DetectIETF()
	return language.MustParse(lang)
}

// ConfigInfo describes the configuration file that should be loaded,
// including its name, type, path and the viper instance to use.
type ConfigInfo struct {
	// Name of the configuration name.
	Name string
	// ConfigType specifies the configuration type.
	ConfigType string
	// ConfigPath is the path to the configuration.
	ConfigPath string
	// Viper instance to load the configuration.
	Viper cfg.ViperConfig
}

// ConfigureOptions groups configuration options that influence how
// Bootstrap initialises localisation and configuration.
type ConfigureOptions struct {
	// Detector for locale identification.
	Detector LocaleDetector
	// Config describing the configuration.
	Config ConfigInfo
}

// ConfigureOptionFn is a functional option used to modify
// ConfigureOptions before Bootstrap performs its setup.
type ConfigureOptionFn func(*ConfigureOptions)

// Bootstrap represents construct that performs start up of the cli
// without resorting to the use of Go's init() mechanism and minimal
// use of package global variables.
type Bootstrap struct {
	container *assist.CobraContainer
	options   ConfigureOptions
}

func (b *Bootstrap) prepare() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	b.options = ConfigureOptions{
		Detector: &Jabber{},
		Config: ConfigInfo{
			Name:       ApplicationName,
			ConfigType: "yaml",
			ConfigPath: home,
			Viper:      &cfg.GlobalViperConfig{},
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

	// all these string literals here should be made translate-able
	//

	b.container = assist.NewCobraContainer(
		&cobra.Command{
			Use:     ApplicationName,
			Short:   li18ngo.Text(locale.RootCmdShortDescTemplData{}),
			Long:    li18ngo.Text(locale.RootCmdLongDescTemplData{}),
			Version: fmt.Sprintf("'%v'", Version),
			Run: func(_ *cobra.Command, _ []string) {
				fmt.Println("=== jay ===")
			},
		},
	)

	b.buildRootCommand(b.container)

	return b.container.Root()
}

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
				SourceID: li18ngo.TranslationSource{Name: ApplicationName},

				// By adding in the source for mamba, we relieve the client from having
				// to do this. After-all, it should be taken as read that since any
				// instantiation of jay (ie a project using this template) is by
				// necessity dependent on mamba, it's source should be loaded so that a
				// localizer can be created for it.
				//
				// The client app has to make sure that when their app is deployed,
				// the translations file(s) for mamba are named as 'mamba', as you
				// can see below, that is the name assigned to the app name of the
				// source.
				//
				si18n.MambaSourceID: li18ngo.TranslationSource{Name: "mamba"},
			},
		}
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (b *Bootstrap) buildRootCommand(container *assist.CobraContainer) {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	//
	root := container.Root()
	paramSet := assist.NewParamSet[RootParameterSet](root)

	container.MustRegisterParamSet(RootPsName, paramSet)
}
