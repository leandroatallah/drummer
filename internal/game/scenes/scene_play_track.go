package gamescene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leandroatallah/firefly/internal/config"
)

type Track struct{}

type MainTrack struct {
	left  Track
	down  Track
	up    Track
	right Track
	scene *PlayScene
}

func NewMainTrack(scene *PlayScene) *MainTrack {
	return &MainTrack{scene: scene}
}

func (t *MainTrack) Update() error {
	return nil
}

// FIX: Long method.
func (t *MainTrack) Draw(screen *ebiten.Image) {
	s := t.scene

	cfg := config.Get()

	track := ebiten.NewImage(s.ui.trackWidth, s.ui.innerHeight)
	track.Fill(cfg.Colors.Dark)

	trackColWidth := (s.ui.trackWidth - (paddingX/2)*3) / 4
	trackLeft := ebiten.NewImage(trackColWidth, s.ui.innerHeight)
	trackDown := ebiten.NewImage(trackColWidth, s.ui.innerHeight)
	trackUp := ebiten.NewImage(trackColWidth, s.ui.innerHeight)
	trackRight := ebiten.NewImage(trackColWidth, s.ui.innerHeight)

	trackLeft.Fill(cfg.Colors.Light)
	trackDown.Fill(cfg.Colors.Light)
	trackUp.Fill(cfg.Colors.Light)
	trackRight.Fill(cfg.Colors.Light)

	// arrows bottom
	bottomBorderHeight := paddingY / 2
	bottomBorder := ebiten.NewImage(s.ui.trackWidth, bottomBorderHeight)
	bottomBorder.Fill(cfg.Colors.Dark)
	bottomBorderOp := &ebiten.DrawImageOptions{}
	bottomBorderOp.GeoM.Translate(0, float64(s.ui.innerHeight-trackColWidth-bottomBorderHeight))

	arrowsLightLeft, arrowsLightDown, arrowsLightUp, arrowsLightRight := GetArrows(arrowsLightImg)
	arrowsDarkLeft, arrowsDarkDown, arrowsDarkUp, arrowsDarkRight := GetArrows(arrowsDarkImg)

	arrowThemeMap := map[string][]*ebiten.Image{
		"left":  {arrowsLightLeft, arrowsDarkLeft},
		"down":  {arrowsLightDown, arrowsDarkDown},
		"up":    {arrowsLightUp, arrowsDarkUp},
		"right": {arrowsLightRight, arrowsDarkRight},
	}

	getArrowImageAndOp := func(direction string, isPressed bool) (*ebiten.Image, *ebiten.DrawImageOptions) {
		arrow := ebiten.NewImage(trackColWidth, trackColWidth)
		arrowsImg, found := arrowThemeMap[direction]
		if !found {
			log.Panic()
		}
		light := arrowsImg[0]
		dark := arrowsImg[1]

		if isPressed {
			arrow.Fill(cfg.Colors.Dark)
			DrawCenteredImage(arrow, dark)
		} else {
			arrow.Fill(cfg.Colors.Medium)
			DrawCenteredImage(arrow, light)
		}
		arrowOp := &ebiten.DrawImageOptions{}
		arrowOp.GeoM.Translate(0, float64(s.ui.innerHeight-trackColWidth))
		return arrow, arrowOp
	}

	arrowLeft, arrowLeftOp := getArrowImageAndOp("left", s.keyControl.isLeftPressed)
	arrowDown, arrowDownOp := getArrowImageAndOp("down", s.keyControl.isDownPressed)
	arrowUp, arrowUpOp := getArrowImageAndOp("up", s.keyControl.isUpPressed)
	arrowRight, arrowRightOp := getArrowImageAndOp("right", s.keyControl.isRightPressed)

	trackLeft.DrawImage(arrowLeft, arrowLeftOp)
	trackDown.DrawImage(arrowDown, arrowDownOp)
	trackUp.DrawImage(arrowUp, arrowUpOp)
	trackRight.DrawImage(arrowRight, arrowRightOp)

	trackLeftOp := &ebiten.DrawImageOptions{}
	trackDownOp := &ebiten.DrawImageOptions{}
	trackUpOp := &ebiten.DrawImageOptions{}
	trackRightOp := &ebiten.DrawImageOptions{}

	for i, op := range []*ebiten.DrawImageOptions{trackLeftOp, trackDownOp, trackUpOp, trackRightOp} {
		padding := 0
		if i > 0 {
			padding = paddingX / 2
		}
		op.GeoM.Translate(float64(trackColWidth*i+padding*i), 0)
	}

	// Draw moving arrows
	for _, n := range s.song.PlayingNotes {
		var lane *ebiten.Image
		var arrow *ebiten.Image

		switch n.Direction {
		case "left":
			arrow = arrowsDarkLeft
			lane = trackLeft
		case "down":
			arrow = arrowsDarkDown
			lane = trackDown
		case "up":
			arrow = arrowsDarkUp
			lane = trackUp
		case "right":
			arrow = arrowsDarkRight
			lane = trackRight
		default:
			continue
		}
		movingKey, movingKeyOp := newMovingKey(trackColWidth, arrow)

		progress := (s.song.GetPositionInBPM() - (float64(n.Onset) - s.song.offsetBpm)) / s.song.offsetBpm
		offsetY := progress*float64(s.ui.innerHeight) - float64(trackColWidth)

		movingKeyOp.GeoM.Translate(0, offsetY)
		lane.DrawImage(movingKey, movingKeyOp)
	}

	// Draw tracks to track container
	track.DrawImage(trackLeft, trackLeftOp)
	track.DrawImage(trackRight, trackRightOp)
	track.DrawImage(trackUp, trackUpOp)
	track.DrawImage(trackDown, trackDownOp)
	track.DrawImage(bottomBorder, bottomBorderOp)

	trackOp := &ebiten.DrawImageOptions{}
	trackOp.GeoM.Translate(paddingX, topRowHeight+paddingY*2)

	screen.DrawImage(track, trackOp)
}
