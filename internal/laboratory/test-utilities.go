package lab

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/snivilised/traverse/pref"
	"github.com/snivilised/traverse/test/hydra"
)

type (
	NamedFunc func(name string) string
)

var (
	Reasons = struct {
		Node NamedFunc
	}{
		Node: func(name string) string {
			return fmt.Sprintf("❌ for node named: '%v'", name)
		},
	}
)

func Normalise(p string) string {
	return strings.ReplaceAll(p, "/", string(filepath.Separator))
}

func Because(name, because string) string {
	return fmt.Sprintf("❌ for node named: '%v', because: '%v'", name, because)
}

func Reason(name string) string {
	return fmt.Sprintf("❌ for node named: '%v'", name)
}

// yoke is similar to filepath.Join but it is meant specifically for relative file
// systems where the rules of a path are different; see fs.ValidPath
func yoke(segments ...string) string {
	return strings.Join(segments, "/")
}

func GetJSONPath() string {
	jroot := hydra.Repo(filepath.Join("test", "json"))

	return yoke(jroot, "unmarshal", Static.JSONFile)
}

func IgnoreFault(_ *pref.NavigationFault) error {
	return nil
}
