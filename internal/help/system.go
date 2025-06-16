package help

import (
	"fmt"
	"strings"
)

// HelpSystem provides context-sensitive help and documentation
type HelpSystem struct {
	content        map[string]*HelpTopic
	shortcuts      map[string]string
	currentContext string
}

// HelpTopic represents a help topic with content and related information
type HelpTopic struct {
	Title     string
	Content   string
	Shortcuts []ShortcutInfo
	SeeAlso   []string
	Category  string
	Context   string
}

// ShortcutInfo represents a keyboard shortcut and its description
type ShortcutInfo struct {
	Key         string
	Description string
	Context     string
}

// NewHelpSystem creates a new help system with built-in content
func NewHelpSystem() *HelpSystem {
	hs := &HelpSystem{
		content:   make(map[string]*HelpTopic),
		shortcuts: make(map[string]string),
	}

	hs.initializeContent()
	return hs
}

// GetTopic returns help content for a specific topic
func (hs *HelpSystem) GetTopic(topic string) (*HelpTopic, error) {
	if content, exists := hs.content[topic]; exists {
		return content, nil
	}
	return nil, fmt.Errorf("help topic '%s' not found", topic)
}

// GetContextualHelp returns help content for the current context
func (hs *HelpSystem) GetContextualHelp(context string) (*HelpTopic, error) {
	hs.currentContext = context

	// Look for context-specific help first
	contextTopic := fmt.Sprintf("%s_help", context)
	if content, exists := hs.content[contextTopic]; exists {
		return content, nil
	}

	// Fall back to general help
	return hs.GetTopic("general_help")
}

// ListTopics returns all available help topics
func (hs *HelpSystem) ListTopics() []string {
	topics := make([]string, 0, len(hs.content))
	for topic := range hs.content {
		topics = append(topics, topic)
	}
	return topics
}

// SearchHelp searches for help topics containing the query
func (hs *HelpSystem) SearchHelp(query string) []*HelpTopic {
	query = strings.ToLower(query)
	results := make([]*HelpTopic, 0)

	for _, topic := range hs.content {
		if strings.Contains(strings.ToLower(topic.Title), query) ||
			strings.Contains(strings.ToLower(topic.Content), query) {
			results = append(results, topic)
		}
	}

	return results
}

// GetShortcuts returns keyboard shortcuts for a specific context
func (hs *HelpSystem) GetShortcuts(context string) []ShortcutInfo {
	shortcuts := make([]ShortcutInfo, 0)

	for _, topic := range hs.content {
		if topic.Context == context || context == "" {
			shortcuts = append(shortcuts, topic.Shortcuts...)
		}
	}

	return shortcuts
}

// FormatHelp formats help content for display
func (hs *HelpSystem) FormatHelp(topic *HelpTopic) string {
	var builder strings.Builder

	// Title
	builder.WriteString(fmt.Sprintf("=== %s ===\n\n", topic.Title))

	// Content
	builder.WriteString(topic.Content)
	builder.WriteString("\n\n")

	// Shortcuts
	if len(topic.Shortcuts) > 0 {
		builder.WriteString("KEYBOARD SHORTCUTS:\n")
		for _, shortcut := range topic.Shortcuts {
			builder.WriteString(fmt.Sprintf("  %-12s - %s\n", shortcut.Key, shortcut.Description))
		}
		builder.WriteString("\n")
	}

	// See Also
	if len(topic.SeeAlso) > 0 {
		builder.WriteString("SEE ALSO:\n")
		for _, related := range topic.SeeAlso {
			builder.WriteString(fmt.Sprintf("  - %s\n", related))
		}
	}

	return builder.String()
}

