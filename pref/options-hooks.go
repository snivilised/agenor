package pref

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/lo"
)

func WithHookQueryStatus(hook core.QueryStatusHook) Option {
	return func(o *Options) error {
		o.Hooks.QueryStatus.Tap(hook)

		return nil
	}
}

func WithHookReadDirectory(hook core.ReadDirectoryHook) Option {
	return func(o *Options) error {
		o.Hooks.ReadDirectory.Tap(hook)

		return nil
	}
}

func WithHookSortCase(icase bool) Option {
	return func(o *Options) error {
		hook := lo.Ternary(icase,
			CaseInSensitiveSortHook,
			CaseSensitiveSortHook,
		)
		o.Hooks.Sort.Tap(hook)

		return nil
	}
}

func WithHookSort(hook core.SortHook) Option {
	return func(o *Options) error {
		o.Hooks.Sort.Tap(hook)

		return nil
	}
}

func WithHookFileSubPath(hook core.SubPathHook) Option {
	return func(o *Options) error {
		o.Hooks.FileSubPath.Tap(hook)

		return nil
	}
}

func WithHookFolderSubPath(hook core.SubPathHook) Option {
	return func(o *Options) error {
		o.Hooks.FolderSubPath.Tap(hook)

		return nil
	}
}
