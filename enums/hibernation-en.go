package enums

//go:generate stringer -type=Hibernation -linecomment -trimprefix=Hibernation -output hibernation-en-auto.go

// Hibernation denotes whether user defined callback is being invoked.
type Hibernation uint

const (
	HibernationUndefined Hibernation = iota // undefined

	// HibernationSleep listen not active, callback always invoked (subject to filtering)
	//
	HibernationSleep // sleep-hibernation

	// HibernationFastward listen used to resume by fast-forwarding
	//
	HibernationFastward // fastward-hibernation

	// HibernationPending conditional listening is awaiting activation
	//
	HibernationPending // pending-hibernation

	// HibernationAwake conditional listening is active (callback is invoked)
	//
	HibernationAwake // awake-hibernation

	// HibernationRetired conditional listening is now deactivated
	//
	HibernationRetired // retired-hibernation
)
