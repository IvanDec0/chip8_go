package export

import (
	"chip8/internal/browser"
	"chip8/internal/menu"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ExportFormat represents supported export formats
type ExportFormat string

const (
	FormatJSON ExportFormat = "json"
	FormatCSV  ExportFormat = "csv"
	FormatTXT  ExportFormat = "txt"
)

// ExportOptions configures export behavior
type ExportOptions struct {
	Format          ExportFormat `json:"format"`
	IncludeMetadata bool         `json:"include_metadata"`
	IncludePaths    bool         `json:"include_paths"`
	SortBy          string       `json:"sort_by"` // "name", "date", "size"
	FilterBy        string       `json:"filter_by"`
}

// ExportManager handles data export functionality
type ExportManager struct {
	outputDir string
}

// NewExportManager creates a new export manager
func NewExportManager(outputDir string) *ExportManager {
	return &ExportManager{
		outputDir: outputDir,
	}
}

// ExportROMs exports ROM list to specified format
func (em *ExportManager) ExportROMs(roms []menu.ROMInfo, filename string, options ExportOptions) error {
	if err := os.MkdirAll(em.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	outputPath := filepath.Join(em.outputDir, filename)

	// Filter ROMs if specified
	filteredROMs := em.filterROMs(roms, options.FilterBy)

	// Sort ROMs if specified
	sortedROMs := em.sortROMs(filteredROMs, options.SortBy)

	switch options.Format {
	case FormatJSON:
		return em.exportROMsJSON(sortedROMs, outputPath, options)
	case FormatCSV:
		return em.exportROMsCSV(sortedROMs, outputPath, options)
	case FormatTXT:
		return em.exportROMsTXT(sortedROMs, outputPath, options)
	default:
		return fmt.Errorf("unsupported export format: %s", options.Format)
	}
}

// ExportFavorites exports favorites list
func (em *ExportManager) ExportFavorites(favManager *browser.FavoritesManager, filename string, options ExportOptions) error {
	favorites := favManager.GetFavorites()

	// Convert favorites to ROMInfo format
	var roms []menu.ROMInfo
	for _, fav := range favorites {
		roms = append(roms, menu.ROMInfo{
			Name:        fav.Name,
			Path:        fav.Path,
			Size:        0, // Size not available in favorites
			ModTime:     fav.DateAdded,
			IsDirectory: false,
		})
	}

	if !strings.Contains(filename, "favorites") {
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)
		filename = base + "_favorites" + ext
	}

	return em.ExportROMs(roms, filename, options)
}

// ExportRecent exports recent ROMs list
func (em *ExportManager) ExportRecent(recentManager *browser.RecentManager, filename string, options ExportOptions) error {
	recent := recentManager.GetRecentROMs()

	// Convert recent to ROMInfo format
	var roms []menu.ROMInfo
	for _, rec := range recent {
		roms = append(roms, menu.ROMInfo{
			Name:        rec.Name,
			Path:        rec.Path,
			Size:        0, // Size not available in recent
			ModTime:     rec.LastPlayed,
			IsDirectory: false,
		})
	}

	if !strings.Contains(filename, "recent") {
		ext := filepath.Ext(filename)
		base := strings.TrimSuffix(filename, ext)
		filename = base + "_recent" + ext
	}

	return em.ExportROMs(roms, filename, options)
}

// ExportCollections exports ROM collections
func (em *ExportManager) ExportCollections(collectionManager *browser.CollectionManager, filename string, options ExportOptions) error {
	collections := collectionManager.GetAllCollections()

	if err := os.MkdirAll(em.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	outputPath := filepath.Join(em.outputDir, filename)

	switch options.Format {
	case FormatJSON:
		return em.exportCollectionsJSON(collections, outputPath)
	case FormatCSV:
		return em.exportCollectionsCSV(collections, outputPath)
	case FormatTXT:
		return em.exportCollectionsTXT(collections, outputPath)
	default:
		return fmt.Errorf("unsupported export format: %s", options.Format)
	}
}

// ExportStatistics exports usage statistics
func (em *ExportManager) ExportStatistics(statsManager *menu.StatisticsManager, filename string, options ExportOptions) error {
	if err := os.MkdirAll(em.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	outputPath := filepath.Join(em.outputDir, filename)

	stats, err := statsManager.ExportStatistics()
	if err != nil {
		return fmt.Errorf("failed to get statistics: %v", err)
	}

	switch options.Format {
	case FormatJSON:
		return em.exportStatisticsJSON(stats, outputPath)
	case FormatCSV:
		return em.exportStatisticsCSV(stats, outputPath)
	case FormatTXT:
		return em.exportStatisticsTXT(stats, outputPath)
	default:
		return fmt.Errorf("unsupported export format: %s", options.Format)
	}
}

// filterROMs filters ROMs based on criteria
func (em *ExportManager) filterROMs(roms []menu.ROMInfo, filterBy string) []menu.ROMInfo {
	if filterBy == "" {
		return roms
	}

	var filtered []menu.ROMInfo
	filterLower := strings.ToLower(filterBy)

	for _, rom := range roms {
		if strings.Contains(strings.ToLower(rom.Name), filterLower) ||
			strings.Contains(strings.ToLower(rom.Path), filterLower) {
			filtered = append(filtered, rom)
		}
	}

	return filtered
}

// sortROMs sorts ROMs based on criteria
func (em *ExportManager) sortROMs(roms []menu.ROMInfo, sortBy string) []menu.ROMInfo {
	if sortBy == "" {
		return roms
	}

	sorted := make([]menu.ROMInfo, len(roms))
	copy(sorted, roms)

	switch strings.ToLower(sortBy) {
	case "name":
		// Already sorted by name in most cases
	case "date":
		// Sort by modification time (newest first)
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].ModTime.Before(sorted[j].ModTime) {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
	case "size":
		// Sort by size (largest first)
		for i := 0; i < len(sorted)-1; i++ {
			for j := i + 1; j < len(sorted); j++ {
				if sorted[i].Size < sorted[j].Size {
					sorted[i], sorted[j] = sorted[j], sorted[i]
				}
			}
		}
	}

	return sorted
}

// exportROMsJSON exports ROMs to JSON format
func (em *ExportManager) exportROMsJSON(roms []menu.ROMInfo, outputPath string, options ExportOptions) error {
	data := map[string]interface{}{
		"export_time":    time.Now(),
		"export_options": options,
		"rom_count":      len(roms),
		"roms":           roms,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}

// exportROMsCSV exports ROMs to CSV format
func (em *ExportManager) exportROMsCSV(roms []menu.ROMInfo, outputPath string, options ExportOptions) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Name", "Size", "Modified"}
	if options.IncludePaths {
		header = append(header, "Path")
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write data
	for _, rom := range roms {
		if rom.IsDirectory {
			continue // Skip directories in CSV export
		}

		record := []string{
			rom.Name,
			strconv.FormatInt(rom.Size, 10),
			rom.ModTime.Format("2006-01-02 15:04:05"),
		}

		if options.IncludePaths {
			record = append(record, rom.Path)
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	return nil
}

// exportROMsTXT exports ROMs to text format
func (em *ExportManager) exportROMsTXT(roms []menu.ROMInfo, outputPath string, options ExportOptions) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create text file: %v", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "ROM Export - Generated at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "Total ROMs: %d\n\n", len(roms))

	// Write ROM list
	for i, rom := range roms {
		if rom.IsDirectory {
			continue // Skip directories
		}

		fmt.Fprintf(file, "%d. %s\n", i+1, rom.Name)
		fmt.Fprintf(file, "   Size: %d bytes\n", rom.Size)
		fmt.Fprintf(file, "   Modified: %s\n", rom.ModTime.Format("2006-01-02 15:04:05"))

		if options.IncludePaths {
			fmt.Fprintf(file, "   Path: %s\n", rom.Path)
		}

		fmt.Fprintln(file)
	}

	return nil
}

// exportCollectionsJSON exports collections to JSON
func (em *ExportManager) exportCollectionsJSON(collections []*browser.Collection, outputPath string) error {
	data := map[string]interface{}{
		"export_time":      time.Now(),
		"collection_count": len(collections),
		"collections":      collections,
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal collections JSON: %v", err)
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}

// exportCollectionsCSV exports collections to CSV
func (em *ExportManager) exportCollectionsCSV(collections []*browser.Collection, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Collection", "Description", "ROM Count", "Created", "Updated"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write data
	for _, collection := range collections {
		record := []string{
			collection.Name,
			collection.Description,
			strconv.Itoa(len(collection.ROMs)),
			collection.CreatedAt.Format("2006-01-02 15:04:05"),
			collection.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	return nil
}

// exportCollectionsTXT exports collections to text format
func (em *ExportManager) exportCollectionsTXT(collections []*browser.Collection, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create text file: %v", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "Collections Export - Generated at %s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(file, "Total Collections: %d\n\n", len(collections))

	// Write collections
	for i, collection := range collections {
		fmt.Fprintf(file, "%d. %s\n", i+1, collection.Name)
		fmt.Fprintf(file, "   Description: %s\n", collection.Description)
		fmt.Fprintf(file, "   ROM Count: %d\n", len(collection.ROMs))
		fmt.Fprintf(file, "   Created: %s\n", collection.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Fprintf(file, "   Updated: %s\n", collection.UpdatedAt.Format("2006-01-02 15:04:05"))

		if len(collection.ROMs) > 0 {
			fmt.Fprintf(file, "   ROMs:\n")
			for _, rom := range collection.ROMs {
				fmt.Fprintf(file, "     - %s\n", rom.Name)
			}
		}

		fmt.Fprintln(file)
	}

	return nil
}

// exportStatisticsJSON exports statistics to JSON
func (em *ExportManager) exportStatisticsJSON(stats map[string]interface{}, outputPath string) error {
	jsonData, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal statistics JSON: %v", err)
	}

	return os.WriteFile(outputPath, jsonData, 0644)
}

// exportStatisticsCSV exports statistics to CSV (simplified)
func (em *ExportManager) exportStatisticsCSV(stats map[string]interface{}, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Metric", "Value"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write system stats
	if systemStats, ok := stats["system_stats"]; ok {
		if statsMap, ok := systemStats.(map[string]interface{}); ok {
			for key, value := range statsMap {
				record := []string{key, fmt.Sprintf("%v", value)}
				if err := writer.Write(record); err != nil {
					return fmt.Errorf("failed to write CSV record: %v", err)
				}
			}
		}
	}

	return nil
}

// exportStatisticsTXT exports statistics to text format
func (em *ExportManager) exportStatisticsTXT(stats map[string]interface{}, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create text file: %v", err)
	}
	defer file.Close()

	// Write header
	fmt.Fprintf(file, "Usage Statistics Export - Generated at %s\n\n", time.Now().Format("2006-01-02 15:04:05"))

	// Write system statistics
	if systemStats, ok := stats["system_stats"]; ok {
		fmt.Fprintf(file, "=== System Statistics ===\n")
		if statsMap, ok := systemStats.(map[string]interface{}); ok {
			for key, value := range statsMap {
				fmt.Fprintf(file, "%s: %v\n", strings.Title(strings.ReplaceAll(key, "_", " ")), value)
			}
		}
		fmt.Fprintln(file)
	}

	// Write top ROMs
	if topROMs, ok := stats["top_roms"]; ok {
		fmt.Fprintf(file, "=== Most Played ROMs ===\n")
		if roms, ok := topROMs.([]interface{}); ok {
			for i, rom := range roms {
				if romMap, ok := rom.(map[string]interface{}); ok {
					fmt.Fprintf(file, "%d. %s\n", i+1, romMap["name"])
					fmt.Fprintf(file, "   Play Count: %v\n", romMap["play_count"])
					fmt.Fprintf(file, "   Total Play Time: %v\n", romMap["total_play_time"])
					fmt.Fprintln(file)
				}
			}
		}
	}

	return nil
}

// GetSupportedFormats returns list of supported export formats
func (em *ExportManager) GetSupportedFormats() []ExportFormat {
	return []ExportFormat{FormatJSON, FormatCSV, FormatTXT}
}

// ValidateOptions validates export options
func (em *ExportManager) ValidateOptions(options ExportOptions) error {
	validFormats := map[ExportFormat]bool{
		FormatJSON: true,
		FormatCSV:  true,
		FormatTXT:  true,
	}

	if !validFormats[options.Format] {
		return fmt.Errorf("unsupported format: %s", options.Format)
	}

	return nil
}
