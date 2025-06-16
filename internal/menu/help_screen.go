package menu

import (
	"chip8/internal/help"
	"fmt"
	"strings"
)

// HelpScreen manages the in-game help system
type HelpScreen struct {
	helpSystem     *help.HelpSystem
	currentTopic   string
	currentContent string
	searchQuery    string
	searchResults  []*help.HelpTopic
	scrollOffset   int
	maxLines       int
	searchMode     bool
	showIndex      bool
}

// NewHelpScreen creates a new help screen
func NewHelpScreen() *HelpScreen {
	return &HelpScreen{
		helpSystem: help.NewHelpSystem(),
		maxLines:   20,
		showIndex:  true,
	}
}

// ShowHelp displays help for a specific topic or context
func (hs *HelpScreen) ShowHelp(context string) error {
	topic, err := hs.helpSystem.GetContextualHelp(context)
	if err != nil {
		// Fall back to general help
		topic, err = hs.helpSystem.GetTopic("general_help")
		if err != nil {
			return fmt.Errorf("help system unavailable: %w", err)
		}
	}

	hs.currentTopic = context
	hs.currentContent = hs.helpSystem.FormatHelp(topic)
	hs.scrollOffset = 0
	hs.searchMode = false
	hs.showIndex = false

	return nil
}

// ShowTopic displays a specific help topic
func (hs *HelpScreen) ShowTopic(topicName string) error {
	topic, err := hs.helpSystem.GetTopic(topicName)
	if err != nil {
		return err
	}

	hs.currentTopic = topicName
	hs.currentContent = hs.helpSystem.FormatHelp(topic)
	hs.scrollOffset = 0
	hs.searchMode = false
	hs.showIndex = false

	return nil
}

// ShowIndex displays the help index
func (hs *HelpScreen) ShowIndex() {
	hs.currentContent = hs.helpSystem.GetHelpIndex()
	hs.currentTopic = "index"
	hs.scrollOffset = 0
	hs.searchMode = false
	hs.showIndex = true
}

// ShowQuickReference displays the quick reference card
func (hs *HelpScreen) ShowQuickReference() {
	hs.currentContent = hs.helpSystem.GetQuickReference()
	hs.currentTopic = "quick_reference"
	hs.scrollOffset = 0
	hs.searchMode = false
	hs.showIndex = false
}

// Search searches help content
func (hs *HelpScreen) Search(query string) {
	hs.searchQuery = query
	hs.searchResults = hs.helpSystem.SearchHelp(query)
	hs.searchMode = true
	hs.showIndex = false
	hs.scrollOffset = 0

	// Format search results
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("=== SEARCH RESULTS FOR '%s' ===\n\n", query))

	if len(hs.searchResults) == 0 {
		builder.WriteString("No help topics found matching your search.\n\n")
		builder.WriteString("Try:\n")
		builder.WriteString("- Using fewer or different keywords\n")
		builder.WriteString("- Checking the help index (I key)\n")
		builder.WriteString("- Browsing by category\n")
	} else {
		builder.WriteString(fmt.Sprintf("Found %d matching topics:\n\n", len(hs.searchResults)))
		for i, topic := range hs.searchResults {
			builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, topic.Title))

			// Show snippet of content
			content := strings.ReplaceAll(topic.Content, "\n", " ")
			if len(content) > 100 {
				content = content[:100] + "..."
			}
			builder.WriteString(fmt.Sprintf("   %s\n\n", content))
		}

		builder.WriteString("Press number key (1-9) to view a topic, or ESC to return.\n")
	}

	hs.currentContent = builder.String()
}

// ProcessInput handles input for the help screen
func (hs *HelpScreen) ProcessInput(input string) HelpAction {
	switch strings.ToLower(input) {
	case "escape", "esc":
		return HelpActionClose

	case "i", "index":
		hs.ShowIndex()
		return HelpActionUpdate

	case "q", "quick":
		hs.ShowQuickReference()
		return HelpActionUpdate

	case "up", "k":
		hs.ScrollUp()
		return HelpActionUpdate

	case "down", "j":
		hs.ScrollDown()
		return HelpActionUpdate

	case "page_up":
		hs.PageUp()
		return HelpActionUpdate

	case "page_down":
		hs.PageDown()
		return HelpActionUpdate

	case "home":
		hs.ScrollToTop()
		return HelpActionUpdate

	case "end":
		hs.ScrollToBottom()
		return HelpActionUpdate

	case "/":
		// Start search mode
		return HelpActionStartSearch

	case "f1":
		// Show help about help
		hs.ShowTopic("general_help")
		return HelpActionUpdate

	default:
		// Check for search result selection (1-9)
		if hs.searchMode && len(input) == 1 && input >= "1" && input <= "9" {
			index := int(input[0] - '1')
			if index < len(hs.searchResults) {
				topic := hs.searchResults[index]
				hs.ShowTopic(hs.getTopicName(topic))
				return HelpActionUpdate
			}
		}

		// Check for topic shortcuts
		if topic := hs.getTopicByShortcut(input); topic != "" {
			hs.ShowTopic(topic)
			return HelpActionUpdate
		}

		return HelpActionNone
	}
}

// ScrollUp scrolls the help content up
func (hs *HelpScreen) ScrollUp() {
	if hs.scrollOffset > 0 {
		hs.scrollOffset--
	}
}

// ScrollDown scrolls the help content down
func (hs *HelpScreen) ScrollDown() {
	lines := strings.Split(hs.currentContent, "\n")
	maxScroll := len(lines) - hs.maxLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	if hs.scrollOffset < maxScroll {
		hs.scrollOffset++
	}
}

