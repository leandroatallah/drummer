package gamestate

import "github.com/leandroatallah/drummer/internal/engine/core/game/state"

type GameOverState struct {
	state.BaseState
}

func (s *GameOverState) OnStart() {}
