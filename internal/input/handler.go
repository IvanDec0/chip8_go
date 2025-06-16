package input

import (
	"chip8/internal/menu"

	"github.com/veandco/go-sdl2/sdl"
)

// KeyEvent represents a key event
type KeyEvent struct {
	Key     uint8
	Pressed bool
}

// Handler handles input events
type Handler struct {
	keyMap    map[sdl.Keycode]uint8
	keyStates [16]bool
}

// NewHandler creates a new input handler
func NewHandler() *Handler {
	return &Handler{
		keyMap: map[sdl.Keycode]uint8{
			sdl.K_1: 0x1, sdl.K_2: 0x2, sdl.K_3: 0x3, sdl.K_4: 0xC,
			sdl.K_q: 0x4, sdl.K_w: 0x5, sdl.K_e: 0x6, sdl.K_r: 0xD,
			sdl.K_a: 0x7, sdl.K_s: 0x8, sdl.K_d: 0x9, sdl.K_f: 0xE,
			sdl.K_z: 0xA, sdl.K_x: 0x0, sdl.K_c: 0xB, sdl.K_v: 0xF,
		},
	}
}

// ProcessEvents processes SDL events and returns key events and exit status
func (h *Handler) ProcessEvents() ([]KeyEvent, bool) {
	var keyEvents []KeyEvent
	quit := false

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			quit = true
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN && e.Keysym.Sym == sdl.K_ESCAPE {
				quit = true
			} else if key, ok := h.keyMap[e.Keysym.Sym]; ok {
				pressed := e.Type == sdl.KEYDOWN
				if h.keyStates[key] != pressed {
					h.keyStates[key] = pressed
					keyEvents = append(keyEvents, KeyEvent{
						Key:     key,
						Pressed: pressed,
					})
				}
			}
		}
	}

	return keyEvents, quit
}

// ExtendedHandler extends the basic handler with menu support
type ExtendedHandler struct {
	*Handler
	menuMode       bool
	menuKeyMap     map[sdl.Keycode]menu.MenuInputType
	lastKeyTime    map[sdl.Keycode]uint32
	keyRepeatDelay uint32 // Delay before key repeat starts (in milliseconds)
	keyRepeatRate  uint32 // Rate of key repeat (in milliseconds)
}

// NewExtendedHandler creates a new extended input handler
func NewExtendedHandler() *ExtendedHandler {
	return &ExtendedHandler{
		Handler:        NewHandler(),
		menuMode:       false,
		lastKeyTime:    make(map[sdl.Keycode]uint32),
		keyRepeatDelay: 500, // 500ms delay before repeat starts
		keyRepeatRate:  100, // 100ms between repeats
		menuKeyMap: map[sdl.Keycode]menu.MenuInputType{
			sdl.K_UP:        menu.MenuInputUp,
			sdl.K_DOWN:      menu.MenuInputDown,
			sdl.K_w:         menu.MenuInputUp,   // Alternative WASD
			sdl.K_s:         menu.MenuInputDown, // Alternative WASD
			sdl.K_RETURN:    menu.MenuInputSelect,
			sdl.K_SPACE:     menu.MenuInputSelect,
			sdl.K_ESCAPE:    menu.MenuInputBack,
			sdl.K_BACKSPACE: menu.MenuInputBack,
			sdl.K_f:         menu.MenuInputToggleFavorite, // 'F' key to toggle favorites
		},
	}
}

// SetMenuMode enables or disables menu input mode
func (eh *ExtendedHandler) SetMenuMode(enabled bool) {
	eh.menuMode = enabled
}

// IsMenuMode returns whether menu mode is enabled
func (eh *ExtendedHandler) IsMenuMode() bool {
	return eh.menuMode
}

// ProcessEventsExtended processes SDL events and returns key events, menu events, and exit status
func (eh *ExtendedHandler) ProcessEventsExtended() ([]KeyEvent, []menu.MenuInputEvent, bool) {
	var keyEvents []KeyEvent
	var menuEvents []menu.MenuInputEvent
	quit := false
	currentTime := sdl.GetTicks()

	// Handle key repeat for navigation keys
	if eh.menuMode {
		for keycode, lastTime := range eh.lastKeyTime {
			if menuInputType, ok := eh.menuKeyMap[keycode]; ok {
				// Only repeat navigation keys (up/down)
				if menuInputType == menu.MenuInputUp || menuInputType == menu.MenuInputDown {
					timeSincePress := currentTime - lastTime
					if timeSincePress > eh.keyRepeatDelay {
						// Calculate how many repeats should have occurred
						repeatsSinceDelay := (timeSincePress - eh.keyRepeatDelay) / eh.keyRepeatRate
						if repeatsSinceDelay > 0 {
							menuEvents = append(menuEvents, menu.MenuInputEvent{
								Type: menuInputType,
								Data: nil,
							})
							// Update last time to prevent excessive repeats
							eh.lastKeyTime[keycode] = currentTime - ((timeSincePress - eh.keyRepeatDelay) % eh.keyRepeatRate)
						}
					}
				}
			}
		}
	}

	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch e := event.(type) {
		case *sdl.QuitEvent:
			quit = true
		case *sdl.KeyboardEvent:
			if e.Type == sdl.KEYDOWN {
				// Handle ESC key - behavior depends on context
				if e.Keysym.Sym == sdl.K_ESCAPE {
					if eh.menuMode {
						// In menu mode, ESC means go back/exit
						menuEvents = append(menuEvents, menu.MenuInputEvent{
							Type: menu.MenuInputBack,
							Data: nil,
						})
					} else {
						// In emulator mode, ESC means quit
						quit = true
					}
					continue
				}

				// Handle menu navigation if in menu mode
				if eh.menuMode {
					if menuInputType, ok := eh.menuKeyMap[e.Keysym.Sym]; ok {
						menuEvents = append(menuEvents, menu.MenuInputEvent{
							Type: menuInputType,
							Data: nil,
						})
						// Track key press time for repeat functionality
						eh.lastKeyTime[e.Keysym.Sym] = currentTime
						continue
					}
				}

				// Handle CHIP-8 keys if not in menu mode
				if !eh.menuMode {
					if key, ok := eh.Handler.keyMap[e.Keysym.Sym]; ok {
						if eh.Handler.keyStates[key] != true {
							eh.Handler.keyStates[key] = true
							keyEvents = append(keyEvents, KeyEvent{
								Key:     key,
								Pressed: true,
							})
						}
					}
				}
			} else if e.Type == sdl.KEYUP {
				// Stop key repeat for menu keys
				if eh.menuMode {
					if _, ok := eh.menuKeyMap[e.Keysym.Sym]; ok {
						delete(eh.lastKeyTime, e.Keysym.Sym)
					}
				}

				// Handle CHIP-8 key releases if not in menu mode
				if !eh.menuMode {
					if key, ok := eh.Handler.keyMap[e.Keysym.Sym]; ok {
						if eh.Handler.keyStates[key] != false {
							eh.Handler.keyStates[key] = false
							keyEvents = append(keyEvents, KeyEvent{
								Key:     key,
								Pressed: false,
							})
						}
					}
				}
			}
		}
	}

	return keyEvents, menuEvents, quit
}
