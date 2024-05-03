package pref

import (
	"github.com/snivilised/traverse/enums"
)

type (
	// EntryQuantities contains specification of no of files and folders
	// used in various contexts, but primarily sampling.
	EntryQuantities struct {
		Files   uint
		Folders uint
	}

	// SamplingOptions
	SamplingOptions struct {
		// SampleInReverse determines the direction of iteration for the sampling
		// operation
		SampleInReverse bool

		// SampleType the type of sampling to use
		SampleType enums.SampleType

		// NoOf specifies number of items required in each sample (only applies
		// when not using Custom iterator options)
		NoOf EntryQuantities
	}
)

func WithSamplingOptions(so *SamplingOptions) OptionFn {
	return func(o *Options, _ *Registry) error {
		o.Core.Sampling = *so

		return nil
	}
}

func WithSamplingInReverse() OptionFn {
	return func(o *Options, _ *Registry) error {
		o.Core.Sampling.SampleInReverse = true

		return nil
	}
}

func WithSamplingType(sample enums.SampleType) OptionFn {
	return func(o *Options, _ *Registry) error {
		o.Core.Sampling.SampleType = sample

		return nil
	}
}

func WithSamplingNoOf(noOf *EntryQuantities) OptionFn {
	return func(o *Options, _ *Registry) error {
		o.Core.Sampling.NoOf = *noOf

		return nil
	}
}

func WithSampling(files, folders uint) OptionFn {
	return func(o *Options, _ *Registry) error {
		o.Core.Sampling.NoOf.Files = files
		o.Core.Sampling.NoOf.Folders = folders

		return nil
	}
}
