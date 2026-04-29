package prism

import (
	"io"
	"os"

	"github.com/charmbracelet/lipgloss/v2"
)

// Theme holds all lipgloss styles used by a renderer. Constructed once via
// NewTheme and injected into the renderer - no package-level style variables
// exist. Dark/light palette selection is handled automatically by lipgloss v2.
type Theme struct {
	// BannerStyle is applied to the overture banner block.
	BannerStyle lipgloss.Style

	// BannerTitleStyle is applied to the title text inside the banner.
	BannerTitleStyle lipgloss.Style

	// BannerCaptionStyle is applied to the caption line inside the banner.
	BannerCaptionStyle lipgloss.Style

	// DirStyle is applied to directory name text in Show output.
	DirStyle lipgloss.Style

	// FileStyle is applied to file name text in Show output.
	FileStyle lipgloss.Style

	// DepthStyle is applied to the depth indent prefix in Show output.
	DepthStyle lipgloss.Style

	// SummaryBoxStyle is the outer container style for the closing summary.
	SummaryBoxStyle lipgloss.Style

	// SummaryLabelStyle is applied to label text inside the summary.
	SummaryLabelStyle lipgloss.Style

	// SummaryValueStyle is applied to value text inside the summary.
	SummaryValueStyle lipgloss.Style

	// ErrorStyle is applied to error entries in Show output and the summary.
	ErrorStyle lipgloss.Style

	// MutedStyle is applied to secondary or de-emphasised text.
	MutedStyle lipgloss.Style
}

// NewTheme constructs a Theme by detecting the terminal's dark/light
// background via lipgloss v2 against the provided writer. If detection
// fails (e.g. in a test with a non-TTY writer) dark mode is assumed as
// a safe fallback - colours are still applied but downsampled to nothing
// by lipgloss when output is not a TTY.
//
// Colour downsampling (TrueColor -> ANSI256 -> ANSI -> none) is handled
// automatically by lipgloss.Fprintln and friends at render time.
func NewTheme(w io.Writer) Theme {
	// HasDarkBackground returns (bool, error) in lipgloss v2. Errors occur
	// when the terminal cannot be queried (non-TTY, piped output, tests).
	// We fall back to dark rather than propagating the error - a rendering
	// library should never force callers into error handling for cosmetic
	// decisions.

	hasDark := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	ld := lipgloss.LightDark(hasDark)

	// Colour palette - one declaration per semantic role, light/dark
	// variants resolved at construction time by ld().
	banner := ld(
		lipgloss.Color("#7C3AED"), // light: vivid violet
		lipgloss.Color("#5C4AE4"), // dark:  indigo
	)
	bannerText := lipgloss.Color("#FFFFFF")

	dirFg := ld(
		lipgloss.Color("#0369A1"), // light: ocean blue
		lipgloss.Color("#89DCEB"), // dark:  sky cyan
	)
	fileFg := ld(
		lipgloss.Color("#1E293B"), // light: near-black slate
		lipgloss.Color("#CDD6F4"), // dark:  lavender white
	)
	depthFg := ld(
		lipgloss.Color("#94A3B8"), // light: muted slate
		lipgloss.Color("#6C7086"), // dark:  muted overlay
	)
	summaryBorder := ld(
		lipgloss.Color("#7C3AED"),
		lipgloss.Color("#5C4AE4"),
	)
	summaryLabel := ld(
		lipgloss.Color("#1D4ED8"), // light: strong blue
		lipgloss.Color("#89B4FA"), // dark:  soft blue
	)
	summaryValue := ld(
		lipgloss.Color("#1E293B"),
		lipgloss.Color("#CDD6F4"),
	)
	errorFg := ld(
		lipgloss.Color("#DC2626"), // light: red
		lipgloss.Color("#F38BA8"), // dark:  rose
	)
	muted := ld(
		lipgloss.Color("#94A3B8"),
		lipgloss.Color("#6C7086"),
	)

	return Theme{
		BannerStyle: lipgloss.NewStyle().
			Background(banner).
			Padding(0, 2).
			MarginBottom(1),

		BannerTitleStyle: lipgloss.NewStyle().
			Foreground(bannerText).
			Background(banner).
			Bold(true),

		BannerCaptionStyle: lipgloss.NewStyle().
			Foreground(bannerText).
			Background(banner).
			Faint(true),

		DirStyle: lipgloss.NewStyle().
			Foreground(dirFg).
			Bold(true),

		FileStyle: lipgloss.NewStyle().
			Foreground(fileFg),

		DepthStyle: lipgloss.NewStyle().
			Foreground(depthFg).
			Faint(true),

		SummaryBoxStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(summaryBorder).
			Padding(0, 2).
			MarginTop(1),

		SummaryLabelStyle: lipgloss.NewStyle().
			Foreground(summaryLabel).
			Bold(true).
			Width(16),

		SummaryValueStyle: lipgloss.NewStyle().
			Foreground(summaryValue),

		ErrorStyle: lipgloss.NewStyle().
			Foreground(errorFg).
			Bold(true),

		MutedStyle: lipgloss.NewStyle().
			Foreground(muted).
			Faint(true),
	}
}
