package emulator

import (
	"chip8/internal/audio"
	"chip8/internal/core"
	"chip8/internal/input"
	"chip8/internal/renderer"
	"time"
)

// Config holds emulator configuration
type Config struct {
	WindowTitle     string
	Scale           int
	CyclesPerSecond int
	AudioConfig     audio.Config
	RendererConfig  renderer.Config
}

// Emulator coordinates all emulator components
type Emulator struct {
	cpu      *core.CPU
	audio    audio.AudioSystem
	renderer renderer.Renderer
	input    *input.Handler
	config   Config
}

// New creates a new emulator instance
func New(config Config) (*Emulator, error) {
	// Initialize components
	cpu := core.NewCPU()

	audioSystem, err := audio.NewSDLAudioSystem(config.AudioConfig)
	if err != nil {
		return nil, err
	}

	sdlRenderer := renderer.NewSDLRenderer()
	inputHandler := input.NewHandler()

	return &Emulator{
		cpu:      cpu,
		audio:    audioSystem,
		renderer: sdlRenderer,
		input:    inputHandler,
		config:   config,
	}, nil
}

// Initialize initializes all emulator subsystems
func (e *Emulator) Initialize() error {
	// Initialize renderer
	err := e.renderer.Initialize(
		e.config.RendererConfig.Width,
		e.config.RendererConfig.Height,
		e.config.Scale,
		e.config.WindowTitle,
	)
	if err != nil {
		return err
	}

	// Initialize audio
	return e.audio.Start()
}

// LoadROM loads a ROM file
func (e *Emulator) LoadROM(path string) error {
	return e.cpu.LoadROM(path)
}

// Run starts the emulator main loop
func (e *Emulator) Run() {
	cycleDelay := time.Second / time.Duration(e.config.CyclesPerSecond)
	lastCycleTime := time.Now()
	lastTimerUpdate := time.Now()

	running := true
	for running {
		// Process input
		keyEvents, quit := e.input.ProcessEvents()
		if quit {
			running = false
		}

		// Handle key events
		for _, event := range keyEvents {
			e.cpu.SetKey(event.Key, event.Pressed)
		}

		// Update timers at 60Hz
		if time.Since(lastTimerUpdate) >= time.Second/60 {
			e.cpu.UpdateTimers()

			// Update audio based on sound timer
			if sdlAudio, ok := e.audio.(*audio.SDLAudioSystem); ok {
				sdlAudio.SetBeepState(e.cpu.IsSoundActive())
				sdlAudio.Update()
			}

			lastTimerUpdate = time.Now()
		}

		// Run CPU cycle
		if time.Since(lastCycleTime) >= cycleDelay {
			e.cpu.Cycle()
			lastCycleTime = time.Now()
		}

		// Render if needed
		if e.cpu.ShouldDraw() {
			e.render()
			e.cpu.ClearDrawFlag()
		}

		// Small delay to prevent 100% CPU usage
		time.Sleep(time.Millisecond)
	}
}

// render renders the current display state
func (e *Emulator) render() {
	e.renderer.Clear()

	display := e.cpu.GetDisplay()
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if display[y][x] == 1 {
				e.renderer.DrawPixel(x, y)
			}
		}
	}

	e.renderer.Present()
}

// Close shuts down the emulator
func (e *Emulator) Close() {
	e.audio.Close()
	e.renderer.Close()
}
