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

// ColorTheme represents a complete color theme
type ColorTheme struct {
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Background     menu.Color `json:"background"`
	Text           menu.Color `json:"text"`
	Selected       menu.Color `json:"selected"`
	Border         menu.Color `json:"border"`
	Error          menu.Color `json:"error"`
	Success        menu.Color `json:"success"`
	Warning        menu.Color `json:"warning"`
	Info           menu.Color `json:"info"`
	Highlight      menu.Color `json:"highlight"`
	Disabled       menu.Color `json:"disabled"`
	MenuBackground menu.Color `json:"menu_background"`
	MenuSelected   menu.Color `json:"menu_selected"`
	MenuBorder     menu.Color `json:"menu_border"`
}

// Predefined themes
var (
	// Default theme
	DefaultTheme = ColorTheme{
		Name:           "Default",
		Description:    "Standard CHIP-8 emulator theme",
		Background:     ColorBlack,
		Text:           ColorWhite,
		Selected:       ColorGreen,
		Border:         ColorGray,
		Error:          ColorRed,
		Success:        ColorGreen,
		Warning:        ColorYellow,
		Info:           ColorBlue,
		Highlight:      ColorCyan,
		Disabled:       ColorDarkGray,
		MenuBackground: ColorDarkGray,
		MenuSelected:   ColorGreen,
		MenuBorder:     ColorLightGray,
	}

	// High contrast theme for accessibility
	HighContrastTheme = ColorTheme{
		Name:           "High Contrast",
		Description:    "High contrast theme for better visibility",
		Background:     ColorBlack,
		Text:           ColorWhite,
		Selected:       menu.Color{R: 255, G: 255, B: 0, A: 255}, // Bright Yellow
		Border:         ColorWhite,
		Error:          menu.Color{R: 255, G: 100, B: 100, A: 255}, // Bright Red
		Success:        menu.Color{R: 100, G: 255, B: 100, A: 255}, // Bright Green
		Warning:        menu.Color{R: 255, G: 200, B: 0, A: 255},   // Bright Orange
		Info:           menu.Color{R: 100, G: 200, B: 255, A: 255}, // Bright Blue
		Highlight:      menu.Color{R: 255, G: 255, B: 255, A: 255}, // White
		Disabled:       menu.Color{R: 100, G: 100, B: 100, A: 255}, // Dark Gray
		MenuBackground: ColorBlack,
		MenuSelected:   menu.Color{R: 255, G: 255, B: 0, A: 255}, // Yellow
		MenuBorder:     ColorWhite,
	}

	// Dark theme with softer colors
	DarkTheme = ColorTheme{
		Name:           "Dark",
		Description:    "Dark theme with comfortable colors",
		Background:     menu.Color{R: 25, G: 25, B: 25, A: 255},    // Very Dark Gray
		Text:           menu.Color{R: 220, G: 220, B: 220, A: 255}, // Light Gray
		Selected:       menu.Color{R: 100, G: 200, B: 255, A: 255}, // Light Blue
		Border:         menu.Color{R: 100, G: 100, B: 100, A: 255}, // Medium Gray
		Error:          menu.Color{R: 255, G: 120, B: 120, A: 255}, // Light Red
		Success:        menu.Color{R: 120, G: 255, B: 120, A: 255}, // Light Green
		Warning:        menu.Color{R: 255, G: 200, B: 120, A: 255}, // Light Orange
		Info:           menu.Color{R: 120, G: 180, B: 255, A: 255}, // Light Blue
		Highlight:      menu.Color{R: 180, G: 180, B: 255, A: 255}, // Purple-ish
		Disabled:       menu.Color{R: 80, G: 80, B: 80, A: 255},    // Dark Gray
		MenuBackground: menu.Color{R: 40, G: 40, B: 40, A: 255},    // Dark Gray
		MenuSelected:   menu.Color{R: 100, G: 200, B: 255, A: 255}, // Light Blue
		MenuBorder:     menu.Color{R: 120, G: 120, B: 120, A: 255}, // Medium Gray
	}

	// Light theme
	LightTheme = ColorTheme{
		Name:           "Light",
		Description:    "Light theme for bright environments",
		Background:     menu.Color{R: 250, G: 250, B: 250, A: 255}, // Almost White
		Text:           menu.Color{R: 40, G: 40, B: 40, A: 255},    // Dark Gray
		Selected:       menu.Color{R: 0, G: 120, B: 200, A: 255},   // Blue
		Border:         menu.Color{R: 180, G: 180, B: 180, A: 255}, // Light Gray
		Error:          menu.Color{R: 200, G: 50, B: 50, A: 255},   // Dark Red
		Success:        menu.Color{R: 50, G: 150, B: 50, A: 255},   // Dark Green
		Warning:        menu.Color{R: 200, G: 140, B: 0, A: 255},   // Dark Orange
		Info:           menu.Color{R: 50, G: 100, B: 200, A: 255},  // Dark Blue
		Highlight:      menu.Color{R: 100, G: 150, B: 200, A: 255}, // Light Blue
		Disabled:       menu.Color{R: 160, G: 160, B: 160, A: 255}, // Gray
		MenuBackground: menu.Color{R: 240, G: 240, B: 240, A: 255}, // Light Gray
		MenuSelected:   menu.Color{R: 0, G: 120, B: 200, A: 255},   // Blue
		MenuBorder:     menu.Color{R: 200, G: 200, B: 200, A: 255}, // Medium Gray
	}

	// Retro green theme (classic computer style)
	RetroGreenTheme = ColorTheme{
		Name:           "Retro Green",
		Description:    "Classic green monochrome computer theme",
		Background:     ColorBlack,
		Text:           menu.Color{R: 0, G: 255, B: 100, A: 255},   // Bright Green
		Selected:       menu.Color{R: 100, G: 255, B: 150, A: 255}, // Brighter Green
		Border:         menu.Color{R: 0, G: 200, B: 80, A: 255},    // Medium Green
		Error:          menu.Color{R: 255, G: 100, B: 100, A: 255}, // Red (for contrast)
		Success:        menu.Color{R: 150, G: 255, B: 150, A: 255}, // Light Green
		Warning:        menu.Color{R: 255, G: 255, B: 100, A: 255}, // Yellow
		Info:           menu.Color{R: 100, G: 255, B: 200, A: 255}, // Cyan-ish
		Highlight:      menu.Color{R: 200, G: 255, B: 200, A: 255}, // Very Light Green
		Disabled:       menu.Color{R: 0, G: 100, B: 40, A: 255},    // Dark Green
		MenuBackground: menu.Color{R: 0, G: 40, B: 20, A: 255},     // Very Dark Green
		MenuSelected:   menu.Color{R: 100, G: 255, B: 150, A: 255}, // Bright Green
		MenuBorder:     menu.Color{R: 0, G: 150, B: 60, A: 255},    // Medium Green
	}

	// Amber theme (classic amber monitor style)
	AmberTheme = ColorTheme{
		Name:           "Amber",
		Description:    "Classic amber monochrome monitor theme",
		Background:     ColorBlack,
		Text:           menu.Color{R: 255, G: 180, B: 0, A: 255},   // Amber
		Selected:       menu.Color{R: 255, G: 220, B: 100, A: 255}, // Bright Amber
		Border:         menu.Color{R: 200, G: 140, B: 0, A: 255},   // Medium Amber
		Error:          menu.Color{R: 255, G: 100, B: 100, A: 255}, // Red
		Success:        menu.Color{R: 100, G: 255, B: 100, A: 255}, // Green
		Warning:        menu.Color{R: 255, G: 255, B: 100, A: 255}, // Yellow
		Info:           menu.Color{R: 100, G: 200, B: 255, A: 255}, // Blue
		Highlight:      menu.Color{R: 255, G: 255, B: 200, A: 255}, // Light Amber
		Disabled:       menu.Color{R: 100, G: 70, B: 0, A: 255},    // Dark Amber
		MenuBackground: menu.Color{R: 40, G: 30, B: 0, A: 255},     // Very Dark Amber
		MenuSelected:   menu.Color{R: 255, G: 220, B: 100, A: 255}, // Bright Amber
		MenuBorder:     menu.Color{R: 150, G: 110, B: 0, A: 255},   // Medium Amber
	}
)

