package controller

import (
	"context"
	"regexp"
	"strings"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/locale"
)

// handleServant dispatches to the appropriate per-node handler based on
// whether an action, pipeline, or neither is configured on the request.
func (c *Coordinator) handleServant(
	ctx context.Context,
	servant core.Servant,
	req *Request,
	traversal *report.Traversal,
	peerInfoMap PeerInfoMap,
) error {
	node := servant.Node()
	_ = ctx

	var isLast bool

	if peerInfoMap != nil {
		if info, ok := peerInfoMap[node.Path]; ok {
			isLast = info.IsLast
		}
	}

	switch {
	case req.PipelineName != "":
		e := c.executePipeline(ctx, node, req.PipelineName, req.Root, req.DryRun)
		if e.Skipped {
			traversal.ActionsSkipped.Tick()
			req.UI.OnSkipEvent(&report.SkipEvent{
				DisplayEvent: report.DisplayEvent{
					Node:   node,
					IsLast: isLast,
					Name:   req.PipelineName,
				},
				Placeholder:  e.Placeholder,
				ResolvedPath: e.ResolvedPath,
			})
			return nil
		}
		e.Event.IsLast = isLast
		req.UI.OnPipelineEvent(e.Event)
		return e.Event.Err

	case req.ActionName != "":
		e := c.executeAction(ctx, node, req.ActionName, req.Root, req.DryRun)
		if e.Skipped {
			traversal.ActionsSkipped.Tick()
			req.UI.OnSkipEvent(&report.SkipEvent{
				DisplayEvent: report.DisplayEvent{
					Node:   node,
					IsLast: isLast,
					Name:   req.ActionName,
				},
				Placeholder:  e.Placeholder,
				ResolvedPath: e.ResolvedPath,
			})
			return nil
		}
		e.Event.IsLast = isLast
		req.UI.OnActionEvent(e.Event)
		return e.Event.Err

	default:
		req.UI.OnNodeEvent(&report.NeutralEvent{
			DisplayEvent: report.DisplayEvent{
				Node:   node,
				IsLast: isLast,
			},
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
func (c *Coordinator) executeAction(
	ctx context.Context,
	node *core.Node,
	name, root string,
	dryRun bool,
) actionResult {
	action, ok := c.config.Raw.Actions[name]
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

	event := &report.ActionEvent{
		DisplayEvent:    report.DisplayEvent{Node: node, Name: name},
		ExecutionString: result.Cmd,
		DryRun:          dryRun,
	}

	if !dryRun {
		// cmd := exec.Command("sh", "-c", result.Cmd) //nolint:gosec // this is expected to be a shell command string
		// output, err := cmd.CombinedOutput()
		output, err := c.exec(ctx, event.ExecutionString)
		if err != nil {
			event.Err = err
		}
		event.CommandOutput = c.processOutput(output, c.actionRegexes[name])
	}

	return actionResult{
		Event: event,
	}
}

// executePipeline expands and executes each step in the named pipeline
// in order. The first skip or error stops the pipeline for this node.
func (c *Coordinator) executePipeline(ctx context.Context,
	node *core.Node,
	name, root string,
	dryRun bool,
) pipelineResult {
	pipeline, ok := c.config.Raw.Pipelines[name]
	if !ok {
		return pipelineResult{
			Event: &report.PipelineEvent{
				DisplayEvent: report.DisplayEvent{Node: node, Name: name},
				Err:          locale.NewPipelineNotFoundError(name),
			},
		}
	}

	for _, step := range pipeline.Steps {
		ar := c.executeAction(ctx, node, step, root, dryRun)
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
					CommandOutput:   ar.Event.CommandOutput,
					DryRun:          dryRun,
					Err:             ar.Event.Err,
				},
			}
		}
	}

	return pipelineResult{
		Event: &report.PipelineEvent{
			DisplayEvent: report.DisplayEvent{Node: node, Name: name},
			DryRun:       dryRun,
		},
	}
}

const (
	minLimit     = 20
	maxLimit     = 120
	defaultLimit = 75
	ellipsis     = " ..."
)

// processOutput extracts a single line from the raw command output and applies truncation.
// It removes leading/trailing empty lines. If captureRe is provided, it uses it
// to select the matching line.
func (c *Coordinator) processOutput(output []byte, captureRe *regexp.Regexp) string {
	lines := strings.Split(string(output), "\n")
	var contentLines []string
	for _, ln := range lines {
		trimmed := strings.TrimSpace(ln)
		if trimmed != "" {
			contentLines = append(contentLines, trimmed)
		}
	}

	if len(contentLines) == 0 {
		return ""
	}

	selectedLine := contentLines[0]

	if captureRe != nil {
		for _, ln := range contentLines {
			if captureRe.MatchString(ln) {
				selectedLine = ln
				break
			}
		}
	}

	limit := c.config.Mapped.Advanced.Output.Exec.Truncate
	if limit < minLimit || limit > maxLimit {
		limit = defaultLimit
	}

	if len(selectedLine) > limit {
		if limit > 4 {
			selectedLine = selectedLine[:limit-len(ellipsis)] + ellipsis
		} else {
			selectedLine = selectedLine[:limit]
		}
	}

	return selectedLine
}
