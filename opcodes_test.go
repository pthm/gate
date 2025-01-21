package main

import (
	"testing"
)

func Test_op2NNN(t *testing.T) {
	cpu := NewCPU()
	cpu.LoadROM([]uint8{0x24, 0x00})

	currentPC := cpu.pc
	currentSP := cpu.sp

	cpu.cycle()

	if cpu.stack[currentSP] != currentPC {
		t.Fatalf("stack at sp should contain 0x200, was 0x%X\n", cpu.stack[currentSP])
	}
	if cpu.sp != currentSP+1 {
		t.Fatalf("sp should increment by 1\n")
	}
	if cpu.pc != 0x400 {
		t.Fatalf("pc should equal 0x400, was 0x%X", cpu.opcode)
	}
}

func Test_op8XY4(t *testing.T) {
	cpu := NewCPU()
	cpu.LoadROM([]uint8{
		0x80, 0x14,
		0x82, 0x34,
	})

	cpu.v[0x0] = 255
	cpu.v[0x1] = 1
	cpu.v[0x2] = 4
	cpu.v[0x3] = 4

	currentPC := cpu.pc

	cpu.cycle()
	if cpu.pc != currentPC+2 {
		t.Fatalf("program counter did not increase by two\n")
	}
	// carry should be 1
	if cpu.v[0xF] != 1 {
		t.Fatalf("carry should be set to 1")
	}
	// V0 should be 0 because we are adding VY to VX (0x8014)
	// X = 0, Y = 1, VX = 255, VY = 1
	if cpu.v[0x0] != 0x0 {
		t.Fatalf("V0 should be 0, was 0x%X\n", cpu.v[0x0])
	}

	currentPC = cpu.pc
	cpu.cycle()
	if cpu.pc != currentPC+2 {
		t.Fatalf("program counter did not increase by two\n")
	}
	// carry should be 0
	if cpu.v[0xF] != 0 {
		t.Fatalf("carry should be set to 0")
	}
	// V2 should be 8 because we are adding VY to VX (0x8234)
	// X = 2, Y = 3, VX = 4, VY = 4
	if cpu.v[0x2] != 0x8 {
		t.Fatalf("V2 should be 8, was 0x%X\n", cpu.v[0x0])
	}
}

func Test_opANNN(t *testing.T) {
	cpu := NewCPU()
	cpu.LoadROM([]uint8{0xA4, 0x00})

	currentPC := cpu.pc

	cpu.cycle()

	if cpu.pc != currentPC+2 {
		t.Fatalf("program counter did not increase by two\n")
	}
	if cpu.i != 0x400 {
		t.Fatalf("i was not set to 0x400, was 0x%X\n", cpu.i)
	}
}

func Test_opFX33(t *testing.T) {
	cpu := NewCPU()
	cpu.LoadROM([]uint8{0xF0, 0x33})

	cpu.v[0x0] = 0xFF

	currentPC := cpu.pc

	cpu.cycle()

	if cpu.pc != currentPC+2 {
		t.Fatalf("program counter did not increase by two\n")
	}
	if cpu.memory[0] != 2 {
		t.Fatalf("memory at 0 should be 2")
	}
	if cpu.memory[1] != 5 {
		t.Fatalf("memory at 1 should be 5")
	}
	if cpu.memory[2] != 5 {
		t.Fatalf("memory at 2 should be 5")
	}
}
