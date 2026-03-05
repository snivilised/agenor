package enums

//go:generate stringer -type=PersistenceFormat -linecomment -trimprefix=Persist -output persistence-format-en-auto.go

// PersistenceFormat represents the persistence formats available
type PersistenceFormat uint

const (
	// PersistUndefined undefined
	PersistUndefined PersistenceFormat = iota // persistence-undefined

	// PersistJSON persist in JSON
	PersistJSON // persist-json
)
