package cfg_test

import (
	"strings"

	"github.com/spf13/viper"
)

// viperFromYAML returns a Viper instance populated from the given YAML string.
// It is used by tests to avoid touching the filesystem.
func viperFromYAML(yaml string) *viper.Viper {
	v := viper.New()
	v.SetConfigType("yaml")
	_ = v.ReadConfig(strings.NewReader(yaml))
	return v
}
