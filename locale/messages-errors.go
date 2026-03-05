package locale

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/snivilised/li18ngo"
)

// ❌ FilterIsNil

// FilterIsNilErrorTemplData is the template data for the FilterIsNil error message.
type FilterIsNilErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td FilterIsNilErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "filter-is-nil.age.static-error",
		Description: "filter is nil error",
		Other:       "filter is nil",
	}
}

// FilterIsNilError is the error type for the error when filter is nil.
type FilterIsNilError struct {
	li18ngo.LocalisableError
}

// ErrFilterIsNil is the exported error variable for FilterIsNilError with
// the template data.
var ErrFilterIsNil = FilterIsNilError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterIsNilErrorTemplData{},
	},
}

// ❌ FilterMissingType

// FilterMissingTypeErrorTemplData is the template data for the
// FilterMissingType error
type FilterMissingTypeErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td FilterMissingTypeErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "filter-missing-type.age.static-error",
		Description: "filter missing type",
		Other:       "filter missing type",
	}
}

// FilterMissingTypeError is the error type for the error when filter is missing type.
type FilterMissingTypeError struct {
	li18ngo.LocalisableError
}

// ErrFilterMissingType is the exported error variable for FilterMissingTypeError
// with the template data.
var ErrFilterMissingType = FilterMissingTypeError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterMissingTypeErrorTemplData{},
	},
}

// ❌ FilterCustomNotSupported

// FilterCustomNotSupportedErrorTemplData is the template data for the
// FilterCustomNotSupported error
type FilterCustomNotSupportedErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td FilterCustomNotSupportedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "custom-filter-not-supported-for-children.age.static-error",
		Description: "custom filter not supported for children",
		Other:       "custom filter not supported for children",
	}
}

// FilterCustomNotSupportedError is the error type for the error when custom filter
// is not supported for children.
type FilterCustomNotSupportedError struct {
	li18ngo.LocalisableError
}

// ErrFilterCustomNotSupported is the exported error variable for
// FilterCustomNotSupportedError with the template data.
var ErrFilterCustomNotSupported = FilterCustomNotSupportedError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterCustomNotSupportedErrorTemplData{},
	},
}

// ❌ FilterChildGlobExNotSupported

// FilterChildGlobExNotSupportedErrorTemplData is the template data for the
// FilterChildGlobExNotSupported error message.
type FilterChildGlobExNotSupportedErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td FilterChildGlobExNotSupportedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "glob-ex-filter-not-supported-for-children.age.static-error",
		Description: "glob-ex filter not supported for children",
		Other:       "glob-ex filter not supported for children",
	}
}

// FilterChildGlobExNotSupportedError is the error type for the error when glob-ex filter
// is not supported for children.
type FilterChildGlobExNotSupportedError struct {
	li18ngo.LocalisableError
}

// ErrFilterChildGlobExNotSupported is the exported error variable for
// FilterChildGlobExNotSupportedError with the template data.
var ErrFilterChildGlobExNotSupported = FilterChildGlobExNotSupportedError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterChildGlobExNotSupportedErrorTemplData{},
	},
}

// ❌ FilterUndefined

// FilterUndefinedErrorTemplData is the template data for the FilterUndefined error
// message.
type FilterUndefinedErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td FilterUndefinedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "filter-is-undefined.age.static-error",
		Description: "filter is undefined error",
		Other:       "filter is undefined",
	}
}

// FilterUndefinedError error message
type FilterUndefinedError struct {
	li18ngo.LocalisableError
}

// ErrFilterUndefined is the exported error variable for FilterUndefinedError with
// the template data.
var ErrFilterUndefined = FilterUndefinedError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterUndefinedErrorTemplData{},
	},
}

// ❌ InternalFailedToGetNavigatorDriver

// InternalFailedToGetNavigatorDriverErrorTemplData is the template data for the
// error when failed to get navigator driver.
type InternalFailedToGetNavigatorDriverErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td InternalFailedToGetNavigatorDriverErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "failed-to-get-navigator-driver.age.static-error",
		Description: "failed to get navigator driver",
		Other:       "failed to get navigator driver",
	}
}

