package resume

import (
	"github.com/pkg/errors"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
)

type fastwardGuardianSealer struct {
}

func (g *fastwardGuardianSealer) Seal(top types.Link) error {
	if top.Role() == enums.RoleFastward {
		return errors.New("can't decorate, last item is sealed (fastward)")
	}
	return nil
}

func (g *fastwardGuardianSealer) IsSealed(top types.Link) bool {
	_ = top

	return false
}

type fastwardStrategy struct {
	baseStrategy
}

func (s *fastwardStrategy) init() {

}

func (s *fastwardStrategy) attach() {

}

func (s *fastwardStrategy) detach() {

}

func (s *fastwardStrategy) resume() (*types.NavigateResult, error) {
	return &types.NavigateResult{}, nil
}

func (s *fastwardStrategy) finish() error {
	return nil
}
