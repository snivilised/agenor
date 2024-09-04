package pref

import (
	"github.com/snivilised/traverse/enums"
)

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

		// DirectoryEntryOrder defines whether a folder's files or directories
		// should be navigated first.
		//
		DirectoryEntryOrder enums.DirectoryContentsOrder
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

		// Sort, behaviours relating to sorting of a folder's directory entries.
		//
		Sort SortBehaviour

		// Cascade controls how deep to navigate
		//
		Cascade CascadeBehaviour
	}
)

func WithNavigationBehaviours(nb *NavigationBehaviours) Option {
	return func(o *Options) error {
		o.Behaviours = *nb

		return nil
	}
}

func WithSubPathBehaviour(sb *SubPathBehaviour) Option {
	return func(o *Options) error {
		o.Behaviours.SubPath = *sb

		return nil
	}
}

func WithSortBehaviour(sb *SortBehaviour) Option {
	return func(o *Options) error {
		o.Behaviours.Sort = *sb

		return nil
	}
}

func WithDepth(depth uint) Option {
	return func(o *Options) error {
		o.Behaviours.Cascade.Depth = depth

		return nil
	}
}

func WithNoRecurse() Option {
	return func(o *Options) error {
		o.Behaviours.Cascade.NoRecurse = true

		return nil
	}
}
