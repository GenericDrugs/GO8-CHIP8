package main

import "math/rand"

func (c *chip8) OP_NULL() {
}

func (c *chip8) OP_00E0() {
	c.video = [64 * 32]uint32{}
}

func (c *chip8) OP_00EE() {
	c.sp--
	c.pc = c.stack[c.sp]
}

func (c *chip8) OP_1nnn() {
	address := c.opcode & 0x0FFF
	c.pc = address
}

func (c *chip8) OP_2nnn() {
	address := c.opcode & 0x0FFF
	c.stack[c.sp] = c.pc
	c.sp++
	c.pc = address
}

func (c *chip8) OP_3xkk() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	byte := uint8(c.opcode) & 0x00FF

	if c.registers[vx] == byte {
		c.pc += 2
	}
}

func (c *chip8) OP_4xkk() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	byte := uint8(c.opcode) & 0x00FF

	if c.registers[vx] != byte {
		c.pc += 2
	}
}

func (c *chip8) OP_5xy0() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)

	if c.registers[vx] == c.registers[vy] {
		c.pc += 2
	}
}

func (c *chip8) OP_6xkk() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	byte := uint8(c.opcode) & 0x00FF

	c.registers[vx] = byte
}

func (c *chip8) OP_7xkk() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	byte := uint8(c.opcode) & 0x00FF

	c.registers[vx] += byte
}

func (c *chip8) OP_8xy0() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)

	c.registers[vx] = c.registers[vy]
}

func (c *chip8) OP_8xy1() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)

	c.registers[vx] |= c.registers[vy]
}

func (c *chip8) OP_8xy2() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)

	c.registers[vx] &= c.registers[vy]
}

func (c *chip8) OP_8xy3() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)

	c.registers[vx] ^= c.registers[vy]
}

func (c *chip8) OP_8xy4() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)

	sum := uint16(c.registers[vx]) + uint16(c.registers[vy])

	if sum > 255 {
		c.registers[0xF] = 1
	} else {
		c.registers[0xF] = 0
	}

	c.registers[vx] = uint8(sum)
}

func (c *chip8) OP_8xy5() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)

	if c.registers[vx] > c.registers[vy] {
		c.registers[0xF] = 1
	} else {
		c.registers[0xF] = 0
	}

	c.registers[vx] -= c.registers[vy]
}

func (c *chip8) OP_8xy6() {
	vx := uint8((c.opcode & 0x0F00) >> 8)

	c.registers[0xF] = c.registers[vx] & 0x1

	c.registers[vx] >>= 1
}

func (c *chip8) OP_8xy7() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)

	if c.registers[vx] > c.registers[vy] {
		c.registers[0xF] = 1
	} else {
		c.registers[0xF] = 0
	}

	c.registers[vx] = c.registers[vy] - c.registers[vx]
}

func (c *chip8) OP_8xyE() {
	vx := uint8((c.opcode & 0x0F00) >> 8)

	c.registers[0xF] = (c.registers[vx] & 0x80) >> 7

	c.registers[vx] <<= 1
}

func (c *chip8) OP_9xy0() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)

	if c.registers[vx] != c.registers[vy] {
		c.pc += 2
	}
}

func (c *chip8) OP_Annn() {
	address := c.opcode & 0x0FFF
	c.index = address
}

func (c *chip8) OP_Bnnn() {
	address := c.opcode & 0x0FFF
	c.pc = uint16(c.registers[0]) + address
}

func (c *chip8) OP_Cxkk() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	byte := uint8(c.opcode) & 0x00FF

	c.registers[vx] = uint8(rand.Intn(256)) & byte
}

func (c *chip8) OP_Dxyn() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	vy := uint8((c.opcode & 0x00F0) >> 4)
	height := uint8(c.opcode & 0x000F)

	xPos := c.registers[vx] % VIDEO_WIDTH
	yPos := c.registers[vy] % VIDEO_HEIGHT

	c.registers[0xF] = 0

	for row := 0; row < int(height); row++ {
		spriteByte := c.memory[int(c.index)+row]

		for col := 0; col < 8; col++ {
			spritePixel := spriteByte & (0x80 >> col)
			if spritePixel == 0 {
				continue
			}

			x := (int(xPos) + col) % VIDEO_WIDTH
			y := (int(yPos) + row) % VIDEO_HEIGHT
			screenPixel := &c.video[y*VIDEO_WIDTH+x]

			if *screenPixel != 0 {
				c.registers[0xF] = 1
			}
			*screenPixel ^= 0xFFFFFFFF
		}
	}

}

