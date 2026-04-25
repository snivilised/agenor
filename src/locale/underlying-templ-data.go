package locale

import (
	lingo "github.com/snivilised/li18ngo/locale"
	"github.com/snivilised/li18ngo/locale/enums"
)

//go:generate lingo

// underliers is the single source of truth for all i18n messages in this
// package. Edit this map and run go generate to regenerate the auto files.
var _ = lingo.Underliers{

	// -------------------------------------------------------------------------
	// Cobra messages
	// -------------------------------------------------------------------------

	"root-command-short-description": {
		MessageID:   "root-command-short-description",
		Seed:        "RootCmdShortDesc",
		TypeName:    enums.UnderlyingTypeStaticCobra,
		Description: "short description for the root command",
		Story: "RootCmdShortDesc is the short description shown in" +
			" cobra help output for the root command.",
		Other: "A brief description of your application",
	},

	"root-command-long-description": {
		MessageID:   "root-command-long-description",
		Seed:        "RootCmdLongDesc",
		TypeName:    enums.UnderlyingTypeStaticCobra,
		Description: "long description for the root command",
		Story: "RootCmdLongDesc is the long description shown in" +
			" cobra help output for the root command.",
		Other: `A longer description that spans multiple lines and likely contains
		examples and usage of using your application. For example:

		Cobra is a CLI library for Go that empowers applications.
		This application is a tool to generate the needed files
		to quickly create a Cobra application.`,
	},

	"root-command-config-file-usage": {
		MessageID:   "root-command-config-file-usage",
		Seed:        "RootCmdConfigFileUsage",
		TypeName:    enums.UnderlyingTypeDynamicCobra,
		Description: "root command config flag usage",
		Story: "RootCmdConfigFileUsage is the usage string for the" +
			" config file flag on the root command.",
		Other: "config file (default is $HOME/{{.ConfigFileName}}.yml)",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "ConfigFileName",
				GoType: "string",
				Tale:   "is the base name of the config file without extension",
			},
		},
	},

	"root-command-language-usage": {
		MessageID:   "root-command-language-usage",
		Seed:        "RootCmdLangUsage",
		TypeName:    enums.UnderlyingTypeStaticCobra,
		Description: "root command lang usage",
		Story: "RootCmdLangUsage is the usage string for the" +
			" language flag on the root command.",
		Other: "'lang' defines the language according to IETF BCP 47",
	},

	// -------------------------------------------------------------------------
	// General messages
	// -------------------------------------------------------------------------

	"using-config-file": {
		MessageID:   "using-config-file",
		Seed:        "UsingConfigFile",
		TypeName:    enums.UnderlyingTypeDynamicGeneral,
		Description: "Message to indicate which config is being used",
		Story: "UsingConfigFile is printed on startup to indicate" +
			" which configuration file has been loaded.",
		Other: "Using config file: '{{.ConfigFileName}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "ConfigFileName",
				GoType: "string",
				Tale:   "is the name of the config file being used",
			},
		},
	},

	"node-visited": {
		MessageID:   "node-visited",
		Seed:        "NodeVisited",
		TypeName:    enums.UnderlyingTypeDynamicGeneral,
		Description: "Printed for each node visited during traversal",
		Story: "NodeVisited is printed for each filesystem node" +
			" encountered during traversal.",
		Other: "Node Visited -> '{{.Path}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Path",
				GoType: "string",
				Tale:   "is the path of the node being visited",
			},
		},
	},

	"action-failed": {
		MessageID: "action-failed",
		Seed:      "ActionFailed",
		TypeName:  enums.UnderlyingTypeDynamicGeneral,
		Description: "Printed when an action fails on a node" +
			" during traversal",
		Story: "ActionFailed is printed when a named action returns" +
			" an error while processing a node.",
		Other: "[!] action '{{.Name}}' failed on '{{.Path}}': '{{.Err}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Name",
				GoType: "string",
				Tale:   "is the name of the action that failed",
			},
			{
				Note:   "Path",
				GoType: "string",
				Tale:   "is the path of the node on which the action failed",
			},
			{
				Note:   "Err",
				GoType: "string",
				Tale:   "is the error message",
			},
		},
	},

	"action-visited": {
		MessageID: "action-visited",
		Seed:      "ActionVisited",
		TypeName:  enums.UnderlyingTypeDynamicGeneral,
		Description: "Printed for each node successfully processed" +
			" by an action",
		Story: "ActionVisited is printed for each node successfully" +
			" processed by a named action.",
		Other: "[+] actioned for '{{.Name}}' at '{{.Path}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Name",
				GoType: "string",
				Tale:   "is the name of the action that succeeded",
			},
			{
				Note:   "Path",
				GoType: "string",
				Tale:   "is the path of the node successfully processed",
			},
		},
	},

	"pipeline-failed": {
		MessageID: "pipeline-failed",
		Seed:      "PipelineFailed",
		TypeName:  enums.UnderlyingTypeDynamicGeneral,
		Description: "Printed when a pipeline fails on a node" +
			" during traversal",
		Story: "PipelineFailed is printed when a named pipeline" +
			" returns an error while processing a node.",
		Other: "[!] pipeline '{{.Name}}' failed on '{{.Path}}': '{{.Err}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Name",
				GoType: "string",
				Tale:   "is the name of the pipeline that failed",
			},
			{
				Note:   "Path",
				GoType: "string",
				Tale:   "is the path of the node on which the pipeline failed",
			},
			{
				Note:   "Err",
				GoType: "string",
				Tale:   "is the error message",
			},
		},
	},

	"pipeline-visited": {
		MessageID: "pipeline-visited",
		Seed:      "PipelineVisited",
		TypeName:  enums.UnderlyingTypeDynamicGeneral,
		Description: "Printed for each node successfully processed" +
			" by a pipeline",
		Story: "PipelineVisited is printed for each node successfully" +
			" processed by a named pipeline.",
		Other: "[+] pipeline succeeded for '{{.Name}}' at '{{.Path}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Name",
				GoType: "string",
				Tale:   "is the name of the pipeline that succeeded",
			},
			{
				Note:   "Path",
				GoType: "string",
				Tale:   "is the path of the node successfully processed",
			},
		},
	},

	"traversal-failed": {
		MessageID:   "traversal-failed",
		Seed:        "TraversalFailed",
		TypeName:    enums.UnderlyingTypeDynamicGeneral,
		Description: "Printed when the traversal itself fails",
		Story: "TraversalFailed is printed when the traversal" +
			" operation itself encounters a fatal error.",
		Other: "[!] traversal failed: '{{.Err}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Err",
				GoType: "string",
				Tale:   "is the error message",
			},
		},
	},

	"action-skipped": {
		MessageID:   "action-skipped",
		Seed:        "ActionSkipped",
		TypeName:    enums.UnderlyingTypeDynamicGeneral,
		Description: "Printed when an action is skipped",
		Story:       "ActionSkipped is printed when an action is skipped",
		Other:       "[!] action '{{.Name}}' skipped at '{{.Path}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Name",
				GoType: "string",
				Tale:   "Name of the action that was skipped",
			},
			{
				Note:   "Path",
				GoType: "string",
				Tale:   "Full path of the node that action that was skipped for",
			},
		},
	},

	"placeholder-breach": {
		MessageID:   "placeholder-breach",
		Seed:        "PlaceholderBreach",
		TypeName:    enums.UnderlyingTypeDynamicGeneral,
		Description: "Printed when placeholder breach occurs during expansion",
		Story:       "Printed when a placeholder breach occurs during expansion.",
		Other:       "Placeholder: '{{.Placeholder}}' resolved to '{{.ResolvedPath}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Placeholder",
				GoType: "string",
				Tale:   "The placeholder that breached above the traversal root",
			},
			{
				Note:   "ResolvedPath",
				GoType: "string",
				Tale:   "The resolved path of the placeholder that breached above the traversal root",
			},
		},
	},

	"traversal-complete": {
		MessageID:   "traversal-complete",
		Seed:        "TraversalComplete",
		TypeName:    enums.UnderlyingTypeDynamicGeneral,
		Description: "Printed on successful completion of a traversal",
		Story: "TraversalComplete is printed when a traversal finishes" +
			" successfully, summarising the nodes visited and time elapsed.",
		Other: "[+] traversal complete successfully: {{.Files}} files, {{.Dirs}} dirs" +
			" visited in {{.Elapsed}}",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Files",
				GoType: "uint",
				Tale:   "is the total number of files visited",
			},
			{
				Note:   "Dirs",
				GoType: "uint",
				Tale:   "is the total number of directories visited",
			},
			{
				Note:   "Elapsed",
				GoType: "string",
				Tale:   "is the total time elapsed during the traversal",
			},
		},
	},

	// -------------------------------------------------------------------------
	// Error messages
	// -------------------------------------------------------------------------

	"filter-is-nil.static-error": {
		MessageID:   "filter-is-nil.static-error",
		Seed:        "FilterIsNil",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "filter is nil error",
		Story: "FilterIsNil indicates that the caller passed a nil" +
			" filter reference where a concrete filter implementation" +
			" was required.",
		Other: "filter is nil",
	},

	"filter-missing-type.static-error": {
		MessageID:   "filter-missing-type.static-error",
		Seed:        "FilterMissingType",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "filter missing type",
		Story: "FilterMissingType indicates that the filter definition" +
			" is missing a required type field.",
		Other: "filter missing type",
	},

	"custom-filter-not-supported-for-children.static-error": {
		MessageID: "custom-filter-not-supported-for-children.static-error" +
			".static-error",
		Seed:        "FilterCustomNotSupported",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "custom filter not supported for children",
		Story: "FilterCustomNotSupported indicates that custom filters" +
			" cannot be applied to child nodes in this context.",
		Other: "custom filter not supported for children",
	},

	"glob-ex-filter-not-supported-for-children.static-error": {
		MessageID: "glob-ex-filter-not-supported-for-children.static-error" +
			".static-error",
		Seed:        "FilterChildGlobExNotSupported",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "glob-ex filter not supported for children",
		Story: "FilterChildGlobExNotSupported indicates that glob-ex" +
			" filters cannot be applied to child nodes in this context.",
		Other: "glob-ex filter not supported for children",
	},

	"filter-is-undefined.static-error": {
		MessageID:   "filter-is-undefined.static-error",
		Seed:        "FilterUndefined",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "filter is undefined error",
		Story: "FilterUndefined indicates that the filter referenced" +
			" in the traversal options has not been defined.",
		Other: "filter is undefined",
	},

	"failed-to-get-navigator-driver.static-error": {
		MessageID:   "failed-to-get-navigator-driver.static-error",
		Seed:        "InternalFailedToGetNavigatorDriver",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "failed to get navigator driver",
		Story: "InternalFailedToGetNavigatorDriver indicates an" +
			" internal failure when resolving the navigator driver." +
			" This is not expected during normal operation.",
		Other: "failed to get navigator driver",
	},

	// TODO: Need to to check where this is invoked from, may not be required.
	"invalid-incase-filter-definition.dynamic-error": {
		MessageID: "invalid-incase-filter-definition.dynamic-error",
		Seed:      "InvalidInCaseFilterDef",
		TypeName:  enums.UnderlyingTypeDynamicErrorWrapper,
		Description: "invalid incase filter definition; pattern is" +
			" missing separator wrapper error",
		Story: "InvalidInCaseFilterDef indicates that a case-insensitive" +
			" filter definition is invalid because the pattern is missing" +
			" the required separator.",
		Other: "'{{.Wrapped}}', pattern: '{{.Pattern}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Wrapped",
				GoType: "error",
				Tale:   "is the underlying error",
			},
			{
				Note:   "Pattern",
				GoType: "string",
				Tale:   "is the invalid filter pattern",
			},
		},
	},

	// TODO: Need to to check where this is invoked from, may not be required.
	"invalid-incase-filter-definition.sentinel-error": {
		MessageID: "invalid-incase-filter-definition.sentinel-error",
		Seed:      "CoreInvalidInCaseFilterDef",
		TypeName:  enums.UnderlyingTypeSentinelError,
		Description: "invalid incase filter definition; pattern is" +
			" missing separator core error",
		Story: "CoreInvalidInCaseFilterDef is the sentinel core error" +
			" for an invalid case-insensitive filter definition. Wrap" +
			" this error using NewInvalidInCaseFilterDefError.",
		Other: "invalid incase filter definition;" +
			" pattern is missing separator",
	},

	"failed-to-create-worker-pool.static-error": {
		MessageID:   "failed-to-create-worker-pool.static-error",
		Seed:        "WorkerPoolCreationFailed",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "failed to create worker pool",
		Story: "WorkerPoolCreationFailed indicates that the worker" +
			" pool could not be initialised.",
		Other: "failed to create worker pool",
	},

	"invalid-file-sampling-spec-missing-files.static-error": {
		MessageID: "invalid-file-sampling-spec-missing-files" +
			".static-error",
		Seed:     "InvalidFileSamplingSpecMissingFiles",
		TypeName: enums.UnderlyingTypeStaticError,
		Description: "invalid file sampling specification," +
			" missing no of files",
		Story: "InvalidFileSamplingSpecMissingFiles indicates that" +
			" the file sampling specification is invalid because the" +
			" required number-of-files field is absent.",
		Other: "invalid file sampling specification," +
			" missing no of files",
	},

	"invalid-file-sampling-spec-missing-directories.static-error": {
		MessageID: "invalid-file-sampling-spec-missing-directories" +
			".static-error",
		Seed:     "InvalidSamplingSpecMissingDirectories",
		TypeName: enums.UnderlyingTypeStaticError,
		Description: "invalid file sampling specification," +
			" missing no of directories",
		Story: "InvalidSamplingSpecMissingDirectories indicates that" +
			" the file sampling specification is invalid because the" +
			" required number-of-directories field is absent.",
		Other: "invalid file sampling specification," +
			" missing no of directories",
	},

	"missing-custom-filter-definition.static-error": {
		MessageID:   "missing-custom-filter-definition.static-error",
		Seed:        "MissingCustomFilterDefinition",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "config error missing-custom-filter-definition",
		Story: "MissingCustomFilterDefinition indicates that the" +
			" traversal configuration references a custom filter but" +
			" no definition for it was found.",
		Other: "missing custom filter definition (config error)",
	},

	"invalid-glob-ex-filter-missing-separator.dynamic-error": {
		MessageID: "invalid-glob-ex-filter-missing-separator" +
			".dynamic-error",
		Seed:     "InvalidExtGlobFilterMissingSeparator",
		TypeName: enums.UnderlyingTypeDynamicError,
		Description: "invalid glob ex filter definition;" +
			" pattern is missing separator",
		Story: "InvalidExtGlobFilterMissingSeparator indicates that" +
			" an extended glob filter definition is invalid because" +
			" the pattern is missing the required separator character.",
		Other: "extended glob pattern missing separator, pattern: '{{.Pattern}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Pattern",
				GoType: "string",
				Tale:   "is the invalid filter pattern",
			},
		},
	},

	"invalid-extended-glob-filter-missing-separator.sentinel-error": {
		MessageID: "invalid-extended-glob-filter-missing-separator" +
			".sentinel-error",
		Seed:     "CoreInvalidExtGlobFilterMissingSeparator",
		TypeName: enums.UnderlyingTypeSentinelError,
		Description: "invalid glob ex filter definition;" +
			" pattern is missing separator",
		Story: "CoreInvalidExtGlobFilterMissingSeparator is the" +
			" sentinel core error for an invalid extended glob filter" +
			" definition. Wrap this error using" +
			" NewInvalidExtGlobFilterMissingSeparatorError.",
		Other: "invalid glob ex filter definition;" +
			" pattern is missing separator",
	},

	"poly-filter-is-invalid.static-error": {
		MessageID:   "poly-filter-is-invalid.static-error",
		Seed:        "PolyFilterIsInvalid",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "poly filter definition is invalid error",
		Story: "PolyFilterIsInvalid indicates that a poly filter" +
			" definition fails validation.",
		Other: "poly filter definition is invalid",
	},

	"usage-missing-tree-path.static-error": {
		MessageID:   "usage-missing-tree-path.static-error",
		Seed:        "UsageMissingTreePath",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "usage missing tree path",
		Story: "UsageMissingTreePath indicates that the command was" +
			" invoked without the required tree path argument.",
		Other: "usage missing tree path",
	},

	"usage-missing-restore-path.static-error": {
		MessageID:   "usage-missing-restore-path.static-error",
		Seed:        "UsageMissingRestorePath",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "usage missing restore path",
		Story: "UsageMissingRestorePath indicates that the command was" +
			" invoked without the required restore path argument.",
		Other: "usage missing restore path",
	},

	"usage-missing-subscription.static-error": {
		MessageID:   "usage-missing-subscription.static-error",
		Seed:        "UsageMissingSubscription",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "usage missing subscription",
		Story: "UsageMissingSubscription indicates that the command" +
			" was invoked without specifying a subscription type.",
		Other: "usage missing subscription",
	},

	"usage-missing-handler.static-error": {
		MessageID:   "usage-missing-handler.static-error",
		Seed:        "UsageMissingHandler",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "usage missing handler",
		Story: "UsageMissingHandler indicates that the command was" +
			" invoked without registering a required handler.",
		Other: "usage missing handler",
	},

	"id-generator-func-cant-be-nil.static-error": {
		MessageID:   "id-generator-func-cant-be-nil.static-error",
		Seed:        "IDGeneratorFuncCantBeNil",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "id generator func is nil, should be defined",
		Story: "IDGeneratorFuncCantBeNil indicates that a nil function" +
			" was supplied where an ID generator func is required.",
		Other: "id generator func can't be nil",
	},

	"un-equal-conversion.sentinel-error": {
		MessageID:   "un-equal-conversion.sentinel-error",
		Seed:        "UnEqualJSONConversion",
		TypeName:    enums.UnderlyingTypeSentinelError,
		Description: "JSON options conversion error",
		Story: "UnEqualJSONConversion indicates that a round-trip" +
			" JSON conversion produced a result that does not equal" +
			" the original value.",
		Other: "unequal JSON conversion",
	},

	"invalid-path.dynamic-error": {
		MessageID:   "invalid-path.dynamic-error",
		Seed:        "InvalidPath",
		TypeName:    enums.UnderlyingTypeDynamicError,
		Description: "invalid path (dynamic error)",
		Story: "InvalidPath indicates that a path supplied by the" +
			" caller fails validation.",
		Other: "path: '{{.Path}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Path",
				GoType: "string",
				Tale:   "is the invalid path",
			},
		},
	},

	"traverse-fs-mismatch.static-error": {
		MessageID:   "traverse-fs-mismatch.static-error",
		Seed:        "TraverseFsMismatch",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "traverse fs mismatch error",
		Story: "TraverseFsMismatch indicates that the filesystem" +
			" passed to the traversal does not match the filesystem" +
			" recorded at the point the session was saved.",
		Other: "traverse-fs mismatch",
	},

	"resume-fs-mismatch.static-error": {
		MessageID:   "resume-fs-mismatch.static-error",
		Seed:        "ResumeFsMismatch",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "resume fs mismatch error",
		Story: "ResumeFsMismatch indicates that the filesystem passed" +
			" to a resume operation does not match the filesystem" +
			" recorded at the point the session was saved.",
		Other: "resume-fs mismatch",
	},

	"resume-fs-mismatch.sentinel-error": {
		MessageID:   "resume-fs-mismatch.sentinel-error",
		Seed:        "CoreResumeFsMismatch",
		TypeName:    enums.UnderlyingTypeSentinelError,
		Description: "core resume file system mismatch error",
		Story: "CoreResumeFsMismatch is the sentinel core error for" +
			" a filesystem mismatch detected during traversal or resume." +
			" Wrap using NewTraverseFsMismatchError or" +
			" NewResumeFsMismatchError.",
		Other: "resume file system mismatch",
	},

	"panic-occurred.sentinel-error": {
		MessageID:   "panic-occurred.sentinel-error",
		Seed:        "CorePanicOccurred",
		TypeName:    enums.UnderlyingTypeSentinelError,
		Description: "core error",
		Story: "CorePanicOccurred is the sentinel core error indicating" +
			" that a panic was intercepted during traversal.",
		Other: "panic occurred",
	},

	"invalid-subscription.static-error": {
		MessageID:   "invalid-subscription.static-error",
		Seed:        "InvalidSubscription",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "invalid subscription type",
		Story: "InvalidSubscription indicates that the subscription" +
			" type supplied by the caller is not one of the accepted values.",
		Other: "invalid subscription type," +
			" must be one of: [files, dirs, all]",
	},

	"invalid-resume-strategy.static-error": {
		MessageID:   "invalid-resume-strategy.static-error",
		Seed:        "InvalidResumeStrategy",
		TypeName:    enums.UnderlyingTypeStaticError,
		Description: "invalid resume strategy type",
		Story: "InvalidResumeStrategy indicates that the resume strategy" +
			" supplied by the caller is not one of the accepted values.",
		Other: "invalid resume strategy, must be one of: [spawn, fast]",
	},

	// TODO: rename bedrock errors; shouldn't use bedrock in the name
	"bedrock-load-viper-setup.jaywalk.static-error": {
		MessageID: "bedrock-load-viper-setup.jaywalk.static-error",
		Seed:      "BedrockLoadViperSetup",
		TypeName:  enums.UnderlyingTypeStaticErrorWrapper,
		Description: "Error returned when viper setup fails" +
			" during bedrock.Load",
		Story: "BedrockLoadViperSetup indicates that viper could not" +
			" be configured during the bedrock.Load call.",
		Other: "bedrock.Load: viper setup",
	},

	"bedrock-load-reading-config.jaywalk.static-error": {
		MessageID: "bedrock-load-reading-config.jaywalk.static-error",
		Seed:      "BedrockLoadReadingConfig",
		TypeName:  enums.UnderlyingTypeStaticErrorWrapper,
		Description: "Error returned when reading the config file" +
			" fails during bedrock.Load",
		Story: "BedrockLoadReadingConfig indicates that the" +
			" configuration file could not be read during the" +
			" bedrock.Load call.",
		Other: "bedrock.Load: reading config",
	},

	"bedrock-load-decoding.jaywalk.static-error": {
		MessageID: "bedrock-load-decoding.jaywalk.static-error",
		Seed:      "BedrockLoadDecoding",
		TypeName:  enums.UnderlyingTypeStaticErrorWrapper,
		Description: "Error returned when config decoding fails" +
			" during bedrock.Load",
		Story: "BedrockLoadDecoding indicates that the configuration" +
			" could not be decoded into the target struct during the" +
			" bedrock.Load call.",
		Other: "bedrock.Load: decoding",
	},

	"bedrock-load-validation.jaywalk.static-error": {
		MessageID: "bedrock-load-validation.jaywalk.static-error",
		Seed:      "BedrockLoadValidation",
		TypeName:  enums.UnderlyingTypeStaticErrorWrapper,
		Description: "Error returned when config validation fails" +
			" during bedrock.Load",
		Story: "BedrockLoadValidation indicates that the decoded" +
			" configuration failed validation during the bedrock.Load call.",
		Other: "bedrock.Load: validation",
	},

	"unsupported-format.jaywalk.dynamic-error": {
		MessageID: "unsupported-format.jaywalk.dynamic-error",
		Seed:      "UnsupportedFormat",
		TypeName:  enums.UnderlyingTypeDynamicError,
		Description: "Error returned when an unregistered config" +
			" format is requested",
		Story: "UnsupportedFormat indicates that the configuration" +
			" format requested by the caller has not been registered" +
			" with the format registry.",
		Other: "unsupported format '{{.Format}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Format",
				GoType: "string",
				Tale:   "is the unsupported format",
			},
		},
	},

	"creating-decoder-for.jaywalk.dynamic-error": {
		MessageID: "creating-decoder-for.jaywalk.dynamic-error",
		Seed:      "CreatingDecoderFor",
		TypeName:  enums.UnderlyingTypeDynamicErrorWrapper,
		Description: "Error returned when a mapstructure decoder" +
			" cannot be created for a config section",
		Story: "CreatingDecoderFor indicates that a mapstructure" +
			" decoder could not be constructed for the named" +
			" configuration section.",
		Other: "creating decoder for '{{.Key}}': '{{.Wrapped}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Wrapped",
				GoType: "error",
				Tale:   "is the underlying error",
			},
			{
				Note:   "Key",
				GoType: "string",
				Tale:   "is the configuration section name",
			},
		},
	},

	"decoding-section.jaywalk.dynamic-error": {
		MessageID: "decoding-section.jaywalk.dynamic-error",
		Seed:      "DecodingSection",
		TypeName:  enums.UnderlyingTypeDynamicErrorWrapper,
		Description: "Error returned when mapstructure decoding of" +
			" a config section fails",
		Story: "DecodingSection indicates that mapstructure failed" +
			" to decode the named configuration section into its" +
			" target struct.",
		Other: "decoding section '{{.Key}}': '{{.Wrapped}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Wrapped",
				GoType: "error",
				Tale:   "is the underlying error",
			},
			{
				Note:   "Key",
				GoType: "string",
				Tale:   "is the configuration section name",
			},
		},
	},

	"flags-section-unexpected-type.jaywalk.dynamic-error": {
		MessageID: "flags-section-unexpected-type.jaywalk.dynamic-error",
		Seed:      "FlagsSectionUnexpectedType",
		TypeName:  enums.UnderlyingTypeDynamicError,
		Description: "Error returned when the flags config section" +
			" has an unexpected type",
		Story: "FlagsSectionUnexpectedType indicates that the flags" +
			" section of the configuration file was decoded into an" +
			" unexpected Go type.",
		Other: "flags section has unexpected type '{{.TypeName}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "TypeName",
				GoType: "string",
				Tale:   "is the unexpected type",
			},
		},
	},

	"traversal-saved.dynamic-error": {
		MessageID: "traversal-saved.dynamic-error",
		Seed:      "TraversalSaved",
		TypeName:  enums.UnderlyingTypeDynamicErrorWrapper,
		Description: "Error returned when a panic occurs during traversal" +
			" and save of the traversal state succeeds",
		Story: "Traversal state saved as a result of a panic occurring during traversal",
		Other: "'{{.Wrapped}}' Traversal state saved successfully to: '{{.SavedTo}}' ",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "SavedTo",
				GoType: "string",
				Tale:   "The path the traversal state was saved to",
			},
			{
				Note:   "Wrapped",
				GoType: "error",
				Tale:   "The underlying error representing the panic that occurred",
			},
		},
	},

	"traversal-not-saved.dynamic-error": {
		MessageID: "traversal-not-saved.dynamic-error",
		Seed:      "TraversalNotSaved",
		TypeName:  enums.UnderlyingTypeDynamicErrorWrapper,
		Description: "Error returned when a panic occurs during traversal" +
			" and save of the traversal state fails",
		Story: "Traversal state not saved as a result of a panic occurring during traversal",
		Other: "'{{.Wrapped}}' Failed to save traversal state to: '{{.SavedTo}}' ",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "SavedTo",
				GoType: "string",
				Tale:   "The path the traversal state was attempted to be saved to",
			},
			{
				Note:   "Wrapped",
				GoType: "error",
				Tale: "The underlying error representing the panic that occurred and the" +
					" failure to save",
			},
		},
	},

	"action-not-found.dynamic-error": {
		MessageID:   "action-not-found.dynamic-error",
		Seed:        "ActionNotFound",
		TypeName:    enums.UnderlyingTypeDynamicError,
		Description: "Action not found in config error",
		Story: "ActionNotFound indicates that an action with the specified" +
			" name was not found in the traversal configuration.",
		Other: "Action not found: '{{.Action}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Action",
				GoType: "string",
				Tale:   "is the action name that was not found",
			},
		},
	},

	"pipeline-not-found.dynamic-error": {
		MessageID:   "pipeline-not-found.dynamic-error",
		Seed:        "PipelineNotFound",
		TypeName:    enums.UnderlyingTypeDynamicError,
		Description: "Pipeline not found in config error",
		Story: "PipelineNotFound indicates that a pipeline with the specified" +
			" name was not found in the traversal configuration.",
		Other: "Pipeline not found: '{{.Pipeline}}'",
		Fields: []lingo.UnderlyingField{
			{
				Note:   "Pipeline",
				GoType: "string",
				Tale:   "is the pipeline name that was not found",
			},
		},
	},
}
