package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"fmt"
	"os"
)

type opcodeFunc func()

type chip8 struct {
	registers  [16]uint8
	memory     [4096]uint8
	index      uint16
	pc         uint16
	stack      [16]uint16
	sp         uint8
	delayTimer uint8
	soundTimer uint8
	keypad     [16]uint8
	video      [64 * 32]uint32
	opcode     uint16

	table  [16]opcodeFunc
	table0 [15]opcodeFunc
	table8 [15]opcodeFunc
	tableE [15]opcodeFunc
	tableF [256]opcodeFunc
}

const VIDEO_HEIGHT = 32
const VIDEO_WIDTH = 64

const START_ADDRESS = 0x200
const FONTSET_START_ADDRESS = 0x50

const FontsetSize = 80

var fontset = [FontsetSize]uint8{
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

func newChip8() *chip8 {
	c := &chip8{
		pc: START_ADDRESS,
	}

	copy(c.memory[FONTSET_START_ADDRESS:], fontset[:])

	c.table = [16]opcodeFunc{
		c.Table0, c.OP_1nnn, c.OP_2nnn, c.OP_3xkk,
		c.OP_4xkk, c.OP_5xy0, c.OP_6xkk, c.OP_7xkk,
		c.Table8, c.OP_9xy0, c.OP_Annn, c.OP_Bnnn,
		c.OP_Cxkk, c.OP_Dxyn, c.TableE, c.TableF,
	}

	for i := range c.table0 {
		c.table0[i] = c.OP_NULL
		c.table8[i] = c.OP_NULL
		c.tableE[i] = c.OP_NULL
	}
	for i := range c.tableF {
		c.tableF[i] = c.OP_NULL
	}

	c.table0[0x0] = c.OP_00E0
	c.table0[0xE] = c.OP_00EE

	c.table8[0x0] = c.OP_8xy0
	c.table8[0x1] = c.OP_8xy1
	c.table8[0x2] = c.OP_8xy2
	c.table8[0x3] = c.OP_8xy3
	c.table8[0x4] = c.OP_8xy4
	c.table8[0x5] = c.OP_8xy5
	c.table8[0x6] = c.OP_8xy6
	c.table8[0x7] = c.OP_8xy7
	c.table8[0xE] = c.OP_8xyE

	c.tableE[0x1] = c.OP_ExA1
	c.tableE[0xE] = c.OP_Ex9E

	c.tableF[0x07] = c.OP_Fx07
	c.tableF[0x0A] = c.OP_Fx0A
	c.tableF[0x15] = c.OP_Fx15
	c.tableF[0x18] = c.OP_Fx18
	c.tableF[0x1E] = c.OP_Fx1E
	c.tableF[0x29] = c.OP_Fx29
	c.tableF[0x33] = c.OP_Fx33
	c.tableF[0x55] = c.OP_Fx55
	c.tableF[0x65] = c.OP_Fx65

	return c
}

func (c *chip8) Cycle() {
	// Fetch
	c.opcode = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])

	// Increment the PC before executing the instruction.
	c.pc += 2

	// Decode and execute.
	c.table[(c.opcode&0xF000)>>12]()

	// Decrement timers if they are set.
	if c.delayTimer > 0 {
		c.delayTimer--
	}

	if c.soundTimer > 0 {
		c.soundTimer--
	}
}

func (c *chip8) LoadROM(filename string) error {
	rom, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read ROM: %w", err)
	}

	start := int(START_ADDRESS)
	if len(rom) > len(c.memory)-start {
		return fmt.Errorf("ROM is too large to fit in memory")
	}

	copy(c.memory[start:], rom)
	return nil
}
