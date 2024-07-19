package nfs

import (
	"io/fs"
)

// 📚 package: nfs contains file system abstractions for navigation. Since
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
)
