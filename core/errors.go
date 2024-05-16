package core

import (
	"fmt"
)

type UsingError struct {
	message string
}

func (e UsingError) Error() string {
	// TODO: i18n
	return fmt.Sprintf("using error: %v", e.message)
}
