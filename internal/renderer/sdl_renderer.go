package renderer

import (
	"github.com/veandco/go-sdl2/sdl"
)

// SDLRenderer implements Renderer using SDL2
type SDLRenderer struct {
	window   *sdl.Window
	renderer *sdl.Renderer
	config   Config
}

// NewSDLRenderer creates a new SDL-based renderer
func NewSDLRenderer() *SDLRenderer {
	return &SDLRenderer{}
}

// Initialize sets up the SDL window and renderer
func (r *SDLRenderer) Initialize(width, height, scale int, title string) error {
	r.config = Config{
		Width:  width,
		Height: height,
		Scale:  scale,
		Title:  title,
	}

	windowWidth := int32(width * scale)
	windowHeight := int32(height * scale)

	window, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight,
		sdl.WINDOW_SHOWN,
	)
	if err != nil {
		return err
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		window.Destroy()
		return err
	}

	r.window = window
	r.renderer = renderer
	return nil
}

// Clear clears the screen with black color
func (r *SDLRenderer) Clear() {
	r.renderer.SetDrawColor(0, 0, 0, 255)
	r.renderer.Clear()
}

// DrawPixel draws a pixel at the specified coordinates
func (r *SDLRenderer) DrawPixel(x, y int) {
	r.renderer.SetDrawColor(255, 255, 255, 255)
	rect := sdl.Rect{
		X: int32(x * r.config.Scale),
		Y: int32(y * r.config.Scale),
		W: int32(r.config.Scale),
		H: int32(r.config.Scale),
	}
	r.renderer.FillRect(&rect)
}

// Present presents the rendered frame
func (r *SDLRenderer) Present() {
	r.renderer.Present()
}

// Close shuts down the renderer and window
func (r *SDLRenderer) Close() {
	if r.renderer != nil {
		r.renderer.Destroy()
	}
	if r.window != nil {
		r.window.Destroy()
	}
}
