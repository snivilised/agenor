package pref

type (
	// PeerBehaviour contains options concerning peer characteristics
	PeerBehaviour struct {
		// IsActive peer buffering is enabled
		IsActive bool
	}

	// ViewBehaviours contains rendering options required by views
	ViewBehaviours struct {
		// Peer indicates wether the navigator should indicate
		// wether the current node is the last of its peers.
		Peer PeerBehaviour
	}
)

// WithTraversalConfigurer enables adhoc options configuration
func WithTraversalConfigurer(c TraversalConfigurer) Option {
	return func(o *Options) error {
		o.Configurer = c
		return nil
	}
}
