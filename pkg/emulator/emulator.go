package emulator

import (
	"chip8/internal/audio"
	"chip8/internal/browser"
	"chip8/internal/core"
	"chip8/internal/input"
	"chip8/internal/menu"
	"chip8/internal/renderer"
	"fmt"
	"time"
)

// MenuConfig holds menu-specific configuration
type MenuConfig struct {
	Enabled       bool
	DefaultROMDir string
	Theme         menu.MenuTheme
}

// BrowserConfig holds browser-specific configuration
type BrowserConfig struct {
	DefaultDirectory string
	ShowHiddenFiles  bool
	SortBy           string // "name", "date", "size"
}

// Config holds emulator configuration
type Config struct {
	WindowTitle     string
	Scale           int
	CyclesPerSecond int
	AudioConfig     audio.Config
	RendererConfig  renderer.Config
	MenuConfig      MenuConfig    // New
	BrowserConfig   BrowserConfig // New
}

// Emulator coordinates all emulator components
type Emulator struct {
	cpu      *core.CPU
	audio    audio.AudioSystem
	renderer renderer.Renderer
	input    *input.ExtendedHandler // Changed to ExtendedHandler
	config   Config
	// New state management fields
	stateManager *menu.StateManager
	menuManager  *menu.Manager
	menuRenderer *menu.MenuRenderer // Add menu renderer
	browser      browser.Browser    // Add browser instance
	initialized  bool
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
	inputHandler := input.NewExtendedHandler() // Use ExtendedHandler

	// Create browser instance
	fileBrowser := browser.NewFileBrowser(config.BrowserConfig.DefaultDirectory)

	// Create state and menu managers
	stateManager := menu.NewStateManager()
	menuManager := menu.NewManager()

	// Set the shared state manager in menu manager
	menuManager.SetStateManager(stateManager)

	// Create menu renderer
	var menuRenderer *menu.MenuRenderer
	if config.MenuConfig.Enabled {
		menuRenderer = menu.NewMenuRenderer(
			sdlRenderer,
			config.RendererConfig.Width*config.Scale,
			config.RendererConfig.Height*config.Scale,
			config.Scale,
			config.MenuConfig.Theme,
		)
	}

	// Set browser reference in menu manager
	menuManager.SetBrowser(fileBrowser)

	// Set initial state based on configuration
	if config.MenuConfig.Enabled {
		stateManager.TransitionTo(menu.StateMenu)
		inputHandler.SetMenuMode(true)
	} else {
		stateManager.TransitionTo(menu.StateEmulator)
		inputHandler.SetMenuMode(false)
	}

	return &Emulator{
		cpu:          cpu,
		audio:        audioSystem,
		renderer:     sdlRenderer,
		input:        inputHandler,
		config:       config,
		stateManager: stateManager,
		menuManager:  menuManager,
		menuRenderer: menuRenderer,
		browser:      fileBrowser,
		initialized:  false,
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
	err = e.audio.Start()
	if err != nil {
		return err
	}

	// Initialize menu system
	err = e.initializeMenu()
	if err != nil {
		return err
	}

	e.initialized = true
	return nil
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
		currentState := e.stateManager.GetCurrentState()

		// Process input based on current state
		keyEvents, menuEvents, quit := e.input.ProcessEventsExtended()
		if quit {
			running = false
			continue
		}

		// Handle state-specific logic
		switch currentState {
		case menu.StateMenu:
			// Ensure menu has items
			if len(e.menuManager.GetMenuItems()) == 0 {
				// Force reload main menu items
				e.menuManager.SetMenuItems([]menu.MenuItem{
					{
						Text:   "Browse ROMs",
						Action: menu.ActionBrowseROMs,
						Data:   nil,
					},
					{
						Text:   "Exit",
						Action: menu.ActionExit,
						Data:   nil,
					},
				})
			}

			// Process menu events
			for _, event := range menuEvents {
				e.menuManager.ProcessMenuInput(event)
			}

			// Check for state transitions
			if e.stateManager.ShouldExit() {
				running = false
				continue
			}

			if e.stateManager.ShouldLoadROM() {
				// Load the selected ROM and transition to emulator
				romPath := e.stateManager.GetSelectedROM()
				if err := e.LoadROM(romPath); err == nil {
					e.stateManager.TransitionTo(menu.StateEmulator)
					e.input.SetMenuMode(false)
				} else {
					fmt.Printf("Failed to load ROM: %v\n", err)
				}
				e.stateManager.SetShouldLoadROM(false)
			}

			// Render menu
			if e.menuRenderer != nil {
				e.renderMenu()
			} else {
				// Fallback rendering if no menu renderer
				e.renderer.ClearWithColor(menu.Color{R: 0, G: 0, B: 0, A: 255})
				e.renderer.Present()
			}

		case menu.StateEmulator:
			// Handle CHIP-8 key events
			for _, event := range keyEvents {
				e.cpu.SetKey(event.Key, event.Pressed)
			}

			// Handle menu events (ESC to return to menu)
			for _, event := range menuEvents {
				if event.Type == menu.MenuInputBack && e.config.MenuConfig.Enabled {
					fmt.Println("ESC pressed, returning to menu...")
					e.stateManager.TransitionTo(menu.StateMenu)
					e.input.SetMenuMode(true)
				}
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

		case menu.StatePaused:
			// Handle paused state (future implementation)
			// For now, just process menu events
			for _, event := range menuEvents {
				e.menuManager.ProcessMenuInput(event)
			}
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

// renderMenu renders the current menu state
func (e *Emulator) renderMenu() {
	currentState := e.stateManager.GetCurrentState()

	switch currentState {
	case menu.StateMenu:
		// Render main menu
		items := e.menuManager.GetMenuItems()
		selected := e.menuManager.GetSelectedItem()
		e.menuRenderer.RenderMainMenu(items, selected)

	default:
		// Check if we're in ROM browser mode by looking at menu items
		items := e.menuManager.GetMenuItems()
		if len(items) > 0 && items[0].Action == menu.ActionGoBack {
			// We're in ROM browser mode
			browser := e.menuManager.GetBrowser()
			if browser != nil {
				roms := browser.GetROMs()
				selected := e.menuManager.GetSelectedItem()
				currentPath := browser.GetCurrentDirectory()
				e.menuRenderer.RenderROMBrowser(roms, selected, currentPath)
			}
		} else {
			// Fallback to main menu
			selected := e.menuManager.GetSelectedItem()
			e.menuRenderer.RenderMainMenu(items, selected)
		}
	}
}

// Close shuts down the emulator
func (e *Emulator) Close() {
	e.audio.Close()
	e.renderer.Close()
}

// GetStateManager returns the state manager
func (e *Emulator) GetStateManager() *menu.StateManager {
	return e.stateManager
}

// GetMenuManager returns the menu manager
func (e *Emulator) GetMenuManager() *menu.Manager {
	return e.menuManager
}

// GetBrowser returns the ROM browser instance
func (e *Emulator) GetBrowser() menu.ROMBrowser {
	return e.browser
}

// GetMenuRenderer returns the menu renderer instance
func (e *Emulator) GetMenuRenderer() *menu.MenuRenderer {
	return e.menuRenderer
}

// SetMenuRenderer sets the menu renderer (useful for testing)
func (e *Emulator) SetMenuRenderer(renderer *menu.MenuRenderer) {
	e.menuRenderer = renderer
}

// SetMenuEnabled enables or disables menu functionality
func (e *Emulator) SetMenuEnabled(enabled bool) {
	e.config.MenuConfig.Enabled = enabled
	if enabled {
		e.input.SetMenuMode(true)
		e.stateManager.TransitionTo(menu.StateMenu)
	} else {
		e.input.SetMenuMode(false)
		e.stateManager.TransitionTo(menu.StateEmulator)
	}
}

// loadSelectedROM loads the ROM selected in the menu
func (e *Emulator) loadSelectedROM() error {
	romPath := e.stateManager.GetSelectedROM()
	if romPath == "" {
		return fmt.Errorf("no ROM selected")
	}

	err := e.LoadROM(romPath)
	if err != nil {
		return fmt.Errorf("failed to load ROM %s: %w", romPath, err)
	}

	e.stateManager.TransitionTo(menu.StateEmulator)
	e.input.SetMenuMode(false)
	return nil
}

// initializeMenu initializes the menu system
func (e *Emulator) initializeMenu() error {
	if !e.config.MenuConfig.Enabled {
		return nil
	}

	// Set the default directory for the browser
	if e.config.BrowserConfig.DefaultDirectory != "" {
		err := e.browser.SetCurrentDirectory(e.config.BrowserConfig.DefaultDirectory)
		if err != nil {
			// Try fallback with ScanDirectoryWithFallback
			if fb, ok := e.browser.(*browser.FileBrowser); ok {
				_, fallbackErr := fb.ScanDirectoryWithFallback(e.config.BrowserConfig.DefaultDirectory)
				if fallbackErr != nil {
					// Log error but don't fail - we can still show main menu
					fmt.Printf("Warning: Failed to load ROM directory: %v\n", fallbackErr)
				}
			}
		}
	}

	return nil
}
