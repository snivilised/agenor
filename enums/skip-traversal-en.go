package enums

//go:generate stringer -type=SkipTraversal -linecomment -trimprefix=Skip -output skip-traversal-en-auto.go

type SkipTraversal uint

const (
	SkipNoneTraversal SkipTraversal = iota // skip-none
	SkipDirTraversal                       // skip-dir
	SkipAllTraversal                       // skip-all
)
