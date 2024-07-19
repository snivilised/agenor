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

	// FaultHandler
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

func WithFaultHandler(handler FaultHandler) Option {
	return func(o *Options) error {
		o.Defects.Fault = handler

		return nil
	}
}

func WithPanicHandler(handler PanicHandler) Option {
	return func(o *Options) error {
		o.Defects.Panic = handler

		return nil
	}
}

func WithSkipHandler(handler SkipHandler) Option {
	return func(o *Options) error {
		o.Defects.Skip = handler

		return nil
	}
}
