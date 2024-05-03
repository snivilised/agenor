package enums

//go:generate stringer -type=PersistenceFormat -linecomment -trimprefix=Persist -output persistence-format-en-auto.go

type PersistenceFormat uint

const (
	PersistUndefined PersistenceFormat = iota // persistence-undefined
	PersistJSON                               // persist-json
)
