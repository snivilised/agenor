package locale

import (
	"fmt"

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
		ID:          "filter-is-nil.static-error",
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
		ID:          "filter-missing-type.static-error",
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
		ID:          "custom-filter-not-supported-for-children.static-error",
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
		ID:          "filter-is-undefined.static-error",
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
		ID:          "failed-to-get-navigator-driver.static-error",
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

// ❌ InvalidIncaseFilterDefError

type InvalidIncaseFilterDefTemplData struct {
	traverseTemplData
	Pattern string
}

func (td InvalidIncaseFilterDefTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-incase-filter-definition.dynamic-error",
		Description: "invalid incase filter definition; pattern is missing separator wrapper error",
		Other:       "pattern: {{.Pattern}}",
	}
}

type InvalidIncaseFilterDefError struct {
	li18ngo.LocalisableError
	Wrapped error
}

func (e InvalidIncaseFilterDefError) Error() string {
	return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

func (e InvalidIncaseFilterDefError) Unwrap() error {
	return e.Wrapped
}

func NewInvalidIncaseFilterDefError(pattern string) error {
	return &InvalidIncaseFilterDefError{
		LocalisableError: li18ngo.LocalisableError{
			Data: InvalidIncaseFilterDefTemplData{
				Pattern: pattern,
			},
		},
		Wrapped: errCoreInvalidIncaseFilterDef,
	}
}

type CoreInvalidIncaseFilterDefErrorTemplData struct {
	traverseTemplData
}

func IsInvalidIncaseFilterDefError(err error) bool {
	return errors.Is(err, errCoreInvalidIncaseFilterDef)
}

func (td CoreInvalidIncaseFilterDefErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-incase-filter-definition.core-error",
		Description: "invalid incase filter definition; pattern is missing separator core error",
		Other:       "invalid incase filter definition; pattern is missing separator",
	}
}

type CoreInvalidIncaseFilterDefError struct {
	li18ngo.LocalisableError
}

var errCoreInvalidIncaseFilterDef = CoreInvalidIncaseFilterDefError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreInvalidIncaseFilterDefErrorTemplData{},
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
		ID:          "failed-to-create-worker-pool.static-error",
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

// InvalidSamplingSpecMissingFilesErrorTemplData
type InvalidSamplingSpecMissingFilesErrorTemplData struct {
	traverseTemplData
}

// Message
func (td InvalidSamplingSpecMissingFilesErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-file-sampling-spec-missing-files.static-error",
		Description: "invalid file sampling specification, missing no of files",
		Other:       "invalid file sampling specification, missing no of files",
	}
}

type InvalidFileSamplingSpecificationError struct {
	li18ngo.LocalisableError
}

var ErrInvalidFileSamplingSpecMissingFiles = InvalidFileSamplingSpecificationError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidSamplingSpecMissingFilesErrorTemplData{},
	},
}

// ❌ InvalidSamplingSpecMissingDirectoriesError

// InvalidSamplingSpecMissingDirectoriesErrorTemplData
type InvalidSamplingSpecMissingDirectoriesErrorTemplData struct {
	traverseTemplData
}

// Message
func (td InvalidSamplingSpecMissingDirectoriesErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-file-sampling-spec-missing-directories.static-error",
		Description: "invalid file sampling specification, missing no of directories",
		Other:       "invalid file sampling specification, missing no of directories",
	}
}

type InvalidSamplingSpecMissingDirectoriesError struct {
	li18ngo.LocalisableError
}

var ErrInvalidSamplingSpecMissingDirectories = InvalidSamplingSpecMissingDirectoriesError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidSamplingSpecMissingDirectoriesErrorTemplData{},
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
		ID:          "missing-custom-filter-definition.static-error",
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
	Pattern string
}

// Message
func (td InvalidExtGlobFilterMissingSeparatorErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-glob-ex-filter-missing-separator.dynamic-error",
		Description: "invalid glob ex filter definition; pattern is missing separator",
		Other:       "pattern: {{.Pattern}}",
	}
}

