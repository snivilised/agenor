package enums

//go:generate stringer -type=ConfigFormat -linecomment -trimprefix=ConfigFormat -output config-format-en-auto.go

// ConfigFormat used to enable selecting directory entry type.
type ConfigFormat uint

const (
	// ConfigFormatDirectory represents a directory entry.
	//
	ConfigFormatYaml ConfigFormat = iota // config-format-yaml

	// ConfigFormatFile represents a file entry.
	//
	ConfigFormatJson // config-format-json
)
