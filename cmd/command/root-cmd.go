// Package command provides CLI commands for the jay application.
package command

const (
	AppEmoji        = "🍒"
	ApplicationName = "jay"
	RootPsName      = "root-ps"
	SourceID        = "github.com/snivilised/agenor"
)

func Execute() error {
	return (&Bootstrap{}).Root().Execute()
}

// RootParameterSet defines the configuration options exposed on the
// root command's parameter set (CLIENT-TODO: refine these properties).
type RootParameterSet struct {
	// Language defines the IETF BCP 47 language tag.
	Language string
}
