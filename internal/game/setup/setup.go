package gamesetup

import (
	_ "image/png"
	"io/fs"
	"log"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/drummer/internal/config"
	"github.com/leandroatallah/drummer/internal/engine/actors"
	"github.com/leandroatallah/drummer/internal/engine/core"
	"github.com/leandroatallah/drummer/internal/engine/core/game"
	"github.com/leandroatallah/drummer/internal/engine/core/levels"
	"github.com/leandroatallah/drummer/internal/engine/core/scene"
	"github.com/leandroatallah/drummer/internal/engine/systems/audiomanager"
	"github.com/leandroatallah/drummer/internal/engine/systems/datamanager"
	"github.com/leandroatallah/drummer/internal/engine/systems/imagemanager"
	"github.com/leandroatallah/drummer/internal/engine/systems/input"
	gamescene "github.com/leandroatallah/drummer/internal/game/scenes"
)

func Setup(assets fs.FS) {
	// Basic Ebiten setup
	ebiten.SetWindowSize(config.Get().ScreenWidth*6, config.Get().ScreenHeight*6)
	ebiten.SetWindowTitle("The Drummer")

	// Initialize all systems and managers
	inputManager := input.NewManager()
	audioManager := audiomanager.NewAudioManager()
	imageManager := imagemanager.NewImageManager()
	dataManager := datamanager.NewDataManager()
	sceneManager := scene.NewSceneManager()
	levelManager := levels.NewManager()
	actorManager := actors.NewManager()

	// Load assets
	loadAudioAssetsFromFS(assets, audioManager)
	loadImageAssetsFromFS(assets, imageManager)
	loadDataAssetsFromFS(assets, dataManager)

	appContext := &core.AppContext{
		InputManager:    inputManager,
		AudioManager:    audioManager,
		ImageManager:    imageManager,
		DataManager:     dataManager,
		DialogueManager: nil,
		ActorManager:    actorManager,
		SceneManager:    sceneManager,
		LevelManager:    levelManager,
		// TODO: Rename this
		Assets: assets,
	}

	sceneFactory := scene.NewDefaultSceneFactory(gamescene.InitSceneMap(appContext))
	sceneFactory.SetAppContext(appContext)

	sceneManager.SetFactory(sceneFactory)
	sceneManager.SetAppContext(appContext)

	// Create and run the game
	game := game.NewGame(appContext)

	// Set initial game scene
	game.AppContext.SceneManager.NavigateTo(gamescene.SceneMenu, nil, false)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// loadAudioAssetsFromFS is a helper function to load all audio files from an fs.FS.
func loadAudioAssetsFromFS(assets fs.FS, am *audiomanager.AudioManager) {
	dir := "assets/audio"
	files, err := fs.ReadDir(assets, dir)
	if err != nil {
		log.Fatalf("error reading embedded audio dir: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		// Filter for supported audio types
		if !(strings.HasSuffix(fileName, ".ogg") || strings.HasSuffix(fileName, ".wav") || strings.HasSuffix(fileName, ".mp3")) {
			continue
		}

		fullPath := dir + "/" + fileName
		data, err := fs.ReadFile(assets, fullPath)
		if err != nil {
			log.Printf("failed to read embedded file %s: %v", fullPath, err)
			continue
		}

		// Use the existing Add method to process and store the player.
		am.Add(dir+"/"+fileName, data)
	}
}

// loadImageAssetsFromFS is a helper function to load all images files from an fs.FS.
func loadImageAssetsFromFS(assets fs.FS, m *imagemanager.ImageManager) {
	dir := "assets/images"
	files, err := fs.ReadDir(assets, dir)
	if err != nil {
		log.Fatalf("error reading embedded images dir: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			path := dir + "/" + file.Name()
			ebitenImg, _, err := ebitenutil.NewImageFromFileSystem(assets, path) // Use NewImageFromFileSystem directly
			if err != nil {
				log.Printf("error loading image file %s from FS: %v", file.Name(), err)
				continue
			}
			m.Add(file.Name(), ebitenImg)
		}
	}
}

// loadDataAssetsFromFS loads all .json files from the assets directory into the DataManager.
func loadDataAssetsFromFS(assets fs.FS, dm *datamanager.Manager) {
	err := fs.WalkDir(assets, "assets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".json") {
			data, err := fs.ReadFile(assets, path)
			if err != nil {
				log.Printf("error reading data file %s: %v", path, err)
				return nil // continue walking
			}
			fileName := filepath.Base(path)
			dm.Add(fileName, data)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("error walking data assets directory: %v", err)
	}
}
