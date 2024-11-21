package age

import (
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/pants"
)

type (
	// TraverseJobStream represents the core channel type of the worker pool's
	// input stream. The client owns this channel and is responsible for
	// closing it when done or invoking Conclude directly on the pool (See
	// boost for more details).
	TraverseJobStream = pants.JobStream[core.TraverseInput]

	// TraverseJobStreamR worker pool's read stream, pool reads from this channel.
	TraverseJobStreamR = pants.JobStreamR[core.TraverseInput]

	// TraverseJobStreamW worker pool's write stream, client writes to this channel.
	TraverseJobStreamW = pants.JobStreamW[core.TraverseInput]

	// TraverseOutputStream represents the core channel type of the worker pool's
	// output stream. The pool owns this stream and will be closed only when
	// safe to do so, which will be anytime after navigation is complete.
	// The channel is only closed when there are no remaining outstanding jobs
	// and all workers are idle.
	TraverseOutputStream = pants.JobOutputStream[core.TraverseOutput]

	// TraverseOutputStreamR worker pool's output stream read by the client.
	TraverseOutputStreamR = pants.JobOutputStreamR[core.TraverseOutput]

	// TraverseOutputStreamW worker pool's output stream written to by the pool.
	TraverseOutputStreamW = pants.JobOutputStreamW[core.TraverseOutput]
)
