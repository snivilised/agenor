package nfs

import (
	"io/fs"
	"os"
)

// 📦 pkg: nfs - contains file system abstractions for navigation. Since
// there are no standard write-able file system interfaces,
// we need to define proprietary ones here in this package.
// This is a low level package that should not use anything else in
// traverse.

type (
	// FileSystems contains the logical file systems required
	// for navigation.
	FileSystems struct {
		// N represents the read only navigation file system. Uses
		// of the shelf interface as defined by the standard library.
		N fs.ReadDirFS

		// Q represents the file system instance that can perform a query
		// status via an Lstat call.
		Q fs.StatFS

		// R represents the resume/save file system that requires
		// write access and whose path should be outside of the path
		// represented by N, the navigation file system.
		R fs.FS
	}

	// ExistsInFS contains methods that check the existence of file system items.
	ExistsInFS interface {
		// FileExists does file exist at the path specified
		FileExists(path string) bool

		// DirectoryExists does directory exist at the path specified
		DirectoryExists(path string) bool
	}

	// MkDirAllFS is a file system with a MkDirAll method.
	MkDirAllFS interface {
		ExistsInFS
		MkDirAll(path string, perm os.FileMode) error
	}
)