// InternalFailedToGetNavigatorDriverError is the error type for the error when failed
// to get navigator driver.
type InternalFailedToGetNavigatorDriverError struct {
	li18ngo.LocalisableError
}

// ErrInternalFailedToGetNavigatorDriver is the exported error variable for
// InternalFailedToGetNavigatorDriverError with the template data.
var ErrInternalFailedToGetNavigatorDriver = InternalFailedToGetNavigatorDriverError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InternalFailedToGetNavigatorDriverErrorTemplData{},
	},
}

// ❌ InvalidInCaseFilterDefError

// InvalidInCaseFilterDefTemplData is the template data for the InvalidInCaseFilterDef
// error message.
type InvalidInCaseFilterDefTemplData struct {
	agenorTemplData
	Pattern string
}

// Message creates a new i18n message using the template data.
func (td InvalidInCaseFilterDefTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-incase-filter-definition.age.dynamic-error",
		Description: "invalid incase filter definition; pattern is missing separator wrapper error",
		Other:       "pattern: {{.Pattern}}",
	}
}

// InvalidInCaseFilterDefError is the error type for the error when incase filter
// definition is invalid.
type InvalidInCaseFilterDefError struct {
	li18ngo.LocalisableError

	// Wrapped error
	Wrapped error
}

// Error returns the error message for InvalidInCaseFilterDefError by combining
// the wrapped error message and the i18n message created from the template data.
func (e InvalidInCaseFilterDefError) Error() string {
	return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

func (e InvalidInCaseFilterDefError) Unwrap() error {
	return e.Wrapped
}

// NewInvalidInCaseFilterDefError creates a new InvalidInCaseFilterDefError with the given pattern
// and wraps the core error: CoreInvalidInCaseFilterDefError
func NewInvalidInCaseFilterDefError(pattern string) error {
	return &InvalidInCaseFilterDefError{
		LocalisableError: li18ngo.LocalisableError{
			Data: InvalidInCaseFilterDefTemplData{
				Pattern: pattern,
			},
		},
		Wrapped: errCoreInvalidInCaseFilterDef,
	}
}

// CoreInvalidInCaseFilterDefErrorTemplData is the template data for the
// CoreInvalidInCaseFilterDefError core error message.
type CoreInvalidInCaseFilterDefErrorTemplData struct {
	agenorTemplData
}

// IsInvalidInCaseFilterDefError uses errors.Is to check if the err's error
// tree contains the core error: CoreInvalidInCaseFilterDefError
func IsInvalidInCaseFilterDefError(err error) bool {
	return errors.Is(err, errCoreInvalidInCaseFilterDef)
}

// Message creates a new i18n message using the template data.
func (td CoreInvalidInCaseFilterDefErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-incase-filter-definition.age.core-error",
		Description: "invalid incase filter definition; pattern is missing separator core error",
		Other:       "invalid incase filter definition; pattern is missing separator",
	}
}

// CoreInvalidInCaseFilterDefError is the core error for invalid incase filter definition
// when pattern is missing separator.
type CoreInvalidInCaseFilterDefError struct {
	li18ngo.LocalisableError
}

var errCoreInvalidInCaseFilterDef = CoreInvalidInCaseFilterDefError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreInvalidInCaseFilterDefErrorTemplData{},
	},
}

// ❌ WorkerPoolCreationFailed

// WorkerPoolCreationFailedErrorTemplData is the template data for the
// WorkerPoolCreationFailed error message.
type WorkerPoolCreationFailedErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td WorkerPoolCreationFailedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "failed-to-create-worker-pool.age.static-error",
		Description: "failed to create worker pool",
		Other:       "failed to create worker pool",
	}
}

// WorkerPoolCreationFailedError is the error type for the error when failed
// to create worker pool.
type WorkerPoolCreationFailedError struct {
	li18ngo.LocalisableError
}

// ErrWorkerPoolCreationFailed is the exported error variable for
// WorkerPoolCreationFailedError with the template data.
var ErrWorkerPoolCreationFailed = WorkerPoolCreationFailedError{
	LocalisableError: li18ngo.LocalisableError{
		Data: WorkerPoolCreationFailedErrorTemplData{},
	},
}

