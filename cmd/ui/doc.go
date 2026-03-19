// Package ui defines the user interface abstraction for jay. All output
// to the terminal is routed through a UI implementation, chosen at startup
// via the --tui flag. This makes it trivial to swap in richer Charm-based
// renderers later without touching any command or traversal logic.
//
// The only implementation shipped initially is Linear, which writes plain
// text to stdout via fmt.Println. Future implementations (e.g. a Bubble Tea
// TUI) satisfy the same Manager interface and are selected by name.
package ui
