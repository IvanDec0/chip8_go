package browser

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

// TextParser handles parsing of ROM metadata text files
type TextParser struct {
	patterns map[string]*regexp.Regexp
}

// ParsedFields represents extracted metadata fields
type ParsedFields struct {
	Title       string
	Author      string
	Description []string
	Controls    []string
	Year        string
	System      string
}

// NewTextParser creates a new text parser
func NewTextParser() *TextParser {
	patterns := map[string]*regexp.Regexp{
		"title":       regexp.MustCompile(`(?i)^(?:title|name|game):\s*(.+)$`),
		"author":      regexp.MustCompile(`(?i)^(?:author|by|programmer|made by):\s*(.+)$`),
		"year":        regexp.MustCompile(`(?i)^(?:year|date|copyright):\s*(.+)$`),
		"system":      regexp.MustCompile(`(?i)^(?:system|platform|for):\s*(.+)$`),
		"controls":    regexp.MustCompile(`(?i)^(?:controls|keys|input):\s*(.+)$`),
		"description": regexp.MustCompile(`(?i)^(?:description|about|info):\s*(.+)$`),
	}

	return &TextParser{
		patterns: patterns,
	}
}

// ParseFile parses a text file and extracts metadata
func (tp *TextParser) ParseFile(filePath string) (*ParsedFields, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var content strings.Builder
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return tp.ParseContent(content.String()), nil
}

// ParseContent parses text content and extracts metadata
func (tp *TextParser) ParseContent(content string) *ParsedFields {
	lines := strings.Split(content, "\n")
	fields := &ParsedFields{
		Description: make([]string, 0),
		Controls:    make([]string, 0),
	}

	// First pass: extract structured data
	tp.extractStructuredFields(lines, fields)

	// Second pass: extract title from first line if not found
	if fields.Title == "" {
		fields.Title = tp.extractTitle(lines)
	}

	// Third pass: extract author from common patterns if not found
	if fields.Author == "" {
		fields.Author = tp.extractAuthor(lines)
	}

	// Fourth pass: extract description paragraphs
	if len(fields.Description) == 0 {
		fields.Description = tp.extractDescription(lines)
	}

	// Fifth pass: extract controls information
	if len(fields.Controls) == 0 {
		fields.Controls = tp.extractControls(lines)
	}

	return fields
}

// extractStructuredFields extracts fields using regex patterns
func (tp *TextParser) extractStructuredFields(lines []string, fields *ParsedFields) {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Try each pattern
		if match := tp.patterns["title"].FindStringSubmatch(line); match != nil {
			fields.Title = strings.TrimSpace(match[1])
		} else if match := tp.patterns["author"].FindStringSubmatch(line); match != nil {
			fields.Author = strings.TrimSpace(match[1])
		} else if match := tp.patterns["year"].FindStringSubmatch(line); match != nil {
			fields.Year = strings.TrimSpace(match[1])
		} else if match := tp.patterns["system"].FindStringSubmatch(line); match != nil {
			fields.System = strings.TrimSpace(match[1])
		} else if match := tp.patterns["controls"].FindStringSubmatch(line); match != nil {
			fields.Controls = append(fields.Controls, strings.TrimSpace(match[1]))
		} else if match := tp.patterns["description"].FindStringSubmatch(line); match != nil {
			fields.Description = append(fields.Description, strings.TrimSpace(match[1]))
		}
	}
}

// extractTitle extracts title from the first non-empty line
func (tp *TextParser) extractTitle(lines []string) string {
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// Skip obvious non-title lines
		if strings.Contains(strings.ToLower(line), "readme") ||
			strings.Contains(strings.ToLower(line), "instructions") ||
			strings.Contains(strings.ToLower(line), "controls") {
			continue
		}

		// Clean up common patterns
		title := line

		// Remove common prefixes
		title = strings.TrimPrefix(title, "Game: ")
		title = strings.TrimPrefix(title, "Title: ")
		title = strings.TrimPrefix(title, "Name: ")

		// Remove file extensions if somehow included
		if strings.Contains(title, ".ch8") {
			title = strings.Split(title, ".ch8")[0]
		}

		return strings.TrimSpace(title)
	}

	return ""
}

