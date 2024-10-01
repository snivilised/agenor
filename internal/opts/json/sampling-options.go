package json

import (
	"github.com/snivilised/traverse/enums"
)

type (
	// EntryQuantities contains specification of no of files and folders
	// used in various contexts, but primarily sampling.
	EntryQuantities struct {
		Files   uint `json:"no-of-files"`
		Folders uint `json:"no-of-folders"`
	}

	// SamplingOptions
	SamplingOptions struct {
		// SampleType the type of sampling to use
		SampleType enums.SampleType `json:"sample-type"`

		// SampleInReverse determines the direction of iteration for the sampling
		// operation
		SampleInReverse bool `json:"sample-in-reverse"`

		// NoOf specifies number of items required in each sample (only applies
		// when not using Custom iterator options)
		NoOf EntryQuantities
	}
)
