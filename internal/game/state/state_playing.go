package gamestate

import "github.com/leandroatallah/drummer/internal/engine/core/game/state"

type PlayingState struct {
	state.BaseState
}

func (s *PlayingState) OnStart() {}