// ❌ InvalidFileSamplingSpecMissingFilesError

// InvalidSamplingSpecMissingFilesErrorTemplData is the template data for the
// error when file sampling specification is invalid due to missing no of files.
type InvalidSamplingSpecMissingFilesErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td InvalidSamplingSpecMissingFilesErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-file-sampling-spec-missing-files.age.static-error",
		Description: "invalid file sampling specification, missing no of files",
		Other:       "invalid file sampling specification, missing no of files",
	}
}

// InvalidFileSamplingSpecificationError is the error type for the error when file
// sampling specification is invalid due to missing no of files.
type InvalidFileSamplingSpecificationError struct {
	li18ngo.LocalisableError
}

// ErrInvalidFileSamplingSpecMissingFiles is the exported error variable for
// InvalidFileSamplingSpecificationError with the template data.
var ErrInvalidFileSamplingSpecMissingFiles = InvalidFileSamplingSpecificationError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidSamplingSpecMissingFilesErrorTemplData{},
	},
}

// ❌ InvalidSamplingSpecMissingDirectoriesError

// InvalidSamplingSpecMissingDirectoriesErrorTemplData is the template data for
// the error when file sampling specification is invalid due to missing no of
// directories.
type InvalidSamplingSpecMissingDirectoriesErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td InvalidSamplingSpecMissingDirectoriesErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-file-sampling-spec-missing-directories.age.static-error",
		Description: "invalid file sampling specification, missing no of directories",
		Other:       "invalid file sampling specification, missing no of directories",
	}
}

// InvalidSamplingSpecMissingDirectoriesError is the error type for the error when
// file sampling specification is invalid due to missing no of directories.
type InvalidSamplingSpecMissingDirectoriesError struct {
	li18ngo.LocalisableError
}

// ErrInvalidSamplingSpecMissingDirectories is the exported error variable for
// InvalidSamplingSpecMissingDirectoriesError with the template data.
var ErrInvalidSamplingSpecMissingDirectories = InvalidSamplingSpecMissingDirectoriesError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidSamplingSpecMissingDirectoriesErrorTemplData{},
	},
}

// ❌ MissingCustomFilterDefinition

// MissingCustomFilterDefinitionErrorTemplData is the template data for the error
// when custom filter definition is missing in config.
type MissingCustomFilterDefinitionErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td MissingCustomFilterDefinitionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "missing-custom-filter-definition.age.static-error",
		Description: "config error missing-custom-filter-definition",
		Other:       "missing custom filter definition (config error)",
	}
}

// MissingCustomFilterDefinitionError is the error type for the error when custom
// filter definition is missing in config.
type MissingCustomFilterDefinitionError struct {
	li18ngo.LocalisableError
}

// ErrMissingCustomFilterDefinition is the exported error variable for
// MissingCustomFilterDefinitionError with the template data.
var ErrMissingCustomFilterDefinition = MissingCustomFilterDefinitionError{
	LocalisableError: li18ngo.LocalisableError{
		Data: MissingCustomFilterDefinitionErrorTemplData{},
	},
}

// 🍀 Pattern

// PatternFieldTemplData is the template data for the Pattern error message.
type PatternFieldTemplData struct {
	agenorTemplData

	// Pattern is the user provided pattern string containing the error.
	Pattern string
}

// Message creates a new i18n message using the template data.
func (td PatternFieldTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "pattern.field",
		Description: "pattern",
		Other:       "pattern: {{.Pattern}}",
	}
}

// ❌ InvalidExtGlobFilterMissingSeparator

// InvalidExtGlobFilterMissingSeparatorErrorTemplData is the template data for the
// error when glob-ex filter definition is invalid due to missing separator.
type InvalidExtGlobFilterMissingSeparatorErrorTemplData struct {
	agenorTemplData

	// Pattern is the user defined pattern string containing the error.
	Pattern string
}

