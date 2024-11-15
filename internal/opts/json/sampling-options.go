package json

import (
	"github.com/snivilised/agenor/enums"
)

type (
	// EntryQuantities contains specification of no of files and directories
	// used in various contexts, but primarily sampling.
	EntryQuantities struct {
		Files       uint `json:"no-of-files"`
		Directories uint `json:"no-of-directories"`
	}

	// SamplingOptions
	SamplingOptions struct {
		// Type the type of sampling to use
		Type enums.SampleType `json:"sample-type"`

		// InReverse determines the direction of iteration for the sampling
		// operation
		InReverse bool `json:"in-reverse"`

		// NoOf specifies number of items required in each sample (only applies
		// when not using Custom iterator options)
		NoOf EntryQuantities
	}
)
