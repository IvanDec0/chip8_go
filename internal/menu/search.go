package menu

import (
	"strings"
	"unicode"
)

// SearchManager handles ROM search functionality
type SearchManager struct {
	query         string
	caseSensitive bool
	isActive      bool
	lastResults   []ROMInfo
}

// SearchFilter represents different search criteria
type SearchFilter struct {
	ByName        bool
	ByAuthor      bool
	ByDescription bool
	ByTags        bool
}

// SearchResult represents a search result with additional context
type SearchResult struct {
	ROM           ROMInfo
	MatchedFields []string // Which fields matched the search
	Score         int      // Relevance score
}

// NewSearchManager creates a new search manager
func NewSearchManager() *SearchManager {
	return &SearchManager{
		query:         "",
		caseSensitive: false,
		isActive:      false,
		lastResults:   make([]ROMInfo, 0),
	}
}

// SetQuery sets the search query
func (sm *SearchManager) SetQuery(query string) error {
	sm.query = strings.TrimSpace(query)
	sm.isActive = sm.query != ""
	return nil
}

// GetQuery returns the current search query
func (sm *SearchManager) GetQuery() string {
	return sm.query
}

// SetCaseSensitive sets whether search should be case sensitive
func (sm *SearchManager) SetCaseSensitive(caseSensitive bool) {
	sm.caseSensitive = caseSensitive
}

// IsCaseSensitive returns whether search is case sensitive
func (sm *SearchManager) IsCaseSensitive() bool {
	return sm.caseSensitive
}

// Search performs search on ROM list with given filter
func (sm *SearchManager) Search(roms []ROMInfo, filter SearchFilter) []ROMInfo {
	if sm.query == "" {
		sm.lastResults = roms
		return roms
	}

	var results []ROMInfo
	searchTerm := sm.query
	if !sm.caseSensitive {
		searchTerm = strings.ToLower(searchTerm)
	}

	for _, rom := range roms {
		if sm.matchesROM(rom, searchTerm, filter) {
			results = append(results, rom)
		}
	}

	sm.lastResults = results
	return results
}

// SearchWithScoring performs search with relevance scoring
func (sm *SearchManager) SearchWithScoring(roms []ROMInfo, filter SearchFilter) []SearchResult {
	if sm.query == "" {
		results := make([]SearchResult, len(roms))
		for i, rom := range roms {
			results[i] = SearchResult{
				ROM:           rom,
				MatchedFields: []string{},
				Score:         0,
			}
		}
		return results
	}

	var results []SearchResult
	searchTerm := sm.query
	if !sm.caseSensitive {
		searchTerm = strings.ToLower(searchTerm)
	}

	for _, rom := range roms {
		result := sm.searchROMWithScoring(rom, searchTerm, filter)
		if result.Score > 0 {
			results = append(results, result)
		}
	}

	// Sort by relevance score (highest first)
	sm.sortByScore(results)

	// Convert to ROMInfo for compatibility
	romResults := make([]ROMInfo, len(results))
	for i, result := range results {
		romResults[i] = result.ROM
	}
	sm.lastResults = romResults

	return results
}

// searchROMWithScoring performs detailed search with scoring
func (sm *SearchManager) searchROMWithScoring(rom ROMInfo, searchTerm string, filter SearchFilter) SearchResult {
	result := SearchResult{
		ROM:           rom,
		MatchedFields: []string{},
		Score:         0,
	}

	// Search in ROM name
	if filter.ByName {
		romName := rom.Name
		if !sm.caseSensitive {
			romName = strings.ToLower(romName)
		}
		if strings.Contains(romName, searchTerm) {
			result.MatchedFields = append(result.MatchedFields, "name")
			result.Score += 10 // High score for name match

			// Bonus for exact match
			if romName == searchTerm {
				result.Score += 5
			}

			// Bonus for prefix match
			if strings.HasPrefix(romName, searchTerm) {
				result.Score += 3
			}
		}
	}

	return result
}

// matchesROM checks if a ROM matches the search criteria (simple version)
func (sm *SearchManager) matchesROM(rom ROMInfo, searchTerm string, filter SearchFilter) bool {
	// Search in name (always enabled for basic search)
	romName := rom.Name
	if !sm.caseSensitive {
		romName = strings.ToLower(romName)
	}

	return strings.Contains(romName, searchTerm)
}

