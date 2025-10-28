package gamescene

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/core/transition"
	gameplayer "github.com/leandroatallah/firefly/internal/game/actors/player"
)

const (
	thermometerLimit = 25

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
	streak         int
	thermometer    int // thermometer starts on 0 and can range from -10 to 10
	ui             *ScreenUI
	keyControl     *KeyControl
	mainTrack      *MainTrack
	song           *Song
	speed          float64

	// Cached images
	containerImage *CachedImage
}

func NewPlayScene(context *core.AppContext) *PlayScene {
	mainText, err := font.NewFontText(config.Get().MainFontFace)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: Should receive from somewhere (maybe context)
	scene := &PlayScene{
		BaseScene:   *scene.NewScene(),
		mainText:    mainText,
		ui:          NewScreenUI(),
		keyControl:  NewKeyControl(),
		speed:       2.0,
		thermometer: 0,
	}
	song := NewSong("internal/game/songs/smell-like-teen-spirit.json", scene)
	scene.song = song

	song.offsetBpm = 4 / scene.speed
	scene.mainTrack = NewMainTrack(scene)

	// scene.SetAppContext(context)
	return scene
}

func (s *PlayScene) OnStart() {
	s.BaseScene.OnStart()

	container := s.DrawScreen()
	containerOp := &ebiten.DrawImageOptions{}
	containerOp.GeoM.Translate(float64(s.ui.margin), float64(s.ui.margin))

	s.drawDrummer(container)
	s.mainTrack.Draw(container)
	s.drawStatusColumn(container)

	s.drawIllustration(container)
	s.drawScore(container)
	s.drawThermometer(container)

	s.containerImage = &CachedImage{container, containerOp}

	// Play sound
	s.AudioManager().PlaySound("assets/audio/" + s.song.Filename)
}

func (s *PlayScene) Update() error {
	s.count++

	s.handleKeyPress()

	s.handleRightKeys()

	s.mainTrack.Update()
	s.song.Update()

	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	cfg := config.Get()
	screen.Fill(cfg.Colors.Light)

	s.containerImage.DrawTo(screen)

	// Dynamic content
	container, containerOp := s.DrawContainer()

	s.drawIllustration(container)

	if true {
		s.drawDrummer(container)
		s.drawScore(container)
		s.drawThermometer(container)
	}

	s.mainTrack.Draw(container)

	screen.DrawImage(container, containerOp)
	s.containerImage.DrawOver(container)
}

func (s *PlayScene) OnFinish() {
	s.AudioManager().PauseMusic("assets/audio/" + s.song.Filename)
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

func (s *PlayScene) handleKeyPress() {
	s.keyControl.Reset()

	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		s.keyControl.PressLeft()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		s.keyControl.PressDown()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		s.keyControl.PressUp()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		s.keyControl.PressRight()
	}
}

func (s *PlayScene) handleRightKeys() {
	tolerance := 0.5

	if !s.keyControl.IsSomeKeyPressed() {
		return
	}

	hasAnyCorrect := false
	for _, n := range s.song.PlayingNotes {
		if n.skip {
			continue
		}

		position := float64(n.Onset) - s.song.GetPositionInBPM()
		if math.Abs(position) > tolerance {
			continue
		}

		switch {
		case s.keyControl.isLeftPressed && n.Direction == "left":
			s.IncreaseScore()
			hasAnyCorrect = true
			n.skip = true
		case s.keyControl.isDownPressed && n.Direction == "down":
			s.IncreaseScore()
			hasAnyCorrect = true
			n.skip = true
		case s.keyControl.isUpPressed && n.Direction == "up":
			s.IncreaseScore()
			hasAnyCorrect = true
			n.skip = true
		case s.keyControl.isRightPressed && n.Direction == "right":
			s.IncreaseScore()
			hasAnyCorrect = true
			n.skip = true
		}
	}

	if !hasAnyCorrect {
		s.handleMistake()
	}
}

func (s *PlayScene) IncreaseScore() {
	s.score += 5
	s.streak++
	s.thermometer++
	if s.thermometer > thermometerLimit {
		s.thermometer = thermometerLimit
	}
}

func (s *PlayScene) handleMistake() {
	s.streak = 0
	s.thermometer -= 2
	if s.thermometer < 0 {
		s.thermometer = 0
	}
}
