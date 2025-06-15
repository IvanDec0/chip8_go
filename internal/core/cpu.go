package core

import (
	"math/rand"
	"os"
	"time"
)

var fontSet = [80]byte{
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

// CPU represents the CHIP-8 CPU state
type CPU struct {
	Memory     [4096]byte
	V          [16]byte
	I          uint16
	PC         uint16
	Stack      [16]uint16
	SP         uint8
	DelayTimer byte
	SoundTimer byte
	Display    [32][64]byte
	DrawFlag   bool
	Key        [16]bool
	rng        *rand.Rand
}

// NewCPU creates a new CHIP-8 CPU instance
func NewCPU() *CPU {
	cpu := &CPU{
		PC:  0x200,
		rng: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Load font set into memory
	for i, b := range fontSet {
		cpu.Memory[0x050+i] = b
	}

	return cpu
}

// LoadROM loads a ROM file into memory
func (c *CPU) LoadROM(path string) error {
	rom, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if len(rom) > 4096-0x200 {
		rom = rom[:4096-0x200]
	}

	for i, b := range rom {
		c.Memory[0x200+i] = b
	}

	return nil
}

// UpdateTimers decrements the delay and sound timers
func (c *CPU) UpdateTimers() {
	if c.DelayTimer > 0 {
		c.DelayTimer--
	}
	if c.SoundTimer > 0 {
		c.SoundTimer--
	}
}

// IsSoundActive returns true if sound should be playing
func (c *CPU) IsSoundActive() bool {
	return c.SoundTimer > 0
}

// GetDisplay returns a copy of the display buffer
func (c *CPU) GetDisplay() [32][64]byte {
	return c.Display
}

// ShouldDraw returns true if the display needs to be redrawn
func (c *CPU) ShouldDraw() bool {
	return c.DrawFlag
}

// ClearDrawFlag clears the draw flag
func (c *CPU) ClearDrawFlag() {
	c.DrawFlag = false
}

// SetKey sets the state of a key
func (c *CPU) SetKey(key uint8, pressed bool) {
	if key < 16 {
		c.Key[key] = pressed
	}
}

// Cycle executes one CPU cycle
func (c *CPU) Cycle() {
	if c.PC >= 4096-2 {
		c.PC = 0x200
		return
	}

	opcode := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])
	c.PC += 2
	c.executeOpcode(opcode)
}
