package pref

import (
	"github.com/snivilised/agenor/enums"
)

type (
	// PersistOptions defines the options for persisting data.
	PersistOptions struct {
		// Format specifies the format to use for persistence.
		Format enums.PersistenceFormat
	}
)
