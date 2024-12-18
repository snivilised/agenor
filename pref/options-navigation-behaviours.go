package pref

type (
	// SubPathBehaviour
	SubPathBehaviour struct {
		KeepTrailingSep bool
	}
	// SortBehaviour
	SortBehaviour struct {
		// case sensitive traversal order
		//
		IsCaseSensitive bool

		// SortFilesFirst defines whether a directory's files or directories
		// should be navigated first.
		//
		SortFilesFirst bool
	}

	// HibernationBehaviour
	HibernationBehaviour struct {
		InclusiveWake  bool
		InclusiveSleep bool
	}

	CascadeBehaviour struct {
		// Depth sets a maximum traversal depth
		//
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

	// NavigationBehaviours
	NavigationBehaviours struct {
		// SubPath, behaviours relating to handling of sub-path calculation
		//
		SubPath SubPathBehaviour

		// Sort, behaviours relating to sorting of a directory's entries.
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
