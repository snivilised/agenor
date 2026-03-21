package locale

const (
	// SourceID defines the ID (by convention the repo URL) required
	// for i18n translation purposes.
	SourceID = "github.com/snivilised/jaywalk/src/agenor"
)

type agenorTemplData struct{}

// SourceID returns the source ID for the template data.
func (td agenorTemplData) SourceID() string {
	return SourceID
}
