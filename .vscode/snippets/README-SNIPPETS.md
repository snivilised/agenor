# 🐷 Snippets

For further help with snippets see:

- [How To Create A VSCode Snippet (webdevsimplified)](https://blog.webdevsimplified.com/2022-03/vscode-snippet/)

## i18n Field

Creates a simple text field. The field is a key value pair, where the key is static text subject to translation.

- prefix: f18
- $1: the name of the field
- $2: the name of the source package from which to import the field
- example: locale.PatternFieldTemplData

## i18n Field and Variable i18n Error

Creates a variable i18n error (an error with state). This snippet uses a copy of `Simple i18n Field` to implement the variable part; ie the native variable error contains within it a field that implements the variable.

- prefix: v18e
- $1: the core name of the error (do not include Err or Error in name)
...

## Ginko Variable Native Error

Creates Ginkgo unit tests corresponding to `Variable Native Error`.

- prefix: gverrt
- $1: the core name of the error (do not include Err or Error in name)
- $2: the name of the source package from which to import the error
- example: core.errors_test

Creates the following test case:

- check content of variable error
- affirmative case; check that an error matches the target error being defined
- negative case; check that a different error does not match the target error being defined

## Simple i18n Error

Creates a simple i18n error message (ie a message without any variables)

- prefix: s18e
- $1: the core name of the error (do not include Err or Error in name)
- $2: the `TemplData` prefix, eg ___traverse___ to create `traverseTemplData`
- example: locale.FilterIsNilErrorTemplData

## Variable Native Error

Creates a native error (un-translated error, not meant for the end user). Note that by default it creates a single variable `value` of type string in the New function. This needs to be changed to be the state required for this message, or indeed, further state can be defined.

- prefix: verr
- $1: the core name of the error (do not include Err or Error in name)
- example: core.NewInvalidNotificationMuteRequestedError
