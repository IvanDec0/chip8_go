package unit

import (
	"chip8/internal/menu"
	"fmt"
	"strings"
	"testing"
)

func TestManager_NavigateUp(t *testing.T) {
	manager := menu.NewManager()

	// Create test menu items
	items := []menu.MenuItem{
		{Text: "Item 1", Action: menu.ActionBrowseROMs, Data: nil},
		{Text: "Item 2", Action: menu.ActionLoadROM, Data: nil},
		{Text: "Item 3", Action: menu.ActionExit, Data: nil},
	}
	manager.SetMenuItems(items)

	// Initially selected should be 0
	if manager.GetSelectedItem() != 0 {
		t.Errorf("Expected initial selection 0, got %d", manager.GetSelectedItem())
	}

	// Navigate up should wrap to bottom
	manager.NavigateUp()
	if manager.GetSelectedItem() != 2 {
		t.Errorf("Expected selection to wrap to 2, got %d", manager.GetSelectedItem())
	}

	// Navigate up again
	manager.NavigateUp()
	if manager.GetSelectedItem() != 1 {
		t.Errorf("Expected selection 1, got %d", manager.GetSelectedItem())
	}
}

func TestManager_NavigateDown(t *testing.T) {
	manager := menu.NewManager()

	// Create test menu items
	items := []menu.MenuItem{
		{Text: "Item 1", Action: menu.ActionBrowseROMs, Data: nil},
		{Text: "Item 2", Action: menu.ActionLoadROM, Data: nil},
		{Text: "Item 3", Action: menu.ActionExit, Data: nil},
	}
	manager.SetMenuItems(items)

	// Navigate down
	manager.NavigateDown()
	if manager.GetSelectedItem() != 1 {
		t.Errorf("Expected selection 1, got %d", manager.GetSelectedItem())
	}

	// Navigate down again
	manager.NavigateDown()
	if manager.GetSelectedItem() != 2 {
		t.Errorf("Expected selection 2, got %d", manager.GetSelectedItem())
	}

	// Navigate down should wrap to top
	manager.NavigateDown()
	if manager.GetSelectedItem() != 0 {
		t.Errorf("Expected selection to wrap to 0, got %d", manager.GetSelectedItem())
	}
}

func TestManager_SetMenuItems(t *testing.T) {
	manager := menu.NewManager()

	// Test setting menu items
	items := []menu.MenuItem{
		{Text: "Item 1", Action: menu.ActionBrowseROMs, Data: nil},
		{Text: "Item 2", Action: menu.ActionLoadROM, Data: nil},
	}
	manager.SetMenuItems(items)

	retrievedItems := manager.GetMenuItems()
	if len(retrievedItems) != len(items) {
		t.Errorf("Expected %d items, got %d", len(items), len(retrievedItems))
	}

	for i, item := range retrievedItems {
		if item.Text != items[i].Text || item.Action != items[i].Action {
			t.Errorf("Item %d mismatch: expected %+v, got %+v", i, items[i], item)
		}
	}
}

func TestManager_SetSelectedItem(t *testing.T) {
	manager := menu.NewManager()

	// Create test menu items
	items := []menu.MenuItem{
		{Text: "Item 1", Action: menu.ActionBrowseROMs, Data: nil},
		{Text: "Item 2", Action: menu.ActionLoadROM, Data: nil},
		{Text: "Item 3", Action: menu.ActionExit, Data: nil},
	}
	manager.SetMenuItems(items)

	// Test valid selection
	manager.SetSelectedItem(1)
	if manager.GetSelectedItem() != 1 {
		t.Errorf("Expected selection 1, got %d", manager.GetSelectedItem())
	}

	// Test invalid selection (too high)
	manager.SetSelectedItem(10)
	if manager.GetSelectedItem() != 1 {
		t.Errorf("Expected selection to remain 1 for invalid index, got %d", manager.GetSelectedItem())
	}

	// Test invalid selection (negative)
	manager.SetSelectedItem(-1)
	if manager.GetSelectedItem() != 1 {
		t.Errorf("Expected selection to remain 1 for negative index, got %d", manager.GetSelectedItem())
	}
}

func TestManager_ProcessMenuInput(t *testing.T) {
	manager := menu.NewManager()
	stateManager := menu.NewStateManager()
	manager.SetStateManager(stateManager)

	// Create test menu items
	items := []menu.MenuItem{
		{Text: "Item 1", Action: menu.ActionBrowseROMs, Data: nil},
		{Text: "Item 2", Action: menu.ActionLoadROM, Data: "test.ch8"},
		{Text: "Exit", Action: menu.ActionExit, Data: nil},
	}
	manager.SetMenuItems(items)

	// Test navigation up
	manager.ProcessMenuInput(menu.MenuInputEvent{Type: menu.MenuInputUp})
	if manager.GetSelectedItem() != 2 {
		t.Errorf("Expected selection to wrap to 2, got %d", manager.GetSelectedItem())
	}

	// Test navigation down
	manager.ProcessMenuInput(menu.MenuInputEvent{Type: menu.MenuInputDown})
	if manager.GetSelectedItem() != 0 {
		t.Errorf("Expected selection to wrap to 0, got %d", manager.GetSelectedItem())
	}

	// Test exit
	manager.ProcessMenuInput(menu.MenuInputEvent{Type: menu.MenuInputExit})
	if !stateManager.ShouldExit() {
		t.Error("Expected exit flag to be set")
	}
}

func TestManager_GetCurrentScreen(t *testing.T) {
	manager := menu.NewManager()

	// Initially should be main screen
	if manager.GetCurrentScreen() != menu.ScreenMain {
		t.Errorf("Expected initial screen %v, got %v", menu.ScreenMain, manager.GetCurrentScreen())
	}
}

