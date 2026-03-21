# Jaywalk / Jay - Project Context

## Module Identity

- Module: `github.com/snivilised/jaywalk`
- CLI binary: `jay`
- Organisation: `github.com/snivilised`
- Navigation library: `agenor` (package name), lives at `src/agenor`

## What This Project Is

Jaywalk is both a library and a CLI application:

- **agenor** - the core navigation/traversal library (public API)
- **jay** - a CLI tool built on top of agenor, entry point at `cmd/jay/main.go`

## Directory Structure

```txt
github.com/snivilised/jaywalk/
├── cmd/
│   └── jay/
│       └── main.go                 # entry point - wires and launches only
├── src/
│   ├── agenor/                     # navigation library - core traversal services
│   │   ├── internal/               # private to agenor
│   │   │   ├── enclave/
│   │   │   ├── feat/
│   │   │   ├── filtering/
│   │   │   ├── kernel/
│   │   │   ├── laboratory/         # test helpers - private to agenor
│   │   │   ├── level/
│   │   │   ├── opts/
│   │   │   └── persist/
│   │   ├── test/                   # test programs
│   │   ├── collections/
│   │   ├── core/
│   │   ├── enums/
│   │   ├── life/
│   │   ├── pref/
│   │   ├── stock/
│   │   ├── tapable/
│   │   └── tfs/
│   ├── app/
│   │   ├── command/                # CLI adapter layer - cobra/mamba
│   │   │   └── internal/
│   │   │       └── cfg/            # config - internal to command
│   │   ├── controller/             # mediator/coordination layer
│   │   ├── dispatch/               # domain layer - core jay logic
│   │   └── ui/                     # UI concerns
│   └── internal/                   # visible to all of src/
│       ├── third/                  # third party utilities
│       └── services/               # shared services
├── locale/
├── src/examples/
├── scripts/
└── resources/
```

## Three Layer Architecture

The jay CLI follows a strict three-layer architecture with a one-way dependency rule:

```txt
cmd/jay/main.go
  → src/app/command      (CLI adapter - cobra/mamba, thin command handlers)
    → src/app/controller (mediator/coordinator - owns UI, coordinates layers)
      → src/app/dispatch (domain - core jay execution logic)
      → src/agenor       (navigation services)
```

**The dependency rule is absolute - nothing in dispatch or agenor imports from app or cmd.**

## Package Layer Responsibilities

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

## Design Principles

- Single Responsibility Principle
- Dependency Inversion - domain defines interfaces, outer layers satisfy them
- Ports and Adapters (Hexagonal Architecture) - strict layering
- Low coupling, high cohesion
- Package names describe the service they provide, not what they contain
- Go internal visibility rules used to enforce layer boundaries

## Key Naming Conventions

- Command package: `package command`
- Controller package: `package controller`
- Dispatch package: `package dispatch`
- UI package: `package ui`

## Template Lineage

- Jaywalk was created from `astrolib` (library template) with jay CLI inserted
- Templates: `arcadia` (CLI apps), `astrolib` (libraries) - plan to consolidate into one
