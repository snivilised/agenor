package pref

import (
	"github.com/snivilised/traverse/core"
)

// WithHookCaseSensitiveSort specifies that a folder's contents
// should be sorted with case sensitivity.
func WithHookCaseSensitiveSort() Option {
	return func(o *Options) error {
		o.Hooks.Sort.Tap(CaseSensitiveSortHook)

		return nil
	}
}

// WithHookFileSubPath defines an custom hook to override the
// default behaviour for obtaining the sub-path of a file.
func WithHookFileSubPath(hook core.SubPathHook) Option {
	return func(o *Options) error {
		o.Hooks.FileSubPath.Tap(hook)

		return nil
	}
}

// WithHookFolderSubPath defines an custom hook to override the
// default behaviour for obtaining the sub-path of a folder.
func WithHookFolderSubPath(hook core.SubPathHook) Option {
	return func(o *Options) error {
		o.Hooks.FolderSubPath.Tap(hook)

		return nil
	}
}

// WithHookQueryStatus defines an custom hook to override the
// default behaviour for Stating a folder.
func WithHookQueryStatus(hook core.QueryStatusHook) Option {
	return func(o *Options) error {
		o.Hooks.QueryStatus.Tap(hook)

		return nil
	}
}

// WithHookReadDirectory defines an custom hook to override the
// default behaviour for reading a folder's contents.
func WithHookReadDirectory(hook core.ReadDirectoryHook) Option {
	return func(o *Options) error {
		o.Hooks.ReadDirectory.Tap(hook)

		return nil
	}
}

// WithHookSort defines an custom hook to override the
// default behaviour for sorting a folder's contents.
func WithHookSort(hook core.SortHook) Option {
	return func(o *Options) error {
		o.Hooks.Sort.Tap(hook)

		return nil
	}
}
