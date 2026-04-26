package enums

//go:generate stringer -type=ShellKind -linecomment -trimprefix=ShellKind -output shell-kind-en-auto.go

// ShellKind identifies the shell environment in which jay is executing.
type ShellKind uint

const (
	// ShellKindNativeUnix covers Linux, macOS, and any other POSIX platform.
	// command -v is used for all lookups.
	ShellKindNativeUnix ShellKind = iota // native-unix

	// ShellKindPowerShell covers Windows hosting jay inside PowerShell 5.1
	// (powershell.exe) or PowerShell 7+ (pwsh.exe).
	// Get-Command is used for all lookups.
	ShellKindPowerShell // powershell

	// ShellKindCmdExe covers Windows hosting jay inside the legacy cmd.exe
	// console. where.exe is used for binary lookups; a hardcoded builtin
	// set covers cmd.exe internal commands.
	ShellKindCmdExe // cmd.exe

	// ShellKindCygwin covers Windows hosting jay inside a Cygwin bash session.
	// /bin/sh -c "command -v" is used for all lookups, identical to Unix.
	ShellKindCygwin // cygwin

	// ShellKindMSYS2 covers Windows hosting jay inside MSYS2 or Git Bash.
	// /bin/sh -c "command -v" is used for all lookups, identical to Unix.
	ShellKindMSYS2 // msys2
)
