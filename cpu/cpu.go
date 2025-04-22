package cpu

import (
	"context"
	"fmt"
	"time"
)

type Renderer interface {
	Render(gfx [64][32]uint8) error
}

type CPU struct {
	opcode uint16      // Current opcode - Two bytes
	memory [4096]uint8 // Memory - 4KB
	v      [16]uint8   // Registers - 16 8-bit registers, V0-VE, VF (16th) is carry flag
	i      uint16      // Index register - 16-bit register
	pc     uint16      // Program counter - 16-bit register

	gfx      [64][32]uint8 // Graphics - 64x32 monochrome display
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

		gfx: [64][32]uint8{},

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
func (cpu *CPU) LoadROM(rom []uint8) error {
	if len(rom) > len(cpu.memory)-0x200 {
		return fmt.Errorf("ROM size exceeds available memory: %d bytes", len(rom))
	}
	for i, b := range rom {
		cpu.memory[0x200+i] = b
	}
	fmt.Printf("Successfully loaded ROM (%d bytes) into memory\n", len(rom))
	return nil
}

func (cpu *CPU) Reset() {
	// Reset the CPU state
}

func (cpu *CPU) Run(ctx context.Context) {
	clockTick := time.NewTicker(time.Second / 60) // 60Hz

	for {
		select {
		case <-ctx.Done():
			clockTick.Stop()
			return
		case <-clockTick.C:
			cpu.cycle()
			cpu.updateTimers()
			if cpu.drawFlag {
				if cpu.renderer != nil {
					cpu.renderer.Render(cpu.gfx)
				} else {
					for y := 0; y < 32; y++ {
						for x := 0; x < 64; x++ {
							fmt.Printf("%d", cpu.gfx[x][y])
						}
						fmt.Println()
					}
					fmt.Println()
				}
				cpu.drawFlag = false
			}
		}
	}
}

func (cpu *CPU) cycle() {
	// Fetch the opcode
	// TODO: understand if there is a way of doing this without the cast, or if it impacts performance
	cpu.opcode = uint16(cpu.memory[cpu.pc])<<8 | uint16(cpu.memory[cpu.pc+1])

	// The opcodes are 2 bytes long and are stored in big-endian format, this means that the most significant byte is stored first

	// We use binary masks to extract the different parts of the opcode
	// The opcode is 16 bits long, so we can use a 16-bit mask to extract the different parts
	// In the hexadecimal representation each digit is 4 bits long, so we can use a 4-bit mask to extract the different parts
	switch cpu.opcode & 0xF000 { // Get the first 4 bits
	case 0x0000: // The opcode starts with 0x0
		switch cpu.opcode & 0x000F { // Get the last 4 bits
		case 0x0000: // 0x00E0 - Clears the screen
			Op00E0(cpu)
		case 0x000E: // 0x00EE - Returns from a subroutine
			Op00EE(cpu)
		default:
			fmt.Printf("Unknown opcode [0x0000]: 0x%X\n", cpu.opcode)
		}
	case 0x1000: // 1NNN - Jumps to address NNN
		Op1NNN(cpu)
	case 0x2000: // 2NNN - Calls subroutine at NNN
		Op2NNN(cpu)
	case 0x3000: // 3XNN - Skips the next instruction if VX equals NN (usually the next instruction is a jump to skip a code block).
		Op3XNN(cpu)
	case 0x4000: // 4XNN - Skips the next instruction if VX does not equal NN (usually the next instruction is a jump to skip a code block)
		Op4XNN(cpu)
	case 0x5000: // 5XY0 - Skips the next instruction if VX equals VY (usually the next instruction is a jump to skip a code block)
		Op5XY0(cpu)
	case 0x6000: // 6XNN - Sets VX to NN
		Op6XNN(cpu)
	case 0x7000: // 7XNN - Adds NN to VX (carry flag is not changed)
		Op7XNN(cpu)
	case 0x8000: // Opcodes beginning with 8
		switch cpu.opcode & 0x000F { // Get the last 4 bits
		case 0x0000: // 0x8XY0 - Sets VX to the value of VY
			Op8XY0(cpu)
		case 0x0001: // 0x8XY1 - Sets VX to VX or VY
			Op8XY1(cpu)
		case 0x0002: // 0x8XY2 - Sets VX to VX and VY
			Op8XY2(cpu)
		case 0x0003: // 0x8XY3 - Sets VX to VX xor VY
			Op8XY3(cpu)
		case 0x0004: // 0x8XY4 - Adds VY to VF, when there is an overflow set the carry to 1
			Op8XY4(cpu)
		case 0x0005: // 0x8XY5 - VY is subtracted from VX. VF is set to 0 when there's an underflow, and 1 when there is not. (i.e. VF set to 1 if VX >= VY and 0 if not)
			Op8XY5(cpu)
		case 0x0006: // 0x8XY6 - Shifts VX right by one. VF is set to the least significant bit of VX before the shift.
			Op8XY6(cpu)
		case 0x0007: // 0x8XY7 - Sets VX to VY minus VX. VF is set to 0 when there's an underflow, and 1 when there is not. (i.e. VF set to 1 if VY >= VX and 0 if not)
			Op8XY7(cpu)
		case 0x000E: // 0x8XYE - Shifts VX left by one. VF is set to the most significant bit of VX before the shift.
			Op8XYE(cpu)
		default:
			fmt.Printf("Unknown opcode [0x8000]: 0x%X\n", cpu.opcode)
		}
	case 0x9000:
		Op9XY0(cpu)
	case 0xA000: //ANNN - Sets I to the address NNN
		OpANNN(cpu)
	case 0xB000:
		OpBNNN(cpu) // BNNN - Jumps to the address NNN plus V0
	case 0xC000:
		OpCXNN(cpu) // CXNN - Sets VX to the result of a bitwise and operation on a random number (Typically: 0 to 255) and NN
	case 0xD000: // DXYN - Draw a sprite at coordinate XY
		OpDXYN(cpu)
	case 0xF000: // Opcodes starting with 0xF
		switch cpu.opcode & 0x00FF {
		case 0x001E:
			OpFX1E(cpu)
		case 0x0033:
			OpFX33(cpu)
		default:
			fmt.Printf("Unknown opcode [0xF000]: 0x%X\n", cpu.opcode)
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

func (cpu *CPU) SetRenderer(renderer Renderer) {
	cpu.renderer = renderer
}
