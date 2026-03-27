# SKILL: go-i18n Template Data Structs

## Purpose

This skill guides a coding agent in converting raw string literals into
properly structured `go-i18n` template data types, following the conventions
established in this codebase. It covers all five message scenarios, file
placement rules, ID naming conventions, emoji markers, and wrapping-error
patterns.

---

## 1. Prerequisites — Locate the Base Embed Struct

Before generating any code, identify the project-specific base embed struct.
It lives in the `locale` package (typically `src/app/locale/`) and follows
the naming pattern `<projectName>TemplData`.
```go
type <projectName>TemplData struct{}

func (td <projectName>TemplData) SourceID() string {
    return SourceID
}
```

Substitute `<projectName>` with the actual project name (e.g. `agenor`,
`jaywalk`). Every template data struct in this codebase embeds this type.
Never define a template data struct without it.

---

## 2. File Placement

| Message kind          | Target file                          |
|-----------------------|--------------------------------------|
| Non-error messages    | `src/app/locale/messages-general.go` |
| Error messages        | `src/app/locale/messages-errors.go`  |

Always append to the appropriate existing file. Do not create new files for
individual messages.

---

## 3. Emoji Markers

Every message block is introduced by a comment containing a canonical emoji
that signals its scenario. Use exactly the emoji shown — do not substitute.

| Emoji | Meaning                               |
|-------|---------------------------------------|
| 🧊    | Non-error message (static or dynamic) |
| ❌    | Static error (no variable content)    |
| 🍒    | Dynamic error (contains variables)    |

The section comment format is:
```go
// 🧊 <HumanReadableName>
```

---

## 4. Message ID Naming Convention

Message IDs are **strictly enforced** based on message type:

| Type             | ID pattern                              | Example                                    |
|------------------|-----------------------------------------|--------------------------------------------|
| Non-error (any)  | `<kebab-slug>`                          | `root-command-short-description`           |
| Static error     | `<kebab-slug>.<project>.static-error`  | `filter-missing-type.agenor.static-error`  |
| Dynamic error    | `<kebab-slug>.<project>.dynamic-error` | `traversal-saved.agenor.dynamic-error`     |

Rules:
- `<kebab-slug>` is derived from the struct name: strip the suffix
  `TemplData`, then convert `PascalCase` → `kebab-case`.
- `<project>` is the project name in lowercase (same as the prefix in the
  base embed struct name).
- Never omit the `.static-error` / `.dynamic-error` suffix for error messages.
- Never add those suffixes to non-error messages.

---

## 5. Scenario Reference

### Scenario 1 — Static Non-Error Message (🧊)

No variable fields. The `Other` string is a plain literal.
```go
// 🧊 <HumanReadableName>

// <StructName>TemplData provides template data for <description>.
type <StructName>TemplData struct {
    <projectName>TemplData
}

func (td <StructName>TemplData) Message() *i18n.Message {
    return &i18n.Message{
        ID:          "<kebab-slug>",
        Description: "<short description>",
        Other:       "<user-facing message>",
    }
}
```

---

### Scenario 2 — Dynamic Non-Error Message (🧊)

One or more exported fields supply variable content. Each field is referenced
in `Other` as `{{.<FieldName>}}`.
```go
// 🧊 <HumanReadableName>

// <StructName>TemplData supplies template data for <description>.
type <StructName>TemplData struct {
    <projectName>TemplData
    // <FieldName> represents <what this field describes>.
    <FieldName> <Type>
}

func (td <StructName>TemplData) Message() *i18n.Message {
    return &i18n.Message{
        ID:          "<kebab-slug>",
        Description: "<short description>",
        Other:       "<message with {{.<FieldName>}} interpolation>",
    }
}
```

Rules:
- Add one exported field per variable token in the message.
- Field names must be exported (capitalised).
- Document each field with a comment.
- Template references must use `{{.<FieldName>}}` — exact capitalisation match.

---

### Scenario 3 — Static Error Message (❌)

No variable fields. Produces a reusable sentinel error variable.
```go
// ❌ <HumanReadableName>

// <StructName>ErrorTemplData is the template data for the <StructName> error.
type <StructName>ErrorTemplData struct {
    <projectName>TemplData
}

// Message creates a new i18n message using the template data.
func (td <StructName>ErrorTemplData) Message() *i18n.Message {
    return &i18n.Message{
        ID:          "<kebab-slug>.<project>.static-error",
        Description: "<short description>",
        Other:       "<user-facing error message>",
    }
}

// <StructName>Error is the error type for <description>.
type <StructName>Error struct {
    li18ngo.LocalisableError
}

// Err<StructName> is the exported sentinel error for <StructName>Error.
var Err<StructName> = <StructName>Error{
    LocalisableError: li18ngo.LocalisableError{
        Data: <StructName>ErrorTemplData{},
    },
}
```

Rules:
- Template data struct is named `<StructName>ErrorTemplData`.
- Error struct is named `<StructName>Error` and embeds `li18ngo.LocalisableError`.
- Sentinel var is named `Err<StructName>`.
- No constructor function is needed.
- No `Error() string` override is needed (inherited from `LocalisableError`).

