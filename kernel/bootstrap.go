package kernel

// BootstrapPluginName is the name of a plugin that needs to hook
// the bootstrap process
type BootstrapPluginName = string

// BootStrapHook provides access to an individual bootstrap hook
type BootStrapHook interface {
	Tap(name string)
}

// BootStrapHooks provides access to all available bootstrap hooks
type BootStrapHooks interface {
}

type BootStrapper interface {
	Hooks() BootStrapHooks
}

// bootstrap is the coordinator of initialisation events

type bootstrap struct {
	// should allow other entities to customise the boot phase, eg
	// if there is a resumer in place, it should change the boot
	// process to read the options from file, rather than getting
	// the options.
	//
}

// bs.on("some-event", func() { the "some-event", will be an enum, not string
//		do something
// })

// could use a similar model for life-cycle events

// apply plugin; each plugin has an apply method. when Tap is invoked,
// apply is called on the plugin.

// tap?

/* example:
compilation.hooks.runtimeRequirementInTree
	.for(RuntimeGlobals.systemContext)
	.tap("RuntimePlugin", chunk => {
		const { outputOptions } = compilation;
		const { library: globalLibrary } = outputOptions;
		const entryOptions = chunk.getEntryOptions();
		const libraryType =
			entryOptions && entryOptions.library !== undefined
				? entryOptions.library.type
				: globalLibrary.type;

		if (libraryType === "system") {
			compilation.addRuntimeModule(
				chunk,
				new SystemContextRuntimeModule()
			);
		}
		return true;
	});


	in this context, there are plugins of different types; this one
	is a compilation plugin: compilation.hooks

	so we could define:
	- navigation.hooks
	- bootstrap.hooks


	a plugin can be created with a client provided func handler

	I think that tap is like Register


	- so an internal entity, eg the navigator (nav-ctrl), exposes a set of hooks
	- a plugin registers to tap this hooks. to achieve this, the plugin defines an
	apply method, which is invoked by the source entry. the source entity provides an
	interface required by the plugin in order for them to implement apply.
	- hooks are defined for the life-cycle events
*/
