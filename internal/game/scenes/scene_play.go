package gamescene

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/core/transition"
	gameplayer "github.com/leandroatallah/firefly/internal/game/actors/player"
)

const (
	bgSound = "assets/audio/Sketchbook.ogg"

	// UI
	screenMargin      = 4
	paddingX          = 4
	paddingY          = 6
	topRowHeight      = 21
	leftColumnWidth   = 48
	statusBoxPadding  = 2
	scoreHeight       = 22
	thermometerHeight = 22
)

var (
	illustrationDark  *ebiten.Image
	illustrationLight *ebiten.Image
)

func init() {
	var err error
	illustrationLight, _, err = ebitenutil.NewImageFromFile("assets/images/illustration-light.png")
	if err != nil {
		log.Fatal(err)
	}
	illustrationDark, _, err = ebitenutil.NewImageFromFile("assets/images/illustration-dark.png")
	if err != nil {
		log.Fatal(err)
	}
}

type CachedImage struct {
	image   *ebiten.Image
	imageOp *ebiten.DrawImageOptions
}

func (i *CachedImage) DrawTo(screen *ebiten.Image) {
	screen.DrawImage(i.image, i.imageOp)
}

func (i *CachedImage) DrawOver(img *ebiten.Image) {
	i.image.DrawImage(img, nil)
}

type ScreenUI struct {
	margin          int
	containerWidth  int
	containerHeight int
	innerWidth      int
	innerHeight     int
	trackWidth      int

	textsImg *ebiten.Image
}

func NewScreenUI() *ScreenUI {
	cfg := config.Get()

	margin := 4
	width := cfg.ScreenWidth - (margin * 2)
	height := cfg.ScreenHeight - (margin * 2)
	innerWidth := width - (paddingX * 2)

	textsPath := "assets/images/texts.png"
	textsImg, _, err := ebitenutil.NewImageFromFile(textsPath)
	if err != nil {
		log.Fatal(err)
	}

	return &ScreenUI{
		margin:          margin,
		containerWidth:  width,
		containerHeight: height,
		innerWidth:      innerWidth,
		innerHeight:     height - (paddingY * 2) - topRowHeight - paddingY,
		trackWidth:      innerWidth - paddingY - leftColumnWidth,
		textsImg:        textsImg,
	}
}

type PlayScene struct {
	scene.BaseScene
	count          int
	mainText       *font.FontText
	levelCompleted bool
	score          int
	ui             *ScreenUI

	// Cached images
	containerImage *CachedImage
}

func NewPlayScene(context *core.AppContext) *PlayScene {
	mainText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}
	scene := PlayScene{
		BaseScene: *scene.NewScene(),
		mainText:  mainText,
		ui:        NewScreenUI(),
	}
	// scene.SetAppContext(context)
	return &scene
}

func (s *PlayScene) OnStart() {
	s.BaseScene.OnStart()

	container := DrawScreen(s.ui.containerWidth, s.ui.containerHeight)
	containerOp := &ebiten.DrawImageOptions{}
	containerOp.GeoM.Translate(float64(s.ui.margin), float64(s.ui.margin))

	s.drawDrummer(container)
	// TODO: Should draw each track separated?
	s.drawTrack(container)
	s.drawStatusColumn(container)

	s.drawIllustration(container, false)
	s.drawScore(container)
	s.drawThermometer(container)

	s.containerImage = &CachedImage{container, containerOp}

	// TODO: Is it working?
	// Play BG sound
	go func() {
		time.Sleep(1 * time.Second)
		s.AudioManager().PlaySound(bgSound)
	}()
}

func (s *PlayScene) Update() error {
	s.count++

	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	// TODO: Move all to OnStart method
	cfg := config.Get()
	screen.Fill(cfg.Colors.Light)

	// TODO: Split in different images to update when necessary
	s.containerImage.DrawTo(screen)

	// Dynamic content
	container, containerOp := s.DrawContainer()

	// TODO: Implement a solution to check if it should be updated
	if true {
		s.drawDrummer(container)
		s.drawIllustration(container, false)
		s.drawScore(container)
		s.drawThermometer(container)
		s.containerImage.DrawOver(container)
	}
	screen.DrawImage(container, containerOp)
}

func (s *PlayScene) OnFinish() {
	s.AudioManager().PauseMusic(bgSound)
}

func (s *PlayScene) finishLevel() {
	if s.levelCompleted {
		return
	}

	s.levelCompleted = true
	s.AppContext.SceneManager.NavigateTo(SceneMenu, transition.NewFader())
}

func createPlayer(appContext *core.AppContext) (actors.PlayerEntity, error) {
	p, err := gameplayer.NewCherryPlayer(appContext)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PlayScene) DrawContainer() (*ebiten.Image, *ebiten.DrawImageOptions) {
	container := ebiten.NewImage(s.ui.containerWidth, s.ui.containerHeight)
	containerOp := &ebiten.DrawImageOptions{}
	containerOp.GeoM.Translate(float64(s.ui.margin), float64(s.ui.margin))

	return container, containerOp
}
