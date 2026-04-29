# Jaywalk Architecture

This document records architectural decisions, constraints, and design
principles for the jaywalk project. It exists to preserve intent across
development sessions and to guide contributors toward consistent choices.

---

## Package Layout

```text
src/
  agenor/         - file system traversal (reusable library)
  app/
    command/      - cobra command definitions, thin adapters only
    controller/   - ApplicationController, coordinates layers
    dispatch/     - core execution logic, observer interfaces
  prism/          - terminal rendering (reusable library, peer of agenor)
```

### Dependency Rule

Dependencies flow strictly downward. No package may import a package above
it in this hierarchy:

```text
command -> controller -> dispatch -> prism -> (terminal)
```

`agenor` and `prism` are reusable libraries. They must never import
`command`, `controller`, `dispatch`, or any other jaywalk-specific package.

---

## Prism - Terminal Rendering Library

`prism` is a reusable rendering library importable by any CLI project at
`github.com/snivilised/jaywalk/src/prism`. It has no knowledge of jaywalk
internals. All jaywalk-specific translation happens in an adapter inside
the controller layer.

### Public API

```go
type Renderer interface {
    Begin(overture Overture)
    Show(motif Motif)
    End(summary Summary)
}

func New(kind ViewKind, w io.Writer) Renderer
```

### Core Types

| Type | Role |
| --- | --- |
| `Overture` | Metadata passed to `Begin` - root path, caption, start time, |
| | navigation kind, survey results |
| `Motif` | Single render-able item passed to `Show` per traversal node |
| `Summary` | Traversal result passed to `End` - counts, elapsed, errors |
| `ViewKind` | Typed string identifying the view - open set, not an iota |
| `Theme` | All lipgloss styles, constructed by `NewTheme`, never globals |

### ViewKind

`ViewKind` is a typed string (`type ViewKind string`) rather than an iota
enum because the set of views is open - new views will be added without
modifying the core type.

```go
const (
    StreamView   ViewKind = "stream"
    PortholeView ViewKind = "porthole"
    LanesView    ViewKind = "lanes"
)
```

### NavigationKind

```go
const (
    PrimeNavigation  NavigationKind = "prime"
    ResumeNavigation NavigationKind = "resume"
)
```

`Overture.Kind` signals to the renderer whether this is a fresh traversal
or a continuation from a checkpoint. The banner and summary labels adapt
accordingly.

### Survey Results

A two-phase navigation may precede execution:

- Phase 1 (survey) - full traversal with no invocation; counts selected
  nodes and records maximum depth. This phase is entirely internal to
  `controller/dispatch` and is never surfaced to prism directly.
- Phase 2 (execute) - full traversal with invocation, using survey results.

Survey results are passed to prism via `Overture.Survey`:

```go
type SurveyResult struct {
    NodeCount uint  // total nodes to be visited; enables accurate progress
    MaxDepth  uint  // deepest level seen; used for layout calculations
}

type Overture struct {
    Root       string
    Caption    string
    StartedAt  time.Time
    Kind       NavigationKind
    ResumeFrom string         // populated only when Kind == ResumeNavigation
    Survey     *SurveyResult  // nil when no survey was performed (single-phase)
}
```

`Survey` is a pointer so that `nil` unambiguously means no survey was
performed, avoiding zero-value ambiguity.

### Progress Ownership

Progress calculation is owned by prism, not the caller. When
`Overture.Survey` is non-nil, prism derives progress internally by
incrementing a counter on each `Show` call against the known
`Survey.NodeCount`. Callers do not compute or pass progress values.

This design reduces caller burden, which is especially important for
third-party consumers of prism.

### Depth

`Motif.Depth` is sourced from `node.Extension.Depth` in agenor. Prism
does not track or compute depth itself.

### Colour and Theme

- lipgloss v2 (`github.com/charmbracelet/lipgloss/v2`) is the styling
  dependency.
- Dark/light detection uses `lipgloss.HasDarkBackground(os.Stdin, os.Stdout)`.
- All styles are fields on `Theme`, constructed in `NewTheme`. No
  package-level style variables exist anywhere in prism.
- Colour downsampling and TTY detection (including `NO_COLOR` and piped
  output) are handled automatically by `lipgloss.Fprintln` and friends.

### Views

#### Stream (linear)

- Single-phase navigation only.
- Output written immediately on each `Show` call - no buffering.
- No progress bar; progress tracking not required.
- `Overture.Survey` will be `nil` for this view.

#### Porthole (planned)

- Static header and footer with vertically scrolling content between them.
- Requires bubbletea.
- Progress displayed in footer using `Survey.NodeCount`.
- `Survey.MaxDepth` may inform layout.

#### Lanes (planned)

- Parallel horizontal lanes, one per worker.
- Suited to concurrent execution (jay run) but not exclusive to it.
- Single-lane mode available for jay walk.
- Requires bubbletea.
- Progress displayed per lane using `Survey.NodeCount`.
- `Survey.MaxDepth` informs lane column width.

### Tree Glyphs (planned)

Indentation currently uses plain spaces. Tree-branch glyphs (`├──`, `└──`)
are deferred until sibling tracking is available from agenor. NerdFont
glyph variants are also planned, selectable based on terminal font
capability.

---

## Navigation

### Prime Navigation

Standard full traversal from the root. `Overture.Kind == PrimeNavigation`.

### Resume Navigation

Traversal continues from a previously saved checkpoint.
`Overture.Kind == ResumeNavigation` and `Overture.ResumeFrom` is populated
with the resume path.

### Two-Phase Navigation

An optional survey phase precedes execution to enable accurate progress
reporting. The survey is entirely internal to `controller/dispatch` -
prism is not invoked during the survey phase. Survey results reach prism
only via `Overture.Survey` at the start of the execute phase.

---

## Command Layer

### Principles

- Command handlers are thin adapters: parse flags, build a typed request,
  call the controller. No business logic.
- `Bootstrap` is a pure composition root - wiring only.
- `buildOptions` lives in `src/app/command/options.go`, not in individual
  command files.

### Requests

```text
src/app/controller/requests.go
```

`PrimeRequest` and `ResumeRequest` are defined here. Command handlers
build one of these and pass it to the controller.

---

## Controller

`ApplicationController` owns the UI (prism renderer) and shared state. It
translates between the command layer's requests and the dispatch layer's
operations. It is also responsible for:

- Constructing the prism `Renderer` via `prism.New`.
- Populating `Overture` (including `SurveyResult` when applicable) from
  request data and survey phase output.
- Translating agenor node events into `prism.Motif` values (the adapter).

---

## Dispatch

Core execution logic. Defines its own observer interfaces. No upward
dependencies - dispatch does not import controller or command.

The survey phase runs entirely within dispatch. Its output (`NodeCount`,
`MaxDepth`) is returned to the controller, which packages it into
`SurveyResult` before calling `prism.New`.

---

## General Principles

- Package names describe the service they provide, not their contents.
- Strict one-way dependency rule enforced at all times.
- Domain defines interfaces; outer layers satisfy them.
- No global variables - everything flows through constructed types.
- `errors.Is` / `errors.As` used directly by callers; no generated helper
  wrappers.
- Em-dashes must not appear in code comments - use a regular dash (`-`).
