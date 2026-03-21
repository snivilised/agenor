// Package tfs provides a local file system implementation for traversal.
// It includes functions to create both relative and absolute file systems
// using the nefilim library. The NewFS function creates a relative local
// file system based on a provided relative path, while the New function
// creates an absolute local file system. These functions are essential for
// setting up the file system context required for traversal operations in
// the agenor library. Should not depend on anything else in agenor.
package tfs
