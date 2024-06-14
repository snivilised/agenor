package pref

type CoreOptions struct {
	// Behaviours collection of behaviours that adjust the way navigation occurs,
	// that can be tweaked by the client.
	//
	Behaviours NavigationBehaviours

	// Sampling options
	//
	Sampling SamplingOptions

	// Filter
	//
	Filter FilterOptions

	// Hibernation
	//
	Hibernate HibernateOptions

	// Concurrency contains options relating concurrency
	//
	Concurrency ConcurrencyOptions

	// Persist contains options for persisting traverse options
	//
	Persist PersistOptions
}
