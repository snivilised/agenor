package locale

const (
	// SourceID defines the ID (by convemtion the repo URL) required
	// for i18n translation pureposes.
	SourceID = "github.com/snivilised/agenor"
)

type agenorTemplData struct{}

// SourceID returns the source ID for the template data.
func (td agenorTemplData) SourceID() string {
	return SourceID
}
