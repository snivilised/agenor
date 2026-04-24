package controller

import (
	"fmt"
	"strings"
)

// PreFlight validates that the executable named in the first token of
// each action's cmd string can be located on PATH before the traversal
// begins. This prevents the same "executable not found" error from being
// emitted once per matched node during a walk.
//
// Rules:
//   - If neither ActionName nor PipelineName is set, PreFlight is a no-op.
//   - If ActionName is set, the action is looked up in cfg.Raw.Actions and
//     its cmd token[0] is verified via lookPath.
//   - If PipelineName is set, every step in the pipeline is looked up in
//     cfg.Raw.Actions and each cmd token[0] is verified. The first failure
//     aborts the check immediately.
func (c *Coordinator) PreFlight(req *Request) error {
	switch {
	case req.PipelineName != "":
		return c.preFlightPipeline(req.PipelineName)

	case req.ActionName != "":
		return c.preFlightAction(req.ActionName)

	default:
		return nil
	}
}

// preFlightAction verifies a single named action.
func (c *Coordinator) preFlightAction(name string) error {
	action, ok := c.cfg.Raw.Actions[name]
	if !ok {
		return fmt.Errorf("action %q is not defined in config", name)
	}

	return c.preFlightCmd(name, action.Cmd)
}

// preFlightPipeline verifies every step in a named pipeline.
func (c *Coordinator) preFlightPipeline(name string) error {
	pipeline, ok := c.cfg.Raw.Pipelines[name]
	if !ok {
		return fmt.Errorf("pipeline %q is not defined in config", name)
	}

	for _, step := range pipeline.Steps {
		if err := c.preFlightAction(step); err != nil {
			return fmt.Errorf("pipeline %q: %w", name, err)
		}
	}

	return nil
}

// preFlightCmd extracts the first token from a cmd string and verifies
// it can be found on PATH. An empty cmd string is an immediate error.
func (c *Coordinator) preFlightCmd(actionName, cmd string) error {
	fields := strings.Fields(cmd)
	if len(fields) == 0 {
		return fmt.Errorf("action %q has an empty cmd string", actionName)
	}

	executable := fields[0]

	if _, err := c.lookPath(executable); err != nil {
		return fmt.Errorf(
			"action %q: executable %q not found on PATH: %w",
			actionName, executable, err,
		)
	}

	return nil
}
