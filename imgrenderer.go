package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageRenderer struct {
	image *ebiten.Image
}

func NewImageRenderer() *ImageRenderer {
	return &ImageRenderer{
		image: ebiten.NewImage(64, 32),
	}
}

func (g *ImageRenderer) Render(gfx [64][32]uint8) error {
	width := 64
	height := 32
	pixels := make([]byte, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if gfx[x][y] == 0 {
				// Set the pixel to black, row-major order by index
				pixels[(y*width+x)*4] = 0
				pixels[(y*width+x)*4+1] = 0
				pixels[(y*width+x)*4+2] = 0
				pixels[(y*width+x)*4+3] = 255
			} else {
				// Set the pixel to white, row-major order by index
				pixels[(y*width+x)*4] = 255
				pixels[(y*width+x)*4+1] = 255
				pixels[(y*width+x)*4+2] = 255
				pixels[(y*width+x)*4+3] = 255
			}
		}
	}
	g.image.WritePixels(pixels)
	return nil
}

func (g *ImageRenderer) Image() *ebiten.Image {
	return g.image
}
