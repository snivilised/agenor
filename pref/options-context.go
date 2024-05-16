package pref

import (
	"context"
	"runtime"
)

type AccelerationOptions struct {
	ctx    context.Context
	cancel context.CancelFunc
	now    int
}

func (ao *AccelerationOptions) Cancellation() (context.Context, context.CancelFunc) {
	return ao.ctx, ao.cancel
}

func WithContext(ctx context.Context) Option {
	return func(o *Options) error {
		o.Acceleration.ctx = ctx

		return nil
	}
}

func WithCancel(cancel context.CancelFunc) Option {
	return func(o *Options) error {
		o.Acceleration.cancel = cancel

		return nil
	}
}

func WithCPU() Option {
	return func(o *Options) error {
		o.Acceleration.now = runtime.NumCPU()

		return nil
	}
}

func WithNoW(now int) Option {
	return func(o *Options) error {
		o.Acceleration.now = now

		return nil
	}
}
