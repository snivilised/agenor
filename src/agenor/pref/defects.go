package pref

import (
	"io/fs"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
)

type (
	// NavigationFault represents an error that occurs during navigation. It
	// contains the error that indicates the nature of the fault, the file
	// system path that was being navigated when the fault occurred, and the
	// file information if available. This struct is used to provide context
	// when handling navigation errors.
	NavigationFault struct {
		// Err is the error that indicates the nature of the fault that
		// occurred during navigation.
		Err error

		// Path is the file system path that was being navigated when the
		// fault occurred.
		Path string

		// Info is the file information if available. This may be nil if the error
		// occurred before file information could be retrieved.
		Info fs.FileInfo
	}

	// PanicHandler is called to handle a panic. The handler specified allows custom
	// functionality.
	PanicHandler interface {
		Rescue(r Recovery, data RescueData) (string, error)
	}

	// Rescuer is a function type that implements the PanicHandler interface.
	Rescuer func(r Recovery, data RescueData) (string, error)

	// FaultHandler is called to handle an error that occurs when Stating
	// the tree directory. When an error occurs, traversal terminates
	// immediately. The handler specified allows custom functionality
	// when an error occurs here.
	FaultHandler interface {
		// Accept is called to handle a fault that occurs during navigation. The handler
		// specified allows custom functionality to be invoked when a fault occurs.
		Accept(fault *NavigationFault) error
	}

	// Accepter is a function type that implements the FaultHandler interface.
	Accepter func(fault *NavigationFault) error

	// SkipHandler is called to determine if traversal should be skipped when the
	// client callback returns an error during traversal. The client can control
	// if traversal is either terminated early (fs.SkipAll) or the remaining items
	// in a directory are skipped (fs.SkipDir).
	SkipHandler interface {
		// Ask is called to determine if traversal should be skipped when the client
		// callback returns an error during traversal. The client can control if
		// traversal is either terminated early (fs.SkipAll) or the remaining items
		// in a directory are skipped (fs.SkipDir).
		Ask(current *core.Node,
			contents core.DirectoryContents,
			err error,
		) (enums.SkipTraversal, error)
	}

	// Asker is a function type that implements the SkipHandler interface.
	Asker func(current *core.Node,
		contents core.DirectoryContents,
		err error,
	) (enums.SkipTraversal, error)

	// DefectOptions contains the handlers for handling faults, panics, and
	// skip decisions during traversal.
	DefectOptions struct {
		// Fault is the handler for handling errors that occur when Stating the tree directory.
		// When an error occurs, traversal terminates immediately. The handler specified
		// allows custom functionality to be invoked when an error occurs here.
		Fault FaultHandler

		// Panic is the handler for handling panics. The handler specified allows custom
		// functionality to be invoked when a panic occurs.
		Panic PanicHandler

		// Skip is the handler that will be invoked if the client callback returns an error
		// during traversal. The client can control if traversal is either terminated early
		// (fs.SkipAll) or the remaining items in a directory are skipped (fs.SkipDir).
		Skip SkipHandler
	}
)

// Accept is called to handle a fault that occurs during navigation. The handler
// specified allows custom functionality to be invoked when a fault occurs.
func (fn Accepter) Accept(fault *NavigationFault) error {
	return fn(fault)
}

// Rescue is called to handle a panic. The handler specified allows custom
// functionality to be invoked when a panic occurs.
func (fn Rescuer) Rescue(r Recovery, data RescueData) (string, error) {
	return fn(r, data)
}

// Ask is called to determine if traversal should be skipped when the client
// callback returns an error during traversal. The client can control if traversal
// is either terminated early (fs.SkipAll) or the remaining items in a directory
// are skipped (fs.SkipDir).
func (fn Asker) Ask(current *core.Node, contents core.DirectoryContents,
	err error,
) (enums.SkipTraversal, error) {
	return fn(current, contents, err)
}

// WithFaultHandler defines a custom handler to handle an error that occurs
// when 'Stat'ing the tree directory. When an error occurs, traversal terminates
// immediately. The handler specified allows custom functionality to be invoked
// when an error occurs here.
func WithFaultHandler(handler FaultHandler) Option {
	return func(o *Options) error {
		o.Defects.Fault = handler

		return nil
	}
}

// WithPanicHandler defines a custom handler to handle a panic.
func WithPanicHandler(handler PanicHandler) Option {
	return func(o *Options) error {
		o.Defects.Panic = handler

		return nil
	}
}

// WithSkipHandler defines a handler that will be invoked if the
// client callback returns an error during traversal. The client
// can control if traversal is either terminated early (fs.SkipAll)
// or the remaining items in a directory are skipped (fs.SkipDir).
func WithSkipHandler(handler SkipHandler) Option {
	return func(o *Options) error {
		o.Defects.Skip = handler

		return nil
	}
}
