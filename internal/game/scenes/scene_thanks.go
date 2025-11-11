package gamescene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/drummer/internal/engine/assets"
	"github.com/leandroatallah/drummer/internal/engine/core"
	"github.com/leandroatallah/drummer/internal/engine/core/scene"
	"github.com/leandroatallah/drummer/internal/engine/core/transition"
)

var bgImg *ebiten.Image

type ThanksScene struct {
	scene.BaseScene
}

func NewThanksScene(context *core.AppContext) *ThanksScene {
	scene := ThanksScene{}
	scene.SetAppContext(context)
	return &scene
}

func (s *ThanksScene) OnStart() {
	bgImg = assets.LoadImageFromFs(s.AppContext, "assets/images/thank-you.png")

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
