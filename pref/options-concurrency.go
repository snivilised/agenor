package pref

import (
	"runtime"
	"time"

	"github.com/snivilised/agenor/core"
)

// ConcurrencyOptions specifies options used for current traversal sessions
type (
	InputOptions struct {
		// Size specifies the size of the input channel, if not specified, this
		// will default to the number of workers requested via NoW(WithNow). The
		// main Go routine under which the navigator runs will be able to submit
		// jobs to the pool via this input channel. Size determines how many jobs can
		// can be inside the input channel, before this main Go routine will start to
		// block.
		Size uint
	}

	OutputOptions struct {
		// Size specifies the size of the output channel. If set to 0, the
		// size will revert to the number of workers requested via NoW(WithNow).
		// As each worker completes, they will send their output to the output
		// channel. They will start to block once the output channel reaches
		// capacity. The output channel must be consumed by the client in order
		// for the channel to be depleted, thus allowing workers to continue to
		// send their outputs.
		Size uint

		// CheckCloseInterval is used to assist in the process of the pool closing
		// the output channel. The navigator will close the pool on a polling
		// basis and CheckCloseInterval is the time in between successive
		// close attempts.
		CheckCloseInterval time.Duration

		// TimeoutOnSend is the duration of time that must elapse after each
		// worker attempts to send their output to the output channel, before
		// being abandoned resulting in eventual pool closure in error. This
		// is used to prevent the possibility of deadlock
		TimeoutOnSend time.Duration

		// On is a callback which is invoked if the an output channel
		// is requested. This callback will be provided with the output channel
		// required so that the output of the pool can be consumed.
		On core.OutputFunc
	}

	ConcurrencyOptions struct {
		// NoW specifies the number of go-routines to use in the worker
		// pool used for concurrent traversal sessions requested by using
		// the Run function.
		NoW uint

		// Input contains input channel properties
		Input InputOptions

		// Output contains output properties
		Output OutputOptions
	}
)

// WithCPU configures the worker pool used for concurrent traversal sessions
// in the Run function to utilise a number of go-routines equal to the available
// CPU count, optimising performance based on the system's processing capabilities.
func WithCPU() Option {
	return func(o *Options) error {
		o.Concurrency.NoW = uint(runtime.NumCPU())

		return nil
	}
}

// WithNoW sets the number of go-routines to use in the worker
// pool used for concurrent traversal sessions requested by using
// the Run function.
func WithNoW(now uint) Option {
	return func(o *Options) error {
		o.Concurrency.NoW = now

		return nil
	}
}

// WithOutput requests that the worker pool emits outputs
func WithOutput(output *OutputOptions) Option {
	return func(o *Options) error {
		o.Concurrency.Output = *output

		return nil
	}
}