// initializeContent sets up the built-in help content
func (hs *HelpSystem) initializeContent() {
	// General help
	hs.content["general_help"] = &HelpTopic{
		Title:   "CHIP-8 Emulator Help",
		Context: "general",
		Content: `Welcome to the CHIP-8 Emulator with Advanced ROM Browser!

This emulator provides a full-featured CHIP-8 gaming experience with an enhanced
ROM browser, metadata support, and comprehensive configuration options.

GETTING STARTED:
1. Navigate the ROM browser using arrow keys
2. Press ENTER to launch a ROM or enter a directory
3. Use F1 for context-sensitive help
4. Press F2 to access settings

MAIN FEATURES:
- Advanced ROM browser with search and sorting
- Metadata extraction and display
- Favorites and recent ROMs tracking
- Customizable themes and key bindings
- Performance optimization for large collections`,
		Shortcuts: []ShortcutInfo{
			{Key: "F1", Description: "Show help", Context: "general"},
			{Key: "F2", Description: "Settings", Context: "general"},
			{Key: "F3", Description: "Favorites", Context: "general"},
			{Key: "F4", Description: "Recent ROMs", Context: "general"},
			{Key: "F5", Description: "Refresh", Context: "general"},
			{Key: "Ctrl+Q", Description: "Quit", Context: "general"},
		},
		SeeAlso:  []string{"browser_help", "emulator_help", "settings_help"},
		Category: "General",
	}

	// ROM Browser help
	hs.content["browser_help"] = &HelpTopic{
		Title:   "ROM Browser Help",
		Context: "browser",
		Content: `The ROM browser allows you to navigate and organize your CHIP-8 ROM collection.

NAVIGATION:
- Use arrow keys (↑/↓) to move through the ROM list
- Press Page Up/Down to jump by screen pages
- Press Home/End to go to first/last item
- Press ENTER to launch a ROM or enter a directory
- Press Backspace to go to parent directory

SEARCH:
- Start typing to search for ROMs by name
- Search is case-insensitive and supports partial matches
- Press ESC to clear the search
- Search includes metadata (author, description)

SORTING:
- Press S + N to sort by name
- Press S + D to sort by date
- Press S + Z to sort by size
- Press SPACE to toggle sort direction
- Directories always appear first

FEATURES:
- Automatic metadata extraction from .txt files
- Caching for improved performance
- Recent ROMs tracking
- Favorites management`,
		Shortcuts: []ShortcutInfo{
			{Key: "↑/↓", Description: "Navigate list", Context: "browser"},
			{Key: "Page Up/Down", Description: "Jump pages", Context: "browser"},
			{Key: "Home/End", Description: "First/last item", Context: "browser"},
			{Key: "ENTER", Description: "Select/launch", Context: "browser"},
			{Key: "Backspace", Description: "Parent directory", Context: "browser"},
			{Key: "ESC", Description: "Clear search", Context: "browser"},
			{Key: "F", Description: "Add to favorites", Context: "browser"},
			{Key: "S+N", Description: "Sort by name", Context: "browser"},
			{Key: "S+D", Description: "Sort by date", Context: "browser"},
			{Key: "S+Z", Description: "Sort by size", Context: "browser"},
			{Key: "SPACE", Description: "Toggle sort direction", Context: "browser"},
		},
		SeeAlso:  []string{"emulator_help", "settings_help", "metadata_help"},
		Category: "Browser",
	}

	// Emulator help
	hs.content["emulator_help"] = &HelpTopic{
		Title:   "CHIP-8 Emulator Help",
		Context: "emulator",
		Content: `The CHIP-8 emulator provides accurate emulation of the CHIP-8 system.

CHIP-8 KEYPAD LAYOUT:
The CHIP-8 system uses a 16-key hexadecimal keypad (0-F).
Default keyboard mapping:

CHIP-8 Keypad    Your Keyboard
┌─┬─┬─┬─┐        ┌─┬─┬─┬─┐
│1│2│3│C│   →    │1│2│3│4│
├─┼─┼─┼─┤        ├─┼─┼─┼─┤
│4│5│6│D│        │Q│W│E│R│
├─┼─┼─┼─┤        ├─┼─┼─┼─┤
│7│8│9│E│        │A│S│D│F│
├─┼─┼─┼─┤        ├─┼─┼─┼─┤
│A│0│B│F│        │Z│X│C│V│
└─┴─┴─┴─┘        └─┴─┴─┴─┘

EMULATOR CONTROLS:
- ESC: Return to ROM browser
- F11: Toggle fullscreen
- Ctrl+R: Reset game
- Ctrl+P: Pause/resume

GAME TIPS:
- Each game uses different keys - check the game's documentation
- Some games require specific timing - adjust clock speed in settings
- If a game runs too fast/slow, modify the clock speed setting`,
		Shortcuts: []ShortcutInfo{
			{Key: "1,2,3,4", Description: "CHIP-8 keys 1,2,3,C", Context: "emulator"},
			{Key: "Q,W,E,R", Description: "CHIP-8 keys 4,5,6,D", Context: "emulator"},
			{Key: "A,S,D,F", Description: "CHIP-8 keys 7,8,9,E", Context: "emulator"},
			{Key: "Z,X,C,V", Description: "CHIP-8 keys A,0,B,F", Context: "emulator"},
			{Key: "ESC", Description: "Return to browser", Context: "emulator"},
			{Key: "F11", Description: "Toggle fullscreen", Context: "emulator"},
			{Key: "Ctrl+R", Description: "Reset game", Context: "emulator"},
			{Key: "Ctrl+P", Description: "Pause/resume", Context: "emulator"},
		},
		SeeAlso:  []string{"browser_help", "settings_help"},
		Category: "Emulator",
	}

	// Settings help
	hs.content["settings_help"] = &HelpTopic{
		Title:   "Settings and Configuration",
		Context: "settings",
		Content: `The settings system allows you to customize the emulator and browser.

ACCESSING SETTINGS:
- Press F2 from any screen to open settings
- Use arrow keys to navigate options
- Press ENTER to modify a setting
- Changes are saved automatically

MAIN SETTINGS CATEGORIES:

DISPLAY SETTINGS:
- Theme: Visual appearance (default, dark, light, classic)
- Scale: Size of the emulator display
- Fullscreen: Start games in fullscreen mode
- Color Palette: Colors for CHIP-8 graphics

BROWSER SETTINGS:
- Default ROM Directory: Where to look for ROMs
- Sort Order: How to sort ROM lists
- Page Size: Number of ROMs per screen
- Show Hidden Files: Display hidden files and directories

PERFORMANCE SETTINGS:
- Clock Speed: CHIP-8 CPU speed (affects game speed)
- Cache Settings: Performance optimization options
- Memory Limits: Control memory usage

AUDIO SETTINGS:
- Sound Enable: Turn sound on/off
- Volume: Audio volume level

The configuration file is stored at ~/.chip8/config.json and can be
manually edited for advanced customization.`,
		Shortcuts: []ShortcutInfo{
			{Key: "F2", Description: "Open settings", Context: "settings"},
			{Key: "↑/↓", Description: "Navigate options", Context: "settings"},
			{Key: "ENTER", Description: "Modify setting", Context: "settings"},
			{Key: "ESC", Description: "Close settings", Context: "settings"},
			{Key: "Tab", Description: "Switch categories", Context: "settings"},
		},
		SeeAlso:  []string{"general_help", "browser_help"},
		Category: "Configuration",
	}

	// Metadata help
	hs.content["metadata_help"] = &HelpTopic{
		Title:   "ROM Metadata System",
		Context: "metadata",
		Content: `The metadata system provides detailed information about your ROMs.

AUTOMATIC METADATA:
The browser automatically extracts metadata from:
- Filename patterns (e.g., "Pong [David Winter, 1990].ch8")
- Accompanying .txt files with the same name as the ROM

CREATING METADATA FILES:
Create a .txt file with the same name as your ROM file:

Example: For "Pong.ch8", create "Pong.txt" with content:
  Title: Pong
  Author: Paul Vervalin  
  Year: 1990
  Description: Classic paddle and ball game
  Genre: Arcade
  Instructions: Use 1 and Q for left paddle, 4 and R for right

SUPPORTED METADATA FIELDS:
- Title: Game name
- Author: Creator/developer
- Year: Release year
- Description: Game description
- Genre: Game category
- Instructions: How to play
- Rating: Game rating
- Notes: Additional information

BENEFITS:
- Enhanced ROM information display
- Better search capabilities
- Organized ROM collections
- Historical preservation`,
		Shortcuts: []ShortcutInfo{
			{Key: "Ctrl+I", Description: "Show ROM info", Context: "metadata"},
		},
		SeeAlso:  []string{"browser_help"},
		Category: "Metadata",
	}

	// Troubleshooting help
	hs.content["troubleshooting_help"] = &HelpTopic{
		Title:   "Troubleshooting",
		Context: "troubleshooting",
		Content: `Common issues and solutions for the CHIP-8 emulator.

COMMON PROBLEMS:

NO ROMS FOUND:
- Check that your ROM directory contains .ch8 files
- Verify the ROM directory path in settings (F2)
- Ensure you have read permissions on the directory
- Press F5 to refresh the directory scan

GAME RUNS TOO FAST/SLOW:
- Adjust the clock speed in settings (F2)
- Try values between 500-1000 Hz
- Some games require specific speeds

NO SOUND:
- Check sound settings (F2)
- Verify system audio is working
- Try toggling sound on/off

POOR PERFORMANCE:
- Reduce cache size in settings
- Disable background scanning
- Organize ROMs into smaller directories
- Close other applications

CONTROLS NOT WORKING:
- Check key bindings in settings
- Ensure correct CHIP-8 key mapping
- Try resetting to default controls

CRASHES OR ERRORS:
- Check application logs in ~/.chip8/logs/
- Try deleting config file to reset settings
- Verify ROM files are not corrupted

For additional help, check the documentation files or report
issues on the project's GitHub repository.`,
		Shortcuts: []ShortcutInfo{
			{Key: "F5", Description: "Refresh/reload", Context: "troubleshooting"},
		},
		SeeAlso:  []string{"settings_help", "general_help"},
		Category: "Support",
	}

	// Performance help
	hs.content["performance_help"] = &HelpTopic{
		Title:   "Performance Optimization",
		Context: "performance",
		Content: `Tips and settings for optimal performance with large ROM collections.

PERFORMANCE SETTINGS:

CACHE OPTIMIZATION:
- Enable caching for faster directory scans
- Adjust cache size based on available memory
- Set appropriate cache expiration times

DIRECTORY ORGANIZATION:
- Keep directories under 100 ROMs for best performance
- Use subdirectories to organize collections
- Remove duplicate or unwanted ROMs

DISPLAY SETTINGS:
- Reduce page size for faster rendering
- Disable animations on slower systems
- Use simpler themes

MEMORY MANAGEMENT:
- Set memory limits to prevent excessive usage
- Disable metadata for very large collections
- Clear cache periodically by restarting

SYSTEM REQUIREMENTS:
- Minimum: 512MB RAM, 50MB disk space
- Recommended: 1GB RAM, 100MB disk space
- SSD storage improves scan performance

BENCHMARKS:
- Directory scan: <100ms for 1000 ROMs
- Search response: <50ms for real-time search  
- Memory usage: <50MB additional overhead
- Startup time: <500ms from launch to menu`,
		SeeAlso:  []string{"settings_help", "troubleshooting_help"},
		Category: "Performance",
	}
}

