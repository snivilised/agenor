package pref

import (
	"fmt"
)

type UsageError struct {
	message string
}

func (e UsageError) Error() string {
	// TODO: i18n
	return fmt.Sprintf("usage error: %v", e.message)
}
