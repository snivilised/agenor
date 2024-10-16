package lab

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Join creates a path from the parent combined with the relative path. The relative
// path is a file system path so should only contain forward slashes, not the standard
// file path separator as denoted by filepath.Separator, typically used when interacting
// with the local file system. Do not use trailing "/".
func Join(parent, relative string) string {
	if relative == "" {
		return parent
	}

	return parent + "/" + relative
}

// Repo gets the path of the repo with relative joined on
func Repo(relative string) string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, _ := cmd.Output()
	repo := strings.TrimSpace(string(output))

	return Join(repo, relative)
}

func Normalise(p string) string {
	return strings.ReplaceAll(p, "/", string(filepath.Separator))
}

func Because(name, because string) string {
	return fmt.Sprintf("❌ for item named: '%v', because: '%v'", name, because)
}

func Reason(name string) string {
	return fmt.Sprintf("❌ for item named: '%v'", name)
}

func Log() string {
	return Repo("Test/test.log")
}
