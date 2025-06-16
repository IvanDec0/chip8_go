package browser

import (
	"sync"
	"time"
)

// CacheType represents different types of cached data
type CacheType int

const (
	CacheTypeMetadata CacheType = iota
	CacheTypeDirectory
	CacheTypeSearch
)

// CacheEntry represents a single cache entry
type CacheEntry struct {
	Value       interface{}
	Expiry      time.Time
	AccessTime  time.Time
	AccessCount int
}

// CacheStats provides cache statistics
type CacheStats struct {
	TotalEntries     int
	MetadataEntries  int
	DirectoryEntries int
	SearchEntries    int
	HitRate          float64
	TotalHits        int64
	TotalMisses      int64
}

// CacheConfig configures caching behavior
type CacheConfig struct {
	MaxMetadataEntries  int
	MaxDirectoryEntries int
	MaxSearchEntries    int
	MetadataTTL         time.Duration
	DirectoryTTL        time.Duration
	SearchTTL           time.Duration
	CleanupInterval     time.Duration
}

// DefaultCacheConfig returns a default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		MaxMetadataEntries:  1000,
		MaxDirectoryEntries: 50,
		MaxSearchEntries:    100,
		MetadataTTL:         time.Hour * 24,
		DirectoryTTL:        time.Minute * 5,
		SearchTTL:           time.Minute * 10,
		CleanupInterval:     time.Minute * 15,
	}
}

// CacheManager handles various caching strategies
type CacheManager struct {
	metadataCache  map[string]*CacheEntry
	directoryCache map[string]*CacheEntry
	searchCache    map[string]*CacheEntry
	config         CacheConfig
	mutex          sync.RWMutex
	stats          CacheStats
	cleanupTicker  *time.Ticker
	stopCleanup    chan bool
}

// NewCacheManager creates a new cache manager
func NewCacheManager(config CacheConfig) *CacheManager {
	return &CacheManager{
		metadataCache:  make(map[string]*CacheEntry),
		directoryCache: make(map[string]*CacheEntry),
		searchCache:    make(map[string]*CacheEntry),
		config:         config,
		stopCleanup:    make(chan bool),
	}
}

// Get retrieves a value from cache
func (cm *CacheManager) Get(key string, cacheType CacheType) (interface{}, bool) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	var entry *CacheEntry
	var exists bool

	switch cacheType {
	case CacheTypeMetadata:
		entry, exists = cm.metadataCache[key]
	case CacheTypeDirectory:
		entry, exists = cm.directoryCache[key]
	case CacheTypeSearch:
		entry, exists = cm.searchCache[key]
	default:
		cm.stats.TotalMisses++
		return nil, false
	}

	if !exists {
		cm.stats.TotalMisses++
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.Expiry) {
		cm.stats.TotalMisses++
		return nil, false
	}

	// Update access statistics
	entry.AccessTime = time.Now()
	entry.AccessCount++
	cm.stats.TotalHits++

	return entry.Value, true
}

// Set stores a value in cache
func (cm *CacheManager) Set(key string, value interface{}, cacheType CacheType, ttl time.Duration) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	entry := &CacheEntry{
		Value:       value,
		Expiry:      time.Now().Add(ttl),
		AccessTime:  time.Now(),
		AccessCount: 0,
	}

	switch cacheType {
	case CacheTypeMetadata:
		cm.metadataCache[key] = entry
		cm.evictLRU(CacheTypeMetadata, cm.config.MaxMetadataEntries)
	case CacheTypeDirectory:
		cm.directoryCache[key] = entry
		cm.evictLRU(CacheTypeDirectory, cm.config.MaxDirectoryEntries)
	case CacheTypeSearch:
		cm.searchCache[key] = entry
		cm.evictLRU(CacheTypeSearch, cm.config.MaxSearchEntries)
	}
}

// Clear removes all entries of a specific type
func (cm *CacheManager) Clear(cacheType CacheType) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	switch cacheType {
	case CacheTypeMetadata:
		cm.metadataCache = make(map[string]*CacheEntry)
	case CacheTypeDirectory:
		cm.directoryCache = make(map[string]*CacheEntry)
	case CacheTypeSearch:
		cm.searchCache = make(map[string]*CacheEntry)
	}
}

