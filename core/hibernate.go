package core

type (
	// HibernationBehaviour
	HibernationBehaviour struct {
		// InclusiveWake when wake occurs, permit client callback to
		// be invoked for the current node. Inclusive, true by default
		InclusiveWake bool

		// InclusiveSleep when sleep occurs, permit client callback to
		// be invoked for the current node. Exclusive, false by default.
		InclusiveSleep bool
	}

	HibernateOptions struct {
		// WakeAt defines a filter for hibernation wake condition
		WakeAt *FilterDef

		// SleepAt defines a filter for hibernation sleep condition
		SleepAt *FilterDef

		// Behaviour contains hibernation behavioural aspects
		Behaviour HibernationBehaviour
	}
)

func (o *HibernateOptions) IsHibernateActive() bool {
	return o.WakeAt != nil || o.SleepAt != nil
}
