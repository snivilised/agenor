package filter

import (
	"github.com/pkg/errors"
)

// ❌ InvalidNotificationMuteRequested error: this scenario should
// never happen, if it does then its an internal issue to be
// addressed here in the filter plugin.

// IsNoSubordinateHybridSchemesDefinedError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidNotificationMuteRequestedError
func IsNoSubordinateHybridSchemesDefinedError(err error) bool {
	return errors.Is(err, ErrNoSubordinateHybridSchemesDefined)
}

var ErrNoSubordinateHybridSchemesDefined = errors.New(
	"invalid filter scheme, both primary and child not set",
)