type InvalidExtGlobFilterMissingSeparatorError struct {
	li18ngo.LocalisableError
	Wrapped error
}

func (e InvalidExtGlobFilterMissingSeparatorError) Error() string {
	return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

func (e InvalidExtGlobFilterMissingSeparatorError) Unwrap() error {
	return e.Wrapped
}

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

// InvalidExtGlobFilterMissingSeparatorTemplData
type CoreInvalidExtGlobFilterMissingSeparatorErrorTemplData struct {
	traverseTemplData
}

// IsInvalidExtGlobFilterMissingSeparatorError uses errors.Is to check
// if the err's error tree contains the core error:
// CoreInvalidExtGlobFilterMissingSeparatorError
func IsInvalidExtGlobFilterMissingSeparatorError(err error) bool {
	return errors.Is(err, errCoreInvalidExtGlobFilterMissingSeparator)
}

// Message
func (td CoreInvalidExtGlobFilterMissingSeparatorErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-extended-glob-filter-missing-separator.core-error",
		Description: "invalid glob ex filter definition; pattern is missing separator",
		Other:       "invalid glob ex filter definition; pattern is missing separator",
	}
}

type CoreInvalidExtGlobFilterMissingSeparatorError struct {
	li18ngo.LocalisableError
}

var errCoreInvalidExtGlobFilterMissingSeparator = CoreInvalidExtGlobFilterMissingSeparatorError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreInvalidExtGlobFilterMissingSeparatorErrorTemplData{},
	},
}

// ❌ PolyFilterIsInvalid

// FilterIsNilTemplData
type PolyFilterIsInvalidTemplData struct {
	traverseTemplData
}

// Message
func (td PolyFilterIsInvalidTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "poly-filter-is-invalid.static-error",
		Description: "poly filter definition is invalid error",
		Other:       "poly filter definition is invalid",
	}
}

type PolyFilterIsInvalidError struct {
	li18ngo.LocalisableError
}

var ErrPolyFilterIsInvalid = PolyFilterIsInvalidError{
	LocalisableError: li18ngo.LocalisableError{
		Data: FilterIsNilErrorTemplData{},
	},
}

// ❌ UsageMissingTreePath

// UsageMissingRootPathTemplData
type UsageMissingTreePathErrorTemplData struct {
	traverseTemplData
}

// Message
func (td UsageMissingTreePathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-tree-path.static-error",
		Description: "usage missing tree path",
		Other:       "usage missing tree path",
	}
}

type UsageMissingTreePathError struct {
	li18ngo.LocalisableError
}

var ErrUsageMissingTreePath = UsageMissingTreePathError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UsageMissingTreePathErrorTemplData{},
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
		ID:          "usage-missing-restore-path.static-error",
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
		ID:          "usage-missing-subscription.static-error",
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
		ID:          "usage-missing-handler.static-error",
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
		ID:          "id-generator-func-cant-be-nil.static-error",
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

// ❌ UnEqualJSONConversion

// UnEqualConversionTemplData
type UnEqualJSONConversionErrorTemplData struct {
	traverseTemplData
}

// Message
func (td UnEqualJSONConversionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "un-equal-conversion.static-error",
		Description: "JSON options conversion error",
		Other:       "unequal JSON conversion",
	}
}

type UnEqualConversionError struct {
	li18ngo.LocalisableError
}

var ErrUnEqualConversion = UnEqualConversionError{
	LocalisableError: li18ngo.LocalisableError{
		Data: UnEqualJSONConversionErrorTemplData{},
	},
}

// ❌ InvalidPath

type InvalidPathTemplData struct {
	traverseTemplData
	Path string
}

func (td InvalidPathTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-path.dynamic-error",
		Description: "invalid path dynamic error",
		Other:       "path: {{.Path}}",
	}
}
