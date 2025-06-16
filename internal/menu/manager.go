package menu

import (
	"chip8/internal/config"
	"fmt"
	"strings"
)

// Manager coordinates menu operations
type Manager struct {
	stateManager  *StateManager
	theme         MenuTheme
	selectedItem  int
	menuItems     []MenuItem
	browser       ROMBrowser
	currentScreen MenuScreen
	searchManager *SearchManager
	sortManager   *SortManager
	configManager *config.ConfigManager
	searchQuery   string
	searchActive  bool
	currentHelp   string
	exportStatus  string
}

// NewManager creates a new menu manager
func NewManager() *Manager {
	return &Manager{
		stateManager:  nil, // Will be set externally
		theme:         DefaultMenuTheme(),
		selectedItem:  0,
		menuItems:     createMainMenuItems(),
		currentScreen: ScreenMain,
		searchManager: NewSearchManager(),
		sortManager:   NewSortManager(),
		configManager: nil, // Will be set externally
		searchQuery:   "",
		searchActive:  false,
		currentHelp:   "",
		exportStatus:  "",
	}
}

// SetConfigManager sets the configuration manager
func (m *Manager) SetConfigManager(configManager *config.ConfigManager) {
	m.configManager = configManager
}

// GetStateManager returns the state manager
func (m *Manager) GetStateManager() *StateManager {
	return m.stateManager
}

// SetStateManager sets the state manager
func (m *Manager) SetStateManager(stateManager *StateManager) {
	m.stateManager = stateManager
}

// GetSelectedItem returns the currently selected menu item index
func (m *Manager) GetSelectedItem() int {
	return m.selectedItem
}

// SetSelectedItem sets the currently selected menu item index
func (m *Manager) SetSelectedItem(index int) {
	if index >= 0 && index < len(m.menuItems) {
		m.selectedItem = index
	}
}

// GetMenuItems returns the current menu items
func (m *Manager) GetMenuItems() []MenuItem {
	return m.menuItems
}

// SetMenuItems sets the current menu items
func (m *Manager) SetMenuItems(items []MenuItem) {
	m.menuItems = items
	if m.selectedItem >= len(items) {
		m.selectedItem = 0
	}
}

// NavigateUp moves selection up in the menu
func (m *Manager) NavigateUp() {
	if m.selectedItem > 0 {
		m.selectedItem--
	} else {
		// Wrap to bottom
		m.selectedItem = len(m.menuItems) - 1
	}
}

// NavigateDown moves selection down in the menu
func (m *Manager) NavigateDown() {
	if m.selectedItem < len(m.menuItems)-1 {
		m.selectedItem++
	} else {
		// Wrap to top
		m.selectedItem = 0
	}
}

// SelectCurrentItem handles selection of the current menu item
func (m *Manager) SelectCurrentItem() {
	if m.selectedItem >= 0 && m.selectedItem < len(m.menuItems) {
		item := m.menuItems[m.selectedItem]
		m.handleMenuAction(item.Action, item.Data)
	}
}

