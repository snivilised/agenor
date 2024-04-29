package i18n

// CLIENT-TODO: Should be updated to use url of the implementing project,
// so should not be left as astrolib. (this should be set by auto-check)
const TraverseSourceID = "github.com/snivilised/traverse"

type traverseTemplData struct{}

func (td traverseTemplData) SourceID() string {
	return TraverseSourceID
}
