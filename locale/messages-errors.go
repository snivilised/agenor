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
	agenorTemplData
}

// Message
func (td FilterIsNilErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "filter-is-nil.age.static-error",
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
	agenorTemplData
}

// Message
func (td FilterMissingTypeErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "filter-missing-type.age.static-error",
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
	agenorTemplData
}

// Message
func (td FilterCustomNotSupportedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "custom-filter-not-supported-for-children.age.static-error",
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
	agenorTemplData
}

// Message
func (td FilterUndefinedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "filter-is-undefined.age.static-error",
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
	agenorTemplData
}

// Message
func (td InternalFailedToGetNavigatorDriverErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "failed-to-get-navigator-driver.age.static-error",
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
	agenorTemplData
	Pattern string
}

func (td InvalidIncaseFilterDefTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-incase-filter-definition.age.dynamic-error",
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
	agenorTemplData
}

func IsInvalidIncaseFilterDefError(err error) bool {
	return errors.Is(err, errCoreInvalidIncaseFilterDef)
}

func (td CoreInvalidIncaseFilterDefErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-incase-filter-definition.age.core-error",
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
	agenorTemplData
}

// Message
func (td WorkerPoolCreationFailedErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "failed-to-create-worker-pool.age.static-error",
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
	agenorTemplData
}

// Message
func (td InvalidSamplingSpecMissingFilesErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-file-sampling-spec-missing-files.age.static-error",
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
	agenorTemplData
}

// Message
func (td InvalidSamplingSpecMissingDirectoriesErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-file-sampling-spec-missing-directories.age.static-error",
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
	agenorTemplData
}

// Message
func (td MissingCustomFilterDefinitionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "missing-custom-filter-definition.age.static-error",
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
	agenorTemplData
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
	agenorTemplData
	Pattern string
}

// Message
func (td InvalidExtGlobFilterMissingSeparatorErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-glob-ex-filter-missing-separator.age.dynamic-error",
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
	agenorTemplData
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
		ID:          "invalid-extended-glob-filter-missing-separator.age.core-error",
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
	agenorTemplData
}

// Message
func (td PolyFilterIsInvalidTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "poly-filter-is-invalid.age.static-error",
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
	agenorTemplData
}

// Message
func (td UsageMissingTreePathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-tree-path.age.static-error",
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
	agenorTemplData
}

// Message
func (td UsageMissingRestorePathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-restore-path.age.static-error",
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
	agenorTemplData
}

// Message
func (td UsageMissingSubscriptionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-subscription.age.static-error",
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
	agenorTemplData
}

// Message
func (td UsageMissingHandlerErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-handler.age.static-error",
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
	agenorTemplData
}

// Message
func (td IDGeneratorFuncCantBeNilErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "id-generator-func-cant-be-nil.age.static-error",
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
	agenorTemplData
}

// Message
func (td UnEqualJSONConversionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "un-equal-conversion.age.static-error",
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
	agenorTemplData
	Path string
}

func (td InvalidPathTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-path.age.dynamic-error",
		Description: "invalid path (dynamic error)",
		Other:       "path: {{.Path}}",
	}
}

// 🍒 ResumeTraverseFsMismatch (dynamic i18n error)
type TraverseFsMismatchTemplData struct {
	agenorTemplData
}

func (td TraverseFsMismatchTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "traverse-fs-mismatch.age.dynamic-error",
		Description: "traverse fs mismatch (dynamic error) (prefix for core mismatch error)",
		Other:       "traverse-fs",
	}
}

type TraverseFsMismatchError struct {
	li18ngo.LocalisableError
	Wrapped error
}

func (e TraverseFsMismatchError) Error() string {
	return fmt.Sprintf("%v, %v", li18ngo.Text(e.Data), e.Wrapped.Error())
}

func (e TraverseFsMismatchError) Unwrap() error {
	return e.Wrapped
}

func NewTraverseFsMismatchError() error {
	return &TraverseFsMismatchError{
		LocalisableError: li18ngo.LocalisableError{
			Data: TraverseFsMismatchTemplData{},
		},
		Wrapped: ErrCoreResumeFsMismatch,
	}
}

