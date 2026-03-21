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

	// TraversePool is the type of the worker pool that will execute the client's
	// handler for each eligible Node. It is defined as a ManifoldFuncPool that
	// takes TraverseInput and produces TraverseOutput.
	TraversePool = pants.ManifoldFuncPool[*TraverseInput, *TraverseOutput]

	// OutputStream the traverse output channel
	OutputStream = pants.JobOutputStreamR[*TraverseOutput]

	// Cancellation information required for the cancellation monitor
	Cancellation struct {
		// Cancel is the function that can be called to cancel the ongoing traversal.
		Cancel context.CancelFunc

		// On is a callback function that can be used to register a handler that will
		// be invoked when the cancellation is triggered. This allows the client to
		// perform any necessary cleanup or finalization when the traversal is cancelled.
		On pants.OnCancel
	}

	// OutputFunc is the type of the client defined callback function that will be invoked
	// for each output produced by the worker pool. It takes an OutputStream as an argument,
	// which allows the client to receive and process the outputs of the traversal.
	OutputFunc func(outs OutputStream)
)
