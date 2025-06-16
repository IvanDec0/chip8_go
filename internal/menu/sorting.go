package menu

import (
	"sort"
	"strings"
)

// SortManager handles ROM list sorting
type SortManager struct {
	currentSort SortCriteria
	ascending   bool
}

// SortCriteria represents different sorting options
type SortCriteria int

const (
	SortByName SortCriteria = iota
	SortByDate
	SortBySize
	SortByType
	SortByLastPlayed
	SortByPlayCount
)

// String returns string representation of sort criteria
func (sc SortCriteria) String() string {
	switch sc {
	case SortByName:
		return "Name"
	case SortByDate:
		return "Date"
	case SortBySize:
		return "Size"
	case SortByType:
		return "Type"
	case SortByLastPlayed:
		return "Last Played"
	case SortByPlayCount:
		return "Play Count"
	default:
		return "Unknown"
	}
}

// NewSortManager creates a new sort manager
func NewSortManager() *SortManager {
	return &SortManager{
		currentSort: SortByName,
		ascending:   true,
	}
}

// SetSort sets the current sort criteria and direction
func (sm *SortManager) SetSort(criteria SortCriteria, ascending bool) {
	sm.currentSort = criteria
	sm.ascending = ascending
}

// SortROMs sorts a slice of ROMs according to current criteria
func (sm *SortManager) SortROMs(roms []ROMInfo) []ROMInfo {
	// Create a copy to avoid modifying original slice
	sorted := make([]ROMInfo, len(roms))
	copy(sorted, roms)

	// Always put directories first, then sort files
	sort.Slice(sorted, func(i, j int) bool {
		// Directories always come first
		if sorted[i].IsDirectory && !sorted[j].IsDirectory {
			return true
		}
		if !sorted[i].IsDirectory && sorted[j].IsDirectory {
			return false
		}

		// Both are same type, apply sorting criteria
		return sm.compareROMs(sorted[i], sorted[j])
	})

	return sorted
}

// compareROMs compares two ROMs based on current sort criteria
func (sm *SortManager) compareROMs(a, b ROMInfo) bool {
	var result bool

	switch sm.currentSort {
	case SortByName:
		result = sm.compareByName(a, b)
	case SortByDate:
		result = sm.compareByDate(a, b)
	case SortBySize:
		result = sm.compareBySize(a, b)
	case SortByType:
		result = sm.compareByType(a, b)
	default:
		result = sm.compareByName(a, b) // Default fallback
	}

	// Reverse if descending order
	if !sm.ascending {
		result = !result
	}

	return result
}

// compareByName compares ROMs by name
func (sm *SortManager) compareByName(a, b ROMInfo) bool {
	nameA := strings.ToLower(a.Name)
	nameB := strings.ToLower(b.Name)

	// Special handling for ".." parent directory
	if a.Name == ".." {
		return true
	}
	if b.Name == ".." {
		return false
	}

	return nameA < nameB
}

// compareByDate compares ROMs by modification date
func (sm *SortManager) compareByDate(a, b ROMInfo) bool {
	return a.ModTime.Before(b.ModTime)
}

// compareBySize compares ROMs by file size
func (sm *SortManager) compareBySize(a, b ROMInfo) bool {
	return a.Size < b.Size
}

// compareByType compares ROMs by file type/extension
func (sm *SortManager) compareByType(a, b ROMInfo) bool {
	extA := strings.ToLower(getFileExtension(a.Name))
	extB := strings.ToLower(getFileExtension(b.Name))

	if extA == extB {
		// Same extension, sort by name
		return sm.compareByName(a, b)
	}

	return extA < extB
}

// getFileExtension extracts file extension from filename
func getFileExtension(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) > 1 {
		return "." + parts[len(parts)-1]
	}
	return ""
}

// GetCurrentSort returns current sort criteria and direction
func (sm *SortManager) GetCurrentSort() (SortCriteria, bool) {
	return sm.currentSort, sm.ascending
}

// ToggleDirection toggles sort direction (ascending/descending)
func (sm *SortManager) ToggleDirection() {
	sm.ascending = !sm.ascending
}

// NextSort cycles to the next sort criteria
func (sm *SortManager) NextSort() {
	switch sm.currentSort {
	case SortByName:
		sm.currentSort = SortByDate
	case SortByDate:
		sm.currentSort = SortBySize
	case SortBySize:
		sm.currentSort = SortByType
	case SortByType:
		sm.currentSort = SortByName
	default:
		sm.currentSort = SortByName
	}
}