// handleMenuAction processes menu actions
func (m *Manager) handleMenuAction(action MenuAction, data interface{}) {
	switch action {
	case ActionBrowseROMs:
		if data != nil {
			// Navigate to directory
			if dirPath, ok := data.(string); ok && m.browser != nil {
				err := m.browser.SetCurrentDirectory(dirPath)
				if err == nil {
					m.currentScreen = ScreenBrowser
					m.LoadROMBrowser() // Reload the browser with new directory
				} else {
					fmt.Printf("Error setting directory: %v\n", err)
				}
			}
		} else {
			// Initial ROM browser load
			m.currentScreen = ScreenBrowser
			err := m.LoadROMBrowser()
			if err != nil {
				fmt.Printf("Error loading ROM browser: %v\n", err)
			}
		}
	case ActionLoadROM:
		if romPath, ok := data.(string); ok {
			// Add to recent before loading
			if m.browser != nil {
				romName := extractROMName(romPath)
				m.browser.AddToRecent(romPath, romName)
			}

			if m.stateManager != nil {
				m.stateManager.SetSelectedROM(romPath)
				m.stateManager.SetShouldLoadROM(true)
				m.stateManager.TransitionTo(StateEmulator)
			} else {
				fmt.Println("Error: StateManager is nil!")
			}
		}
	case ActionShowRecent:
		m.currentScreen = ScreenRecent
		m.LoadRecentROMs()
	case ActionShowFavorites:
		m.currentScreen = ScreenFavorites
		m.LoadFavorites()
	case ActionToggleFavorite:
		if romPath, ok := data.(string); ok && m.browser != nil {
			romName := extractROMName(romPath)
			m.browser.ToggleFavorite(romPath, romName)
			// Refresh current screen
			m.refreshCurrentScreen()
		}
	case ActionStartSearch:
		m.currentScreen = ScreenSearch
		m.searchActive = true
		m.searchQuery = ""
		m.LoadSearchScreen()
	case ActionClearSearch:
		m.searchActive = false
		m.searchQuery = ""
		m.searchManager.ClearSearch()
		m.currentScreen = ScreenBrowser
		m.LoadROMBrowser()
	case ActionNextSort:
		m.sortManager.NextSort()
		m.refreshCurrentScreen()
	case ActionToggleSort:
		m.sortManager.ToggleDirection()
		m.refreshCurrentScreen()
	case ActionRefresh:
		if m.browser != nil {
			m.browser.Refresh()
		}
		m.refreshCurrentScreen()
	case ActionShowSettings:
		m.currentScreen = ScreenSettings
		m.LoadSettings()
	case ActionShowHelp:
		m.currentScreen = ScreenHelp
		m.LoadHelpScreen()
	case ActionShowExport:
		m.currentScreen = ScreenExport
		m.LoadExportScreen()
	case ActionExportFavorites:
		m.handleExportAction("favorites")
	case ActionExportRecent:
		m.handleExportAction("recent")
	case ActionExportStatistics:
		m.handleExportAction("statistics")
	case ActionExit:
		if m.stateManager != nil {
			m.stateManager.SetShouldExit(true)
		} else {
			fmt.Println("Error: StateManager is nil!")
		}
	case ActionGoBack:
		m.handleGoBack()
	case ActionShowHelpContent:
		if topic, ok := data.(string); ok {
			m.showHelpContent(topic)
		}
	}
}

// createMainMenuItems creates the default main menu items
func createMainMenuItems() []MenuItem {
	return []MenuItem{
		{
			Text:   "Browse ROMs",
			Action: ActionBrowseROMs,
			Data:   nil,
		},
		{
			Text:   "Recent ROMs",
			Action: ActionShowRecent,
			Data:   nil,
		},
		{
			Text:   "Favorites",
			Action: ActionShowFavorites,
			Data:   nil,
		},
		{
			Text:   "Search ROMs",
			Action: ActionStartSearch,
			Data:   nil,
		},
		{
			Text:   "Settings",
			Action: ActionShowSettings,
			Data:   nil,
		},
		{
			Text:   "Help",
			Action: ActionShowHelp,
			Data:   nil,
		},
		{
			Text:   "Export Data",
			Action: ActionShowExport,
			Data:   nil,
		},
		{
			Text:   "Exit",
			Action: ActionExit,
			Data:   nil,
		},
	}
}

// ProcessMenuInput processes menu input events
func (m *Manager) ProcessMenuInput(event MenuInputEvent) {
	switch event.Type {
	case MenuInputUp:
		m.NavigateUp()
	case MenuInputDown:
		m.NavigateDown()
	case MenuInputSelect:
		m.SelectCurrentItem()
	case MenuInputBack:
		if len(m.menuItems) > 0 && m.menuItems[0].Action == ActionGoBack {
			// We're in a submenu, go back
			m.handleMenuAction(ActionGoBack, nil)
		} else {
			// We're in main menu, exit
			m.handleMenuAction(ActionExit, nil)
		}
	case MenuInputExit:
		m.handleMenuAction(ActionExit, nil)
	case MenuInputToggleFavorite:
		m.HandleToggleFavorite()
	}
}

// SetBrowser sets the ROM browser instance
func (m *Manager) SetBrowser(browser ROMBrowser) {
	m.browser = browser
}

// GetBrowser returns the ROM browser instance
func (m *Manager) GetBrowser() ROMBrowser {
	return m.browser
}

