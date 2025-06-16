package browser

import (
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Metadata holds ROM metadata information
type Metadata struct {
	Title       string    `json:"title"`
	Author      string    `json:"author,omitempty"`
	Description []string  `json:"description,omitempty"`
	Controls    []string  `json:"controls,omitempty"`
	Year        string    `json:"year,omitempty"`
	System      string    `json:"system,omitempty"`
	LastUpdated time.Time `json:"last_updated"`
}

// MetadataCache provides caching for ROM metadata
type MetadataCache struct {
	cache      map[string]*Metadata
	expiration time.Duration
	mutex      sync.RWMutex
}

// NewMetadataCache creates a new metadata cache
func NewMetadataCache(expiration time.Duration) *MetadataCache {
	return &MetadataCache{
		cache:      make(map[string]*Metadata),
		expiration: expiration,
	}
}

// Get retrieves metadata from cache
func (mc *MetadataCache) Get(key string) (*Metadata, bool) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	metadata, exists := mc.cache[key]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Since(metadata.LastUpdated) > mc.expiration {
		return nil, false
	}

	return metadata, true
}

// Set stores metadata in cache
func (mc *MetadataCache) Set(key string, metadata *Metadata) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	metadata.LastUpdated = time.Now()
	mc.cache[key] = metadata
}

// Clear removes all entries from cache
func (mc *MetadataCache) Clear() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.cache = make(map[string]*Metadata)
}

// Size returns the number of entries in cache
func (mc *MetadataCache) Size() int {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	return len(mc.cache)
}

// MetadataExtractor handles metadata extraction from text files
type MetadataExtractor struct {
	cache  *MetadataCache
	parser *TextParser
}

// NewMetadataExtractor creates a new metadata extractor
func NewMetadataExtractor() *MetadataExtractor {
	return &MetadataExtractor{
		cache:  NewMetadataCache(time.Hour * 24), // Cache for 24 hours
		parser: NewTextParser(),
	}
}

// ExtractMetadata extracts metadata for a ROM file
func (me *MetadataExtractor) ExtractMetadata(romPath string) (*Metadata, error) {
	// Check cache first
	if metadata, exists := me.cache.Get(romPath); exists {
		return metadata, nil
	}

	// Look for accompanying text file
	textPath := me.findTextFile(romPath)
	var metadata *Metadata

	if textPath != "" {
		// Parse metadata from text file
		parsedFields, err := me.parser.ParseFile(textPath)
		if err == nil {
			metadata = &Metadata{
				Title:       parsedFields.Title,
				Author:      parsedFields.Author,
				Description: parsedFields.Description,
				Controls:    parsedFields.Controls,
				Year:        parsedFields.Year,
				System:      parsedFields.System,
			}
		}
	}

	// If no metadata found, generate from filename
	if metadata == nil || metadata.Title == "" {
		metadata = me.generateFallbackMetadata(romPath)
	}

	// Cache the result
	me.cache.Set(romPath, metadata)
	return metadata, nil
}

// GetCachedMetadata retrieves metadata from cache without extraction
func (me *MetadataExtractor) GetCachedMetadata(romPath string) (*Metadata, bool) {
	return me.cache.Get(romPath)
}

// ClearCache clears the metadata cache
func (me *MetadataExtractor) ClearCache() {
	me.cache.Clear()
}

// findTextFile looks for accompanying text files for a ROM
func (me *MetadataExtractor) findTextFile(romPath string) string {
	// Remove extension and try different text file extensions
	basePath := strings.TrimSuffix(romPath, filepath.Ext(romPath))

	// Common text file extensions
	extensions := []string{".txt", ".TXT", ".doc", ".DOC"}

	for _, ext := range extensions {
		textPath := basePath + ext
		if fileExists(textPath) {
			return textPath
		}
	}

	return ""
}

// generateFallbackMetadata generates basic metadata from filename
func (me *MetadataExtractor) generateFallbackMetadata(romPath string) *Metadata {
	filename := filepath.Base(romPath)
	name := strings.TrimSuffix(filename, filepath.Ext(filename))

	// Try to extract information from common naming patterns
	title, author, year := parseFilename(name)

	return &Metadata{
		Title:  title,
		Author: author,
		Year:   year,
		System: "CHIP-8",
	}
}

// parseFilename extracts title, author, and year from common filename patterns
func parseFilename(name string) (title, author, year string) {
	title = name

	// Pattern: "Game Name [Author, Year]"
	if strings.Contains(name, "[") && strings.Contains(name, "]") {
		parts := strings.Split(name, "[")
		if len(parts) >= 2 {
			title = strings.TrimSpace(parts[0])
			authorYear := strings.TrimSpace(strings.TrimSuffix(parts[1], "]"))

			// Try to extract author and year
			if strings.Contains(authorYear, ",") {
				authorParts := strings.Split(authorYear, ",")
				if len(authorParts) >= 2 {
					author = strings.TrimSpace(authorParts[0])
					year = strings.TrimSpace(authorParts[len(authorParts)-1])
				}
			} else {
				// Could be just author or just year
				if isYear(authorYear) {
					year = authorYear
				} else {
					author = authorYear
				}
			}
		}
	}

	// Clean up title - remove extra spaces
	title = strings.Join(strings.Fields(title), " ")

	return
}

// isYear checks if a string represents a year (4 digits starting with 19 or 20)
func isYear(s string) bool {
	if len(s) != 4 {
		return false
	}

	return strings.HasPrefix(s, "19") || strings.HasPrefix(s, "20")
}
