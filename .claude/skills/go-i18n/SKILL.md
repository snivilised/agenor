---
name: go-i18n
description: Guides a coding agent in adding new i18n messages and errors using li18ngo and lingo. Covers the Underliers map input format, UnderlyingType enum values, placeholder patterns for un-generated messages, and the programming error exemption.
---

# SKILL: Adding i18n Messages via lingo

## Purpose

This skill guides a coding agent in adding new i18n messages and errors to a
project that uses li18ngo and go-i18n. All user-facing text and translatable
errors are defined by adding entries to the `Underliers` map in
`underlying-templ-data.go`, then running lingo to regenerate the auto files.

**Never hand-write the generated output.** The only files a coding agent
should touch are:
- `underlying-templ-data.go` - add new `UnderlyingTemplData` entries here
- Run lingo: `go run ./cmd/lingo/...`

---

## 1. The Two Key Types

### `UnderlyingTemplData` - one entry per message

```go
type UnderlyingTemplData struct {
    MessageID   string                // unique go-i18n message ID
    Seed        string                // PascalCase base name for all generated identifiers
    TypeName    enums.UnderlyingType  // controls what code is generated
    Description string                // short human summary; used as struct doc comment
    Story       string                // longer narrative; word-wrapped at 80 chars in banner
    Other       string                // go-i18n translation string; may contain {{.Token}}
    Fields      []UnderlyingField     // variable fields for dynamic messages; empty for static
    File        string                // optional output file prefix; leave empty for defaults
}
```

### `UnderlyingField` - one entry per variable token

```go
type UnderlyingField struct {
    Note   string // must match a {{.Note}} token in Other exactly
    GoType string // valid Go type: "string", "int", "uint", "error"
    Tale   string // doc comment for the generated field; 🔥 TODO if empty
}
```

---

## 2. The `UnderlyingType` Enum - Full Reference

Choose the `TypeName` that matches what you need:

| TypeName | Use when |
|---|---|
| `UnderlyingTypeStaticCobra` | Static cobra command/flag description, no variables |
| `UnderlyingTypeDynamicCobra` | Cobra description with `{{.Token}}` variables |
| `UnderlyingTypeStaticGeneral` | Static non-error user-facing message |
| `UnderlyingTypeDynamicGeneral` | Non-error message with `{{.Token}}` variables |
| `UnderlyingTypeStaticError` | Static error, no variables, no wrapping |
| `UnderlyingTypeSentinelError` | Static sentinel error designed to be wrapped by outer errors; generates `ErrXxx` |
| `UnderlyingTypeStaticErrorWrapper` | Static error that wraps another for the error chain only; wrapped message does NOT appear in translated output |
| `UnderlyingTypeStaticErrorWrapperMsg` | Static error that wraps another AND shows `{{.Wrapped}}` in the translated output |
| `UnderlyingTypeDynamicError` | Dynamic error with variables, no wrapping |
| `UnderlyingTypeDynamicErrorWrapper` | Dynamic error with variables that also wraps another error |

Import path for the enum:
```go
import "github.com/snivilised/li18ngo/locale/enums"
```

---

## 3. MessageID and Seed Conventions

- `MessageID` must be unique across the entire `Underliers` map.
- For non-errors: use `kebab-case` slug, e.g. `"using-config-file"`.
- For errors: append `.dynamic-error` or `.static-error` (or the wrapper
  variant), e.g. `"path-not-found.dynamic-error"`.
- `Seed` is `PascalCase` derived from the human name, e.g. `"PathNotFound"`.
  lingo uses `Seed` to derive all generated identifiers:
  `PathNotFoundTemplData`, `PathNotFoundError`, `ErrPathNotFound`,
  `NewPathNotFoundError`.

---

## 4. Fields Rules

- `Fields` must be **empty** for all static `TypeName` values.
- `Fields` must be **non-empty** for all dynamic `TypeName` values.
- Every `Note` in `Fields` must have a matching `{{.Note}}` token in `Other`
  and vice versa - lingo validates this and refuses to generate on mismatch.
- For wrapper types, exactly one `Fields` entry must have `GoType: "error"`
  and `Note: "Wrapped"`. No other entry may use `GoType: "error"`.
- `GoType` must be a valid native Go type: `"string"`, `"int"`, `"uint"`,
  `"error"`.
- Always populate `Tale` - if empty lingo emits a 🔥 TODO in the generated
  doc comment instead.

