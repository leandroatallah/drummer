package gamescene

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leandroatallah/firefly/internal/config"
)

func DrawScreen(width, height int) *ebiten.Image {
	cfg := config.Get()

	// container
	// TODO: Duplicated
	innerWidth := width - (paddingX * 2)
	container := ebiten.NewImage(width, height)
	container.Fill(cfg.Colors.Dark)

	// inner
	// TODO: Duplicated
	innerHeight := height - (paddingY * 2) - topRowHeight - paddingY

	inner := ebiten.NewImage(innerWidth, innerHeight)
	innerOp := &ebiten.DrawImageOptions{}
	innerOp.GeoM.Translate(
		float64(paddingX),
		float64(topRowHeight+(paddingY*2)),
	)

	container.DrawImage(inner, innerOp)

	return container
}

func (s *PlayScene) drawDrummer(screen *ebiten.Image) {
	cfg := config.Get()

	drummer := ebiten.NewImage(s.ui.innerWidth, topRowHeight)
	drummer.Fill(cfg.Colors.Medium)
	drummerOp := &ebiten.DrawImageOptions{}
	drummerOp.GeoM.Translate(float64(paddingX), float64(paddingY))

	drummerImgPath := "assets/images/drummer-idle.png"
	drummerImg, _, err := ebitenutil.NewImageFromFile(drummerImgPath)
	if err != nil {
		log.Fatal(err)
	}

	frameOX, frameOY := 0, 0
	frameRate := 20
	frameSprites := 2
	width := drummerImg.Bounds().Dx() / frameSprites
	height := drummerImg.Bounds().Dy()

	elementWidth := drummerImg.Bounds().Dx()
	frameCount := elementWidth / width
	i := (s.count / frameRate) % frameCount
	sx, sy := frameOX+i*width, frameOY

	res := drummerImg.SubImage(
		image.Rect(sx, sy, sx+width, sy+height),
	).(*ebiten.Image)

	drummer.DrawImage(res, nil)

	screen.DrawImage(drummer, drummerOp)
}

func (s *PlayScene) drawTrack(screen *ebiten.Image) {
	// s.ui.trackWidth, s.ui.innerHeight
	arrowsLightPath := "assets/images/light-arrows.png"
	arrowsDarkPath := "assets/images/dark-arrows.png"

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

	// arrows images
	arrowsLightImg, _, err := ebitenutil.NewImageFromFile(arrowsLightPath)
	if err != nil {
		log.Fatal(err)
	}
	arrowsDarkImg, _, err := ebitenutil.NewImageFromFile(arrowsDarkPath)
	if err != nil {
		log.Fatal(err)
	}

	arrowsLightLeft, arrowsLightDown, arrowsLightUp, arrowsLightRight := GetArrows(arrowsLightImg)
	arrowsDarkLeft, arrowsDarkDown, arrowsDarkUp, arrowsDarkRight := GetArrows(arrowsDarkImg)

	// arrows bottom
	bottomBorderHeight := paddingY / 2
	bottomBorder := ebiten.NewImage(s.ui.trackWidth, bottomBorderHeight)
	bottomBorder.Fill(cfg.Colors.Dark)
	bottomBorderOp := &ebiten.DrawImageOptions{}
	bottomBorderOp.GeoM.Translate(0, float64(s.ui.innerHeight-trackColWidth-bottomBorderHeight))

	arrowLeft := ebiten.NewImage(trackColWidth, trackColWidth)
	arrowLeft.Fill(cfg.Colors.Medium)
	arrowLeftOp := &ebiten.DrawImageOptions{}
	arrowLeftOp.GeoM.Translate(0, float64(s.ui.innerHeight-trackColWidth))
	DrawCenteredImage(arrowLeft, arrowsLightLeft)

	arrowDown := ebiten.NewImage(trackColWidth, trackColWidth)
	arrowDown.Fill(cfg.Colors.Medium)
	arrowDownOp := &ebiten.DrawImageOptions{}
	arrowDownOp.GeoM.Translate(0, float64(s.ui.innerHeight-trackColWidth))
	DrawCenteredImage(arrowDown, arrowsLightDown)

	arrowUp := ebiten.NewImage(trackColWidth, trackColWidth)
	arrowUp.Fill(cfg.Colors.Medium)
	arrowUpOp := &ebiten.DrawImageOptions{}
	arrowUpOp.GeoM.Translate(0, float64(s.ui.innerHeight-trackColWidth))
	DrawCenteredImage(arrowUp, arrowsLightUp)

	arrowRight := ebiten.NewImage(trackColWidth, trackColWidth)
	arrowRight.Fill(cfg.Colors.Medium)
	arrowRightOp := &ebiten.DrawImageOptions{}
	arrowRightOp.GeoM.Translate(0, float64(s.ui.innerHeight-trackColWidth))
	DrawCenteredImage(arrowRight, arrowsLightRight)

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
	movingLeftKey := newMovingKey(trackColWidth)
	movingDownKey := newMovingKey(trackColWidth)
	movingUpKey := newMovingKey(trackColWidth)
	movingRightKey := newMovingKey(trackColWidth)
	DrawCenteredImage(movingLeftKey, arrowsDarkLeft)
	DrawCenteredImage(movingDownKey, arrowsDarkDown)
	DrawCenteredImage(movingUpKey, arrowsDarkUp)
	DrawCenteredImage(movingRightKey, arrowsDarkRight)
	movingLeftKeyOp := &ebiten.DrawImageOptions{}
	movingDownKeyOp := &ebiten.DrawImageOptions{}
	movingUpKeyOp := &ebiten.DrawImageOptions{}
	movingRightKeyOp := &ebiten.DrawImageOptions{}

	movingDownKeyOp.GeoM.Translate(0, 20)
	movingRightKeyOp.GeoM.Translate(0, 40)

	trackLeft.DrawImage(movingLeftKey, movingLeftKeyOp)
	trackDown.DrawImage(movingDownKey, movingDownKeyOp)
	trackUp.DrawImage(movingUpKey, movingUpKeyOp)
	trackRight.DrawImage(movingRightKey, movingRightKeyOp)

	track.DrawImage(trackLeft, trackLeftOp)
	track.DrawImage(trackRight, trackRightOp)
	track.DrawImage(trackUp, trackUpOp)
	track.DrawImage(trackDown, trackDownOp)
	track.DrawImage(bottomBorder, bottomBorderOp)

	trackOp := &ebiten.DrawImageOptions{}
	trackOp.GeoM.Translate(paddingX, topRowHeight+paddingY*2)

	screen.DrawImage(track, trackOp)
}

