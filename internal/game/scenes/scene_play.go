package gamescene

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
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
	margin   = 4
	paddingX = 4
	paddingY = 6
)

type CachedImage struct {
	image   *ebiten.Image
	imageOp *ebiten.DrawImageOptions
}

func (i *CachedImage) Draw(screen *ebiten.Image) {
	screen.DrawImage(i.image, i.imageOp)
}

type PlayScene struct {
	scene.BaseScene
	count          int
	mainText       *font.FontText
	levelCompleted bool
	score          int

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
	}
	// scene.SetAppContext(context)
	return &scene
}

func (s *PlayScene) OnStart() {
	s.BaseScene.OnStart()

	container, containerOp := DrawScreen()
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
	s.containerImage.Draw(screen)
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
