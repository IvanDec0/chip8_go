package ui

import "chip8/internal/menu"

// Predefined colors
var (
	ColorBlack     = menu.Color{R: 0, G: 0, B: 0, A: 255}
	ColorWhite     = menu.Color{R: 255, G: 255, B: 255, A: 255}
	ColorGreen     = menu.Color{R: 0, G: 255, B: 0, A: 255}
	ColorRed       = menu.Color{R: 255, G: 0, B: 0, A: 255}
	ColorBlue      = menu.Color{R: 0, G: 0, B: 255, A: 255}
	ColorYellow    = menu.Color{R: 255, G: 255, B: 0, A: 255}
	ColorMagenta   = menu.Color{R: 255, G: 0, B: 255, A: 255}
	ColorCyan      = menu.Color{R: 0, G: 255, B: 255, A: 255}
	ColorGray      = menu.Color{R: 128, G: 128, B: 128, A: 255}
	ColorDarkGray  = menu.Color{R: 64, G: 64, B: 64, A: 255}
	ColorLightGray = menu.Color{R: 192, G: 192, B: 192, A: 255}
)

// Theme colors for different UI elements
var (
	ThemeBackground = ColorBlack
	ThemeText       = ColorWhite
	ThemeSelected   = ColorGreen
	ThemeBorder     = ColorGray
	ThemeError      = ColorRed
	ThemeSuccess    = ColorGreen
	ThemeWarning    = ColorYellow
	ThemeInfo       = ColorBlue
)

// GetThemeColor returns a themed color by name
func GetThemeColor(name string) menu.Color {
	switch name {
	case "background":
		return ThemeBackground
	case "text":
		return ThemeText
	case "selected":
		return ThemeSelected
	case "border":
		return ThemeBorder
	case "error":
		return ThemeError
	case "success":
		return ThemeSuccess
	case "warning":
		return ThemeWarning
	case "info":
		return ThemeInfo
	default:
		return ColorWhite
	}
}
