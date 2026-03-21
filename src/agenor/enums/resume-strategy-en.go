package enums

//go:generate stringer -type=ResumeStrategy -linecomment -trimprefix=ResumeStrategy -output resume-strategy-en-auto.go

// ResumeStrategy represents the different resume strategies
type ResumeStrategy uint

// If these enum definitions change, the test data (eg, resume-fastward.json) also needs
// to be updated.

const (
	// ResumeStrategyUndefined strategy undefined
	ResumeStrategyUndefined ResumeStrategy = iota // undefined-resume-strategy

	// ResumeStrategySpawn spawn resume
	ResumeStrategySpawn // spawn-resume-strategy

	// ResumeStrategyFastward fastforward resume
	ResumeStrategyFastward // fastward-resume-strategy
)
