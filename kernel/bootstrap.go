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
// ==> handshake phase, allows other features to see how bootstrap has been established
// ==> use bus (maybe we create a local version of bus) ??
// (https://dev.to/mustafaturan/decoupled-package-communication-in-go-g39)
//
// ==> perhaps we have a new life cycle event "initialised" => indicates things like,
// * "l10n language used": this will confirm what language is in play (language.Tag,
// not string), useful when client
// * session type (standard or resume)
// * sampling active
// * requests an unsupported language, and informs the lang defaulted to
// these are the features that need to be considered:
//
// i18n: defined by the client
// log setup: defined by the client
// options: primary or resume
// hibernation
// filter
// sampling
// options
// resume
//
// * remember, hooks are different to events and there is a distinction
// between life-cycle events, which only the kernel can emit and the events
// that can also be emitted by the hooks or features(plugins)
//
// to make this a bit clearer, perhaps we define a bus and comms to and
// from it use messages instead of events. That allows us to have a clearer
// vocabulary within traverse. Messages are dynamic and can be sent to a topic.
// Whereas events are fixed in nature. We have a clear definition of each message,
// which is bound to a topic.
//
// * chain of responsibility? built into messaging structure; ie a message maybe consumed
// or ignored (return a bool from the function to determine if other listeners are
// able to see the message)
//
// if you build your own bus, then is no need to use an asynchronous model
//
// the term "callback" always refers to the client navigation callback function
//
// * we really ought to strive building this as a series of layers. The only
// problem with the features(plugins) and the core functionality in kernel
// is that it seems that they depend on each other. It's gonna be a fight
// to prevent this from happening.
// 		To achieve this, we must place common definitions in core and ensure
// that any feature does not depend on the kernel. Each feature must subscribe
// to topics they're interested in, and emit a message to indicate completion.

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
