package enums

//go:generate stringer -type=EntryType -linecomment -trimprefix=EntryType -output entry-type-en-auto.go

// EntryType used to enable selecting directory entry type.
type EntryType uint

const (
	// EntryTypeDirectory represents a directory entry.
	//
	EntryTypeDirectory EntryType = iota // directory-entry

	// EntryTypeFile represents a file entry.
	//
	EntryTypeFile // file-entry

	// EntryTypeAll represents both directory and file entries.
	//
	EntryTypeAll // all-entries
)
