package filtering

import (
	"fmt"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
)

// Poly is a dual filter that allows files and directories to be filtered
// independently. The directory filter only applies when the current node
// is a file. This is because, filtering doesn't affect navigation, it only
// controls wether the client callback is invoked or not. That is to say, if
// a particular directory fails to pass a filter, the callback will not be
// invoked for that directory, but we still descend into it and navigate its
// children. This is the reason why the poly filter is only active when the
// the current node is a filter as the client callback will only be invoked
// for the file if its parent directory passes the poly directory filter and
// the file passes the poly file filter.
type Poly struct {
	// File is the filter that applies to a file. Note that the client does
	// not have to set the File scope as this is enforced automatically as
	// well as ensuring that the Directory scope has not been set. The client is
	// still free to set other scopes.
	File core.TraverseFilter

	// Directory is the filter that applies to a directory. Note that the client does
	// not have to set the Directory scope as this is enforced automatically as
	// well as ensuring that the File scope has not been set. The client is
	// still free to set other scopes.
	Directory core.TraverseFilter
}

// Description
func (f *Poly) Description() string {
	return fmt.Sprintf("Poly - FILE: '%v', DIRECTORY: '%v'",
		f.File.Description(), f.Directory.Description(),
	)
}

// Validate ensures that both filters definition are valid, panics when invalid
func (f *Poly) Validate() error {
	if err := f.File.Validate(); err != nil {
		return err
	}

	return f.Directory.Validate()
}

// Source returns the Sources of both the File and directory filters separated
// by a '##'
func (f *Poly) Source() string {
	return fmt.Sprintf("%v##%v",
		f.File.Source(), f.Directory.Source(),
	)
}

// IsMatch returns true if the current node is a file and both the current
// file matches the poly file filter and the file's parent directory matches
// the poly directory filter. Returns true of the current node is a directory.
func (f *Poly) IsMatch(node *core.Node) bool {
	if !node.IsDirectory() {
		return f.Directory.IsMatch(node.Parent) && f.File.IsMatch(node)
	}

	return true
}

// IsApplicable returns the result of applying IsApplicable to
// the poly Filter filter if the current node is a file, returns false
// for directories.
func (f *Poly) IsApplicable(node *core.Node) bool {
	if !node.IsDirectory() {
		return f.File.IsApplicable(node)
	}

	return false
}

// Scope is a bitwise OR combination of both filters
func (f *Poly) Scope() enums.FilterScope {
	return f.File.Scope() | f.Directory.Scope()
}
