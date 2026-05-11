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
- `Seed` must also be unique across the entire `Underliers` map. Because
  lingo derives every generated identifier from `Seed`, a duplicate value
  produces duplicate Go declarations and the build will fail. Verify
  uniqueness by scanning the existing map before inserting, exactly as you
  would for `MessageID`.
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

## 6. Adding New Messages and Errors - Agent Workflow

Whenever a coding agent needs to display user-facing text or return a
translatable error, it must **never** use `fmt.Printf`, `fmt.Println`,
`errors.New`, or `fmt.Errorf` for that content. Instead the agent must
complete the full lingo cycle itself before writing any call site.

This applies to:

- Any text shown to the user (output, progress, status, warnings)
- Any error that could be triggered by end-user behaviour

### Step 1 - Add entries to `Underliers`

Open `underlying-templ-data.go` and add one `UnderlyingTemplData` entry per
new message or error to the `Underliers` map. Use sections 1-5 of this skill
to determine the correct field values. Do not guess - verify both `MessageID`
and `Seed` uniqueness by scanning the existing map before inserting. Both
fields must be unique across the entire map; a duplicate in either will cause
lingo to refuse to generate or the build to fail.

Example entry (dynamic error):

```go
"path-not-found.dynamic-error": {
    MessageID:   "path-not-found.dynamic-error",
    Seed:        "PathNotFound",
    TypeName:    enums.UnderlyingTypeDynamicError,
    Description: "Directory or file path does not exist",
    Story:       "PathNotFoundError is used when a path does not exist.",
    Other:       "{{.Name}} path not found ({{.Path}})",
    Fields: []UnderlyingField{
        {
            Note:   "Name",
            GoType: "string",
            Tale:   "is the label for the path that was not found (e.g. 'Config')",
        },
        {
            Note:   "Path",
            GoType: "string",
            Tale:   "is the actual path value that was not found",
        },
    },
},
```

### Step 2 - Delete stale auto files

Before running lingo, delete every `*-auto.go` file that lingo owns. This
prevents stale generated identifiers from surviving into the new output and
avoids duplicate-declaration compile errors.

**Deletion scope - non-negotiable:** only files inside the `locale/`
directory are ever candidates for deletion. No file outside `locale/` may
be deleted under any circumstances, regardless of its name or extension.
This constraint is absolute and must never be relaxed.

The files lingo owns are determined by the `TypeName` values present in
`Underliers` (and the optional `File` field - see section 5). Always list
the directory first to confirm which auto files are actually present before
deleting anything:

```bash
ls locale/*-auto.go
```

Then delete only the files that appear in that listing:

```bash
# Default auto files
rm -f locale/messages-cobra-auto.go
rm -f locale/messages-general-auto.go
rm -f locale/messages-errors-auto.go

# Custom-file auto files - only if an Underliers entry uses File: "..."
# e.g. rm -f locale/system-automation-errors-auto.go
```

### Step 3 - Run lingo

```bash
go run ./cmd/lingo/...
```

Lingo regenerates all `*-auto.go` files from the current `Underliers` map.
If lingo exits with an error (duplicate `MessageID`, mismatched `{{.Token}}`
vs `Fields`, empty `Tale`, etc.), fix the entry in
`underlying-templ-data.go` and re-run. Do not proceed to step 4 until lingo
completes without error.

### Step 4 - Write call sites using generated types

Only after lingo has succeeded may call sites be written. The generated
identifiers follow directly from `Seed`:

| What you need | Generated identifier |
| --- | --- |
| Template data struct | `locale.<Seed>TemplData{}` |
| Error constructor | `locale.New<Seed>Error(...)` |
| Sentinel error value | `locale.Err<Seed>` |

**User-facing text** - call `li18ngo.Text`:

```go
li18ngo.Text(locale.TraversalCompleteTemplData{})
```

**Dynamic message** - populate fields on the struct:

```go
li18ngo.Text(locale.UsingConfigFileTemplData{
    ConfigFileName: cfgPath,
})
```

**Returning an error** - call the generated constructor:

```go
return locale.NewPathNotFoundError(err, "Config", path)
```

**Checking or wrapping a sentinel**:

```go
if errors.Is(err, locale.ErrSomeSentinel) { ... }
```

### Rules

- Complete all four steps before writing any call site. Never write a call
  site that references a generated type that does not yet exist on disk.
- Never leave a bare `fmt.Println` or `fmt.Errorf` for user-facing content.
- Never leave placeholder comments in source files - the agent performs the
  full lingo cycle, so there is nothing left for the developer to do.
- If lingo fails, fix `underlying-templ-data.go` and re-run; do not work
  around the failure by hand-editing an auto file.
- Always verify `MessageID` uniqueness before inserting - lingo will reject
  duplicates, but catching it before running saves a round-trip.

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
| --- | --- |
| Hand-writing generated structs or error types | Only edit `underlying-templ-data.go`; run lingo to generate |
| `Fields` non-empty on a static `TypeName` | Static types must have empty `Fields` |
| `Fields` empty on a dynamic `TypeName` | Dynamic types must have non-empty `Fields` |
| `Note` in `Fields` not matching a `{{.Token}}` in `Other` | Every `Note` must have a matching token and vice versa |
| More than one `Fields` entry with `GoType: "error"` | Exactly one error field permitted, named `"Wrapped"` |
| `GoType: "error"` field named anything other than `"Wrapped"` | The wrapped error field must always be named `"Wrapped"` |
| Omitting `Tale` | Always provide `Tale`; empty emits a 🔥 TODO in generated code |
| Duplicate `MessageID` across the map | lingo refuses to generate on duplicate IDs |
| Duplicate `Seed` across the map | lingo produces duplicate Go declarations; the build fails - scan the map before inserting |
| Using `fmt.Println` or `fmt.Errorf` for user-facing content | Complete the full lingo cycle (steps 1-4) instead |
| Writing a call site before lingo has run successfully | Always run lingo and confirm it exits cleanly before writing call sites |
| Skipping the auto-file deletion before re-running lingo | Always `rm -f locale/*-auto.go` before re-running to avoid stale declarations |
| Deleting any file outside the `locale/` directory | Deletion scope is `locale/` only - this constraint is absolute and non-negotiable |
| Hand-editing an auto file to work around a lingo error | Fix `underlying-templ-data.go` and re-run; never touch auto files directly |
| Using lingo for a programming error | Programming errors are exempt; use `fmt.Errorf` / `errors.New` with a `programmingError:` comment |
| Omitting the `programmingError:` comment on an exempt error | Always document why the error is exempt |
| Passing a `[]string` directly to a generated constructor | Pre-format with `strings.Join` before passing; lingo generates a `string` field |
