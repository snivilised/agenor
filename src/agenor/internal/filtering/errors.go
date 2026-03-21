package filtering

import (
	"fmt"

	"github.com/pkg/errors"
)

// ❌ InvalidResumeStateTransition error

// NewInvalidResumeStateTransitionError creates an untranslated error to
// indicate in invalid resume state transition
func NewInvalidResumeStateTransitionError(state string) error {
	return errors.Wrap(
		errInvalidResumeStateTransition,
		fmt.Sprintf("state: %v", state),
	)
}

// IsInvalidResumeStateTransitionError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidResumeStateTransitionNativeError
func IsInvalidResumeStateTransitionError(err error) bool {
	return errors.Is(err, errInvalidResumeStateTransition)
}

var errInvalidResumeStateTransition = errors.New(
	"invalid resume state transition detected",
)
