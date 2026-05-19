# CLAUDE.md - jay

## Project Overview

`jay` is a Go CLI application that acts as a companion tool to the `agenor` directory-walking library. It uses `cobra`/`mamba` for the CLI layer and `viper` for configuration management.

- **Repo**: `github.com/snivilised/jaywalk` (`jay` lives at `src/cmd/jay` within it)
- **Module**: `github.com/snivilised/jaywalk` (jay is the CLI frontend, located at `src/cmd/jay` within the agenor module)
- **Entry point for jay**: `./src/cmd/jay/main.go`

## Build & Test Commands

- **Build**: `go build -o jay ./src/cmd/jay/main.go`
- **Test**: `go test ./...`
- **Dependencies**: `go mod tidy`

## Unit testing

- implement all unit tests using ginkgo and gomega
- when using environment variables in code, make sure this is mockable. Also ensure that unit tests do not interact with the real environment. The real environment should always be isolated from unit tests.

## lint

- run lint using:
golangci-lint run

Fix lint issues. However, do not go overboard in fixing lint issues as some lints issues are ok, such type conversions. If you disable a lint error using nolint, then add an appropriate comment to the nolint statement.
If you are not sure, do not fix the lint issue. I will see the lint error and I will make a judgement call.

## Programming

- assume you are an expert Golang developer.
- you are aware of the go best practices and employ good and well established software engineering techniques
- use Go design patterns when appropriate
- using single responsible philosophy
- use liskov substitution principles
- do not use Go's init() function
- avoid complicated and excessive if statements when these can be replaced by better abstractions
- when using charm's lipgloss and other charm packages, check to see if a component already exists, typically in bubbles, bubbletea or lipgloss. Don't replicate functionality locally, if it already implemented in one of the charm packages.

### os.XXX calls

All code using methods defined in os, eg os.UserHomeDir, should not do so directly. Instead follow the pattern already established in core (see core.core-defs). If there is no definition in core for the os.XXX function call you need to make, please add a new one.

The purpose of this is to make production code testable without giving direct access to os.XXX and allows these to be mocked out.

## Architecture

| Package | Responsibility |
| --- | --- |
| `src/app/bedrock` | Configuration loading; `ViperInstance` is the test-isolation seam for `Load` |
| `cmd/ui` | UI abstraction; `Manager` interface accepts `*core.Node`; `agenor.Servant` is unwrapped at the command layer before crossing into `ui` |
| `cmd/command` | cobra/mamba wiring; `Bootstrap` owns all param-set stashes |

All flags are defined in `src/app/bedrock/flags.go`.

## Viper & Configuration

- Use `viper.GetViper()` to obtain the global viper instance for `bedrock.Load`
- Use `viper.Get()` to access configuration values
- In tests, use the `viperFromYAML` helper (defined in `./cmd/internal/config/helpers_test.go`) for in-memory viper fixtures instead of reading from disk

## agenor Integration

`jay` uses `agenor` (`github.com/snivilised/jaywalk`) as its directory-walking backend. Follow these conventions:

- Construct facades as named variables before passing to `Tortoise`/`Hare` - never inline:

```go
  using := &pref.Using{...}
  relic := &pref.Relic{...}
```

- **Synchronous walk**: `agenor.Tortoise(isPrime)(facade, opts...).Navigate(ctx)`
- **Concurrent run**: `agenor.Hare(isPrime, &wg)(facade, opts...).Navigate(ctx)`, followed by `wg.Wait()`
- Enum values are defined in the `enums` package - use `enums.MetricNoFilesInvoked`, not `agenor.MetricNoFilesInvoked`

## mamba/assist

- Use `NewFlagInfo` for local flags; use `NewFlagInfoOnFlagSet` for persistent flags
- The first word of the `usage` string becomes the flag name
- Use the appropriate BindXXX depending on the type of variable the flag represents

## i18n

- Translation structs are defined in `github.com/snivilised/jaywalk/src/locale`
- Follow the i18n conventions in `GO-USER-CONFIG.md`; locale struct placement is per the package above

### Defining user facing content

When defining content. that is presented to the end user, usually via a print to the console or a log message, then do not use raw strings as that does not support i18n translation, instead make use of the 'lingo' tool and proceed as follows.

- first note that go-i18n is the underlying external dependency that handle i18n. It defines generic structures that we can use to interact with it.
- mamba is an helper for go-i18n that makes using go-i18n a bit easier, handling some of the concerns relating to build command structures and defining flags.
- use of lingo requires adding entries into lingo.Underliers global variable defined in underlying-templ-data
- please consult .claude/skills/go-i18n/SKILL.md

## File References

@./.claude/COMMON-COMMANDS.md
