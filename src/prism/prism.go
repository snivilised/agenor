package prism

import (
	"fmt"
	"io"
	"time"

	"github.com/snivilised/jaywalk/src/agenor/core"
)

// ViewKind identifies the rendering view to use. Defined as a typed
// string rather than an iota so that the set remains open - new views
// can be added in future without modifying this file.
type ViewKind string

const (
	// StreamView is a linear scrolling output view rendered with lipgloss.
	StreamView ViewKind = "stream"

	// PortholeView is a bubbletea view with a static header and footer
	// and vertically scrolling content between them.
	PortholeView ViewKind = "porthole"

	// LanesView is a bubbletea view showing parallel lanes of activity,
	// suited to concurrent worker output.
	LanesView ViewKind = "lanes"
)

// NavigationKind identifies whether a traversal is fresh or a
// continuation from a checkpoint.
type NavigationKind string

const (
	// PrimeNavigation is a fresh traversal from the root.
	PrimeNavigation NavigationKind = "prime"

	// ResumeNavigation is a continuation from a saved checkpoint.
	ResumeNavigation NavigationKind = "resume"
)

// SurveyResult carries the output of a two-phase navigation survey
// pass. Populated by controller/dispatch after the survey phase and
// passed to the renderer via Overture. Nil means single-phase
// navigation - no survey was performed.
type SurveyResult struct {
	// NodeCount is the total nodes to be visited in the execute phase.
	// Enables accurate progress reporting.
	NodeCount uint

	// MaxDepth is the deepest level seen during the survey phase.
	// Used by views for layout calculations.
	MaxDepth uint
}

// Overture carries the metadata known at the start of a traversal.
// Passed to Renderer.Begin to render the opening display.
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
// Depth is sourced from node.Extension.Level in agenor.
type Motif struct {
	// Path is the full path of the item.
	Path string

	// Name is the base name of the item.
	Name string

	// IsDir indicates whether the item is a directory.
	IsDir bool

	// Depth is the number of levels below the traversal root, sourced
	// from node.Extension.Level in agenor.
	Depth core.TraversalDepth

	// VisualDepth is the visual indent level for this item. For directories
	// this is the same as Depth, but for files it is Depth+1 since they are
	// visually one level deeper than their parent directory.
	VisualDepth core.TraversalDepth

	// ActionName is the name of the action executed against this node.
	// Empty when no action was configured.
	ActionName string

	// PipelineName is the name of the pipeline executed against this node.
	// Empty when no pipeline was configured.
	PipelineName string

	// ExecutionString is the expanded command string for dry-run display.
	ExecutionString string

	// CommandOutput is the captured output of the command execution.
	CommandOutput string

	// DryRun indicates if this is a dry-run execution.
	DryRun bool

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

	// IsLast is true when this is the last motif to be rendered in the traversal.
	// This can be used by the renderer to apply special styling or
	// behavior to the final item, such as a distinctive end marker or
	// summary display. Populated by the controller when it detects that
	// the current node is the last one in the traversal, which may be
	// determined by tracking the number of nodes rendered against the total
	// node count from a prior survey phase or by detecting completion of
	// the traversal process.
	IsLast bool
}

// Summary carries the result of a completed traversal. Passed to
// Renderer.End to render the closing display.
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
	// Begin is called once before any traversal events.
	Begin(overture Overture)

	// Show is called for each item encountered during traversal.
	Show(motif Motif)

	// End is called once when traversal completes.
	End(summary Summary)
}

// RendererOption configures renderer behavior at construction time.
type RendererOption func(*streamRenderer)

// WithIcons configures tree glyphs and item icons for stream renderers.
// The map keys are the standard token names used by prism tree rendering.
func WithIcons(icons map[string]string) RendererOption {
	return func(r *streamRenderer) {
		if r.treeIcons == nil {
			r.treeIcons = make(map[string]string)
		}

		for key, value := range icons {
			if value != "" {
				r.treeIcons[key] = value
			}
		}
	}
}

// New constructs a Renderer for the requested view kind using the given
// Palette. The Writer writer is the output destination; pass os.Stdout for
// production use or a bytes.Buffer in tests.
//
// Returns an error if the Palette contains unrecognised colour names,
// which would indicate a malformed user theme file. Bootstrap should
// treat this as a startup failure.
func New(kind ViewKind, palette Palette, writer io.Writer, opts ...RendererOption) (Renderer, error) {
	theme, err := NewTheme(palette, writer)
	if err != nil {
		return nil, fmt.Errorf("prism.New: %w", err)
	}

	//nolint:exhaustive // prism.PortholeView, prism.LanesView
	switch kind {
	case StreamView:
		return newStreamRenderer(theme, writer, opts...), nil
	default:
		return newStreamRenderer(theme, writer, opts...), nil
	}
}
