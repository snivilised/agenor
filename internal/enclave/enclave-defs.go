package enclave

import (
	"context"
	"io/fs"

	nef "github.com/snivilised/nefilim"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/life"
	"github.com/snivilised/traverse/pref"
)

// 📦 pkg: enclave - defines internal types

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
		O          *pref.Options
		Kontroller KernelController
		Controls   *life.Controls
		Resources  *Resources
	}

	// Mediator controls interactions between different entities of
	// of the navigator
	Mediator interface {
		Guardian
		Arrangeable
		Navigate(ctx context.Context) (*KernelResult, error)
		Read(path string) ([]fs.DirEntry, error)
		Resume(ctx context.Context, active *core.ActiveState) (*KernelResult, error)
		Spawn(ctx context.Context, active *core.ActiveState) (*KernelResult, error)
		Bridge(tree, current string)
		Supervisor() *core.Supervisor
	}

	// Resources are dependencies required for navigation
	Resources struct {
		Forest     *core.Forest
		Supervisor *core.Supervisor
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

	// RestoreState defines properties required in order to instigate
	// a resume.
	RestoreState struct {
		Path   string
		FS     nef.ReadFileFS
		Resume enums.ResumeStrategy
	}

	// StateHandler defines a method that allows am internal client to modify active
	// state after it has been marshalled in. This is meant for use by unit tests to
	// easy the process of defining their constraints. Without this, there would be
	// a need to provide separate json files for each test.
	StateHandler interface {
		OnLoad(active *core.ActiveState)
	}

	// Loader is to be defined by a unit test and should modify the loaded active state
	// for the test's own purposes. This allows the unit tests to be isolated from the
	// content of the loaded active state.
	Loader func(active *core.ActiveState)

	// OptionHarvest
	OptionHarvest interface {
		Options() *pref.Options
		Binder() *opts.Binder
		Loaded() *opts.LoadInfo
	}
)

func (fn Loader) OnLoad(active *core.ActiveState) {
	fn(active)
}
