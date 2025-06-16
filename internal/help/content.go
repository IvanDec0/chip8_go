package help

import (
	"strings"
)

// Content provides additional help content and utilities
type Content struct {
	topics map[string]string
	tips   []string
}

// GetTutorialContent returns step-by-step tutorials
func GetTutorialContent() map[string]string {
	return map[string]string{
		"first_time_setup": `=== FIRST TIME SETUP TUTORIAL ===

Welcome to the CHIP-8 emulator! Let's get you started in just a few steps.

STEP 1: Organize Your ROMs
1. Create a folder for your CHIP-8 ROMs (e.g., ~/chip8-roms/)
2. Place your .ch8 ROM files in this folder
3. Optionally, create subfolders (games/, demos/, etc.)

STEP 2: Configure the ROM Directory
1. Start the emulator
2. Press F2 to open settings
3. Navigate to "Default ROM Directory"
4. Press ENTER and select your ROM folder
5. Press ESC to save and return

STEP 3: Try Your First Game
1. Navigate the ROM list with ↑/↓ arrows
2. Select a ROM and press ENTER
3. Use the keyboard as the CHIP-8 keypad:
   - 1,2,3,4 for top row
   - Q,W,E,R for second row
   - A,S,D,F for third row
   - Z,X,C,V for bottom row

STEP 4: Explore Features
1. Press F3 to see favorites
2. Press F4 to see recent ROMs
3. Type to search for specific games
4. Press S+N to sort by name

You're all set! Press F1 anytime for help.`,

		"rom_organization": `=== ROM ORGANIZATION TUTORIAL ===

Learn how to organize your ROM collection for the best experience.

DIRECTORY STRUCTURE:
Organize your ROMs like this:
  chip8-roms/
  ├── games/
  │   ├── arcade/
  │   ├── puzzle/
  │   └── action/
  ├── demos/
  └── homebrew/

NAMING CONVENTIONS:
Use descriptive filenames:
  ✓ Good: "Pong [Paul Vervalin, 1990].ch8"
  ✓ Good: "Space Invaders [David Winter].ch8"
  ✗ Poor: "game1.ch8"
  ✗ Poor: "si.ch8"

METADATA FILES:
Create .txt files for detailed information:
  
  For "Pong.ch8", create "Pong.txt":
  Title: Pong
  Author: Paul Vervalin
  Year: 1990
  Description: The classic paddle and ball game
  Instructions: Player 1: 1 and Q keys, Player 2: 4 and R keys
  
TIPS:
- Keep directories under 100 ROMs for best performance
- Use consistent naming for better sorting
- Include author and year in filenames when known
- Remove duplicate or broken ROMs regularly`,

		"advanced_search": `=== ADVANCED SEARCH TUTORIAL ===

Master the search system to quickly find any ROM in your collection.

BASIC SEARCH:
1. Start typing in the ROM browser
2. Search box appears at top-right
3. Results filter in real-time
4. Press ESC to clear search

SEARCH FEATURES:
- Case insensitive: "PONG" finds "pong"
- Partial matching: "pac" finds "Pac-Man"
- Multi-word: "space invaders" finds "Space Invaders"
- Metadata search: Includes author, year, description

SEARCH EXAMPLES:
  "david"       → Find all ROMs by David Winter
  "1990"        → Find ROMs from 1990
  "arcade"      → Find arcade games
  "puzzle"      → Find puzzle games
  "pong space"  → Find ROMs with both "pong" and "space"

SEARCH TIPS:
- Use author names to find ROM collections
- Search by year to find retro classics
- Use genre terms for game categories
- Search partial words for quick filtering

COMBINING WITH SORT:
1. Search for a category: "arcade"
2. Sort results: S+N for name, S+D for date
3. Navigate filtered results with ↑/↓

The search system makes large ROM collections manageable!`,

		"customization": `=== CUSTOMIZATION TUTORIAL ===

Personalize your CHIP-8 emulator for the best experience.

VISUAL THEMES:
1. Press F2 for settings
2. Navigate to "Theme" option
3. Try different themes:
   - Default: Standard colors
   - Dark: Dark mode for low light
   - Light: Bright mode for daytime
   - Classic: Retro terminal style

EMULATOR SETTINGS:
Adjust game performance:
- Clock Speed: 500-1000 Hz (higher = faster games)
- Scale: 8-15 (larger = bigger display)
- Color Palette: Choose game colors
- Sound: Enable/disable beep sounds

BROWSER SETTINGS:
Customize the ROM browser:
- Page Size: ROMs per screen (10-25)
- Sort Order: Default sorting preference
- Show Hidden: Display hidden files
- Max Recent: Number of recent ROMs to track

KEY BINDINGS:
Customize controls in config file (~/.chip8/config.json):
{
  "input": {
    "keyBindings": {
      "menu_keys": {
        "navigate_up": "k",    // Vim style
        "navigate_down": "j",
        "select": "Return"
      }
    }
  }
}

PERFORMANCE TUNING:
For large collections:
- Enable caching
- Reduce cache size if low memory
- Disable background scanning if slow
- Use SSD storage for ROM directory

Save your perfect setup and enjoy!`,

		"troubleshooting_guide": `=== TROUBLESHOOTING GUIDE ===

Solve common issues with the CHIP-8 emulator.

PROBLEM: No ROMs appear in browser
SOLUTION:
1. Check ROM directory setting (F2)
2. Verify .ch8 files exist in directory
3. Check file permissions (read access)
4. Press F5 to refresh directory
5. Look for error messages in terminal

PROBLEM: Game runs too fast or slow
SOLUTION:
1. Press F2 for settings
2. Adjust "Clock Speed" setting
3. Try these speeds:
   - 500 Hz: Slower games
   - 700 Hz: Standard speed
   - 1000 Hz: Faster games
4. Some games need specific speeds

PROBLEM: Controls don't work in game
SOLUTION:
1. Check CHIP-8 keypad mapping:
   1,2,3,4 / Q,W,E,R / A,S,D,F / Z,X,C,V
2. Try different keys for the game
3. Check game instructions/metadata
4. Reset key bindings to default

PROBLEM: No sound
SOLUTION:
1. Check sound setting (F2)
2. Test system audio with other apps
3. Toggle sound off and on
4. Check volume levels

PROBLEM: Poor performance/slow browsing
SOLUTION:
1. Reduce directory size (<100 ROMs)
2. Enable caching in settings
3. Close other applications
4. Use SSD storage
5. Reduce page size setting

PROBLEM: Settings not saved
SOLUTION:
1. Check ~/.chip8/ directory permissions
2. Verify disk space available
3. Try running with admin rights
4. Delete config file to recreate

GETTING MORE HELP:
- Check log files in ~/.chip8/logs/
- Report bugs on GitHub repository
- Ask community for game-specific help`,
	}
}

