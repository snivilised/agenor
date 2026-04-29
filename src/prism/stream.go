package prism

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss/v2"
)

// streamRenderer is the linear scrolling view. Output is written
// immediately as events arrive - no internal buffering. Implements
// Renderer.
type streamRenderer struct {
	theme Theme
	w     io.Writer
}

// newStreamRenderer constructs a streamRenderer. Called by New() when
// ViewKind is StreamView. Unexported - callers use New().
func newStreamRenderer(theme Theme, w io.Writer) Renderer {
	return &streamRenderer{
		theme: theme,
		w:     w,
	}
}

// Begin renders the opening banner using the Overture metadata. The
// banner adapts to indicate whether this is a prime or resume traversal.
func (r *streamRenderer) Begin(overture Overture) {
	var title string

	if overture.Kind == ResumeNavigation {
		title = fmt.Sprintf("  jay  resuming %s", overture.Root)
	} else {
		title = fmt.Sprintf("  jay  %s", overture.Root)
	}

	caption := fmt.Sprintf("  %s  -  %s",
		overture.Caption,
		overture.StartedAt.Format(time.RFC1123),
	)

	if overture.Kind == ResumeNavigation && overture.ResumeFrom != "" {
		caption += fmt.Sprintf("  -  from: %s", overture.ResumeFrom)
	}

	banner := r.theme.BannerStyle.Render(
		r.theme.BannerTitleStyle.Render(title) +
			"\n" +
			r.theme.BannerCaptionStyle.Render(caption),
	)

	_, _ = lipgloss.Fprintln(r.w, banner)
}

// Show renders a single Motif immediately to the output writer.
// Directories and files are styled differently. Actions, pipelines,
// skips, and errors each get distinct visual treatment.
func (r *streamRenderer) Show(motif Motif) {
	indent := r.indentFor(motif.Depth)
	depth := r.theme.DepthStyle.Render(indent)

	var name string

	switch {
	case motif.Err != nil:
		name = r.theme.ErrorStyle.Render(
			fmt.Sprintf("! %s  %s", motif.Name, motif.Err.Error()),
		)

	case motif.Skipped:
		name = r.theme.MutedStyle.Render(
			fmt.Sprintf("~ %s  [skipped: %s -> %s]",
				motif.Name,
				motif.Placeholder,
				motif.ResolvedPath,
			),
		)

	case motif.ActionName != "":
		name = r.theme.FileStyle.Render(motif.Name) +
			r.theme.MutedStyle.Render("  via "+motif.ActionName)

	case motif.PipelineName != "":
		name = r.theme.FileStyle.Render(motif.Name) +
			r.theme.MutedStyle.Render("  via "+motif.PipelineName)

	case motif.IsDir:
		name = r.theme.DirStyle.Render(motif.Name + "/")

	default:
		name = r.theme.FileStyle.Render(motif.Name)
	}

	_, _ = lipgloss.Fprintf(r.w, "%s%s\n", depth, name)
}

// End renders the closing summary box with traversal counts, elapsed
// time, and any errors encountered. Labels adapt for resume traversals.
func (r *streamRenderer) End(summary Summary) {
	fileLabel := "Files"
	dirLabel := "Directories"

	if summary.Kind == ResumeNavigation {
		fileLabel = "Files (resumed)"
		dirLabel = "Dirs (resumed)"
	}

	lines := []string{
		r.summaryRow(fileLabel, fmt.Sprintf("%d", summary.FilesVisited)),
		r.summaryRow(dirLabel, fmt.Sprintf("%d", summary.DirsVisited)),
		r.summaryRow("Elapsed", formatElapsed(summary.Elapsed)),
	}

	if len(summary.Errors) > 0 {
		lines = append(lines,
			r.summaryRow("Errors",
				r.theme.ErrorStyle.Render(fmt.Sprintf("%d", len(summary.Errors))),
			),
		)

		for _, err := range summary.Errors {
			lines = append(lines,
				r.theme.ErrorStyle.Render("  ! "+err.Error()),
			)
		}
	}

	box := r.theme.SummaryBoxStyle.Render(strings.Join(lines, "\n"))
	_, _ = lipgloss.Fprintln(r.w, box)
}

// summaryRow renders a label/value pair aligned inside the summary box.
func (r *streamRenderer) summaryRow(label, value string) string {
	return r.theme.SummaryLabelStyle.Render(label) +
		r.theme.SummaryValueStyle.Render(value)
}

// indentFor returns the indent prefix string for the given depth level.
// Depth 0 (the root) produces no indent. Each subsequent level adds two
// spaces. Tree-branch glyphs (such as, ├── and └──) are deferred until
// sibling tracking is available from the agenor side.
func (r *streamRenderer) indentFor(depth uint) string {
	if depth == 0 {
		return ""
	}

	return strings.Repeat("  ", int(depth)) //nolint:gosec // overflow not likely
}

// formatElapsed produces a human-readable elapsed duration string.
// Free function - no dependency on renderer state.
func formatElapsed(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}

	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}

	m := int(d.Minutes())
	s := int(d.Seconds()) % 60

	return fmt.Sprintf("%dm%ds", m, s)
}
