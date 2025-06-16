package browser

import (
	"chip8/internal/menu"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// FileBrowser implements the Browser interface
type FileBrowser struct {
	currentDir        string
	roms              []menu.ROMInfo
	lastScan          time.Time
	scanCache         map[string][]menu.ROMInfo
	cacheMutex        sync.RWMutex // Protects scanCache and lastScan
	maxCacheAge       time.Duration
	metadataExtractor *MetadataExtractor
	recentManager     *RecentManager
	favoritesManager  *FavoritesManager
	cacheManager      *CacheManager
}

// NewFileBrowser creates a new file browser
func NewFileBrowser(initialDir string) *FileBrowser {
	if initialDir == "" {
		initialDir = "."
	}

	// Resolve to absolute path
	absDir, err := filepath.Abs(initialDir)
	if err != nil {
		absDir = initialDir
	}

	// Initialize data directory for user data
	dataDir := getDataDirectory()

	return &FileBrowser{
		currentDir:        absDir,
		roms:              make([]menu.ROMInfo, 0),
		scanCache:         make(map[string][]menu.ROMInfo),
		maxCacheAge:       time.Minute * 5, // Cache for 5 minutes
		metadataExtractor: NewMetadataExtractor(),
		recentManager:     NewRecentManager(filepath.Join(dataDir, "recent.json"), 20),
		favoritesManager:  NewFavoritesManager(filepath.Join(dataDir, "favorites.json")),
		cacheManager:      NewCacheManager(DefaultCacheConfig()),
	}
}

// ScanDirectory scans a directory for ROM files and subdirectories
func (fb *FileBrowser) ScanDirectory(path string) ([]menu.ROMInfo, error) {
	// Check cache first with read lock
	fb.cacheMutex.RLock()
	if cached, exists := fb.scanCache[path]; exists {
		if time.Since(fb.lastScan) < fb.maxCacheAge {
			fb.cacheMutex.RUnlock()
			return cached, nil
		}
	}
	fb.cacheMutex.RUnlock()

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, &BrowserError{
			Type:    ErrorFileAccess,
			Message: "Failed to read directory",
			Cause:   err,
		}
	}

	var roms []menu.ROMInfo

	// Always add parent directory option unless we're at root
	if path != "/" && path != "." {
		parent := filepath.Dir(path)
		roms = append(roms, menu.ROMInfo{
			Path:        parent,
			Name:        "..",
			Size:        0,
			ModTime:     time.Now(),
			IsDirectory: true,
		})
	}

	for _, entry := range entries {
		fullPath := filepath.Join(path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue // Skip files we can't read
		}

		// Include directories and .ch8 files
		if info.IsDir() || isROMFile(entry.Name()) {
			roms = append(roms, menu.ROMInfo{
				Path:        fullPath,
				Name:        entry.Name(),
				Size:        info.Size(),
				ModTime:     info.ModTime(),
				IsDirectory: info.IsDir(),
			})
		}
	}

	// Sort: directories first, then files, both alphabetically
	sort.Slice(roms, func(i, j int) bool {
		// Special case: .. should always be first
		if roms[i].Name == ".." {
			return true
		}
		if roms[j].Name == ".." {
			return false
		}

		// Directories come before files
		if roms[i].IsDirectory != roms[j].IsDirectory {
			return roms[i].IsDirectory
		}

		// Alphabetical within the same type
		return strings.ToLower(roms[i].Name) < strings.ToLower(roms[j].Name)
	})

	// Cache the results with write lock
	fb.cacheMutex.Lock()
	fb.scanCache[path] = roms
	fb.lastScan = time.Now()
	fb.cacheMutex.Unlock()

	return roms, nil
}

// GetROMs returns the currently loaded ROM list
func (fb *FileBrowser) GetROMs() []menu.ROMInfo {
	return fb.roms
}

