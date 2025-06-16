package unit

import (
	"chip8/internal/menu"
	"fmt"
	"testing"
	"time"
)

func TestSortManager_DefaultSort(t *testing.T) {
	sm := menu.NewSortManager()

	// Test default sort criteria and direction
	criteria, ascending := sm.GetCurrentSort()
	if criteria.String() == "" {
		t.Error("Default sort criteria should not be empty")
	}

	if !ascending {
		t.Error("Default sort should be ascending")
	}
}

func TestSortManager_SortByName(t *testing.T) {
	sm := menu.NewSortManager()

	roms := []menu.ROMInfo{
		{Name: "Zetris.ch8", Path: "/test/Zetris.ch8", IsDirectory: false},
		{Name: "Breakout.ch8", Path: "/test/Breakout.ch8", IsDirectory: false},
		{Name: "Asteroids.ch8", Path: "/test/Asteroids.ch8", IsDirectory: false},
		{Name: "Pong.ch8", Path: "/test/Pong.ch8", IsDirectory: false},
	}

	sorted := sm.SortROMs(roms)

	// Check if sorted alphabetically
	expectedOrder := []string{"Asteroids.ch8", "Breakout.ch8", "Pong.ch8", "Zetris.ch8"}
	for i, rom := range sorted {
		if rom.Name != expectedOrder[i] {
			t.Errorf("Sort by name failed: expected %s at position %d, got %s",
				expectedOrder[i], i, rom.Name)
		}
	}
}

func TestSortManager_SortByDate(t *testing.T) {
	sm := menu.NewSortManager()

	// Create ROMs with different modification times
	now := time.Now()
	roms := []menu.ROMInfo{
		{Name: "New.ch8", Path: "/test/New.ch8", ModTime: now, IsDirectory: false},
		{Name: "Old.ch8", Path: "/test/Old.ch8", ModTime: now.Add(-24 * time.Hour), IsDirectory: false},
		{Name: "Medium.ch8", Path: "/test/Medium.ch8", ModTime: now.Add(-12 * time.Hour), IsDirectory: false},
	}

	// Set sort to date and ensure ascending order
	sm.SetSort(menu.SortByDate, true) // Ascending (oldest first)

	sorted := sm.SortROMs(roms)

	// Check if sorted by date (oldest first)
	if sorted[0].Name != "Old.ch8" {
		t.Errorf("Expected oldest ROM first, got %s", sorted[0].Name)
	}

	if sorted[2].Name != "New.ch8" {
		t.Errorf("Expected newest ROM last, got %s", sorted[2].Name)
	}
}

func TestSortManager_SortBySize(t *testing.T) {
	sm := menu.NewSortManager()

	roms := []menu.ROMInfo{
		{Name: "Big.ch8", Path: "/test/Big.ch8", Size: 4096, IsDirectory: false},
		{Name: "Small.ch8", Path: "/test/Small.ch8", Size: 512, IsDirectory: false},
		{Name: "Medium.ch8", Path: "/test/Medium.ch8", Size: 2048, IsDirectory: false},
	}

	// Set sort to size
	sm.SetSort(menu.SortBySize, true) // Ascending (smallest first)

	sorted := sm.SortROMs(roms)

	// Check if sorted by size (smallest first)
	if sorted[0].Name != "Small.ch8" {
		t.Errorf("Expected smallest ROM first, got %s", sorted[0].Name)
	}

	if sorted[2].Name != "Big.ch8" {
		t.Errorf("Expected largest ROM last, got %s", sorted[2].Name)
	}
}

func TestSortManager_DirectoriesFirst(t *testing.T) {
	sm := menu.NewSortManager()

	roms := []menu.ROMInfo{
		{Name: "Game.ch8", Path: "/test/Game.ch8", IsDirectory: false},
		{Name: "folder", Path: "/test/folder", IsDirectory: true},
		{Name: "Another.ch8", Path: "/test/Another.ch8", IsDirectory: false},
	}

	sorted := sm.SortROMs(roms)

	// Directories should come first
	if !sorted[0].IsDirectory {
		t.Error("Expected directory to be sorted first")
	}

	if sorted[0].Name != "folder" {
		t.Errorf("Expected 'folder' to be first, got %s", sorted[0].Name)
	}
}

