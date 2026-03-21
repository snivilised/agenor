package kernel

import (
	"context"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/internal/enclave"
	"github.com/snivilised/jaywalk/src/agenor/pref"
)

// HadesNav creates a new navigator hades which is used when the
// navigator cannot be built.
func HadesNav(o *pref.Options, err error) enclave.KernelController {
	return &navigatorHades{
		o:   o,
		err: err,
	}
}

type navigatorHades struct {
	o   *pref.Options
	err error
}

// Rank is a no-op.
func (n *navigatorHades) Rank() {
}

// Ignite is a no-op.
func (n *navigatorHades) Ignite(*enclave.Ignition) {
}

// Navigate returns the error result of the navigation.
func (n *navigatorHades) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	return n.Result(ctx), n.err
}

// Result returns the error result of the navigation.
func (n *navigatorHades) Result(_ context.Context) *enclave.KernelResult {
	return enclave.NewFailed()
}

// Snooze is a no-op.
func (n *navigatorHades) Snooze(context.Context,
	*core.ActiveState,
) (*enclave.KernelResult, error) {
	return &enclave.KernelResult{}, nil
}

// Bye is a no-op.
func (n *navigatorHades) Bye(_ core.TraverseResult) {
}
