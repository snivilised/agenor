package tapable

import (
	"github.com/snivilised/traverse/core"
)

type (
	Hooks struct {
		FileSubPath   Hook[core.SubPathHook]
		FolderSubPath Hook[core.SubPathHook]
		ReadDirectory Hook[core.ReadDirectoryHook]
		QueryStatus   Hook[core.QueryStatusHook]
		Sort          Hook[core.SortHook]
	}

	// HookCtrl contains the handler function to be invoked. The control
	// is agnostic to the handler's signature and therefore can not invoke it.
	HookCtrl[F any] struct {
		handler F
		def     F
	}
)

func NewHookCtrl[F any](handler F) *HookCtrl[F] {
	return &HookCtrl[F]{
		handler: handler,
		def:     handler,
	}
}

func (c *HookCtrl[F]) Tap(handler F) {
	c.handler = handler
}

func (c *HookCtrl[F]) Default() F {
	return c.def
}

func (c *HookCtrl[F]) Invoke() F {
	return c.handler
}