// SetCurrentDirectory sets the current directory and scans it
func (fb *FileBrowser) SetCurrentDirectory(path string) error {
	// Resolve to absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return &BrowserError{
			Type:    ErrorInvalidDirectory,
			Message: "Failed to resolve directory path",
			Cause:   err,
		}
	}

	// Check if directory exists and is accessible
	info, err := os.Stat(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &BrowserError{
				Type:    ErrorNotFound,
				Message: "Directory does not exist",
				Cause:   err,
			}
		}
		return &BrowserError{
			Type:    ErrorFileAccess,
			Message: "Cannot access directory",
			Cause:   err,
		}
	}

	if !info.IsDir() {
		return &BrowserError{
			Type:    ErrorInvalidDirectory,
			Message: "Path is not a directory",
			Cause:   nil,
		}
	}

	fb.currentDir = absPath

	// Scan the new directory
	roms, err := fb.ScanDirectory(absPath)
	if err != nil {
		return err
	}

	fb.roms = roms
	return nil
}

// GetCurrentDirectory returns the current directory path
func (fb *FileBrowser) GetCurrentDirectory() string {
	return fb.currentDir
}

// NavigateUp navigates to the parent directory
func (fb *FileBrowser) NavigateUp() error {
	parent := filepath.Dir(fb.currentDir)
	if parent == fb.currentDir {
		// Already at root
		return nil
	}
	return fb.SetCurrentDirectory(parent)
}

// Refresh rescans the current directory, clearing cache
func (fb *FileBrowser) Refresh() error {
	// Clear cache for current directory with write lock
	fb.cacheMutex.Lock()
	delete(fb.scanCache, fb.currentDir)
	fb.cacheMutex.Unlock()
	return fb.SetCurrentDirectory(fb.currentDir)
}

// isROMFile checks if a filename appears to be a ROM file
func isROMFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".ch8" || ext == ".c8"
}

// ValidateROM validates a ROM file before loading
func (fb *FileBrowser) ValidateROM(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &BrowserError{
				Type:    ErrorNotFound,
				Message: "ROM file not found",
				Cause:   err,
			}
		}
		return &BrowserError{
			Type:    ErrorFileAccess,
			Message: "Cannot access ROM file",
			Cause:   err,
		}
	}

	// Check if it's a file (not a directory)
	if info.IsDir() {
		return &BrowserError{
			Type:    ErrorInvalidDirectory,
			Message: "Path is a directory, not a ROM file",
			Cause:   nil,
		}
	}

	// Check file size (typical ROM is 512-4096 bytes, but allow up to 64KB for larger ROMs)
	size := info.Size()
	if size < 1 {
		return &BrowserError{
			Type:    ErrorFileAccess,
			Message: "ROM file is empty",
			Cause:   nil,
		}
	}
	if size > 65536 { // 64KB limit
		return &BrowserError{
			Type:    ErrorFileAccess,
			Message: "ROM file is too large (max 64KB)",
			Cause:   nil,
		}
	}

	// Verify file extension
	if !isROMFile(info.Name()) {
		return &BrowserError{
			Type:    ErrorFileAccess,
			Message: "File is not a recognized ROM format (.ch8 or .c8)",
			Cause:   nil,
		}
	}

	return nil
}

// ScanDirectoryWithFallback scans directory with fallback options
func (fb *FileBrowser) ScanDirectoryWithFallback(path string) ([]menu.ROMInfo, error) {
	// Try primary directory
	roms, err := fb.ScanDirectory(path)
	if err == nil {
		return roms, nil
	}

	// Try executable directory as fallback
	execDir, execErr := os.Executable()
	if execErr == nil {
		execDir = filepath.Dir(execDir)
		romsDir := filepath.Join(execDir, "roms")
		if roms, fallbackErr := fb.ScanDirectory(romsDir); fallbackErr == nil {
			return roms, nil
		}
	}

	// Try current directory as last resort
	if currentRoms, currentErr := fb.ScanDirectory("."); currentErr == nil {
		return currentRoms, nil
	}

	// Return original error if all fallbacks fail
	return nil, &BrowserError{
		Type:    ErrorNotFound,
		Message: "No ROM directories found. Tried: " + path,
		Cause:   err,
	}
}

// getDataDirectory returns the application data directory
func getDataDirectory() string {
	//home := os.Getenv("HOME")
	home := "."
	if home == "" {
		return "."
	}
	dataDir := filepath.Join(home, ".chip8")
	os.MkdirAll(dataDir, 0755)
	return dataDir
}

// GetMetadataExtractor returns the metadata extractor
func (fb *FileBrowser) GetMetadataExtractor() *MetadataExtractor {
	return fb.metadataExtractor
}

// GetRecentManager returns the recent manager
func (fb *FileBrowser) GetRecentManager() *RecentManager {
	return fb.recentManager
}

