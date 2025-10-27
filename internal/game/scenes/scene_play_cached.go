package gamescene

import "github.com/hajimehoshi/ebiten/v2"

type CachedImage struct {
	image   *ebiten.Image
	imageOp *ebiten.DrawImageOptions
}

func (i *CachedImage) DrawTo(screen *ebiten.Image) {
	screen.DrawImage(i.image, i.imageOp)
}

func (i *CachedImage) DrawOver(img *ebiten.Image) {
	i.image.DrawImage(img, nil)
}
