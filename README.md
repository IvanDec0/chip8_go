# Chip-8 Emulator

A modern, configurable Chip-8 emulator written in Go with SDL2, featuring robust error handling, memory safety, and customizable display/performance settings.

_[También disponible en Español](README_ES.md)_

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration Options](#configuration-options)
- [Controls](#controls)
- [Project Structure](#project-structure)
- [Technical Details](#technical-details)
- [ROM Compatibility](#rom-compatibility)
- [Troubleshooting](#troubleshooting)

## Overview

The Chip-8 is an interpreted programming language developed in the 1970s for use on 8-bit microcomputers. This emulator faithfully recreates the Chip-8 system, allowing you to run classic games and programs originally designed for the platform.

### What is Chip-8?

- **Platform**: Originally designed for the COSMAC VIP computer
- **Display**: 64×32 monochrome pixel display
- **Memory**: 4KB RAM
- **Registers**: 16 8-bit general-purpose registers (V0-VF)
- **Input**: 16-key hexadecimal keypad
- **Sound**: Single-tone beeper (440Hz)

## Features

### Core Emulation

- ✅ Complete Chip-8 instruction set implementation
- ✅ Accurate timing and display rendering
- ✅ Sound timer support (beep indication)
- ✅ 16-key input handling
- ✅ Memory-safe operation with bounds checking

### Modern Enhancements

- ✅ **Configurable resolution scaling** (1x to 20x)
- ✅ **Adjustable CPU speed** (100-2000 cycles/second)
- ✅ **Custom window titles**
- ✅ **Sound support** with SDL2 audio (440Hz beep)
- ✅ **Clean architecture** with separated concerns
- ✅ **Interface-based design** for extensibility
- ✅ **Robust error handling** and input validation
- ✅ **Memory protection** against buffer overflows
- ✅ **Stack overflow/underflow protection**
- ✅ **Cross-platform support** (Linux, Windows, macOS)

### Quality Assurance

- ✅ Comprehensive bounds checking for all memory operations
- ✅ Safe sprite rendering with coordinate validation
- ✅ Protected subroutine call/return mechanisms
- ✅ Input parameter validation with clear error messages

## Prerequisites

- **Go**: Version 1.19 or later
- **SDL2**: Development libraries
  - **Ubuntu/Debian**: `sudo apt-get install libsdl2-dev`
  - **Fedora/CentOS**: `sudo dnf install SDL2-devel`
  - **macOS**: `brew install sdl2`
  - **Windows**: Download from [SDL2 releases](https://github.com/libsdl-org/SDL/releases)

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/IvanDec0/chip8_go
   cd chip8_go
   ```

2. **Install dependencies:**

   ```bash
   go mod download
   ```

3. **Build the emulator:**
   ```bash
   go build -o chip8-emulator
   ```

## Usage

### Basic Usage

```bash
# Run with default settings
./chip8-emulator <ROM-file>

# Example
./chip8-emulator roms/games/Pong.ch8
```

### Advanced Usage

```bash
# Custom configuration
./chip8-emulator [options] <ROM-file>

# Examples
./chip8-emulator -scale 15 -speed 500 -title "Pong Game" roms/games/Pong.ch8
./chip8-emulator -scale 5 -speed 200 roms/demos/Maze\ \[David\ Winter,\ 199x\].ch8
```

### Help

```bash
./chip8-emulator --help
```

## Configuration Options

| Flag     | Description           | Range      | Default |
| -------- | --------------------- | ---------- | ------- |
| `-scale` | Window scale factor   | 1-20       | 10      |
| `-speed` | CPU cycles per second | 100-2000   | 700     |
| `-title` | Custom window title   | Any string | "Chip8" |

### Scale Examples

- **Scale 1**: 64×32 pixels (original size)
- **Scale 10**: 640×320 pixels (default)
- **Scale 20**: 1280×640 pixels (maximum)

### Speed Examples

- **100 cycles/sec**: Very slow, good for debugging
- **700 cycles/sec**: Default speed, works well for most games
- **2000 cycles/sec**: Fast execution for demos

## Controls

The emulator maps your keyboard to the Chip-8's 16-key hexadecimal keypad:

```
Chip-8 Keypad    Keyboard Mapping
┌─┬─┬─┬─┐       ┌─┬─┬─┬─┐
│1│2│3│C│       │1│2│3│4│
├─┼─┼─┼─┤       ├─┼─┼─┼─┤
│4│5│6│D│       │Q│W│E│R│
├─┼─┼─┼─┤   =>  ├─┼─┼─┼─┤
│7│8│9│E│       │A│S│D│F│
├─┼─┼─┼─┤       ├─┼─┼─┼─┤
│A│0│B│F│       │Z│X│C│V│
└─┴─┴─┴─┘       └─┴─┴─┴─┘
```

### Special Keys

- **ESC**: Exit emulator
- **Window Close (X)**: Exit emulator

### Menu Controls

When browsing ROMs in the menu:

- **↑/↓ or W/S**: Navigate up/down through ROM list
- **Enter or Space**: Select ROM to load
- **F**: Toggle ROM as favorite ⭐
- **ESC**: Go back or exit to system

The favorites system allows you to mark your favorite ROMs with a star (⭐) for quick access.

## Advanced Features

### Sound System

The emulator includes a full audio implementation:

- **Frequency**: 440Hz (A4 note)
- **Automatic playback**: Sound plays when CHIP-8 programs use the sound timer

### Configuration Examples

Perfect for different use cases:

- **Development**: `./chip8-emulator -scale 5 -speed 200 rom.ch8` (small, slow)
- **Gaming**: `./chip8-emulator -scale 15 -speed 700 rom.ch8` (large, normal speed)
- **Demos**: `./chip8-emulator -scale 20 -speed 1000 rom.ch8` (full screen, fast)

## Project Structure

```
chip8_go/
├── main.go                    # Simple entry point and configuration
├── internal/                  # Private implementation packages
│   ├── core/                  # Pure CHIP-8 CPU logic (no dependencies)
│   │   ├── cpu.go             # CPU state and core functionality
│   │   └── opcodes.go         # Instruction set implementation
│   ├── audio/                 # Audio system interface & implementation
│   │   ├── audio.go           # Audio interface definition
│   │   └── sdl_audio.go       # SDL2 audio implementation
│   ├── renderer/              # Rendering interface & implementation
│   │   ├── renderer.go        # Renderer interface definition
│   │   └── sdl_renderer.go    # SDL2 renderer implementation
│   └── input/                 # Input handling
│       └── handler.go         # Keyboard input processing
├── pkg/                       # Public packages
│   └── emulator/              # High-level emulator coordination
│       └── emulator.go        # Main emulator orchestration
├── roms/                      # ROM collection
│   ├── demos/                 # Demo programs
│   ├── games/                 # Classic games
│   └── programs/              # Utility programs
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency checksums
└──  README.md                 # This documentation
```

### Key Components

- **`main.go`**: Clean entry point with configuration parsing
- **`internal/core/`**: Pure CHIP-8 CPU logic, platform-independent
- **`internal/audio/`**: Interface-based audio system with SDL2 implementation
- **`internal/renderer/`**: Interface-based rendering system with SDL2 implementation
- **`internal/input/`**: Event-driven input handling
- **`pkg/emulator/`**: Public API that coordinates all components
- **`roms/`**: Collection of CHIP-8 programs and games

## Technical Details

### Memory Layout

```
0x000-0x1FF: Chip-8 interpreter (font data stored at 0x50-0x9F)
0x200-0xFFF: Program ROM and RAM (3584 bytes)
```

### Registers

- **V0-VF**: 16 general-purpose 8-bit registers
- **I**: 16-bit address register
- **PC**: Program counter
- **SP**: Stack pointer
- **DT**: Delay timer (decrements at 60Hz)
- **ST**: Sound timer (decrements at 60Hz, beeps when > 0)

### Display

- **Resolution**: 64×32 pixels
- **Colors**: Monochrome (black/white)
- **Rendering**: XOR-based sprite drawing
- **Refresh**: 60 FPS

### Instruction Set

The emulator implements all 35 standard Chip-8 opcodes:

- Memory operations (load, store, copy)
- Arithmetic and logic operations
- Control flow (jump, call, return, skip)
- Graphics (clear screen, draw sprite)
- Input handling (key press detection)
- Timer operations

### Safety Features

1. **Memory Protection**:

   - Bounds checking for all memory access
   - Prevention of buffer overflows
   - Safe ROM loading with size validation

2. **Stack Protection**:

   - Stack overflow prevention (max 16 levels)
   - Stack underflow protection
   - Safe subroutine handling

3. **Display Safety**:
   - Coordinate wrapping and bounds checking
   - Protected sprite rendering
   - Safe pixel manipulation

## ROM Compatibility

The emulator is compatible with standard Chip-8 ROMs and has been tested with:

### Included Games

- **Pong**: Classic paddle game
- **Breakout**: Brick-breaking game
- **Tetris**: Block puzzle game
- **Space Invaders**: Arcade shooter
- **Pac-Man**: Maze navigation game

### Included Demos

- **Maze**: Procedural maze generation
- **Particle Demo**: Visual effects demonstration
- **Sierpinski**: Fractal pattern generation
- **Stars**: Animated starfield

### File Formats

- **Extension**: `.ch8` (standard)
- **Size**: Up to 3584 bytes
- **Format**: Raw binary data

## Troubleshooting

### Common Issues

1. **"Error initializing SDL"**

   - Ensure SDL2 development libraries are installed
   - Check your graphics drivers are up to date

2. **"Error loading ROM"**

   - Verify the ROM file exists and is readable
   - Check file permissions
   - Ensure ROM is a valid Chip-8 file

3. **Window too small/large**

   - Adjust the `-scale` parameter (1-20)
   - Try `-scale 10` for a good default size

4. **Game running too fast/slow**

   - Adjust the `-speed` parameter (100-2000)
   - Most games work well with 500-1000 cycles/second

### Performance Tips

- Use lower scale values for better performance on older hardware
- Reduce CPU speed for complex ROMs that run too fast
- Close other applications if experiencing stuttering

---

## Quick Start Examples

```bash
# Download and run a classic game
./chip8-emulator roms/games/Pong.ch8

# Run a demo with custom settings
./chip8-emulator -scale 15 -speed 800 -title "Chip8 Demo" roms/demos/Maze\ \[David\ Winter,\ 199x\].ch8

# Debug mode (slow speed, small window)
./chip8-emulator -scale 5 -speed 200 your-rom.ch8
```

Enjoy exploring the world of Chip-8 programming and gaming! 🎮
