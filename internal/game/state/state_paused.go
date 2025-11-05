package gamestate

import "github.com/leandroatallah/drummer/internal/engine/core/game/state"

type PausedState struct {
	state.BaseState
}

func (s *PausedState) OnStart() {}
