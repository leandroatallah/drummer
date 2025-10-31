package gamescene

import (
	"github.com/leandroatallah/firefly/internal/engine/contracts/navigation"
	"github.com/leandroatallah/firefly/internal/engine/core"
)

const (
	SceneIntro navigation.SceneType = iota
	SceneMenu
	ScenePlay
	SceneTrackSelection
)

func InitSceneMap(context *core.AppContext) map[navigation.SceneType]navigation.Scene {
	SceneMap := map[navigation.SceneType]navigation.Scene{
		SceneIntro:          NewIntroScene(context),
		SceneMenu:           NewMenuScene(context),
		ScenePlay:           NewPlayScene(context),
		SceneTrackSelection: NewTrackSelectionScene(context),
	}
	return SceneMap
}
