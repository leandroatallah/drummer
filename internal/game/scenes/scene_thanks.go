package gamescene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/drummer/internal/engine/core"
	"github.com/leandroatallah/drummer/internal/engine/core/scene"
	"github.com/leandroatallah/drummer/internal/engine/core/transition"
)

var bgImg *ebiten.Image

func init() {
	var err error
	bgImg, _, err = ebitenutil.NewImageFromFile("assets/images/thank-you.png")
	if err != nil {
		log.Fatal(err)
	}
}

type ThanksScene struct {
	scene.BaseScene
}

func NewThanksScene(context *core.AppContext) *ThanksScene {
	scene := ThanksScene{}
	scene.SetAppContext(context)
	return &scene
}

func (s *ThanksScene) OnStart() {
	s.AudioManager().PauseAll()
	s.AudioManager().PlaySound(bgSound)
}

func (s *ThanksScene) Update() error {
	if !s.IsKeysDisabled && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.DisableKeys()
		s.Manager.NavigateTo(SceneMenu, transition.NewFader(), true)
	}

	return nil
}

func (s *ThanksScene) Draw(screen *ebiten.Image) {
	DrawCenteredImage(screen, bgImg)
}

func (s *ThanksScene) OnFinish() {}