func (s *PlayScene) drawStatusColumn(screen *ebiten.Image) {
	status := ebiten.NewImage(leftColumnWidth, s.ui.innerHeight)
	statusOp := &ebiten.DrawImageOptions{}
	statusOp.GeoM.Translate(float64(s.ui.trackWidth+paddingY+paddingX), topRowHeight+(paddingY*2))

	screen.DrawImage(status, statusOp)
}

func (s *PlayScene) drawScore(screen *ebiten.Image) {
	score := DrawStatusRectangle(leftColumnWidth, scoreHeight)
	scoreTitle := s.ui.textsImg.SubImage(image.Rect(0, 0, 29, 7)).(*ebiten.Image)
	scoreTitleOp := &ebiten.DrawImageOptions{}
	scoreTitleOp.GeoM.Translate(float64(statusBoxPadding), float64(statusBoxPadding))
	scoreAmount := ebiten.NewImage(35, 7)
	scoreStr := fmt.Sprintf("%06d", s.score)
	for i := 0; i < 6; i++ {
		v := scoreStr[i : i+1]
		n := GetImageNumber(s.ui.textsImg, string(v))
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(i*6), 0)
		scoreAmount.DrawImage(n, op)
	}
	scoreAmountOp := &ebiten.DrawImageOptions{}
	scoreAmountOp.GeoM.Translate(float64(statusBoxPadding), 11)
	score.DrawImage(scoreTitle, scoreTitleOp)
	score.DrawImage(scoreAmount, scoreAmountOp)
	scoreOp := &ebiten.DrawImageOptions{}
	scoreOp.GeoM.Translate(float64(s.ui.trackWidth+paddingX+paddingY), float64(topRowHeight+(paddingY*2)))

	screen.DrawImage(score, scoreOp)
}

func (s *PlayScene) drawThermometer(screen *ebiten.Image) {
	cfg := config.Get()
	thermometer := DrawStatusRectangle(leftColumnWidth, thermometerHeight)
	thermometerOp := &ebiten.DrawImageOptions{}
	thermometerOp.GeoM.Translate(
		float64(s.ui.trackWidth+paddingX+paddingY),
		float64(topRowHeight+scoreHeight+(paddingY*3)),
	)
	thermTitle := s.ui.textsImg.SubImage(image.Rect(0, 8, 29, 15)).(*ebiten.Image)
	thermOp := &ebiten.DrawImageOptions{}
	thermOp.GeoM.Translate(float64(statusBoxPadding), float64(statusBoxPadding))
	thermometer.DrawImage(thermTitle, thermOp)

	therm := s.count % 60 / 10
	blockWidth := 8
	blockHeight := 5
	block := ebiten.NewImage(blockWidth, blockHeight)
	blockHollow := ebiten.NewImage(blockWidth-2, blockHeight-2)
	blockHollow.Fill(cfg.Colors.Light)
	for i := 0; i < 5; i++ {
		if i < therm {
			block.Fill(cfg.Colors.Dark)
		} else {
			block.Fill(cfg.Colors.Medium)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(1, 1)
			block.DrawImage(blockHollow, op)
		}
		blockOp := &ebiten.DrawImageOptions{}
		blockOp.GeoM.Translate(float64(statusBoxPadding+blockWidth*i+i*1), float64(statusBoxPadding+9))
		thermometer.DrawImage(block, blockOp)
	}

	screen.DrawImage(thermometer, thermometerOp)
}

