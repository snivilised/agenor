package json

import "fmt"

type (
	// HibernationBehaviour
	HibernationBehaviour struct {
		// InclusiveWake when wake occurs, permit client callback to
		// be invoked for the current node. Inclusive, true by default
		InclusiveWake bool `json:"hibernate-inclusive-wake"`

		// InclusiveSleep when sleep occurs, permit client callback to
		// be invoked for the current node. Exclusive, false by default.
		InclusiveSleep bool `json:"hibernate-inclusive-sleep"`
	}

	// HibernateOptions
	HibernateOptions struct {
		// WakeAt defines a filter for hibernation wake condition
		WakeAt *FilterDef

		// SleepAt defines a filter for hibernation sleep condition
		SleepAt *FilterDef

		// Behaviour contains hibernation behavioural aspects
		Behaviour HibernationBehaviour
	}
)

func (b *HibernationBehaviour) String() string {
	return fmt.Sprintf("[HibernationBehaviour] inclusive wake: %v, inclusive sleep: %v",
		b.InclusiveWake, b.InclusiveSleep,
	)
}
