// Package shell detects the runtime shell environment in which jay is
// executing and provides a platform-appropriate LocateFunc for validating
// whether a named executable, builtin, or shell function is invokable
// before a traversal begins.
//
// On Unix-like systems (Linux, macOS, etc.) the detection is trivial -
// jay always runs under a POSIX-compatible shell and command -v is used
// for all lookups.
//
// On Windows the environment is ambiguous: jay may be hosted inside
// PowerShell, cmd.exe, Cygwin, or MSYS2/Git Bash. Each environment has
// a different set of builtins and a different mechanism for locating
// executables and functions. shell.Detect() inspects environment variables
// to determine the correct strategy and returns an Environment whose
// Locate field is ready to use.
package shell
