package gamescene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/drummer/internal/engine/assets/font"
	"github.com/leandroatallah/drummer/internal/engine/core"
	"github.com/leandroatallah/drummer/internal/engine/core/scene"
	"github.com/leandroatallah/drummer/internal/engine/core/transition"
)

const (
	bgSound = "assets/audio/black-sabbath-paranoid.mp3"
)

var pressStartImg *ebiten.Image

func init() {
	var err error
	pressStartImg, _, err = ebitenutil.NewImageFromFile("assets/images/press-start.png")
	if err != nil {
		log.Fatal(err)
	}
}

type MenuScene struct {
	scene.BaseScene

	count          int
	fontText       *font.FontText
	showPressStart bool
}

func NewMenuScene(context *core.AppContext) *MenuScene {
	scene := MenuScene{}
	scene.SetAppContext(context)
	return &scene
}

func (s *MenuScene) OnStart() {
	// Init audio
	s.AudioManager().PauseAll()
	if !s.AudioManager().IsPlayingSomething() {
		s.AudioManager().PlayMusic(bgSound)
	}
}

func (s *MenuScene) Update() error {
	s.count++

	if s.count%40 == 0 {
		s.showPressStart = !s.showPressStart
	}

	if !s.IsKeysDisabled && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.DisableKeys()
		s.Manager.NavigateTo(SceneTrackSelection, transition.NewFader(), false)
	}

	return nil
}

func (s *MenuScene) Draw(screen *ebiten.Image) {
	DrawCenteredImage(screen, pressStartImg)
}

func (s *MenuScene) OnFinish() {}
