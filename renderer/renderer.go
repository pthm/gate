package renderer

import (
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"math"
)

type RaylibRenderer struct {
	gfx          [64][32]uint8
	scaleFactorX int32
	scaleFactorY int32
}

func NewRaylibRenderer(width, height int32) *RaylibRenderer {
	rl.InitWindow(width, height, "gate - chip8 emulator")
	rl.SetTargetFPS(60)

	scaleFactorX := math.Floor(float64(width) / 64.0)
	scaleFactorY := math.Floor(float64(height) / 32.0)

	return &RaylibRenderer{
		gfx:          [64][32]uint8{},
		scaleFactorX: int32(scaleFactorX),
		scaleFactorY: int32(scaleFactorY),
	}
}

func (r *RaylibRenderer) Run() {
	for !rl.WindowShouldClose() {
		fps := rl.GetFPS()

		rl.BeginDrawing()

		for y := 0; y < 32; y++ {
			for x := 0; x < 64; x++ {
				pixel := r.gfx[x][y]

				posX := int32(x) * r.scaleFactorX
				posY := int32(y) * r.scaleFactorY

				width := r.scaleFactorX
				height := r.scaleFactorY

				if pixel == 0 {
					rl.DrawRectangle(posX, posY, width, height, rl.Black)
				} else {
					rl.DrawRectangle(posX, posY, width, height, rl.White)
				}
			}
		}

		rl.DrawText(fmt.Sprintf("FPS: %d", fps), 10, 10, 10, rl.LightGray)

		rl.EndDrawing()
	}
}

func (r *RaylibRenderer) Render(gfx [64][32]uint8) error {
	r.gfx = gfx
	return nil
}

func (r *RaylibRenderer) Close() {
	rl.CloseWindow()
}