// GetHelpIndex returns a formatted index of all available help topics
func (hs *HelpSystem) GetHelpIndex() string {
	var builder strings.Builder

	builder.WriteString("=== HELP INDEX ===\n\n")

	// Group by category
	categories := make(map[string][]string)
	for topicName, topic := range hs.content {
		category := topic.Category
		if category == "" {
			category = "Other"
		}
		categories[category] = append(categories[category], topicName)
	}

	// Display by category
	categoryOrder := []string{"General", "Browser", "Emulator", "Configuration", "Metadata", "Performance", "Support", "Other"}

	for _, category := range categoryOrder {
		if topics, exists := categories[category]; exists {
			builder.WriteString(fmt.Sprintf("%s:\n", category))
			for _, topicName := range topics {
				topic := hs.content[topicName]
				builder.WriteString(fmt.Sprintf("  %-20s - %s\n", topicName, topic.Title))
			}
			builder.WriteString("\n")
		}
	}

	builder.WriteString("Use 'help <topic>' for detailed information on any topic.\n")
	builder.WriteString("Press F1 for context-sensitive help.\n")

	return builder.String()
}

// GetQuickReference returns a quick reference card
func (hs *HelpSystem) GetQuickReference() string {
	return `=== QUICK REFERENCE ===

NAVIGATION:
↑/↓         Navigate ROM list
Page Up/Down Jump by pages
ENTER       Launch ROM/enter directory
Backspace   Parent directory
ESC         Clear search/back

SEARCH & SORT:
Type        Start search
S+N         Sort by name
S+D         Sort by date  
S+Z         Sort by size
SPACE       Toggle sort direction

FEATURES:
F           Add to favorites
F1          Help
F2          Settings
F3          Favorites
F4          Recent ROMs
F5          Refresh

EMULATOR:
1234/QWER/ASDF/ZXCV  CHIP-8 keypad
ESC         Return to browser
F11         Toggle fullscreen
Ctrl+R      Reset game
Ctrl+P      Pause/resume

For detailed help, press F1 or type 'help' followed by a topic name.`
}
