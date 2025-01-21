package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

type Display struct {
	image *image.RGBA
}

func NewDisplay(width, height int) *Display {
	return &Display{
		image: image.NewRGBA(image.Rect(0, 0, width, height)),
	}
}

func (g *Display) Update() error {
	return nil
}

func (g *Display) Draw(screen *ebiten.Image) {
	screen.WritePixels(g.image.Pix)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Display) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
