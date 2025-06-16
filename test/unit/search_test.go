package unit

import (
	"chip8/internal/menu"
	"fmt"
	"testing"
	"time"
)

func TestSearchManager_BasicSearch(t *testing.T) {
	sm := menu.NewSearchManager()

	// Create test ROM data
	roms := []menu.ROMInfo{
		{Name: "Pong.ch8", Path: "/test/Pong.ch8", IsDirectory: false},
		{Name: "Breakout.ch8", Path: "/test/Breakout.ch8", IsDirectory: false},
		{Name: "Space Invaders.ch8", Path: "/test/Space Invaders.ch8", IsDirectory: false},
		{Name: "Tetris.ch8", Path: "/test/Tetris.ch8", IsDirectory: false},
		{Name: "games", Path: "/test/games", IsDirectory: true},
	}

	// Test basic search
	sm.SetQuery("Pong")
	results := sm.Search(roms, sm.GetDefaultFilter())

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'Pong', got %d", len(results))
	}

	if results[0].Name != "Pong.ch8" {
		t.Errorf("Expected 'Pong.ch8', got '%s'", results[0].Name)
	}
}

func TestSearchManager_CaseInsensitiveSearch(t *testing.T) {
	sm := menu.NewSearchManager()

	roms := []menu.ROMInfo{
		{Name: "Pong.ch8", Path: "/test/Pong.ch8", IsDirectory: false},
		{Name: "BREAKOUT.ch8", Path: "/test/BREAKOUT.ch8", IsDirectory: false},
	}

	// Test case insensitive search
	testCases := []string{"pong", "PONG", "Pong", "PonG"}

	for _, query := range testCases {
		sm.SetQuery(query)
		results := sm.Search(roms, sm.GetDefaultFilter())

		if len(results) != 1 {
			t.Errorf("Case insensitive search failed for '%s': expected 1 result, got %d", query, len(results))
		}
	}
}

func TestSearchManager_PartialMatching(t *testing.T) {
	sm := menu.NewSearchManager()

	roms := []menu.ROMInfo{
		{Name: "Space Invaders.ch8", Path: "/test/Space Invaders.ch8", IsDirectory: false},
		{Name: "Space War.ch8", Path: "/test/Space War.ch8", IsDirectory: false},
		{Name: "Tetris.ch8", Path: "/test/Tetris.ch8", IsDirectory: false},
	}

	// Test partial matching
	sm.SetQuery("Space")
	results := sm.Search(roms, sm.GetDefaultFilter())

	if len(results) != 2 {
		t.Errorf("Expected 2 results for 'Space', got %d", len(results))
	}

	// Test more specific partial match
	sm.SetQuery("Invaders")
	results = sm.Search(roms, sm.GetDefaultFilter())

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'Invaders', got %d", len(results))
	}
}

func TestSearchManager_EmptyQuery(t *testing.T) {
	sm := menu.NewSearchManager()

	roms := []menu.ROMInfo{
		{Name: "Pong.ch8", Path: "/test/Pong.ch8", IsDirectory: false},
		{Name: "Tetris.ch8", Path: "/test/Tetris.ch8", IsDirectory: false},
	}

	// Empty query should return all ROMs
	sm.SetQuery("")
	results := sm.Search(roms, sm.GetDefaultFilter())

	if len(results) != len(roms) {
		t.Errorf("Empty query should return all ROMs: expected %d, got %d", len(roms), len(results))
	}
}

func TestSearchManager_NoMatches(t *testing.T) {
	sm := menu.NewSearchManager()

	roms := []menu.ROMInfo{
		{Name: "Pong.ch8", Path: "/test/Pong.ch8", IsDirectory: false},
		{Name: "Tetris.ch8", Path: "/test/Tetris.ch8", IsDirectory: false},
	}

	// Search for non-existent ROM
	sm.SetQuery("NonExistentGame")
	results := sm.Search(roms, sm.GetDefaultFilter())

	if len(results) != 0 {
		t.Errorf("Expected 0 results for non-existent game, got %d", len(results))
	}
}

