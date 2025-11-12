package gamescene

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/audio"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/leandroatallah/drummer/internal/config"
	"github.com/leandroatallah/drummer/internal/engine/actors"
	"github.com/leandroatallah/drummer/internal/engine/assets"
	"github.com/leandroatallah/drummer/internal/engine/assets/font"
	"github.com/leandroatallah/drummer/internal/engine/core"
	"github.com/leandroatallah/drummer/internal/engine/core/scene"
	"github.com/leandroatallah/drummer/internal/engine/core/transition"
	gameplayer "github.com/leandroatallah/drummer/internal/game/actors/player"
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
	drummerIdleImg    *ebiten.Image
	drummerRockImg    *ebiten.Image
	arrowsLightImg    *ebiten.Image
	arrowsDarkImg     *ebiten.Image
	textsImg          *ebiten.Image
)

type ScreenUI struct {
	margin          int
	containerWidth  int
	containerHeight int
	innerWidth      int
	innerHeight     int
	trackWidth      int

	textsImg *ebiten.Image
}

func NewScreenUI(ctx *core.AppContext) *ScreenUI {
	cfg := config.Get()

	margin := 4
	width := cfg.ScreenWidth - (margin * 2)
	height := cfg.ScreenHeight - (margin * 2)
	innerWidth := width - (paddingX * 2)

	return &ScreenUI{
		margin:          margin,
		containerWidth:  width,
		containerHeight: height,
		innerWidth:      innerWidth,
		innerHeight:     height - (paddingY * 2) - topRowHeight - paddingY,
		trackWidth:      innerWidth - paddingY - leftColumnWidth,
		textsImg:        assets.LoadImageFromFs(ctx, "assets/images/texts.png"),
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
	songPlayer     *audio.Player
	isOver         bool

	// Caching layers for draw optimization
	staticLayer         *ebiten.Image
	scoreLayer          *ebiten.Image
	thermometerLayer    *ebiten.Image
	illustrationLayer   *ebiten.Image
	isScoreDirty        bool
	isThermometerDirty  bool
	isIllustrationDirty bool
}

func NewPlayScene(context *core.AppContext) *PlayScene {
	// TODO: Should receive from somewhere (maybe context)
	scene := &PlayScene{
		BaseScene:   *scene.NewScene(),
		ui:          NewScreenUI(context),
		keyControl:  NewKeyControl(),
		speed:       2.0,
		thermometer: 0,
	}

	songData := context.DataManager.Get("smell-like-teen-spirit.json")
	song := NewSongFromData(songData, scene)

	scene.song = song

	song.offsetBpm = 4 / scene.speed
	scene.mainTrack = NewMainTrack(scene)

	// scene.SetAppContext(context)
	return scene
}

func (s *PlayScene) OnStart() {
	s.BaseScene.OnStart()
	cfg := config.Get()

	// Init images
	illustrationLight = assets.LoadImageFromFs(s.AppContext, "assets/images/illustration-light.png")
	illustrationDark = assets.LoadImageFromFs(s.AppContext, "assets/images/illustration-dark.png")
	drummerIdleImg = assets.LoadImageFromFs(s.AppContext, "assets/images/drummer-idle.png")
	drummerRockImg = assets.LoadImageFromFs(s.AppContext, "assets/images/drummer-rock.png")
	arrowsLightImg = assets.LoadImageFromFs(s.AppContext, arrowsLightPath)
	arrowsDarkImg = assets.LoadImageFromFs(s.AppContext, arrowsDarkPath)

	// --- Initialize Layers ---
	s.staticLayer = ebiten.NewImage(cfg.ScreenWidth, cfg.ScreenHeight)
	s.scoreLayer = ebiten.NewImage(leftColumnWidth, scoreHeight)
	s.thermometerLayer = ebiten.NewImage(leftColumnWidth, thermometerHeight)
	illustrationHeight := s.ui.innerHeight - scoreHeight - thermometerHeight - (paddingY * 2)
	s.illustrationLayer = ebiten.NewImage(leftColumnWidth, illustrationHeight)

	// --- Pre-render Static Backgrounds ---
	s.staticLayer.Fill(cfg.Colors.Medium) // Main screen background

	// The main UI container
	container := s.DrawScreen()
	containerOp := &ebiten.DrawImageOptions{}
	containerOp.GeoM.Translate(float64(s.ui.margin), float64(s.ui.margin))

	// The drummer area background
	drummerBg := ebiten.NewImage(s.ui.innerWidth, topRowHeight)
	drummerBg.Fill(cfg.Colors.Medium)
	drummerBgOp := &ebiten.DrawImageOptions{}
	drummerBgOp.GeoM.Translate(float64(paddingX), float64(paddingY))
	container.DrawImage(drummerBg, drummerBgOp)

	// The status column background
	s.drawStatusColumn(container)

	// Draw the fully prepared static container to the static layer
	s.staticLayer.DrawImage(container, containerOp)

	// --- Set Dirty Flags for First Render ---
	s.isScoreDirty = true
	s.isThermometerDirty = true
	s.isIllustrationDirty = true
}

func (s *PlayScene) Update() error {
	s.count++

	// Wait menu sound end before start
	if s.songPlayer == nil && !s.AudioManager().IsPlayingSomething() {
		s.AudioManager().SetVolume(1)
		s.songPlayer = s.AudioManager().PlaySound("assets/audio/" + s.song.Filename)
	}

	// The soung is over
	if !s.isOver && s.songPlayer != nil && !s.songPlayer.IsPlaying() {
		s.isOver = true
		s.DisableKeys()
		s.AppContext.SceneManager.NavigateTo(SceneThanks, transition.NewFader(), true)
	}

	if s.songPlayer != nil && s.songPlayer.IsPlaying() {
		s.handleKeyPress()
		s.handleRightKeys()
		s.mainTrack.Update()
		s.song.Update()
	}

	return nil
}

func (s *PlayScene) Draw(screen *ebiten.Image) {
	// 1. Draw the static background, which is already composed.
	screen.DrawImage(s.staticLayer, nil)

	// --- Handle semi-static layers that need updating ---

	// Redraw score only if it has changed.
	if s.isScoreDirty {
		s.redrawScoreLayer()
		s.isScoreDirty = false
	}

	// Redraw thermometer only if it has changed.
	if s.isThermometerDirty {
		s.redrawThermometerLayer()
		s.isThermometerDirty = false
	}

	// Redraw illustration only if it has changed.
	s.redrawIllustrationLayer()

	// --- Draw the cached layers to the screen at their correct positions ---
	containerOriginX := float64(s.ui.margin)
	containerOriginY := float64(s.ui.margin)

	scoreOp := &ebiten.DrawImageOptions{}
	scoreOp.GeoM.Translate(containerOriginX+float64(s.ui.trackWidth+paddingX+paddingY), containerOriginY+float64(topRowHeight+(paddingY*2)))
	screen.DrawImage(s.scoreLayer, scoreOp)

	thermometerOp := &ebiten.DrawImageOptions{}
	thermometerOp.GeoM.Translate(containerOriginX+float64(s.ui.trackWidth+paddingX+paddingY), containerOriginY+float64(topRowHeight+scoreHeight+(paddingY*3)))
	screen.DrawImage(s.thermometerLayer, thermometerOp)

	illustrationOp := &ebiten.DrawImageOptions{}
	illustrationOp.GeoM.Translate(containerOriginX+float64(s.ui.trackWidth+paddingX+paddingY), containerOriginY+float64(topRowHeight+scoreHeight+thermometerHeight+(paddingY*4)))
	screen.DrawImage(s.illustrationLayer, illustrationOp)

	// --- Draw fully dynamic elements directly on top ---
	// We draw them into a temporary transparent image so their local coordinates match the container.
	dynamicContainer := ebiten.NewImage(s.ui.containerWidth, s.ui.containerHeight)
	s.drawDrummer(dynamicContainer)
	s.mainTrack.Draw(dynamicContainer)

	// Draw the dynamic container onto the screen.
	dynamicContainerOp := &ebiten.DrawImageOptions{}
	dynamicContainerOp.GeoM.Translate(float64(s.ui.margin), float64(s.ui.margin))
	screen.DrawImage(dynamicContainer, dynamicContainerOp)
}

func (s *PlayScene) OnFinish() {
	if s.songPlayer != nil {
		s.songPlayer.Pause()
	}
}

func createPlayer(appContext *core.AppContext) (actors.PlayerEntity, error) {
	p, err := gameplayer.NewCherryPlayer(appContext)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// DrawContainer was removed as it's no longer used by the optimized Draw method.

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

	s.isScoreDirty = true
	s.isThermometerDirty = true
	s.isIllustrationDirty = true
}

func (s *PlayScene) handleMistake() {
	s.streak = 0
	s.thermometer -= 2
	if s.thermometer < 0 {
		s.thermometer = 0
	}

	s.isThermometerDirty = true
	s.isIllustrationDirty = true
}
