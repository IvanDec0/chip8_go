package integration

import (
	"chip8/internal/browser"
	"chip8/internal/config"
	"chip8/internal/menu"
	"os"
	"path/filepath"
	"testing"
)

func TestBasicIntegration(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "chip8_integration_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test ROM files
	testFiles := []string{
		"test1.ch8",
		"test2.c8",
		"readme.txt",
	}

	for _, file := range testFiles {
		filePath := filepath.Join(tmpDir, file)
		if err := os.WriteFile(filePath, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Test browser functionality
	browser := browser.NewFileBrowser(tmpDir)
	roms, err := browser.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Failed to scan directory: %v", err)
	}

	// Verify expected ROM files found (+ parent directory)
	romFileCount := 0
	for _, rom := range roms {
		if !rom.IsDirectory && (filepath.Ext(rom.Name) == ".ch8" || filepath.Ext(rom.Name) == ".c8") {
			romFileCount++
		}
	}

	if romFileCount != 2 {
		t.Errorf("Expected 2 ROM files, found %d", romFileCount)
	}

	// Test menu manager
	menuManager := menu.NewManager()
	menuItems := make([]menu.MenuItem, 0, len(roms))
	for _, rom := range roms {
		action := menu.ActionBrowseROMs
		if !rom.IsDirectory {
			action = menu.ActionLoadROM
		}
		menuItems = append(menuItems, menu.MenuItem{
			Text:   rom.Name,
			Action: action,
			Data:   rom.Path,
		})
	}
	menuManager.SetMenuItems(menuItems)

	// Test search functionality
	searchManager := menu.NewSearchManager()
	searchManager.SetQuery("test")
	filter := searchManager.GetDefaultFilter()
	results := searchManager.Search(roms, filter)

	if len(results) < 2 {
		t.Errorf("Expected at least 2 search results, got %d", len(results))
	}

	// Test configuration
	configPath := filepath.Join(tmpDir, "config.json")
	configManager := config.NewConfigManager(configPath)

	// Test default config
	userConfig := configManager.GetUserConfig()
	userConfig.MenuConfig.DefaultROMDirectory = tmpDir
	err = configManager.SetUserConfig(userConfig)
	if err != nil {
		t.Fatalf("Failed to set user config: %v", err)
	}

	err = configManager.SaveConfig()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Test config reload
	configManager2 := config.NewConfigManager(configPath)
	err = configManager2.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to reload config: %v", err)
	}

	reloadedConfig := configManager2.GetUserConfig()
	if reloadedConfig.MenuConfig.DefaultROMDirectory != tmpDir {
		t.Errorf("Configuration not persisted correctly")
	}
}

func TestBrowserNavigationIntegration(t *testing.T) {
	// Create temporary test directory structure
	tmpDir, err := os.MkdirTemp("", "chip8_nav_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create subdirectory with ROM files
	gamesDir := filepath.Join(tmpDir, "games")
	if err := os.Mkdir(gamesDir, 0755); err != nil {
		t.Fatalf("Failed to create games directory: %v", err)
	}

	// Create test files
	testFiles := map[string]string{
		filepath.Join(tmpDir, "main.ch8"):     "main rom",
		filepath.Join(gamesDir, "game1.ch8"):  "game 1",
		filepath.Join(gamesDir, "game2.c8"):   "game 2",
		filepath.Join(gamesDir, "readme.txt"): "readme",
	}

	for filePath, content := range testFiles {
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filePath, err)
		}
	}

	// Test browser navigation
	browser := browser.NewFileBrowser(tmpDir)

	// Scan root directory
	roms, err := browser.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Failed to scan root directory: %v", err)
	}

	// Should find main.ch8 and games directory
	hasMainROM := false
	hasGamesDir := false
	for _, rom := range roms {
		if rom.Name == "main.ch8" && !rom.IsDirectory {
			hasMainROM = true
		}
		if rom.Name == "games" && rom.IsDirectory {
			hasGamesDir = true
		}
	}

	if !hasMainROM {
		t.Error("main.ch8 not found in root directory")
	}
	if !hasGamesDir {
		t.Error("games directory not found in root directory")
	}

	// Navigate into games directory
	err = browser.SetCurrentDirectory(gamesDir)
	if err != nil {
		t.Fatalf("Failed to navigate to games directory: %v", err)
	}

	gamesRoms, err := browser.ScanDirectory(gamesDir)
	if err != nil {
		t.Fatalf("Failed to scan games directory: %v", err)
	}

	// Count ROM files in games directory
	gameROMCount := 0
	for _, rom := range gamesRoms {
		if !rom.IsDirectory && (filepath.Ext(rom.Name) == ".ch8" || filepath.Ext(rom.Name) == ".c8") {
			gameROMCount++
		}
	}

	if gameROMCount != 2 {
		t.Errorf("Expected 2 ROM files in games directory, found %d", gameROMCount)
	}

	// Test navigation back to parent
	err = browser.NavigateUp()
	if err != nil {
		t.Fatalf("Failed to navigate back to parent: %v", err)
	}

	if browser.GetCurrentDirectory() != tmpDir {
		t.Errorf("Expected to be back in root directory")
	}
}

func TestMetadataIntegration(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "chip8_metadata_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create ROM file and accompanying text file
	romPath := filepath.Join(tmpDir, "testgame.ch8")
	textPath := filepath.Join(tmpDir, "testgame.txt")

	romContent := []byte("fake rom content")
	textContent := []byte(`Test Game
Author: Test Developer

This is a test game for integration testing.
`)

	if err := os.WriteFile(romPath, romContent, 0644); err != nil {
		t.Fatalf("Failed to create ROM file: %v", err)
	}
	if err := os.WriteFile(textPath, textContent, 0644); err != nil {
		t.Fatalf("Failed to create text file: %v", err)
	}

	// Test metadata extraction
	extractor := browser.NewMetadataExtractor()
	metadata, err := extractor.ExtractMetadata(romPath)
	if err != nil {
		t.Fatalf("Failed to extract metadata: %v", err)
	}

	if metadata.Title != "Test Game" {
		t.Errorf("Expected title 'Test Game', got '%s'", metadata.Title)
	}
	if metadata.Author != "Test Developer" {
		t.Errorf("Expected author 'Test Developer', got '%s'", metadata.Author)
	}

	// Test caching
	metadata2, err := extractor.ExtractMetadata(romPath)
	if err != nil {
		t.Fatalf("Failed to extract cached metadata: %v", err)
	}

	if metadata.Title != metadata2.Title {
		t.Error("Metadata caching not working correctly")
	}
}
