//go:build windows

package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/snivilised/jaywalk/src/agenor/enums"
	"github.com/snivilised/jaywalk/src/locale"
)

// detect inspects Windows environment variables to determine which shell
// is hosting jay and returns the appropriate Environment.
//
// Detection priority (first match wins):
//  1. CYGWIN env var set      -> Cygwin bash   -> command -v via /bin/sh
//  2. MSYSTEM env var set     -> MSYS2/Git Bash -> command -v via /bin/sh
//  3. PSModulePath env var set -> PowerShell    -> Get-Command via pwsh/powershell
//  4. fallback                 -> cmd.exe       -> where.exe + builtin set
func detect() (Environment, error) {
	if os.Getenv("CYGWIN") != "" {
		return Environment{
			Kind:   enums.ShellKindCygwin,
			Locate: unixStyleLocate,
		}, nil
	}

	if os.Getenv("MSYSTEM") != "" {
		return Environment{
			Kind:   enums.ShellKindMSYS2,
			Locate: unixStyleLocate,
		}, nil
	}

	if os.Getenv("PSModulePath") != "" {
		psExe, err := resolvePowerShellExe()
		if err != nil {
			return Environment{}, locale.NewPSMSetNoPowerShellExeFoundError(err)
		}

		return Environment{
			Kind:   enums.ShellKindPowerShell,
			Locate: powerShellLocate(psExe),
		}, nil
	}

	// Fallback: assume cmd.exe.
	return Environment{
		Kind:   enums.ShellKindCmdExe,
		Locate: cmdExeLocate,
	}, nil
}

// ---------------------------------------------------------------------------
// Cygwin / MSYS2 - delegate to POSIX command -v via /bin/sh
// ---------------------------------------------------------------------------

// unixStyleLocate is used for Cygwin and MSYS2 environments, both of
// which provide a fully functional /bin/sh with POSIX command -v.
func unixStyleLocate(name string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", "command -v "+shellQuote(name))
	out, err := cmd.Output()

	if err != nil {
		return "", locale.NewCmdNotFoundInEnvError(name)
	}

	return strings.TrimSpace(string(out)), nil
}

// shellQuote wraps name in single quotes for safe embedding in a -c argument.
func shellQuote(name string) string {
	return "'" + strings.ReplaceAll(name, "'", `'\''`) + "'"
}

// ---------------------------------------------------------------------------
// PowerShell
// ---------------------------------------------------------------------------

// resolvePowerShellExe returns the path to the best available PowerShell
// executable. pwsh.exe (PowerShell 7+) is preferred over the legacy
// powershell.exe (Windows PowerShell 5.1).
func resolvePowerShellExe() (string, error) {
	if path, err := exec.LookPath("pwsh.exe"); err == nil {
		return path, nil
	}

	if path, err := exec.LookPath("powershell.exe"); err == nil {
		return path, nil
	}

	return "", locale.ErrNeitherPwshOrPowerShellExeFound
}

// powerShellLocate returns a LocateFunc that uses Get-Command to check
// whether a token is invokable inside PowerShell. The profile is loaded
// (no -NoProfile flag) so that user-defined functions are visible.
// psExe must be the resolved path to pwsh.exe or powershell.exe.
func powerShellLocate(psExe string) LocateFunc {
	return func(name string) (string, error) {
		// Get-Command writes the command info object to stdout on success
		// and throws a terminating error (exit code 1) when not found.
		// Piping to Out-Null suppresses output; we care only about the
		// exit code.
		script := fmt.Sprintf("Get-Command %s | Out-Null", name)

		cmd := exec.Command(psExe, "-Command", script)
		if err := cmd.Run(); err != nil {
			return "", locale.NewCmdNotFoundInEnvError(name)
		}

		return name, nil
	}
}

// ---------------------------------------------------------------------------
// cmd.exe
// ---------------------------------------------------------------------------

// cmdExeLocate implements LocateFunc for cmd.exe environments. It tries
// where.exe first (covers binaries on PATH), then falls back to the
// hardcoded set of cmd.exe internal commands.
func cmdExeLocate(name string) (string, error) {
	// where.exe resolves binaries on PATH, respecting PATHEXT.
	cmd := exec.Command("where.exe", name)
	out, err := cmd.Output()

	if err == nil {
		return strings.TrimSpace(string(out)), nil
	}

	// where.exe failed - check the hardcoded cmd.exe builtin set.
	lower := strings.ToLower(name)
	if _, ok := cmdExeBuiltins[lower]; ok {
		return name, nil
	}

	return "", locale.NewCmdNotFoundAsPathBinaryOrBuiltinError(name)
}

// cmdExeBuiltins is the complete set of commands that are internal to
// cmd.exe and therefore not resolvable by where.exe. This list is stable;
// Microsoft has not added internal commands since Windows XP.
//
// Reference:
// https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/windows-commands
var cmdExeBuiltins = map[string]struct{}{
	"assoc":    {},
	"break":    {},
	"call":     {},
	"cd":       {},
	"chdir":    {},
	"cls":      {},
	"color":    {},
	"copy":     {},
	"date":     {},
	"del":      {},
	"dir":      {},
	"echo":     {},
	"endlocal": {},
	"erase":    {},
	"exit":     {},
	"for":      {},
	"ftype":    {},
	"goto":     {},
	"if":       {},
	"md":       {},
	"mkdir":    {},
	"mklink":   {},
	"move":     {},
	"path":     {},
	"pause":    {},
	"popd":     {},
	"prompt":   {},
	"pushd":    {},
	"rd":       {},
	"rem":      {},
	"ren":      {},
	"rename":   {},
	"rmdir":    {},
	"set":      {},
	"setlocal": {},
	"shift":    {},
	"start":    {},
	"time":     {},
	"title":    {},
	"type":     {},
	"ver":      {},
	"verify":   {},
	"vol":      {},
}