// Message creates a new i18n message using the template data.
func (td InvalidExtGlobFilterMissingSeparatorErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-glob-ex-filter-missing-separator.age.dynamic-error",
		Description: "invalid glob ex filter definition; pattern is missing separator",
		Other:       "pattern: {{.Pattern}}",
	}
}

// InvalidExtGlobFilterMissingSeparatorError is the error type for the error when
// glob-ex filter definition is invalid due to missing separator.
type InvalidExtGlobFilterMissingSeparatorError struct {
	li18ngo.LocalisableError
	Wrapped error
}

// Error returns the error message for InvalidExtGlobFilterMissingSeparatorError
// by combining the wrapped error message and the i18n message created from the
// template data.
func (e InvalidExtGlobFilterMissingSeparatorError) Error() string {
	return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

// Unwrap returns the wrapped error for InvalidExtGlobFilterMissingSeparatorError.
func (e InvalidExtGlobFilterMissingSeparatorError) Unwrap() error {
	return e.Wrapped
}

// NewInvalidExtGlobFilterMissingSeparatorError creates a new
// InvalidExtGlobFilterMissingSeparatorError with the given pattern and wraps
// the core error: CoreInvalidExtGlobFilterMissingSeparatorError
func NewInvalidExtGlobFilterMissingSeparatorError(pattern string) error {
	return &InvalidExtGlobFilterMissingSeparatorError{
		Wrapped: errCoreInvalidExtGlobFilterMissingSeparator,
		LocalisableError: li18ngo.LocalisableError{
			Data: InvalidExtGlobFilterMissingSeparatorErrorTemplData{
				Pattern: pattern,
			},
		},
	}
}

// ❌ CoreInvalidExtGlobFilterMissingSeparator

// CoreInvalidExtGlobFilterMissingSeparatorErrorTemplData is the template data
// for the error when glob-ex filter definition is invalid due to
// missing separator.
type CoreInvalidExtGlobFilterMissingSeparatorErrorTemplData struct {
	agenorTemplData
}

// IsInvalidExtGlobFilterMissingSeparatorError uses errors.Is to check
// if the err's error tree contains the core error:
// CoreInvalidExtGlobFilterMissingSeparatorError
func IsInvalidExtGlobFilterMissingSeparatorError(err error) bool {
	return errors.Is(err, errCoreInvalidExtGlobFilterMissingSeparator)
}

// Message creates a new i18n message using the template data.
func (td CoreInvalidExtGlobFilterMissingSeparatorErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-extended-glob-filter-missing-separator.age.core-error",
		Description: "invalid glob ex filter definition; pattern is missing separator",
		Other:       "invalid glob ex filter definition; pattern is missing separator",
	}
}

// CoreInvalidExtGlobFilterMissingSeparatorError is the core error for invalid glob-ex filter
// definition when pattern is missing separator.
type CoreInvalidExtGlobFilterMissingSeparatorError struct {
	li18ngo.LocalisableError
}

var errCoreInvalidExtGlobFilterMissingSeparator = CoreInvalidExtGlobFilterMissingSeparatorError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreInvalidExtGlobFilterMissingSeparatorErrorTemplData{},
	},
}

// ❌ PolyFilterIsInvalid

// PolyFilterIsInvalidTemplData is the template data for the error when poly
// filter definition is invalid.
type PolyFilterIsInvalidTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td PolyFilterIsInvalidTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "poly-filter-is-invalid.age.static-error",
		Description: "poly filter definition is invalid error",
		Other:       "poly filter definition is invalid",
	}
}

// PolyFilterIsInvalidError is the error type for the error when poly filter
// definition is invalid.
type PolyFilterIsInvalidError struct {
	li18ngo.LocalisableError
}

// ErrPolyFilterIsInvalid is the exported error variable for
// PolyFilterIsInvalidError with the template data.
var ErrPolyFilterIsInvalid = PolyFilterIsInvalidError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterIsNilErrorTemplData{},
	},
}

// ❌ UsageMissingTreePath

// UsageMissingTreePathErrorTemplData is the template data for the error
// when tree path is missing.
type UsageMissingTreePathErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td UsageMissingTreePathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-tree-path.age.static-error",
		Description: "usage missing tree path",
		Other:       "usage missing tree path",
	}
}

