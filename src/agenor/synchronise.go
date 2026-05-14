package agenor

import (
	"context"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/pants"
)

type synchroniser interface {
	enclave.KernelNavigator
	Ignite(*enclave.Ignition)
	IsComplete() bool
	Bye(result *enclave.KernelResult)
}

type trunk struct {
	kc  enclave.KernelController
	o   *pref.Options
	ext extent
	err error
}

func (t *trunk) IsComplete() bool {
	return t.ext.complete()
}

func (t *trunk) Ignite(ignition *enclave.Ignition) {
	t.kc.Ignite(ignition)
}

func (t *trunk) Bye(result *enclave.KernelResult) {
	t.kc.Bye(result)
}

type concurrent struct {
	trunk
	wg pants.WaitGroup
}

func (c *concurrent) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	if c.err != nil {
		return c.kc.Result(ctx), c.err
	}

	return c.kc.Navigate(ctx)
}

type sequential struct {
	trunk
}

func (s *sequential) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	if s.err != nil {
		return s.kc.Result(ctx), s.err
	}

	return s.kc.Navigate(ctx)
}
