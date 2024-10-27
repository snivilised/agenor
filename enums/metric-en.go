package enums

//go:generate stringer -type=Metric -linecomment -trimprefix=Metric -output metric-en-auto.go

type Metric uint

// if new metrics are added, ensure that navigationMetricsFactory.new is kept
// in sync.
const (
	_ Metric = iota

	// MetricNoFilesInvoked represents the no of files invoked for during traversal
	//
	MetricNoFilesInvoked // metric-no-of-files

	// MetricNoFilesFilteredOut represents the no of files filtered out
	//
	MetricNoFilesFilteredOut // metric-no-of-files-filtered-out

	// MetricNoDirectoriesInvoked represents the no of directories invoked for during traversal
	//
	MetricNoDirectoriesInvoked // metric-no-of-directories

	// MetricNoDirectoriesFilteredOut represents the no of directories filtered out
	//
	MetricNoDirectoriesFilteredOut // metric-no-of-directories-filtered-out

	// MetricNoChildFilesFound represents the number of children files
	// of a particular directory that pass the compound filter when using the directories
	// with files subscription
	//
	MetricNoChildFilesFound // metric-no-of-child-files-found

	// MetricNoChildFilesFilteredOut represents the number of children files
	// of a particular directory that fail to pass the compound filter when using
	// the directories with files subscription
	//
	MetricNoChildFilesFilteredOut // metric-no-of-child-files-found
)
