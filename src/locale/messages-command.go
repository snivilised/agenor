package locale

import "github.com/nicksnyder/go-i18n/v2/i18n"

// 🧊 Root Cmd Short Description

// RootCmdShortDescTemplData provides template data for the short
// description of the root command.
type RootCmdShortDescTemplData struct {
	agenorTemplData
}

func (td RootCmdShortDescTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "root-command-short-description",
		Description: "short description for the root command",
		Other:       "A brief description of your application",
	}
}

// 🧊 Root Cmd Long Description

// RootCmdLongDescTemplData provides template data for the long
// description of the root command.
type RootCmdLongDescTemplData struct {
	agenorTemplData
}

func (td RootCmdLongDescTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "root-command-long-description",
		Description: "long description for the root command",
		Other: `A longer description that spans multiple lines and likely contains
		examples and usage of using your application. For example:
		
		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
	}
}

// 🧊 Root Cmd Config File Usage

// RootCmdConfigFileUsageTemplData supplies template data for the
// usage text of the root command's config file flag.
type RootCmdConfigFileUsageTemplData struct {
	agenorTemplData
	// ConfigFileName represents the configuration file name.
	ConfigFileName string
}

func (td RootCmdConfigFileUsageTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "root-command-config-file-usage",
		Description: "root command config flag usage",
		Other:       "config file (default is $HOME/{{.ConfigFileName}}.yml)",
	}
}

// 🧊 Root Cmd Lang Usage

// RootCmdLangUsageTemplData provides template data for the language
// flag usage text on the root command.
type RootCmdLangUsageTemplData struct {
	agenorTemplData
}

func (td RootCmdLangUsageTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "root-command-language-usage",
		Description: "root command lang usage",
		Other:       "'lang' defines the language according to IETF BCP 47",
	}
}
