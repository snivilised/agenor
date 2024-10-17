package pref

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
)

type (
	// NavigationFault
	NavigationFault struct {
		Err  error
		Path string
		Info fs.FileInfo
	}

	// PanicHandler
	PanicHandler interface {
		Rescue()
	}

	Rescuer func()

	// FaultHandler is called to handle an error that occurs when Stating
	// the tree folder. When an error occurs, traversal terminates
	// immediately. The handler specified allows custom functionality
	// when an error occurs here.
	FaultHandler interface {
		Accept(fault *NavigationFault) error
	}

	Accepter func(fault *NavigationFault) error

	// SkipHandler
	SkipHandler interface {
		Ask(current *core.Node,
			contents core.DirectoryContents,
			err error,
		) (enums.SkipTraversal, error)
	}

	Asker func(current *core.Node,
		contents core.DirectoryContents,
		err error,
	) (enums.SkipTraversal, error)

	// DefectOptions
	DefectOptions struct {
		Fault FaultHandler
		Panic PanicHandler
		Skip  SkipHandler
	}
)

func (fn Accepter) Accept(fault *NavigationFault) error {
	return fn(fault)
}

func (fn Rescuer) Rescue() {
	fn()
}

func (fn Asker) Ask(current *core.Node, contents core.DirectoryContents, err error) (enums.SkipTraversal, error) {
	return fn(current, contents, err)
}

// WithFaultHandler defines a custom handler to handle an error that occurs
// when 'Stat'ing the tree folder. When an error occurs, traversal terminates
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
