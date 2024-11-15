package locale

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// 🌷 prohibitive Word

// PatternTemplData
type ProhibitiveWordTemplData struct {
	agenorTemplData
}

// Message
func (td ProhibitiveWordTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "prohibitive.word",
		Description: "prohibitive",
		Other:       "prohibitive",
	}
}

// 🌷 permissive Word

// PatternTemplData
type PermissiveWordTemplData struct {
	agenorTemplData
}

// Message
func (td PermissiveWordTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "permissive.word",
		Description: "permissive",
		Other:       "permissive",
	}
}
