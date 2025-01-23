package cpu

import (
	"fmt"
)

// Op00E0 - Clear screen
func Op00E0(cpu *CPU) {
	cpu.gfx = [64][32]uint8{}
	cpu.pc += 2
}

// Op00EE - Returns from a subroutine
func Op00EE(cpu *CPU) {
	if cpu.sp == 0 {
		fmt.Println("Stack underflow!")
		return // Prevent underflow
	}
	cpu.sp--                   // Decrement the stack pointer so we are at the "top" of the stack
	cpu.pc = cpu.stack[cpu.sp] // Set the PC to the value in at the "top" of the stack
}

// Op1NNN - Jumps to address NNN
func Op1NNN(cpu *CPU) {
	cpu.pc = cpu.opcode & 0x0FFF // Set the program counter to the address NNN, we use the mask 0x0FFF to extract NNN
}

// Op2NNN - Calls subroutine at NNN
func Op2NNN(cpu *CPU) {
	if cpu.sp >= uint16(len(cpu.stack)) {
		fmt.Println("Stack overflow")
		return // Prevent overflow
	}
	cpu.stack[cpu.sp] = cpu.pc   // Store the program counter value in the stack at the current stack pointer
	cpu.sp++                     // Increment the stack pointer
	cpu.pc = cpu.opcode & 0x0FFF // Set the program counter to the address NNN, we use the mask 0x0FFF to extract NNN
}

// Op3XNN - Skips the next instruction if VX equals NN (usually the next instruction is a jump to skip a code block).
func Op3XNN(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8  // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	nn := uint8(cpu.opcode & 0x00FF) // Fetch NN from the opcode
	vx := cpu.v[x]
	if vx == nn {
		cpu.pc += 2
	}
	cpu.pc += 2
}

// Op4XNN - Skips the next instruction if VX does not equal NN (usually the next instruction is a jump to skip a code block).
func Op4XNN(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8  // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	nn := uint8(cpu.opcode & 0x00FF) // Fetch NN from the opcode
	vx := cpu.v[x]
	if vx != nn {
		cpu.pc += 2
	}
	cpu.pc += 2
}

// Op5XY0 - Skips the next instruction if VX equals VY (usually the next instruction is a jump to skip a code block)
func Op5XY0(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := (cpu.opcode & 0x00F0) >> 4 // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	vx := cpu.v[x]
	vy := cpu.v[y]

	if vx == vy {
		cpu.pc += 2
	}

	cpu.pc += 2
}

// Op6XNN - Sets VX to NN.
func Op6XNN(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	nn := cpu.opcode & 0x00FF       // Fetch NN from the opcode

	cpu.v[x] = uint8(nn)
	cpu.pc += 2
}

// Op7XNN - Adds NN to VX (carry flag is not changed).
func Op7XNN(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	nn := cpu.opcode & 0x00FF       // Fetch NN from the opcode

	cpu.v[x] += uint8(nn)
	cpu.pc += 2
}

// Op8XY0 - Sets VX to the value of VY
func Op8XY0(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := (cpu.opcode & 0x00F0) >> 4 // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	cpu.v[x] = cpu.v[y]
	cpu.pc += 2
}

// Op8XY1 - Sets VX to VX or VY. (bitwise OR operation)
func Op8XY1(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := (cpu.opcode & 0x00F0) >> 4 // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	cpu.v[x] = cpu.v[x] | cpu.v[y]
	cpu.pc += 2
}

// Op8XY2 - Sets VX to VX and VY. (bitwise AND operation)
func Op8XY2(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := (cpu.opcode & 0x00F0) >> 4 // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	cpu.v[x] = cpu.v[x] & cpu.v[y]
	cpu.pc += 2
}

// Op8XY3 - Sets VX to VX xor VY
func Op8XY3(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := (cpu.opcode & 0x00F0) >> 4 // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	cpu.v[x] = cpu.v[x] ^ cpu.v[y]
	cpu.pc += 2
}

// Op8XY4 - Adds VY to VX. VF is set to 1 when there's an overflow, and to 0 when there is not
func Op8XY4(cpu *CPU) {
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

// Op8XY5 - VY is subtracted from VX. VF is set to 0 when there's an underflow, and 1 when there is not. (i.e. VF set to 1 if VX >= VY and 0 if not)
func Op8XY5(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := (cpu.opcode & 0x00F0) >> 4 // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	if cpu.v[x] >= cpu.v[y] {
		cpu.v[0xF] = 1 // Set carry to 1
	} else {
		cpu.v[0xF] = 0
	}
	cpu.v[x] -= cpu.v[y]
	cpu.pc += 2
}

// Op8XY6 - Shifts VX to the right by 1, then stores the least significant bit of VX prior to the shift into VF
func Op8XY6(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit

	// Store the least significant bit in VF
	cpu.v[0xF] = cpu.v[x] & 0x01
	// Shift VX to the right
	cpu.v[x] = cpu.v[x] >> 0x0001

	cpu.pc += 2
}

// Op8XY7 - Sets VX to VY minus VX. VF is set to 0 when there's an underflow, and 1 when there is not. (i.e. VF set to 1 if VY >= VX)
func Op8XY7(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := (cpu.opcode & 0x00F0) >> 4 // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	if cpu.v[y] >= cpu.v[x] {
		cpu.v[0xF] = 1
	} else {
		cpu.v[0xF] = 0
	}

	cpu.v[x] = cpu.v[y] - cpu.v[x]
	cpu.pc += 2
}

// Op8XYE - Shifts VX to the left by 1, then sets VF to 1 if the most significant bit of VX prior to that shift was set, or to 0 if it was unset
func Op8XYE(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit

	cpu.v[0xF] = (cpu.v[x] & 0x80) >> 7 // Set VF to the most significant bit of VX (prior to the shift)
	cpu.v[x] <<= 1                      // Shift VX left by 1

	cpu.pc += 2
}

// Op9XY0 - Skips the next instruction if VX does not equal VY. (Usually the next instruction is a jump to skip a code block)
func Op9XY0(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := (cpu.opcode & 0x00F0) >> 4 // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	if cpu.v[x] != cpu.v[y] {
		cpu.pc += 2
	}

	cpu.pc += 2
}

// OpANNN - Sets I to the address NNN
func OpANNN(cpu *CPU) {
	cpu.i = cpu.opcode & 0x0FFF // Set I to NNN
	cpu.pc += 2                 // Increment the program counter by two
}

// OpDXYN - Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.
// Each row of 8 pixels is read as bit-coded starting from memory location I; I value does not change after the execution of this instruction.
// As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that does not happen
func OpDXYN(cpu *CPU) {
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
	cpu.pc += 2
}

// FX1E - Adds VX to I. VF is not affected
func OpFX1E(cpu *CPU) {
	x := (cpu.opcode & 0x0F00) >> 8 // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	vx := cpu.v[x]

	cpu.i += uint16(vx)

	cpu.pc += 2
}

// OpFX33 - Stores the binary-coded decimal representation of VX, with the hundreds digit in memory at location in I, the tens digit at location I+1, and the ones digit at location I+2.
func OpFX33(cpu *CPU) {
	x := uint8((cpu.opcode & 0x0F00) >> 8)
	vx := cpu.v[x]
	cpu.memory[cpu.i] = vx / 100
	cpu.memory[cpu.i+1] = (vx / 10) % 10
	cpu.memory[cpu.i+2] = vx % 10
	cpu.pc += 2
}
