package locale

const (
	TraverseSourceID = "github.com/snivilised/traverse"
)

type traverseTemplData struct{}

func (td traverseTemplData) SourceID() string {
	return TraverseSourceID
}
