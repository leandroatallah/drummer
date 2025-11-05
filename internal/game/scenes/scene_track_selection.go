package gamescene

import (
	"image"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/drummer/internal/config"
	"github.com/leandroatallah/drummer/internal/engine/assets/font"
	"github.com/leandroatallah/drummer/internal/engine/core"
	"github.com/leandroatallah/drummer/internal/engine/core/scene"
	"github.com/leandroatallah/drummer/internal/engine/core/transition"
	"github.com/leandroatallah/drummer/internal/engine/systems/audiomanager"
)

var selectionImg *ebiten.Image

func init() {
	var err error
	selectionImg, _, err = ebitenutil.NewImageFromFile("assets/images/track-selection.png")
	if err != nil {
		log.Fatal(err)
	}
}

type TrackSelectionScene struct {
	scene.BaseScene

	count          int
	audiomanager   *audiomanager.AudioManager
	fontText       *font.FontText
	showPressStart bool
}

func NewTrackSelectionScene(context *core.AppContext) *TrackSelectionScene {
	fontText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	scene := TrackSelectionScene{fontText: fontText}
	scene.SetAppContext(context)
	return &scene
}

func (s *TrackSelectionScene) OnStart() {
	s.audiomanager = s.Manager.AudioManager()

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
