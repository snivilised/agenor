package controller

import (
	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/locale"
)

// handleNode dispatches to the appropriate per-node handler based on
// whether an action, pipeline, or neither is configured on the request.
func (c *Coordinator) handleNode(node *core.Node, req *Request, t *report.Traversal) error {
	switch {
	case req.PipelineName != "":
		e := c.executePipeline(node, req.PipelineName, req.Root)
		if e.Skipped {
			t.ActionsSkipped++
			req.UI.OnSkipEvent(&report.SkipEvent{
				DisplayEvent: report.DisplayEvent{Node: node, Name: req.PipelineName},
				Placeholder:  e.Placeholder,
				ResolvedPath: e.ResolvedPath,
			})
			return nil
		}
		req.UI.OnPipelineEvent(e.Event)
		return e.Event.Err

	case req.ActionName != "":
		e := c.executeAction(node, req.ActionName, req.Root)
		if e.Skipped {
			t.ActionsSkipped++
			req.UI.OnSkipEvent(&report.SkipEvent{
				DisplayEvent: report.DisplayEvent{Node: node, Name: req.ActionName},
				Placeholder:  e.Placeholder,
				ResolvedPath: e.ResolvedPath,
			})
			return nil
		}
		req.UI.OnActionEvent(e.Event)
		return e.Event.Err

	default:
		req.UI.OnNodeEvent(&report.NeutralEvent{
			DisplayEvent: report.DisplayEvent{Node: node},
		})
		return nil
	}
}

// actionResult is the outcome of executeAction. Either Skipped is true
// (and Placeholder/ResolvedPath carry the breach details), or Event
// carries the ActionEvent to hand to the UI.
type actionResult struct {
	Skipped      bool
	Placeholder  string
	ResolvedPath string
	Event        *report.ActionEvent
}

// pipelineResult mirrors actionResult for pipeline execution.
type pipelineResult struct {
	Skipped      bool
	Placeholder  string
	ResolvedPath string
	Event        *report.PipelineEvent
}

// executeAction expands the cmd string for the named action and returns
// the result. If a placeholder breaches root the result is marked as
// skipped and no shell execution is attempted.
func (c *Coordinator) executeAction(node *core.Node, name, root string) actionResult {
	action, ok := c.cfg.Raw.Actions[name]
	if !ok {
		// PreFlight should have caught this - treat as an action error.
		return actionResult{
			Event: &report.ActionEvent{
				DisplayEvent: report.DisplayEvent{Node: node, Name: name},
				Err:          locale.NewActionNotFoundError(name),
			},
		}
	}

	result := Expand(action.Cmd, root, node)
	if result.Skipped {
		return actionResult{
			Skipped:      true,
			Placeholder:  result.Placeholder,
			ResolvedPath: result.ResolvedPath,
		}
	}

	return actionResult{
		Event: &report.ActionEvent{
			DisplayEvent:    report.DisplayEvent{Node: node, Name: name},
			ExecutionString: result.Cmd,
		},
	}
}

// executePipeline expands and executes each step in the named pipeline
// in order. The first skip or error stops the pipeline for this node.
func (c *Coordinator) executePipeline(node *core.Node, name, root string) pipelineResult {
	pipeline, ok := c.cfg.Raw.Pipelines[name]
	if !ok {
		return pipelineResult{
			Event: &report.PipelineEvent{
				DisplayEvent: report.DisplayEvent{Node: node, Name: name},
				Err:          locale.NewPipelineNotFoundError(name),
			},
		}
	}

	for _, step := range pipeline.Steps {
		ar := c.executeAction(node, step, root)
		if ar.Skipped {
			return pipelineResult{
				Skipped:      true,
				Placeholder:  ar.Placeholder,
				ResolvedPath: ar.ResolvedPath,
			}
		}
		if ar.Event.Err != nil {
			return pipelineResult{
				Event: &report.PipelineEvent{
					DisplayEvent:    report.DisplayEvent{Node: node, Name: name},
					ExecutionString: ar.Event.ExecutionString,
					Err:             ar.Event.Err,
				},
			}
		}
	}

	return pipelineResult{
		Event: &report.PipelineEvent{
			DisplayEvent: report.DisplayEvent{Node: node, Name: name},
		},
	}
}
