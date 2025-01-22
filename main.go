package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pthm/gate/cpu"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	chip8 := cpu.NewCPU()
	renderer := NewImageRenderer()
	chip8.Renderer(renderer)

	flag.Parse()
	romPath := flag.Arg(0)

	if romPath == "" {
		fmt.Println("Must supply a path to a ROM")
		return
	}

	romBytes, err := os.ReadFile(romPath)
	if err != nil {
		fmt.Printf("Could not read ROM file at (%s): %v", romPath, err)
		return
	}
	chip8.LoadROM(romBytes)

	go chip8.Run(context.Background())

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Gate (CHIP-8 Emulator)")
	g := NewDisplay(renderer)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
