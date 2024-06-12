package i18n

import "errors"

// these errors are to be converted into proper i18n errors

var (
	ErrWorkerPoolCreationFailed           = errors.New("failed to create worker pool")
	ErrInternalFailedToGetNavigatorDriver = errors.New("failed to get navigator driver (internal)")
)
