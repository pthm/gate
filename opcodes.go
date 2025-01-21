package main

// 2NNN - Calls subroutine at NNN
func op2NNN(cpu *CPU) {
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
	cpu.memory[cpu.i+2] = (vx % 100) % 10
	cpu.pc += 2
}

// 00E0 - Clear screen
func op00E0(cpu *CPU) {
	cpu.gfx = [64 * 32]uint8{}
	cpu.pc += 2
}

// 1NNN - Jump
func op1NNN() {

}

// 6XNN - Set register VX
func op6XNN() {

}

// 7XNN - Add value to register VX
func op7XNN() {

}

// DXYN - Draws a sprite at coordinate (VX, VY) that has a width of 8 pixels and a height of N pixels.
// Each row of 8 pixels is read as bit-coded starting from memory location I; I value does not change after the execution of this instruction.
// As described above, VF is set to 1 if any screen pixels are flipped from set to unset when the sprite is drawn, and to 0 if that does not happen
func opDXYN(cpu *CPU) {
	x := uint8((cpu.opcode & 0x0F00) >> 8) // Fetch X from the opcode, shift it 8 bits so its in the most significant bit
	y := uint8((cpu.opcode & 0x00F0) >> 4) // Fetch Y from the opcode, shift it 4 bits so its in the most significant bit

	numRows := uint8(cpu.opcode & 0x000F) // Fetch the number of sprite rows or the height from the opcode, its already in the most significant bit

	// Set the carry to 0
	cpu.v[0xF] = 0
	for row := uint8(0); row < numRows; row++ {
		pixel := cpu.memory[cpu.i+uint16(row)]
		// Each "pixel" of a sprite is 8 bits wide
		for col := uint8(0); col < 8; col++ {
			// Use a mask to select the "sub-pixel" from the sprite at bit position col
			// We adjust the mask with the col value, its shifting 0x80 right
			if (pixel & (0x80 >> col)) != 0 { // If the sub-pixel is a 1 we need to draw it
				pos := (x + row) + ((y + col) * 64)
				if cpu.gfx[pos] == 1 { // If the sub pixel in the gfx array is already a 1, we are going to flip it and need to set the carry
					cpu.v[0xF] = 1
				}
				cpu.gfx[pos] ^= 1 // XOR the sub pixel at pos
			}
		}
	}

	cpu.drawFlag = true
	cpu.pc += 2
}
