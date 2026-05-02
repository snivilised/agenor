package command

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/snivilised/jaywalk/src/locale"
)

// MarkInheritedFlagRequired checks that the named flag has been set,
// searching both local and inherited persistent flags. Returns an error
// if the flag was not supplied by the user.
//
// Use in place of cobra's MarkFlagRequired when the flag is defined on
// a parent command and inherited via PersistentFlags.
// See https://github.com/spf13/cobra/issues/921.
func MarkInheritedFlagRequired(cmd *cobra.Command, name string) error {
	cmd.InheritedFlags()

	if f := cmd.Flags().Lookup(name); f != nil && f.Changed {
		return nil
	}

	return locale.NewMarkInheritedFlagsRequiredError(cmd.Name(), name)
}

// MarkInheritedFlagsOneRequired checks that at least one of the named
// flags has been set, searching both local and inherited persistent flags.
// Returns an error if none of the flags were supplied by the user.
//
// Use in place of cobra's MarkFlagsOneRequired when the flags are defined
// on a parent command and inherited via PersistentFlags.
// See https://github.com/spf13/cobra/issues/921.
func MarkInheritedFlagsOneRequired(cmd *cobra.Command, names ...string) error {
	cmd.InheritedFlags()

	for _, name := range names {
		if f := cmd.Flags().Lookup(name); f != nil && f.Changed {
			return nil
		}
	}

	return locale.NewMarkInheritedFlagsOneRequiredError(
		cmd.Name(),
		strings.Join(names, ", "),
	)
}

// MarkInheritedFlagsMutuallyExclusive checks that at most one of the named
// flags has been set, searching both local and inherited persistent flags.
// Returns an error if more than one of the flags were supplied by the user.
//
// Use in place of cobra's MarkFlagsMutuallyExclusive when the flags are
// defined on a parent command and inherited via PersistentFlags.
// See https://github.com/spf13/cobra/issues/921.
func MarkInheritedFlagsMutuallyExclusive(cmd *cobra.Command, names ...string) error {
	cmd.InheritedFlags()

	var changed []string

	for _, name := range names {
		if f := cmd.Flags().Lookup(name); f != nil && f.Changed {
			changed = append(changed, name)
		}
	}

	if len(changed) > 1 {
		return locale.NewMutuallyExclusiveFlagsPresentError(
			cmd.Name(),
			strings.Join(changed, ", "),
		)
	}

	return nil
}
