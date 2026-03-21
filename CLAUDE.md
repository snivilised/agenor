# CLAUDE.md - jaywalk

## Project Overview

`jaywalk` is a file system traversal library that navigates directory trees and notifies
callers of events at each node. It extends the standard `filepath.Walk` with: regex/glob
filtering, hibernation (deferred activation of callbacks until a condition is met),
resume from a previously interrupted session, concurrent navigation via pants worker pool,
and hook-able traversal behaviour.

- **Module**: `github.com/snivilised/jaywalk`
- **Docs**: <https://pkg.go.dev/github.com/snivilised/jaywalk>

## Build & Test Commands

- **Test all**: `go test ./...`
- **Dependencies**: `go mod tidy`

## Package Architecture

The dependency rule is: packages may only depend on packages in layers below them.
This rule does not apply to unit tests.

```txt
🔆 user interface layer
  src/agenor (root package)                - public API; may use everything

🔆 feature layer
  src/agenor/internal/feat/resume          - depends on pref, opts, kernel
  src/agenor/internal/feat/sampling        - depends on filter
  src/agenor/internal/feat/hiber           - depends on filter, services
  src/agenor/internal/feat/filter          - no internal deps

🔆 central layer
  src/agenor/internal/kernel               - no internal deps
  src/agenor/internal/enclave              - depends on pref, override
  src/agenor/internal/opts                 - depends on pref
  src/agenor/internal/override             - depends on tapable; must not use enclave

🔆 support layer
  src/agenor/pref                          - depends on life, services, persist
  src/agenor/internal/persist              - no internal deps
  src/internal/services                    - no internal deps

🔆 intermediary layer
  src/agenor/life                          - no internal deps; must not use pref

🔆 platform layer
  src/agenor/tapable                       - depends on core
  src/agenor/core                          - no internal deps
  src/agenor/enums                         - no deps
  src/agenor/tfs                           - no internal deps

🔆 shared internal layer
  src/internal/third                       - no internal deps; visible to all of src/
  src/internal/services                    - no internal deps; visible to all of src/

🔆 cli layer
  src/app/command                          - depends on src/app/controller
  src/app/controller                       - depends on src/app/dispatch, src/agenor
  src/app/dispatch                         - depends on src/agenor; no upward deps
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
agenor.Walk().Configure().Extent(agenor.Prime(facade, opts...)).Navigate(ctx)

// Run/Resume
agenor.Run(wg).Configure().Extent(agenor.Resume(facade, opts...)).Navigate(ctx)
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
agenor.Tortoise(isPrime)(facade, opts...).Navigate(ctx)

var wg sync.WaitGroup
agenor.Hare(isPrime, &wg)(facade, opts...).Navigate(ctx)
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

All enum values are in the `enums` package. Do not use `agenor.` prefixed aliases
for enum values - use `enums.` directly:

```go
enums.SubscribeFiles            // not agenor.SubscribeFiles
enums.MetricNoFilesInvoked      // not agenor.MetricNoFilesInvoked
enums.ResumeStrategyFastward
```

### Options (With* functions)

Options are passed as variadic `...pref.Option` to `Prime`/`Resume` or to a composite.
All `With*` option constructors are re-exported from the root `agenor` package:

```go
agenor.WithFilter(...)
agenor.WithDepth(5)
agenor.WithOnBegin(handler)
agenor.WithCPU              // use all available CPUs for Run
agenor.WithNoW(n)           // use n workers for Run
```

Use `agenor.IfOption` / `agenor.IfOptionF` / `agenor.IfElseOptionF` for conditional options.

## Key Types

| Type | Package | Purpose |
| --- | --- | --- |
| `agenor.Node` | `core` | A file system node passed to the client callback |
| `agenor.Servant` | `core` | Provides the client with traversal properties |
| `agenor.Client` | `core` | The callback signature: `func(node *agenor.Node) error` |
| `agenor.Navigator` | `core` | Returned by `Extent()`; call `.Navigate(ctx)` on it |
| `agenor.Options` | `pref` | Full options struct available inside `With*` constructors |
| `agenor.Using` | `pref` | Alias for `pref.Using` (Prime facade) |
| `agenor.Relic` | `pref` | Alias for `pref.Relic` (Resume facade) |
| `agenor.TraversalFS` | `tfs` | File system interface required for traversal |

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

- **`test/hanno`** (`github.com/snivilised/jaywalk/test/hanno`) - utilities for building
  virtual file system trees; see `GO-USER-CONFIG.md` for `Nuxx` usage
- **`test/data/musico-index.xml`** - standard XML fixture representing a sample music
  directory tree, used by `Nuxx` to populate an in-memory file system
- **`internal/laboratory`** - internal-only test utilities; do not use from outside the module

## i18n

- Translation structs are defined in `github.com/snivilised/jaywalk/locale`
- Follow the i18n conventions in `GO-USER-CONFIG.md`

## File References

@~/.claude/GO-USER-CONFIG.md
