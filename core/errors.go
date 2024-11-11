package core

import (
	"fmt"

	"github.com/pkg/errors"
)

// Errors defined here are internal errors that are of no value to end
// users (hence not l10n). There are usually programming errors which
// means they only have meaning for client developers as opposed to the
// end user.
//
// A Variable error is any error that contains extra state that
// is used to decorate the underlying core error. So for example
// 'InvalidNotificationMuteRequestedError' has a 'notification' value
// associated with it, that is required when creating this error. The
// core error (unexported), is decorated by the value to create a custom
// version of the core. But because of the way wrapping works, when we
// display the error via the Error method, the information unfortunately
// comes out backwards. So creating an instance of
// 'InvalidNotificationMuteRequestedError', with a notification  value
// 'OnBegin' is displayed as:
// notification: OnBegin: invalid notification mute requested
// So what this is saying is that the underlying error is
// 'invalid notification mute requested', but the detail is that the
// notification value is 'OnBegin'. Intuitively, we would expect to see
// this arranged the other way around, but this is not possible, without
// an unjustifiable amount of re-work.
//
// Any non variable errors, eg 'ErrGuardianCantDecorateItemSealed' do not
// have an IsXXX function defined, because the client can invoke errors.Is
// on the error directly themselves;
// ie: errors.Is(err, ErrGuardianCantDecorateItemSealed)
//
// These errors are generated from snippet "Variable Native Error" (verr)
//

// ❌ InvalidNotificationMuteRequested error

// NewInvalidNotificationMuteRequestedError creates an untranslated error to
// indicate invalid notification mute requested (internal error)
func NewInvalidNotificationMuteRequestedError(notification string) error {
	return errors.Wrap(
		errInvalidNotificationMuteRequested,
		fmt.Sprintf("notification: %v", notification),
	)
}

// IsInvalidNotificationMuteRequestedError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidNotificationMuteRequestedError
func IsInvalidNotificationMuteRequestedError(err error) bool {
	return errors.Is(err, errInvalidNotificationMuteRequested)
}

var errInvalidNotificationMuteRequested = errors.New(
	"invalid notification mute requested",
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

// ❌ NewItemAlreadyExtended error

// NewItemAlreadyExtendedError creates an untranslated error to
// indicate the node has already been extended
func NewItemAlreadyExtendedError(path string) error {
	return errors.Wrap(
		errNewItemAlreadyExtended,
		fmt.Sprintf("path: %v", path),
	)
}

// IsNewItemAlreadyExtendedError uses errors.Is to check
// if the err's error tree contains the core error:
// NewItemAlreadyExtendedError
func IsNewItemAlreadyExtendedError(err error) bool {
	return errors.Is(err, errNewItemAlreadyExtended)
}

var errNewItemAlreadyExtended = errors.New(
	"item already extended for item",
)

// ❌ MissingHibernationDetacherFunction error

// NewMissingHibernationDetacherFunctionError creates an untranslated error to
// indicate hibernation detacher function nt defined
func NewMissingHibernationDetacherFunctionError(state string) error {
	return errors.Wrap(
		errMissingListenDetacherFunction,
		fmt.Sprintf("state: %v", state),
	)
}

// IsMissingHibernationDetacherFunctionError uses errors.Is to check
// if the err's error tree contains the core error:
// MissingListenDetacherFunctionError
func IsMissingHibernationDetacherFunctionError(err error) bool {
	return errors.Is(err, errMissingListenDetacherFunction)
}

var errMissingListenDetacherFunction = errors.New(
	"missing listen detacher function",
)

// ❌ InvalidPeriscopeRootPath error

// NewInvalidPeriscopeRootPathError creates an untranslated error to
// indicate invalid periscope tree path, ie the current path is
// not a child directory relative to the tree path.
func NewInvalidPeriscopeRootPathError(tree, current string) error {
	return errors.Wrap(
		errInvalidPeriscopeTreePath,
		fmt.Sprintf("tree: '%v', current: '%v'", tree, current),
	)
}

// IsInvalidPeriscopeRootPathError uses errors.Is to check
// if the err's error tree contains the core error:
// InvalidPeriscopeRootPathError
func IsInvalidPeriscopeRootPathError(err error) bool {
	return errors.Is(err, errInvalidPeriscopeTreePath)
}

var errInvalidPeriscopeTreePath = errors.New(
	"tree path can't be longer than current",
)

// ❌ ResumeControllerNotSet error

// NewResumeControllerNotSetError creates an untranslated error to
// indicate resume controller not set
func NewResumeControllerNotSetError(from string) error {
	return errors.Wrap(
		errResumeControllerNotSet,
		fmt.Sprintf("from: %v", from),
	)
}

// IsResumeControllerNotSetError uses errors.Is to check
// if the err's error tree contains the core error:
// ResumeControllerNotSetError
func IsResumeControllerNotSetError(err error) bool {
	return errors.Is(err, errResumeControllerNotSet)
}

var errResumeControllerNotSet = errors.New(
	"resume controller not set",
)

// ❌ GuardianCantDecorateItemSealed error

// ErrGuardianCantDecorateItemSealedError creates an untranslated error to
// indicate last item is sealed
var ErrGuardianCantDecorateItemSealed = errors.New(
	"can't decorate, last item is sealed",
)

// ❌ BrokerTopicNotFound error

// NewBrokerTopicNotFoundError creates an untranslated error to
// indicate the topic requested was not found; ie it wasn't registered.
func NewBrokerTopicNotFoundError(topic string) error {
	return errors.Wrap(
		errBrokerTopicNotFound,
		fmt.Sprintf("topic: %v", topic),
	)
}

// IsBrokerTopicNotFoundError uses errors.Is to check
// if the err's error tree contains the core error:
// BrokerTopicNotFoundError
func IsBrokerTopicNotFoundError(err error) bool {
	return errors.Is(err, errBrokerTopicNotFound)
}

var errBrokerTopicNotFound = errors.New(
	"broker topic not found",
)

// ❌ DetectedSpawnStackOverflow error

var ErrDetectedSpawnStackOverflow = errors.New(
	"spawn resume stack-overflow protection",
)

// ❌ WrongPrimaryFacade error

var ErrWrongPrimaryFacade = errors.New(
	"wrong primary facade",
)

// ❌ WrongResumeFacade error

var ErrWrongResumeFacade = errors.New(
	"wrong resume facade",
)

// ❌ Nil Forest error

var ErrNilForest = errors.New(
	"forest is nil",
)
