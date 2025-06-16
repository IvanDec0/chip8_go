package unit

import (
	"chip8/internal/browser"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMetadataExtractor_ParseTextFile(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "chip8_metadata_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test metadata file
	metadataContent := `15 Puzzle
By: Roger Ivie

This is a sliding puzzle game where you arrange numbered tiles.

Controls:
2, 4, 6, 8: Move tiles

Instructions:
Use the arrow keys to slide tiles into the empty space.
Try to arrange the numbers in order from 1 to 15.`

	metadataFile := filepath.Join(tmpDir, "15 Puzzle [Roger Ivie].txt")
	err = os.WriteFile(metadataFile, []byte(metadataContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create metadata file: %v", err)
	}

	// Test metadata extraction
	extractor := browser.NewMetadataExtractor()
	metadata, err := extractor.ExtractMetadata(metadataFile)
	if err != nil {
		t.Fatalf("Failed to extract metadata: %v", err)
	}

	// Verify extracted metadata
	if metadata.Title != "15 Puzzle" {
		t.Errorf("Expected title '15 Puzzle', got '%s'", metadata.Title)
	}

	if metadata.Author != "Roger Ivie" {
		t.Errorf("Expected author 'Roger Ivie', got '%s'", metadata.Author)
	}

	if len(metadata.Description) == 0 {
		t.Error("Expected description to be extracted")
	}

	if len(metadata.Controls) == 0 {
		t.Error("Expected controls to be extracted")
	}

	// Check if controls contain expected information
	controlsText := strings.Join(metadata.Controls, " ")
	if !strings.Contains(controlsText, "2, 4, 6, 8") {
		t.Errorf("Expected controls to contain '2, 4, 6, 8', got: %s", controlsText)
	}
}

func TestMetadataExtractor_ParseTextFileVariousFormats(t *testing.T) {
	testCases := []struct {
		name     string
		content  string
		expected struct {
			title  string
			author string
		}
	}{
		{
			name: "Standard Format",
			content: `Game Title
Author: John Doe

Game description here.`,
			expected: struct {
				title  string
				author string
			}{
				title:  "Game Title",
				author: "John Doe",
			},
		},
		{
			name: "By Format",
			content: `Amazing Game
By John Smith

This is an amazing game.`,
			expected: struct {
				title  string
				author string
			}{
				title:  "Amazing Game",
				author: "John Smith",
			},
		},
		{
			name: "Alternate Format",
			content: `Cool ROM
Written by: Jane Developer

Cool description.`,
			expected: struct {
				title  string
				author string
			}{
				title:  "Cool ROM",
				author: "Jane Developer",
			},
		},
	}

	tmpDir, err := os.MkdirTemp("", "chip8_metadata_formats_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	extractor := browser.NewMetadataExtractor()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test file
			filename := filepath.Join(tmpDir, tc.name+".txt")
			err := os.WriteFile(filename, []byte(tc.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Extract metadata
			metadata, err := extractor.ExtractMetadata(filename)
			if err != nil {
				t.Fatalf("Failed to extract metadata: %v", err)
			}

			// Verify results
			if metadata.Title != tc.expected.title {
				t.Errorf("Expected title '%s', got '%s'", tc.expected.title, metadata.Title)
			}

			if metadata.Author != tc.expected.author {
				t.Errorf("Expected author '%s', got '%s'", tc.expected.author, metadata.Author)
			}
		})
	}
}

func TestMetadataExtractor_FallbackGeneration(t *testing.T) {
	extractor := browser.NewMetadataExtractor()

	testCases := []struct {
		filename string
		expected struct {
			title  string
			author string
		}
	}{
		{
			filename: "Pong [David Winter].ch8",
			expected: struct {
				title  string
				author string
			}{
				title:  "Pong",
				author: "David Winter",
			},
		},
		{
			filename: "Simple Game.ch8",
			expected: struct {
				title  string
				author string
			}{
				title:  "Simple Game",
				author: "",
			},
		},
		{
			filename: "test-rom.c8",
			expected: struct {
				title  string
				author string
			}{
				title:  "test-rom",
				author: "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.filename, func(t *testing.T) {
			// Extract metadata which will use fallback for non-existent files
			metadata, err := extractor.ExtractMetadata(tc.filename)
			if err != nil {
				t.Fatalf("Failed to extract metadata: %v", err)
			}

			if metadata.Title != tc.expected.title {
				t.Errorf("Expected title '%s', got '%s'", tc.expected.title, metadata.Title)
			}

			if metadata.Author != tc.expected.author {
				t.Errorf("Expected author '%s', got '%s'", tc.expected.author, metadata.Author)
			}
		})
	}
}

func TestMetadataExtractor_Caching(t *testing.T) {
	// Create temporary directory and file
	tmpDir, err := os.MkdirTemp("", "chip8_metadata_cache_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	metadataFile := filepath.Join(tmpDir, "test.txt")
	content := "Test Game\nBy: Test Author\n\nTest description."
	err = os.WriteFile(metadataFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	extractor := browser.NewMetadataExtractor()

	// First extraction
	metadata1, err := extractor.ExtractMetadata(metadataFile)
	if err != nil {
		t.Fatalf("First extraction failed: %v", err)
	}

	// Second extraction (should use cache)
	metadata2, err := extractor.ExtractMetadata(metadataFile)
	if err != nil {
		t.Fatalf("Second extraction failed: %v", err)
	}

	// Results should be identical
	if metadata1.Title != metadata2.Title {
		t.Errorf("Cached metadata title mismatch: '%s' vs '%s'", metadata1.Title, metadata2.Title)
	}

	if metadata1.Author != metadata2.Author {
		t.Errorf("Cached metadata author mismatch: '%s' vs '%s'", metadata1.Author, metadata2.Author)
	}
}

func TestMetadataExtractor_InvalidFile(t *testing.T) {
	extractor := browser.NewMetadataExtractor()

	// Test with non-existent file - should return fallback metadata, not error
	metadata, err := extractor.ExtractMetadata("/nonexistent/file.txt")
	if err != nil {
		t.Errorf("ExtractMetadata should not return error for non-existent file, got: %v", err)
	}
	if metadata == nil {
		t.Error("Expected fallback metadata for non-existent file")
	}

	// Test with empty file
	tmpDir, err := os.MkdirTemp("", "chip8_metadata_invalid_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	emptyFile := filepath.Join(tmpDir, "empty.txt")
	err = os.WriteFile(emptyFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	metadata, err = extractor.ExtractMetadata(emptyFile)
	if err != nil {
		t.Fatalf("Failed to handle empty file: %v", err)
	}

	// Should still return some metadata even for empty file
	if metadata == nil {
		t.Error("Expected metadata to be returned even for empty file")
	}
}

func BenchmarkMetadataExtractor_ExtractFromTextFile(b *testing.B) {
	// Create temporary test file
	tmpDir, err := os.MkdirTemp("", "chip8_metadata_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	content := `Test Game
By: Test Author

This is a test game with multiple lines of description.
It has several paragraphs to test parsing performance.

Controls:
1,2,3,4: Player controls
Q,W,E,R: Menu controls

Instructions:
Use the controls to play the game.
Have fun!`

	testFile := filepath.Join(tmpDir, "test.txt")
	err = os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	extractor := browser.NewMetadataExtractor()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := extractor.ExtractMetadata(testFile)
		if err != nil {
			b.Fatalf("Extraction failed: %v", err)
		}
	}
}
