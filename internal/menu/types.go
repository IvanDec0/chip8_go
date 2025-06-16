package menu

import "time"

// EmulatorState represents the current state of the emulator
type EmulatorState int

const (
	StateMenu     EmulatorState = iota // Browsing ROMs in menu
	StateEmulator                      // Running CHIP-8 emulator
	StatePaused                        // Emulator paused (for future use)
)

// String returns a string representation of the emulator state
func (s EmulatorState) String() string {
	switch s {
	case StateMenu:
		return "Menu"
	case StateEmulator:
		return "Emulator"
	case StatePaused:
		return "Paused"
	default:
		return "Unknown"
	}
}

// MenuScreen represents different menu screens
type MenuScreen int

const (
	ScreenMain MenuScreen = iota
	ScreenBrowser
	ScreenRecent
	ScreenFavorites
	ScreenSearch
	ScreenSettings
	ScreenHelp
	ScreenExport
)

// String returns a string representation of the menu screen
func (s MenuScreen) String() string {
	switch s {
	case ScreenMain:
		return "Main Menu"
	case ScreenBrowser:
		return "ROM Browser"
	case ScreenRecent:
		return "Recent ROMs"
	case ScreenFavorites:
		return "Favorites"
	case ScreenSearch:
		return "Search"
	case ScreenSettings:
		return "Settings"
	case ScreenHelp:
		return "Help"
	case ScreenExport:
		return "Export"
	default:
		return "Unknown"
	}
}

// MenuData holds menu-specific information
type MenuData struct {
	SelectedROM   string
	ShouldLoadROM bool
	ShouldExit    bool
}

// MenuItem represents a menu item
type MenuItem struct {
	Text   string
	Action MenuAction
	Data   interface{}
}

// MenuAction represents what happens when item is selected
type MenuAction int

const (
	ActionBrowseROMs MenuAction = iota
	ActionLoadROM
	ActionExit
	ActionGoBack
	ActionShowRecent
	ActionShowFavorites
	ActionToggleFavorite
	ActionAddToRecent
	ActionStartSearch
	ActionClearSearch
	ActionShowSettings
	ActionToggleSort
	ActionNextSort
	ActionRefresh
	ActionShowHelp
	ActionShowHelpContent
	ActionShowExport
	ActionExportFavorites
	ActionExportRecent
	ActionExportStatistics
)

// String returns a string representation of the menu action
func (a MenuAction) String() string {
	switch a {
	case ActionBrowseROMs:
		return "BrowseROMs"
	case ActionLoadROM:
		return "LoadROM"
	case ActionExit:
		return "Exit"
	case ActionGoBack:
		return "GoBack"
	case ActionShowRecent:
		return "ShowRecent"
	case ActionShowFavorites:
		return "ShowFavorites"
	case ActionToggleFavorite:
		return "ToggleFavorite"
	case ActionAddToRecent:
		return "AddToRecent"
	case ActionStartSearch:
		return "StartSearch"
	case ActionClearSearch:
		return "ClearSearch"
	case ActionShowSettings:
		return "ShowSettings"
	case ActionToggleSort:
		return "ToggleSort"
	case ActionNextSort:
		return "NextSort"
	case ActionRefresh:
		return "Refresh"
	case ActionShowHelp:
		return "ShowHelp"
	case ActionShowHelpContent:
		return "ShowHelpContent"
	case ActionShowExport:
		return "ShowExport"
	case ActionExportFavorites:
		return "ExportFavorites"
	case ActionExportRecent:
		return "ExportRecent"
	case ActionExportStatistics:
		return "ExportStatistics"
	default:
		return "Unknown"
	}
}

// MenuInputEvent represents menu navigation events
type MenuInputEvent struct {
	Type MenuInputType
	Data interface{}
}

// MenuInputType represents different menu input types
type MenuInputType int

const (
	MenuInputUp MenuInputType = iota
	MenuInputDown
	MenuInputSelect
	MenuInputBack
	MenuInputExit
	MenuInputToggleFavorite
)

// String returns a string representation of the menu input type
func (t MenuInputType) String() string {
	switch t {
	case MenuInputUp:
		return "Up"
	case MenuInputDown:
		return "Down"
	case MenuInputSelect:
		return "Select"
	case MenuInputBack:
		return "Back"
	case MenuInputExit:
		return "Exit"
	case MenuInputToggleFavorite:
		return "ToggleFavorite"
	default:
		return "Unknown"
	}
}

// Color represents an RGBA color
type Color struct {
	R, G, B, A uint8
}

// ROMInfo contains information about a ROM file
type ROMInfo struct {
	Path        string    // Full file path
	Name        string    // Display name
	Size        int64     // File size in bytes
	ModTime     time.Time // Last modification time
	IsDirectory bool      // Whether this is a directory
}

// MenuTheme holds theme configuration for the menu
type MenuTheme struct {
	BackgroundColor Color
	TextColor       Color
	SelectedColor   Color
	BorderColor     Color
}

// DefaultMenuTheme returns a default menu theme
func DefaultMenuTheme() MenuTheme {
	return MenuTheme{
		BackgroundColor: Color{0, 0, 0, 255},       // Black
		TextColor:       Color{255, 255, 255, 255}, // White
		SelectedColor:   Color{0, 255, 0, 255},     // Green
		BorderColor:     Color{128, 128, 128, 255}, // Gray
	}
}

// RecentROM represents a recently played ROM for the interface
type RecentROM struct {
	Path       string
	Name       string
	LastPlayed string
	PlayCount  int
	IsFavorite bool
}

// FavoriteROM represents a favorited ROM for the interface
type FavoriteROM struct {
	Path       string
	Name       string
	CustomName string
	Rating     int
	Tags       []string
	DateAdded  string
}

// ROMMetadata represents ROM metadata for the interface
type ROMMetadata struct {
	Title       string
	Author      string
	Description []string
	Controls    []string
	Year        string
	System      string
}

// ROMBrowser interface for ROM file browsing functionality
type ROMBrowser interface {
	ScanDirectory(path string) ([]ROMInfo, error)
	GetROMs() []ROMInfo
	SetCurrentDirectory(path string) error
	GetCurrentDirectory() string
	NavigateUp() error
	Refresh() error

	// Phase 3 enhancements
	ScanDirectoryWithMetadata(path string) ([]ROMInfo, error)
	GetROMMetadata(romPath string) (*ROMMetadata, error)
	AddToRecent(romPath, romName string) error
	GetRecentROMs() []RecentROM
	GetFavorites() []FavoriteROM
	ToggleFavorite(romPath, romName string) error
	IsFavorite(romPath string) bool
	LoadUserData() error
}
