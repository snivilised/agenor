package pref

type (
	// RescueData is the data to be saved in case of a rescue.
	RescueData interface {
		// Data returns the data to be saved.
		Data() interface{}
	}

	// Recovery is the interface for saving rescue data.
	Recovery interface {
		// Save saves the rescue data and returns the path write to.
		Save(data RescueData) (string, error)
	}
)
