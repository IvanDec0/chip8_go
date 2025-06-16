package browser

import (
	"chip8/internal/menu"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Collection represents a user-defined collection of ROMs
type Collection struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	ROMs        []menu.ROMInfo    `json:"roms"`
	Tags        []string          `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	Metadata    map[string]string `json:"metadata"`
}

// CollectionManager manages ROM collections
type CollectionManager struct {
	collections    map[string]*Collection
	dataPath       string
	maxCollections int
}

// NewCollectionManager creates a new collection manager
func NewCollectionManager(dataPath string) *CollectionManager {
	manager := &CollectionManager{
		collections:    make(map[string]*Collection),
		dataPath:       filepath.Join(dataPath, "collections"),
		maxCollections: 50, // Reasonable limit
	}

	// Ensure data directory exists
	if err := os.MkdirAll(manager.dataPath, 0755); err == nil {
		manager.loadCollections()
	}

	return manager
}

// CreateCollection creates a new ROM collection
func (cm *CollectionManager) CreateCollection(name, description string) (*Collection, error) {
	if name == "" {
		return nil, fmt.Errorf("collection name cannot be empty")
	}

	if len(cm.collections) >= cm.maxCollections {
		return nil, fmt.Errorf("maximum number of collections (%d) reached", cm.maxCollections)
	}

	// Check for duplicate names
	for _, collection := range cm.collections {
		if collection.Name == name {
			return nil, fmt.Errorf("collection with name '%s' already exists", name)
		}
	}

	id := generateCollectionID()
	collection := &Collection{
		ID:          id,
		Name:        name,
		Description: description,
		ROMs:        make([]menu.ROMInfo, 0),
		Tags:        make([]string, 0),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Metadata:    make(map[string]string),
	}

	cm.collections[id] = collection

	if err := cm.saveCollection(collection); err != nil {
		delete(cm.collections, id)
		return nil, fmt.Errorf("failed to save collection: %v", err)
	}

	return collection, nil
}

// AddToCollection adds a ROM to a collection
func (cm *CollectionManager) AddToCollection(collectionID string, rom menu.ROMInfo) error {
	collection, exists := cm.collections[collectionID]
	if !exists {
		return fmt.Errorf("collection not found: %s", collectionID)
	}

	// Check if ROM is already in collection
	for _, existingROM := range collection.ROMs {
		if existingROM.Path == rom.Path {
			return fmt.Errorf("ROM already exists in collection")
		}
	}

	collection.ROMs = append(collection.ROMs, rom)
	collection.UpdatedAt = time.Now()

	return cm.saveCollection(collection)
}

// RemoveFromCollection removes a ROM from a collection
func (cm *CollectionManager) RemoveFromCollection(collectionID string, romPath string) error {
	collection, exists := cm.collections[collectionID]
	if !exists {
		return fmt.Errorf("collection not found: %s", collectionID)
	}

	for i, rom := range collection.ROMs {
		if rom.Path == romPath {
			// Remove ROM from slice
			collection.ROMs = append(collection.ROMs[:i], collection.ROMs[i+1:]...)
			collection.UpdatedAt = time.Now()
			return cm.saveCollection(collection)
		}
	}

	return fmt.Errorf("ROM not found in collection")
}

// DeleteCollection deletes a collection
func (cm *CollectionManager) DeleteCollection(collectionID string) error {
	_, exists := cm.collections[collectionID]
	if !exists {
		return fmt.Errorf("collection not found: %s", collectionID)
	}

	// Remove from memory
	delete(cm.collections, collectionID)

	// Remove file
	filename := filepath.Join(cm.dataPath, fmt.Sprintf("%s.json", collectionID))
	return os.Remove(filename)
}

// GetCollection gets a collection by ID
func (cm *CollectionManager) GetCollection(collectionID string) (*Collection, error) {
	collection, exists := cm.collections[collectionID]
	if !exists {
		return nil, fmt.Errorf("collection not found: %s", collectionID)
	}
	return collection, nil
}

// GetAllCollections returns all collections
func (cm *CollectionManager) GetAllCollections() []*Collection {
	collections := make([]*Collection, 0, len(cm.collections))
	for _, collection := range cm.collections {
		collections = append(collections, collection)
	}

	// Sort by creation date (newest first)
	sort.Slice(collections, func(i, j int) bool {
		return collections[i].CreatedAt.After(collections[j].CreatedAt)
	})

	return collections
}

// SearchCollections searches collections by name or description
func (cm *CollectionManager) SearchCollections(query string) []*Collection {
	if query == "" {
		return cm.GetAllCollections()
	}

	var results []*Collection
	queryLower := strings.ToLower(query)

	for _, collection := range cm.collections {
		if strings.Contains(strings.ToLower(collection.Name), queryLower) ||
			strings.Contains(strings.ToLower(collection.Description), queryLower) {
			results = append(results, collection)
		}
	}

	// Sort by relevance (name matches first)
	sort.Slice(results, func(i, j int) bool {
		nameMatchI := strings.Contains(strings.ToLower(results[i].Name), queryLower)
		nameMatchJ := strings.Contains(strings.ToLower(results[j].Name), queryLower)

		if nameMatchI && !nameMatchJ {
			return true
		}
		if !nameMatchI && nameMatchJ {
			return false
		}

		return results[i].CreatedAt.After(results[j].CreatedAt)
	})

	return results
}

// UpdateCollection updates collection metadata
func (cm *CollectionManager) UpdateCollection(collectionID, name, description string, tags []string) error {
	collection, exists := cm.collections[collectionID]
	if !exists {
		return fmt.Errorf("collection not found: %s", collectionID)
	}

	if name != "" {
		// Check for duplicate names (excluding current collection)
		for id, c := range cm.collections {
			if id != collectionID && c.Name == name {
				return fmt.Errorf("collection with name '%s' already exists", name)
			}
		}
		collection.Name = name
	}

	if description != "" {
		collection.Description = description
	}

	if tags != nil {
		collection.Tags = tags
	}

	collection.UpdatedAt = time.Now()

	return cm.saveCollection(collection)
}

// ExportCollection exports a collection to JSON
func (cm *CollectionManager) ExportCollection(collectionID string) ([]byte, error) {
	collection, exists := cm.collections[collectionID]
	if !exists {
		return nil, fmt.Errorf("collection not found: %s", collectionID)
	}

	return json.MarshalIndent(collection, "", "  ")
}

// ImportCollection imports a collection from JSON
func (cm *CollectionManager) ImportCollection(data []byte) (*Collection, error) {
	var collection Collection
	if err := json.Unmarshal(data, &collection); err != nil {
		return nil, fmt.Errorf("failed to parse collection: %v", err)
	}

	// Generate new ID to avoid conflicts
	collection.ID = generateCollectionID()
	collection.CreatedAt = time.Now()
	collection.UpdatedAt = time.Now()

	// Check for duplicate names
	originalName := collection.Name
	counter := 1
	for {
		nameExists := false
		for _, existing := range cm.collections {
			if existing.Name == collection.Name {
				nameExists = true
				break
			}
		}
		if !nameExists {
			break
		}
		collection.Name = fmt.Sprintf("%s (%d)", originalName, counter)
		counter++
	}

	cm.collections[collection.ID] = &collection

	if err := cm.saveCollection(&collection); err != nil {
		delete(cm.collections, collection.ID)
		return nil, fmt.Errorf("failed to save imported collection: %v", err)
	}

	return &collection, nil
}

// GetCollectionStats returns statistics about collections
func (cm *CollectionManager) GetCollectionStats() map[string]interface{} {
	totalCollections := len(cm.collections)
	totalROMs := 0
	avgROMsPerCollection := 0.0

	for _, collection := range cm.collections {
		totalROMs += len(collection.ROMs)
	}

	if totalCollections > 0 {
		avgROMsPerCollection = float64(totalROMs) / float64(totalCollections)
	}

	return map[string]interface{}{
		"total_collections":           totalCollections,
		"total_roms":                  totalROMs,
		"average_roms_per_collection": avgROMsPerCollection,
		"max_collections":             cm.maxCollections,
	}
}

// saveCollection saves a collection to disk
func (cm *CollectionManager) saveCollection(collection *Collection) error {
	filename := filepath.Join(cm.dataPath, fmt.Sprintf("%s.json", collection.ID))

	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// loadCollections loads all collections from disk
func (cm *CollectionManager) loadCollections() {
	files, err := os.ReadDir(cm.dataPath)
	if err != nil {
		return
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			cm.loadCollection(file.Name())
		}
	}
}

// loadCollection loads a single collection from disk
func (cm *CollectionManager) loadCollection(filename string) {
	path := filepath.Join(cm.dataPath, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	var collection Collection
	if err := json.Unmarshal(data, &collection); err != nil {
		return
	}

	cm.collections[collection.ID] = &collection
}

// generateCollectionID generates a unique collection ID
func generateCollectionID() string {
	return fmt.Sprintf("col_%d", time.Now().UnixNano())
}
