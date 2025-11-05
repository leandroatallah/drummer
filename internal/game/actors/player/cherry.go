package gameplayer

import (
	"github.com/leandroatallah/drummer/internal/engine/actors"
	"github.com/leandroatallah/drummer/internal/engine/systems/physics"
)

type CherryPlayer struct {
	actors.Player

	coinCount int
}

func NewCherryPlayer(
	movementBlocker physics.PlayerMovementBlocker,
) (actors.PlayerEntity, error) {
	spriteData, statData, err := actors.ParseJsonPlayer("internal/game/actors/player/cherry.json")
	if err != nil {
		return nil, err
	}

	character, err := CreateAnimatedCharacter(spriteData)
	if err != nil {
		return nil, err
	}

	player := &CherryPlayer{
		Player: actors.Player{Character: *character},
	}
	SetPlayerBodies(player, spriteData)
	SetPlayerStats(player, statData)
	SetMovementModel(player, physics.Platform, movementBlocker)

	return player, nil
}

func (p *CherryPlayer) AddCoinCount(amount int) {
	p.coinCount += amount
}
func (p *CherryPlayer) CoinCount() int {
	return p.coinCount
}