---

### Scenario 4 — Dynamic Error Message (🍒, no wrapping)

Contains variable content but does **not** wrap another error.
```go
// 🍒 <HumanReadableName>

// <StructName>TemplData is the template data for the <StructName> error.
type <StructName>TemplData struct {
    <projectName>TemplData
    // <FieldName> represents <what this field describes>.
    <FieldName> <Type>
}

// Message creates a new i18n message using the template data.
func (td <StructName>TemplData) Message() *i18n.Message {
    return &i18n.Message{
        ID:          "<kebab-slug>.<project>.dynamic-error",
        Description: "<short description>",
        Other:       "<message with {{.<FieldName>}} interpolation>",
    }
}

// <StructName>Error is the error type for <description>.
type <StructName>Error struct {
    li18ngo.LocalisableError
    <FieldName> <Type>
}
```

Rules:
- Template data fields and error struct fields are **independent**. The error
  struct may carry fields that are not referenced in the template, and vice
  versa.
- No constructor function is required unless wrapping is also needed.

---

### Scenario 5 — Dynamic Error that Wraps Another Error (🍒)

Extends Scenario 4 with a `Wrapped error` field, `Error() string`,
`Unwrap() error`, and a constructor.
```go
// 🍒 <HumanReadableName>

// <StructName>TemplData is the template data for the <StructName> error.
type <StructName>TemplData struct {
    <projectName>TemplData
    // (add exported fields here if the template uses variables)
}

// Message creates a new i18n message using the template data.
func (td <StructName>TemplData) Message() *i18n.Message {
    return &i18n.Message{
        ID:          "<kebab-slug>.<project>.dynamic-error",
        Description: "<short description>",
        Other:       "<message, optionally with {{.<FieldName>}} tokens>",
    }
}

// <StructName>Error is the error type for <description>.
type <StructName>Error struct {
    li18ngo.LocalisableError
    Wrapped     error
    <FieldName> <Type>  // additional fields as required
}

// Error returns the combined wrapped and localised error message.
func (e <StructName>Error) Error() string {
    return fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))
}

// Unwrap returns the wrapped error.
func (e <StructName>Error) Unwrap() error {
    return e.Wrapped
}

// New<StructName>Error creates a new <StructName>Error wrapping <wrappedErrVar>.
func New<StructName>Error(<params>) error {
    return &<StructName>Error{
        LocalisableError: li18ngo.LocalisableError{
            Data: <StructName>TemplData{},
        },
        Wrapped:     <wrappedErrVar>,
        <FieldName>: <value>,
    }
}
```

Rules:
- `Wrapped error` field must always be present when another error is wrapped.
- `Unwrap() error` is **always** required when `Wrapped error` is present.
- `Error() string` must always use `li18ngo.Text(e.Data)` — no other formatter.
- `fmt.Sprintf("%v, %v", e.Wrapped.Error(), li18ngo.Text(e.Data))` is the
  canonical format string; do not deviate from it.
- A `New<StructName>Error(...)` constructor is **required** for wrapping errors
  and **only** for wrapping errors.
- Template data struct fields and error struct fields are independent — the
  error struct may hold fields (including `Destination`, context values, etc.)
  that are not referenced in the template `Other` string.

---

## 6. Decision Tree

Use this to classify a string literal before generating code:
```
Is the message an error?
├── No  → 🧊  Does Other contain {{.Field}} tokens?
│           ├── No  → Scenario 1 (static non-error)
│           └── Yes → Scenario 2 (dynamic non-error)
└── Yes → Does Other contain {{.Field}} tokens?
          ├── No  → ❌ Scenario 3 (static error)
          └── Yes → 🍒 Does it wrap another error?
                    ├── No  → Scenario 4 (dynamic error, no wrap)
                    └── Yes → Scenario 5 (dynamic error, wrapping)
```

---

## 7. Common Mistakes to Avoid

| Mistake | Correct behaviour |
|---|---|
| Using `.static-error` / `.dynamic-error` suffix on a non-error message ID | Non-error IDs are plain kebab slugs only |
| Omitting the suffix on an error message ID | Always append `.<project>.static-error` or `.<project>.dynamic-error` |
| Adding a constructor for a static error | Constructors are only for wrapping errors (Scenario 5) |
| Omitting `Unwrap()` when `Wrapped error` is present | Always add `Unwrap()` when the field exists |
| Mirroring error struct fields into template data struct | They are independent — only add template data fields that appear in `Other` |
| Placing non-error messages in `messages-errors.go` | Non-error → `messages-general.go`; errors → `messages-errors.go` |
| Using the wrong emoji for the section comment | 🧊 non-error · ❌ static error · 🍒 dynamic error |
| Forgetting to embed `<projectName>TemplData` | Every template data struct must embed the base struct |

---

## 8. Imports Required
```go
// messages-general.go and messages-errors.go typically require:
import (
    "fmt"                          // only in files with Error() string methods

    "github.com/nicksnyder/go-i18n/v2/i18n"
    "github.com/snivilised/li18ngo"
)
```

Confirm the actual import paths against `go.mod` before writing — paths vary
slightly between projects.
