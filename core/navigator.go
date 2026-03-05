package core

import (
	"context"
)

// Navigator represents the interface for navigating through a file
// system. It defines a single method,
// Navigate, which takes a context and returns a TraverseResult and an
// error. The Navigate method is responsible for performing the traversal
// and returning the results of the navigation process.
type Navigator interface {
	// Navigate performs the traversal through the file system,
	// using the provided context for cancellation and timeout control. It returns a
	// TraverseResult containing the results of the navigation process, as well as an
	// error if any issues were encountered during the traversal.
	Navigate(ctx context.Context) (TraverseResult, error)
}
