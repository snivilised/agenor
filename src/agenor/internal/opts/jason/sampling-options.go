package jason

import (
	"github.com/snivilised/jaywalk/src/agenor/enums"
)

type (
	// EntryQuantities contains specification of no of files and directories
	// used in various contexts, but primarily sampling.
	EntryQuantities struct {
		// Files no of files
		Files uint `json:"no-of-files"`

		// Directories no of directories
		Directories uint `json:"no-of-directories"`
	}

	// SamplingOptions contains options relating to sampling, which is the
	// process of selecting a subset of entries from a directory's contents
	// during traversal. Sampling can be used to limit the number of entries
	// processed, or to select specific entries based on certain criteria.
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
