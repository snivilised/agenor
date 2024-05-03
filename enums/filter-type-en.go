package enums

//go:generate stringer -type=FilterType -linecomment -trimprefix=Filter -output filter-type-en-auto.go

type FilterType uint

const (
	FilterTypeUndefined FilterType = iota // undefined-filter

	// FilterTypeExtendedGlob is the preferred filter type as it the most
	// user friendly. The base part of the name is filtered by a glob
	// and the suffix is filtered by a list of defined extensions. The pattern
	// for the extended filter type is composed of 2 parts; the first is a
	// glob, which is applied to the base part of the name. The second part
	// is a csv of required extensions to filter for. The pattern is specified
	// in the form: "<base-glob>|ext1,ext2...". Each extension may include a
	// a leading dot. An example pattern definition would be:
	// "cover.*|.jpg,jpeg"
	//
	FilterTypeExtendedGlob // extended-glob-filter

	// FilterTypeRegex regex filter
	//
	FilterTypeRegex // regex-filter

	// FilterTypeGlob glob filter
	//
	FilterTypeGlob // glob-filter

	// FilterTypeCustom client definable filter
	//
	FilterTypeCustom // custom-filter

	// FilterTypePoly poly filter
	//
	FilterTypePoly // poly-filter
)