// 🍒 ResumeFsMismatch (dynamic i18n error)
type ResumeFsMismatchTemplData struct {
	agenorTemplData
}

func (td ResumeFsMismatchTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "resume-fs-mismatch.age.dynamic-error",
		Description: "resume fs mismatch (dynamic error) (prefix for core mismatch error)",
		Other:       "resume-fs",
	}
}

type ResumeFsMismatchError struct {
	li18ngo.LocalisableError
	Wrapped error
}

func (e ResumeFsMismatchError) Error() string {
	return fmt.Sprintf("%v, %v", li18ngo.Text(e.Data), e.Wrapped.Error())
}

func (e ResumeFsMismatchError) Unwrap() error {
	return e.Wrapped
}

func NewResumeFsMismatchError() error {
	return &ResumeFsMismatchError{
		LocalisableError: li18ngo.LocalisableError{
			Data: ResumeFsMismatchTemplData{},
		},
		Wrapped: ErrCoreResumeFsMismatch,
	}
}

// 🥥 CoreResumeFsMismatch (core i18n error)
type CoreResumeFsMismatchErrorTemplData struct {
	agenorTemplData
}

func IsCoreResumeFsMismatchError(err error) bool {
	return errors.Is(err, ErrCoreResumeFsMismatch)
}

func (td CoreResumeFsMismatchErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "core-resume-fs-mismatch.age.core-error",
		Description: "core resume file system mismatch error",
		Other:       "resume file system mismatch",
	}
}

type CoreResumeFsMismatchError struct {
	li18ngo.LocalisableError
}

var ErrCoreResumeFsMismatch = CoreResumeFsMismatchError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreResumeFsMismatchErrorTemplData{},
	},
}

// 🍒 TraversalSaved (dynamic i18n error)
type TraversalSavedTemplData struct {
	agenorTemplData
}

func (td TraversalSavedTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "traversal-saved.age.dynamic-error",
		Description: "traversal saved due to panic (dynamic error)",
		Other:       "traversal saved: {{.Field}}",
	}
}

type TraversalSavedError struct {
	li18ngo.LocalisableError
	Wrapped     error
	Destination string
}

func (e TraversalSavedError) Error() string {
	return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

func (e TraversalSavedError) Unwrap() error {
	return e.Wrapped
}

func NewTraversalSavedError(destination string, _ error) error {
	return &TraversalSavedError{
		LocalisableError: li18ngo.LocalisableError{
			Data: TraversalSavedTemplData{},
		},
		Wrapped:     ErrCorePanicOccurred,
		Destination: destination,
	}
}

// 🍒 TraversalNotSaved (dynamic i18n error)
type TraversalNotSavedTemplData struct {
	agenorTemplData
}

func (td TraversalNotSavedTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "traversal-not-saved.age.dynamic-error",
		Description: "panic induced traversal not saved (dynamic error)",
		Other:       "field: {{.Reason}}",
	}
}

type TraversalNotSavedError struct {
	li18ngo.LocalisableError
	Wrapped error
	Reason  error
}

func (e TraversalNotSavedError) Error() string {
	return fmt.Sprintf("%v, %v (%v)",
		e.Wrapped.Error(), li18ngo.Text(e.Data), e.Reason.Error(),
	)
}

func (e TraversalNotSavedError) Unwrap() error {
	return e.Wrapped
}

func NewTraversalNotSavedError(_, reason error) error {
	return &TraversalNotSavedError{
		LocalisableError: li18ngo.LocalisableError{
			Data: TraversalNotSavedTemplData{},
		},
		Wrapped: ErrCorePanicOccurred,
		Reason:  reason,
	}
}

// 🥥 CorePanicOccurred (core i18n error)
type CorePanicOccurredErrorTemplData struct {
	agenorTemplData
}

func IsCorePanicOccurredError(err error) bool {
	return errors.Is(err, ErrCorePanicOccurred)
}

func (td CorePanicOccurredErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "core-panic-occurred.age.core-error",
		Description: "core error",
		Other:       "panic occurred",
	}
}

type CorePanicOccurredError struct {
	li18ngo.LocalisableError
}

var ErrCorePanicOccurred = CorePanicOccurredError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CorePanicOccurredErrorTemplData{},
	},
}