func (s *PlayScene) drawIllustration(screen *ebiten.Image, dark bool) {
	illustrationHeight := s.ui.innerHeight - scoreHeight - thermometerHeight - (paddingY * 2)

	cfg := config.Get()
	illustration := DrawStatusRectangle(leftColumnWidth, illustrationHeight)
	illustrationOp := &ebiten.DrawImageOptions{}
	illustrationOp.GeoM.Translate(
		float64(s.ui.trackWidth+paddingX+paddingY),
		float64(topRowHeight+scoreHeight+thermometerHeight+(paddingY*4)),
	)

	if s.count%2 == 0 {
		dark = true
	}

	img := illustrationLight
	if dark {
		illustration.Fill(cfg.Colors.Medium)
		img = illustrationDark
	}
	DrawCenteredImage(illustration, img)

	screen.DrawImage(illustration, illustrationOp)
}

func DrawStatusRectangle(width, height int) *ebiten.Image {
	cfg := config.Get()
	container := ebiten.NewImage(width, height)
	container.Fill(cfg.Colors.Light)
	squareSize := 2
	square := ebiten.NewImage(squareSize, squareSize)
	square.Fill(cfg.Colors.Dark)
	squareOp := &ebiten.DrawImageOptions{}
	squareOp.GeoM.Translate(float64(width-squareSize), 0)
	container.DrawImage(square, squareOp)

	return container
}

func DrawCenteredImage(screen *ebiten.Image, image *ebiten.Image) {
	imageOp := &ebiten.DrawImageOptions{}
	screenW := screen.Bounds().Dx()
	screenH := screen.Bounds().Dy()

	imgW := image.Bounds().Dx()
	imgH := image.Bounds().Dy()
	imageOp.GeoM.Translate(
		float64(screenW/2-imgW/2), float64(screenH/2-imgH/2),
	)
	screen.DrawImage(image, imageOp)
}

func GetImageNumber(src *ebiten.Image, n string) *ebiten.Image {
	var res image.Image
	switch n {
	case "0":
		res = src.SubImage(image.Rect(0, 16, 5, 23))
	case "1":
		res = src.SubImage(image.Rect(6, 16, 11, 23))
	case "2":
		res = src.SubImage(image.Rect(12, 16, 17, 23))
	case "3":
		res = src.SubImage(image.Rect(18, 16, 23, 23))
	case "4":
		res = src.SubImage(image.Rect(24, 16, 29, 23))
	case "5":
		res = src.SubImage(image.Rect(0, 24, 5, 31))
	case "6":
		res = src.SubImage(image.Rect(6, 24, 11, 31))
	case "7":
		res = src.SubImage(image.Rect(12, 24, 17, 31))
	case "8":
		res = src.SubImage(image.Rect(18, 24, 23, 31))
	case "9":
		res = src.SubImage(image.Rect(24, 24, 29, 31))
	default:
		return nil
	}

	return res.(*ebiten.Image)
}

func GetArrows(img *ebiten.Image) (*ebiten.Image, *ebiten.Image, *ebiten.Image, *ebiten.Image) {
	verticalArrowWidth, verticalArrowHeight := 8, 15
	horizontalArrowWidth, horizontalArrowHeight := 14, 9

	subX0, subX1 := 0, 0
	subX1 += verticalArrowWidth
	left := img.SubImage(image.Rect(
		subX0, 0, subX1, verticalArrowHeight,
	)).(*ebiten.Image)

	subX0 += verticalArrowWidth
	subX1 += horizontalArrowWidth
	down := img.SubImage(image.Rect(
		subX0, 0, subX1, horizontalArrowHeight,
	)).(*ebiten.Image)

	subX0 += horizontalArrowWidth
	subX1 += horizontalArrowWidth
	up := img.SubImage(image.Rect(
		subX0, 0, subX1, horizontalArrowHeight,
	)).(*ebiten.Image)

	subX0 += horizontalArrowWidth
	subX1 += verticalArrowWidth
	right := img.SubImage(image.Rect(
		subX0, 0, subX1, verticalArrowHeight,
	)).(*ebiten.Image)

	return left, down, up, right
}

func newMovingKey(size int) *ebiten.Image {
	cfg := config.Get()
	key := ebiten.NewImage(size, size)
	key.Fill(cfg.Colors.Medium)
	keyInner := ebiten.NewImage(size-2, size-2)
	keyInner.Fill(cfg.Colors.Dark)
	keyInnerOp := &ebiten.DrawImageOptions{}
	keyInnerOp.GeoM.Translate(1, 1)
	key.DrawImage(keyInner, keyInnerOp)

	return key
}
