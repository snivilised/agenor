package command

import (
	"github.com/snivilised/mamba/store"
)

// ---------------------------------------------------------------------------
// Subscription flag values - what the user types on the command line
// ---------------------------------------------------------------------------

const (
	// SubscribeFlagFiles subscribes to file nodes only.
	SubscribeFlagFiles = "files"

	// SubscribeFlagDirs subscribes to directory nodes only.
	SubscribeFlagDirs = "dirs"

	// SubscribeFlagAll subscribes to all nodes (files and directories).
	SubscribeFlagAll = "all"

	// SubscribeFlagDefault is the default subscription if not specified.
	SubscribeFlagDefault = SubscribeFlagAll
)

// ---------------------------------------------------------------------------
// Root parameter set
// ---------------------------------------------------------------------------

// RootParameterSet holds flags defined on the root command that are
// inherited by all sub-commands via PersistentFlags.
type RootParameterSet struct {
	store.ParameterSetWithOverrides

	// Language sets the IETF BCP 47 language tag for i18n output.
	Language string

	// TUI selects the display mode. Corresponds to --tui(-t) <mode>.
	// Defaults to "linear". Future values: "porthole", "lanes".
	TUI string

	// Theme selects the colour theme. Corresponds to --theme <name>.
	// Defaults to "system" which uses ANSI-16 colours set by the
	// user's terminal theme. Any other value is resolved to a YAML
	// file in the themes directory (~/.config/jay/themes/<name>.yaml
	// or $JAY_THEMES_DIR/<name>.yaml).
	Theme string
}

// ---------------------------------------------------------------------------
// Walk parameter set
// ---------------------------------------------------------------------------

// WalkParameterSet holds the jay-specific flags for the walk command.
// Shared families (preview, cascade, sampling, poly-filter) are
// registered separately and are not embedded here.
type WalkParameterSet struct {
	store.ParameterSetWithOverrides

	// Subscribe controls which node types are visited.
	// Valid values: "files", "dirs", "all". Maps to --subscribe(-s).
	Subscribe string

	// Action names the config-defined action to invoke for each node.
	// Maps to --action(-a).
	Action string

	// Pipeline names the config-defined pipeline to execute.
	// Maps to --pipeline(-p).
	Pipeline string

	// Resume defines the resume strategy for interrupted traversals.
	// Maps to --resume(-r).
	// Valid values: "spawn", "fastward". Empty means prime (no resume).
	Resume string
}

// ---------------------------------------------------------------------------
// Run parameter set
// ---------------------------------------------------------------------------

// RunParameterSet holds the jay-specific flags for the run command.
// It mirrors WalkParameterSet; run additionally gets
// WorkerPoolParameterSet via a separate family registration.
type RunParameterSet struct {
	store.ParameterSetWithOverrides

	// Subscribe controls which node types are visited.
	// Valid values: "files", "dirs", "all". Maps to --subscribe(-s).
	Subscribe string

	// Action names the config-defined action to invoke for each node.
	// Maps to --action(-a).
	Action string

	// Pipeline names the config-defined pipeline to execute.
	// Maps to --pipeline(-p).
	Pipeline string

	// Resume defines the resume strategy for interrupted traversals.
	// Maps to --resume(-r).
	// Valid values: "spawn", "fastward". Empty means prime (no resume).
	Resume string
}
