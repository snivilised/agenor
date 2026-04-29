// Package ui defines the user interface abstraction for jay. All output
// to the terminal is routed through a UI implementation, chosen at startup
// via the --tui flag. This makes it trivial to swap in richer Charm-based
// renderers later without touching any command or traversal logic.
package ui
