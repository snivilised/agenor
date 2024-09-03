package filtering

import (
	"fmt"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
)

func createPolyFilter(polyDef *core.PolyFilterDef) (core.TraverseFilter, error) {
	// enforce the correct filter scopes
	//
	polyDef.File.Scope.Set(enums.ScopeFile)
	polyDef.File.Scope.Clear(enums.ScopeFolder)

	polyDef.Folder.Scope.Set(enums.ScopeFolder)
	polyDef.Folder.Scope.Clear(enums.ScopeFile)

	var (
		file, folder core.TraverseFilter
		err          error
	)

	if file, err = NewNodeFilter(&polyDef.File, nil); err != nil {
		return nil, err
	}

	if folder, err = NewNodeFilter(&polyDef.Folder, nil); err != nil {
		return nil, err
	}

	filter := &Poly{
		File:   file,
		Folder: folder,
	}

	return filter, nil
}

// Poly is a dual filter that allows files and folders to be filtered
// independently. The Folder filter only applies when the current node
// is a file. This is because, filtering doesn't affect navigation, it only
// controls wether the client callback is invoked or not. That is to say, if
// a particular folder fails to pass a filter, the callback will not be
// invoked for that folder, but we still descend into it and navigate its
// children. This is the reason why the poly filter is only active when the
// the current node is a filter as the client callback will only be invoked
// for the file if its parent folder passes the poly folder filter and
// the file passes the poly file filter.
type Poly struct {
	// File is the filter that applies to a file. Note that the client does
	// not have to set the File scope as this is enforced automatically as
	// well as ensuring that the Folder scope has not been set. The client is
	// still free to set other scopes.
	File core.TraverseFilter

	// Folder is the filter that applies to a folder. Note that the client does
	// not have to set the Folder scope as this is enforced automatically as
	// well as ensuring that the File scope has not been set. The client is
	// still free to set other scopes.
	Folder core.TraverseFilter
}

// Description
func (f *Poly) Description() string {
	return fmt.Sprintf("Poly - FILE: '%v', FOLDER: '%v'",
		f.File.Description(), f.Folder.Description(),
	)
}

// Validate ensures that both filters definition are valid, panics when invalid
func (f *Poly) Validate() error {
	if err := f.File.Validate(); err != nil {
		return err
	}

	return f.Folder.Validate()
}

// Source returns the Sources of both the File and Folder filters separated
// by a '##'
func (f *Poly) Source() string {
	return fmt.Sprintf("%v##%v",
		f.File.Source(), f.Folder.Source(),
	)
}

// IsMatch returns true if the current node is a file and both the current
// file matches the poly file filter and the file's parent folder matches
// the poly folder filter. Returns true of the current node is a folder.
func (f *Poly) IsMatch(node *core.Node) bool {
	if !node.IsFolder() {
		return f.Folder.IsMatch(node.Parent) && f.File.IsMatch(node)
	}

	return true
}

// IsApplicable returns the result of applying IsApplicable to
// the poly Filter filter if the current node is a file, returns false
// for folders.
func (f *Poly) IsApplicable(node *core.Node) bool {
	if !node.IsFolder() {
		return f.File.IsApplicable(node)
	}

	return false
}

// Scope is a bitwise OR combination of both filters
func (f *Poly) Scope() enums.FilterScope {
	return f.File.Scope() | f.Folder.Scope()
}
