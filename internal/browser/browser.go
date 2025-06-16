package browser

import (
	"chip8/internal/menu"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// FileBrowser implements the Browser interface
type FileBrowser struct {
	currentDir  string
	roms        []menu.ROMInfo
	lastScan    time.Time
	scanCache   map[string][]menu.ROMInfo
	maxCacheAge time.Duration
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

	return &FileBrowser{
		currentDir:  absDir,
		roms:        make([]menu.ROMInfo, 0),
		scanCache:   make(map[string][]menu.ROMInfo),
		maxCacheAge: time.Minute * 5, // Cache for 5 minutes
	}
}

// ScanDirectory scans a directory for ROM files and subdirectories
func (fb *FileBrowser) ScanDirectory(path string) ([]menu.ROMInfo, error) {
	// Check cache first
	if cached, exists := fb.scanCache[path]; exists {
		if time.Since(fb.lastScan) < fb.maxCacheAge {
			return cached, nil
		}
	}

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

	// Cache the results
	fb.scanCache[path] = roms
	fb.lastScan = time.Now()

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
	// Clear cache for current directory
	delete(fb.scanCache, fb.currentDir)
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
