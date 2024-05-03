package enums

//go:generate stringer -type=ResumeStrategy -linecomment -trimprefix=ResumeStrategy -output resume-strategy-en-auto.go

type ResumeStrategy uint

// If these enum definitions change, the test data (eg, resume-fastward.json) also needs
// to be updated.

const (
	ResumeStrategyUndefined ResumeStrategy = iota // undefined-resume-strategy
	ResumeStrategySpawn                           // spawn-resume-strategy
	ResumeStrategyFastward                        // fastward-resume-strategy
)
