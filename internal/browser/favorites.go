package browser

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// FavoriteROM represents a favorited ROM with additional metadata
type FavoriteROM struct {
	Path       string    `json:"path"`
	Name       string    `json:"name"`
	CustomName string    `json:"custom_name,omitempty"` // User-defined name
	Rating     int       `json:"rating,omitempty"`      // 1-5 stars
	Tags       []string  `json:"tags,omitempty"`        // User-defined tags
	Notes      string    `json:"notes,omitempty"`       // User notes
	DateAdded  time.Time `json:"date_added"`
	LastPlayed time.Time `json:"last_played,omitempty"`
}

// FavoritesManager handles ROM favorites system
type FavoritesManager struct {
	favorites map[string]FavoriteROM
	dataFile  string
}

// NewFavoritesManager creates a new favorites manager
func NewFavoritesManager(dataFile string) *FavoritesManager {
	return &FavoritesManager{
		favorites: make(map[string]FavoriteROM),
		dataFile:  dataFile,
	}
}

// AddFavorite adds a ROM to favorites
func (fm *FavoritesManager) AddFavorite(romPath, romName string) error {
	favorite := FavoriteROM{
		Path:      romPath,
		Name:      romName,
		Rating:    0,
		Tags:      make([]string, 0),
		DateAdded: time.Now(),
	}

	fm.favorites[romPath] = favorite
	return fm.SaveToFile()
}

// RemoveFavorite removes a ROM from favorites
func (fm *FavoritesManager) RemoveFavorite(romPath string) error {
	delete(fm.favorites, romPath)
	return fm.SaveToFile()
}

// IsFavorite checks if a ROM is favorited
func (fm *FavoritesManager) IsFavorite(romPath string) bool {
	_, exists := fm.favorites[romPath]
	return exists
}

// GetFavorites returns all favorited ROMs
func (fm *FavoritesManager) GetFavorites() []FavoriteROM {
	favorites := make([]FavoriteROM, 0, len(fm.favorites))
	for _, favorite := range fm.favorites {
		favorites = append(favorites, favorite)
	}
	return favorites
}

// GetFavorite returns a specific favorited ROM
func (fm *FavoritesManager) GetFavorite(romPath string) (FavoriteROM, bool) {
	favorite, exists := fm.favorites[romPath]
	return favorite, exists
}

// SetRating sets the rating for a favorited ROM
func (fm *FavoritesManager) SetRating(romPath string, rating int) error {
	if rating < 0 || rating > 5 {
		return nil // Invalid rating, ignore
	}

	if favorite, exists := fm.favorites[romPath]; exists {
		favorite.Rating = rating
		fm.favorites[romPath] = favorite
		return fm.SaveToFile()
	}

	return nil // ROM not in favorites
}

// GetRating gets the rating for a ROM
func (fm *FavoritesManager) GetRating(romPath string) int {
	if favorite, exists := fm.favorites[romPath]; exists {
		return favorite.Rating
	}
	return 0
}

// SetCustomName sets a custom name for a favorited ROM
func (fm *FavoritesManager) SetCustomName(romPath, customName string) error {
	if favorite, exists := fm.favorites[romPath]; exists {
		favorite.CustomName = customName
		fm.favorites[romPath] = favorite
		return fm.SaveToFile()
	}
	return nil
}

// GetCustomName gets the custom name for a ROM
func (fm *FavoritesManager) GetCustomName(romPath string) string {
	if favorite, exists := fm.favorites[romPath]; exists {
		return favorite.CustomName
	}
	return ""
}

// AddTag adds a tag to a favorited ROM
func (fm *FavoritesManager) AddTag(romPath string, tag string) error {
	if tag == "" {
		return nil
	}

	if favorite, exists := fm.favorites[romPath]; exists {
		// Check if tag already exists
		for _, existingTag := range favorite.Tags {
			if existingTag == tag {
				return nil // Tag already exists
			}
		}

		favorite.Tags = append(favorite.Tags, tag)
		fm.favorites[romPath] = favorite
		return fm.SaveToFile()
	}

	return nil
}