// PreviousSort cycles to the previous sort criteria
func (sm *SortManager) PreviousSort() {
	switch sm.currentSort {
	case SortByName:
		sm.currentSort = SortByType
	case SortByDate:
		sm.currentSort = SortByName
	case SortBySize:
		sm.currentSort = SortByDate
	case SortByType:
		sm.currentSort = SortBySize
	default:
		sm.currentSort = SortByName
	}
}

// GetSortIndicator returns a string indicator for current sort
func (sm *SortManager) GetSortIndicator() string {
	direction := "↑"
	if !sm.ascending {
		direction = "↓"
	}

	return sm.currentSort.String() + " " + direction
}

// GetAllSortCriteria returns all available sort criteria
func (sm *SortManager) GetAllSortCriteria() []SortCriteria {
	return []SortCriteria{
		SortByName,
		SortByDate,
		SortBySize,
		SortByType,
	}
}

// SetSortFromString sets sort criteria from string representation
func (sm *SortManager) SetSortFromString(sortStr string, ascending bool) bool {
	sortStr = strings.ToLower(strings.TrimSpace(sortStr))

	switch sortStr {
	case "name":
		sm.currentSort = SortByName
	case "date", "time", "modified":
		sm.currentSort = SortByDate
	case "size":
		sm.currentSort = SortBySize
	case "type", "extension", "ext":
		sm.currentSort = SortByType
	default:
		return false // Invalid sort criteria
	}

	sm.ascending = ascending
	return true
}

// ApplySort applies current sort to ROM list and returns sorted copy
func (sm *SortManager) ApplySort(roms []ROMInfo) []ROMInfo {
	return sm.SortROMs(roms)
}

// SortByCustom sorts ROMs using a custom comparison function
func (sm *SortManager) SortByCustom(roms []ROMInfo, less func(i, j ROMInfo) bool) []ROMInfo {
	sorted := make([]ROMInfo, len(roms))
	copy(sorted, roms)

	sort.Slice(sorted, func(i, j int) bool {
		return less(sorted[i], sorted[j])
	})

	return sorted
}

// GroupByType groups ROMs by file type
func (sm *SortManager) GroupByType(roms []ROMInfo) map[string][]ROMInfo {
	groups := make(map[string][]ROMInfo)

	for _, rom := range roms {
		var group string

		if rom.IsDirectory {
			group = "Directories"
		} else {
			ext := getFileExtension(rom.Name)
			if ext == "" {
				group = "No Extension"
			} else {
				group = strings.ToUpper(ext) + " Files"
			}
		}

		groups[group] = append(groups[group], rom)
	}

	// Sort within each group
	for groupName, groupRoms := range groups {
		groups[groupName] = sm.SortROMs(groupRoms)
	}

	return groups
}

// FilterAndSort applies filtering and sorting in one operation
func (sm *SortManager) FilterAndSort(roms []ROMInfo, filter func(ROMInfo) bool) []ROMInfo {
	// First filter
	var filtered []ROMInfo
	for _, rom := range roms {
		if filter == nil || filter(rom) {
			filtered = append(filtered, rom)
		}
	}

	// Then sort
	return sm.SortROMs(filtered)
}

// GetRecommendedSort returns recommended sort for given ROM count
func (sm *SortManager) GetRecommendedSort(romCount int) SortCriteria {
	if romCount > 100 {
		return SortByName // Large collections benefit from alphabetical
	} else if romCount > 20 {
		return SortByDate // Medium collections might want recent first
	} else {
		return SortByName // Small collections, alphabetical is fine
	}
}

// IsSortStable checks if current sort is stable (deterministic)
func (sm *SortManager) IsSortStable() bool {
	// Name sorting is always stable
	return sm.currentSort == SortByName
}

// CreateSortedIndex creates an index of sorted positions
func (sm *SortManager) CreateSortedIndex(roms []ROMInfo) []int {
	// Create index array
	indices := make([]int, len(roms))
	for i := range indices {
		indices[i] = i
	}

	// Sort indices based on ROM comparison
	sort.Slice(indices, func(i, j int) bool {
		return sm.compareROMs(roms[indices[i]], roms[indices[j]])
	})

	return indices
}