func TestManager_SearchQuery(t *testing.T) {
	manager := menu.NewManager()

	// Test setting search query
	testQuery := "test search"
	manager.SetSearchQuery(testQuery)

	if manager.GetSearchQuery() != testQuery {
		t.Errorf("Expected search query '%s', got '%s'", testQuery, manager.GetSearchQuery())
	}

	// Search is not active until we enter search screen via action
	if manager.IsSearchActive() {
		t.Error("Search should not be active by just setting query")
	}

	// Note: Search activation normally happens through menu navigation
	// and input handling rather than directly calling ProcessMenuInput
	// In a real scenario, this would be done through specific key presses
}

func TestManager_HelpFunctionality(t *testing.T) {
	manager := menu.NewManager()

	// Test help screen loading
	manager.LoadHelpScreen()
	items := manager.GetMenuItems()

	// Check that help topics are loaded
	if len(items) < 5 {
		t.Errorf("Expected at least 5 help items, got %d", len(items))
	}

	// Test navigation to a help topic by setting selection and calling SelectCurrentItem
	// Find the navigation help item
	var navigationIndex = -1
	for i, item := range items {
		if item.Action == menu.ActionShowHelpContent && item.Data == "navigation" {
			navigationIndex = i
			break
		}
	}

	if navigationIndex == -1 {
		t.Error("Could not find navigation help item")
	} else {
		// Select the navigation help item and activate it
		manager.SetSelectedItem(navigationIndex)
		manager.SelectCurrentItem()

		helpContent := manager.GetCurrentHelp()
		if helpContent == "" {
			t.Errorf("Expected help content for 'navigation', got empty string")
		}

		if !strings.Contains(helpContent, "Navigation Controls") {
			t.Errorf("Expected help content to contain 'Navigation Controls', got: %s", helpContent)
		}

		// Test clearing help content
		manager.ClearCurrentHelp()
		if manager.GetCurrentHelp() != "" {
			t.Errorf("Expected empty help content after clearing, got: %s", manager.GetCurrentHelp())
		}
	}
}

func TestManager_ExportFunctionality(t *testing.T) {
	manager := menu.NewManager()

	// Test export screen loading
	manager.LoadExportScreen()
	items := manager.GetMenuItems()

	// Check that export options are loaded
	if len(items) < 3 {
		t.Errorf("Expected at least 3 export items, got %d", len(items))
	}

	// Check that export actions exist
	foundFavorites := false
	foundRecent := false
	foundStats := false

	for _, item := range items {
		switch item.Action {
		case menu.ActionExportFavorites:
			foundFavorites = true
		case menu.ActionExportRecent:
			foundRecent = true
		case menu.ActionExportStatistics:
			foundStats = true
		}
	}

	if !foundFavorites {
		t.Error("Expected to find ActionExportFavorites in export menu")
	}
	if !foundRecent {
		t.Error("Expected to find ActionExportRecent in export menu")
	}
	if !foundStats {
		t.Error("Expected to find ActionExportStatistics in export menu")
	}

	// Test export status (without browser, should get error status)
	// Find the export favorites item and select it
	for i, item := range items {
		if item.Action == menu.ActionExportFavorites {
			manager.SetSelectedItem(i)
			manager.SelectCurrentItem()
			break
		}
	}

	status := manager.GetExportStatus()
	if status == "" {
		t.Error("Expected export status message, got empty string")
	}
}

func TestManager_ScreenTransitions(t *testing.T) {
	manager := menu.NewManager()

	// Start at main menu, navigate to help
	items := manager.GetMenuItems()
	for i, item := range items {
		if item.Action == menu.ActionShowHelp {
			manager.SetSelectedItem(i)
			manager.SelectCurrentItem()
			break
		}
	}

	if manager.GetCurrentScreen() != menu.ScreenHelp {
		t.Errorf("Expected ScreenHelp, got %v", manager.GetCurrentScreen())
	}

	// Go back to main and then navigate to export
	items = manager.GetMenuItems()
	for i, item := range items {
		if item.Action == menu.ActionGoBack {
			manager.SetSelectedItem(i)
			manager.SelectCurrentItem()
			break
		}
	}

	items = manager.GetMenuItems()
	for i, item := range items {
		if item.Action == menu.ActionShowExport {
			manager.SetSelectedItem(i)
			manager.SelectCurrentItem()
			break
		}
	}

	if manager.GetCurrentScreen() != menu.ScreenExport {
		t.Errorf("Expected ScreenExport, got %v", manager.GetCurrentScreen())
	}

	// Test back navigation from export
	items = manager.GetMenuItems()
	for i, item := range items {
		if item.Action == menu.ActionGoBack {
			manager.SetSelectedItem(i)
			manager.SelectCurrentItem()
			break
		}
	}

	if manager.GetCurrentScreen() != menu.ScreenMain {
		t.Errorf("Expected ScreenMain after going back from export, got %v", manager.GetCurrentScreen())
	}
}

func BenchmarkManager_Navigation(b *testing.B) {
	manager := menu.NewManager()

	// Create many menu items
	items := make([]menu.MenuItem, 1000)
	for i := 0; i < 1000; i++ {
		items[i] = menu.MenuItem{
			Text:   fmt.Sprintf("Item %d", i),
			Action: menu.ActionLoadROM,
			Data:   fmt.Sprintf("rom%d.ch8", i),
		}
	}
	manager.SetMenuItems(items)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		manager.NavigateDown()
	}
}
