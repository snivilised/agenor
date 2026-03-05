package pref

type (
	// SubPathBehaviour behaviours relating to handling of sub-path calculation
	SubPathBehaviour struct {
		// KeepTrailingSep defines whether to keep the trailing separator in
		// the sub-path. By default, the trailing separator is removed from
		// the sub-path. This option can be useful in scenarios where the presence
		// of a trailing separator is significant, such as when distinguishing between
		// directories and files. When KeepTrailingSep is set to true, the sub-path
		// will retain the trailing separator if it was present in the original path.
		// When set to false (the default), the trailing separator will be removed
		// from the sub-path.
		//
		// Example:
		// - Original path: "/path/to/directory/"
		//   - Sub-path with KeepTrailingSep = true: "directory/"
		//   - Sub-path with KeepTrailingSep = false: "directory"
		KeepTrailingSep bool
	}

	// SortBehaviour behaviours relating to sorting of a directory's entries.
	SortBehaviour struct {
		// IsCaseSensitive defines whether the traversal order is case sensitive.
		//
		IsCaseSensitive bool

		// SortFilesFirst defines whether a directory's files or directories
		// should be navigated first.
		//
		SortFilesFirst bool
	}

	// HibernationBehaviour defines behaviours relating to hibernation of the navigator.
	HibernationBehaviour struct {
		// InclusiveWake when wake occurs, permit client callback to
		// be invoked for the current node. Inclusive, true by default
		InclusiveWake bool

		// InclusiveSleep when sleep occurs, permit client callback to
		// be invoked for the current node. Exclusive, false by default.
		InclusiveSleep bool
	}

	// CascadeBehaviour behaviours relating to how deep to navigate
	CascadeBehaviour struct {
		// Depth sets a maximum traversal depth
		Depth uint

		// NoRecurse is an alternative to using Depth, but limits the traversal
		// to just the path specified by the user. Since the raison d'etre
		// of the navigator is to recursively process a directory tree, using
		// NoRecurse would appear to be contrary to its natural behaviour. However
		// there are clear usage scenarios where a client needs to process
		// only the files in a specified directory.
		//
		NoRecurse bool
	}

	// NavigationBehaviours defines all navigation behaviours for the navigator.
	NavigationBehaviours struct {
		// SubPath, behaviours relating to handling of sub-path calculation
		//
		SubPath SubPathBehaviour

		// Sort behaviours relating to sorting of a directory's entries.
		//
		Sort SortBehaviour

		// Cascade controls how deep to navigate
		//
		Cascade CascadeBehaviour
	}
)

// WithNavigationBehaviours defines all navigation behaviours
func WithNavigationBehaviours(nb *NavigationBehaviours) Option {
	return func(o *Options) error {
		o.Behaviours = *nb

		return nil
	}
}

// WithSubPathBehaviour defines all sub-path behaviours.
func WithSubPathBehaviour(sb *SubPathBehaviour) Option {
	return func(o *Options) error {
		o.Behaviours.SubPath = *sb

		return nil
	}
}

// WithSortBehaviour enabling setting of all sorting behaviours.
func WithSortBehaviour(sb *SortBehaviour) Option {
	return func(o *Options) error {
		o.Behaviours.Sort = *sb

		return nil
	}
}

// WithDepth sets the maximum number of directories deep the navigator
// will traverse to.
func WithDepth(depth uint) Option {
	return func(o *Options) error {
		o.Behaviours.Cascade.Depth = depth

		return nil
	}
}

// WithNoRecurse sets the navigator to not descend sub-directories.
func WithNoRecurse() Option {
	return func(o *Options) error {
		o.Behaviours.Cascade.NoRecurse = true

		return nil
	}
}
