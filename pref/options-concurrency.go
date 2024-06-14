package pref

import (
	"runtime"
)

type ConcurrencyOptions struct {
	NoW uint
}

func WithCPU() Option {
	return func(o *Options) error {
		o.Core.Concurrency.NoW = uint(runtime.NumCPU())

		return nil
	}
}

func WithNoW(now uint) Option {
	return func(o *Options) error {
		o.Core.Concurrency.NoW = now

		return nil
	}
}
