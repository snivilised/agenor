package age

import (
	"context"

	"github.com/pkg/errors"
	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/pants"
)

type synchroniser interface {
	enclave.KernelNavigator
	Ignite(*enclave.Ignition)
	IsComplete() bool
	Conclude(result *enclave.KernelResult)
}

type trunk struct {
	kc  enclave.KernelController
	o   *pref.Options
	ext extent
	err error
}

func (t *trunk) extent() extent {
	return t.ext
}

func (t *trunk) IsComplete() bool {
	return t.ext.complete()
}

func (t *trunk) Ignite(ignition *enclave.Ignition) {
	t.kc.Ignite(ignition)
}

func (t *trunk) Conclude(result *enclave.KernelResult) {
	t.kc.Conclude(result)
}

type concurrent struct {
	trunk
	wg        pants.WaitGroup
	pool      *pants.ManifoldFuncPool[*TraverseInput, *TraverseOutput]
	decorator core.Client
	inputCh   pants.SourceStreamW[*TraverseInput]
}

func (c *concurrent) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	defer c.close()

	if c.err != nil {
		return c.kc.Result(ctx), c.err
	}

	c.decorator = func(servant Servant) error {
		// c.decorator is the function we register with the navigator,
		// so instead of invoking the client handler, the navigator
		// will invoke the decorator, which will send a job to the pool
		// containing the current file system node. The navigator is
		// not aware that its invoking the decorator ...
		//
		// TODO: later, we need to be able to decorate the client callback,
		// either by a Tap or a bus event...
		//
		input := &TraverseInput{
			Servant: servant,
			Handler: c.ext.facade().Client(),
		}

		c.inputCh <- input // support for timeout (TimeoutOnSendInput) ???

		return nil
	}

	c.pool, c.err = pants.NewManifoldFuncPool(
		ctx, func(input *TraverseInput) (*TraverseOutput, error) {
			err := input.Handler(input.Servant)

			return &TraverseOutput{
				Servant: input.Servant,
				Error:   err,
			}, err
		}, c.wg,
		pants.WithSize(c.o.Concurrency.NoW),
		pants.WithOutput(OutputChSize, CheckCloseInterval, TimeoutOnSend),
	)

	if c.err != nil {
		err := errors.Wrap(c.err, locale.ErrWorkerPoolCreationFailed.Error())
		return c.kc.Result(ctx), err
	}
	c.open(ctx)

	return c.kc.Navigate(ctx)
}

func (c *concurrent) open(ctx context.Context) {
	c.inputCh = c.pool.Source(ctx, c.wg)
}

func (c *concurrent) close() {
	if c.inputCh != nil {
		close(c.inputCh)
	}
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
