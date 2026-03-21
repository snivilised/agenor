package locale

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// UsingConfigFileTemplData is the template data for the message indicating
// which config file is being used.
type UsingConfigFileTemplData struct {
	agenorTemplData
	ConfigFileName string
}

// Message returns the i18n message for UsingConfigFileTemplData.
func (td UsingConfigFileTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "using-config-file",
		Description: "Message to indicate which config is being used",
		Other:       "Using config file: {{.ConfigFileName}}",
	}
}