// PageUp scrolls up by a page
func (hs *HelpScreen) PageUp() {
	hs.scrollOffset -= hs.maxLines - 2
	if hs.scrollOffset < 0 {
		hs.scrollOffset = 0
	}
}

// PageDown scrolls down by a page
func (hs *HelpScreen) PageDown() {
	lines := strings.Split(hs.currentContent, "\n")
	maxScroll := len(lines) - hs.maxLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	hs.scrollOffset += hs.maxLines - 2
	if hs.scrollOffset > maxScroll {
		hs.scrollOffset = maxScroll
	}
}

// ScrollToTop scrolls to the beginning
func (hs *HelpScreen) ScrollToTop() {
	hs.scrollOffset = 0
}

// ScrollToBottom scrolls to the end
func (hs *HelpScreen) ScrollToBottom() {
	lines := strings.Split(hs.currentContent, "\n")
	maxScroll := len(lines) - hs.maxLines
	if maxScroll < 0 {
		maxScroll = 0
	}
	hs.scrollOffset = maxScroll
}

// GetDisplayContent returns the content to display on screen
func (hs *HelpScreen) GetDisplayContent() []string {
	lines := strings.Split(hs.currentContent, "\n")

	// Calculate display window
	start := hs.scrollOffset
	end := start + hs.maxLines

	if start >= len(lines) {
		return []string{"(End of content)"}
	}

	if end > len(lines) {
		end = len(lines)
	}

	return lines[start:end]
}

// GetStatusLine returns the status line for the help screen
func (hs *HelpScreen) GetStatusLine() string {
	lines := strings.Split(hs.currentContent, "\n")
	totalLines := len(lines)

	var status strings.Builder

	// Current topic/mode
	if hs.searchMode {
		status.WriteString(fmt.Sprintf("Search: '%s' (%d results)", hs.searchQuery, len(hs.searchResults)))
	} else if hs.showIndex {
		status.WriteString("Help Index")
	} else {
		status.WriteString(fmt.Sprintf("Help: %s", hs.currentTopic))
	}

	// Scroll position
	if totalLines > hs.maxLines {
		percentage := int(float64(hs.scrollOffset) / float64(totalLines-hs.maxLines) * 100)
		status.WriteString(fmt.Sprintf(" [%d%%]", percentage))
	}

	// Navigation hints
	status.WriteString(" | ↑/↓: Scroll | I: Index | Q: Quick Ref | /: Search | ESC: Exit")

	return status.String()
}

// GetShortcuts returns keyboard shortcuts for current context
func (hs *HelpScreen) GetShortcuts() []help.ShortcutInfo {
	context := "help"
	if hs.searchMode {
		context = "search"
	}
	return hs.helpSystem.GetShortcuts(context)
}

// getTopicName extracts topic name from HelpTopic (helper function)
func (hs *HelpScreen) getTopicName(topic *help.HelpTopic) string {
	// This would need to be implemented based on how topics are stored
	// For now, return the context or a derived name
	return strings.ToLower(strings.ReplaceAll(topic.Title, " ", "_"))
}

// getTopicByShortcut returns topic name for shortcut key
func (hs *HelpScreen) getTopicByShortcut(shortcut string) string {
	shortcuts := map[string]string{
		"b": "browser_help",
		"e": "emulator_help",
		"s": "settings_help",
		"m": "metadata_help",
		"t": "troubleshooting_help",
		"p": "performance_help",
	}

	return shortcuts[strings.ToLower(shortcut)]
}

// HelpAction represents possible actions from help screen
type HelpAction int

const (
	HelpActionNone HelpAction = iota
	HelpActionClose
	HelpActionUpdate
	HelpActionStartSearch
)

// SetMaxLines sets the maximum number of lines to display
func (hs *HelpScreen) SetMaxLines(maxLines int) {
	if maxLines > 0 {
		hs.maxLines = maxLines
	}
}

// GetCurrentTopic returns the current help topic
func (hs *HelpScreen) GetCurrentTopic() string {
	return hs.currentTopic
}

// IsSearchMode returns whether help is in search mode
func (hs *HelpScreen) IsSearchMode() bool {
	return hs.searchMode
}

// ClearSearch clears the current search
func (hs *HelpScreen) ClearSearch() {
	hs.searchMode = false
	hs.searchQuery = ""
	hs.searchResults = nil
	hs.ShowIndex() // Return to index after clearing search
}

// ShowTutorial displays a specific tutorial
func (hs *HelpScreen) ShowTutorial(tutorialName string) {
	tutorials := help.GetTutorialContent()
	if content, exists := tutorials[tutorialName]; exists {
		hs.currentContent = content
		hs.currentTopic = "tutorial_" + tutorialName
		hs.scrollOffset = 0
		hs.searchMode = false
		hs.showIndex = false
	}
}

// ShowTipsAndTricks displays helpful tips
func (hs *HelpScreen) ShowTipsAndTricks() {
	tips := help.GetTipsAndTricks()

	var builder strings.Builder
	builder.WriteString("=== TIPS AND TRICKS ===\n\n")

	for i, tip := range tips {
		builder.WriteString(fmt.Sprintf("%d. %s\n\n", i+1, tip))
	}

	builder.WriteString("Press ESC to return to help index.\n")

	hs.currentContent = builder.String()
	hs.currentTopic = "tips_and_tricks"
	hs.scrollOffset = 0
	hs.searchMode = false
	hs.showIndex = false
}
