package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 64
	screenHeight = 32
)

func main() {

	cpu := NewCPU()

	flag.Parse()
	romPath := flag.Arg(0)

	if romPath == "" {
		fmt.Println("Must supply a path to a ROM")
		return
	}

	romBytes, err := os.ReadFile(romPath)
	if err != nil {
		fmt.Printf("Could not read ROM file at (%s): %w", romPath, err)
		return
	}
	cpu.LoadROM(romBytes)

	cpu.Run()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Gate (CHIP-8 Emulator)")
	g := NewDisplay(screenWidth, screenHeight)
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
