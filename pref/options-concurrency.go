package pref

import (
	"runtime"
)

// ConcurrencyOptions specifies options used for current traversal sessions
type ConcurrencyOptions struct {
	// NoW specifies the number of go-routines to use in the worker
	// pool used for concurrent traversal sessions requested by using
	// the Run function.
	NoW uint
}

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
