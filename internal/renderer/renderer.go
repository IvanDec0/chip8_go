package renderer

// Renderer defines the interface for display renderers
type Renderer interface {
	Initialize(width, height, scale int, title string) error
	Clear()
	DrawPixel(x, y int)
	Present()
	Close()
}

// Config holds renderer configuration
type Config struct {
	Width  int
	Height int
	Scale  int
	Title  string
}
