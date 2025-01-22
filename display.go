package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image"
)

const (
	chip8Width   = 64
	chip8Height  = 32
	scaleFactor  = 8
	screenWidth  = chip8Width * scaleFactor
	screenHeight = chip8Height * scaleFactor
)

type Display struct {
	image       *image.RGBA
	cpuRenderer *ImageRenderer
}

func NewDisplay(renderer *ImageRenderer) *Display {
	return &Display{
		image:       image.NewRGBA(image.Rect(0, 0, chip8Width, chip8Height)),
		cpuRenderer: renderer,
	}
}

func (g *Display) Update() error {
	// Add update logic if needed
	return nil
}

func (g *Display) Draw(screen *ebiten.Image) {
	img := g.cpuRenderer.Image() // Get the latest image.RGBA

	// Scale the image and draw it on the screen
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleFactor, scaleFactor) // Scale CHIP-8 resolution to 512x256
	screen.DrawImage(img, op)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.ActualTPS(), ebiten.ActualFPS()))
}

func (g *Display) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
