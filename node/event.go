package node

import (
	"io/fs"

	"github.com/snivilised/extendio/xfs/utils"
	"github.com/snivilised/traverse/enums"
)

// package: node represents a single file system entity and in the world
// of observables is like an event. The observable fluency chain will be
// based upon the node.

// Envelope is a wrapper around the Event
type Envelope struct {
	E *Event
}

// Seal wraps an Event inside an Envelope
func Seal(e *Event) *Envelope {
	// When using rx, we must instantiate it with a struct, not a pointer
	// to struct. So we create a distinction between what is the unit of
	// transfer (Envelope) and it's payload, the Event. This means that the
	// Envelope can only have non pointer receivers, but the Event can
	// be defined with pointer receivers.
	//
	return &Envelope{
		E: e,
	}
}

// Event represents a file system node event and represents each file system
// entity encountered during traversal. The event representing the root node
// does not have a DirEntry because it is not created as a result of a readDir
// invoke. Therefore, the client has to know that when its function is called back,
// they will be no DirEntry for the root node.
type Event struct {
	Path        string
	Entry       fs.DirEntry // contains a FileInfo via Info() function
	Info        fs.FileInfo // optional file info instance
	Extension   Extension   // extended information about the directory entry
	Error       error
	Children    []fs.DirEntry
	Parent      *Event
	filteredOut bool
	dir         bool
}

// Extension provides extended information if the client requests
// it by setting the DoExtend boolean in the traverse options.
type Extension struct {
	Depth   int               // traversal depth relative to the root
	IsLeaf  bool              // defines whether this node a leaf node
	Name    string            // derived as the leaf segment from filepath.Split
	Parent  string            // derived as the directory from filepath.Split
	SubPath string            // represents the path between the root and the current item
	Scope   enums.FilterScope // type of folder corresponding to the Filter Scope
	Custom  any               // to be set and used by the client
}

// New create a new node Event
func New(
	path string, entry fs.DirEntry, info fs.FileInfo, parent *Event, err error,
) *Event {
	event := &Event{
		Path:     path,
		Entry:    entry,
		Info:     info,
		Parent:   parent,
		Children: []fs.DirEntry{},
		Error:    err,
	}
	event.dir = isDir(event)

	return event
}

// Root creates a new node Event which represents the root of directory
// tree to traverse.
func Root() *Event {
	event := &Event{
		// TODO: complete
	}

	return event
}

// Clone makes shallow copy of Event (excluding the error).
func (e *Event) Clone() *Event {
	c := *e
	c.Error = nil

	return &c
}

// IsFolder indicates wether this event is a folder.
func (e *Event) IsFolder() bool {
	return e.dir
}

func (e *Event) key() string {
	// ti.Extension.SubPath
	return "ti.Extension.SubPath"
}

func isDir(e *Event) bool {
	if !utils.IsNil(e.Entry) {
		return e.Entry.IsDir()
	} else if !utils.IsNil(e.Info) {
		return e.Info.IsDir()
	}

	// only get here in error scenario, because neither Entry or Info is set
	//
	return false
}
