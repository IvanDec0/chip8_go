package input

import (
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
