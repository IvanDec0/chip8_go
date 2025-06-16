package unit

import (
	"chip8/internal/browser"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileBrowser_ScanDirectory(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "chip8_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	testFiles := []string{
		"game1.ch8",
		"game2.c8",
		"readme.txt",
		"invalid.exe",
	}

	for _, file := range testFiles {
		path := filepath.Join(tmpDir, file)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Test browser scanning
	fb := browser.NewFileBrowser(tmpDir)
	roms, err := fb.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	// Verify results
	expectedROMs := 2 // game1.ch8, game2.c8
	expectedDirs := 2 // subdir + ".." parent directory

	romCount := 0
	dirCount := 0
	for _, rom := range roms {
		if rom.IsDirectory {
			dirCount++
		} else {
			romCount++
		}
	}

	if romCount != expectedROMs {
		t.Errorf("Expected %d ROM files, got %d", expectedROMs, romCount)
	}
	if dirCount != expectedDirs {
		t.Errorf("Expected %d directories, got %d", expectedDirs, dirCount)
	}
}

func TestFileBrowser_SetCurrentDirectory(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "chip8_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fb := browser.NewFileBrowser(tmpDir)

	// Test setting valid directory
	err = fb.SetCurrentDirectory(tmpDir)
	if err != nil {
		t.Errorf("SetCurrentDirectory failed for valid directory: %v", err)
	}

	currentDir := fb.GetCurrentDirectory()
	if currentDir != tmpDir {
		t.Errorf("Expected current directory %s, got %s", tmpDir, currentDir)
	}

	// Test setting invalid directory
	invalidDir := filepath.Join(tmpDir, "nonexistent")
	err = fb.SetCurrentDirectory(invalidDir)
	if err == nil {
		t.Error("Expected error for nonexistent directory, got nil")
	}
}

func TestFileBrowser_NavigateUp(t *testing.T) {
	// Create temporary test directory structure
	tmpDir, err := os.MkdirTemp("", "chip8_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	fb := browser.NewFileBrowser(subDir)

	// Navigate up from subdirectory
	err = fb.NavigateUp()
	if err != nil {
		t.Errorf("NavigateUp failed: %v", err)
	}

	currentDir := fb.GetCurrentDirectory()
	if currentDir != tmpDir {
		t.Errorf("Expected to navigate to parent directory %s, got %s", tmpDir, currentDir)
	}
}

func TestFileBrowser_Caching(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "chip8_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	testFile := filepath.Join(tmpDir, "test.ch8")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	fb := browser.NewFileBrowser(tmpDir)

	// First scan should cache results
	start := time.Now()
	roms1, err := fb.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("First scan failed: %v", err)
	}
	firstScanTime := time.Since(start)

	// Second scan should be faster due to caching
	start = time.Now()
	roms2, err := fb.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("Second scan failed: %v", err)
	}
	secondScanTime := time.Since(start)

	// Verify same results
	if len(roms1) != len(roms2) {
		t.Errorf("Cache returned different results: %d vs %d", len(roms1), len(roms2))
	}

	// Second scan should be significantly faster (cached)
	if secondScanTime >= firstScanTime {
		t.Logf("Warning: Second scan (%v) not faster than first (%v) - caching may not be working", secondScanTime, firstScanTime)
	}
}