// Available themes registry
var AvailableThemes = map[string]ColorTheme{
	"default":       DefaultTheme,
	"high_contrast": HighContrastTheme,
	"dark":          DarkTheme,
	"light":         LightTheme,
	"retro_green":   RetroGreenTheme,
	"amber":         AmberTheme,
}

// Current active theme
var ActiveTheme = DefaultTheme

// SetTheme sets the active theme
func SetTheme(themeName string) bool {
	if theme, exists := AvailableThemes[themeName]; exists {
		ActiveTheme = theme
		updateThemeColors()
		return true
	}
	return false
}

// GetCurrentTheme returns the current active theme
func GetCurrentTheme() ColorTheme {
	return ActiveTheme
}

// GetAvailableThemes returns list of available theme names
func GetAvailableThemes() []string {
	var themes []string
	for name := range AvailableThemes {
		themes = append(themes, name)
	}
	return themes
}

// GetThemeByName returns a theme by name
func GetThemeByName(name string) (ColorTheme, bool) {
	theme, exists := AvailableThemes[name]
	return theme, exists
}

// updateThemeColors updates the global theme colors
func updateThemeColors() {
	ThemeBackground = ActiveTheme.Background
	ThemeText = ActiveTheme.Text
	ThemeSelected = ActiveTheme.Selected
	ThemeBorder = ActiveTheme.Border
	ThemeError = ActiveTheme.Error
	ThemeSuccess = ActiveTheme.Success
	ThemeWarning = ActiveTheme.Warning
	ThemeInfo = ActiveTheme.Info
}