// UsageMissingTreePathError is the error type for the error when usage is
// missing tree path.
type UsageMissingTreePathError struct {
	li18ngo.LocalisableError
}

// ErrUsageMissingTreePath is the exported error variable for
// UsageMissingTreePathError with the template data.
var ErrUsageMissingTreePath = UsageMissingTreePathError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UsageMissingTreePathErrorTemplData{},
	},
}

// ❌ UsageMissingRestorePath

// UsageMissingRestorePathErrorTemplData is the template data for the error when
// usage is missing restore path.
type UsageMissingRestorePathErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td UsageMissingRestorePathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-restore-path.age.static-error",
		Description: "usage missing restore path",
		Other:       "usage missing restore path",
	}
}

// UsageMissingRestorePathError is the error type for the error when usage is
// missing restore path.
type UsageMissingRestorePathError struct {
	li18ngo.LocalisableError
}

// ErrUsageMissingRestorePath is the exported error variable for UsageMissingRestorePathError
// with the template data.
var ErrUsageMissingRestorePath = UsageMissingRestorePathError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UsageMissingRestorePathErrorTemplData{},
	},
}

// ❌ UsageMissingSubscription

// UsageMissingSubscriptionErrorTemplData is the template data for the error when
// usage is missing subscription.
type UsageMissingSubscriptionErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td UsageMissingSubscriptionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-subscription.age.static-error",
		Description: "usage missing subscription",
		Other:       "usage missing subscription",
	}
}

// UsageMissingSubscriptionError is the error type for the error when usage is missing subscription.
type UsageMissingSubscriptionError struct {
	li18ngo.LocalisableError
}

// ErrUsageMissingSubscription is the exported error variable for UsageMissingSubscriptionError
// with the template data.
var ErrUsageMissingSubscription = UsageMissingSubscriptionError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UsageMissingSubscriptionErrorTemplData{},
	},
}

// ❌ UsageMissingHandler

// UsageMissingHandlerErrorTemplData is the template data for the error when usage is missing handler.
type UsageMissingHandlerErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td UsageMissingHandlerErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-handler.age.static-error",
		Description: "usage missing handler",
		Other:       "usage missing handler",
	}
}

// UsageMissingHandlerError is the error type for the error when usage
// is missing handler.
type UsageMissingHandlerError struct {
	li18ngo.LocalisableError
}

// ErrUsageMissingHandler is the exported error variable for
// UsageMissingHandlerError with the template data.
var ErrUsageMissingHandler = UsageMissingHandlerError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UsageMissingHandlerErrorTemplData{},
	},
}

// ❌ IDGeneratorFuncCantBeNil

// IDGeneratorFuncCantBeNilErrorTemplData is the template data for the error
// when id generator func is nil.
type IDGeneratorFuncCantBeNilErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td IDGeneratorFuncCantBeNilErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "id-generator-func-cant-be-nil.age.static-error",
		Description: "id generator func is nil, should be defined",
		Other:       "id generator func can't be nil",
	}
}

// IDGeneratorFuncCantBeNilError is the error type for the error when id generator func is nil.
type IDGeneratorFuncCantBeNilError struct {
	li18ngo.LocalisableError
}

// ErrIDGeneratorFuncCantBeNil is the exported error variable for IDGeneratorFuncCantBeNilError
// with the template data.
var ErrIDGeneratorFuncCantBeNil = IDGeneratorFuncCantBeNilError{
	LocalisableError: li18ngo.LocalisableError{
		Data: IDGeneratorFuncCantBeNilErrorTemplData{},
	},
}

// ❌ UnEqualJSONConversion

// UnEqualJSONConversionErrorTemplData is the template data for the error
// when JSON conversion results are not equal.
type UnEqualJSONConversionErrorTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td UnEqualJSONConversionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "un-equal-conversion.age.static-error",
		Description: "JSON options conversion error",
		Other:       "unequal JSON conversion",
	}
}

// UnEqualConversionError is the error type for the error when JSON conversion
// results are not equal.
type UnEqualConversionError struct {
	li18ngo.LocalisableError
}

