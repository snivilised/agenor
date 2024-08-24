package locale

import (
	"github.com/pkg/errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/snivilised/li18ngo"
)

// ❌ FilterIsNil

// FilterIsNilTemplData
type FilterIsNilErrorTemplData struct {
	traverseTemplData
}

// Message
func (td FilterIsNilErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "filter-is-nil.error",
		Description: "filter is nil error",
		Other:       "filter is nil",
	}
}

type FilterIsNilError struct {
	li18ngo.LocalisableError
}

var ErrFilterIsNil = FilterIsNilError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterIsNilErrorTemplData{},
	},
}

// ❌ FilterMissingType

// FilterMissingTypeTemplData
type FilterMissingTypeErrorTemplData struct {
	traverseTemplData
}

// Message
func (td FilterMissingTypeErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "filter-missing-type.error",
		Description: "filter missing type",
		Other:       "filter missing type",
	}
}

type FilterMissingTypeError struct {
	li18ngo.LocalisableError
}

var ErrFilterMissingType = FilterMissingTypeError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterMissingTypeErrorTemplData{},
	},
}

// ❌ FilterCustomNotSupported

// FilterCustomNotSupportedTemplData
type FilterCustomNotSupportedErrorTemplData struct {
	traverseTemplData
}

// Message
func (td FilterCustomNotSupportedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "custom-filter-not-supported-for-children.error",
		Description: "custom filter not supported for children",
		Other:       "custom filter not supported for children",
	}
}

type FilterCustomNotSupportedError struct {
	li18ngo.LocalisableError
}

var ErrFilterCustomNotSupported = FilterCustomNotSupportedError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterCustomNotSupportedErrorTemplData{},
	},
}

// ❌ FilterUndefined

// FilterUndefinedTemplData
type FilterUndefinedErrorTemplData struct {
	traverseTemplData
}

// Message
func (td FilterUndefinedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "filter-is-undefined.error",
		Description: "filter is undefined error",
		Other:       "filter is undefined",
	}
}

type FilterUndefinedError struct {
	li18ngo.LocalisableError
}

var ErrFilterUndefined = FilterUndefinedError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterUndefinedErrorTemplData{},
	},
}

// ❌ InternalFailedToGetNavigatorDriver

// InternalFailedToGetNavigatorDriverTemplData
type InternalFailedToGetNavigatorDriverErrorTemplData struct {
	traverseTemplData
}

// Message
func (td InternalFailedToGetNavigatorDriverErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "failed-to-get-navigator-driver.error",
		Description: "failed to get navigator driver",
		Other:       "failed to get navigator driver",
	}
}

type InternalFailedToGetNavigatorDriverError struct {
	li18ngo.LocalisableError
}

var ErrInternalFailedToGetNavigatorDriver = InternalFailedToGetNavigatorDriverError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InternalFailedToGetNavigatorDriverErrorTemplData{},
	},
}

// ❌ InvalidIncaseFilterDef error

// InvalidIncaseFilterDefTemplData
type InvalidIncaseFilterDefErrorTemplData struct {
	traverseTemplData
}

// Message
func (td InvalidIncaseFilterDefErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-incase-filter-definition.error",
		Description: "invalid incase filter definition; pattern is missing separator",
		Other:       "invalid incase filter definition; pattern is missing separator",
	}
}

type InvalidIncaseFilterDefError struct {
	li18ngo.LocalisableError
}

// IsInvalidIncaseFilterDefError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidIncaseFilterDefError
func IsInvalidIncaseFilterDefError(err error) bool {
	return errors.Is(err, errInvalidIncaseFilterDef)
}

func NewInvalidIncaseFilterDefError(pattern string) error {
	return errors.Wrap(
		errInvalidIncaseFilterDef,
		li18ngo.Text(PatternFieldTemplData{
			Pattern: pattern,
		}),
	)
}

var errInvalidIncaseFilterDef = InvalidIncaseFilterDefError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidIncaseFilterDefErrorTemplData{},
	},
}

// ❌ WorkerPoolCreationFailed

// WorkerPoolCreationFailedTemplData
type WorkerPoolCreationFailedErrorTemplData struct {
	traverseTemplData
}

// Message
func (td WorkerPoolCreationFailedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "failed-to-create-worker-pool.error",
		Description: "failed to create worker pool",
		Other:       "failed to create worker pool",
	}
}

