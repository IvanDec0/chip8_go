package renderer

import "chip8/internal/menu"

// Renderer defines the interface for display renderers
type Renderer interface {
	Initialize(width, height, scale int, title string) error
	Clear()
	DrawPixel(x, y int)
	Present()
	Close()
	// New methods for menu support
	RenderText(text string, x, y int, color menu.Color) error
	SetBackgroundColor(color menu.Color)
	ClearWithColor(color menu.Color)
}

// MenuRenderer interface for menu-specific rendering
type MenuRenderer interface {
	RenderMainMenu(items []menu.MenuItem, selected int) error
	RenderROMBrowser(roms []menu.ROMInfo, selected int, currentPath string) error
	RenderText(text string, x, y int, color menu.Color) error
	Clear() error
	Present() error
}

// Config holds renderer configuration
type Config struct {
	Width  int
	Height int
	Scale  int
	Title  string
}
