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
	// nav ghost command param-set
	NavPsName = "nav-ps"

	// exec ghost command param-set
	ExecPsName = "exec-ps"

	// shared families (registered on nav ghost, inherited by walk/sprint/query)
	PreviewFamName     = "preview-fam"
	CascadeFamName     = "cascade-fam"
	InteractionFamName = "interaction-fam"
	SamplingFamName    = "sampling-fam"

	// poly-filter family (registered on nav ghost, inherited by walk/sprint/query)
	PolyFamName = "poly-fam"

	// sprint-only family
	WorkerPoolFamName = "worker-pool-fam"

	// jay-specific param sets
	WalkPsName   = "walk-ps"
	SprintPsName = "sprint-ps"
	QueryPsName  = "query-ps"
)

// ---------------------------------------------------------------------------
// Resume strategy values
// ---------------------------------------------------------------------------

const (
	ResumeStrategySpawn    = "spawn"
	ResumeStrategyFastward = "fast"
)
