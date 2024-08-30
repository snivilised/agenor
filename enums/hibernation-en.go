package enums

//go:generate stringer -type=Hibernation -linecomment -trimprefix=Hibernation -output hibernation-en-auto.go

// Hibernation denotes whether user defined callback is being invoked.
type Hibernation uint

const (
	HibernationUndefined Hibernation = iota // undefined

	// HibernationPending conditional listening is awaiting activation
	//
	HibernationPending // pending-hibernation

	// HibernationActive conditional listening is active (callback is invoked)
	//
	HibernationActive // active-hibernation

	// HibernationRetired conditional listening is now deactivated
	//
	HibernationRetired // retired-hibernation
)