// GetTipsAndTricks returns helpful tips for users
func GetTipsAndTricks() []string {
	return []string{
		"Tip: Use descriptive ROM filenames for better search results",
		"Tip: Create .txt metadata files for detailed ROM information",
		"Tip: Organize ROMs in subdirectories for better performance",
		"Tip: Press F to add frequently played ROMs to favorites",
		"Tip: Use S+N, S+D, S+Z to quickly sort your ROM collection",
		"Tip: Type immediately to search - no need to press a search key",
		"Tip: Adjust clock speed if games run too fast or slow",
		"Tip: Use F4 to quickly access recently played ROMs",
		"Tip: Press F5 to refresh if new ROMs don't appear",
		"Tip: Enable caching in settings for faster large directory browsing",
		"Tip: Use Page Up/Down to quickly navigate long ROM lists",
		"Tip: ESC clears search and returns to full ROM list",
		"Tip: F11 toggles fullscreen mode during gameplay",
		"Tip: Check ROM metadata for game instructions and controls",
		"Tip: Dark theme is easier on the eyes for extended play sessions",
		"Tip: Backup your config file to preserve custom settings",
		"Tip: Use Ctrl+I to view detailed ROM information",
		"Tip: Sort by date to find recently added ROMs",
		"Tip: Home/End keys jump to first/last ROM in list",
		"Tip: Multiple subdirectories help organize large collections",
	}
}

