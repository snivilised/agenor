package enums

//go:generate stringer -type=DirectoryContentsOrder -linecomment -trimprefix=DirectoryContentsOrder -output directory-contents-en-auto.go

// DirectoryContentsOrder determines what order a directories
// entries are invoked for.
type DirectoryContentsOrder uint

const (
	// DirectoryContentsOrderFoldersFirst invoke folders first
	//
	DirectoryContentsOrderFoldersFirst DirectoryContentsOrder = iota // folders-first

	// DirectoryContentsOrderFilesFirst invoke files first
	//
	DirectoryContentsOrderFilesFirst // files-first
)
