package core

// executeOpcode executes a CHIP-8 opcode
func (c *CPU) executeOpcode(opcode uint16) {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := opcode & 0x000F
	nn := byte(opcode & 0x00FF)
	nnn := opcode & 0x0FFF

	switch opcode & 0xF000 {
	case 0x0000:
		switch opcode {
		case 0x00E0:
			// 00E0: Clear screen
			for i := range c.Display {
				for j := range c.Display[i] {
					c.Display[i][j] = 0
				}
			}
			c.DrawFlag = true
		case 0x00EE:
			// 00EE: Return from subroutine
			if c.SP > 0 {
				c.SP--
				c.PC = c.Stack[c.SP]
			}
		}
	case 0x1000:
		// 1NNN: Jump to address NNN
		c.PC = nnn
	case 0x2000:
		// 2NNN: Call subroutine at NNN
		if c.SP < 16 {
			c.Stack[c.SP] = c.PC
			c.SP++
			c.PC = nnn
		}
	case 0x3000:
		// 3XNN: Skip if VX == NN
		if c.V[x] == nn {
			c.PC += 2
		}
	case 0x4000:
		// 4XNN: Skip if VX != NN
		if c.V[x] != nn {
			c.PC += 2
		}
	case 0x5000:
		// 5XY0: Skip if VX == VY
		if c.V[x] == c.V[y] {
			c.PC += 2
		}
	case 0x6000:
		// 6XNN: Set VX = NN
		c.V[x] = nn
	case 0x7000:
		// 7XNN: VX += NN
		c.V[x] += nn
	case 0x8000:
		switch n {
		case 0x0:
			// 8XY0: VX = VY
			c.V[x] = c.V[y]
		case 0x1:
			// 8XY1: VX |= VY
			c.V[x] |= c.V[y]
		case 0x2:
			// 8XY2: VX &= VY
			c.V[x] &= c.V[y]
		case 0x3:
			// 8XY3: VX ^= VY
			c.V[x] ^= c.V[y]
		case 0x4:
			// 8XY4: VX += VY; VF = carry
			if c.V[x]+c.V[y] > 255 {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[x] += c.V[y]
		case 0x5:
			// 8XY5: VX -= VY; VF = !borrow
			if c.V[x] > c.V[y] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[x] -= c.V[y]
		case 0x6:
			// 8XY6: VX >>= 1; VF = LSB
			c.V[0xF] = c.V[x] & 0x1
			c.V[x] >>= 1
		case 0x7:
			// 8XY7: VX = VY - VX; VF = !borrow
			if c.V[y] > c.V[x] {
				c.V[0xF] = 1
			} else {
				c.V[0xF] = 0
			}
			c.V[x] = c.V[y] - c.V[x]
		case 0xE:
			// 8XYE: VX <<= 1; VF = MSB
			c.V[0xF] = (c.V[x] >> 7) & 0x1
			c.V[x] <<= 1
		}
	case 0x9000:
		// 9XY0: Skip if VX != VY
		if c.V[x] != c.V[y] {
			c.PC += 2
		}
	case 0xA000:
		// ANNN: Set I = NNN
		c.I = nnn
	case 0xB000:
		// BNNN: Jump to V0 + NNN
		c.PC = uint16(c.V[0]) + nnn
	case 0xC000:
		// CXNN: VX = random() & NN
		c.V[x] = byte(c.rng.Intn(256)) & nn
	case 0xD000:
		// DXYN: Draw sprite at (VX, VY)
		c.V[0xF] = 0
		for row := byte(0); row < byte(n); row++ {
			spriteAddr := c.I + uint16(row)
			if spriteAddr >= 4096 {
				continue
			}
			sprite := c.Memory[spriteAddr]
			for col := byte(0); col < 8; col++ {
				if sprite&(0x80>>col) != 0 {
					xPos := (c.V[x] + col) % 64
					yPos := (c.V[y] + row) % 32
					if xPos < 64 && yPos < 32 {
						if c.Display[yPos][xPos] == 1 {
							c.V[0xF] = 1
						}
						c.Display[yPos][xPos] ^= 1
					}
				}
			}
		}
		c.DrawFlag = true
	case 0xE000:
		switch nn {
		case 0x9E:
			// EX9E: Skip if key VX is pressed
			if c.V[x] < 16 && c.Key[c.V[x]] {
				c.PC += 2
			}
		case 0xA1:
			// EXA1: Skip if key VX is not pressed
			if c.V[x] < 16 && !c.Key[c.V[x]] {
				c.PC += 2
			}
		}
	case 0xF000:
		switch nn {
		case 0x07:
			// FX07: VX = delay timer
			c.V[x] = c.DelayTimer
		case 0x0A:
			// FX0A: Wait for key press
			keyPressed := false
			for i, key := range c.Key {
				if key {
					c.V[x] = byte(i)
					keyPressed = true
					break
				}
			}
			if !keyPressed {
				c.PC -= 2 // Stay on this instruction
			}
		case 0x15:
			// FX15: Set delay timer = VX
			c.DelayTimer = c.V[x]
		case 0x18:
			// FX18: Set sound timer = VX
			c.SoundTimer = c.V[x]
		case 0x1E:
			// FX1E: I += VX
			c.I += uint16(c.V[x])
		case 0x29:
			// FX29: Set I = sprite address for VX
			if c.V[x] <= 0xF {
				c.I = 0x50 + uint16(c.V[x])*5
			}
		case 0x33:
			// FX33: Store BCD of VX at I
			if c.I+2 < 4096 {
				c.Memory[c.I] = c.V[x] / 100
				c.Memory[c.I+1] = (c.V[x] / 10) % 10
				c.Memory[c.I+2] = c.V[x] % 10
			}
		case 0x55:
			// FX55: Store V0-VX in memory at I
			for i := 0; i <= int(x); i++ {
				if c.I+uint16(i) < 4096 {
					c.Memory[c.I+uint16(i)] = c.V[i]
				}
			}
		case 0x65:
			// FX65: Load V0-VX from memory at I
			for i := 0; i <= int(x); i++ {
				if c.I+uint16(i) < 4096 {
					c.V[i] = c.Memory[c.I+uint16(i)]
				}
			}
		}
	}
}
