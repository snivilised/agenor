package enums

//go:generate stringer -type=SampleType -linecomment -trimprefix=SampleType -output sample-type-en-auto.go

// SampleType determines the type of sampling to use
type SampleType uint

const (
	SampleTypeUndefined SampleType = iota // undefined-sample
	SampleTypeSlice                       // slice-sample
	SampleTypeFilter                      // filter-sample
	SampleTypeCustom                      // custom-sample
)
