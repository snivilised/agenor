package core

type (
	// HibernationBehaviour defines hibernation behaviours
	HibernationBehaviour struct {
		// InclusiveWake when wake occurs, permit client callback to
		// be invoked for the current node. Inclusive, true by default
		InclusiveWake bool

		// InclusiveSleep when sleep occurs, permit client callback to
		// be invoked for the current node. Exclusive, false by default.
		InclusiveSleep bool
	}

	// HibernateOptions defines hibernation options
	HibernateOptions struct {
		// WakeAt defines a filter for hibernation wake condition
		WakeAt *FilterDef

		// SleepAt defines a filter for hibernation sleep condition
		SleepAt *FilterDef

		// Behaviour contains hibernation behavioural aspects
		Behaviour HibernationBehaviour
	}
)

// IsHibernateActive returns true if either WakeAt or SleepAt is defined,
// indicating that hibernation is active, and false otherwise. This can be
// used to determine whether the hibernation functionality should be engaged
// during traversal based on the presence of these filters.
func (o *HibernateOptions) IsHibernateActive() bool {
	return o.WakeAt != nil || o.SleepAt != nil
}
