package enums

//go:generate stringer -type=Metric -linecomment -trimprefix=Metric -output metric-en-auto.go

type Metric uint

// if new metrics are added, ensure that navigationMetricsFactory.new is kept
// in sync.
const (
	// MetricNoFilesInvoked represents the no of files invoked for during traversal
	//
	MetricNoFilesInvoked Metric = iota // metric-no-of-files

	// MetricNoFilesFilteredOut represents the no of files filtered out
	//
	MetricNoFilesFilteredOut // metric-no-of-files-filtered-out

	// MetricNoFoldersInvoked represents the no of folders invoked for during traversal
	//
	MetricNoFoldersInvoked // metric-no-of-folders

	// MetricNoFoldersFilteredOut represents the no of folders filtered out
	//
	MetricNoFoldersFilteredOut // metric-no-of-folders-filtered-out

	// MetricNoChildFilesFound represents the number of children files
	// of a particular directory that pass the compound filter when using the folders
	// with files subscription
	//
	MetricNoChildFilesFound // metric-no-of-child-files-found

	// MetricNoChildFilesFilteredOut represents the number of children files
	// of a particular directory that fail to pass the compound filter when using
	// the folders with files subscription
	//
	MetricNoChildFilesFilteredOut // metric-no-of-child-files-found
)