// extractAuthor extracts author information from various patterns
func (tp *TextParser) extractAuthor(lines []string) string {
	authorPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)by\s+([^,\n]+)`),
		regexp.MustCompile(`(?i)author:\s*([^\n]+)`),
		regexp.MustCompile(`(?i)programmer:\s*([^\n]+)`),
		regexp.MustCompile(`(?i)made by\s+([^\n]+)`),
		regexp.MustCompile(`(?i)written by:?\s*([^\n]+)`),
		regexp.MustCompile(`(?i)created by\s+([^\n]+)`),
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		for _, pattern := range authorPatterns {
			if match := pattern.FindStringSubmatch(line); match != nil {
				author := strings.TrimSpace(match[1])
				// Clean up common suffixes
				author = strings.TrimSuffix(author, ",")
				author = strings.TrimSuffix(author, ".")
				return author
			}
		}
	}

	return ""
}

// extractDescription extracts description paragraphs
func (tp *TextParser) extractDescription(lines []string) []string {
	var description []string
	inDescription := false

	skipPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)^(controls?|keys?|input):`),
		regexp.MustCompile(`(?i)^(running|installation|requirements):`),
		regexp.MustCompile(`(?i)^(credits?|thanks?|acknowledgments?):`),
		regexp.MustCompile(`(?i)^(distribution|license|copyright):`),
		regexp.MustCompile(`(?i)^\s*[-=_]{3,}`), // Separator lines
	}

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}

		// Skip title line (usually first meaningful line)
		if i < 3 && (strings.Contains(strings.ToLower(line), "by ") ||
			len(strings.Fields(line)) < 8) {
			continue
		}

		// Skip lines that match skip patterns
		shouldSkip := false
		for _, pattern := range skipPatterns {
			if pattern.MatchString(line) {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			inDescription = false
			continue
		}

		// Look for description indicators
		if strings.Contains(strings.ToLower(line), "description") ||
			strings.Contains(strings.ToLower(line), "about") ||
			strings.Contains(strings.ToLower(line), "game") {
			inDescription = true
			// Extract description from same line if present
			if strings.Contains(line, ":") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
					description = append(description, strings.TrimSpace(parts[1]))
				}
			}
			continue
		}

		// If we're in description mode or line looks like description
		if inDescription || (len(strings.Fields(line)) > 3 &&
			!strings.Contains(line, ":") &&
			!strings.HasPrefix(line, "-") &&
			!strings.HasPrefix(line, "*")) {

			description = append(description, line)
		}
	}

	return description
}

// extractControls extracts control information
func (tp *TextParser) extractControls(lines []string) []string {
	var controls []string
	inControls := false

	controlPatterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)(controls?|keys?|input|how to play):`),
		regexp.MustCompile(`(?i)(use|press|key|button)`),
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if this line indicates start of controls section
		for _, pattern := range controlPatterns {
			if pattern.MatchString(line) {
				inControls = true
				// Extract control info from same line if present
				if strings.Contains(line, ":") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
						controls = append(controls, strings.TrimSpace(parts[1]))
					}
				}
				break
			}
		}

		// If we're in controls section, collect relevant lines
		if inControls {
			// Stop at next section
			if strings.Contains(strings.ToLower(line), "credits") ||
				strings.Contains(strings.ToLower(line), "distribution") ||
				strings.Contains(strings.ToLower(line), "notes") {
				break
			}

			// Add lines that look like control instructions
			if strings.Contains(line, "key") ||
				strings.Contains(line, "button") ||
				strings.Contains(line, "press") ||
				strings.Contains(line, "use") ||
				regexp.MustCompile(`[0-9A-F]`).MatchString(line) {
				controls = append(controls, line)
			}
		}
	}

	return controls
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
