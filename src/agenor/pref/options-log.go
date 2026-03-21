package pref

import (
	"log/slog"
)

type (
	// AdminOptions defines options for admin related configurations.
	AdminOptions struct {
		// Path specifies the path for admin related files.
		Path string
	}

	// MonitorOptions represents the options for monitoring the traversal process.
	MonitorOptions struct {
		// Log is the logger used for logging messages during the traversal process.
		Log *slog.Logger

		// Admin specifies the options for admin related configurations.
		Admin AdminOptions
	}
)

// WithAdminPath defines the path for admin related files
func WithAdminPath(path string) Option {
	return func(o *Options) error {
		o.Monitor.Admin.Path = path

		return nil
	}
}

// WithLogger defines a structure logger
func WithLogger(logger *slog.Logger) Option {
	return func(o *Options) error {
		o.Monitor.Log = logger

		return nil
	}
}
