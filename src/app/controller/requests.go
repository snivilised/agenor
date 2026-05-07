package controller

import (
	"github.com/snivilised/jaywalk/src/agenor"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/agenor/pref"
	"github.com/snivilised/jaywalk/src/app/report"
)

// Request holds the fields common to all traversal requests.
// It is embedded as the first field in PrimeRequest and ResumeRequest.
type Request struct {
	// Subscription controls which node types are visited.
	Subscription enums.Subscription

	// Settings are the agenor option functions derived from shared flags.
	Settings []pref.Option

	// ActionName is the name of the configured action to execute per node.
	// Empty when no action is configured.
	ActionName string

	// PipelineName is the name of the configured pipeline to execute per node.
	// Empty when no pipeline is configured.
	PipelineName string

	// Scenario is the agenor scenario provided by the command adapter.
	// It encapsulates the walk/run distinction so the coordinator is
	// unaware of it. Set to agenor.Tortoise(isPrime) for walk or
	// agenor.Hare(isPrime, wg) for run.
	Scenario agenor.Scenario

	// Root is the traversal root directory. For prime traversals this is
	// sourced from the --tree argument. For resume traversals it will be
	// sourced from the restored checkpoint state. The expand function uses
	// this to detect placeholder breaches.
	Root string

	// UI is the Presenter injected by Bootstrap via PersistentPreRunE.
	// It receives traversal events and decides how to render them.
	// The controller never formats output directly.
	UI report.Presenter

	// GetForest is the pref.BuildForest used to construct the file system
	GetForest pref.BuildForest
}

// PrimeRequest carries everything the coordinator needs to execute a
// fresh directory traversal.
type PrimeRequest struct {
	Request

	// Tree is the root directory path to traverse.
	Tree string
}

// ResumeRequest carries everything the coordinator needs to resume an
// interrupted traversal. Strategy is agenor's internal resume strategy,
// distinct from the app-layer Scenario field on Request.
type ResumeRequest struct {
	Request

	// Strategy controls how the resume proceeds within agenor.
	Strategy enums.ResumeStrategy

	// ResumeFrom is the path from which the traversal continues. Sourced
	// from the restored checkpoint state. Passed to OnBegin so the
	// renderer can display the resume point in the opening banner.
	ResumeFrom string
}