func TestSearchManager_PerformanceTest(t *testing.T) {
	sm := menu.NewSearchManager()

	// Create large ROM dataset
	numROMs := 1000
	roms := make([]menu.ROMInfo, numROMs)
	for i := 0; i < numROMs; i++ {
		roms[i] = menu.ROMInfo{
			Name:        fmt.Sprintf("Game%04d.ch8", i),
			Path:        fmt.Sprintf("/test/Game%04d.ch8", i),
			IsDirectory: false,
		}
	}

	// Add target ROM
	roms[500].Name = "TargetGame.ch8"

	// Performance test
	start := time.Now()
	sm.SetQuery("TargetGame")
	results := sm.Search(roms, sm.GetDefaultFilter())
	duration := time.Since(start)

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Performance target: < 50ms for 1000 ROMs
	if duration > 50*time.Millisecond {
		t.Errorf("Search took too long: %v (target: <50ms)", duration)
	}

	t.Logf("Searched %d ROMs in %v", numROMs, duration)
}

func TestSearchManager_FilterDirectories(t *testing.T) {
	sm := menu.NewSearchManager()

	roms := []menu.ROMInfo{
		{Name: "games", Path: "/test/games", IsDirectory: true},
		{Name: "Game.ch8", Path: "/test/Game.ch8", IsDirectory: false},
		{Name: "demos", Path: "/test/demos", IsDirectory: true},
	}

	// Default filter should find both directory and file that match "game"
	sm.SetQuery("game")
	results := sm.Search(roms, sm.GetDefaultFilter())

	if len(results) != 2 {
		t.Errorf("Expected 2 results (games dir + Game.ch8), got %d", len(results))
	}

	// Test with exact filename search
	sm.SetQuery("Game.ch8")
	results = sm.Search(roms, sm.GetDefaultFilter())

	// Should find exactly the .ch8 file
	if len(results) != 1 {
		t.Errorf("Expected 1 result with exact filename search, got %d", len(results))
	}

	if results[0].IsDirectory {
		t.Error("Expected ROM file, got directory")
	}
}

func TestSearchManager_ClearSearch(t *testing.T) {
	sm := menu.NewSearchManager()

	sm.SetQuery("test query")
	if sm.GetQuery() != "test query" {
		t.Error("Query was not set correctly")
	}

	sm.ClearSearch()
	if sm.GetQuery() != "" {
		t.Error("Query was not cleared")
	}
}

func TestSearchManager_MultipleWordSearch(t *testing.T) {
	sm := menu.NewSearchManager()

	roms := []menu.ROMInfo{
		{Name: "Space Invaders.ch8", Path: "/test/Space Invaders.ch8", IsDirectory: false},
		{Name: "Space War.ch8", Path: "/test/Space War.ch8", IsDirectory: false},
		{Name: "Invaders From Mars.ch8", Path: "/test/Invaders From Mars.ch8", IsDirectory: false},
		{Name: "Tetris.ch8", Path: "/test/Tetris.ch8", IsDirectory: false},
	}

	// Test multi-word search
	sm.SetQuery("Space Invaders")
	results := sm.Search(roms, sm.GetDefaultFilter())

	if len(results) != 1 {
		t.Errorf("Expected 1 result for 'Space Invaders', got %d", len(results))
	}

	if results[0].Name != "Space Invaders.ch8" {
		t.Errorf("Expected 'Space Invaders.ch8', got '%s'", results[0].Name)
	}
}

func TestSearchManager_SpecialCharacters(t *testing.T) {
	sm := menu.NewSearchManager()

	roms := []menu.ROMInfo{
		{Name: "Game-1.ch8", Path: "/test/Game-1.ch8", IsDirectory: false},
		{Name: "Game_2.ch8", Path: "/test/Game_2.ch8", IsDirectory: false},
		{Name: "Game (v3).ch8", Path: "/test/Game (v3).ch8", IsDirectory: false},
	}

	// Test searching with special characters
	testCases := map[string]int{
		"Game-1":    1,
		"Game_2":    1,
		"Game (v3)": 1,
		"Game":      3, // Should match all
	}

	for query, expectedCount := range testCases {
		sm.SetQuery(query)
		results := sm.Search(roms, sm.GetDefaultFilter())

		if len(results) != expectedCount {
			t.Errorf("Query '%s': expected %d results, got %d", query, expectedCount, len(results))
		}
	}
}
