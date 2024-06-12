package resume

import (
	"github.com/snivilised/traverse/internal/types"
)

type spawnStrategy struct {
	baseStrategy
}

func (s *spawnStrategy) init() {

}

func (s *spawnStrategy) attach() {

}

func (s *spawnStrategy) detach() {

}

func (s *spawnStrategy) resume() (*types.NavigateResult, error) {
	return &types.NavigateResult{}, nil
}

func (s *spawnStrategy) finish() error {
	return nil
}
