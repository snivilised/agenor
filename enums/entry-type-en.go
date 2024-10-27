package enums

//go:generate stringer -type=EntryType -linecomment -trimprefix=EntryType -output entry-type-en-auto.go

// EntryType used to enable selecting directory entry type.
type EntryType uint

const (
	// EntryTypeDirectory
	//
	EntryTypeDirectory EntryType = iota // directory-entry

	// EntryTypeFile
	//
	EntryTypeFile // file-entry

	// EntryTypeAll
	//
	EntryTypeAll // all-entries
)
