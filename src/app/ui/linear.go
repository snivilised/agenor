package ui

import (
	"sync"

	"github.com/snivilised/jaywalk/src/app/report"
	"github.com/snivilised/jaywalk/src/prism"
)

// linear is the stream-view display implementation. It translates
// report events into prism.Motif calls and delegates all formatting
// and output to the prism.Renderer. It contains no formatting logic.
//
// Safe for concurrent use - all renderer calls are serialised through
// a mutex so interleaved output from the run command's worker pool is
// avoided.
type linear struct {
	mu       sync.Mutex
	renderer prism.Renderer
	kind     prism.NavigationKind // remembered from OnBegin for use in OnComplete
}

// OnBegin translates the BeginEvent into a prism.Overture and calls
// renderer.Begin to render the opening banner.
func (l *linear) OnBegin(e *report.BeginEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	kind := prism.PrimeNavigation
	if !e.IsPrime {
		kind = prism.ResumeNavigation
	}

	l.kind = kind

	l.renderer.Begin(prism.Overture{
		Root:       e.Root,
		Caption:    e.Caption,
		StartedAt:  e.StartedAt,
		Kind:       kind,
		ResumeFrom: e.ResumeFrom,
	})
}

// OnNodeEvent translates a neutral node visit into a prism.Motif.
// Depth is sourced from node.Extension.Depth as provided by agenor.
func (l *linear) OnNodeEvent(e *report.NeutralEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.renderer.Show(prism.Motif{
		Path:  e.Node.Path,
		Name:  e.Node.Extension.Name,
		IsDir: e.Node.IsDirectory(),
		Depth: uint(e.Node.Extension.Depth), //nolint: gosec // overflow
	})
}

// OnActionEvent translates an action event into a prism.Motif.
func (l *linear) OnActionEvent(e *report.ActionEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.renderer.Show(prism.Motif{
		Path:       e.Node.Path,
		Name:       e.Node.Extension.Name,
		IsDir:      e.Node.IsDirectory(),
		Depth:      uint(e.Node.Extension.Depth), //nolint: gosec // overflow
		ActionName: e.Name,
		Err:        e.Err,
	})
}

// OnPipelineEvent translates a pipeline event into a prism.Motif.
func (l *linear) OnPipelineEvent(e *report.PipelineEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.renderer.Show(prism.Motif{
		Path:         e.Node.Path,
		Name:         e.Node.Extension.Name,
		IsDir:        e.Node.IsDirectory(),
		Depth:        uint(e.Node.Extension.Depth), //nolint: gosec // overflow
		PipelineName: e.Name,
		Err:          e.Err,
	})
}

// OnSkipEvent translates a skip event into a prism.Motif flagged as
// skipped so the renderer can apply warning styling.
func (l *linear) OnSkipEvent(e *report.SkipEvent) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.renderer.Show(prism.Motif{
		Path:         e.Node.Path,
		Name:         e.Node.Extension.Name,
		IsDir:        e.Node.IsDirectory(),
		Depth:        uint(e.Node.Extension.Depth), //nolint: gosec // overflow
		ActionName:   e.Name,
		Skipped:      true,
		Placeholder:  e.Placeholder,
		ResolvedPath: e.ResolvedPath,
	})
}

// OnComplete translates the Traversal outcome into a prism.Summary and
// calls renderer.End to render the closing summary box. Kind is carried
// from OnBegin so the summary labels correctly for resume traversals.
func (l *linear) OnComplete(t *report.Traversal) {
	l.mu.Lock()
	defer l.mu.Unlock()

	errs := []error{}
	if t.Err != nil {
		errs = append(errs, t.Err)
	}

	l.renderer.End(prism.Summary{
		FilesVisited: t.FilesVisited,
		DirsVisited:  t.DirsVisited,
		Elapsed:      t.Elapsed,
		Errors:       errs,
		Kind:         l.kind,
	})
}

// NewLinearWithRenderer constructs a linear presenter backed by the
// given renderer. Intended for use in tests only - production code
// constructs linear via the ui registry using New(). This allows a
// spy or stub renderer to be injected without going through prism.New
// and without requiring a real terminal.
func NewLinearWithRenderer(r prism.Renderer) report.Presenter {
	return &linear{renderer: r}
}
