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
	SceneThanks
)

func InitSceneMap(context *core.AppContext) navigation.SceneMap {
	sceneMap := navigation.SceneMap{
		SceneIntro: func() navigation.Scene {
			return NewIntroScene(context)
		},
		SceneMenu: func() navigation.Scene {
			return NewMenuScene(context)
		},
		ScenePlay: func() navigation.Scene {
			return NewPlayScene(context)
		},
		SceneTrackSelection: func() navigation.Scene {
			return NewTrackSelectionScene(context)
		},
		SceneThanks: func() navigation.Scene {
			return NewThanksScene(context)
		},
	}
	return sceneMap
}
