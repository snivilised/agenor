# CLAUDE.md - jay

## Project Overview

`jay` is a Go CLI application that acts as a companion tool to the `agenor` directory-walking library. It uses `cobra`/`mamba` for the CLI layer and `viper` for configuration management.

- **Repo**: `github.com/snivilised/jaywalk` (`jay` lives at `cmd/jay` within it)
- **Module**: `github.com/snivilised/jaywalk` (jay is the CLI frontend, located at `cmd/jay` within the agenor module)
- **Entry point for jay**: `./cmd/jay/main.go`

## Build & Test Commands

- **Build**: `go build -o jay ./cmd/jay/main.go`
- **Test**: `go test ./...`
- **Dependencies**: `go mod tidy`

## Architecture

| Package | Responsibility |
| --- | --- |
| `cmd/internal/cfg` | Configuration loading; `ViperInstance` is the test-isolation seam for `Load` |
| `cmd/ui` | UI abstraction; `Manager` interface accepts `*core.Node`; `age.Servant` is unwrapped at the command layer before crossing into `ui` |
| `cmd/command` | cobra/mamba wiring; `Bootstrap` owns all param-set stashes |

All flags are defined in `cmd/internal/cfg/flags.go`.

## Viper & Configuration

- Use `viper.GetViper()` to obtain the global viper instance for `cfg.Load`
- Use `viper.Get()` to access configuration values
- In tests, use the `viperFromYAML` helper (defined in `./cmd/internal/config/helpers_test.go`) for in-memory viper fixtures instead of reading from disk

## agenor Integration

`jay` uses `agenor` (`github.com/snivilised/jaywalk`) as its directory-walking backend. Follow these conventions:

- Construct facades as named variables before passing to `Tortoise`/`Hare` - never inline:

```go
  using := &pref.Using{...}
  relic := &pref.Relic{...}
```

- **Synchronous walk**: `age.Tortoise(isPrime)(facade, opts...).Navigate(ctx)`
- **Concurrent run**: `age.Hare(isPrime, &wg)(facade, opts...).Navigate(ctx)`, followed by `wg.Wait()`
- Enum values are defined in the `enums` package - use `enums.MetricNoFilesInvoked`, not `age.MetricNoFilesInvoked`

## mamba/assist

- Use `NewFlagInfo` for local flags; use `NewFlagInfoOnFlagSet` for persistent flags
- The first word of the `usage` string becomes the flag name
- Use `BindString`, not `BindValidatedString`, unless validation is explicitly required

## i18n

- Translation structs are defined in `github.com/snivilised/jaywalk/locale`
- Follow the i18n conventions in `GO-USER-CONFIG.md`; locale struct placement is per the package above

## File References

@./.claude/COMMON-COMMANDS.md
@~/.claude/GO-USER-CONFIG.md
