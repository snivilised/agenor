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

	// SubPathInfo represents the information about a subpath during traversal, including
	// the tree it belongs to, the current node being processed, and whether to keep
	// the trailing separator. This information can be used by the SubPathHook to
	// generate a subpath string based on the current traversal context.
	SubPathInfo struct {
		// Tree represents the tree to which the current node belongs. This can be used
		// to provide context about the file system being traversed and can be useful for
		// generating sub-paths that are relevant to the specific tree.
		Tree string

		// Node represents the current node being processed during traversal. This can be used
		// to access information about the current file or directory, such as its name, path,
		// and other attributes, which can be useful for generating sub-paths that are specific
		// to the current node.
		Node *Node

		// KeepTrailingSep indicates whether to keep the trailing separator in the generated
		// subpath. This can be useful for distinguishing between directories and files in the
		// generated sub-paths, as directories typically have a trailing separator while files
		// do not.
		KeepTrailingSep bool
	}

	// SubPathHook function signature for generating sub-paths during traversal. It takes a
	// SubPathInfo struct as input and returns a string representing the generated sub-path.
	// This hook can be used to customize how sub-paths are generated based on the current
	// traversal context, allowing for flexible and dynamic sub-path generation during the
	// traversal process.
	SubPathHook func(info *SubPathInfo) string

	// ChainSubPathHook chainable version of SubPathHook
	ChainSubPathHook func(result string,
		info *SubPathInfo,
	) string
)
