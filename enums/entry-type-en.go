package enums

//go:generate stringer -type=EntryType -linecomment -trimprefix=EntryType -output entry-type-en-auto.go

// EntryType used to enable selecting directory entry type.
type EntryType uint

const (
	// EntryTypeFolder
	//
	EntryTypeFolder EntryType = iota // folder-entry

	// EntryTypeFile
	//
	EntryTypeFile // file-entry

	// EntryTypeAll
	//
	EntryTypeAll // all-entries
)