// LoadROMBrowser switches to ROM browser mode
func (m *Manager) LoadROMBrowser() error {
	if m.browser == nil {
		fmt.Println("No browser available, returning nil")
		return nil // No browser available
	}

	currentDir := m.browser.GetCurrentDirectory()
	// Scan current directory with metadata
	roms, err := m.browser.ScanDirectoryWithMetadata(currentDir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return err
	}

	// Apply sorting
	if m.sortManager != nil {
		roms = m.sortManager.SortROMs(roms)
	}

	// Create menu items from ROM list
	var items []MenuItem

	// Add back option
	items = append(items, MenuItem{
		Text:   "← Back to Main Menu",
		Action: ActionGoBack,
		Data:   nil,
	})

	// Add ROM items
	for _, rom := range roms {
		icon := "📁"
		if !rom.IsDirectory {
			icon = "🎮"
			if m.browser.IsFavorite(rom.Path) {
				icon = "⭐"
			}
		}

		items = append(items, MenuItem{
			Text:   icon + " " + rom.Name,
			Action: getActionForROM(rom),
			Data:   rom.Path,
		})
	}

	// Add action items at the bottom
	items = append(items, MenuItem{
		Text:   "🔍 Search",
		Action: ActionStartSearch,
		Data:   nil,
	})

	items = append(items, MenuItem{
		Text:   "🔄 Refresh",
		Action: ActionRefresh,
		Data:   nil,
	})

	m.SetMenuItems(items)
	return nil
}

// extractROMName extracts ROM name from path
func extractROMName(romPath string) string {
	// Simple implementation - extract filename without extension
	parts := strings.Split(romPath, "/")
	if len(parts) == 0 {
		return romPath
	}
	filename := parts[len(parts)-1]
	// Remove extension
	if dotIndex := strings.LastIndex(filename, "."); dotIndex > 0 {
		return filename[:dotIndex]
	}
	return filename
}

// refreshCurrentScreen refreshes the current menu screen
func (m *Manager) refreshCurrentScreen() {
	switch m.currentScreen {
	case ScreenMain:
		m.menuItems = createMainMenuItems()
	case ScreenBrowser:
		m.LoadROMBrowser()
	case ScreenRecent:
		m.LoadRecentROMs()
	case ScreenFavorites:
		m.LoadFavorites()
	case ScreenSearch:
		m.LoadSearchScreen()
	case ScreenSettings:
		m.LoadSettings()
	}
}

// handleGoBack handles back navigation based on current screen
func (m *Manager) handleGoBack() {
	switch m.currentScreen {
	case ScreenMain:
		// Already at main, exit
		m.handleMenuAction(ActionExit, nil)
	case ScreenBrowser:
		// Try to navigate up in directory, or go to main
		if m.browser != nil {
			if err := m.browser.NavigateUp(); err == nil {
				m.LoadROMBrowser()
			} else {
				m.currentScreen = ScreenMain
				m.menuItems = createMainMenuItems()
				m.selectedItem = 0
			}
		} else {
			m.currentScreen = ScreenMain
			m.menuItems = createMainMenuItems()
			m.selectedItem = 0
		}
	case ScreenHelp:
		// If showing help content, go back to help menu; otherwise go to main
		if m.currentHelp != "" {
			m.ClearCurrentHelp()
			m.LoadHelpScreen()
		} else {
			m.currentScreen = ScreenMain
			m.menuItems = createMainMenuItems()
			m.selectedItem = 0
		}
	case ScreenExport:
		// Clear export status and go back to main
		m.exportStatus = ""
		m.currentScreen = ScreenMain
		m.menuItems = createMainMenuItems()
		m.selectedItem = 0
	default:
		// Go back to main menu
		m.currentScreen = ScreenMain
		m.menuItems = createMainMenuItems()
		m.selectedItem = 0
	}
}

