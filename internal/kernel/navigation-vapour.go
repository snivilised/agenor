package kernel

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
)

type (
	// navigationVapour represents short-lived navigation data whose state relates
	// only to the current Node. (equivalent to inspection in extendio)
	navigationVapour struct { // after content has been read
		ns      *navigationStatic
		present *core.Node
		cargo   *Contents
		ents    []fs.DirEntry
	}
)

func (v *navigationVapour) static() *navigationStatic {
	return v.ns
}

func (v *navigationVapour) Current() *core.Node {
	return v.present
}

func (v *navigationVapour) Contents() core.DirectoryContents {
	return v.cargo
}

func (v *navigationVapour) Entries() []fs.DirEntry {
	return v.ents
}

func (v *navigationVapour) clear() {
	if v.cargo != nil {
		v.cargo.clear()
	} else {
		newEmptyContents()
	}
}

func (v *navigationVapour) Sort(et enums.EntryType) []fs.DirEntry {
	v.cargo.Sort(et)

	// change SortHook to return entries so we don't have to do this switch?
	switch et {
	case enums.EntryTypeAll:
		return v.cargo.All()
	case enums.EntryTypeFolder:
		return v.cargo.folders
	case enums.EntryTypeFile:
		return v.cargo.files
	}

	return nil
}

func (v *navigationVapour) Pick(et enums.EntryType) {
	switch et {
	case enums.EntryTypeAll:
		v.ents = v.cargo.All()
	case enums.EntryTypeFolder:
		v.ents = v.cargo.folders
	case enums.EntryTypeFile:
		v.ents = v.cargo.files
	}
}

func (v *navigationVapour) AssignChildren(children []fs.DirEntry) {
	v.present.Children = children
}
