package i18n

import "errors"

// these errors are to be converted into proper i18n errors

var (
	ErrFilterIsNil                        = errors.New("filter is nil")
	ErrFilterMissingType                  = errors.New("filter missing type")
	ErrFilterCustomNotSupported           = errors.New("custom filter not supported for children")
	ErrFilterUndefined                    = errors.New("filter is undefined")
	ErrInternalFailedToGetNavigatorDriver = errors.New("failed to get navigator driver (internal)")
	ErrInvalidIncaseFilterDef             = errors.New("invalid incase filter definition; pattern is missing separator")
	ErrWorkerPoolCreationFailed           = errors.New("failed to create worker pool")
)
