package enums

//go:generate stringer -type=SkipTraversal -linecomment -trimprefix=Skip -output skip-traversal-en-auto.go

// SkipTraversal represents the different skip traversal strategies
type SkipTraversal uint

const (
	// SkipNoneTraversal no skip traversal
	SkipNoneTraversal SkipTraversal = iota // skip-none

	// SkipDirTraversal skip directory only traversal
	SkipDirTraversal // skip-dir

	// SkipAllTraversal skip entire traversal
	SkipAllTraversal // skip-all
)
