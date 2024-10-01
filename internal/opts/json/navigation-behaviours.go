package json

type (
	SubPathBehaviour struct {
		KeepTrailingSep bool
	}

	SortBehaviour struct {
		// case sensitive traversal order
		//
		IsCaseSensitive bool

		// SortFilesFirst defines whether a folder's files or directories
		// should be navigated first.
		//
		SortFilesFirst bool
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