// ErrUnEqualConversion is the exported error variable for UnEqualConversionError.
var ErrUnEqualConversion = UnEqualConversionError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UnEqualJSONConversionErrorTemplData{},
	},
}

// ❌ InvalidPath

// InvalidPathTemplData is the template data for the InvalidPath error message.
type InvalidPathTemplData struct {
	agenorTemplData
	Path string
}

// Message creates a new i18n message using the template data.
func (td InvalidPathTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-path.age.dynamic-error",
		Description: "invalid path (dynamic error)",
		Other:       "path: {{.Path}}",
	}
}

// TraverseFsMismatchTemplData 🍒 (dynamic i18n error)
type TraverseFsMismatchTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td TraverseFsMismatchTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "traverse-fs-mismatch.age.dynamic-error",
		Description: "traverse fs mismatch (dynamic error) (prefix for core mismatch error)",
		Other:       "traverse-fs",
	}
}

// TraverseFsMismatchError is the error type for the error when traverse fs
// mismatch occurs during traversal.
type TraverseFsMismatchError struct {
	li18ngo.LocalisableError
	Wrapped error
}

// Error returns the error message for TraverseFsMismatchError by combining the wrapped error
// message and the i18n message created from the template data.
func (e TraverseFsMismatchError) Error() string {
	return fmt.Sprintf("%v, %v", li18ngo.Text(e.Data), e.Wrapped.Error())
}

// Unwrap returns the wrapped error for TraverseFsMismatchError.
func (e TraverseFsMismatchError) Unwrap() error {
	return e.Wrapped
}

// NewTraverseFsMismatchError creates a new TraverseFsMismatchError and wraps
// the core error: CoreResumeFsMismatchError
func NewTraverseFsMismatchError() error {
	return &TraverseFsMismatchError{
		LocalisableError: li18ngo.LocalisableError{
			Data: TraverseFsMismatchTemplData{},
		},
		Wrapped: ErrCoreResumeFsMismatch,
	}
}

// ResumeFsMismatchTemplData 🍒 (dynamic i18n error)
type ResumeFsMismatchTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td ResumeFsMismatchTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "resume-fs-mismatch.age.dynamic-error",
		Description: "resume fs mismatch (dynamic error) (prefix for core mismatch error)",
		Other:       "resume-fs",
	}
}

// ResumeFsMismatchError is the error type for the error when resume fs mismatch
// occurs during resume.
type ResumeFsMismatchError struct {
	li18ngo.LocalisableError
	Wrapped error
}

// Error returns the error message for ResumeFsMismatchError by combining the wrapped error
// message and the i18n message created from the template data.
func (e ResumeFsMismatchError) Error() string {
	return fmt.Sprintf("%v, %v", li18ngo.Text(e.Data), e.Wrapped.Error())
}

// Unwrap returns the wrapped error for ResumeFsMismatchError.
func (e ResumeFsMismatchError) Unwrap() error {
	return e.Wrapped
}

// NewResumeFsMismatchError creates a new ResumeFsMismatchError and wraps
// the core error: CoreResumeFsMismatchError
func NewResumeFsMismatchError() error {
	return &ResumeFsMismatchError{
		LocalisableError: li18ngo.LocalisableError{
			Data: ResumeFsMismatchTemplData{},
		},
		Wrapped: ErrCoreResumeFsMismatch,
	}
}

// CoreResumeFsMismatchErrorTemplData 🥥 (core i18n error)
type CoreResumeFsMismatchErrorTemplData struct {
	agenorTemplData
}

// IsCoreResumeFsMismatchError uses errors.Is to check if the err's error tree
// contains the core error: CoreResumeFsMismatchError
func IsCoreResumeFsMismatchError(err error) bool {
	return errors.Is(err, ErrCoreResumeFsMismatch)
}

// Message creates a new i18n message using the template data.
func (td CoreResumeFsMismatchErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "core-resume-fs-mismatch.age.core-error",
		Description: "core resume file system mismatch error",
		Other:       "resume file system mismatch",
	}
}

