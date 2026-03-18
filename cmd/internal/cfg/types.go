package cfg

import (
	"time"
)

// ---------------------------------------------------------------------------
// Mapped sections - decoded via mapstructure
// ---------------------------------------------------------------------------

// TUIConfig holds settings for the terminal user interface.
type TUIConfig struct {
	PerItemDelay time.Duration `mapstructure:"per-item-delay"`
}

// InteractionConfig groups all user-interaction knobs.
type InteractionConfig struct {
	TUI TUIConfig `mapstructure:"tui"`
}

// ExtensionsConfig controls how file extensions are normalised.
type ExtensionsConfig struct {
	SuffixesCSV   string            `mapstructure:"suffixes-csv"`
	TransformsCSV string            `mapstructure:"transforms-csv"`
	Map           map[string]string `mapstructure:"map"`
}

// AdvancedConfig holds low-level behavioural switches.
type AdvancedConfig struct {
	AbortOnError         bool             `mapstructure:"abort-on-error"`
	OverwriteOnCollision bool             `mapstructure:"overwrite-on-collision"`
	Extensions           ExtensionsConfig `mapstructure:"extensions"`
}

// LoggingConfig controls the jay log file.
type LoggingConfig struct {
	LogPath      string `mapstructure:"log-path"`
	MaxSizeInMB  int    `mapstructure:"max-size-in-mb"`
	MaxBackups   int    `mapstructure:"max-backups"`
	MaxAgeInDays int    `mapstructure:"max-age-in-days"`
	Level        string `mapstructure:"level"`
	TimeFormat   string `mapstructure:"time-format"`
}

// MappedConfig bundles all sections that are decoded into concrete types.
type MappedConfig struct {
	Interaction InteractionConfig `mapstructure:"interaction"`
	Advanced    AdvancedConfig    `mapstructure:"advanced"`
	Logging     LoggingConfig     `mapstructure:"logging"`
}

// ---------------------------------------------------------------------------
// Raw sections - consumer-driven, arbitrary user content
// ---------------------------------------------------------------------------

// RawAction is one entry from the actions block.  The cmd and when strings
// are kept verbatim; jay's action-runner is responsible for interpreting them.
type RawAction struct {
	Cmd  string `mapstructure:"cmd"`
	When string `mapstructure:"when"`
}

// RawPipeline is one entry from the pipelines block.
type RawPipeline struct {
	Steps []string `mapstructure:"steps"`
}

// FlagShortOverride captures per-command short-flag remappings.
//
//	flags.short.overrides.cmds.<cmd>.<flag> = <letter>
type FlagShortOverride map[string]map[string]string

// FlagInvokeDefaults captures command-level flag defaults.
//
//	flags.invoke.cmds.<cmd>.<flag> = <value>
type FlagInvokeDefaults map[string]map[string]any

// FlagComponentDefaults captures component-level flag defaults.
//
//	flags.component.<component>.<flag> = <value>
type FlagComponentDefaults map[string]map[string]any

// FlagsConfig aggregates all flags sub-sections.
type FlagsConfig struct {
	Short     FlagShortOverride     `mapstructure:"short"`
	Invoke    FlagInvokeDefaults    `mapstructure:"invoke"`
	Component FlagComponentDefaults `mapstructure:"component"`
}

// RawConfig holds all unstructured sections verbatim.
type RawConfig struct {
	Actions   map[string]RawAction   `mapstructure:"actions"`
	Pipelines map[string]RawPipeline `mapstructure:"pipelines"`
	Flags     FlagsConfig            `mapstructure:"flags"`
}

// ---------------------------------------------------------------------------
// Top-level unified config
// ---------------------------------------------------------------------------

// Config is the root configuration object handed to callers after a
// successful Load + Validate cycle.
type Config struct {
	Mapped MappedConfig
	Raw    RawConfig
}