// GetContrastRatio calculates the contrast ratio between two colors
func GetContrastRatio(c1, c2 menu.Color) float64 {
	// Calculate relative luminance
	l1 := getRelativeLuminance(c1)
	l2 := getRelativeLuminance(c2)

	// Ensure l1 is the lighter color
	if l1 < l2 {
		l1, l2 = l2, l1
	}

	return (l1 + 0.05) / (l2 + 0.05)
}

// getRelativeLuminance calculates the relative luminance of a color
func getRelativeLuminance(c menu.Color) float64 {
	// Convert to sRGB
	r := float64(c.R) / 255.0
	g := float64(c.G) / 255.0
	b := float64(c.B) / 255.0

	// Apply gamma correction
	if r <= 0.03928 {
		r = r / 12.92
	} else {
		r = ((r + 0.055) / 1.055)
		r = r * r // Simplified power of 2.4
	}

	if g <= 0.03928 {
		g = g / 12.92
	} else {
		g = ((g + 0.055) / 1.055)
		g = g * g
	}

	if b <= 0.03928 {
		b = b / 12.92
	} else {
		b = ((b + 0.055) / 1.055)
		b = b * b
	}

	return 0.2126*r + 0.7152*g + 0.0722*b
}

// IsAccessible checks if a color combination meets WCAG AA accessibility standards
func IsAccessible(foreground, background menu.Color) bool {
	ratio := GetContrastRatio(foreground, background)
	return ratio >= 4.5 // WCAG AA standard for normal text
}

// IsAccessibleLarge checks if a color combination meets WCAG AA accessibility standards for large text
func IsAccessibleLarge(foreground, background menu.Color) bool {
	ratio := GetContrastRatio(foreground, background)
	return ratio >= 3.0 // WCAG AA standard for large text
}

// ValidateThemeAccessibility validates that a theme meets accessibility standards
func ValidateThemeAccessibility(theme ColorTheme) map[string]bool {
	results := make(map[string]bool)

	results["text_on_background"] = IsAccessible(theme.Text, theme.Background)
	results["selected_on_background"] = IsAccessible(theme.Selected, theme.Background)
	results["border_on_background"] = IsAccessible(theme.Border, theme.Background)
	results["error_on_background"] = IsAccessible(theme.Error, theme.Background)
	results["success_on_background"] = IsAccessible(theme.Success, theme.Background)
	results["warning_on_background"] = IsAccessible(theme.Warning, theme.Background)
	results["info_on_background"] = IsAccessible(theme.Info, theme.Background)
	results["menu_text_on_menu_bg"] = IsAccessible(theme.Text, theme.MenuBackground)
	results["menu_selected_on_menu_bg"] = IsAccessible(theme.MenuSelected, theme.MenuBackground)

	return results
}