type WorkerPoolCreationFailedError struct {
	li18ngo.LocalisableError
}

var ErrWorkerPoolCreationFailed = WorkerPoolCreationFailedError{
	LocalisableError: li18ngo.LocalisableError{
		Data: WorkerPoolCreationFailedErrorTemplData{},
	},
}

// ❌ InvalidFileSamplingSpecMissingFilesError

// InvalidFileSamplingSpecMissingFilesErrorTemplData
type InvalidFileSamplingSpecMissingFilesErrorTemplData struct {
	traverseTemplData
}

// Message
func (td InvalidFileSamplingSpecMissingFilesErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-file-sampling-spec-missing-files.error",
		Description: "invalid file sampling specification, missing no of files",
		Other:       "invalid file sampling specification, missing no of files",
	}
}

type InvalidFileSamplingSpecificationError struct {
	li18ngo.LocalisableError
}

var ErrInvalidFileSamplingSpecMissingFiles = InvalidFileSamplingSpecificationError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidFileSamplingSpecMissingFilesErrorTemplData{},
	},
}

// ❌ InvalidFolderSamplingSpecMissingFolders

// InvalidFolderSamplingSpecMissingFoldersTemplData
type InvalidFolderSamplingSpecMissingFoldersErrorTemplData struct {
	traverseTemplData
}

// Message
func (td InvalidFolderSamplingSpecMissingFoldersErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-file-sampling-spec-missing-folders.error",
		Description: "invalid file sampling specification, missing no of folders",
		Other:       "invalid file sampling specification, missing no of folders",
	}
}

type InvalidFolderSamplingSpecMissingFoldersError struct {
	li18ngo.LocalisableError
}

var ErrInvalidFolderSamplingSpecMissingFolders = InvalidFolderSamplingSpecMissingFoldersError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidFolderSamplingSpecMissingFoldersErrorTemplData{},
	},
}

// ❌ MissingCustomFilterDefinition

// MissingCustomFilterDefinitionTemplData
type MissingCustomFilterDefinitionErrorTemplData struct {
	traverseTemplData
}

// Message
func (td MissingCustomFilterDefinitionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "missing-custom-filter-definition.error",
		Description: "config error missing-custom-filter-definition",
		Other:       "missing custom filter definition (config error)",
	}
}

type MissingCustomFilterDefinitionError struct {
	li18ngo.LocalisableError
}

var ErrMissingCustomFilterDefinition = MissingCustomFilterDefinitionError{
	LocalisableError: li18ngo.LocalisableError{
		Data: MissingCustomFilterDefinitionErrorTemplData{},
	},
}

// to define variable error with simple field - "Field and Variable error/fv18e"
// "Simple i18n Field"
// followed by

// 🍀 Pattern

// PatternTemplData
type PatternFieldTemplData struct {
	traverseTemplData
	Pattern string
}

// Message
func (td PatternFieldTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "pattern.field",
		Description: "pattern",
		Other:       "pattern: {{.Pattern}}",
	}
}

// ❌ InvalidExtGlobFilterMissingSeparator

// InvalidExtGlobFilterMissingSeparatorTemplData
type InvalidExtGlobFilterMissingSeparatorErrorTemplData struct {
	traverseTemplData
}

// Message
func (td InvalidExtGlobFilterMissingSeparatorErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-extended-glob-filter-missing-separator.error",
		Description: "invalid extended glob filter definition; pattern is missing separator",
		Other:       "invalid extended glob filter definition; pattern is missing separator",
	}
}

type InvalidExtGlobFilterMissingSeparatorError struct {
	li18ngo.LocalisableError
}

// IsInvalidExtGlobFilterMissingSeparatorError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidExtGlobFilterMissingSeparatorError
func IsInvalidExtGlobFilterMissingSeparatorError(err error) bool {
	return errors.Is(err, errInvalidExtGlobFilterMissingSeparator)
}

func NewInvalidExtGlobFilterMissingSeparatorError(pattern string) error {
	return errors.Wrap(
		errInvalidExtGlobFilterMissingSeparator,
		li18ngo.Text(PatternFieldTemplData{
			Pattern: pattern,
		}),
	)
}

var errInvalidExtGlobFilterMissingSeparator = InvalidExtGlobFilterMissingSeparatorError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidExtGlobFilterMissingSeparatorErrorTemplData{},
	},
}

