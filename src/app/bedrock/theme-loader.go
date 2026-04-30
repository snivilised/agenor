package bedrock

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"

	"github.com/snivilised/jaywalk/src/prism"
)

const (
	// ThemesDirEnvVar is the environment variable that overrides the
	// default XDG themes directory. When set, jay loads themes from
	// this path instead of ~/.config/jay/themes/.
	ThemesDirEnvVar = "JAY_THEMES_DIR"

	// ThemeSystemName is the reserved name for the built-in ANSI-16
	// palette that respects the user's terminal theme. Selecting this
	// name bypasses file loading entirely.
	ThemeSystemName = "system"
)

// ThemeLoader resolves and loads a named theme from the themes directory.
// Constructed once by Bootstrap and shared across the application lifetime.
type ThemeLoader struct {
	themesDir string
}

// NewThemeLoader constructs a ThemeLoader. The themes directory is
// resolved from JAY_THEMES_DIR if set, otherwise from the XDG config
// base at ~/.config/jay/themes/.
func NewThemeLoader() *ThemeLoader {
	dir := os.Getenv(ThemesDirEnvVar)

	if dir == "" {
		dir = filepath.Join(xdg.ConfigHome, "jay", "themes")
	}

	return &ThemeLoader{
		themesDir: dir,
	}
}

// Load returns the prism.Palette for the named theme. The name
// "system" returns SystemPalette() without reading any file.
// Any other name is resolved to <themesDir>/<name>.yaml and decoded
// via mapstructure. Returns an error if the file does not exist or
// cannot be decoded.
func (tl *ThemeLoader) Load(name string) (prism.Palette, error) {
	if name == "" || name == ThemeSystemName {
		return prism.SystemPalette(), nil
	}

	path := filepath.Join(tl.themesDir, name+".yaml")

	v := viper.New()
	v.SetConfigType("yaml")
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return prism.Palette{}, fmt.Errorf(
				"theme %q not found - expected file at %s",
				name,
				path,
			)
		}

		return prism.Palette{}, fmt.Errorf(
			"reading theme %q: %w",
			name,
			err,
		)
	}

	raw := v.Sub("palette")
	if raw == nil {
		return prism.Palette{}, fmt.Errorf(
			"theme %q: missing required 'palette' key",
			name,
		)
	}

	var palette prism.Palette

	if err := raw.Unmarshal(&palette, mapstructureTagOption()); err != nil {
		return prism.Palette{}, fmt.Errorf(
			"theme %q: decoding palette: %w",
			name,
			err,
		)
	}

	return palette, nil
}

// ThemesDir returns the resolved themes directory path. Useful for
// error messages and the --theme flag's help text.
func (tl *ThemeLoader) ThemesDir() string {
	return tl.themesDir
}

// mapstructureTagOption returns a viper.DecoderConfigOption that
// instructs mapstructure to use the "mapstructure" struct tag when
// decoding YAML keys into Go field names. This enables kebab-case YAML
// keys to map correctly into the palette structs via their
// mapstructure tags.
//
// The explicit viper.DecoderConfigOption cast is required because Go
// does not implicitly convert a func literal to a named function type
// even when the underlying signatures are identical.
func mapstructureTagOption() viper.DecoderConfigOption {
	return viper.DecoderConfigOption(func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "mapstructure"
	})
}
