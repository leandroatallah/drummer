package gamescene

import (
	"image"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/drummer/internal/engine/assets"
	"github.com/leandroatallah/drummer/internal/engine/assets/font"
	"github.com/leandroatallah/drummer/internal/engine/core"
	"github.com/leandroatallah/drummer/internal/engine/core/scene"
	"github.com/leandroatallah/drummer/internal/engine/core/transition"
	"github.com/leandroatallah/drummer/internal/engine/systems/audiomanager"
)

var selectionImg *ebiten.Image

type TrackSelectionScene struct {
	scene.BaseScene

	count          int
	audiomanager   *audiomanager.AudioManager
	fontText       *font.FontText
	showPressStart bool
}

func NewTrackSelectionScene(context *core.AppContext) *TrackSelectionScene {
	scene := TrackSelectionScene{}
	scene.SetAppContext(context)
	return &scene
}

func (s *TrackSelectionScene) OnStart() {
	s.audiomanager = s.Manager.AudioManager()

	selectionImg = assets.LoadImageFromFs(s.AppContext, "assets/images/track-selection.png")

	s.EnableKeys()
}

func (s *TrackSelectionScene) Update() error {
	s.count++

	if !s.IsKeysDisabled && ebiten.IsKeyPressed(ebiten.KeyEnter) {
		s.DisableKeys()
		s.Manager.NavigateTo(ScenePlay, transition.NewFader(), true)
	}

	return nil
}

func (s *TrackSelectionScene) Draw(screen *ebiten.Image) {
	frameOX, frameOY := 0, 0
	frameSprites := 3
	frameRate := 20
	width := selectionImg.Bounds().Dx() / frameSprites
	height := selectionImg.Bounds().Dy()

	elementWidth := selectionImg.Bounds().Dx()
	frameCount := elementWidth / width
	i := (s.count / frameRate) % frameCount
	sx, sy := frameOX+i*width, frameOY

	res := selectionImg.SubImage(
		image.Rect(sx, sy, sx+width, sy+height),
	).(*ebiten.Image)

	DrawCenteredImage(screen, res)
}

func (s *TrackSelectionScene) OnFinish() {
	s.audiomanager.FadeOut(bgSound, 1*time.Second)
}