// ClearAll removes all cache entries
func (cm *CacheManager) ClearAll() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.metadataCache = make(map[string]*CacheEntry)
	cm.directoryCache = make(map[string]*CacheEntry)
	cm.searchCache = make(map[string]*CacheEntry)
	cm.stats = CacheStats{}
}

// evictLRU removes least recently used entries when cache exceeds max size
func (cm *CacheManager) evictLRU(cacheType CacheType, maxEntries int) {
	var cache map[string]*CacheEntry

	switch cacheType {
	case CacheTypeMetadata:
		cache = cm.metadataCache
	case CacheTypeDirectory:
		cache = cm.directoryCache
	case CacheTypeSearch:
		cache = cm.searchCache
	default:
		return
	}

	if len(cache) <= maxEntries {
		return
	}

	// Find least recently used entries
	type keyTime struct {
		key        string
		accessTime time.Time
	}

	var entries []keyTime
	for key, entry := range cache {
		entries = append(entries, keyTime{key, entry.AccessTime})
	}

	// Sort by access time (oldest first)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[i].accessTime.After(entries[j].accessTime) {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	// Remove oldest entries
	toRemove := len(cache) - maxEntries
	for i := 0; i < toRemove && i < len(entries); i++ {
		delete(cache, entries[i].key)
	}
}

// cleanupExpired removes expired entries
func (cm *CacheManager) cleanupExpired() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	now := time.Now()

	// Clean metadata cache
	for key, entry := range cm.metadataCache {
		if now.After(entry.Expiry) {
			delete(cm.metadataCache, key)
		}
	}

	// Clean directory cache
	for key, entry := range cm.directoryCache {
		if now.After(entry.Expiry) {
			delete(cm.directoryCache, key)
		}
	}

	// Clean search cache
	for key, entry := range cm.searchCache {
		if now.After(entry.Expiry) {
			delete(cm.searchCache, key)
		}
	}
}

// StartCleanup starts the periodic cleanup routine
func (cm *CacheManager) StartCleanup() {
	if cm.cleanupTicker != nil {
		return // Already started
	}

	cm.cleanupTicker = time.NewTicker(cm.config.CleanupInterval)
	go func() {
		for {
			select {
			case <-cm.cleanupTicker.C:
				cm.cleanupExpired()
			case <-cm.stopCleanup:
				return
			}
		}
	}()
}

// StopCleanup stops the periodic cleanup routine
func (cm *CacheManager) StopCleanup() {
	if cm.cleanupTicker != nil {
		cm.cleanupTicker.Stop()
		cm.cleanupTicker = nil
		close(cm.stopCleanup)
		cm.stopCleanup = make(chan bool)
	}
}

// GetStats returns current cache statistics
func (cm *CacheManager) GetStats() CacheStats {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	stats := cm.stats
	stats.TotalEntries = len(cm.metadataCache) + len(cm.directoryCache) + len(cm.searchCache)
	stats.MetadataEntries = len(cm.metadataCache)
	stats.DirectoryEntries = len(cm.directoryCache)
	stats.SearchEntries = len(cm.searchCache)

	if stats.TotalHits+stats.TotalMisses > 0 {
		stats.HitRate = float64(stats.TotalHits) / float64(stats.TotalHits+stats.TotalMisses)
	}

	return stats
}

// SetTTL updates the TTL for a cache type
func (cm *CacheManager) SetTTL(cacheType CacheType, ttl time.Duration) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	switch cacheType {
	case CacheTypeMetadata:
		cm.config.MetadataTTL = ttl
	case CacheTypeDirectory:
		cm.config.DirectoryTTL = ttl
	case CacheTypeSearch:
		cm.config.SearchTTL = ttl
	}
}

// GetSize returns the number of entries in a specific cache
func (cm *CacheManager) GetSize(cacheType CacheType) int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()

	switch cacheType {
	case CacheTypeMetadata:
		return len(cm.metadataCache)
	case CacheTypeDirectory:
		return len(cm.directoryCache)
	case CacheTypeSearch:
		return len(cm.searchCache)
	default:
		return 0
	}
}
