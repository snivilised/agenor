//go:build !windows

package shell

import (
	"os/exec"
	"strings"

	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/locale"
)

// detect returns the Unix environment. On all non-Windows platforms jay
// always runs under a POSIX-compatible shell, so no environment variable
// inspection is needed.
func detect() (Environment, error) {
	env := Environment{
		Kind:   enums.ShellKindNativeUnix,
		Locate: unixLocate,
	}

	return env, nil
}

// unixLocate implements LocateFunc for Unix-like systems. It delegates
// entirely to the POSIX command -v built into /bin/sh, which correctly
// resolves binaries on PATH, shell builtins, and exported shell functions.
func unixLocate(name string) (string, error) {
	// command -v is POSIX-specified and handles all three token kinds:
	//   - binaries:         prints the resolved path, exits 0
	//   - builtins:         prints the name, exits 0
	//   - exported funcs:   prints the name, exits 0
	//   - not found:        prints nothing, exits non-zero
	//nolint:gosec // This is not a security issue as the input is controlled
	// by the user and not used in a context where injection is possible.
	cmd := exec.Command("/bin/sh", "-c", "command -v "+shellQuote(name))
	out, err := cmd.Output()

	if err != nil {
		return "", locale.NewCmdNotFoundInEnvError(name)
	}

	return strings.TrimSpace(string(out)), nil
}

// shellQuote wraps name in single quotes so that a token containing
// spaces or special characters does not break the -c argument. Single
// quotes inside the name are escaped by ending the quoted string,
// inserting an escaped quote, and reopening.
func shellQuote(name string) string {
	return "'" + strings.ReplaceAll(name, "'", `'\''`) + "'"
}