func TestFileBrowser_EmptyDirectory(t *testing.T) {
	// Create empty temporary directory
	tmpDir, err := os.MkdirTemp("", "chip8_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	fb := browser.NewFileBrowser(tmpDir)
	roms, err := fb.ScanDirectory(tmpDir)
	if err != nil {
		t.Fatalf("ScanDirectory failed for empty directory: %v", err)
	}

	if len(roms) != 1 || !roms[0].IsDirectory || roms[0].Name != ".." {
		t.Errorf("Expected 1 directory entry (..) in empty directory, got %d entries", len(roms))
	}
}

func TestFileBrowser_LargeCollection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large collection test in short mode")
	}

	// Create temporary directory with many files
	tmpDir, err := os.MkdirTemp("", "chip8_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create 100 test ROM files
	for i := 0; i < 100; i++ {
		filename := filepath.Join(tmpDir, fmt.Sprintf("game%03d.ch8", i))
		if err := os.WriteFile(filename, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	fb := browser.NewFileBrowser(tmpDir)

	// Benchmark directory scanning
	start := time.Now()
	roms, err := fb.ScanDirectory(tmpDir)
	scanTime := time.Since(start)

	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if len(roms) != 101 { // 100 ROMs + 1 parent directory
		t.Errorf("Expected 101 entries (100 ROMs + parent dir), got %d", len(roms))
	}

	// Count actual ROM files (excluding "..")
	romCount := 0
	for _, rom := range roms {
		if !rom.IsDirectory {
			romCount++
		}
	}
	if romCount != 100 {
		t.Errorf("Expected 100 ROM files, got %d", romCount)
	}

	// Performance target: <100ms for 100 ROMs
	if scanTime > 100*time.Millisecond {
		t.Errorf("Directory scan took too long: %v (target: <100ms)", scanTime)
	}

	t.Logf("Scanned %d ROMs in %v", len(roms), scanTime)
}

func BenchmarkFileBrowser_ScanDirectory(b *testing.B) {
	// Create temporary directory with test files
	tmpDir, err := os.MkdirTemp("", "chip8_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test ROM files
	for i := 0; i < 50; i++ {
		filename := filepath.Join(tmpDir, fmt.Sprintf("game%03d.ch8", i))
		if err := os.WriteFile(filename, []byte("test"), 0644); err != nil {
			b.Fatalf("Failed to create test file: %v", err)
		}
	}

	fb := browser.NewFileBrowser(tmpDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := fb.ScanDirectory(tmpDir)
		if err != nil {
			b.Fatalf("ScanDirectory failed: %v", err)
		}
	}
}

// Additional performance and edge case tests for FileBrowser
func TestFileBrowser_PerformanceWithManyFiles(t *testing.T) {
	testDir, err := os.MkdirTemp("", "chip8_perf_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	// Create many test files
	numFiles := 1000
	for i := 0; i < numFiles; i++ {
		filename := filepath.Join(testDir, fmt.Sprintf("rom%04d.ch8", i))
		if err := os.WriteFile(filename, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	fb := browser.NewFileBrowser(testDir)

	start := time.Now()
	roms, err := fb.ScanDirectory(testDir)
	duration := time.Since(start)

	if err != nil {
		t.Fatal(err)
	}

	if len(roms) != numFiles+1 { // +1 for parent directory
		t.Errorf("Expected %d entries (%d ROMs + parent dir), got %d", numFiles+1, numFiles, len(roms))
	}

	// Count actual ROM files (excluding "..")
	romCount := 0
	for _, rom := range roms {
		if !rom.IsDirectory {
			romCount++
		}
	}
	if romCount != numFiles {
		t.Errorf("Expected %d ROM files, got %d", numFiles, romCount)
	}

	// Performance target: < 100ms for 1000 ROMs
	if duration > 100*time.Millisecond {
		t.Errorf("Scan took too long: %v (target: <100ms)", duration)
	}

	t.Logf("Scanned %d ROMs in %v", numFiles, duration)
}

func TestFileBrowser_NonExistentDirectory(t *testing.T) {
	nonExistentDir := "/non/existent/directory/chip8/test"
	fb := browser.NewFileBrowser(nonExistentDir)

	_, err := fb.ScanDirectory(nonExistentDir)
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestFileBrowser_CachePerformance(t *testing.T) {
	testDir, err := os.MkdirTemp("", "chip8_cache_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	// Create test files
	for i := 0; i < 100; i++ {
		filename := filepath.Join(testDir, fmt.Sprintf("rom%04d.ch8", i))
		if err := os.WriteFile(filename, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	fb := browser.NewFileBrowser(testDir)

	// First scan (cache miss)
	start1 := time.Now()
	_, err = fb.ScanDirectory(testDir)
	duration1 := time.Since(start1)
	if err != nil {
		t.Fatal(err)
	}

	// Second scan (should use cache)
	start2 := time.Now()
	_, err = fb.ScanDirectory(testDir)
	duration2 := time.Since(start2)
	if err != nil {
		t.Fatal(err)
	}

	// Cache should make second scan faster
	if duration2 > duration1 {
		t.Logf("Warning: Second scan (%v) was not faster than first (%v)", duration2, duration1)
	}

	t.Logf("First scan: %v, Second scan: %v", duration1, duration2)
}

func TestFileBrowser_NavigationEdgeCases(t *testing.T) {
	// Test navigation at root level
	fb := browser.NewFileBrowser("/")
	err := fb.NavigateUp()
	if err != nil {
		t.Errorf("NavigateUp from root should not return error, got: %v", err)
	}

	// Test setting invalid directory
	err = fb.SetCurrentDirectory("/invalid/path/that/does/not/exist")
	if err == nil {
		t.Error("Expected error when setting invalid directory")
	}
}

func TestFileBrowser_FileExtensionHandling(t *testing.T) {
	testDir, err := os.MkdirTemp("", "chip8_ext_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(testDir)

	// Create files with various extensions
	testFiles := map[string]bool{
		"game.ch8": true,  // Valid
		"game.c8":  true,  // Valid
		"game.CH8": true,  // Valid (case insensitive)
		"game.txt": false, // Invalid
		"game.bin": false, // Invalid
		"game":     false, // No extension
	}

	for filename := range testFiles {
		filePath := filepath.Join(testDir, filename)
		if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	fb := browser.NewFileBrowser(testDir)
	roms, err := fb.ScanDirectory(testDir)
	if err != nil {
		t.Fatal(err)
	}

	// Count valid ROM files found (excluding directories)
	romFileCount := 0
	for _, rom := range roms {
		if !rom.IsDirectory {
			romFileCount++
		}
	}

	validCount := 0
	for _, shouldBeValid := range testFiles {
		if shouldBeValid {
			validCount++
		}
	}

	if romFileCount != validCount {
		t.Errorf("Expected %d valid ROM files, got %d", validCount, romFileCount)
	}
}
