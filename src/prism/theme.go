package prism

import (
	"fmt"
	"image/color"
	"io"

	"github.com/charmbracelet/lipgloss/v2"
)

// Theme holds all lipgloss styles used by a renderer. Constructed once
// via NewTheme and injected into the renderer. No package-level style
// variables exist anywhere in prism.
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

	// ActionStyle is applied to action names shown alongside nodes.
	ActionStyle lipgloss.Style

	// PipelineStyle is applied to pipeline names shown alongside nodes.
	PipelineStyle lipgloss.Style

	// SkippedStyle is applied to nodes whose action was skipped.
	SkippedStyle lipgloss.Style

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

	// ProgressStyle is applied to progress indicators.
	ProgressStyle lipgloss.Style

	// WorkerStyle is applied to active concurrent worker indicators.
	WorkerStyle lipgloss.Style

	// WorkerIdleStyle is applied to idle concurrent worker indicators.
	WorkerIdleStyle lipgloss.Style

	// LaneHeaderStyle is applied to per-worker lane identity headers.
	LaneHeaderStyle lipgloss.Style
}

// NewTheme constructs a Theme from the given Palette. The colour profile
// is detected from the environment via colorprofile.Detect so that
// lipgloss can select the best available colour tier (TrueColor,
// ANSI-256, ANSI-16, or none) at style construction time.
//
// Returns an error if any palette entry contains an unrecognised ANSI-16
// colour name. Bootstrap should treat this as a startup failure.
func NewTheme(palette Palette, w io.Writer) (Theme, error) {
	// Detect the terminal colour profile from the writer's environment.
	// This determines whether TrueColor, ANSI-256, ANSI-16, or plain
	// text output is used. colorprofile.Detect handles NO_COLOR, isatty,
	// COLORTERM and TERM inspection automatically.

	resolve := func(sc SemanticColour, field string) (color.Color, error) {
		ansi, ansi256, trueCol, err := sc.Resolve()
		if err != nil {
			return nil, fmt.Errorf("palette.%s: %w", field, err)
		}

		switch {
		case trueCol != nil:
			return trueCol, nil
		case ansi256 != nil:
			return ansi256, nil
		default:
			return ansi, nil
		}
	}
	banner, err := resolve(palette.Banner, "banner")
	if err != nil {
		return Theme{}, err
	}

	bannerText, err := resolve(palette.BannerText, "banner-text")
	if err != nil {
		return Theme{}, err
	}

	dir, err := resolve(palette.Directory, "directory")
	if err != nil {
		return Theme{}, err
	}

	file, err := resolve(palette.File, "file")
	if err != nil {
		return Theme{}, err
	}

	action, err := resolve(palette.Action, "action")
	if err != nil {
		return Theme{}, err
	}

	pipeline, err := resolve(palette.Pipeline, "pipeline")
	if err != nil {
		return Theme{}, err
	}

	skipped, err := resolve(palette.Skipped, "skipped")
	if err != nil {
		return Theme{}, err
	}

	errorCol, err := resolve(palette.Error, "error")
	if err != nil {
		return Theme{}, err
	}

	muted, err := resolve(palette.Muted, "muted")
	if err != nil {
		return Theme{}, err
	}

	progress, err := resolve(palette.Progress, "progress")
	if err != nil {
		return Theme{}, err
	}

	summaryBorder, err := resolve(palette.SummaryBorder, "summary-border")
	if err != nil {
		return Theme{}, err
	}

	summaryLabel, err := resolve(palette.SummaryLabel, "summary-label")
	if err != nil {
		return Theme{}, err
	}

	summaryValue, err := resolve(palette.SummaryValue, "summary-value")
	if err != nil {
		return Theme{}, err
	}

	worker, err := resolve(palette.Worker, "worker")
	if err != nil {
		return Theme{}, err
	}

	workerIdle, err := resolve(palette.WorkerIdle, "worker-idle")
	if err != nil {
		return Theme{}, err
	}

	laneHeader, err := resolve(palette.LaneHeader, "lane-header")
	if err != nil {
		return Theme{}, err
	}

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
			Foreground(dir).
			Bold(true),

		FileStyle: lipgloss.NewStyle().
			Foreground(file),

		DepthStyle: lipgloss.NewStyle().
			Foreground(muted).
			Faint(true),

		ActionStyle: lipgloss.NewStyle().
			Foreground(action),

		PipelineStyle: lipgloss.NewStyle().
			Foreground(pipeline),

		SkippedStyle: lipgloss.NewStyle().
			Foreground(skipped).
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
			Foreground(errorCol).
			Bold(true),

		MutedStyle: lipgloss.NewStyle().
			Foreground(muted).
			Faint(true),

		ProgressStyle: lipgloss.NewStyle().
			Foreground(progress),

		WorkerStyle: lipgloss.NewStyle().
			Foreground(worker).
			Bold(true),

		WorkerIdleStyle: lipgloss.NewStyle().
			Foreground(workerIdle).
			Faint(true),

		LaneHeaderStyle: lipgloss.NewStyle().
			Foreground(laneHeader).
			Bold(true),
	}, nil
}