// GetFavoritesManager returns the favorites manager
func (fb *FileBrowser) GetFavoritesManager() *FavoritesManager {
	return fb.favoritesManager
}

// GetCacheManager returns the cache manager
func (fb *FileBrowser) GetCacheManager() *CacheManager {
	return fb.cacheManager
}

// ScanDirectoryWithMetadata scans directory and includes metadata
func (fb *FileBrowser) ScanDirectoryWithMetadata(path string) ([]menu.ROMInfo, error) {
	roms, err := fb.ScanDirectory(path)
	if err != nil {
		return nil, err
	}

	// Preload metadata for ROM files in background
	go fb.preloadMetadata(roms)

	return roms, nil
}

// preloadMetadata preloads metadata for ROMs in background
func (fb *FileBrowser) preloadMetadata(roms []menu.ROMInfo) {
	for _, rom := range roms {
		if !rom.IsDirectory && isROMFile(rom.Name) {
			// Extract metadata (will be cached)
			fb.metadataExtractor.ExtractMetadata(rom.Path)
		}
	}
}

// GetRecentROMs returns recent ROMs in interface format
func (fb *FileBrowser) GetRecentROMs() []menu.RecentROM {
	recent := fb.recentManager.GetRecentROMs()
	result := make([]menu.RecentROM, len(recent))

	for i, rom := range recent {
		result[i] = menu.RecentROM{
			Path:       rom.Path,
			Name:       rom.Name,
			LastPlayed: rom.LastPlayed.Format("2006-01-02 15:04"),
			PlayCount:  rom.PlayCount,
			IsFavorite: rom.IsFavorite,
		}
	}

	return result
}

// GetFavorites returns favorites in interface format
func (fb *FileBrowser) GetFavorites() []menu.FavoriteROM {
	favorites := fb.favoritesManager.GetFavorites()
	result := make([]menu.FavoriteROM, len(favorites))

	for i, fav := range favorites {
		result[i] = menu.FavoriteROM{
			Path:       fav.Path,
			Name:       fav.Name,
			CustomName: fav.CustomName,
			Rating:     fav.Rating,
			Tags:       fav.Tags,
			DateAdded:  fav.DateAdded.Format("2006-01-02"),
		}
	}

	return result
}

// IsFavorite checks if a ROM is favorited
func (fb *FileBrowser) IsFavorite(romPath string) bool {
	return fb.favoritesManager.IsFavorite(romPath)
}

// GetROMMetadata returns metadata in interface format
func (fb *FileBrowser) GetROMMetadata(romPath string) (*menu.ROMMetadata, error) {
	metadata, err := fb.metadataExtractor.ExtractMetadata(romPath)
	if err != nil {
		return nil, err
	}

	return &menu.ROMMetadata{
		Title:       metadata.Title,
		Author:      metadata.Author,
		Description: metadata.Description,
		Controls:    metadata.Controls,
		Year:        metadata.Year,
		System:      metadata.System,
	}, nil
}

// AddToRecent adds a ROM to recent list
func (fb *FileBrowser) AddToRecent(romPath, romName string) error {
	return fb.recentManager.AddRecentROM(romPath, romName)
}

// ToggleFavorite toggles favorite status for a ROM
func (fb *FileBrowser) ToggleFavorite(romPath, romName string) error {
	if fb.favoritesManager.IsFavorite(romPath) {
		return fb.favoritesManager.RemoveFavorite(romPath)
	} else {
		return fb.favoritesManager.AddFavorite(romPath, romName)
	}
}

// LoadUserData loads recent ROMs and favorites from disk
func (fb *FileBrowser) LoadUserData() error {
	// Load recent ROMs
	if err := fb.recentManager.LoadFromFile(); err != nil {
		// Non-fatal error, just log it
		fmt.Printf("Warning: Failed to load recent ROMs: %v\n", err)
	}

	// Load favorites
	if err := fb.favoritesManager.LoadFromFile(); err != nil {
		// Non-fatal error, just log it
		fmt.Printf("Warning: Failed to load favorites: %v\n", err)
	}

	// Start cache cleanup
	fb.cacheManager.StartCleanup()

	return nil
}

// Close cleans up resources
func (fb *FileBrowser) Close() {
	if fb.cacheManager != nil {
		fb.cacheManager.StopCleanup()
	}
}
