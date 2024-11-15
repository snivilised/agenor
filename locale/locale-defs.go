package locale

const (
	SourceID = "github.com/snivilised/agenor"
)

type agenorTemplData struct{}

func (td agenorTemplData) SourceID() string {
	return SourceID
}
