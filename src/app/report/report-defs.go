package report

import (
	"time"

	"github.com/snivilised/jaywalk/src/agenor/core"
)

// DisplayEvent is the base event embedded into all UI events. It carries
// the node that triggered the event and an optional name identifying the
// action or pipeline that was invoked. Name is empty for NodeEvent.
type DisplayEvent struct {
	Node *core.Node
	Name string
}

// NeutralEvent is emitted per node visit when no action or pipeline is
// configured. The UI decides how to render the node information.
type NeutralEvent struct {
	DisplayEvent
}

// ActionEvent is emitted when a configured action has been executed
// against a node. ExecutionString is the composed CLI string that was
// (or would be) run - population of this field is a future concern.
type ActionEvent struct {
	DisplayEvent
	ExecutionString string
	Err             error
}

// PipelineEvent is emitted when a configured pipeline has been executed
// against a node. A pipeline is a sequence of actions against the same
// node. ExecutionString is the composed CLI string - population of this
// field is a future concern.
type PipelineEvent struct {
	DisplayEvent
	ExecutionString string
	Err             error
}

// SkipEvent is emitted when an action is skipped for a node because a
// placeholder in the action's cmd string resolved to a path at or above
// the traversal root. The navigation continues; this node is not counted
// as a successful action invocation.
type SkipEvent struct {
	DisplayEvent

	// Placeholder is the token that caused the breach, e.g. "{{.grand}}".
	Placeholder string

	// ResolvedPath is the path the offending placeholder resolved to.
	ResolvedPath string
}

// Traversal captures the outcome of a completed directory traversal.
// It is populated by the controller and handed to the UI via OnComplete.
// The UI decides how to present each field - colour, layout, and
// formatting are entirely its own concern.
type Traversal struct {
	// FilesVisited is the number of file nodes invoked during traversal.
	FilesVisited core.MetricValue

	// DirsVisited is the number of directory nodes invoked during traversal.
	DirsVisited core.MetricValue

	// ActionsSkipped is the number of nodes for which an action was skipped
	// because a placeholder breached the traversal root. This field is
	// populated by jay; it is not sourced from agenor metrics.
	ActionsSkipped core.MetricValue

	// Elapsed is the total wall-clock time taken for the traversal.
	Elapsed time.Duration

	// Err holds the traversal error, if any. Nil on success.
	Err error
}
