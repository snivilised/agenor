package prism

import (
	"fmt"
	"io"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/snivilised/jaywalk/src/third/lo"
)

// streamRenderer is the linear scrolling view. Output is written
// immediately as events arrive - no internal buffering. Implements
// Renderer.
type streamRenderer struct {
	theme Theme
	w     io.Writer

	// treeIcons holds configured tree glyphs, either from the theme or
	// from renderer options such as WithIcons.
	treeIcons TreeIcons

	// branchStack tracks ancestor continuation state for tree branch
	// rendering.
	branchStack []bool

	previousDepth  uint
	previousIsLast bool
}

// newStreamRenderer constructs a streamRenderer. Called by New() when
// ViewKind is StreamView. Unexported - callers use New().
func newStreamRenderer(theme Theme, w io.Writer, opts ...RendererOption) Renderer {
	r := &streamRenderer{
		theme:     theme,
		w:         w,
		treeIcons: theme.TreeIcons,
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

// Begin renders the opening banner using the Overture metadata. The
// banner adapts to indicate prime or resume traversal.
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
// Errors, skips, actions, pipelines, directories and files each receive
// distinct visual treatment.
func (r *streamRenderer) Show(motif Motif) {
	prefix := r.branchPrefix(motif)
	depth := r.theme.DepthStyle.Render(prefix)

	var name string

	switch {
	case motif.Err != nil:
		name = r.theme.ErrorStyle.Render(
			fmt.Sprintf("! %s  %s", motif.Name, motif.Err.Error()),
		)

	case motif.Skipped:
		skippedInfo := fmt.Sprintf("  [skipped: %s -> %s]",
			motif.Placeholder,
			motif.ResolvedPath,
		)
		if motif.IsDir {
			name = r.theme.DirStyle.Render("~ "+motif.Name+"/") +
				r.theme.SkippedStyle.Render(skippedInfo)
		} else {
			name = r.theme.FileStyle.Render("~ "+motif.Name) +
				r.theme.SkippedStyle.Render(skippedInfo)
		}

	case motif.Depth == 0:
		name = r.theme.RootStyle.Render(
			r.treeIcons[TreeIconRoot] +
				lo.Ternary(r.treeIcons[TreeIconRoot] != "", " ", "") +
				motif.Name +
				lo.Ternary(motif.IsDir, "/", ""),
		)

	case motif.IsDir:
		name = r.theme.DirStyle.Render(r.itemLabel(motif))

		if motif.ActionName != "" {
			name += r.theme.ActionStyle.Render("  • via " + motif.ActionName)
		} else if motif.PipelineName != "" {
			name += r.theme.PipelineStyle.Render("  • via " + motif.PipelineName)
		}

	default:
		name = r.theme.FileStyle.Render(r.itemLabel(motif))

		if motif.ActionName != "" {
			name += r.theme.ActionStyle.Render("  • via " + motif.ActionName)
		} else if motif.PipelineName != "" {
			name += r.theme.PipelineStyle.Render("  • via " + motif.PipelineName)
		}
	}

	_, _ = lipgloss.Fprintf(r.w, "%s%s\n", depth, name)

	r.updateBranchStack(motif)
}

func (r *streamRenderer) itemLabel(motif Motif) string {
	icon := r.treeIcons[TreeIconFile]
	if motif.IsDir {
		icon = r.treeIcons[TreeIconDirectory]
	}

	label := ""
	if icon != "" {
		label = icon + " "
	}
	label += motif.Name
	if motif.IsDir {
		label += "/"
	}

	return label
}

func (r *streamRenderer) branchPrefix(motif Motif) string {
	if motif.VisualDepth == 0 {
		return ""
	}

	var b strings.Builder
	//nolint:gosec // ok - branchStack is only modified by updateBranchStack based on motif.VisualDepth
	for level := 1; level < int(motif.VisualDepth); level++ {
		if level-1 < len(r.branchStack) && r.branchStack[level-1] {
			b.WriteString(r.treeIcons[TreeIconBranchVertical])
			b.WriteString(r.treeIcons[TreeIconBranchIndent])
		} else {
			b.WriteString(
				strings.Repeat(" ",
					len(r.treeIcons[TreeIconBranchVertical])+len(r.treeIcons[TreeIconBranchIndent]),
				),
			)
		}
	}

	branchIcon := lo.Ternary(motif.IsLast, TreeIconBranchLast, TreeIconBranchJoint)
	b.WriteString(r.treeIcons[branchIcon])

	return b.String()
}

func (r *streamRenderer) updateBranchStack(motif Motif) {
	if motif.VisualDepth == 0 {
		r.branchStack = nil
		r.previousDepth = motif.VisualDepth
		r.previousIsLast = motif.IsLast
		return
	}

	if motif.VisualDepth > r.previousDepth {
		for d := r.previousDepth; d < motif.VisualDepth; d++ {
			r.branchStack = append(r.branchStack, !motif.IsLast)
		}
	} else if motif.VisualDepth < r.previousDepth {
		r.branchStack = r.branchStack[:motif.VisualDepth]
	}

	if motif.VisualDepth > 0 {
		r.branchStack[motif.VisualDepth-1] = !motif.IsLast
	}

	r.previousDepth = motif.VisualDepth
	r.previousIsLast = motif.IsLast
}

// End renders the closing summary box with traversal counts, elapsed
// time, and any errors. Labels adapt for resume traversals.
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
				r.theme.ErrorStyle.Render(
					fmt.Sprintf("%d", len(summary.Errors)),
				),
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
