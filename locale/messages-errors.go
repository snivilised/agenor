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

// InvalidFileSamplingSpecMissingFilesErrorTemplData
type InvalidFileSamplingSpecMissingFilesErrorTemplData struct {
	traverseTemplData
}

// Message
func (td InvalidFileSamplingSpecMissingFilesErrorTemplData) Message() *i18n.Message {
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
		ID:          "invalid-file-sampling-spec-missing-folders.static-error",
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
		ID:          "invalid-extended-glob-filter-missing-separator.dynamic-error",
		Description: "invalid extended glob filter definition; pattern is missing separator",
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
		Description: "invalid extended glob filter definition; pattern is missing separator",
		Other:       "invalid extended glob filter definition; pattern is missing separator",
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

type InvalidPathError struct {
	li18ngo.LocalisableError
	Wrapped error
}

func (e InvalidPathError) Error() string {
	wrapped := e.Wrapped.Error()
	path := li18ngo.Text(e.Data)
	return fmt.Sprintf("%v, %v", wrapped, path)
}

func (e InvalidPathError) Unwrap() error {
	return e.Wrapped
}

func NewInvalidPathError(path string) error {
	return &InvalidPathError{
		LocalisableError: li18ngo.LocalisableError{
			Data: InvalidPathTemplData{
				Path: path,
			},
		},
		Wrapped: errCoreInvalidPath,
	}
}

// ❌ CoreInvalidPathError

type CoreInvalidPathErrorTemplData struct {
	traverseTemplData
}

func IsInvalidPathError(err error) bool {
	return errors.Is(err, errCoreInvalidPath)
}

func (td CoreInvalidPathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-path.core-error",
		Description: "invalid path core error",
		Other:       "invalid path",
	}
}

type CoreInvalidPathError struct {
	li18ngo.LocalisableError
}

var errCoreInvalidPath = CoreInvalidPathError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreInvalidPathErrorTemplData{},
	},
}

// BOOKMARK

// ❌ InvalidBinaryFsOp, can be used to adapt the non i18n nefilim.InvalidBinaryFsOp
// error into an i18n one.

type InvalidBinaryFsOpTemplData struct {
	traverseTemplData
	From string
	To   string
	Op   string
}

func (td InvalidBinaryFsOpTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-binary-op.dynamic-error",
		Description: "invalid binary op dynamic error",
		Other:       "tbd: {{.Field}}",
	}
}

type InvalidBinaryFsOpError struct {
	li18ngo.LocalisableError
	Wrapped error
}

func (e InvalidBinaryFsOpError) Error() string {
	return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

func (e InvalidBinaryFsOpError) Unwrap() error {
	return e.Wrapped
}

func NewInvalidBinaryFsOpError(op, from, to string) error {
	return &InvalidBinaryFsOpError{
		LocalisableError: li18ngo.LocalisableError{
			Data: InvalidBinaryFsOpTemplData{
				From: from,
				To:   to,
				Op:   op,
			},
		},
		Wrapped: errCoreInvalidBinaryFsOp, // replace with nefilim version
	}
}

// ❌ CoreInvalidBinaryFsOp

type CoreInvalidBinaryFsOpErrorTemplData struct {
	traverseTemplData
}

func IsInvalidBinaryFsOpError(err error) bool {
	return errors.Is(err, errCoreInvalidBinaryFsOp)
}

func (td CoreInvalidBinaryFsOpErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-binary-op.core-error",
		Description: "invalid binary op core error",
		Other:       "invalid binary op",
	}
}

// will be replaced by nefilim.CoreInvalidBinaryFsOpError
type CoreInvalidBinaryFsOpError struct {
	li18ngo.LocalisableError
}

// this will be replaced by the nefilim.errCoreInvalidBinaryFsOp error
var errCoreInvalidBinaryFsOp = CoreInvalidBinaryFsOpError{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreInvalidBinaryFsOpErrorTemplData{},
	},
}

// ❌ RejectSameDirMoveError; replace with nefilim version

// RejectSameDirMoveErrorTemplDataL invalid file system operation
// that involves 2 paths, typically a source and destination.
// The error also indicates which operation is at fault.
type RejectSameDirMoveErrorTemplDataL struct {
	traverseTemplData
	From string
	To   string
	Op   string
}

// IsInvalidExtGlobFilterMissingSeparatorError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidExtGlobFilterMissingSeparatorError
func IsRejectSameDirMoveErrorL(err error) bool {
	return errors.Is(err, errCoreRejectSameDirMoveErrorL)
}

func NewRejectSameDirMoveErrorL(op, from, to string) error {
	return errors.Wrap(
		errCoreRejectSameDirMoveErrorL,
		li18ngo.Text(RejectSameDirMoveErrorTemplDataL{
			From: from,
			To:   to,
			Op:   op,
		}),
	)
}

// Message
func (td RejectSameDirMoveErrorTemplDataL) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "reject-same-dir-move.error",
		Description: "Reject same directory move operation as this is just a rename",
		Other:       "{{.Op}}, from: {{.From}}, to {{.To}}",
	}
}

// ❌ CoreRejectSameDirMoveError; replace with nefilim version

type CoreRejectSameDirMoveErrorTemplDataL struct {
	traverseTemplData
	li18ngo.LocalisableError
}

func (td CoreRejectSameDirMoveErrorTemplDataL) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "core-reject-same-dir-move.error",
		Description: "Core reject same directory move operation as this is just a rename",
		Other:       "same directory move rejected; use rename instead",
	}
}

// replace with nefilim version
var errCoreRejectSameDirMoveErrorL = CoreRejectSameDirMoveErrorTemplDataL{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreRejectSameDirMoveErrorTemplDataL{},
	},
}
