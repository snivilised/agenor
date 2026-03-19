# CLAUDE.md - agenor

## Project Overview

`agenor` is a file system traversal library that navigates directory trees and notifies
callers of events at each node. It extends the standard `filepath.Walk` with: regex/glob
filtering, hibernation (deferred activation of callbacks until a condition is met),
resume from a previously interrupted session, concurrent navigation via pants worker pool,
and hook-able traversal behaviour.

- **Module**: `github.com/snivilised/agenor`
- **Package alias**: `age` (import as `age "github.com/snivilised/agenor"`)
- **Docs**: <https://pkg.go.dev/github.com/snivilised/agenor>

## Build & Test Commands

- **Test all**: `go test ./...`
- **Dependencies**: `go mod tidy`

## Package Architecture

The dependency rule is: packages may only depend on packages in layers below them.
This rule does not apply to unit tests.

```txt
🔆 user interface layer
  age (root package)      - public API; may use everything

🔆 feature layer
  internal/feat/resume    - depends on pref, opts, kernel
  internal/feat/sampling  - depends on filter
  internal/feat/hiber     - depends on filter, services
  internal/feat/filter    - no internal deps

🔆 central layer
  internal/kernel         - no internal deps
  internal/enclave        - depends on pref, override
  internal/opts           - depends on pref
  internal/override       - depends on tapable; must not use enclave

🔆 support layer
  pref                    - depends on life, services, persist
  internal/persist        - no internal deps
  internal/services       - no internal deps

🔆 intermediary layer
  life                    - no internal deps; must not use pref

🔆 platform layer
  tapable                 - depends on core
  core                    - no internal deps
  enums                   - no deps
  tfs                     - no internal deps
```

## Core API

### Traversal modes

There are two traversal modes and two extents, giving four possible scenarios:

| Mode | Extent | Description |
| --- | --- | --- |
| Walk | Prime | Sequential traversal from root |
| Walk | Resume | Sequential traversal resuming from a saved session |
| Run | Prime | Concurrent traversal from root |
| Run | Resume | Concurrent traversal resuming from a saved session |

The low-level API composes these explicitly:

```go
// Walk/Prime
age.Walk().Configure().Extent(age.Prime(facade, opts...)).Navigate(ctx)

// Run/Resume
age.Run(wg).Configure().Extent(age.Resume(facade, opts...)).Navigate(ctx)
```

### Scenario composites

To avoid conditional duplication at the call site, use the scenario composites:

| Composite | Fixes | Selects by |
| --- | --- | --- |
| `Tortoise(isPrime)` | Walk | `isPrime bool` → Prime or Resume |
| `Hare(isPrime, wg)` | Run | `isPrime bool` → Prime or Resume |
| `Goldfish(isWalk, wg)` | Prime | `isWalk bool` → Walk or Run |
| `Hydra(isWalk, isPrime, wg)` | neither | both `isWalk` and `isPrime` |

Usage pattern - always pass `isPrime`/`isWalk` as named `const bool` values to avoid
lint warnings from bare literals:

```go
const isPrime = true
age.Tortoise(isPrime)(facade, opts...).Navigate(ctx)

var wg sync.WaitGroup
age.Hare(isPrime, &wg)(facade, opts...).Navigate(ctx)
wg.Wait()
```

### Facades

Construct facades as named variables, never inline:

```go
using := &pref.Using{...}
relic := &pref.Relic{...}   // resume sessions only
```

- `pref.Using` - dependencies for a Prime session
- `pref.Relic` - saved state for a Resume session

### Enums

All enum values are in the `enums` package. Do not use `age.` prefixed aliases
for enum values - use `enums.` directly:

```go
enums.SubscribeFiles            // not age.SubscribeFiles
enums.MetricNoFilesInvoked      // not age.MetricNoFilesInvoked
enums.ResumeStrategyFastward
```

### Options (With* functions)

Options are passed as variadic `...pref.Option` to `Prime`/`Resume` or to a composite.
All `With*` option constructors are re-exported from the root `age` package:

```go
age.WithFilter(...)
age.WithDepth(5)
age.WithOnBegin(handler)
age.WithCPU              // use all available CPUs for Run
age.WithNoW(n)           // use n workers for Run
```

Use `age.IfOption` / `age.IfOptionF` / `age.IfElseOptionF` for conditional options.

## Key Types

| Type | Package | Purpose |
| --- | --- | --- |
| `age.Node` | `core` | A file system node passed to the client callback |
| `age.Servant` | `core` | Provides the client with traversal properties |
| `age.Client` | `core` | The callback signature: `func(node *age.Node) error` |
| `age.Navigator` | `core` | Returned by `Extent()`; call `.Navigate(ctx)` on it |
| `age.Options` | `pref` | Full options struct available inside `With*` constructors |
| `age.Using` | `pref` | Alias for `pref.Using` (Prime facade) |
| `age.Relic` | `pref` | Alias for `pref.Relic` (Resume facade) |
| `age.TraversalFS` | `tfs` | File system interface required for traversal |

## Internal Packages (do not import directly)

- `internal/kernel` - core traversal engine
- `internal/feat/*` - feature plugins (filter, hiber, resume, sampling, nanny)
- `internal/enclave` - supervisor and kernel result types
- `internal/opts` - options loading and binding
- `internal/persist` - session state marshalling for resume
- `internal/services` - cross-cutting concerns (message bus)
- `internal/filtering` - shared filter implementations used by multiple plugins
- `internal/laboratory` - internal test helpers (not for external use)

## Test Helpers

- **`test/hanno`** (`github.com/snivilised/agenor/test/hanno`) - utilities for building
  virtual file system trees; see `GO-USER-CONFIG.md` for `Nuxx` usage
- **`test/data/musico-index.xml`** - standard XML fixture representing a sample music
  directory tree, used by `Nuxx` to populate an in-memory file system
- **`internal/laboratory`** - internal-only test utilities; do not use from outside the module

## i18n

- Translation structs are defined in `github.com/snivilised/agenor/locale`
- Follow the i18n conventions in `GO-USER-CONFIG.md`

## File References

@~/.claude/GO-USER-CONFIG.md
