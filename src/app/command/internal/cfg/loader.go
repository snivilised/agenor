package cfg

import (
	"fmt"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// ---------------------------------------------------------------------------
// Format registry - extensible without touching loader logic
// ---------------------------------------------------------------------------

// Format names a supported configuration file format.
type Format string

const (
	FormatYAML Format = "yaml"
	FormatJSON Format = "json"
	FormatTOML Format = "toml"
	// FormatLua can be registered once a Viper remote/custom source exists.
)

// registeredFormats is the set of formats Viper is told to accept.
// Add new entries here as support is implemented.
var registeredFormats = map[Format]struct{}{
	FormatYAML: {},
	FormatJSON: {},
	FormatTOML: {},
}

// RegisterFormat adds a new format to the accepted set.  Call from init()
// in the package that provides the format-specific codec.
func RegisterFormat(f Format) {
	registeredFormats[f] = struct{}{}
}

// ---------------------------------------------------------------------------
// LoadOptions controls where and how the config is read.
// ---------------------------------------------------------------------------

// LoadOptions is the input value object for Load.
type LoadOptions struct {
	// ConfigFile is an explicit path to the config file.  When non-empty,
	// Format is inferred from the extension unless Format is also set.
	ConfigFile string

	// ConfigName is the base name (without extension) when no explicit
	// ConfigFile is given.  Defaults to "jay".
	ConfigName string

	// ConfigPaths lists directories to search, in priority order.
	// Defaults to [".", "$HOME/.config/jay", "/etc/jay"].
	ConfigPaths []string

	// Format overrides automatic format detection.
	Format Format

	// EnvPrefix is used to bind environment variables.  Defaults to "JAY".
	EnvPrefix string

	// ViperInstance allows callers (including tests) to supply a pre-configured
	// Viper instance, enabling full isolation.
	ViperInstance *viper.Viper
}

// isPreloaded reports whether the caller supplied a ready-to-use Viper
// instance (already loaded from memory or file).  When true, Load must not
// attempt file discovery or call ReadInConfig.
func (o *LoadOptions) isPreloaded() bool {
	return o.ViperInstance != nil
}

func (o *LoadOptions) applyDefaults() {
	if o.EnvPrefix == "" {
		o.EnvPrefix = "JAY"
	}

	// File-discovery defaults are only meaningful when the caller has NOT
	// supplied a pre-loaded Viper instance.  Setting them unconditionally was
	// causing Viper to attempt a filesystem read even in unit tests that
	// populate the instance from an in-memory reader.
	if !o.isPreloaded() {
		if o.ConfigName == "" {
			o.ConfigName = "jay"
		}
		if len(o.ConfigPaths) == 0 {
			o.ConfigPaths = []string{".", "$HOME/.config/jay", "/etc/jay"}
		}
		// Only allocate a fresh Viper when no instance was provided at all.
		o.ViperInstance = viper.New()
	}
}

// ---------------------------------------------------------------------------
// Load
// ---------------------------------------------------------------------------

// Load reads, decodes, and validates the configuration for jay.
// It returns a fully populated *Config or an error describing every problem
// found (validation failures are aggregated, not returned one at a time).
//
// Callers should call Load once on startup, then use the returned *Config
// for the lifetime of the process.
func Load(opts LoadOptions) (*Config, error) {
	preloaded := opts.isPreloaded()
	opts.applyDefaults()
	v := opts.ViperInstance

	if preloaded {
		// Instance was pre-populated by the caller (e.g. from an in-memory
		// reader in a test).  We still apply env-var bindings but must not
		// touch file-discovery settings or call ReadInConfig - doing so would
		// discard the already-loaded content and attempt a filesystem read.
		applyEnvBindings(v, opts.EnvPrefix)
	} else {
		if err := configureViper(v, opts); err != nil {
			return nil, fmt.Errorf("cfg.Load: viper setup: %w", err)
		}
		if err := v.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("cfg.Load: reading config: %w", err)
		}
	}

	cfg, err := decode(v)
	if err != nil {
		return nil, fmt.Errorf("cfg.Load: decoding: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("cfg.Load: validation: %w", err)
	}

	return cfg, nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// applyEnvBindings wires environment-variable support onto an existing Viper
// instance.  It is called for both file-based and pre-loaded instances so
// that env-var overrides work consistently in all cases.
func applyEnvBindings(v *viper.Viper, envPrefix string) {
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	v.AutomaticEnv()
}

func configureViper(v *viper.Viper, opts LoadOptions) error {
	// Environment variable binding
	applyEnvBindings(v, opts.EnvPrefix)

	if opts.ConfigFile != "" {
		v.SetConfigFile(opts.ConfigFile)
		// Honour explicit format override even when a file path is given.
		if opts.Format != "" {
			if _, ok := registeredFormats[opts.Format]; !ok {
				return fmt.Errorf("unsupported format %q", opts.Format)
			}
			v.SetConfigType(string(opts.Format))
		}
		return nil
	}

	// Name-based discovery
	v.SetConfigName(opts.ConfigName)

	if opts.Format != "" {
		if _, ok := registeredFormats[opts.Format]; !ok {
			return fmt.Errorf("unsupported format %q", opts.Format)
		}
		v.SetConfigType(string(opts.Format))
	}

	for _, p := range opts.ConfigPaths {
		v.AddConfigPath(p)
	}

	return nil
}

// decode pulls all sections out of Viper and populates a Config.
func decode(v *viper.Viper) (*Config, error) {
	cfg := &Config{}

	// Shared mapstructure decoder config - handles time.Duration strings,
	// weak typing for numeric/bool values coming from env vars.
	decoderCfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		),
		TagName: "mapstructure",
	}

	if err := decodeSection(v, "interaction", decoderCfg, &cfg.Mapped.Interaction); err != nil {
		return nil, err
	}
	if err := decodeSection(v, "advanced", decoderCfg, &cfg.Mapped.Advanced); err != nil {
		return nil, err
	}
	if err := decodeSection(v, "logging", decoderCfg, &cfg.Mapped.Logging); err != nil {
		return nil, err
	}

	// Raw sections: decode into loosely-typed structs (still mapstructure, but
	// the target types use map[string]any for open-ended sub-keys).
	rawDecoderCfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		),
		TagName: "mapstructure",
	}

	if err := decodeSection(v, "actions", rawDecoderCfg, &cfg.Raw.Actions); err != nil {
		return nil, err
	}
	if err := decodeSection(v, "pipelines", rawDecoderCfg, &cfg.Raw.Pipelines); err != nil {
		return nil, err
	}
	if err := decodeFlagsSection(v, rawDecoderCfg, &cfg.Raw.Flags); err != nil {
		return nil, err
	}

	return cfg, nil
}

