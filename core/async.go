package core

import (
	"context"

	"github.com/snivilised/pants"
)

type (
	// TraverseInput represents the type  of inputs accepted by the worker pool
	TraverseInput struct {
		// Servant represents the file system entity (file or directory) for which
		// a job will execute.
		Servant Servant

		// Handler is the client defined callback function that should be
		// invoked for all eligible Nodes.
		Handler Client
	}

	// TraverseOutput represents the output of a single job executed by the pool.
	TraverseOutput struct {
		// Servant represents the file system entity (file or directory) from
		// which this output was generated via the client defined handler.
		Servant Servant

		// Error error result of client's handler.
		Error error

		// Data is a custom field reserved for the client
		Data any
	}

	TraversePool = pants.ManifoldFuncPool[*TraverseInput, *TraverseOutput]

	// OutputStream the traverse output channel
	OutputStream = pants.JobOutputStreamR[*TraverseOutput]

	// Cancellation information required for the cancellation monitor
	Cancellation struct {
		Cancel context.CancelFunc
		On     pants.OnCancel
	}

	// OutputFunc
	OutputFunc func(outs OutputStream)
)
