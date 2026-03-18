package cfg

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// ---------------------------------------------------------------------------
// FlagResolver
// ---------------------------------------------------------------------------

// FlagResolver merges config-level flag defaults with explicit CLI values.
// Priority (highest → lowest):
//  1. Flag explicitly set on the command line by the user.
//  2. flags.invoke.cmds.<cmd> in config (command-specific override).
//  3. flags.invoke.cmds.any  in config (wildcard command override).
//  4. flags.component.<component> in config (component-level default).
//  5. Cobra / pflag default value (defined in the command itself).
type FlagResolver struct {
	flags FlagsConfig
}

// NewFlagResolver creates a FlagResolver backed by the parsed flags config.
func NewFlagResolver(flags FlagsConfig) *FlagResolver {
	return &FlagResolver{flags: flags}
}

// ResolveInt returns the effective int value for a flag on a named command,
// optionally scoped to a component.  The caller still owns the pflag.FlagSet;
// this function only reads from it (it never mutates Cobra state).
func (r *FlagResolver) ResolveInt(cmd *cobra.Command, flagName, component string) (int, bool) {
	// 1. Explicit CLI value wins unconditionally.
	if f := cmd.Flags().Lookup(flagName); f != nil && f.Changed {
		val, err := cmd.Flags().GetInt(flagName)
		return val, err == nil
	}

	// 2 & 3. flags.invoke.cmds
	if v, ok := r.lookupInvoke(cmd.Name(), flagName); ok {
		return toInt(v)
	}

	// 4. flags.component.<component>
	if component != "" {
		if compDefaults, ok := r.flags.Component[component]; ok {
			if v, ok := compDefaults[flagName]; ok {
				return toInt(v)
			}
		}
	}

	// 5. Cobra default - read back from pflag
	if f := cmd.Flags().Lookup(flagName); f != nil {
		val, err := cmd.Flags().GetInt(flagName)
		return val, err == nil
	}

	return 0, false
}

// ResolveString is the string analogue of ResolveInt.
func (r *FlagResolver) ResolveString(cmd *cobra.Command, flagName, component string) (string, bool) {
	if f := cmd.Flags().Lookup(flagName); f != nil && f.Changed {
		val, err := cmd.Flags().GetString(flagName)
		return val, err == nil
	}

	if v, ok := r.lookupInvoke(cmd.Name(), flagName); ok {
		s, ok := v.(string)
		return s, ok
	}

	if component != "" {
		if compDefaults, ok := r.flags.Component[component]; ok {
			if v, ok := compDefaults[flagName]; ok {
				s, ok := v.(string)
				return s, ok
			}
		}
	}

	if f := cmd.Flags().Lookup(flagName); f != nil {
		val, err := cmd.Flags().GetString(flagName)
		return val, err == nil
	}

	return "", false
}

// ApplyShortOverrides re-registers short flags for a command according to the
// flags.short.overrides section.  Call this before cmd.Execute().
//
// Note: pflag does not allow changing the shorthand of an already-registered
// flag; this helper works by re-registering the flag under the new shorthand
// via a persistent pre-run hook, applied before the cobra tree executes.
func (r *FlagResolver) ApplyShortOverrides(cmd *cobra.Command) {
	// TODO: ApplyShortOverrides is redundant; we can just pass in the overrides
	// into mamba's NewParamSet
	cmdName := cmd.Name()
	overrides, ok := r.flags.Short[cmdName]
	if !ok {
		return
	}

	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		newShort, ok := overrides[f.Name]
		if !ok || len(newShort) != 1 {
			return
		}
		// pflag stores shorthands as a single byte string.
		f.Shorthand = newShort
	})
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// lookupInvoke checks flags.invoke.cmds.<cmdName> then flags.invoke.cmds.any.
func (r *FlagResolver) lookupInvoke(cmdName, flagName string) (any, bool) {
	if cmdDefaults, ok := r.flags.Invoke[cmdName]; ok {
		if v, ok := cmdDefaults[flagName]; ok {
			return v, true
		}
	}
	if anyDefaults, ok := r.flags.Invoke["any"]; ok {
		if v, ok := anyDefaults[flagName]; ok {
			return v, true
		}
	}
	return nil, false
}

func toInt(v any) (int, bool) {
	switch n := v.(type) {
	case int:
		return n, true
	case int64:
		return int(n), true
	case float64:
		return int(n), true
	}
	return 0, false
}