// LoadRecentROMs loads the recent ROMs screen
func (m *Manager) LoadRecentROMs() {
	if m.browser == nil {
		return
	}

	recent := m.browser.GetRecentROMs()
	var items []MenuItem

	// Add back option
	items = append(items, MenuItem{
		Text:   "← Back to Main Menu",
		Action: ActionGoBack,
		Data:   nil,
	})

	// Add recent ROM items
	for _, rom := range recent {
		icon := "🎮"
		if rom.IsFavorite {
			icon = "⭐"
		}

		displayText := fmt.Sprintf("%s %s (played %d times, last: %s)",
			icon, rom.Name, rom.PlayCount, rom.LastPlayed)

		items = append(items, MenuItem{
			Text:   displayText,
			Action: ActionLoadROM,
			Data:   rom.Path,
		})
	}

	if len(recent) == 0 {
		items = append(items, MenuItem{
			Text:   "No recent ROMs",
			Action: ActionGoBack,
			Data:   nil,
		})
	}

	m.SetMenuItems(items)
}

// LoadFavorites loads the favorites screen
func (m *Manager) LoadFavorites() {
	if m.browser == nil {
		return
	}

	favorites := m.browser.GetFavorites()
	var items []MenuItem

	// Add back option
	items = append(items, MenuItem{
		Text:   "← Back to Main Menu",
		Action: ActionGoBack,
		Data:   nil,
	})

	// Add favorite ROM items
	for _, fav := range favorites {
		displayName := fav.Name
		if fav.CustomName != "" {
			displayName = fav.CustomName
		}

		rating := ""
		for i := 0; i < fav.Rating; i++ {
			rating += "⭐"
		}

		displayText := fmt.Sprintf("⭐ %s %s", displayName, rating)

		items = append(items, MenuItem{
			Text:   displayText,
			Action: ActionLoadROM,
			Data:   fav.Path,
		})
	}

	if len(favorites) == 0 {
		items = append(items, MenuItem{
			Text:   "No favorite ROMs",
			Action: ActionGoBack,
			Data:   nil,
		})
	}

	m.SetMenuItems(items)
}

// LoadSearchScreen loads the search screen
func (m *Manager) LoadSearchScreen() {
	var items []MenuItem

	// Add back option
	items = append(items, MenuItem{
		Text:   "← Back to Browser",
		Action: ActionClearSearch,
		Data:   nil,
	})

	// Show search query
	searchText := "Search: " + m.searchQuery
	if m.searchQuery == "" {
		searchText = "Search: (type to search)"
	}

	items = append(items, MenuItem{
		Text:   searchText,
		Action: ActionStartSearch,
		Data:   nil,
	})

	// Show search results if we have a query
	if m.searchQuery != "" && m.browser != nil {
		roms := m.browser.GetROMs()
		filter := m.searchManager.GetDefaultFilter()
		m.searchManager.SetQuery(m.searchQuery)
		results := m.searchManager.Search(roms, filter)

		for _, rom := range results {
			icon := "📁"
			if !rom.IsDirectory {
				icon = "🎮"
				if m.browser.IsFavorite(rom.Path) {
					icon = "⭐"
				}
			}

			items = append(items, MenuItem{
				Text:   icon + " " + rom.Name,
				Action: getActionForROM(rom),
				Data:   rom.Path,
			})
		}

		if len(results) == 0 {
			items = append(items, MenuItem{
				Text:   "No matches found",
				Action: ActionStartSearch,
				Data:   nil,
			})
		}
	}

	m.SetMenuItems(items)
}

// LoadSettings loads the settings screen
func (m *Manager) LoadSettings() {
	var items []MenuItem

	// Add back option
	items = append(items, MenuItem{
		Text:   "← Back to Main Menu",
		Action: ActionGoBack,
		Data:   nil,
	})

	// Show current sort setting
	sortCriteria, ascending := m.sortManager.GetCurrentSort()
	direction := "Ascending"
	if !ascending {
		direction = "Descending"
	}

	items = append(items, MenuItem{
		Text:   fmt.Sprintf("Sort: %s (%s)", sortCriteria.String(), direction),
		Action: ActionNextSort,
		Data:   nil,
	})

	items = append(items, MenuItem{
		Text:   "Toggle Sort Direction",
		Action: ActionToggleSort,
		Data:   nil,
	})

	items = append(items, MenuItem{
		Text:   "Refresh ROMs",
		Action: ActionRefresh,
		Data:   nil,
	})

	m.SetMenuItems(items)
}

