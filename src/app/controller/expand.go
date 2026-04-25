package controller

import (
	"path/filepath"
	"strings"

	"github.com/snivilised/jaywalk/src/agenor/core"
)

// ExpandResult is the outcome of a placeholder expansion attempt.
type ExpandResult struct {
	// Cmd is the fully expanded command string, valid only when skipped is false.
	Cmd string

	// Skipped is true when a placeholder resolved to a path at or above root.
	Skipped bool

	// Placeholder is the token that caused the breach, e.g. "{{.grand}}".
	// Valid only when skipped is true.
	Placeholder string

	// resolvedPath is the path the offending placeholder resolved to.
	// Valid only when skipped is true.
	resolvedPath string
}

// expandData holds the resolved values for every placeholder.
// It is built once per node and reused across multiple expand calls.
type expandData struct {
	path   string
	name   string
	stem   string
	ext    string
	parent string
	grand  string
	great  string
	root   string
}

// buildExpandData resolves all placeholder values for the given node and
// root. Ancestor paths that would breach root are clamped internally for
// the purpose of breach detection only - they are never used in a cmd string.
func buildExpandData(root string, node *core.Node) expandData {
	cleanRoot := filepath.Clean(root)
	cleanPath := filepath.Clean(node.Path)

	name := filepath.Base(cleanPath)
	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	parent := filepath.Dir(cleanPath)
	grand := filepath.Dir(parent)
	great := filepath.Dir(grand)

	return expandData{
		path:   cleanPath,
		name:   name,
		stem:   stem,
		ext:    ext,
		parent: parent,
		grand:  grand,
		great:  great,
		root:   cleanRoot,
	}
}

// breaches reports whether the given path is at or above root.
// Both paths must already be cleaned.
func breaches(root, candidate string) bool {
	// candidate breaches root if it is not rooted within root,
	// i.e. it does not have root as a prefix followed by a separator
	// or it is exactly root itself going further up.
	rel, err := filepath.Rel(root, candidate)
	if err != nil {
		// If Rel fails the paths are on different volumes (Windows).
		// Treat as non-breach - the path cannot be above root.
		return false
	}

	// A relative path starting with ".." exits the root tree.
	return rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator))
}

// Expand substitutes all recognised placeholders in cmd with their
// resolved values for the given node and root. If any placeholder used
// in cmd would resolve to a path at or above root, Expand returns a
// skip result identifying the offending placeholder.
func Expand(cmd, root string, node *core.Node) ExpandResult {
	d := buildExpandData(root, node)

	// Check each ancestor placeholder that could breach root.
	// {{.path}}, {{.name}}, {{.stem}}, {{.ext}} cannot breach root by
	// definition - only ancestor-climbing placeholders can.
	ancestorChecks := []struct {
		placeholder string
		resolved    string
	}{
		{"{{.parent}}", d.parent},
		{"{{.grand}}", d.grand},
		{"{{.great}}", d.great},
	}

	for _, check := range ancestorChecks {
		if strings.Contains(cmd, check.placeholder) {
			if breaches(d.root, check.resolved) {
				return ExpandResult{
					Skipped:      true,
					Placeholder:  check.placeholder,
					resolvedPath: check.resolved,
				}
			}
		}
	}

	replacer := strings.NewReplacer(
		"{{.path}}", d.path,
		"{{.name}}", d.name,
		"{{.stem}}", d.stem,
		"{{.ext}}", d.ext,
		"{{.parent}}", d.parent,
		"{{.grand}}", d.grand,
		"{{.great}}", d.great,
		"{{.root}}", d.root,
	)

	return ExpandResult{
		Cmd: replacer.Replace(cmd),
	}
}
