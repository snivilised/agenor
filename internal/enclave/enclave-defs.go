package enclave

import (
	"context"
	"io/fs"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/life"
	"github.com/snivilised/agenor/pref"
	nef "github.com/snivilised/nefilim"
)

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
		// Seal registers a link as the master sealer. Once sealed, only the master
		// sealer can decorate the guardian. This is to protect against invalid decorations.
		// For example, if a filter was decorated as the master sealer, then it would not
		// be possible to decorate a master sealer later on, which would prevent
		// fastward resume from being used.
		Seal(link Link) error

		// IsSealed checks if the provided link is the master sealer. This is used to determine
		// if a link can be decorated as the master sealer or if it is an ordinary filter.
		IsSealed(top Link) bool
	}

	// Swapper allows the underlying client handler to be decorated
	Swapper interface {
		// Swap replaces the underlying client handler with the provided decorator. This
		// is used to allow the guardian to be decorated with different handlers, such as
		// filters or a master sealer.
		Swap(decorator core.Client)
	}

	// Guardian is the gateway to accessing the invocation chain.
	Guardian interface {
		Swapper
		// Decorate adds a decorator to the invocation chain. The order of decoration is
		// important, as it determines the order in which the decorators are invoked.
		// For example, if a filter is decorated before a master sealer, then the filter
		// will be invoked before the master sealer, which would prevent fastward resume
		// from being used.
		Decorate(link Link) error

		// Unwind removes the most recently added decorator from the invocation chain. This is
		// used to allow decorators to be removed from the chain, such as when a filter is no
		// longer needed.
		Unwind(role enums.Role) error
	}

	// Arrangeable allows the active and order to be arranged for a role. This is used to
	// allow the guardian to arrange the active and order for a role, which is necessary
	// for fastward resume. When resuming, the active and order for the master sealer role
	// need to be arranged in a specific way to ensure that the invocation chain is
	// correctly reconstructed.
	Arrangeable interface {
		// Arrange allows the guardian to arrange the active and order for a role. This is used to
		// allow the guardian to arrange the active and order for a role, which is necessary
		// for fastward resume. When resuming, the active and order for the master sealer role
		// need to be arranged in a specific way to ensure that the invocation chain is
		// correctly reconstructed.
		Arrange(active, order []enums.Role)
	}

	// KernelNavigator defines the method required to navigate the kernel. This is used by the
	// guardian to navigate the kernel when a new session is started or when a session is resumed.
	KernelNavigator interface {
		// Navigate initiates navigation of the kernel. This is used by the guardian to navigate
		// the kernel when a new session is started or when a session is resumed. The context is
		// used to allow for cancellation of the navigation process, such as when a session is
		// closed while navigation is still in progress.
		Navigate(ctx context.Context) (*KernelResult, error)
	}

	// KernelController defines the methods required to control the kernel. This is used by the
	// guardian to control the kernel when a new session is started or when a session is resumed.
	// The guardian needs to be able to ignite the kernel, get the result of the navigation process,
	// snooze the kernel, and say bye to the kernel when the session is closed.
	KernelController interface {
		KernelNavigator
		// Ignite initiates the kernel with the provided ignition. This is used by the guardian
		// to ignite the kernel when a new session is started or when a session is resumed.
		// The ignition contains the session that is being started or resumed, which allows
		// the kernel to access the session and its properties during navigation.
		Ignite(ignition *Ignition)

		// Result retrieves the result of the navigation process. This is used by the
		// guardian to get the result of the navigation process, which is necessary for
		// determining the outcome of the navigation and for performing any necessary cleanup
		// or follow-up actions based on the result.
		Result(ctx context.Context) *KernelResult

		// Snooze allows the kernel to be snoozed for a certain duration.
		Snooze(ctx context.Context, active *core.ActiveState) (*KernelResult, error)

		// Bye allows the kernel to be gracefully shut down when a session is closed. This
		// is used by the guardian to say bye to the kernel when a session is closed, which
		// allows the kernel to perform any necessary cleanup and to ensure that all resources
		// are properly released.
		Bye(result core.TraverseResult)
	}

	// PluginInit defines the properties required to initialize a plugin. This is used
	// by plugins to access the resources and options provided by the guardian during
	// initialization. The plugin can use these properties to set up its internal state
	// and to prepare for handling its role in the invocation chain.
	PluginInit struct {
		// Options provides access to the options defined for the current session. This
		// is used by plugins to access the options defined for the current session, which
		// can be used to configure the plugin's behavior and to allow the plugin to respond
		// to user-defined options.
		O *pref.Options

		// Kontroller provides access to the kernel controller, which allows the plugin to
		// control the kernel during navigation. This is used by plugins to control the kernel
		// during navigation, which can be used to perform actions such as snoozing the kernel
		// or saying bye to the kernel based on certain conditions or events that occur during
		// navigation.
		Kontroller KernelController

		// Controls provides access to the life controls, which allows the plugin to control
		// the flow of the invocation chain. This is used by plugins to control the flow of the
		// invocation chain, which can be used to perform actions such as muting or unmuting
		// the invocation chain based on certain conditions or events that occur during navigation.
		Controls *life.Controls

		// Resources provides access to the resources defined for the current session. This is
		// used by plugins to access the resources defined for the current session, which can
		// be used to perform actions such as accessing the file system or interacting with the
		// supervisor during navigation.
		Resources *Resources
	}

	// Mediator controls interactions between different entities of
	// of the navigator
	Mediator interface {
		Guardian
		Arrangeable
		KernelController
		// Read allows the mediator to read the contents of a directory at the
		// specified path. This is used by the guardian to read the contents of a
		// directory during navigation, which is necessary for determining the next
		// steps in the navigation process based on the structure of the file system.
		Read(path string) ([]fs.DirEntry, error)

		// Spawn allows the mediator to spawn a new child navigation with the specified
		// tree. This is used by the guardian to spawn a new child navigation when a
		// new session is started or when a session is resumed, which allows the kernel
		// to navigate the specified tree and to perform the necessary actions based on
		// the structure of the file system and the options defined for the session.
		Spawn(ctx context.Context, tree string) (*KernelResult, error)

		// Bridge combines information gleaned from the previous traversal that was
		// interrupted, into the resume traversal. This is used by the guardian to bridge
		// the information from the previous traversal into the resume traversal, which
		// allows the kernel to reconstruct the invocation chain and to continue navigation
		// from where it left off when the session was interrupted.
		Bridge(active *core.ActiveState)

		// Supervisor provides access to the supervisor, which allows the guardian to
		// interact with the supervisor during navigation.
		Supervisor() *Supervisor
	}

	// Resources are dependencies required for navigation
	Resources struct {
		// Forest provides access to the forest, which allows plugins and
		// decorators to interact with the file system during navigation. This
		// is used by plugins and decorators to interact with the file system during
		// navigation, which can be used to perform actions such as reading files,
		// checking for the existence of files or directories, or accessing file
		// metadata during navigation.
		Forest *core.Forest

		// Supervisor provides access to the supervisor, which allows plugins
		// and decorators to interact with the supervisor during navigation.
		Supervisor *Supervisor

		// Binder provides access to the options binder, which allows plugins
		// and decorators to bind options during navigation. This is used by
		// plugins and decorators to bind options during navigation, which can
		// be used to allow users to define custom options that can be accessed
		// and used by plugins and decorators during navigation.
		Binder *opts.Binder
	}

	// Plugin used to define interaction with supplementary features
	Plugin interface {
		// Register allows the plugin to register itself with the kernel
		// controller, which is necessary for the plugin to be able to control
		// the kernel during navigation. This is used by plugins to register
		// themselves with the kernel controller, which allows the plugin to
		// perform actions such as snoozing the kernel or saying bye to the
		// kernel based on certain conditions or events that occur during navigation.
		Register(kc KernelController) error

		// Role indicates the identity of the plugin, which is used to determine
		// the order of decoration in the invocation chain. This is used by plugins
		// to indicate their role, which is used to determine the order of decoration
		// in the invocation chain and to ensure that the invocation chain is correctly
		// constructed based on the roles of the plugins and decorators that are
		// registered.
		Role() enums.Role

		// Init initializes the plugin with the given initialization parameters.
		// This is used by plugins to perform any necessary setup or configuration
		// before they can be used during navigation.
		Init(pi *PluginInit) error
	}

	// Restoration defines the method required to inject a previously loaded active
	// state into the kernel. This is used by the guardian to inject a previously
	// loaded active state into the kernel during a resume, which allows the kernel
	// to reconstruct the invocation chain and to continue navigation from where it
	// left off when the session was interrupted.
	Restoration interface {
		// Inject allows the guardian to inject a previously loaded active state into
		// the kernel during a resume, which allows the kernel to reconstruct the invocation
		// chain and to continue navigation from where it left off when the session was
		// interrupted. The active state contains information about the previous traversal,
		// such as the current path, the depth in the file system, and any metrics that were
		// collected during the previous traversal, which can be used by the kernel to
		// reconstruct the invocation chain and to continue navigation based on the state of
		// the previous traversal.
		Inject(state *core.ActiveState)
	}

	// Ignition defines the properties required to ignite the kernel. This is used by the
	// guardian to ignite the kernel when a new session is started or when a session is
	// resumed, which allows the kernel to access the session and its properties during
	// navigation. The session contains information about the user, the options defined
	// for the session, and any resources that are available during navigation, which
	// can be used by the kernel to perform actions based on the session's properties
	// during navigation.
	Ignition struct {
		// Session provides access to the session, which allows the kernel to access the session
		// and its properties during navigation.
		Session core.Session
	}

	// Inspection defines the methods required to inspect the current state of the
	// traversal. This is used by decorators to inspect the current state of the
	// traversal, which can be used to make decisions about how to control the flow of
	// the invocation chain based on the current state of the traversal. The inspection
	// provides methods for accessing the current node, the contents of the current
	// directory, and for picking specific entries from the directory contents based
	// on their type.
	Inspection interface {
		// Current provides access to the current node being traversed.
		Current() *core.Node

		// Contents provides access to the contents of the current directory.
		Contents() core.DirectoryContents

		// Entries provides access to the entries of the current directory.
		Entries() []fs.DirEntry

		// Sort allows the entries of the current directory to be sorted based on their type.
		Sort(et enums.EntryType) []fs.DirEntry

		// Pick allows specific entries from the directory contents to be picked based on their type.
		Pick(et enums.EntryType)

		// AssignChildren allows the entries of the current directory to be assigned as children of
		// the current node.
		AssignChildren(children []fs.DirEntry)
	}

	// RestoreState defines properties required in order to instigate
	// a resume.
	RestoreState struct {
		// Path provides the path to the directory that should be navigated to during the resume.
		Path string

		// FS provides access to the file system, which allows the guardian to access the file
		// system during the resume.
		FS nef.ReadFileFS

		// Strategy indicates the resume strategy that should be used during the resume.
		Strategy enums.ResumeStrategy
	}

	// StateHandler defines a method that allows am internal client to modify active
	// state after it has been marshalled in. This is meant for use by unit tests to
	// easy the process of defining their constraints. Without this, there would be
	// a need to provide separate json files for each test.
	StateHandler interface {
		// OnLoad allows an internal client to modify the active state after it has been
		// marshalled in. This is meant for use by unit tests to easy the process of defining
		// their constraints. Without this, there would be a need to provide separate
		// json files for each test.
		OnLoad(active *core.ActiveState)
	}

	// OptionHarvest provides access to the options and binder during testing. This is
	// used by unit tests to access the options and binder during testing, which allows
	// the unit tests to verify that the options and binder are being used correctly
	// during navigation. This is necessary for ensuring that the plugins and decorators
	// are able to access and use the options and binder as intended during navigation.
	OptionHarvest interface {
		// Options provides access to the options defined for the current session.
		Options() *pref.Options

		// Binder provides access to the options binder.
		Binder() *opts.Binder

		// Loaded provides access to the loaded active state.
		Loaded() *opts.LoadInfo
	}
)

type (
	// Loader is to be defined by a unit test and should modify the loaded active state
	// for the test's own purposes. This allows the unit tests to be isolated from the
	// content of the loaded active state and frees the unit tests from having to provide
	// a separate json file loaded specifically for it.
	Loader func(active *core.ActiveState)
)

// OnLoad allows an internal client to modify the active state after it has been
// marshalled in. This is meant for use by unit tests to easy the process of defining
// their constraints. Without this, there would be a need to provide separate
// json files for each test.
func (fn Loader) OnLoad(active *core.ActiveState) {
	fn(active)
}