// CoreResumeFsMismatchError is the core error for resume file system mismatch.
type CoreResumeFsMismatchError struct {
	li18ngo.LocalisableError
}

// ErrCoreResumeFsMismatch is the exported error variable for CoreResumeFsMismatchError
// with the template data.
var ErrCoreResumeFsMismatch = CoreResumeFsMismatchError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreResumeFsMismatchErrorTemplData{},
	},
}

// TraversalSavedTemplData 🍒 (dynamic i18n error)
type TraversalSavedTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td TraversalSavedTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "traversal-saved.age.dynamic-error",
		Description: "traversal saved due to panic (dynamic error)",
		Other:       "traversal saved: {{.Field}}",
	}
}

// TraversalSavedError is the error type for the error when traversal is
// saved due to panic.
type TraversalSavedError struct {
	li18ngo.LocalisableError
	Wrapped     error
	Destination string
}

// Error returns the error message for TraversalSavedError by combining the wrapped error
// message and the i18n message created from the template data.
func (e TraversalSavedError) Error() string {
	return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

// Unwrap returns the wrapped error for TraversalSavedError.
func (e TraversalSavedError) Unwrap() error {
	return e.Wrapped
}

// NewTraversalSavedError creates a new TraversalSavedError with the given destination and
// wraps the core error: ErrCorePanicOccurred
func NewTraversalSavedError(destination string, _ error) error {
	return &TraversalSavedError{
		LocalisableError: li18ngo.LocalisableError{
			Data: TraversalSavedTemplData{},
		},
		Wrapped:     ErrCorePanicOccurred,
		Destination: destination,
	}
}

// TraversalNotSavedTemplData 🍒 (dynamic i18n error)
type TraversalNotSavedTemplData struct {
	agenorTemplData
}

// Message creates a new i18n message using the template data.
func (td TraversalNotSavedTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "traversal-not-saved.age.dynamic-error",
		Description: "panic induced traversal not saved (dynamic error)",
		Other:       "field: {{.Reason}}",
	}
}

// TraversalNotSavedError is the error type for the error when panic
// induced traversal is not saved.
type TraversalNotSavedError struct {
	li18ngo.LocalisableError
	Wrapped error
	Reason  error
}

// Error returns the error message for TraversalNotSavedError by combining
// the wrapped error message, the reason error message, and the i18n message
// created from the template data.
func (e TraversalNotSavedError) Error() string {
	return fmt.Sprintf("%v, %v (%v)",
		e.Wrapped.Error(), li18ngo.Text(e.Data), e.Reason.Error(),
	)
}

// Unwrap returns the wrapped error for TraversalNotSavedError.
func (e TraversalNotSavedError) Unwrap() error {
	return e.Wrapped
}

// NewTraversalNotSavedError creates a new TraversalNotSavedError with the given
// reason error and wraps the core error: ErrCorePanicOccurred
func NewTraversalNotSavedError(_, reason error) error {
	return &TraversalNotSavedError{
		LocalisableError: li18ngo.LocalisableError{
			Data: TraversalNotSavedTemplData{},
		},
		Wrapped: ErrCorePanicOccurred,
		Reason:  reason,
	}
}

// CorePanicOccurredErrorTemplData 🥥 (core i18n error)
type CorePanicOccurredErrorTemplData struct {
	agenorTemplData
}

// IsCorePanicOccurredError uses errors.Is to check if the err's error tree
// contains the core error: CorePanicOccurredError
func IsCorePanicOccurredError(err error) bool {
	return errors.Is(err, ErrCorePanicOccurred)
}

// Message creates a new i18n message using the template data.
func (td CorePanicOccurredErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "core-panic-occurred.age.core-error",
		Description: "core error",
		Other:       "panic occurred",
	}
}

// CorePanicOccurredError is the core error for panic occurred.
type CorePanicOccurredError struct {
	li18ngo.LocalisableError
}

// ErrCorePanicOccurred is the exported error variable for CorePanicOccurredError
// with the template data.
var ErrCorePanicOccurred = CorePanicOccurredError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CorePanicOccurredErrorTemplData{},
	},
}
