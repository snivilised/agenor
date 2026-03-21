package controller

import (
	"context"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/app/report"
)

// Coordinator coordinates the layers between the command adapters and
// the agenor traversal engine. It is the single place that constructs
// pref.Facade values and calls into agenor via the agenor.Scenario on
// the request. It never imports cobra, mamba, or the command package.
//
// Dependency direction: command → controller → agenor
type Coordinator struct{}

// New returns a ready-to-use Coordinator.
func New() *Coordinator {
	return &Coordinator{}
}

// ExecutePrime runs a fresh directory traversal using the scenario
// provided on the request. The command adapter is responsible for
// constructing the correct agenor.Scenario (Tortoise or Hare).
func (c *Coordinator) ExecutePrime(ctx context.Context, req *PrimeRequest) error {
	facade := &pref.Using{
		Subscription: req.Subscription,
		Head: pref.Head{
			Handler: func(servant agenor.Servant) error {
				return c.handleNode(servant.Node(), &req.Request)
			},
		},
		Tree: req.Tree,
	}

	return c.execute(ctx, &req.Request, facade)
}

// ExecuteResume resumes an interrupted traversal using the scenario
// provided on the request. The command adapter is responsible for
// constructing the correct agenor.Scenario (Tortoise or Hare).
func (c *Coordinator) ExecuteResume(ctx context.Context, req *ResumeRequest) error {
	facade := &pref.Relic{
		Head: pref.Head{
			Handler: func(servant agenor.Servant) error {
				return c.handleNode(servant.Node(), &req.Request)
			},
		},
		Strategy: req.Strategy,
	}

	return c.execute(ctx, &req.Request, facade)
}

// execute is the shared orchestration path for both prime and resume
// traversals. It calls the scenario, collects metrics, and notifies
// the UI via OnComplete.
func (c *Coordinator) execute(ctx context.Context, req *Request, facade pref.Facade) error {
	result, err := req.Scenario(facade, req.Options...).Navigate(ctx)

	t := &report.Traversal{Err: err}
	if result != nil {
		t.FilesVisited = result.Metrics().Count(enums.MetricNoFilesInvoked)
		t.DirsVisited = result.Metrics().Count(enums.MetricNoDirectoriesInvoked)
		t.Elapsed = result.Session().Elapsed()
	}

	req.UI.OnComplete(t)

	return err
}
