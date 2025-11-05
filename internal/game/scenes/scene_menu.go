package gamescene

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/leandroatallah/drummer/internal/config"
	"github.com/leandroatallah/drummer/internal/engine/assets/font"
	"github.com/leandroatallah/drummer/internal/engine/core"
	"github.com/leandroatallah/drummer/internal/engine/core/scene"
	"github.com/leandroatallah/drummer/internal/engine/core/transition"
)

const (
	bgSound = "assets/audio/black-sabbath-paranoid.ogg"
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
	fontText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	scene := MenuScene{fontText: fontText}
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

	if s.showPressStart {
		s.DrawPressStartText(screen)
	}
}

func (s *MenuScene) OnFinish() {}

func (s *MenuScene) DrawPressStartText(screen *ebiten.Image) {
	cfg := config.Get()

	type txt struct {
		offsetX float64
		offsetY float64
		color   color.RGBA
	}

	textMap := []txt{
		{
			offsetX: -1, offsetY: -1,
			color: cfg.Colors.Light,
		},
		{
			offsetX: 1, offsetY: 1,
			color: cfg.Colors.Light,
		},
		{
			offsetX: -1, offsetY: 1,
			color: cfg.Colors.Light,
		},
		{
			offsetX: 1, offsetY: -1,
			color: cfg.Colors.Light,
		},
		{
			offsetX: 1, offsetY: 0,
			color: cfg.Colors.Light,
		},
		{
			offsetX: 0, offsetY: 1,
			color: cfg.Colors.Light,
		},
		{
			offsetX: 0, offsetY: 0,
			color: cfg.Colors.Dark,
		},
	}

	for _, t := range textMap {
		textOp := &text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				PrimaryAlign:   text.AlignCenter,
				SecondaryAlign: text.AlignCenter,
				LineSpacing:    0,
			},
		}
		textOp.GeoM.Translate(
			float64(cfg.ScreenWidth/2)+t.offsetX,
			float64(cfg.ScreenHeight/2)+t.offsetY,
		)
		textOp.ColorScale.ScaleWithColor(t.color)
		s.fontText.Draw(screen, "Press Enter", 8, textOp)
	}
}
