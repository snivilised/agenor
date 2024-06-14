package kernel

import (
	"github.com/snivilised/traverse/internal/types"
)

// navigatorAgent does work on behalf of the navigator. It is distinct
// from navigatorBase and should only be used when the limited polymorphism
// on base is inadequate.
type navigatorAgent struct {
}

func newAgent() *navigatorAgent {
	return &navigatorAgent{}
}

func top(_ navigationStatic) (result types.NavigateResult, err error) {
	return types.NavigateResult{}, nil
}

func traverse(_ navigationStatic) {

}
