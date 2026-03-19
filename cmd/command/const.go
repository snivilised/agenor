package command

// ---------------------------------------------------------------------------
// Application identity
// ---------------------------------------------------------------------------

const (
	AppEmoji        = "🍒"
	ApplicationName = "jay"
	SourceID        = "github.com/snivilised/agenor"
)

// ---------------------------------------------------------------------------
// Param-set registration names
// ---------------------------------------------------------------------------

const (
	// root
	RootPsName = "root-ps"

	// shared families (registered on root, inherited by sub-commands)
	PreviewFamName     = "preview-fam"
	CascadeFamName     = "cascade-fam"
	InteractionFamName = "interaction-fam"
	SamplingFamName    = "sampling-fam"

	// run-only family
	WorkerPoolFamName = "worker-pool-fam"

	// filter family (registered per-command, not inherited)
	PolyFamName = "poly-fam"

	// jay-specific param sets
	WalkPsName = "walk-ps"
	RunPsName  = "run-ps"
)

// ---------------------------------------------------------------------------
// Resume strategy values
// ---------------------------------------------------------------------------

const (
	ResumeStrategySpawn    = "spawn"
	ResumeStrategyFastward = "fastward"
)
