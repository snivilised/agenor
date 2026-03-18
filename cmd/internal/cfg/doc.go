// Package cfg handles all configuration concerns for jay, including loading,
// validation, and access to config values. Mapped sections (interaction,
// advanced, logging) are decoded into strongly-typed Go structs via
// mapstructure. Unstructured sections (actions, pipelines, flags) are
// exposed as raw map values for consumer-driven handling.
package cfg
