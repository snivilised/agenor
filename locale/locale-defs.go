package locale

const (
	SourceID = "github.com/snivilised/traverse"
)

type traverseTemplData struct{}

func (td traverseTemplData) SourceID() string {
	return SourceID
}
