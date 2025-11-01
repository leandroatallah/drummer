package gamesetup

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/config"
	"github.com/leandroatallah/firefly/internal/engine/actors"
	"github.com/leandroatallah/firefly/internal/engine/assets/font"
	"github.com/leandroatallah/firefly/internal/engine/core"
	"github.com/leandroatallah/firefly/internal/engine/core/game"
	"github.com/leandroatallah/firefly/internal/engine/core/levels"
	"github.com/leandroatallah/firefly/internal/engine/core/scene"
	"github.com/leandroatallah/firefly/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/imagemanager"
	"github.com/leandroatallah/firefly/internal/engine/systems/input"
	"github.com/leandroatallah/firefly/internal/engine/systems/speech"
	gamescene "github.com/leandroatallah/firefly/internal/game/scenes"
	gamespeech "github.com/leandroatallah/firefly/internal/game/speech"
)

func Setup() {
	// Basic Ebiten setup
	ebiten.SetWindowSize(config.Get().ScreenWidth*4, config.Get().ScreenHeight*4)
	ebiten.SetWindowTitle("Firefly")

	// Initialize all systems and managers
	inputManager := input.NewManager()
	audioManager := audiomanager.NewAudioManager()
	imageManager := imagemanager.NewImageManager()
	sceneManager := scene.NewSceneManager()
	levelManager := levels.NewManager()
	actorManager := actors.NewManager()

	// Initialize Dialogue Manager
	fontText, err := font.NewFontText("assets/fonts/Silkscreen-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}
	speechFont := speech.NewSpeechFont(fontText, 8, 14)
	speechBubble := gamespeech.NewSpeechBubble(speechFont)
	dialogueManager := speech.NewManager(speechBubble)

	// Load assets
	loadAudioAssets(audioManager)

	appContext := &core.AppContext{
		InputManager:    inputManager,
		AudioManager:    audioManager,
		ImageManager:    imageManager,
		DialogueManager: dialogueManager,
		ActorManager:    actorManager,
		SceneManager:    sceneManager,
		LevelManager:    levelManager,
	}

	sceneFactory := scene.NewDefaultSceneFactory(gamescene.InitSceneMap(appContext))
	sceneFactory.SetAppContext(appContext)

	sceneManager.SetFactory(sceneFactory)
	sceneManager.SetAppContext(appContext)

	// Create and run the game
	game := game.NewGame(appContext)

	// Set initial game scene
	game.AppContext.SceneManager.NavigateTo(gamescene.SceneMenu, nil)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// loadAudioAssets is a helper function to load all audio files from the assets directory.
func loadAudioAssets(am *audiomanager.AudioManager) {
	files, err := os.ReadDir("assets/audio")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".ogg") || strings.HasSuffix(file.Name(), ".wav")) {
			path := "assets/audio/" + file.Name()
			audioItem, err := am.Load(path)
			if err != nil {
				log.Printf("error loading audio file %s: %v", file.Name(), err)
				continue
			}
			am.Add(audioItem.Name(), audioItem.Data())
		}
	}
}

// loadImageAssets is a helper function to load all images files from the assets directory.
func loadImageAssets(m *imagemanager.ImageManager) {
	files, err := os.ReadDir("assets/images")
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			path := "assets/images/" + file.Name()
			img, _, err := ebitenutil.NewImageFromFile(path)
			if err != nil {
				log.Printf("error loading image file %s: %v", file.Name(), err)
				continue
			}
			m.Add(file.Name(), img)
		}
	}
}
