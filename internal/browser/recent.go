package browser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// RecentROM represents a recently accessed ROM
type RecentROM struct {
	Path       string    `json:"path"`
	Name       string    `json:"name"`
	LastPlayed time.Time `json:"last_played"`
	PlayCount  int       `json:"play_count"`
	IsFavorite bool      `json:"is_favorite"`
}

// RecentManager handles recent ROMs tracking
type RecentManager struct {
	recentROMs []RecentROM
	maxRecent  int
	dataFile   string
}

// NewRecentManager creates a new recent manager
func NewRecentManager(dataFile string, maxRecent int) *RecentManager {
	if maxRecent <= 0 {
		maxRecent = 10 // Default to 10 recent ROMs
	}

	return &RecentManager{
		recentROMs: make([]RecentROM, 0, maxRecent),
		maxRecent:  maxRecent,
		dataFile:   dataFile,
	}
}

// AddRecentROM adds a ROM to the recent list
func (rm *RecentManager) AddRecentROM(romPath, romName string) error {
	// Check if ROM already exists in recent list
	for i, rom := range rm.recentROMs {
		if rom.Path == romPath {
			// Update existing entry
			rm.recentROMs[i].LastPlayed = time.Now()
			rm.recentROMs[i].PlayCount++

			// Move to front if not already there
			if i != 0 {
				recent := rm.recentROMs[i]
				rm.recentROMs = append(rm.recentROMs[:i], rm.recentROMs[i+1:]...)
				rm.recentROMs = append([]RecentROM{recent}, rm.recentROMs...)
			}

			return rm.SaveToFile()
		}
	}

	// Add new ROM to front of list
	newROM := RecentROM{
		Path:       romPath,
		Name:       romName,
		LastPlayed: time.Now(),
		PlayCount:  1,
		IsFavorite: false,
	}

	rm.recentROMs = append([]RecentROM{newROM}, rm.recentROMs...)

	// Trim list if too long
	if len(rm.recentROMs) > rm.maxRecent {
		rm.recentROMs = rm.recentROMs[:rm.maxRecent]
	}

	return rm.SaveToFile()
}

// GetRecentROMs returns the list of recent ROMs
func (rm *RecentManager) GetRecentROMs() []RecentROM {
	return rm.recentROMs
}

// GetFavorites returns ROMs marked as favorites
func (rm *RecentManager) GetFavorites() []RecentROM {
	var favorites []RecentROM
	for _, rom := range rm.recentROMs {
		if rom.IsFavorite {
			favorites = append(favorites, rom)
		}
	}
	return favorites
}

// ToggleFavorite toggles the favorite status of a ROM
func (rm *RecentManager) ToggleFavorite(romPath string) error {
	for i, rom := range rm.recentROMs {
		if rom.Path == romPath {
			rm.recentROMs[i].IsFavorite = !rm.recentROMs[i].IsFavorite
			return rm.SaveToFile()
		}
	}

	// ROM not in recent list, we can't favorite it
	return nil
}

// IsFavorite checks if a ROM is marked as favorite
func (rm *RecentManager) IsFavorite(romPath string) bool {
	for _, rom := range rm.recentROMs {
		if rom.Path == romPath {
			return rom.IsFavorite
		}
	}
	return false
}

// RemoveROM removes a ROM from recent list
func (rm *RecentManager) RemoveROM(romPath string) error {
	for i, rom := range rm.recentROMs {
		if rom.Path == romPath {
			rm.recentROMs = append(rm.recentROMs[:i], rm.recentROMs[i+1:]...)
			return rm.SaveToFile()
		}
	}
	return nil
}

// ClearRecent clears all recent ROMs
func (rm *RecentManager) ClearRecent() error {
	rm.recentROMs = rm.recentROMs[:0]
	return rm.SaveToFile()
}

// LoadFromFile loads recent ROMs from file
func (rm *RecentManager) LoadFromFile() error {
	if rm.dataFile == "" {
		return nil
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(rm.dataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(rm.dataFile); os.IsNotExist(err) {
		// File doesn't exist, start with empty list
		return nil
	}

	file, err := os.Open(rm.dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&rm.recentROMs)
}

// SaveToFile saves recent ROMs to file
func (rm *RecentManager) SaveToFile() error {
	if rm.dataFile == "" {
		return nil
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(rm.dataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(rm.dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(rm.recentROMs)
}

// GetMostPlayed returns ROMs sorted by play count
func (rm *RecentManager) GetMostPlayed(limit int) []RecentROM {
	if limit <= 0 || limit > len(rm.recentROMs) {
		limit = len(rm.recentROMs)
	}

	// Create a copy and sort by play count
	roms := make([]RecentROM, len(rm.recentROMs))
	copy(roms, rm.recentROMs)

	// Simple bubble sort by play count (descending)
	for i := 0; i < len(roms)-1; i++ {
		for j := 0; j < len(roms)-i-1; j++ {
			if roms[j].PlayCount < roms[j+1].PlayCount {
				roms[j], roms[j+1] = roms[j+1], roms[j]
			}
		}
	}

	if limit < len(roms) {
		return roms[:limit]
	}
	return roms
}

// CleanupStaleEntries removes entries for ROMs that no longer exist
func (rm *RecentManager) CleanupStaleEntries() error {
	var validROMs []RecentROM

	for _, rom := range rm.recentROMs {
		if _, err := os.Stat(rom.Path); err == nil {
			validROMs = append(validROMs, rom)
		}
	}

	if len(validROMs) != len(rm.recentROMs) {
		rm.recentROMs = validROMs
		return rm.SaveToFile()
	}

	return nil
}
