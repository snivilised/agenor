// Package command provides CLI commands for the jay application.
package command

func Execute() error {
	return (&Bootstrap{}).Root().Execute()
}
