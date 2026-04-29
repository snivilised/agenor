package report

import "time"

// BeginEvent carries the metadata known at the start of a traversal.
// It is passed to Presenter.OnBegin before any node events are fired.
type BeginEvent struct {
	// Root is the top-level path being traversed.
	Root string

	// Caption is a human-readable description of the traversal options.
	Caption string

	// StartedAt is the time the traversal began.
	StartedAt time.Time

	// IsPrime indicates whether this is a fresh traversal. False means
	// this is a resume from a checkpoint.
	IsPrime bool

	// ResumeFrom is the path from which a resume traversal continues.
	// Populated only when IsPrime is false.
	ResumeFrom string
}

// Presenter is the interface all display implementations satisfy.
// It is purely reactive - all methods are event notifications. The
// implementation decides how to render each event; no formatting logic
// lives outside the implementing type.
//
// Presenter is defined here, alongside the event types it handles, so
// that both the ui package and the prism package can depend on report
// without creating a circular dependency between them.
type Presenter interface {
	// OnBegin is called once before any traversal events, with the
	// opening metadata. Implementations should use this to render
	// any opening banner or header.
	OnBegin(e *BeginEvent)

	// OnNodeEvent is called per node visit when no action or pipeline
	// is configured.
	OnNodeEvent(e *NeutralEvent)

	// OnActionEvent is called when a configured action has been executed
	// against a node.
	OnActionEvent(e *ActionEvent)

	// OnPipelineEvent is called when a configured pipeline has been
	// executed against a node.
	OnPipelineEvent(e *PipelineEvent)

	// OnSkipEvent is called when an action is skipped for a node because
	// a placeholder in the action's cmd string resolved to a path at or
	// above the traversal root.
	OnSkipEvent(e *SkipEvent)

	// OnComplete is called once at the end of a traversal with the full
	// structured outcome.
	OnComplete(t *Traversal)
}