func TestSortManager_ToggleDirection(t *testing.T) {
	sm := menu.NewSortManager()

	roms := []menu.ROMInfo{
		{Name: "B.ch8", Path: "/test/B.ch8", IsDirectory: false},
		{Name: "A.ch8", Path: "/test/A.ch8", IsDirectory: false},
		{Name: "C.ch8", Path: "/test/C.ch8", IsDirectory: false},
	}

	// Sort ascending
	sm.SetSort(menu.SortByName, true)
	sorted := sm.SortROMs(roms)

	if sorted[0].Name != "A.ch8" {
		t.Errorf("Ascending sort failed: expected A.ch8 first, got %s", sorted[0].Name)
	}

	// Toggle to descending
	sm.ToggleDirection()
	sorted = sm.SortROMs(roms)

	if sorted[0].Name != "C.ch8" {
		t.Errorf("Descending sort failed: expected C.ch8 first, got %s", sorted[0].Name)
	}
}

func TestSortManager_NextSort(t *testing.T) {
	sm := menu.NewSortManager()

	initialCriteria, _ := sm.GetCurrentSort()

	// Move to next sort criteria
	sm.NextSort()
	newCriteria, _ := sm.GetCurrentSort()

	if newCriteria == initialCriteria {
		t.Error("NextSort should change the sort criteria")
	}
}

func TestSortManager_StableSorting(t *testing.T) {
	sm := menu.NewSortManager()

	// Create ROMs with same names to test stable sorting
	roms := []menu.ROMInfo{
		{Name: "Game.ch8", Path: "/test/dir1/Game.ch8", IsDirectory: false},
		{Name: "Game.ch8", Path: "/test/dir2/Game.ch8", IsDirectory: false},
		{Name: "Game.ch8", Path: "/test/dir3/Game.ch8", IsDirectory: false},
	}

	// Sort multiple times to ensure stability
	sorted1 := sm.SortROMs(roms)
	sorted2 := sm.SortROMs(sorted1)

	// Results should be identical
	for i := range sorted1 {
		if sorted1[i].Path != sorted2[i].Path {
			t.Error("Sorting is not stable")
			break
		}
	}
}

func TestSortManager_PerformanceTest(t *testing.T) {
	sm := menu.NewSortManager()

	// Create large dataset
	numROMs := 1000
	roms := make([]menu.ROMInfo, numROMs)
	for i := 0; i < numROMs; i++ {
		roms[i] = menu.ROMInfo{
			Name:        fmt.Sprintf("ROM%04d.ch8", numROMs-i), // Reverse order
			Path:        fmt.Sprintf("/test/ROM%04d.ch8", i),
			IsDirectory: false,
		}
	}

	start := time.Now()
	sorted := sm.SortROMs(roms)
	duration := time.Since(start)

	if len(sorted) != numROMs {
		t.Errorf("Expected %d ROMs after sorting, got %d", numROMs, len(sorted))
	}

	// Performance target: < 100ms for 1000 ROMs
	if duration > 100*time.Millisecond {
		t.Errorf("Sorting took too long: %v (target: <100ms)", duration)
	}

	t.Logf("Sorted %d ROMs in %v", numROMs, duration)
}

func TestSortManager_EmptyList(t *testing.T) {
	sm := menu.NewSortManager()

	var roms []menu.ROMInfo
	sorted := sm.SortROMs(roms)

	if len(sorted) != 0 {
		t.Error("Sorting empty list should return empty list")
	}
}

func TestSortManager_SingleItem(t *testing.T) {
	sm := menu.NewSortManager()

	roms := []menu.ROMInfo{
		{Name: "Only.ch8", Path: "/test/Only.ch8", IsDirectory: false},
	}

	sorted := sm.SortROMs(roms)

	if len(sorted) != 1 {
		t.Error("Sorting single item should return single item")
	}

	if sorted[0].Name != "Only.ch8" {
		t.Error("Single item was modified during sorting")
	}
}

func TestSortManager_MixedTypes(t *testing.T) {
	sm := menu.NewSortManager()

	roms := []menu.ROMInfo{
		{Name: "zzz.ch8", Path: "/test/zzz.ch8", IsDirectory: false},
		{Name: "aaa-folder", Path: "/test/aaa-folder", IsDirectory: true},
		{Name: "mmm.ch8", Path: "/test/mmm.ch8", IsDirectory: false},
		{Name: "zzz-folder", Path: "/test/zzz-folder", IsDirectory: true},
	}

	sorted := sm.SortROMs(roms)

	// First two should be directories (sorted alphabetically)
	if !sorted[0].IsDirectory || !sorted[1].IsDirectory {
		t.Error("Directories should be sorted first")
	}

	if sorted[0].Name != "aaa-folder" {
		t.Errorf("Expected 'aaa-folder' first among directories, got %s", sorted[0].Name)
	}

	// Last two should be files (sorted alphabetically)
	if sorted[2].IsDirectory || sorted[3].IsDirectory {
		t.Error("Files should be sorted after directories")
	}

	if sorted[2].Name != "mmm.ch8" {
		t.Errorf("Expected 'mmm.ch8' first among files, got %s", sorted[2].Name)
	}
}