func (c *chip8) OP_Ex9E() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	key := c.registers[vx]

	if c.keypad[key] == 1 {
		c.pc += 2
	}
}

func (c *chip8) OP_ExA1() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	key := c.registers[vx]

	if c.keypad[key] != 1 {
		c.pc += 2
	}
}

func (c *chip8) OP_Fx07() {
	vx := uint8((c.opcode & 0x0F00) >> 8)

	c.registers[vx] = c.delayTimer
}

func (c *chip8) OP_Fx0A() {
	vx := uint8((c.opcode & 0x0F00) >> 8)

	if c.keypad[0] != 0 {
		c.registers[vx] = 0
	} else if c.keypad[1] != 0 {
		c.registers[vx] = 1
	} else if c.keypad[2] != 0 {
		c.registers[vx] = 2
	} else if c.keypad[3] != 0 {
		c.registers[vx] = 3
	} else if c.keypad[4] != 0 {
		c.registers[vx] = 4
	} else if c.keypad[5] != 0 {
		c.registers[vx] = 5
	} else if c.keypad[6] != 0 {
		c.registers[vx] = 6
	} else if c.keypad[7] != 0 {
		c.registers[vx] = 7
	} else if c.keypad[8] != 0 {
		c.registers[vx] = 8
	} else if c.keypad[9] != 0 {
		c.registers[vx] = 9
	} else if c.keypad[10] != 0 {
		c.registers[vx] = 10
	} else if c.keypad[11] != 0 {
		c.registers[vx] = 11
	} else if c.keypad[12] != 0 {
		c.registers[vx] = 12
	} else if c.keypad[13] != 0 {
		c.registers[vx] = 13
	} else if c.keypad[14] != 0 {
		c.registers[vx] = 14
	} else if c.keypad[15] != 0 {
		c.registers[vx] = 15
	} else {
		c.pc -= 2
	}
}

func (c *chip8) OP_Fx15() {
	vx := uint8((c.opcode & 0x0F00) >> 8)

	c.delayTimer = c.registers[vx]
}

func (c *chip8) OP_Fx18() {
	vx := uint8((c.opcode & 0x0F00) >> 8)

	c.soundTimer = c.registers[vx]
}

func (c *chip8) OP_Fx1E() {
	vx := uint8((c.opcode & 0x0F00) >> 8)

	c.index += uint16(c.registers[vx])
}

func (c *chip8) OP_Fx29() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	digit := c.registers[vx]
	c.index = FONTSET_START_ADDRESS + 5*uint16(digit)
}

func (c *chip8) OP_Fx33() {
	vx := uint8((c.opcode & 0x0F00) >> 8)
	value := c.registers[vx]

	c.memory[c.index+2] = value % 10
	value /= 10

	c.memory[c.index+1] = value % 10
	value /= 10

	c.memory[c.index] = value % 10
}

func (c *chip8) OP_Fx55() {
	vx := uint8((c.opcode & 0x0F00) >> 8)

	for i := uint8(0); i <= vx; i++ {
		c.memory[c.index+uint16(i)] = c.registers[i]
	}
}

func (c *chip8) OP_Fx65() {
	vx := uint8((c.opcode & 0x0F00) >> 8)

	for i := uint8(0); i <= vx; i++ {
		c.registers[i] = c.memory[c.index+uint16(i)]
	}
}

func (c *chip8) TableE() {
	c.tableE[c.opcode&0x000F]()
}

func (c *chip8) Table0() {
	switch c.opcode & 0x00FF {
	case 0xE0:
		c.OP_00E0()
	case 0xEE:
		c.OP_00EE()
	default:
		c.OP_NULL() // or report an unsupported opcode
	}
}

func (c *chip8) Table8() {
	c.table8[c.opcode&0x000F]()
}

func (c *chip8) TableF() {
	c.tableF[c.opcode&0x00FF]()
}
