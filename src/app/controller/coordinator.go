package controller

import (
	"context"
	"os/exec"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/app/bedrock"
	"github.com/snivilised/jaywalk/src/app/report"
)

// Coordinator coordinates the layers between the command adapters and
// the agenor traversal engine. It is the single place that constructs
// pref.Facade values and calls into agenor via the agenor.Scenario on
// the request. It never imports cobra, mamba, or the command package.
//
// Dependency direction: command -> controller -> agenor
type Coordinator struct {
	cfg      *bedrock.Config
	lookPath func(string) (string, error)
}

// CoordinatorOption is a functional option for Coordinator.
type CoordinatorOption func(*Coordinator)

// WithLookPath overrides the function used to locate executables on PATH.
// The default is exec.LookPath. Use this in tests to inject a stub.
func WithLookPath(fn func(string) (string, error)) CoordinatorOption {
	return func(c *Coordinator) {
		c.lookPath = fn
	}
}

// New returns a ready-to-use Coordinator. cfg must not be nil.
// Optional CoordinatorOptions may be supplied to override defaults.
func New(cfg *bedrock.Config, opts ...CoordinatorOption) *Coordinator {
	c := &Coordinator{
		cfg:      cfg,
		lookPath: exec.LookPath,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// ExecutePrime runs a fresh directory traversal using the scenario
// provided on the request. The command adapter is responsible for
// constructing the correct agenor.Scenario (Tortoise or Hare).
func (c *Coordinator) ExecutePrime(ctx context.Context, req *PrimeRequest) error {
	// Root is sourced from Tree for prime traversals. Resume will source
	// it from restored checkpoint.
	req.Root = req.Tree

	t := &report.Traversal{}

	facade := &pref.Using{
		Subscription: req.Subscription,
		Head: pref.Head{
			Handler: func(servant agenor.Servant) error {
				return c.handleNode(servant.Node(), &req.Request, t)
			},
		},
		Tree: req.Tree,
	}

	return c.execute(ctx, &req.Request, facade, t)
}

// ExecuteResume resumes an interrupted traversal using the scenario
// provided on the request. The command adapter is responsible for
// constructing the correct agenor.Scenario (Tortoise or Hare).
// Note: req.Request.Root must be populated from restored checkpoint
// state before this method is called - that is a resume concern.
func (c *Coordinator) ExecuteResume(ctx context.Context, req *ResumeRequest) error {
	t := &report.Traversal{}

	facade := &pref.Relic{
		Head: pref.Head{
			Handler: func(servant agenor.Servant) error {
				return c.handleNode(servant.Node(), &req.Request, t)
			},
		},
		Strategy: req.Strategy,
	}

	return c.execute(ctx, &req.Request, facade, t)
}

// execute is the shared orchestration path for both prime and resume
// traversals. PreFlight is always the first step - a failure here
// returns immediately before any traversal begins. On success it calls
// the scenario, collects metrics, and notifies the UI via OnComplete.
func (c *Coordinator) execute(
	ctx context.Context,
	req *Request,
	facade pref.Facade,
	t *report.Traversal,
) error {
	if err := c.PreFlight(req); err != nil {
		return err
	}

	result, err := req.Scenario(facade, req.Options...).Navigate(ctx)

	t.Err = err
	if result != nil {
		t.FilesVisited = result.Metrics().Count(enums.MetricNoFilesInvoked)
		t.DirsVisited = result.Metrics().Count(enums.MetricNoDirectoriesInvoked)
		t.Elapsed = result.Session().Elapsed()
	}

	req.UI.OnComplete(t)

	return err
}
