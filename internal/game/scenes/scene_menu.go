package gamescene

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/drummer/internal/engine/assets"
	"github.com/leandroatallah/drummer/internal/engine/assets/font"
	"github.com/leandroatallah/drummer/internal/engine/core"
	"github.com/leandroatallah/drummer/internal/engine/core/scene"
	"github.com/leandroatallah/drummer/internal/engine/core/transition"
)

const (
	bgSound = "assets/audio/black-sabbath-paranoid.ogg"
)

var pressStartImg *ebiten.Image

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
	pressStartImg = assets.LoadImageFromFs(s.AppContext, "assets/images/press-start.png")
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
	frameOX, frameOY := 0, 0
	frameSprites := 2
	frameRate := 30
	width := pressStartImg.Bounds().Dx() / frameSprites
	height := pressStartImg.Bounds().Dy()

	elementWidth := pressStartImg.Bounds().Dx()
	frameCount := elementWidth / width
	i := (s.count / frameRate) % frameCount
	sx, sy := frameOX+i*width, frameOY

	res := pressStartImg.SubImage(
		image.Rect(sx, sy, sx+width, sy+height),
	).(*ebiten.Image)

	DrawCenteredImage(screen, res)
}

func (s *MenuScene) OnFinish() {}
