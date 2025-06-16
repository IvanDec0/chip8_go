package benchmarks

import (
	"chip8/internal/menu"
	"fmt"
	"testing"
)

// generateSearchTestROMs creates test ROM data for benchmarking search
func generateSearchTestROMs(count int) []menu.ROMInfo {
	roms := make([]menu.ROMInfo, count)

	gameNames := []string{
		"Pong", "Space Invaders", "Tetris", "Pac-Man", "Asteroids",
		"Centipede", "Frogger", "Breakout", "Missile Command", "Defender",
	}

	for i := 0; i < count; i++ {
		gameName := gameNames[i%len(gameNames)]

		roms[i] = menu.ROMInfo{
			Name:        fmt.Sprintf("%s-%d.ch8", gameName, i),
			Path:        fmt.Sprintf("/test/roms/%s-%d.ch8", gameName, i),
			Size:        int64(1024 + (i % 4096)), // Vary sizes
			IsDirectory: false,
		}
	}

	return roms
}

// BenchmarkSearchBasic benchmarks basic search operations
func BenchmarkSearchBasic(b *testing.B) {
	roms := generateSearchTestROMs(1000)
	searchManager := menu.NewSearchManager()
	filter := searchManager.GetDefaultFilter()

	queries := []string{
		"pong",
		"space",
		"tetris",
		"pac",
		"game",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		searchManager.SetQuery(query)
		results := searchManager.Search(roms, filter)
		_ = results // Use results to prevent optimization
	}
}

// BenchmarkSearchLargeDataset benchmarks search on large datasets
func BenchmarkSearchLargeDataset(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping large dataset benchmark in short mode")
	}

	roms := generateSearchTestROMs(10000)
	searchManager := menu.NewSearchManager()
	filter := searchManager.GetDefaultFilter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchManager.SetQuery("game")
		results := searchManager.Search(roms, filter)
		_ = results
	}
}

// BenchmarkSearchPartialMatch benchmarks partial matching performance
func BenchmarkSearchPartialMatch(b *testing.B) {
	roms := generateSearchTestROMs(1000)
	searchManager := menu.NewSearchManager()
	filter := searchManager.GetDefaultFilter()

	// Test progressively longer queries
	queries := []string{"p", "po", "pon", "pong", "space", "invaders"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		searchManager.SetQuery(query)
		results := searchManager.Search(roms, filter)
		_ = results
	}
}

// BenchmarkSearchCaseInsensitive benchmarks case-insensitive search
func BenchmarkSearchCaseInsensitive(b *testing.B) {
	roms := generateSearchTestROMs(1000)
	searchManager := menu.NewSearchManager()
	filter := searchManager.GetDefaultFilter()

	queries := []string{
		"PONG",
		"space",
		"TeTrIs",
		"PAC",
		"Game",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query := queries[i%len(queries)]
		searchManager.SetQuery(query)
		results := searchManager.Search(roms, filter)
		_ = results
	}
}

// BenchmarkSearchEmpty benchmarks empty search queries
func BenchmarkSearchEmpty(b *testing.B) {
	roms := generateSearchTestROMs(1000)
	searchManager := menu.NewSearchManager()
	filter := searchManager.GetDefaultFilter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchManager.SetQuery("")
		results := searchManager.Search(roms, filter)
		_ = results
	}
}

// BenchmarkSearchNoResults benchmarks searches with no results
func BenchmarkSearchNoResults(b *testing.B) {
	roms := generateSearchTestROMs(1000)
	searchManager := menu.NewSearchManager()
	filter := searchManager.GetDefaultFilter()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchManager.SetQuery("xyznomatch")
		results := searchManager.Search(roms, filter)
		_ = results
	}
}

// BenchmarkSearchFilter benchmarks different filter combinations
func BenchmarkSearchFilter(b *testing.B) {
	roms := generateSearchTestROMs(1000)
	searchManager := menu.NewSearchManager()

	filters := []menu.SearchFilter{
		{ByName: true, ByAuthor: false, ByDescription: false, ByTags: false},
		{ByName: true, ByAuthor: true, ByDescription: false, ByTags: false},
		{ByName: true, ByAuthor: true, ByDescription: true, ByTags: false},
		{ByName: true, ByAuthor: true, ByDescription: true, ByTags: true},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter := filters[i%len(filters)]
		searchManager.SetQuery("game")
		results := searchManager.Search(roms, filter)
		_ = results
	}
}