// sortByScore sorts search results by relevance score
func (sm *SearchManager) sortByScore(results []SearchResult) {
	// Simple bubble sort by score (descending)
	for i := 0; i < len(results)-1; i++ {
		for j := 0; j < len(results)-i-1; j++ {
			if results[j].Score < results[j+1].Score {
				results[j], results[j+1] = results[j+1], results[j]
			}
		}
	}
}

// ClearSearch clears the current search
func (sm *SearchManager) ClearSearch() {
	sm.query = ""
	sm.isActive = false
	sm.lastResults = sm.lastResults[:0]
}

// IsActive returns whether search is currently active
func (sm *SearchManager) IsActive() bool {
	return sm.isActive
}

// GetResults returns the last search results
func (sm *SearchManager) GetResults() []ROMInfo {
	return sm.lastResults
}

// GetDefaultFilter returns a default search filter
func (sm *SearchManager) GetDefaultFilter() SearchFilter {
	return SearchFilter{
		ByName:        true,
		ByAuthor:      false,
		ByDescription: false,
		ByTags:        false,
	}
}

// FilterByExtension filters ROMs by file extension
func (sm *SearchManager) FilterByExtension(roms []ROMInfo, extensions []string) []ROMInfo {
	if len(extensions) == 0 {
		return roms
	}

	var filtered []ROMInfo
	for _, rom := range roms {
		if rom.IsDirectory {
			filtered = append(filtered, rom)
			continue
		}

		for _, ext := range extensions {
			if strings.HasSuffix(strings.ToLower(rom.Name), strings.ToLower(ext)) {
				filtered = append(filtered, rom)
				break
			}
		}
	}

	return filtered
}

// HighlightMatches returns highlighted version of text with search matches
func (sm *SearchManager) HighlightMatches(text string) string {
	if sm.query == "" || text == "" {
		return text
	}

	searchTerm := sm.query
	if !sm.caseSensitive {
		searchTerm = strings.ToLower(searchTerm)
	}

	// Find all match positions
	var highlighted strings.Builder
	textRunes := []rune(text)
	searchRunes := []rune(searchTerm)
	i := 0

	for i < len(textRunes) {
		// Check if we have a match at current position
		if sm.matchesAt(textRunes, searchRunes, i) {
			// Add highlight markers (simple approach for console)
			highlighted.WriteString("[")
			for j := 0; j < len(searchRunes); j++ {
				highlighted.WriteRune(textRunes[i+j])
			}
			highlighted.WriteString("]")
			i += len(searchRunes)
		} else {
			highlighted.WriteRune(textRunes[i])
			i++
		}
	}

	return highlighted.String()
}

// matchesAt checks if search term matches at given position
func (sm *SearchManager) matchesAt(text, search []rune, pos int) bool {
	if pos+len(search) > len(text) {
		return false
	}

	for i, searchRune := range search {
		textRune := text[pos+i]
		if !sm.caseSensitive {
			textRune = unicode.ToLower(textRune)
			searchRune = unicode.ToLower(searchRune)
		}
		if textRune != searchRune {
			return false
		}
	}

	return true
}

// GetSearchSuggestions returns search suggestions based on available ROMs
func (sm *SearchManager) GetSearchSuggestions(roms []ROMInfo) []string {
	suggestions := make(map[string]bool)

	for _, rom := range roms {
		// Add ROM name words
		sm.addWordsToSuggestions(rom.Name, suggestions)
	}

	// Convert to slice
	result := make([]string, 0, len(suggestions))
	for suggestion := range suggestions {
		if len(suggestion) > 2 { // Only suggest words longer than 2 characters
			result = append(result, suggestion)
		}
	}

	return result
}

// addWordsToSuggestions extracts words from text and adds them to suggestions
func (sm *SearchManager) addWordsToSuggestions(text string, suggestions map[string]bool) {
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsNumber(c)
	})

	for _, word := range words {
		word = strings.TrimSpace(word)
		if word != "" {
			suggestions[word] = true
		}
	}
}
