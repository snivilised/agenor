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
	Bye(result *enclave.KernelResult)
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

func (t *trunk) Bye(result *enclave.KernelResult) {
	t.kc.Bye(result)
}

type concurrent struct {
	trunk
	wg      pants.WaitGroup
	pool    *core.TraversePool
	inputCh pants.SourceStreamW[*core.TraverseInput]
	swapper enclave.Swapper
}

func (c *concurrent) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	defer func() {
		c.close()
		if c.pool != nil {
			c.pool.Release(ctx)
		}
	}()

	if c.err != nil {
		return c.kc.Result(ctx), c.err
	}

	if c.swapper != nil {
		c.swapper.Swap(func(servant Servant) error {
			handler := c.ext.facade().Client()
			input := &core.TraverseInput{
				Servant: servant,
				Handler: handler,
			}

			c.inputCh <- input // TODO: support for timeout (TimeoutOnSendInput) ??? issue #333

			return nil
		})
	}

	c.pool, c.err = pants.NewManifoldFuncPool(
		ctx, func(input *core.TraverseInput) (*core.TraverseOutput, error) {
			err := input.Handler(input.Servant)

			return &core.TraverseOutput{
				Servant: input.Servant,
				Error:   err,
			}, err
		}, c.wg,
		pants.WithSize(c.o.Concurrency.NoW),
		pants.WithInput(c.o.Concurrency.Input.Size),
		pants.IfOptionF(c.o.Concurrency.Output.On != nil, func() pants.Option {
			return pants.WithOutput(
				c.o.Concurrency.Output.Size,
				c.o.Concurrency.Output.CheckCloseInterval,
				c.o.Concurrency.Output.TimeoutOnSend,
			)
		}),
	)

	if c.err != nil {
		err := errors.Wrap(c.err, locale.ErrWorkerPoolCreationFailed.Error())
		return c.kc.Result(ctx), err
	}
	c.open(ctx)

	if c.o.Concurrency.Output.On != nil {
		c.o.Concurrency.Output.On(c.pool.Observe())
	}

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
