package pref

import (
	"log/slog"
)

type (
	AdminOptions struct {
		Path string
	}

	MonitorOptions struct {
		Log   *slog.Logger
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