// LoadHelpScreen loads the help screen
func (m *Manager) LoadHelpScreen() {
	helpItems := []MenuItem{
		{
			Text:   "Navigation Controls",
			Action: ActionShowHelpContent,
			Data:   "navigation",
		},
		{
			Text:   "Game Controls",
			Action: ActionShowHelpContent,
			Data:   "controls",
		},
		{
			Text:   "Menu Features",
			Action: ActionShowHelpContent,
			Data:   "features",
		},
		{
			Text:   "Emulator Settings",
			Action: ActionShowHelpContent,
			Data:   "settings",
		},
		{
			Text:   "ROM Browser Guide",
			Action: ActionShowHelpContent,
			Data:   "browser",
		},
		{
			Text:   "Back to Main Menu",
			Action: ActionGoBack,
			Data:   nil,
		},
	}
	m.menuItems = helpItems
	m.selectedItem = 0
	m.currentHelp = "" // Clear current help content
}

// LoadExportScreen loads the export options screen
func (m *Manager) LoadExportScreen() {
	exportItems := []MenuItem{
		{
			Text:   "Export Favorites List",
			Action: ActionExportFavorites,
			Data:   nil,
		},
		{
			Text:   "Export Recent ROMs",
			Action: ActionExportRecent,
			Data:   nil,
		},
		{
			Text:   "Export Usage Statistics",
			Action: ActionExportStatistics,
			Data:   nil,
		},
		{
			Text:   "Back to Main Menu",
			Action: ActionGoBack,
			Data:   nil,
		},
	}
	m.menuItems = exportItems
	m.selectedItem = 0
}

// handleExportAction handles export operations
func (m *Manager) handleExportAction(exportType string) {
	if m.browser == nil {
		m.exportStatus = "Error: Browser not available for export"
		return
	}

	var filename string

	switch exportType {
	case "favorites":
		favorites := m.browser.GetFavorites()
		if len(favorites) == 0 {
			m.exportStatus = "No favorites to export"
			return
		}

		// Convert to export format
		exportData := make([]map[string]interface{}, len(favorites))
		for i, fav := range favorites {
			exportData[i] = map[string]interface{}{
				"path":        fav.Path,
				"name":        fav.Name,
				"custom_name": fav.CustomName,
				"rating":      fav.Rating,
				"tags":        fav.Tags,
				"date_added":  fav.DateAdded,
			}
		}

		filename = "favorites_export.json"
		m.exportStatus = fmt.Sprintf("Exported %d favorites to %s", len(favorites), filename)

	case "recent":
		recent := m.browser.GetRecentROMs()
		if len(recent) == 0 {
			m.exportStatus = "No recent ROMs to export"
			return
		}

		// Convert to export format
		exportData := make([]map[string]interface{}, len(recent))
		for i, rom := range recent {
			exportData[i] = map[string]interface{}{
				"path":        rom.Path,
				"name":        rom.Name,
				"last_played": rom.LastPlayed,
				"play_count":  rom.PlayCount,
			}
		}

		filename = "recent_roms_export.json"
		m.exportStatus = fmt.Sprintf("Exported %d recent ROMs to %s", len(recent), filename)
	case "statistics":
		// For now, create basic statistics
		favorites := m.browser.GetFavorites()
		recent := m.browser.GetRecentROMs()

		_ = map[string]interface{}{
			"total_favorites":  len(favorites),
			"total_recent":     len(recent),
			"export_date":      "2025-06-16", // TODO: Use current date
			"emulator_version": "1.0.0",      // TODO: Get from version
		}

		filename = "statistics_export.json"
		m.exportStatus = fmt.Sprintf("Exported statistics to %s", filename)

	default:
		m.exportStatus = fmt.Sprintf("Unknown export type: %s", exportType)
		return
	}

}

// getActionForROM returns the appropriate action for a ROM
func getActionForROM(rom ROMInfo) MenuAction {
	if rom.IsDirectory {
		return ActionBrowseROMs
	}
	return ActionLoadROM
}

// GetCurrentScreen returns the current menu screen
func (m *Manager) GetCurrentScreen() MenuScreen {
	return m.currentScreen
}

