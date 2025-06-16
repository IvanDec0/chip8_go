package benchmarks

import (
	"chip8/internal/browser"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// BenchmarkBrowserScanDirectory benchmarks directory scanning performance
func BenchmarkBrowserScanDirectory(b *testing.B) {
	// Create temporary test directory with ROM files
	tmpDir, err := os.MkdirTemp("", "chip8_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test ROM files
	romCount := 100
	for i := 0; i < romCount; i++ {
		romPath := filepath.Join(tmpDir, fmt.Sprintf("rom_%03d.ch8", i))
		err = os.WriteFile(romPath, []byte("FAKE_ROM_DATA"), 0644)
		if err != nil {
			b.Fatalf("Failed to create test ROM %d: %v", i, err)
		}
	}

	// Create file browser
	fileBrowser := browser.NewFileBrowser(tmpDir)

	// Reset timer and run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := fileBrowser.ScanDirectory(tmpDir)
		if err != nil {
			b.Fatalf("Failed to scan directory: %v", err)
		}
	}
}

// BenchmarkBrowserScanLargeDirectory benchmarks large directory scanning
func BenchmarkBrowserScanLargeDirectory(b *testing.B) {
	if testing.Short() {
		b.Skip("Skipping large directory benchmark in short mode")
	}

	tmpDir, err := os.MkdirTemp("", "chip8_large_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create large number of ROM files
	romCount := 1000
	for i := 0; i < romCount; i++ {
		romPath := filepath.Join(tmpDir, fmt.Sprintf("rom_%04d.ch8", i))
		err = os.WriteFile(romPath, []byte("FAKE_ROM_DATA_CONTENT"), 0644)
		if err != nil {
			b.Fatalf("Failed to create test ROM %d: %v", i, err)
		}
	}

	fileBrowser := browser.NewFileBrowser(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		roms, err := fileBrowser.ScanDirectory(tmpDir)
		if err != nil {
			b.Fatalf("Failed to scan directory: %v", err)
		}

		// Verify we got expected number of results
		if len(roms) < romCount {
			b.Errorf("Expected at least %d ROMs, got %d", romCount, len(roms))
		}
	}
}

// BenchmarkBrowserWithMetadata benchmarks directory scanning with metadata
func BenchmarkBrowserWithMetadata(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "chip8_metadata_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create ROM files with metadata
	romCount := 50
	for i := 0; i < romCount; i++ {
		romPath := filepath.Join(tmpDir, fmt.Sprintf("game_%02d.ch8", i))
		metadataPath := filepath.Join(tmpDir, fmt.Sprintf("game_%02d.txt", i))

		// Create ROM file
		err = os.WriteFile(romPath, []byte("FAKE_ROM_DATA"), 0644)
		if err != nil {
			b.Fatalf("Failed to create ROM file: %v", err)
		}

		// Create metadata file
		metadata := fmt.Sprintf(`Title: Test Game %d
Author: Test Author
Year: 202%d
Description: A test game for benchmarking
Genre: Test
`, i, i%10)

		err = os.WriteFile(metadataPath, []byte(metadata), 0644)
		if err != nil {
			b.Fatalf("Failed to create metadata file: %v", err)
		}
	}

	fileBrowser := browser.NewFileBrowser(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Scan with metadata extraction
		roms, err := fileBrowser.ScanDirectoryWithMetadata(tmpDir)
		if err != nil {
			b.Fatalf("Failed to scan directory with metadata: %v", err)
		}

		// Verify metadata was loaded (simplified check)
		metadataCount := 0
		for _, rom := range roms {
			if !rom.IsDirectory && filepath.Ext(rom.Name) == ".ch8" {
				metadataCount++
			}
		}

		if metadataCount < romCount/2 { // At least half should be ROM files
			b.Errorf("Expected more ROM files, got %d", metadataCount)
		}
	}
}

// BenchmarkBrowserCachePerformance benchmarks cache performance
func BenchmarkBrowserCachePerformance(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "chip8_cache_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test ROM files
	romCount := 200
	for i := 0; i < romCount; i++ {
		romPath := filepath.Join(tmpDir, fmt.Sprintf("cached_rom_%03d.ch8", i))
		err = os.WriteFile(romPath, []byte("CACHED_ROM_DATA"), 0644)
		if err != nil {
			b.Fatalf("Failed to create test ROM: %v", err)
		}
	}

	fileBrowser := browser.NewFileBrowser(tmpDir)

	// Prime the cache with initial scan
	_, err = fileBrowser.ScanDirectory(tmpDir)
	if err != nil {
		b.Fatalf("Failed to prime cache: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This should hit the cache
		_, err := fileBrowser.ScanDirectory(tmpDir)
		if err != nil {
			b.Fatalf("Failed to scan cached directory: %v", err)
		}
	}
}

// BenchmarkBrowserMemoryUsage benchmarks memory usage patterns
func BenchmarkBrowserMemoryUsage(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "chip8_memory_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create varying sizes of ROM files
	fileSizes := []int{1024, 4096, 8192, 16384, 32768}
	for i, size := range fileSizes {
		data := make([]byte, size)
		for j := range data {
			data[j] = byte(j % 256)
		}

		romPath := filepath.Join(tmpDir, fmt.Sprintf("memory_test_%d.ch8", i))
		err = os.WriteFile(romPath, data, 0644)
		if err != nil {
			b.Fatalf("Failed to create ROM file: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fileBrowser := browser.NewFileBrowser(tmpDir)

		// Perform multiple operations to test memory management
		for j := 0; j < 10; j++ {
			roms, err := fileBrowser.ScanDirectory(tmpDir)
			if err != nil {
				b.Fatalf("Failed to scan directory: %v", err)
			}

			// Access ROM information to trigger memory allocation
			for _, rom := range roms {
				_ = rom.Name
				_ = rom.Size
				_ = rom.ModTime
			}
		}
	}
}

// BenchmarkBrowserConcurrentAccess benchmarks concurrent directory access
func BenchmarkBrowserConcurrentAccess(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "chip8_concurrent_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test ROM files
	romCount := 100
	for i := 0; i < romCount; i++ {
		romPath := filepath.Join(tmpDir, fmt.Sprintf("concurrent_rom_%03d.ch8", i))
		err = os.WriteFile(romPath, []byte("CONCURRENT_TEST_DATA"), 0644)
		if err != nil {
			b.Fatalf("Failed to create ROM file: %v", err)
		}
	}

	fileBrowser := browser.NewFileBrowser(tmpDir)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := fileBrowser.ScanDirectory(tmpDir)
			if err != nil {
				b.Errorf("Failed to scan directory in parallel: %v", err)
			}
		}
	})
}

// BenchmarkBrowserNavigationOperations benchmarks navigation operations
func BenchmarkBrowserNavigationOperations(b *testing.B) {
	// Create nested directory structure
	tmpDir, err := os.MkdirTemp("", "chip8_nav_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create nested directories with ROMs
	dirs := []string{"games", "demos", "homebrew"}
	for _, dir := range dirs {
		dirPath := filepath.Join(tmpDir, dir)
		err = os.MkdirAll(dirPath, 0755)
		if err != nil {
			b.Fatalf("Failed to create directory: %v", err)
		}

		// Add ROMs to each directory
		for i := 0; i < 20; i++ {
			romPath := filepath.Join(dirPath, fmt.Sprintf("%s_rom_%02d.ch8", dir, i))
			err = os.WriteFile(romPath, []byte("NAV_TEST_ROM"), 0644)
			if err != nil {
				b.Fatalf("Failed to create ROM file: %v", err)
			}
		}
	}

	fileBrowser := browser.NewFileBrowser(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Navigate through directories
		for _, dir := range dirs {
			dirPath := filepath.Join(tmpDir, dir)

			err := fileBrowser.SetCurrentDirectory(dirPath)
			if err != nil {
				b.Fatalf("Failed to navigate to directory: %v", err)
			}

			_, err = fileBrowser.ScanDirectory(dirPath)
			if err != nil {
				b.Fatalf("Failed to scan directory: %v", err)
			}

			// Navigate back
			err = fileBrowser.NavigateUp()
			if err != nil {
				b.Fatalf("Failed to navigate up: %v", err)
			}
		}
	}
}
