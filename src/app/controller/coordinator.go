package controller

import (
	"context"
	"fmt"
	"os/exec"
	"regexp"

	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/app/bedrock"
	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/app/shell"
)

// Coordinator coordinates the layers between the command adapters and
// the agenor traversal engine. It is the single place that constructs
// pref.Facade values and calls into agenor via the agenor.Scenario on
// the request. It never imports cobra, mamba, or the command package.
//
// Dependency direction: command -> controller -> agenor
type Coordinator struct {
	cfg           *bedrock.Config
	locate        shell.LocateFunc
	forestBuilder pref.BuildForest
	actionRegexes map[string]*regexp.Regexp
}

// CoordinatorOption is a functional option for Coordinator.
type CoordinatorOption func(*Coordinator)

// WithLocate overrides the LocateFunc used during PreFlight to validate
// whether action executables are invokable. The default is the
// platform-appropriate function returned by shell.Detect(). Use this
// in tests to inject a stub without spawning real subprocesses.
func WithLocate(fn shell.LocateFunc) CoordinatorOption {
	return func(c *Coordinator) {
		c.locate = fn
	}
}

// WithForest allows injection of a pref.BuildForest, which is used to construct
// the file systems during traversal. This is primarily intended for testing,
// where a stubbed file system can be used to simulate various scenarios without
// relying on the real file system. In production, the Coordinator will use the
// default file system builder if this option is not provided.
func WithForest(forestBuilder pref.BuildForest) CoordinatorOption {
	return func(c *Coordinator) {
		c.forestBuilder = forestBuilder
	}
}

// New returns a ready-to-use Coordinator. cfg must not be nil.
func New(cfg *bedrock.Config, opts ...CoordinatorOption) *Coordinator {
	actionRegexes := make(map[string]*regexp.Regexp)
	if cfg != nil && cfg.Raw.Actions != nil {
		for name, action := range cfg.Raw.Actions {
			if action.Capture != "" {
				if re, err := regexp.Compile(action.Capture); err == nil {
					actionRegexes[name] = re
				}
			}
		}
	}

	c := &Coordinator{
		cfg:           cfg,
		locate:        func(name string) (string, error) { return exec.LookPath(name) },
		actionRegexes: actionRegexes,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// ExecutePrime runs a fresh directory traversal using the scenario
// provided on the request. When the presenter implements PeerAware
// and NeedsPeerInfo returns true, a preview traversal is run first
// to build the PeerInfoMap and collect node counts for the progress
// indicator. The live traversal reuses the options built during the
// preview pass.
func (c *Coordinator) ExecutePrime(ctx context.Context, req *PrimeRequest) error {
	req.Root = req.Tree

	traversal := &report.Traversal{}
	view, isPeerAware := req.UI.(report.PeerAware)

	if isPeerAware && view.NeedsPeerInfo() {
		// Execute the preview traversal to build the PeerInfoMap and collect.
		fmt.Println("🦄 DEBUG: Coordinator.ExecutePrime: peer aware ... 🦄")
		peerInfoMap, builtOptions, result, err := buildPeerInfoMap(
			ctx, req, req.Settings,
		)
		if err != nil {
			return err
		}

		view.OnPeerInfoBegin(
			result.Metrics().Count(enums.MetricNoFilesInvoked),
			result.Metrics().Count(enums.MetricNoDirectoriesInvoked),
		)

		facade := &pref.Using{
			Subscription: req.Subscription,
			Head: pref.Head{
				Handler: func(servant agenor.Servant) error {
					return c.handleServant(servant, &req.Request, traversal, peerInfoMap)
				},
				GetForest: c.forestBuilder,
			},
			Tree: req.Tree,
			O:    builtOptions,
		}

		// Execute the live traversal with the peer info map and options from the
		// preview pass.
		err = c.execute(ctx, &req.Request, facade, traversal, true, "")
		view.OnPeerInfoEnd()

		return err
	}

	fmt.Println("🦁 DEBUG: Coordinator.ExecutePrime: executing live traversal only ... 🦁")
	facade := &pref.Using{
		Subscription: req.Subscription,
		Head: pref.Head{
			Handler: func(servant agenor.Servant) error {
				return c.handleServant(servant, &req.Request, traversal, nil)
			},
			GetForest: c.forestBuilder,
		},
		Tree: req.Tree,
	}

	// Execute the live traversal without peer info.
	return c.execute(ctx, &req.Request, facade, traversal, true, "")
}

// ExecuteResume resumes an interrupted traversal. Peer info is not
// currently supported for resume - this will be addressed in a
// dedicated issue.
func (c *Coordinator) ExecuteResume(ctx context.Context, req *ResumeRequest) error {
	traversal := &report.Traversal{}

	// TODO: implement peer info support for resume traversals.

	facade := &pref.Relic{
		Head: pref.Head{
			Handler: func(servant agenor.Servant) error {
				return c.handleServant(servant, &req.Request, traversal, nil)
			},
		},
		Strategy: req.Strategy,
	}

	// Execute the live traversal without peer info.
	return c.execute(ctx, &req.Request, facade, traversal, false, req.ResumeFrom)
}

// execute is the shared orchestration path for both prime and resume
// traversals.
func (c *Coordinator) execute(
	ctx context.Context,
	req *Request,
	facade pref.Facade,
	traversal *report.Traversal,
	isPrime bool,
	resumeFrom string,
) error {
	if err := c.PreFlight(req); err != nil {
		return err
	}

	req.UI.OnBegin(&report.BeginEvent{
		Root:       req.Root,
		Caption:    c.captionFor(req),
		StartedAt:  core.Now(),
		IsPrime:    isPrime,
		ResumeFrom: resumeFrom,
	})

	result, err := req.Scenario(facade, req.Settings...).Navigate(ctx)

	traversal.Err = err
	if result != nil {
		traversal.FilesVisited = result.Metrics().Count(enums.MetricNoFilesInvoked)
		traversal.DirsVisited = result.Metrics().Count(enums.MetricNoDirectoriesInvoked)
		traversal.Elapsed = result.Session().Elapsed()
	}

	req.UI.OnComplete(traversal)

	return err
}

//nolint:exhaustive // enums.SubscribeDirectoriesWithFiles, enums.SubscribeUniversal
func (c *Coordinator) captionFor(req *Request) string {
	subscription := ""
	switch req.Subscription {
	case enums.SubscribeFiles:
		subscription = "files only"
	case enums.SubscribeDirectories:
		subscription = "folders only"
	default:
		subscription = "files and folders"
	}

	if req.ActionName != "" {
		action, ok := c.cfg.Raw.Actions[req.ActionName]
		if ok {
			return fmt.Sprintf("%s via '%s'", subscription, action.Cmd)
		}
	}

	if req.PipelineName != "" {
		return fmt.Sprintf("%s via pipeline '%s'", subscription, req.PipelineName)
	}

	return subscription
}
