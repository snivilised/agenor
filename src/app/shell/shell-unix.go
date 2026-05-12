//go:build !windows

package shell

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/snivilised/jaywalk/src/agenor/core"
	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/locale"
)

// detect returns the Unix environment. On all non-Windows platforms jay
// always runs under a POSIX-compatible shell, so no environment variable
// inspection is needed.
func detect() (Environment, error) {
	env := Environment{
		Kind:    enums.ShellKindNativeUnix,
		Locate:  unixLocate,
		Execute: unixExec,
	}

	return env, nil
}

// unixLocate implements LocateFunc for Unix-like systems. It delegates
// entirely to the POSIX command -v built into /bin/sh, which correctly
// resolves binaries on PATH, shell builtins, and exported shell functions.
func unixLocate(name string) (string, error) {
	shell := filepath.Base(core.Getenv("SHELL"))
	if shell == "" {
		shell = "sh"
	}

	cmd := exec.Command(shell, "-ic", "command -v -- "+quote(name)) //nolint:gosec // ok
	out, err := cmd.Output()
	if err != nil {
		return "", locale.NewCmdNotFoundInEnvError(name)
	}

	s := strings.TrimSpace(string(out))
	if s == "" {
		return "", locale.NewCmdNotFoundInEnvError(name)
	}

	return s, nil
}

// quote wraps name in single quotes so that a token containing
// spaces or special characters does not break the -c argument. Single
// quotes inside the name are escaped by ending the quoted string,
// inserting an escaped quote, and reopening.
func quote(name string) string {
	return "'" + strings.ReplaceAll(name, "'", `'\''`) + "'"
}

func unixExec(ctx context.Context, cmdStr string) ([]byte, error) {
	shell := filepath.Base(core.Getenv("SHELL"))
	if shell == "" {
		shell = "sh"
	}

	// the -i flag to ensure that the shell's environment is fully initialised,
	// which allows for correct resolution of tokens that are defined in the user's
	// shell profile. This is necessary on Unix-like systems where users commonly
	// define custom functions and aliases in their shell configuration, and we want
	// Locate to reflect the actual environment as closely as possible. The
	// performance impact of spawning an interactive shell is acceptable
	// in this context since these operations are not on the critical path of
	// execution and provide a more accurate representation of the user's environment.
	//
	// the -c flag is used to pass the command string to the shell for execution.
	// This allows the command to be executed in the context of the shell, which is
	// necessary for proper handling of shell features such as variable expansion,
	// command substitution, and built-in commands.
	cmd := exec.CommandContext(ctx, shell, "-ic", cmdStr) //nolint:gosec // ok
	return cmd.CombinedOutput()
}
