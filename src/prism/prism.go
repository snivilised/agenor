package prism

import (
	"io"
	"time"
)

// ViewKind identifies the rendering view to use. Defined as a typed string
// rather than an iota so that the set remains open - new views can be added
// in future without modifying this file.
type ViewKind string

const (
	// StreamView is a linear scrolling output view rendered with lipgloss.
	StreamView ViewKind = "stream"

	// PortholeView is a bubbletea view with a static header and footer and
	// vertically scrolling content between them.
	PortholeView ViewKind = "porthole"

	// LanesView is a bubbletea view showing parallel lanes of activity,
	// suited to concurrent worker output.
	LanesView ViewKind = "lanes"
)

// NavigationKind identifies whether a traversal is fresh or a continuation.
type NavigationKind string

const (
	// PrimeNavigation is a fresh traversal from the root.
	PrimeNavigation NavigationKind = "prime"

	// ResumeNavigation is a continuation from a saved checkpoint.
	ResumeNavigation NavigationKind = "resume"
)

// SurveyResult carries the output of a two-phase navigation survey pass.
// It is populated by the controller/dispatch layer after the survey phase
// completes and passed to the renderer via Overture so that views can
// display accurate progress and use max depth for layout calculations.
// Nil means no survey was performed - single-phase navigation only.
type SurveyResult struct {
	// NodeCount is the total number of nodes that will be visited in the
	// execute phase. Enables accurate progress reporting.
	NodeCount uint

	// MaxDepth is the deepest level encountered during the survey phase.
	// Used by views for layout calculations such as indent budget and
	// lane column widths.
	MaxDepth uint
}

// Overture carries the metadata known at the start of a traversal. It is
// passed to Renderer.Begin and used to render the opening display.
type Overture struct {
	// Root is the top-level path being traversed.
	Root string

	// Caption is a human-readable description of the traversal options,
	// e.g. "files and folders".
	Caption string

	// StartedAt is the time the traversal began.
	StartedAt time.Time

	// Kind indicates whether this is a prime or resume traversal.
	Kind NavigationKind

	// ResumeFrom is the path from which a resume traversal continues.
	// Populated only when Kind == ResumeNavigation.
	ResumeFrom string

	// Survey holds the results of a prior survey phase. Nil for
	// single-phase navigations such as the stream view.
	Survey *SurveyResult
}

// Motif is the unit of render-able content passed to Renderer.Show for
// each item encountered during traversal. Fields are generic filesystem
// and execution concepts - no jaywalk or agenor types appear here.
// Depth is sourced from node.Extension.Depth in agenor.
type Motif struct {
	// Path is the full path of the item.
	Path string

	// Name is the base name of the item.
	Name string

	// IsDir indicates whether the item is a directory.
	IsDir bool

	// Depth is the number of levels below the traversal root, sourced
	// from node.Extension.Depth in agenor.
	Depth uint

	// ActionName is the name of the action executed against this node.
	// Empty when no action was configured.
	ActionName string

	// PipelineName is the name of the pipeline executed against this node.
	// Empty when no pipeline was configured.
	PipelineName string

	// Skipped is true when an action or pipeline was skipped because a
	// placeholder resolved to a path at or above the traversal root.
	Skipped bool

	// Placeholder is the placeholder string that caused the skip.
	// Populated only when Skipped is true.
	Placeholder string

	// ResolvedPath is the path the placeholder resolved to.
	// Populated only when Skipped is true.
	ResolvedPath string

	// Err is any error produced by the action or pipeline for this node.
	// Nil when the node was visited without error.
	Err error
}

// Summary carries the result of a completed traversal. It is passed to
// Renderer.End and used to render the closing display.
type Summary struct {
	// FilesVisited is the count of files encountered.
	FilesVisited uint

	// DirsVisited is the count of directories encountered.
	DirsVisited uint

	// Elapsed is the total duration of the traversal.
	Elapsed time.Duration

	// Errors contains any errors encountered during traversal.
	Errors []error

	// Kind mirrors Overture.Kind so the renderer can label counts
	// appropriately in the closing summary.
	Kind NavigationKind
}

// Renderer is the rendering abstraction for prism views. All view kinds
// implement this interface. Callers depend on Renderer, never on a
// concrete view type.
type Renderer interface {
	// Begin is called once before any traversal events, with the opening
	// metadata.
	Begin(overture Overture)

	// Show is called for each item encountered during traversal.
	Show(motif Motif)

	// End is called once when traversal completes, with the result summary.
	End(summary Summary)
}

// New constructs a Renderer for the requested view kind. The writer w is
// the output destination; pass os.Stdout for production use or a
// bytes.Buffer in tests. Dark/light detection and colour downsampling
// are handled automatically by lipgloss v2 against w - no colour profile
// need be supplied by the caller.
func New(kind ViewKind, w io.Writer) Renderer {
	theme := NewTheme(w)

	//nolint:exhaustive // prism.PortholeView, prism.LanesView
	switch kind {
	case StreamView:
		return newStreamRenderer(theme, w)
	default:
		return newStreamRenderer(theme, w)
	}
}
