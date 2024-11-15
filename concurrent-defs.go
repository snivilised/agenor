package age

import (
	"github.com/snivilised/agenor/core"
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
		Handler core.Client
	}

	// TraverseJobStream represents the core channel type of the worker pool's
	// input stream. The client owns this channel and is responsible for
	// closing it when done or invoking Conclude directly on the pool (See
	// boost for more details).
	TraverseJobStream = pants.JobStream[TraverseInput]

	// TraverseJobStreamR worker pool's read stream, pool reads from this channel.
	TraverseJobStreamR = pants.JobStreamR[TraverseInput]

	// TraverseJobStreamW worker pool's write stream, client writes to this channel.
	TraverseJobStreamW = pants.JobStreamW[TraverseInput]

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

	// TraverseOutputStream represents the core channel type of the worker pool's
	// output stream. The pool owns this stream and will be closed only when
	// safe to do so, which will be anytime after navigation is complete.
	// The channel is only closed when there are no remaining outstanding jobs
	// and all workers are idle.
	TraverseOutputStream = pants.JobOutputStream[TraverseOutput]

	// TraverseOutputStreamR worker pool's output stream read by the client.
	TraverseOutputStreamR = pants.JobOutputStreamR[TraverseOutput]

	// TraverseOutputStreamW worker pool's output stream written to by the pool.
	TraverseOutputStreamW = pants.JobOutputStreamW[TraverseOutput]
)
