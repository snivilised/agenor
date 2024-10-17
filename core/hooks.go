package core

import (
	"io/fs"
)

type (
	// QueryStatusHook function signature that enables the default to be overridden.
	// (By default, uses Lstat)
	QueryStatusHook func(qsys fs.StatFS, path string) (fs.FileInfo, error)

	// ChainQueryStatusHook chainable version of QueryStatusHook
	ChainQueryStatusHook func(result fs.FileInfo, err error,
		qsys fs.StatFS, path string,
	) (fs.FileInfo, error)

	// ReadDirectoryHook hook function to define implementation of how a directory's
	// entries are read. A default implementation is preset, so does not have to be set
	// by the client.
	ReadDirectoryHook func(rsys fs.ReadDirFS, dirname string) ([]fs.DirEntry, error)

	// ChainReadDirectoryHook chainable version of
	ChainReadDirectoryHook func(result []fs.DirEntry, err error,
		rsys fs.ReadDirFS, dirname string,
	) ([]fs.DirEntry, error)

	// SortHook hook function to define how directory entries are sorted. Does not
	// have to be set explicitly. This will be set according to the IsCaseSensitive on
	// the TraverseOptions, but can be overridden if needed.
	SortHook func(entries []fs.DirEntry, custom ...any)

	// ChainSortHook chainable version of SortHook
	ChainSortHook func(
		entries []fs.DirEntry, custom ...any,
	)

	SubPathInfo struct {
		Tree            string
		Node            *Node
		KeepTrailingSep bool
	}

	// SubPathHook
	SubPathHook func(info *SubPathInfo) string

	// ChainSubPathHook chainable version of SubPathHook
	ChainSubPathHook func(result string,
		info *SubPathInfo,
	) string
)
