package main

import (
	"chip8/internal/audio"
	"chip8/internal/renderer"
	"chip8/pkg/emulator"
	"flag"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

type CLIConfig struct {
	WindowTitle     string
	Scale           int
	CyclesPerSecond int
}

func parseCLIFlags() (*CLIConfig, string) {
	config := &CLIConfig{
		WindowTitle:     "Chip8",
		Scale:           10,
		CyclesPerSecond: 700,
	}

	flag.IntVar(&config.Scale, "scale", 10, "Window scale factor (1-20)")
	flag.IntVar(&config.CyclesPerSecond, "speed", 700, "CPU cycles per second (100-2000)")
	flag.StringVar(&config.WindowTitle, "title", "Chip8", "Window title")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <ROM file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s game.ch8\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -scale 15 -speed 500 game.ch8\n", os.Args[0])
	}

	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	// Validate config values
	if config.Scale < 1 || config.Scale > 20 {
		fmt.Fprintf(os.Stderr, "Error: Scale must be between 1 and 20\n")
		os.Exit(1)
	}

	if config.CyclesPerSecond < 100 || config.CyclesPerSecond > 2000 {
		fmt.Fprintf(os.Stderr, "Error: Speed must be between 100 and 2000 cycles per second\n")
		os.Exit(1)
	}

	return config, flag.Arg(0)
}

func main() {
	config, romPath := parseCLIFlags()

	// Initialize SDL
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing SDL: %s\n", err.Error())
		os.Exit(1)
	}
	defer sdl.Quit()

	// Create emulator
	emu, err := emulator.New(emulator.Config{
		WindowTitle:     config.WindowTitle,
		Scale:           config.Scale,
		CyclesPerSecond: config.CyclesPerSecond,
		AudioConfig: audio.Config{
			Frequency:     44100,
			BeepFrequency: 440.0,
			Volume:        0.1,
		},
		RendererConfig: renderer.Config{
			Width:  64,
			Height: 32,
			Scale:  config.Scale,
			Title:  config.WindowTitle,
		},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating emulator: %s\n", err.Error())
		os.Exit(1)
	}
	defer emu.Close()

	// Initialize emulator
	if err := emu.Initialize(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing emulator: %s\n", err.Error())
		os.Exit(1)
	}

	// Load ROM
	if err := emu.LoadROM(romPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error loading ROM: %s\n", err.Error())
		os.Exit(1)
	}

	// Run emulator
	emu.Run()
}
