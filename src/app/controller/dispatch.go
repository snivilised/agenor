package controller

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/app/report"
)

// handleNode dispatches to the appropriate per-node handler based on
// whether an action, pipeline, or neither is configured on the request.
func (c *Coordinator) handleNode(node *core.Node, req *Request) error {
	switch {
	case req.PipelineName != "":
		e := c.executePipeline(node, req.PipelineName)
		req.UI.OnPipelineEvent(e)
		return e.Err

	case req.ActionName != "":
		e := c.executeAction(node, req.ActionName)
		req.UI.OnActionEvent(e)
		return e.Err

	default:
		req.UI.OnNodeEvent(&report.NeutralEvent{
			DisplayEvent: report.DisplayEvent{Node: node},
		})
		return nil
	}
}

// executeAction composes and executes a configured action against a node.
// Composition of the execution string is a future concern — this is a stub.
func (c *Coordinator) executeAction(node *core.Node, name string) *report.ActionEvent {
	return &report.ActionEvent{
		DisplayEvent: report.DisplayEvent{
			Node: node,
			Name: name,
		},
	}
}

// executePipeline composes and executes a configured pipeline against a node.
// Composition of the execution string is a future concern — this is a stub.
func (c *Coordinator) executePipeline(node *core.Node, name string) *report.PipelineEvent {
	return &report.PipelineEvent{
		DisplayEvent: report.DisplayEvent{
			Node: node,
			Name: name,
		},
	}
}