---

## 5. The `File` Field

Leave `File` empty in almost all cases. lingo then writes to the default
output file for the message kind:
- Cobra kinds -> `messages-cobra-auto.go`
- General kinds -> `messages-general-auto.go`
- Error kinds -> `messages-errors-auto.go`

Set `File` only when output must go to a custom file, e.g.
`File: "system-automation"` produces `system-automation-errors-auto.go`.

---

## 6. Placeholder Pattern for Ungenerated Messages

Whenever a coding agent needs to display user-facing text or return a
translatable error, it must **never** use `fmt.Printf`, `fmt.Println`,
`errors.New`, or `fmt.Errorf` for that content. Instead it must emit a
placeholder call site with a structured comment that tells the developer
exactly what lingo entry to create.

This applies to:
- Any text shown to the user (output, progress, status, warnings)
- Any error that could be triggered by end-user behaviour

### Placeholder format - error

```go
// i18n-TODO: add Underliers entry for this error
// Seed: "<PascalCaseName>"
// TypeName: enums.<UnderlyingTypeXxx>
// Other: "<the error message, with {{.Token}} for any variables>"
// Fields: <list field names and Go types, or "none" if static>
return locale.NewPlaceholderError() // replace once lingo-generated
```

### Placeholder format - user-facing text

```go
// i18n-TODO: add Underliers entry for this message
// Seed: "<PascalCaseName>"
// TypeName: enums.<UnderlyingTypeXxx>
// Other: "<the message text, with {{.Token}} for any variables>"
// Fields: <list field names and Go types, or "none" if static>
fmt.Println("<placeholder: replace with li18ngo.Text(locale.XxxTemplData{...})>")
```

### Rules

- The `i18n-TODO:` prefix is mandatory - it makes placeholders greppable
  across the codebase: `grep -r "i18n-TODO" ./src`
- Always specify `TypeName` using the full enum name from section 2 so the
  developer knows exactly which scenario applies.
- Always write out the intended `Other` string, including any `{{.Token}}`
  placeholders, so the developer can copy it directly into the `Underliers`
  entry.
- Always list `Fields` even when there are none - write `"none"` explicitly
  so the developer does not have to infer it.
- Never leave a bare `fmt.Println` or `fmt.Errorf` for user-facing content
  without an `i18n-TODO:` comment above it.

### Example - dynamic error placeholder

```go
// i18n-TODO: add Underliers entry for this error
// Seed: "UnknownDisplayMode"
// TypeName: enums.UnderlyingTypeDynamicError
// Other: "unknown display mode {{.Mode}} (valid modes: {{.ValidModes}})"
// Fields: Mode string, ValidModes string (pre-formatted CSV via strings.Join)
return nil, locale.NewPlaceholderError() // replace once lingo-generated
```

### Example - static user-facing text placeholder

```go
// i18n-TODO: add Underliers entry for this message
// Seed: "TraversalComplete"
// TypeName: enums.UnderlyingTypeStaticGeneral
// Other: "Traversal complete"
// Fields: none
fmt.Println("<placeholder: replace with li18ngo.Text(locale.TraversalCompleteTemplData{})>")
```

---

## 7. Programming Errors - i18n Exemption

Errors that can only be caused by a programmer mistake during development -
and can never be triggered by end-user behaviour at runtime - are exempt from
lingo. These may use `errors.New` (static) or `fmt.Errorf` (dynamic) directly.

Canonical examples:
- Duplicate registration of a named mode or handler
- Registration of a nil factory or constructor
- Incorrect wiring in the composition root detected at startup

Always document the exemption with a `programmingError:` comment:

```go
// programmingError: duplicate mode registration can only be caused by
// incorrect wiring in the composition root, never by end-user input.
// exempt from i18n - English only.
return fmt.Errorf("ui: display mode %q is already registered", name)
```

The test: could an end user trigger this error through normal CLI use?
- Yes -> must go through lingo; use the placeholder pattern from section 6.
- No  -> may use `errors.New` / `fmt.Errorf` directly with a `programmingError:` comment.

---

## 8. Lists of Valid Values in Error Messages

When an error message must include a list of valid values (e.g. valid mode
names), go-i18n represents this as a single `string` field. The caller must
pre-format the list before passing it to the generated constructor:

```go
// ValidModes must be pre-formatted before passing to the constructor.
strings.Join(registeredModes(), ", ")
```

