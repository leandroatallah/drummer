package gamestate

import "github.com/leandroatallah/drummer/internal/engine/core/game/state"

const (
	Intro state.GameStateEnum = iota
	MainMenu
	Playing
	Paused
	GameOver
)
