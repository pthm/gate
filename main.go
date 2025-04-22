package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pthm/gate/cpu"
	"github.com/pthm/gate/renderer"
	"os"
)

func main() {

	chip8 := cpu.NewCPU()

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

	rlRenderer := renderer.NewRaylibRenderer(1400, 450)
	chip8.SetRenderer(rlRenderer)

	go chip8.Run(context.Background())
	rlRenderer.Run()

	defer rlRenderer.Close()
}
