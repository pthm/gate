package main

import (
	"fmt"
	"time"
)

type Renderer interface {
	Render(gfx [64 * 32]uint8)
}

type CPU struct {
	opcode uint16      // Current opcode - Two bytes
	memory [4096]uint8 // Memory - 4KB
	v      [16]uint8   // Registers - 16 8-bit registers, V0-VE, VF (16th) is carry flag
	i      uint16      // Index register - 16-bit register
	pc     uint16      // Program counter - 16-bit register

	gfx      [64 * 32]uint8 // Graphics - 64x32 monochrome display
	drawFlag bool
	renderer Renderer

	delayTimer uint8 // Delay timer, decrements at 60Hz
	soundTimer uint8 // Sound timer, decrements at 60Hz and buzzes when zero

	stack [16]uint16 // Stack - 16 levels
	sp    uint16     // Stack pointer
}

var fontset = [80]uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

func NewCPU() *CPU {
	cpu := &CPU{
		memory: [4096]uint8{},
		v:      [16]uint8{},
		i:      0,
		pc:     0x200, // Program counter starts at 0x200

		gfx: [64 * 32]uint8{},

		delayTimer: 0,
		soundTimer: 0,

		stack: [16]uint16{},
		sp:    0,
	}

	// Initialize memory map
	// 0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
	// 0x050-0x0A0 - Used for the built-in 4x5 pixel font set (0-F)
	// 0x200-0xFFF - Program ROM and work RAM

	// Load fontset
	for i, b := range fontset {
		cpu.memory[i] = b
	}

	return cpu
}

// LoadROM loads a ROM into memory, starting at 0x200
func (cpu *CPU) LoadROM(rom []uint8) {
	fmt.Println("Loading ROM")
	for i, b := range rom {
		cpu.memory[0x200+i] = b
	}
}

func (cpu *CPU) Reset() {
	// Reset the CPU state
}

func (cpu *CPU) Run() {
	clockTick := time.NewTicker(time.Second / 60) // 60Hz
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-clockTick.C:
				cpu.cycle()
				cpu.updateTimers()
				if cpu.drawFlag && cpu.renderer != nil {
					cpu.renderer.Render(cpu.gfx)
					cpu.drawFlag = false
				}
			}
		}
	}()

	// TODO: handle the lifecycle of the CPU externally
	time.Sleep(30 * time.Second)
	clockTick.Stop()
	done <- true
}

func (cpu *CPU) cycle() {
	// Fetch the opcode
	// TODO: understand if there is a way of doing this without the cast, or if it impacts performance
	cpu.opcode = uint16(cpu.memory[cpu.pc])<<8 | uint16(cpu.memory[cpu.pc+1])

	fmt.Printf("= Cycle | PC 0x%X | Opcode 0x%X\n", cpu.pc, cpu.opcode)

	// The opcodes are 2 bytes long and are stored in big-endian format, this means that the most significant byte is stored first

	// We use binary masks to extract the different parts of the opcode
	// The opcode is 16 bits long, so we can use a 16-bit mask to extract the different parts
	// In the hexadecimal representation each digit is 4 bits long, so we can use a 4-bit mask to extract the different parts
	switch cpu.opcode & 0xF000 { // Get the first 4 bits
	case 0x0000: // The opcode starts with 0x0
		switch cpu.opcode & 0x000F { // Get the last 4 bits
		case 0x0000: // 0x00E0 - Clears the screen
			op00E0(cpu)
		case 0x000E: // 0x00EE - Returns from a subroutine
			fmt.Printf("Unimplemented opcode [0x00EE]: 0x%X\n", cpu.opcode)
			//TODO: Implement
		default:
			fmt.Printf("Unknown opcode [0x0000]: 0x%X\n", cpu.opcode)
		}
	case 0x2000: //2NNN - Calls subroutine at NNN
		op2NNN(cpu)
	case 0x8000: // Opcodes beginning with 8
		switch cpu.opcode & 0x000F { // Get the last 4 bits
		case 0x0004: // 0x8XY4 - Adds VY to VF, when there is an overflow set the carry to 1
			op8XY4(cpu)
		}
	case 0xA000: //ANNN - Sets I to the address NNN
		opANNN(cpu)
	case 0xD000: // DXYN - Draw a sprite at coordinate XY
		opDXYN(cpu)
	case 0xF000: // Opcodes starting with 0xF
		switch cpu.opcode & 0x00FF {
		case 0x0033:
			opFX33(cpu)
		}
	default:
		fmt.Printf("Unknown opcode: 0x%X\n", cpu.opcode)
	}

}

func (cpu *CPU) updateTimers() {
	if cpu.delayTimer > 0 {
		cpu.delayTimer--
	}
	if cpu.soundTimer > 0 {
		if cpu.soundTimer == 1 {
			fmt.Println("Beep!")
		}
		cpu.soundTimer--
	}
}

func (cpu *CPU) fetchOpcode() {
}

func (cpu *CPU) decodeOpcode() {

}
