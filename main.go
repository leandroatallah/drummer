package main

import (
	"embed"

	gamesetup "github.com/leandroatallah/drummer/internal/game/setup"
)

//go:embed assets
var assetsFs embed.FS

func main() {
	gamesetup.Setup(assetsFs)
}
