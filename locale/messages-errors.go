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

// ❌ PolyFilterIsInvalid

// FilterIsNilTemplData
type PolyFilterIsInvalidTemplData struct {
	traverseTemplData
}

// Message
func (td PolyFilterIsInvalidTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "poly-filter-is-invalid.error",
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

// ❌ UsageMissingRootPath

// UsageMissingRootPathTemplData
type UsageMissingTreePathErrorTemplData struct {
	traverseTemplData
}

// Message
func (td UsageMissingTreePathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "usage-missing-tree-path.error",
		Description: "usage missing tree path",
		Other:       "usage missing tree path",
	}
}

type UsageMissingRootPathError struct {
	li18ngo.LocalisableError
}

var ErrUsageMissingRootPath = UsageMissingRootPathError{
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

// ❌ UnEqualJSONConversion

// UnEqualConversionTemplData
type UnEqualJSONConversionErrorTemplData struct {
	traverseTemplData
}

// Message
func (td UnEqualJSONConversionErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "un-equal-conversion.error",
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

// InvalidPathErrorTemplData invalid file system path; path must be relative
// relative to the root already defined for this file system so should not
// start or end with a /. Also, only a / should be used to denote a separator
// which applies to all platforms.
type InvalidPathErrorTemplData struct {
	traverseTemplData
	Path string
}

// IsInvalidExtGlobFilterMissingSeparatorError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidExtGlobFilterMissingSeparatorError
func IsInvalidPathError(err error) bool {
	return errors.Is(err, errInvalidPath)
}

func NewInvalidPathError(path string) error {
	return errors.Wrap(
		errInvalidPath,
		li18ngo.Text(InvalidPathErrorTemplData{
			Path: path,
		}),
	)
}

// Message
func (td InvalidPathErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-path.error",
		Description: "Invalid file system path",
		Other:       "invalid path {{.Path}}",
	}
}

type InvalidPathError struct {
	li18ngo.LocalisableError
}

var errInvalidPath = InvalidPathError{
	LocalisableError: li18ngo.LocalisableError{
		Data: InvalidPathErrorTemplData{},
	},
}

// ❌ InvalidBinaryFsOp

// InvalidBinaryFsOpErrorTemplData invalid file system operation
// that involves 2 paths, typically a source and destination.
// The error also indicates which operation is at fault.
type InvalidBinaryFsOpErrorTemplData struct {
	traverseTemplData
	From string
	To   string
	Op   string
}

// IsInvalidExtGlobFilterMissingSeparatorError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidExtGlobFilterMissingSeparatorError
func IsBinaryFsOpError(err error) bool {
	return errors.Is(err, errCoreBinaryFsOp)
}

func NewInvalidBinaryFsOpError(op, from, to string) error {
	return errors.Wrap(
		errCoreBinaryFsOp,
		li18ngo.Text(InvalidBinaryFsOpErrorTemplData{
			From: from,
			To:   to,
			Op:   op,
		}),
	)
}

// Message
func (td InvalidBinaryFsOpErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "invalid-binary-op.error",
		Description: "Invalid file system operation",
		Other:       "{{.Op}}, from: {{.From}}, to {{.To}}",
	}
}

// ❌ CoreBinaryFsOpError

type CoreBinaryFsOpErrorTemplData struct {
	traverseTemplData
	li18ngo.LocalisableError
}

func (td CoreBinaryFsOpErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "core-invalid-binary-op.error",
		Description: "Core Invalid file system operation",
		Other:       "invalid file system operation",
	}
}

var errCoreBinaryFsOp = CoreBinaryFsOpErrorTemplData{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreBinaryFsOpErrorTemplData{},
	},
}

// ❌ RejectSameDirMoveError

// RejectSameDirMoveErrorTemplData invalid file system operation
// that involves 2 paths, typically a source and destination.
// The error also indicates which operation is at fault.
type RejectSameDirMoveErrorTemplData struct {
	traverseTemplData
	From string
	To   string
	Op   string
}

// IsInvalidExtGlobFilterMissingSeparatorError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidExtGlobFilterMissingSeparatorError
func IsRejectSameDirMoveError(err error) bool {
	return errors.Is(err, errCoreRejectSameDirMoveError)
}

func NewRejectSameDirMoveError(op, from, to string) error {
	return errors.Wrap(
		errCoreRejectSameDirMoveError,
		li18ngo.Text(RejectSameDirMoveErrorTemplData{
			From: from,
			To:   to,
			Op:   op,
		}),
	)
}

// Message
func (td RejectSameDirMoveErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "reject-same-dir-move.error",
		Description: "Reject same directory move operation as this is just a rename",
		Other:       "{{.Op}}, from: {{.From}}, to {{.To}}",
	}
}

// ❌ CoreRejectSameDirMoveError

type CoreRejectSameDirMoveErrorTemplData struct {
	traverseTemplData
	li18ngo.LocalisableError
}

func (td CoreRejectSameDirMoveErrorTemplData) Message() *i18n.Message {
	return &i18n.Message{
		ID:          "core-reject-same-dir-move.error",
		Description: "Core reject same directory move operation as this is just a rename",
		Other:       "same directory move rejected; use rename instead",
	}
}

var errCoreRejectSameDirMoveError = CoreRejectSameDirMoveErrorTemplData{
	LocalisableError: li18ngo.LocalisableError{
		Data: CoreRejectSameDirMoveErrorTemplData{},
	},
}
