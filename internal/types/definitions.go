package types

import (
	"context"
	"io/fs"

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/life"
	"github.com/snivilised/traverse/pref"
)

// 📦 pkg: types - defines internal types

type (
	// Link represents a single decorator in the chain
	Link interface {
		// Next invokes this decorator which returns true if
		// next link in the chain can be run or false to stop
		// execution of subsequent links.
		Next(servant core.Servant, inspection Inspection) (bool, error)

		// Role indicates the identity of the link
		Role() enums.Role
	}

	// GuardianSealer protects against invalid decorations. There can only
	// be 1 sealer (the master) and currently that only comes into play
	// for fastward resume. An ordinary filter is decorate-able, so it
	// can't be the sealer. It is not mandatory for a master to be registered.
	// When no master is registered, Benign will be used.
	GuardianSealer interface {
		Seal(link Link) error
		IsSealed(top Link) bool
	}

	// Guardian is the gateway to accessing the invocation chain.
	Guardian interface {
		Decorate(link Link) error
		Unwind(role enums.Role) error
	}

	// Arrangeable
	Arrangeable interface {
		Arrange(active, order []enums.Role)
	}

	// PluginInit
	PluginInit struct {
		O        *pref.Options
		Controls *life.Controls
	}

	// Mediator controls interactions between different entities of
	// of the navigator
	Mediator interface {
		Guardian
		Arrangeable
		Navigate(ctx context.Context) (*KernelResult, error)
		Resume(ctx context.Context, active *core.ActiveState) (*KernelResult, error)
		Spawn(ctx context.Context, active *core.ActiveState) (*KernelResult, error)
		Supervisor() *measure.Supervisor
	}

	// Resources are dependencies required for navigation
	Resources struct {
		FS         *core.Forest
		Supervisor *measure.Supervisor
		Binder     *opts.Binder
	}

	// Plugin used to define interaction with supplementary features
	Plugin interface {
		Register(kc KernelController) error
		Role() enums.Role
		Init(pi *PluginInit) error
	}

	// Restoration; tbd...
	Restoration interface {
		Inject(state *core.ActiveState)
	}

	// Facilities is the interface provided to plugins to enable them
	// to initialise successfully.
	Facilities interface {
		Restoration
		Metrics() *measure.Supervisor
	}

	// Ignition
	Ignition struct {
		Session core.Session
	}

	KernelNavigator interface {
		Navigate(ctx context.Context) (*KernelResult, error)
	}

	// KernelController
	KernelController interface {
		KernelNavigator
		Ignite(ignition *Ignition)
		Result(ctx context.Context, err error) *KernelResult
		Mediator() Mediator
		Resume(ctx context.Context, active *core.ActiveState) (*KernelResult, error)
		Conclude(result core.TraverseResult)
	}

	// Inspection
	Inspection interface {
		Current() *core.Node
		Contents() core.DirectoryContents
		Entries() []fs.DirEntry
		Sort(et enums.EntryType) []fs.DirEntry
		Pick(et enums.EntryType)
		AssignChildren(children []fs.DirEntry)
	}

	SaveState struct {
		Path string
	}

	// RestoreState
	RestoreState struct {
		Path   string
		FS     nef.ReadFileFS
		Resume enums.ResumeStrategy
	}

	// OptionHarvest
	OptionHarvest interface {
		Options() *pref.Options
		Binder() *opts.Binder
		Loaded() *opts.LoadInfo
	}
)

// KernelResult is the internal representation of core.TraverseResult
type KernelResult struct {
	session  core.Session
	reporter core.Reporter
	complete bool
	err      error
}

func NewResult(session core.Session,
	supervisor *measure.Supervisor,
	err error,
	complete bool,
) *KernelResult {
	return &KernelResult{
		session:  session,
		reporter: supervisor,
		err:      err,
		complete: complete,
	}
}

func NewFailed(err error) *KernelResult {
	return &KernelResult{
		err: err,
	}
}

func (r *KernelResult) IsComplete() bool {
	return r.complete
}

func (r *KernelResult) Session() core.Session {
	return r.session
}

func (r *KernelResult) Metrics() core.Reporter {
	return r.reporter
}

func (r *KernelResult) Error() error {
	return r.err
}

type (
	FilterChildren interface { // TODO: is this still needed?
		Matching(files []fs.DirEntry) []fs.DirEntry
	}

	FilterChildrenFunc func(files []fs.DirEntry) []fs.DirEntry
)

func (fn FilterChildrenFunc) Matching(files []fs.DirEntry) []fs.DirEntry {
	return fn(files)
}
