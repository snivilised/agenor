package enums

// EntryType used to enable selecting directory entry type.
type EntryType uint

const (
	// EntryTypeFolder
	//
	EntryTypeFolder EntryType = iota // folder-entry

	// EntryTypeFile
	//
	EntryTypeFile // file-entry
)