// ❌ UsageMissingRootPath

// UsageMissingRootPathTemplData
type UsageMissingRootPathErrorTemplData struct {
	traverseTemplData
}

// Message
func (td UsageMissingRootPathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-root-path.error",
		Description: "usage missing root path",
		Other:       "usage missing root path",
	}
}

type UsageMissingRootPathError struct {
	li18ngo.LocalisableError
}

var ErrUsageMissingRootPath = UsageMissingRootPathError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UsageMissingRootPathErrorTemplData{},
	},
}

// ❌ UsageMissingRestorePath

// UsageMissingRestorePathTemplData
type UsageMissingRestorePathErrorTemplData struct {
	traverseTemplData
}

// Message
func (td UsageMissingRestorePathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-restore-path.error",
		Description: "usage missing restore path",
		Other:       "usage missing restore path",
	}
}

type UsageMissingRestorePathError struct {
	li18ngo.LocalisableError
}

var ErrUsageMissingRestorePath = UsageMissingRestorePathError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UsageMissingRestorePathErrorTemplData{},
	},
}

// ❌ UsageMissingSubscription

// UsageMissingSubscriptionTemplData
type UsageMissingSubscriptionErrorTemplData struct {
	traverseTemplData
}

// Message
func (td UsageMissingSubscriptionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-subscription.error",
		Description: "usage missing subscription",
		Other:       "usage missing subscription",
	}
}

type UsageMissingSubscriptionError struct {
	li18ngo.LocalisableError
}

var ErrUsageMissingSubscription = UsageMissingSubscriptionError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UsageMissingSubscriptionErrorTemplData{},
	},
}

// ❌ UsageMissingHandler

// UsageMissingHandlerTemplData
type UsageMissingHandlerErrorTemplData struct {
	traverseTemplData
}

// Message
func (td UsageMissingHandlerErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-handler.error",
		Description: "usage missing handler",
		Other:       "usage missing handler",
	}
}

type UsageMissingHandlerError struct {
	li18ngo.LocalisableError
}

var ErrUsageMissingHandler = UsageMissingHandlerError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UsageMissingHandlerErrorTemplData{},
	},
}

// ❌ IDGeneratorFuncCantBeNil

// IDGeneratorFuncCantBeNilTemplData
type IDGeneratorFuncCantBeNilErrorTemplData struct {
	traverseTemplData
}

// Message
func (td IDGeneratorFuncCantBeNilErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "id-generator-func-cant-be-nil.error",
		Description: "id generator func is nil, should be defined",
		Other:       "id generator func can't be nil",
	}
}

type IDGeneratorFuncCantBeNilError struct {
	li18ngo.LocalisableError
}

var ErrIDGeneratorFuncCantBeNil = IDGeneratorFuncCantBeNilError{
	LocalisableError: li18ngo.LocalisableError{
		Data: IDGeneratorFuncCantBeNilErrorTemplData{},
	},
}

// ❌ FooBar

// FooBarTemplData - TODO: this is a none existent error that should be
// replaced by the client. Its just defined here to illustrate the pattern
// that should be used to implement i18n with li18ngo. Also note,
// that this message has been removed from the translation files, so
// it is not useable at run time.
type FooBarTemplData struct {
	traverseTemplData
	Path   string
	Reason error
}

// the ID should use spp/library specific code, so replace astrolib with the
// name of the library implementing this template project.
func (td FooBarTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "foo-bar.traverse.nav",
		Description: "Foo Bar description",
		Other:       "foo bar failure '{{.Path}}' (reason: {{.Reason}})",
	}
}

// FooBarErrorBehaviourQuery used to query if an error is:
// "Failed to read directory contents from the path specified"
type FooBarErrorBehaviourQuery interface {
	FooBar() bool
}

type FooBarError struct {
	li18ngo.LocalisableError
}

// FooBar enables the client to check if error is FooBarError
// via FooBarErrorBehaviourQuery
func (e FooBarError) FooBar() bool {
	return true
}

// NewFooBarError creates a FooBarError
func NewFooBarError(path string, reason error) FooBarError {
	return FooBarError{
		LocalisableError: li18ngo.LocalisableError{
			Data: FooBarTemplData{
				Path:   path,
				Reason: reason,
			},
		},
	}
}
