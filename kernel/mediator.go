package kernel

type Mediator interface {
	Register() error
}

// mediator controls traversal events
type mediator struct {
	// there should be a registration phase; but doing so mean that
	// these entities should already exist, which is counter productive.
	// possibly use dependency inject where entities declare their
	// dependencies so that it is easier to orchestrate the boot.
	//
	// Are hooks plugins? (
	// query-status: lstat on the root directory; singular; WithQueryStatus
	// read-dir: ; singular; WithReader
	// folder-sub-path: ; singular; WithFolderSubPath
	// filter-int: ; singular; WithFilter
	// sort: ; singular
	// )
	//
	// Hooks mark specific points in navigation that can be customised
}

// application phases (should we define a state machine?)
//
// --> configuration: OnConfigured
// * ...
// --> i18n: Oni18n
// * ...
// --> log
// * ...
// --> session
// *
// --> get-options (via WithOptions for primary session, or restore with resume session)
