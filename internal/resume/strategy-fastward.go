package resume

import "github.com/snivilised/traverse/internal/types"

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
