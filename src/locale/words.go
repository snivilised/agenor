package locale

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// 🌷 prohibitive Word

// ProhibitiveWordTemplData is the template data for the message indicating a prohibitive word.
type ProhibitiveWordTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td ProhibitiveWordTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "prohibitive.word",
		Description: "prohibitive",
		Other:       "prohibitive",
	}
}

// 🌷 permissive Word

// PermissiveWordTemplData is the template data for the message indicating a permissive word.
type PermissiveWordTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td PermissiveWordTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "permissive.word",
		Description: "permissive",
		Other:       "permissive",
	}
}