The `Fields` entry for such a field uses `GoType: "string"` and the `Tale`
should document that the caller is responsible for formatting.

Note: the lingo generator does not yet handle `[]string` fields natively.
When this limitation is addressed, the generator will produce the
`strings.Join` internally and callers will pass a slice directly.

---

## 9. Worked Examples

### Static non-error message

```go
"localisation.test": {
    MessageID:   "localisation.test",
    Seed:        "Localisation",
    TypeName:    enums.UnderlyingTypeStaticGeneral,
    Description: "Localisation",
    Story:       "A test message for localisation.",
    Other:       "localisation",
},
```

### Dynamic non-error message

```go
"using-config-file": {
    MessageID:   "using-config-file",
    Seed:        "UsingConfigFile",
    TypeName:    enums.UnderlyingTypeDynamicGeneral,
    Description: "Message to indicate which config is being used",
    Story: "UsingConfigFile is printed on startup to indicate" +
        " which configuration file has been loaded.",
    Other: "Using config file: '{{.ConfigFileName}}'",
    Fields: []UnderlyingField{
        {
            Note:   "ConfigFileName",
            GoType: "string",
            Tale:   "is the name of the config file being used",
        },
    },
},
```

### Dynamic error, no wrapping

```go
"path-not-found.dynamic-error": {
    MessageID:   "path-not-found.dynamic-error",
    Seed:        "PathNotFound",
    TypeName:    enums.UnderlyingTypeDynamicError,
    Description: "Directory or file path does not exist",
    Story:       "PathNotFoundError is used when a directory or file path does not exist.",
    Other:       "{{.Name}} path not found ({{.Path}})",
    Fields: []UnderlyingField{
        {
            Note:   "Name",
            GoType: "string",
            Tale:   "is the name of the path that was not found (e.g. 'Config')",
        },
        {
            Note:   "Path",
            GoType: "string",
            Tale:   "is the actual path that was not found (e.g. '/etc/config.yaml')",
        },
    },
},
```

### Static error wrapping another error (message includes wrapped text)

```go
"third-party.error-wrapper-msg": {
    MessageID:   "third-party.error-wrapper-msg",
    Seed:        "ThirdPartyWrapper",
    TypeName:    enums.UnderlyingTypeStaticErrorWrapperMsg,
    Description: "Wrapper for third-party errors",
    Story:       "ThirdPartyErrorWrapper is used to wrap errors from third-party libraries.",
    Other:       "Third party error occurred: '{{.Wrapped}}'",
    Fields: []UnderlyingField{
        {
            Note:   "Wrapped",
            GoType: "error",
            Tale:   "is the original error from the third-party library that is being wrapped",
        },
    },
},
```

---

## 10. Common Mistakes to Avoid

| Mistake | Correct behaviour |
|---|---|
| Hand-writing generated structs or error types | Only edit `underlying-templ-data.go`; run lingo to generate |
| `Fields` non-empty on a static `TypeName` | Static types must have empty `Fields` |
| `Fields` empty on a dynamic `TypeName` | Dynamic types must have non-empty `Fields` |
| `Note` in `Fields` not matching a `{{.Token}}` in `Other` | Every `Note` must have a matching token and vice versa |
| More than one `Fields` entry with `GoType: "error"` | Exactly one error field permitted, named `"Wrapped"` |
| `GoType: "error"` field named anything other than `"Wrapped"` | The wrapped error field must always be named `"Wrapped"` |
| Omitting `Tale` | Always provide `Tale`; empty emits a 🔥 TODO in generated code |
| Duplicate `MessageID` across the map | lingo refuses to generate on duplicate IDs |
| Using `fmt.Println` for user-facing text without an `i18n-TODO:` comment | All user-facing text must go through lingo; use the placeholder pattern |
| Using `fmt.Errorf` / `errors.New` for a user-facing error | User-facing errors must go through lingo; use the placeholder pattern |
| Using lingo for a programming error | Programming errors are exempt; use `fmt.Errorf` / `errors.New` with a `programmingError:` comment |
| Omitting the `programmingError:` comment on an exempt error | Always document why the error is exempt |
| Omitting the `TypeName` from an `i18n-TODO:` comment | Always specify the full enum value so the developer knows which scenario applies |
| Passing a `[]string` directly to a generated constructor | Pre-format with `strings.Join` before passing; lingo generates a `string` field |
