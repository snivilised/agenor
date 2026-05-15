// Package widget contains reusable terminal UI components for prism views.
package widget

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/mattn/go-runewidth"
)

// SummaryCellStyles defines the styles used to render one summary row.
type SummaryCellStyles struct {
	Icon  lipgloss.Style
	Label lipgloss.Style
	Value lipgloss.Style
}

// SummaryStyles defines the styles used to render a boxed summary table.
type SummaryStyles struct {
	Box lipgloss.Style

	Default SummaryCellStyles
}

// SummaryRow is one row in a boxed summary table.
type SummaryRow struct {
	Icon  string
	Label string
	Value string

	Styles *SummaryCellStyles
}

// RenderSummary renders rows as a boxed, three-column summary table. The icon,
// label, and value columns are measured independently using runewidth so that
// emoji glyph widths do not disturb value alignment or the surrounding box.
func RenderSummary(rows []SummaryRow, styles SummaryStyles) string {
	if len(rows) == 0 {
		return styles.Box.Render("")
	}

	widths := summaryWidths(rows)
	lines := make([]string, 0, len(rows))

	for _, row := range rows {
		cs := stylesForRow(row, styles.Default)
		var b strings.Builder

		// Icon cell — pad to widths.icon + gap.
		if widths.icon > 0 {
			b.WriteString(cs.Icon.Render(row.Icon))

			if pad := widths.icon - runewidth.StringWidth(row.Icon) + 1; pad > 0 {
				b.WriteString(strings.Repeat(" ", pad))
			}
		}

		// Label cell — pad to widths.label + gap.
		b.WriteString(cs.Label.Render(row.Label))

		if pad := widths.label - runewidth.StringWidth(row.Label) + 2; pad > 0 {
			b.WriteString(strings.Repeat(" ", pad))
		}

		// Value cell — right-aligned within widths.value.
		if vpad := widths.value - runewidth.StringWidth(row.Value); vpad > 0 {
			b.WriteString(strings.Repeat(" ", vpad))
		}

		b.WriteString(cs.Value.Render(row.Value))

		lines = append(lines, b.String())
	}

	return styles.Box.Render(strings.Join(lines, "\n"))
}

type summaryColumnWidths struct {
	icon  int
	label int
	value int
}

func summaryWidths(rows []SummaryRow) summaryColumnWidths {
	widths := summaryColumnWidths{}
	for _, row := range rows {
		widths.icon = max(widths.icon, runewidth.StringWidth(row.Icon))
		widths.label = max(widths.label, runewidth.StringWidth(row.Label))
		widths.value = max(widths.value, runewidth.StringWidth(row.Value))
	}

	return widths
}

func stylesForRow(row SummaryRow, fallback SummaryCellStyles) SummaryCellStyles {
	cs := fallback
	if row.Styles != nil {
		cs = *row.Styles
	}

	// Clear any Width set on the incoming styles — this widget handles
	// column alignment via runewidth, so letting lipgloss pad internally
	// would produce double-padding.
	cs.Icon = cs.Icon.Width(0)
	cs.Label = cs.Label.Width(0)
	cs.Value = cs.Value.Width(0)

	return cs
}
