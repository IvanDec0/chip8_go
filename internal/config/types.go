package config

import (
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

// MenuConfig holds menu-specific configuration
type MenuConfig struct {
	Enabled             bool   `json:"enabled"`
	DefaultROMDirectory string `json:"default_rom_directory"`
	SortBy              string `json:"sort_by"` // "name", "date", "size", "recent"
	SortAscending       bool   `json:"sort_ascending"`
	ShowMetadata        bool   `json:"show_metadata"`
	ShowRecentROMs      bool   `json:"show_recent_roms"`
	MaxRecentROMs       int    `json:"max_recent_roms"`
	EnableSearch        bool   `json:"enable_search"`
	CaseSensitiveSearch bool   `json:"case_sensitive_search"`
}

// ThemeConfig holds visual theme configuration
type ThemeConfig struct {
	Name            string `json:"name"`
	BackgroundColor Color  `json:"background_color"`
	TextColor       Color  `json:"text_color"`
	SelectedColor   Color  `json:"selected_color"`
	BorderColor     Color  `json:"border_color"`
	AccentColor     Color  `json:"accent_color"`
	FontSize        int    `json:"font_size"`
}

// InputConfig holds input-related configuration
type InputConfig struct {
	KeyRepeatDelay    uint32            `json:"key_repeat_delay"`
	KeyRepeatInterval uint32            `json:"key_repeat_interval"`
	CustomKeyBindings map[string]string `json:"custom_key_bindings"`
	EnableMouseInput  bool              `json:"enable_mouse_input"`
}

// EmulatorConfig holds emulator-specific configuration
type EmulatorConfig struct {
	DefaultScale       int     `json:"default_scale"`
	DefaultSpeed       int     `json:"default_speed"`
	DefaultWindowTitle string  `json:"default_window_title"`
	AudioEnabled       bool    `json:"audio_enabled"`
	AudioVolume        float32 `json:"audio_volume"`
	AudioFrequency     int     `json:"audio_frequency"`
	BeepFrequency      float64 `json:"beep_frequency"`
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	MetadataCacheEnabled  bool          `json:"metadata_cache_enabled"`
	MetadataCacheTTL      time.Duration `json:"metadata_cache_ttl"`
	DirectoryCacheEnabled bool          `json:"directory_cache_enabled"`
	DirectoryCacheTTL     time.Duration `json:"directory_cache_ttl"`
	MaxCacheSize          int           `json:"max_cache_size"`
	CleanupInterval       time.Duration `json:"cleanup_interval"`
}

// Color represents an RGBA color
type Color struct {
	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

// UserConfig represents all user-configurable settings
type UserConfig struct {
	MenuConfig     MenuConfig     `json:"menu_config"`
	ThemeConfig    ThemeConfig    `json:"theme_config"`
	InputConfig    InputConfig    `json:"input_config"`
	EmulatorConfig EmulatorConfig `json:"emulator_config"`
	CacheConfig    CacheConfig    `json:"cache_config"`
	LastUpdated    time.Time      `json:"last_updated"`
	Version        string         `json:"version"`
}

// DefaultUserConfig returns a default user configuration
func DefaultUserConfig() *UserConfig {
	return &UserConfig{
		MenuConfig: MenuConfig{
			Enabled:             true,
			DefaultROMDirectory: "roms",
			SortBy:              "name",
			SortAscending:       true,
			ShowMetadata:        true,
			ShowRecentROMs:      true,
			MaxRecentROMs:       10,
			EnableSearch:        true,
			CaseSensitiveSearch: false,
		},
		ThemeConfig: ThemeConfig{
			Name:            "Default",
			BackgroundColor: Color{0, 0, 0, 255},       // Black
			TextColor:       Color{255, 255, 255, 255}, // White
			SelectedColor:   Color{0, 255, 0, 255},     // Green
			BorderColor:     Color{128, 128, 128, 255}, // Gray
			AccentColor:     Color{0, 150, 255, 255},   // Blue
			FontSize:        8,
		},
		InputConfig: InputConfig{
			KeyRepeatDelay:    250,
			KeyRepeatInterval: 50,
			CustomKeyBindings: make(map[string]string),
			EnableMouseInput:  false,
		},
		EmulatorConfig: EmulatorConfig{
			DefaultScale:       10,
			DefaultSpeed:       700,
			DefaultWindowTitle: "Chip8",
			AudioEnabled:       true,
			AudioVolume:        0.1,
			AudioFrequency:     44100,
			BeepFrequency:      440.0,
		},
		CacheConfig: CacheConfig{
			MetadataCacheEnabled:  true,
			MetadataCacheTTL:      time.Hour * 24,
			DirectoryCacheEnabled: true,
			DirectoryCacheTTL:     time.Minute * 5,
			MaxCacheSize:          1000,
			CleanupInterval:       time.Minute * 15,
		},
		LastUpdated: time.Now(),
		Version:     "1.0",
	}
}

// GetKeyCode converts a string key name to SDL keycode
func (ic *InputConfig) GetKeyCode(keyName string) (sdl.Keycode, bool) {
	keyMap := map[string]sdl.Keycode{
		"1": sdl.K_1, "2": sdl.K_2, "3": sdl.K_3, "4": sdl.K_4,
		"q": sdl.K_q, "w": sdl.K_w, "e": sdl.K_e, "r": sdl.K_r,
		"a": sdl.K_a, "s": sdl.K_s, "d": sdl.K_d, "f": sdl.K_f,
		"z": sdl.K_z, "x": sdl.K_x, "c": sdl.K_c, "v": sdl.K_v,
		"up": sdl.K_UP, "down": sdl.K_DOWN, "left": sdl.K_LEFT, "right": sdl.K_RIGHT,
		"enter": sdl.K_RETURN, "space": sdl.K_SPACE, "escape": sdl.K_ESCAPE,
		"backspace": sdl.K_BACKSPACE, "tab": sdl.K_TAB,
	}

	// Check custom bindings first
	if customKey, exists := ic.CustomKeyBindings[keyName]; exists {
		if keycode, exists := keyMap[customKey]; exists {
			return keycode, true
		}
	}

	// Check default mappings
	if keycode, exists := keyMap[keyName]; exists {
		return keycode, true
	}

	return 0, false
}

// SetCustomKeyBinding sets a custom key binding
func (ic *InputConfig) SetCustomKeyBinding(action, key string) {
	if ic.CustomKeyBindings == nil {
		ic.CustomKeyBindings = make(map[string]string)
	}
	ic.CustomKeyBindings[action] = key
}

// ToSDLColor converts Color to SDL color format
func (c Color) ToSDLColor() (uint8, uint8, uint8, uint8) {
	return c.R, c.G, c.B, c.A
}

// FromSDLColor creates Color from SDL color values
func FromSDLColor(r, g, b, a uint8) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// Clone creates a deep copy of the user config
func (uc *UserConfig) Clone() *UserConfig {
	clone := *uc

	// Deep copy custom key bindings
	clone.InputConfig.CustomKeyBindings = make(map[string]string)
	for k, v := range uc.InputConfig.CustomKeyBindings {
		clone.InputConfig.CustomKeyBindings[k] = v
	}

	return &clone
}

// Validate checks if the configuration is valid and corrects any issues
func (uc *UserConfig) Validate() {
	// Validate menu config
	if uc.MenuConfig.MaxRecentROMs <= 0 {
		uc.MenuConfig.MaxRecentROMs = 10
	}
	if uc.MenuConfig.MaxRecentROMs > 50 {
		uc.MenuConfig.MaxRecentROMs = 50
	}

	// Validate theme config
	if uc.ThemeConfig.FontSize <= 0 {
		uc.ThemeConfig.FontSize = 8
	}

	// Validate input config
	if uc.InputConfig.KeyRepeatDelay < 50 {
		uc.InputConfig.KeyRepeatDelay = 50
	}
	if uc.InputConfig.KeyRepeatInterval < 10 {
		uc.InputConfig.KeyRepeatInterval = 10
	}

	// Validate emulator config
	if uc.EmulatorConfig.DefaultScale <= 0 || uc.EmulatorConfig.DefaultScale > 20 {
		uc.EmulatorConfig.DefaultScale = 10
	}
	if uc.EmulatorConfig.DefaultSpeed < 100 || uc.EmulatorConfig.DefaultSpeed > 2000 {
		uc.EmulatorConfig.DefaultSpeed = 700
	}
	if uc.EmulatorConfig.AudioVolume < 0 || uc.EmulatorConfig.AudioVolume > 1 {
		uc.EmulatorConfig.AudioVolume = 0.1
	}

	// Validate cache config
	if uc.CacheConfig.MaxCacheSize <= 0 {
		uc.CacheConfig.MaxCacheSize = 1000
	}
	if uc.CacheConfig.MetadataCacheTTL <= 0 {
		uc.CacheConfig.MetadataCacheTTL = time.Hour * 24
	}
	if uc.CacheConfig.DirectoryCacheTTL <= 0 {
		uc.CacheConfig.DirectoryCacheTTL = time.Minute * 5
	}
	if uc.CacheConfig.CleanupInterval <= 0 {
		uc.CacheConfig.CleanupInterval = time.Minute * 15
	}
}
