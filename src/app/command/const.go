package command

// ---------------------------------------------------------------------------
// Application identity
// ---------------------------------------------------------------------------

const (
	AppEmoji        = "🍒"
	ApplicationName = "jay"
	SourceID        = "github.com/snivilised/jaywalk/src/agenor"
)

// ---------------------------------------------------------------------------
// Param-set registration names
// ---------------------------------------------------------------------------

const (
	// root
	RootPsName = "root-ps"

	// nav ghost command param-set
	NavPsName = "nav-ps"

	// shared families (registered on nav ghost, inherited by walk/run/query)
	PreviewFamName     = "preview-fam"
	CascadeFamName     = "cascade-fam"
	InteractionFamName = "interaction-fam"
	SamplingFamName    = "sampling-fam"

	// poly-filter family (registered on nav ghost, inherited by walk/run/query)
	PolyFamName = "poly-fam"

	// run-only family
	WorkerPoolFamName = "worker-pool-fam"

	// jay-specific param sets
	WalkPsName  = "walk-ps"
	RunPsName   = "run-ps"
	QueryPsName = "query-ps"
)

// ---------------------------------------------------------------------------
// Resume strategy values
// ---------------------------------------------------------------------------

const (
	ResumeStrategySpawn    = "spawn"
	ResumeStrategyFastward = "fast"
)
