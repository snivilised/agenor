package core

import (
	"io/fs"
	"path/filepath"
	"time"

	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/tfs"
)

// 📦 pkg: core - contains universal definitions and handles user facing cross
// cutting concerns.

type (
	// ResultCompletion used to determine if the result really represents
	// final navigation completion.
	ResultCompletion interface {
		IsComplete() bool
	}

	// Completion
	Completion func() bool

	// Session represents a traversal session and keeps tracks of
	// timing.
	Session interface {
		ResultCompletion
		StartedAt() time.Time
		Elapsed() time.Duration
	}

	// TraverseResult
	TraverseResult interface {
		Metrics() Reporter
		Session() Session
	}

	// Servant provides the client with facility to request properties
	// about the current navigation node.
	Servant interface {
		Node() *Node
	}

	// Forest contains the logical file systems required
	// for navigation.
	Forest struct {
		// T is the file system that contains just the functionality required
		// for traversal. It can also represent other file systems including afero,
		// providing the appropriate adapters are in place.
		T tfs.TraversalFS

		// R is the file system required for resume operations, ie we load
		// and save resume state via this file system instance, which is
		// distinct from the traversal file system.
		R tfs.TraversalFS
	}

	// Client is the callback invoked for each file system node found
	// during traversal.
	Client func(servant Servant) error

	// FsDescription description of a file system
	FsDescription struct {
		IsRelative bool
	}

	Permissions struct {
		File fs.FileMode
		Dir  fs.FileMode
	}

	ActiveState struct {
		Tree                string
		TraverseDescription FsDescription
		ResumeDescription   FsDescription
		Subscription        enums.Subscription
		Hibernation         enums.Hibernation
		CurrentPath         string
		IsDir               bool
		Depth               int
		Metrics             Metrics
	}

	// TimeFunc get time
	TimeFunc func() time.Time

	// SimpleHandler is a function that takes no parameters and can
	// be used by any notification with this signature.
	SimpleHandler func()

	// BeginHandler invoked before traversal begins
	BeginHandler func(tree string)

	// EndHandler invoked at the end of traversal
	EndHandler func(result TraverseResult)

	// HibernateHandler is a generic handler that is used by hibernation
	// to indicate wake or sleep.
	HibernateHandler func(description string)
)

func (fn Completion) IsComplete() bool {
	return fn()
}

func (s *ActiveState) Clone() *ActiveState {
	c := *s
	return &c
}

const (
	FileSystemTimeFormat = "2006-01-02_15-04-05"
	PackageName          = "agenor"
	filePerm             = 0o644
	dirPerm              = 0o755
)

var (
	Now        TimeFunc
	Perms      Permissions
	ResumeTail string
)

func init() {
	Now = time.Now
	Perms = Permissions{
		File: filePerm,
		Dir:  dirPerm,
	}
	ResumeTail = filepath.Join(PackageName, "admin", "resume")
}