// decodeSection extracts key from Viper and mapstructure-decodes it into dest.
func decodeSection(v *viper.Viper, key string, decoderCfg *mapstructure.DecoderConfig, dest any) error {
	raw := v.Get(key)
	if raw == nil {
		// Section absent - leave dest at zero value.
		return nil
	}

	decoderCfg.Result = dest
	decoder, err := mapstructure.NewDecoder(decoderCfg)
	if err != nil {
		return fmt.Errorf("creating decoder for %q: %w", key, err)
	}
	if err := decoder.Decode(raw); err != nil {
		return fmt.Errorf("decoding section %q: %w", key, err)
	}
	return nil
}

// decodeFlagsSection has custom logic because flags.short is nested deeply.
func decodeFlagsSection(v *viper.Viper, decoderCfg *mapstructure.DecoderConfig, destinationCfg *FlagsConfig) error {
	raw := v.Get("flags")
	if raw == nil {
		return nil
	}

	rawMap, ok := raw.(map[string]any)
	if !ok {
		return fmt.Errorf("flags section has unexpected type %T", raw)
	}

	// short.overrides.cmds  ─►  FlagShortOverride = map[cmd]map[flag]letter
	if shortRaw, ok := rawMap["short"]; ok {
		if shortMap, ok := toStringAnyMap(shortRaw); ok {
			if overridesRaw, ok := shortMap["overrides"]; ok {
				if overridesMap, ok := toStringAnyMap(overridesRaw); ok {
					if cmdsRaw, ok := overridesMap["cmds"]; ok {
						destinationCfg.Short = make(FlagShortOverride)
						if cmdsMap, ok := toStringAnyMap(cmdsRaw); ok {
							for cmd, flagsRaw := range cmdsMap {
								if flagsMap, ok := toStringAnyMap(flagsRaw); ok {
									destinationCfg.Short[cmd] = make(map[string]string)
									for flag, val := range flagsMap {
										destinationCfg.Short[cmd][flag] = fmt.Sprintf("%v", val)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// invoke.cmds  ─►  FlagInvokeDefaults = map[cmd]map[flag]any
	if invokeRaw, ok := rawMap["invoke"]; ok {
		if invokeMap, ok := toStringAnyMap(invokeRaw); ok {
			if cmdsRaw, ok := invokeMap["cmds"]; ok {
				destinationCfg.Invoke = make(FlagInvokeDefaults)
				if cmdsMap, ok := toStringAnyMap(cmdsRaw); ok {
					for cmd, flagsRaw := range cmdsMap {
						if flagsMap, ok := toStringAnyMap(flagsRaw); ok {
							destinationCfg.Invoke[cmd] = flagsMap
						}
					}
				}
			}
		}
	}

	// component  ─►  FlagComponentDefaults = map[component]map[flag]any
	if compRaw, ok := rawMap["component"]; ok {
		destinationCfg.Component = make(FlagComponentDefaults)
		if compMap, ok := toStringAnyMap(compRaw); ok {
			for comp, flagsRaw := range compMap {
				if flagsMap, ok := toStringAnyMap(flagsRaw); ok {
					destinationCfg.Component[comp] = flagsMap
				}
			}
		}
	}

	_ = decoderCfg // reserved for future hook composition
	return nil
}

// toStringAnyMap is a safe cast to map[string]any, handling both concrete
// types that Viper/YAML can produce.
func toStringAnyMap(v any) (map[string]any, bool) {
	switch m := v.(type) {
	case map[string]any:
		return m, true
	case map[any]any:
		out := make(map[string]any, len(m))
		for k, val := range m {
			out[fmt.Sprintf("%v", k)] = val
		}
		return out, true
	}
	return nil, false
}

// ---------------------------------------------------------------------------
// Sentinel for zero-value Duration detection
// ---------------------------------------------------------------------------

var _ = time.Duration(0) // ensure time is used
