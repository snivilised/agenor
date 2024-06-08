package tv

import (
	"github.com/snivilised/lorax/boost"
	"github.com/snivilised/traverse/core"
)

type (
	// TraverseInput represents the type  of inputs accepted by the worker pool
	TraverseInput struct {
		// Node represents the file system entity (file or folder) for which
		// a job will execute.
		Node *core.Node

		// Handler is the client defined callback function that should be
		// invoked for all eligible Nodes.
		Handler core.Client
	}

	// TraverseJobStream represents the core channel type of the worker pool's
	// input stream. The client owns this channel and is responsible for
	// closing it when done or invoking Conclude directly on the pool (See
	// boost for more details).
	TraverseJobStream = boost.JobStream[TraverseInput]

	// TraverseJobStreamR worker pool's read stream, pool reads from this channel.
	TraverseJobStreamR = boost.JobStreamR[TraverseInput]

	// TraverseJobStreamW worker pool's write stream, client writes to this channel.
	TraverseJobStreamW = boost.JobStreamW[TraverseInput]

	// TraverseOutput represents the output of a single job executed by the pool.
	TraverseOutput struct {
		// Node represents the file system entity (file or folder) from
		// which this output was generated via the client defined handler.
		Node *core.Node

		// Error error result of client's handler.
		Error error

		// Data is a custom field reserved for the client
		Data any
	}

	// TraverseOutputStream represents the core channel type of the worker pool's
	// output stream. The pool owns this stream and will be closed only when
	// safe to do so, which will be anytime after navigation is complete.
	// The channel is only closed when there are no remaining outstanding jobs
	// and all workers are idle.
	TraverseOutputStream = boost.JobOutputStream[TraverseOutput]

	// TraverseOutputStreamR worker pool's output stream read by the client.
	TraverseOutputStreamR = boost.JobOutputStreamR[TraverseOutput]

	// TraverseOutputStreamW worker pool's output stream written to by the pool.
	TraverseOutputStreamW = boost.JobOutputStreamW[TraverseOutput]
)