// SetSearchQuery sets the search query
func (m *Manager) SetSearchQuery(query string) {
	m.searchQuery = query
	if m.currentScreen == ScreenSearch {
		m.LoadSearchScreen()
	}
}

// GetSearchQuery returns the current search query
func (m *Manager) GetSearchQuery() string {
	return m.searchQuery
}

// IsSearchActive returns whether search is active
func (m *Manager) IsSearchActive() bool {
	return m.searchActive
}

// HandleToggleFavorite handles the favorite toggle input when F key is pressed
func (m *Manager) HandleToggleFavorite() {
	// Only allow toggle favorites when viewing ROMs (not in main menu or other screens)
	if m.currentScreen != ScreenBrowser && m.currentScreen != ScreenRecent && m.currentScreen != ScreenFavorites && m.currentScreen != ScreenSearch {
		return
	}

	// Make sure we have menu items and a selected item
	if len(m.menuItems) == 0 || m.selectedItem < 0 || m.selectedItem >= len(m.menuItems) {
		return
	}

	// Get the currently selected menu item
	currentItem := m.menuItems[m.selectedItem]

	// Only toggle favorites for ROM items (not directories or navigation items)
	if currentItem.Action == ActionLoadROM {
		if romPath, ok := currentItem.Data.(string); ok && m.browser != nil {
			romName := extractROMName(romPath)
			m.browser.ToggleFavorite(romPath, romName)
			// Refresh current screen to update favorite indicators
			m.refreshCurrentScreen()
		}
	}
}

// showHelpContent displays help content for a specific topic
func (m *Manager) showHelpContent(topic string) {
	var content string

	switch topic {
	case "navigation":
		content = `Navigation Controls:

↑/↓ or W/S - Navigate menu items
Enter/Space - Select item
ESC/Backspace - Go back
F - Toggle favorite (in ROM lists)

Menu Navigation:
- Use arrow keys or WASD to navigate
- Press Enter to select items
- ESC to go back or exit
- F key to favorite/unfavorite ROMs`

	case "controls":
		content = `Game Controls (CHIP-8 Keypad):

Original CHIP-8 uses a 16-key keypad:
1 2 3 C
4 5 6 D  
7 8 9 E
A 0 B F

Mapped to your keyboard:
1 2 3 4
Q W E R
A S D F
Z X C V

Most games use:
- 2,4,6,8 for directional movement
- 5 for action/select
- 0 for pause (some games)`

	case "features":
		content = `Menu Features:

ROM Browser:
- Browse and organize ROM files
- View metadata and descriptions
- Mark favorites with F key
- Sort by name, date, or size

Recent ROMs:
- Shows recently played games
- Quick access to your games
- Play count tracking

Search:
- Fast text search across ROMs
- Real-time filtering
- Case-insensitive matching

Export:
- Export favorites list
- Export recent ROMs data  
- Export usage statistics`

	case "settings":
		content = `Emulator Settings:

Video:
- Scale: Window size multiplier (1-20x)
- Speed: CPU cycles per second (100-2000)

Audio:
- Volume control
- Frequency settings
- Beep tone configuration

ROM Browser:
- Default ROM directory
- Sort preferences
- Display options

Controls:
- Key mapping (future feature)
- Input sensitivity`

	case "browser":
		content = `ROM Browser Guide:

Navigation:
- Folders show with 📁 icon
- ROMs show with 🎮 icon
- Favorites show with ⭐ icon

Features:
- F key to toggle favorites
- Search with 🔍 Search option
- Refresh with 🔄 Refresh
- Sort options in Settings

File Support:
- .ch8 files (standard CHIP-8)
- .c8 files (alternative extension)
- Automatic metadata extraction
- Text file descriptions (.txt)

Organization:
- Create folders to organize ROMs
- Use favorites for quick access
- Recent list shows last played`

	default:
		content = "Help topic not found."
	}

	m.currentHelp = content
}

// GetCurrentHelp returns the current help content
func (m *Manager) GetCurrentHelp() string {
	return m.currentHelp
}

// ClearCurrentHelp clears the current help content
func (m *Manager) ClearCurrentHelp() {
	m.currentHelp = ""
}

// GetExportStatus returns the current export status
func (m *Manager) GetExportStatus() string {
	return m.exportStatus
}