// RemoveTag removes a tag from a favorited ROM
func (fm *FavoritesManager) RemoveTag(romPath string, tag string) error {
	if favorite, exists := fm.favorites[romPath]; exists {
		newTags := make([]string, 0, len(favorite.Tags))
		for _, existingTag := range favorite.Tags {
			if existingTag != tag {
				newTags = append(newTags, existingTag)
			}
		}

		favorite.Tags = newTags
		fm.favorites[romPath] = favorite
		return fm.SaveToFile()
	}

	return nil
}

// GetTags gets all tags for a ROM
func (fm *FavoritesManager) GetTags(romPath string) []string {
	if favorite, exists := fm.favorites[romPath]; exists {
		return favorite.Tags
	}
	return nil
}

// SetNotes sets notes for a favorited ROM
func (fm *FavoritesManager) SetNotes(romPath, notes string) error {
	if favorite, exists := fm.favorites[romPath]; exists {
		favorite.Notes = notes
		fm.favorites[romPath] = favorite
		return fm.SaveToFile()
	}
	return nil
}

// GetNotes gets notes for a ROM
func (fm *FavoritesManager) GetNotes(romPath string) string {
	if favorite, exists := fm.favorites[romPath]; exists {
		return favorite.Notes
	}
	return ""
}

// UpdateLastPlayed updates the last played time for a favorited ROM
func (fm *FavoritesManager) UpdateLastPlayed(romPath string) error {
	if favorite, exists := fm.favorites[romPath]; exists {
		favorite.LastPlayed = time.Now()
		fm.favorites[romPath] = favorite
		return fm.SaveToFile()
	}
	return nil
}

// GetFavoritesByRating returns favorites filtered by minimum rating
func (fm *FavoritesManager) GetFavoritesByRating(minRating int) []FavoriteROM {
	favorites := make([]FavoriteROM, 0)
	for _, favorite := range fm.favorites {
		if favorite.Rating >= minRating {
			favorites = append(favorites, favorite)
		}
	}
	return favorites
}

// GetFavoritesByTag returns favorites that have a specific tag
func (fm *FavoritesManager) GetFavoritesByTag(tag string) []FavoriteROM {
	favorites := make([]FavoriteROM, 0)
	for _, favorite := range fm.favorites {
		for _, favoriteTag := range favorite.Tags {
			if favoriteTag == tag {
				favorites = append(favorites, favorite)
				break
			}
		}
	}
	return favorites
}

// GetAllTags returns all unique tags used across favorites
func (fm *FavoritesManager) GetAllTags() []string {
	tagSet := make(map[string]bool)
	for _, favorite := range fm.favorites {
		for _, tag := range favorite.Tags {
			tagSet[tag] = true
		}
	}

	tags := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tags = append(tags, tag)
	}
	return tags
}

// LoadFromFile loads favorites from file
func (fm *FavoritesManager) LoadFromFile() error {
	if fm.dataFile == "" {
		return nil
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fm.dataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Check if file exists
	if _, err := os.Stat(fm.dataFile); os.IsNotExist(err) {
		// File doesn't exist, start with empty favorites
		return nil
	}

	file, err := os.Open(fm.dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&fm.favorites)
}

// SaveToFile saves favorites to file
func (fm *FavoritesManager) SaveToFile() error {
	if fm.dataFile == "" {
		return nil
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fm.dataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.Create(fm.dataFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(fm.favorites)
}

// CleanupStaleEntries removes entries for ROMs that no longer exist
func (fm *FavoritesManager) CleanupStaleEntries() error {
	var changed bool
	for path := range fm.favorites {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			delete(fm.favorites, path)
			changed = true
		}
	}

	if changed {
		return fm.SaveToFile()
	}

	return nil
}
