package menu

import "fmt"

// Manager coordinates menu operations
type Manager struct {
	stateManager *StateManager
	theme        MenuTheme
	selectedItem int
	menuItems    []MenuItem
	browser      ROMBrowser // Add browser reference
}

// NewManager creates a new menu manager
func NewManager() *Manager {
	return &Manager{
		stateManager: nil, // Will be set externally
		theme:        DefaultMenuTheme(),
		selectedItem: 0,
		menuItems:    createMainMenuItems(),
	}
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
					m.LoadROMBrowser() // Reload the browser with new directory
				} else {
					fmt.Printf("Error setting directory: %v\n", err)
				}
			}
		} else {
			// Initial ROM browser load
			err := m.LoadROMBrowser()
			if err != nil {
				fmt.Printf("Error loading ROM browser: %v\n", err)
			}
		}
	case ActionLoadROM:
		if romPath, ok := data.(string); ok {
			if m.stateManager != nil {
				m.stateManager.SetSelectedROM(romPath)
				m.stateManager.SetShouldLoadROM(true)
				m.stateManager.TransitionTo(StateEmulator)
			} else {
				fmt.Println("Error: StateManager is nil!")
			}
		}
	case ActionExit:
		if m.stateManager != nil {
			m.stateManager.SetShouldExit(true)
		} else {
			fmt.Println("Error: StateManager is nil!")
		}
	case ActionGoBack:
		// Go back to main menu
		m.menuItems = createMainMenuItems()
		m.selectedItem = 0
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
	// Scan current directory
	roms, err := m.browser.ScanDirectory(currentDir)
	if err != nil {
		fmt.Printf("Error scanning directory: %v\n", err)
		return err
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
		if rom.IsDirectory {
			items = append(items, MenuItem{
				Text:   "📁 " + rom.Name,
				Action: ActionBrowseROMs,
				Data:   rom.Path,
			})
		} else {
			items = append(items, MenuItem{
				Text:   "🎮 " + rom.Name,
				Action: ActionLoadROM,
				Data:   rom.Path,
			})
		}
	}

	m.SetMenuItems(items)
	return nil
}
