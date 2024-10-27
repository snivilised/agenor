package lab

import (
	"fmt"
	"path/filepath"
	"strings"
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
