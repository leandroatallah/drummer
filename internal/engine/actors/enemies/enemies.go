package enemies

import (
	"github.com/leandroatallah/drummer/internal/engine/actors"
	"github.com/leandroatallah/drummer/internal/engine/contracts/body"
)

type BaseEnemy struct {
	actors.Character
}

func NewBaseEnemy() *BaseEnemy {
	return &BaseEnemy{}
}

// Character Methods
func (e *BaseEnemy) Update(space body.BodiesSpace) error {
	return e.Character.Update(space)
}