// GetGameSpecificHelp returns help for common CHIP-8 games
func GetGameSpecificHelp() map[string]string {
	return map[string]string{
		"pong": `PONG CONTROLS:
Player 1: 1 and Q keys (up/down)
Player 2: 4 and R keys (up/down)

OBJECTIVE:
Classic paddle and ball game. First to score wins!`,

		"space_invaders": `SPACE INVADERS CONTROLS:
A and D: Move left/right
S: Fire

OBJECTIVE:
Destroy all invaders before they reach the bottom!`,

		"tetris": `TETRIS CONTROLS:
Q and E: Rotate piece
A and D: Move left/right
S: Drop faster

OBJECTIVE:
Complete horizontal lines to clear them!`,

		"breakout": `BREAKOUT CONTROLS:
A and D: Move paddle left/right

OBJECTIVE:
Break all blocks with the ball. Don't let the ball fall!`,

		"pac_man": `PAC-MAN CONTROLS:
W, A, S, D: Move up, left, down, right

OBJECTIVE:
Eat all dots while avoiding ghosts!`,

		"frogger": `FROGGER CONTROLS:
W, A, S, D: Move up, left, down, right

OBJECTIVE:
Guide the frog across traffic and rivers to safety!`,
	}
}

// GetPerformanceTips returns performance optimization advice
func GetPerformanceTips() []string {
	return []string{
		"Keep ROM directories under 100 files for optimal performance",
		"Enable caching to speed up repeated directory scans",
		"Use SSD storage for your ROM directory when possible",
		"Organize large collections into category subdirectories",
		"Disable metadata extraction for very large collections",
		"Reduce page size setting on slower systems",
		"Close other resource-intensive applications",
		"Set appropriate memory limits in performance settings",
		"Use background scanning for better responsiveness",
		"Regular cleanup of unused or duplicate ROMs improves speed",
		"Consider using a faster CPU for large ROM collections",
		"Disable animations on older systems for better performance",
		"Use simple themes to reduce rendering overhead",
		"Monitor system memory usage with large collections",
		"Restart the application periodically to clear cache buildup",
	}
}

// FormatTutorial formats tutorial content for display
func FormatTutorial(title, content string) string {
	var builder strings.Builder

	// Add border and title
	border := strings.Repeat("=", 60)
	builder.WriteString(border + "\n")
	builder.WriteString(content)
	builder.WriteString("\n" + border + "\n")

	return builder.String()
}

// GetRandomTip returns a random tip for display
func GetRandomTip() string {
	tips := GetTipsAndTricks()
	if len(tips) == 0 {
		return "Tip: Press F1 for help anytime!"
	}

	// Simple pseudo-random selection (in real implementation, use proper random)
	return tips[0] // For simplicity, return first tip
}

// GetContextualTips returns tips relevant to current context
func GetContextualTips(context string) []string {
	allTips := GetTipsAndTricks()
	contextualTips := make([]string, 0)

	// Filter tips based on context
	contextKeywords := map[string][]string{
		"browser":   {"ROM", "search", "sort", "directory", "metadata"},
		"emulator":  {"game", "speed", "controls", "fullscreen"},
		"settings":  {"config", "theme", "performance", "settings"},
		"search":    {"search", "find", "filter"},
		"favorites": {"favorites", "recent"},
	}

	if keywords, exists := contextKeywords[context]; exists {
		for _, tip := range allTips {
			for _, keyword := range keywords {
				if strings.Contains(strings.ToLower(tip), strings.ToLower(keyword)) {
					contextualTips = append(contextualTips, tip)
					break
				}
			}
		}
	}

	// If no contextual tips found, return general tips
	if len(contextualTips) == 0 {
		return allTips[:5] // Return first 5 general tips
	}

	return contextualTips
}
