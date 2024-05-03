package pref

import (
	"log/slog"
)

type (
	MonitorOptions struct {
		Log *slog.Logger
	}

	LogRotationOptions struct {
		// MaxSizeInMb, max size of a log file, before it is re-cycled
		MaxSizeInMb int

		// MaxNoOfBackups, max number of legacy log files that can exist
		// before being deleted
		MaxNoOfBackups int

		// MaxAgeInDays, max no of days before old log file is deleted
		MaxAgeInDays int
	}
)
