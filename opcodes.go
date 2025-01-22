package main

import "fmt"

// 2NNN - Calls subroutine at NNN
func op2NNN(cpu *CPU) {
	if cpu.sp >= uint16(len(cpu.stack)) {
		fmt.Println("Stack overflow")
		return // Prevent overflow
	}
	cpu.stack[cpu.sp] = cpu.pc   // Store the program counter value in the stack at the current stack pointer
	cpu.sp++                     // Increment the stack pointer
	cpu.pc = cpu.opcode & 0x0FFF // Set the program counter to the address NNN, we use the mask 0x0FFF to extract NNN
}

// 8XY4 - Adds VY to VX. VF is set to 1 when there's an overflow, and to 0 when there is not
func op8XY4(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := (cpu.opcode & 0x00F0) >> 4 // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	// If VY is greater than the result of 255 - VY then there is an overflow, and we should set the carry to 1
	// The registers are 8 bits, by doing 255 minus VY we are calculating how much "remaining space" there is
	// in the register so we can determine if there will be an overflow
	if cpu.v[y] > (0xFF - cpu.v[x]) {
		cpu.v[0xF] = 1 // Set carry to 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] += cpu.v[y]
	cpu.pc += 2
}

// ANNN - Sets I to the address NNN
func opANNN(cpu *CPU) {
	cpu.i = cpu.opcode & 0x0FFF // Set I to NNN
	cpu.pc += 2                 // Increment the program counter by two
}

// FX33 - Stores the binary-coded decimal representation of VX, with the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
func opFX33(cpu *CPU) {
	x := uint8((cpu.opcode & 0x0F00) >> 8)
	vx := cpu.v[x]
	cpu.memory[cpu.i] = vx / 100
	cpu.memory[cpu.i+1] = (vx / 10) % 10
	cpu.memory[cpu.i+2] = vx % 10
	cpu.pc += 2
}

// 00E0 - Clear screen
func op00E0(cpu *CPU) {
	cpu.gfx = [64][32]uint8{}
	cpu.pc += 2
}

// 1NNN - Jumps to address NNN
func op1NNN(cpu *CPU) {
	cpu.pc = cpu.opcode & 0x0FFF // Set the program counter to the address NNN, we use the mask 0x0FFF to extract NNN
}

// 6XNN - Sets VX to NN.
func op6XNN(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	nn := cpu.opcode & 0x00FF       // Fetch NN from the opcode

	cpu.v[x] = uint8(nn)
	cpu.pc += 2
}

// 7XNN - Adds NN to VX (carry flag is not changed).
func op7XNN(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	nn := cpu.opcode & 0x00FF       // Fetch NN from the opcode

	cpu.v[x] += uint8(nn)
	cpu.pc += 2
}

// DXYN - Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.
// Each row of 8 pixels is read as bit-coded starting from memory location I; I value does not change after the execution of this instruction.
// As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that does not happen
func opDXYN(cpu *CPU) {
	vx := cpu.v[(cpu.opcode&0x0F00)>>8]   // VX (x-coordinate from register)
	vy := cpu.v[(cpu.opcode&0x00F0)>>4]   // VY (y-coordinate from register)
	numRows := uint8(cpu.opcode & 0x000F) // Height (N)

	// Reset collision flag
	cpu.v[0xF] = 0

	for row := uint8(0); row < numRows; row++ {
		// Ensure memory access is within bounds
		if int(cpu.i)+int(row) >= len(cpu.memory) {
			fmt.Printf("Memory access out of bounds at I=0x%X, row %d\n", cpu.i, row)
			break
		}

		spriteRow := cpu.memory[cpu.i+uint16(row)] // Fetch sprite row

		for col := uint8(0); col < 8; col++ {
			// Check if the bit at (col) is set
			if (spriteRow & (0x80 >> col)) != 0 {

				// Wrap coordinates
				px := (vx + col) % 64
				py := (vy + row) % 32

				// Check for collision
				if cpu.gfx[px][py] == 1 {
					cpu.v[0xF] = 1 // Set collision flag
				}

				// XOR the pixel
				cpu.gfx[px][py] ^= 1
			}
		}
	}

	// Set draw flag to refresh screen
	cpu.drawFlag = true

	// Increment program counter
	cpu.pc += 2
}
