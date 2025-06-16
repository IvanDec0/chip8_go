package menu

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Layout calculates menu positioning and sizing based on window dimensions
type Layout struct {
	WindowWidth  int
	WindowHeight int
	Scale        int
	CharWidth    int // Width of a single character in pixels
	CharHeight   int // Height of a single character in pixels
}

// LayoutConfig holds configuration for layout calculations
type LayoutConfig struct {
	MarginTop    int
	MarginLeft   int
	MarginRight  int
	MarginBottom int
	ItemSpacing  int // Vertical spacing between menu items
	HeaderHeight int // Height reserved for header
	FooterHeight int // Height reserved for footer
}

// NewLayout creates a new layout calculator
func NewLayout(windowWidth, windowHeight, scale int) *Layout {
	return &Layout{
		WindowWidth:  windowWidth,
		WindowHeight: windowHeight,
		Scale:        scale,
		CharWidth:    8, // 8x8 bitmap font
		CharHeight:   8,
	}
}

// GetDefaultConfig returns default layout configuration
func (l *Layout) GetDefaultConfig() LayoutConfig {
	return LayoutConfig{
		MarginTop:    20,
		MarginLeft:   20,
		MarginRight:  20,
		MarginBottom: 20,
		ItemSpacing:  4,  // Small spacing between items
		HeaderHeight: 32, // 3 lines + spacing
		FooterHeight: 16, // 1 line + spacing
	}
}

// CalculateMenuBounds calculates the main menu area bounds
func (l *Layout) CalculateMenuBounds() (x, y, width, height int) {
	config := l.GetDefaultConfig()

	x = config.MarginLeft
	y = config.MarginTop + config.HeaderHeight
	width = l.WindowWidth - config.MarginLeft - config.MarginRight
	height = l.WindowHeight - config.MarginTop - config.MarginBottom - config.HeaderHeight - config.FooterHeight

	return
}

// CalculateItemBounds calculates bounds for a specific menu item
func (l *Layout) CalculateItemBounds(index, itemHeight int) (x, y, width, height int) {
	config := l.GetDefaultConfig()
	menuX, menuY, menuWidth, _ := l.CalculateMenuBounds()

	x = menuX
	y = menuY + index*(itemHeight+config.ItemSpacing)
	width = menuWidth
	height = itemHeight

	return
}

// GetMaxVisibleItems calculates how many items can be displayed at once
func (l *Layout) GetMaxVisibleItems(itemHeight int) int {
	config := l.GetDefaultConfig()
	_, _, _, menuHeight := l.CalculateMenuBounds()

	// Calculate how many items fit in the available height
	availableHeight := menuHeight
	if availableHeight <= 0 {
		return 1 // At least show one item
	}

	itemsWithSpacing := availableHeight / (itemHeight + config.ItemSpacing)
	if itemsWithSpacing < 1 {
		return 1
	}

	return itemsWithSpacing
}

// GetItemHeight returns the height needed for a menu item
func (l *Layout) GetItemHeight() int {
	return l.CharHeight + 4 // Character height plus padding
}

// CalculateHeaderBounds returns bounds for the header area
func (l *Layout) CalculateHeaderBounds() (x, y, width, height int) {
	config := l.GetDefaultConfig()

	x = config.MarginLeft
	y = config.MarginTop
	width = l.WindowWidth - config.MarginLeft - config.MarginRight
	height = config.HeaderHeight

	return
}

// CalculateFooterBounds returns bounds for the footer area
func (l *Layout) CalculateFooterBounds() (x, y, width, height int) {
	config := l.GetDefaultConfig()

	x = config.MarginLeft
	y = l.WindowHeight - config.MarginBottom - config.FooterHeight
	width = l.WindowWidth - config.MarginLeft - config.MarginRight
	height = config.FooterHeight

	return
}

// CalculateScrollInfo calculates scroll offset for large menus
func (l *Layout) CalculateScrollInfo(totalItems, selectedItem int) (startIndex, endIndex int) {
	maxVisible := l.GetMaxVisibleItems(l.GetItemHeight())

	if totalItems <= maxVisible {
		// All items fit, no scrolling needed
		return 0, totalItems - 1
	}

	// Calculate scroll position to keep selected item visible
	startIndex = selectedItem - maxVisible/2
	if startIndex < 0 {
		startIndex = 0
	}

	endIndex = startIndex + maxVisible - 1
	if endIndex >= totalItems {
		endIndex = totalItems - 1
		startIndex = endIndex - maxVisible + 1
		if startIndex < 0 {
			startIndex = 0
		}
	}

	return
}

// GetTextPosition calculates position for centered text within bounds
func (l *Layout) GetTextPosition(text string, boundsX, boundsY, boundsWidth, boundsHeight int, centered bool) (x, y int) {
	textWidth := len(text) * l.CharWidth

	if centered {
		x = boundsX + (boundsWidth-textWidth)/2
		y = boundsY + (boundsHeight-l.CharHeight)/2
	} else {
		x = boundsX + 4 // Small left padding
		y = boundsY + (boundsHeight-l.CharHeight)/2
	}

	return
}

// TruncateText truncates text to fit within specified width
func (l *Layout) TruncateText(text string, maxWidth int) string {
	maxChars := maxWidth / l.CharWidth
	if len(text) <= maxChars {
		return text
	}

	if maxChars <= 3 {
		return "..."
	}

	return text[:maxChars-3] + "..."
}

// BaseRenderer interface to avoid import cycle
type BaseRenderer interface {
	ClearWithColor(color Color)
	RenderText(text string, x, y int, color Color) error
	Present()
}

// MenuRenderer implements menu-specific rendering functionality
type MenuRenderer struct {
	baseRenderer BaseRenderer
	layout       *Layout
	theme        MenuTheme
}

// RenderConfig holds configuration for menu rendering
type RenderConfig struct {
	ShowHelp      bool
	ShowPath      bool
	MaxItemLength int
}

// NewMenuRenderer creates a new menu renderer
func NewMenuRenderer(baseRenderer BaseRenderer, windowWidth, windowHeight, scale int, theme MenuTheme) *MenuRenderer {
	return &MenuRenderer{
		baseRenderer: baseRenderer,
		layout:       NewLayout(windowWidth, windowHeight, scale),
		theme:        theme,
	}
}

// RenderMainMenu renders the main menu screen
func (mr *MenuRenderer) RenderMainMenu(items []MenuItem, selected int) error {
	// Clear screen with background color
	mr.baseRenderer.ClearWithColor(mr.theme.BackgroundColor)

	// Render header
	err := mr.renderHeader("CHIP-8 Emulator")
	if err != nil {
		return fmt.Errorf("failed to render header: %w", err)
	}

	// Render menu items
	err = mr.renderMenuItems(items, selected, "")
	if err != nil {
		return fmt.Errorf("failed to render menu items: %w", err)
	}

	// Render footer with help
	helpText := "↑/↓: Navigate  Enter: Select  ESC: Exit"
	err = mr.renderFooter(helpText)
	if err != nil {
		return fmt.Errorf("failed to render footer: %w", err)
	}

	mr.baseRenderer.Present()
	return nil
}

// RenderROMBrowser renders the ROM browser screen
func (mr *MenuRenderer) RenderROMBrowser(roms []ROMInfo, selected int, currentPath string) error {
	// Clear screen with background color
	mr.baseRenderer.ClearWithColor(mr.theme.BackgroundColor)

	// Render header with current path
	headerText := fmt.Sprintf("ROM Browser - %s", filepath.Base(currentPath))
	err := mr.renderHeader(headerText)
	if err != nil {
		return fmt.Errorf("failed to render header: %w", err)
	}

	// Convert ROMs to menu items for rendering
	items := make([]MenuItem, len(roms))
	for i, rom := range roms {
		var prefix string
		var action MenuAction

		if rom.IsDirectory {
			prefix = "📁 "
			action = ActionBrowseROMs
		} else {
			prefix = "🎮 "
			action = ActionLoadROM
		}

		items[i] = MenuItem{
			Text:   prefix + rom.Name,
			Action: action,
			Data:   rom.Path,
		}
	}

	// Render menu items with path context
	err = mr.renderMenuItems(items, selected, currentPath)
	if err != nil {
		return fmt.Errorf("failed to render ROM items: %w", err)
	}

	// Render footer with help
	helpText := "↑/↓: Navigate  Enter: Select  ESC: Back"
	err = mr.renderFooter(helpText)
	if err != nil {
		return fmt.Errorf("failed to render footer: %w", err)
	}

	mr.baseRenderer.Present()
	return nil
}

// RenderHelpScreen renders the help screen with detailed content
func (mr *MenuRenderer) RenderHelpScreen(items []MenuItem, selected int, helpContent string) error {
	// Clear screen with background color
	mr.baseRenderer.ClearWithColor(mr.theme.BackgroundColor)

	// Render header
	err := mr.renderHeader("Help - CHIP-8 Emulator")
	if err != nil {
		return fmt.Errorf("failed to render header: %w", err)
	}

	// If we have help content to display, show it instead of menu items
	if helpContent != "" {
		err = mr.renderHelpContent(helpContent)
		if err != nil {
			return fmt.Errorf("failed to render help content: %w", err)
		}
	} else {
		// Render menu items (help topics)
		err = mr.renderMenuItems(items, selected, "")
		if err != nil {
			return fmt.Errorf("failed to render help items: %w", err)
		}
	}

	// Render footer with help
	helpText := "↑/↓: Navigate  Enter: Select  ESC: Back"
	err = mr.renderFooter(helpText)
	if err != nil {
		return fmt.Errorf("failed to render footer: %w", err)
	}

	mr.baseRenderer.Present()
	return nil
}

// RenderExportScreen renders the export options screen
func (mr *MenuRenderer) RenderExportScreen(items []MenuItem, selected int, status string) error {
	// Clear screen with background color
	mr.baseRenderer.ClearWithColor(mr.theme.BackgroundColor)

	// Render header
	err := mr.renderHeader("Export Data - CHIP-8 Emulator")
	if err != nil {
		return fmt.Errorf("failed to render header: %w", err)
	}

	// Show status if available
	if status != "" {
		err = mr.renderExportStatus(status)
		if err != nil {
			return fmt.Errorf("failed to render export status: %w", err)
		}
	}

	// Render menu items
	err = mr.renderMenuItems(items, selected, "")
	if err != nil {
		return fmt.Errorf("failed to render export items: %w", err)
	}

	// Render footer with help
	helpText := "↑/↓: Navigate  Enter: Select  ESC: Back"
	err = mr.renderFooter(helpText)
	if err != nil {
		return fmt.Errorf("failed to render footer: %w", err)
	}

	mr.baseRenderer.Present()
	return nil
}

// renderHeader renders the header section
func (mr *MenuRenderer) renderHeader(title string) error {
	x, y, width, _ := mr.layout.CalculateHeaderBounds()

	// Center the title text
	titleX, titleY := mr.layout.GetTextPosition(title, x, y, width, mr.layout.CharHeight+8, true)

	return mr.baseRenderer.RenderText(title, titleX, titleY, mr.theme.TextColor)
}

// renderFooter renders the footer section with help text
func (mr *MenuRenderer) renderFooter(helpText string) error {
	x, y, width, height := mr.layout.CalculateFooterBounds()

	// Center the help text
	helpX, helpY := mr.layout.GetTextPosition(helpText, x, y, width, height, true)

	return mr.baseRenderer.RenderText(helpText, helpX, helpY, mr.theme.BorderColor)
}

// renderMenuItems renders the list of menu items with scrolling support
func (mr *MenuRenderer) renderMenuItems(items []MenuItem, selected int, currentPath string) error {
	if len(items) == 0 {
		// Render "No items found" message
		menuX, menuY, menuWidth, menuHeight := mr.layout.CalculateMenuBounds()
		emptyText := "No files found"
		textX, textY := mr.layout.GetTextPosition(emptyText, menuX, menuY, menuWidth, menuHeight, true)
		return mr.baseRenderer.RenderText(emptyText, textX, textY, mr.theme.BorderColor)
	}

	itemHeight := mr.layout.GetItemHeight()
	startIndex, endIndex := mr.layout.CalculateScrollInfo(len(items), selected)

	// Render visible items
	for i := startIndex; i <= endIndex && i < len(items); i++ {
		err := mr.renderMenuItem(items[i], i, selected, i-startIndex, itemHeight)
		if err != nil {
			return fmt.Errorf("failed to render menu item %d: %w", i, err)
		}
	}

	// Render scroll indicators if needed
	if len(items) > mr.layout.GetMaxVisibleItems(itemHeight) {
		err := mr.renderScrollIndicators(startIndex, endIndex, len(items))
		if err != nil {
			return fmt.Errorf("failed to render scroll indicators: %w", err)
		}
	}

	return nil
}

// renderMenuItem renders a single menu item
func (mr *MenuRenderer) renderMenuItem(item MenuItem, itemIndex, selectedIndex, displayIndex, itemHeight int) error {
	x, y, width, height := mr.layout.CalculateItemBounds(displayIndex, itemHeight)

	// Choose colors based on selection
	var textColor Color
	if itemIndex == selectedIndex {
		textColor = mr.theme.SelectedColor

		// Draw selection background (simple approach using text background)
		// We'll draw selection indicator instead since we don't have fill rectangle for menu
		selectionIndicator := "> "
		mr.baseRenderer.RenderText(selectionIndicator, x, y+(height-mr.layout.CharHeight)/2, textColor)
		x += len(selectionIndicator) * mr.layout.CharWidth
		width -= len(selectionIndicator) * mr.layout.CharWidth
	} else {
		textColor = mr.theme.TextColor
	}

	// Truncate text if too long
	displayText := mr.layout.TruncateText(item.Text, width)

	// Render the text
	textX, textY := mr.layout.GetTextPosition(displayText, x, y, width, height, false)
	return mr.baseRenderer.RenderText(displayText, textX, textY, textColor)
}

// renderScrollIndicators renders scroll indicators when menu is scrollable
func (mr *MenuRenderer) renderScrollIndicators(startIndex, endIndex, totalItems int) error {
	menuX, menuY, menuWidth, menuHeight := mr.layout.CalculateMenuBounds()

	// Show up arrow if not at top
	if startIndex > 0 {
		upArrow := "↑"
		arrowX := menuX + menuWidth - mr.layout.CharWidth*2
		arrowY := menuY
		err := mr.baseRenderer.RenderText(upArrow, arrowX, arrowY, mr.theme.BorderColor)
		if err != nil {
			return err
		}
	}

	// Show down arrow if not at bottom
	if endIndex < totalItems-1 {
		downArrow := "↓"
		arrowX := menuX + menuWidth - mr.layout.CharWidth*2
		arrowY := menuY + menuHeight - mr.layout.CharHeight
		err := mr.baseRenderer.RenderText(downArrow, arrowX, arrowY, mr.theme.BorderColor)
		if err != nil {
			return err
		}
	}

	return nil
}

// renderHelpContent renders detailed help content text
func (mr *MenuRenderer) renderHelpContent(content string) error {
	if mr.baseRenderer == nil {
		return fmt.Errorf("renderer not initialized")
	}

	layout := mr.layout.GetDefaultConfig()

	// Calculate content area bounds
	contentX := layout.MarginLeft
	contentY := layout.HeaderHeight + layout.MarginTop + 20
	maxWidth := mr.layout.WindowWidth - layout.MarginLeft - layout.MarginRight

	// Split content into lines and wrap long lines
	lines := strings.Split(content, "\n")
	y := contentY
	lineHeight := mr.layout.CharHeight + 2

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			y += lineHeight
			continue
		}

		// Word wrap if line is too long
		wrappedLines := mr.wrapText(line, maxWidth)
		for _, wrappedLine := range wrappedLines {
			if y+lineHeight > mr.layout.WindowHeight-layout.FooterHeight-layout.MarginBottom {
				// Content area full, stop rendering
				break
			}

			err := mr.baseRenderer.RenderText(wrappedLine, contentX, y, mr.theme.TextColor)
			if err != nil {
				return fmt.Errorf("failed to render help line: %w", err)
			}
			y += lineHeight
		}
	}

	return nil
}

// renderExportStatus renders export operation status
func (mr *MenuRenderer) renderExportStatus(status string) error {
	if mr.baseRenderer == nil {
		return fmt.Errorf("renderer not initialized")
	}

	layout := mr.layout.GetDefaultConfig()

	// Position status message below header
	x := layout.MarginLeft
	y := layout.HeaderHeight + layout.MarginTop + 10

	// Add a colored background or border to make it stand out
	err := mr.baseRenderer.RenderText("Status: "+status, x, y, mr.theme.SelectedColor)
	if err != nil {
		return fmt.Errorf("failed to render export status: %w", err)
	}

	return nil
}

// wrapText wraps long text lines to fit within the specified width
func (mr *MenuRenderer) wrapText(text string, maxWidth int) []string {
	if len(text)*mr.layout.CharWidth <= maxWidth {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	var lines []string
	var currentLine string

	for _, word := range words {
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine)*mr.layout.CharWidth <= maxWidth {
			currentLine = testLine
		} else {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	if len(lines) == 0 {
		return []string{text}
	}

	return lines
}

// Clear clears the menu rendering area
func (mr *MenuRenderer) Clear() error {
	mr.baseRenderer.ClearWithColor(mr.theme.BackgroundColor)
	return nil
}

// Present presents the rendered menu
func (mr *MenuRenderer) Present() error {
	mr.baseRenderer.Present()
	return nil
}

// RenderText renders text at specified position (delegated to base renderer)
func (mr *MenuRenderer) RenderText(text string, x, y int, color Color) error {
	return mr.baseRenderer.RenderText(text, x, y, color)
}

// UpdateLayout updates the layout when window size changes
func (mr *MenuRenderer) UpdateLayout(windowWidth, windowHeight, scale int) {
	mr.layout = NewLayout(windowWidth, windowHeight, scale)
}

// SetTheme updates the menu theme
func (mr *MenuRenderer) SetTheme(theme MenuTheme) {
	mr.theme = theme
}

// RenderROMListWithMetadata renders ROM list with metadata display
func (mr *MenuRenderer) RenderROMListWithMetadata(items []ROMInfo, selected int, startIndex int, showMetadata bool) error {
	if mr.baseRenderer == nil {
		return fmt.Errorf("renderer not initialized")
	}

	layout := mr.layout.GetDefaultConfig()
	y := layout.MarginTop + layout.HeaderHeight

	// Calculate visible items
	visibleHeight := mr.layout.WindowHeight - layout.MarginTop - layout.MarginBottom - layout.HeaderHeight - layout.FooterHeight
	maxVisibleItems := visibleHeight / (mr.layout.CharHeight + layout.ItemSpacing)

	// Adjust startIndex to keep selected item visible
	if selected < startIndex {
		startIndex = selected
	} else if selected >= startIndex+maxVisibleItems {
		startIndex = selected - maxVisibleItems + 1
	}

	// Ensure startIndex is within bounds
	if startIndex < 0 {
		startIndex = 0
	}
	if startIndex >= len(items) {
		startIndex = len(items) - 1
	}

	// Render items
	for i := startIndex; i < len(items) && i < startIndex+maxVisibleItems; i++ {
		rom := items[i]
		isSelected := (i == selected)

		// Determine colors
		textColor := mr.theme.TextColor
		if isSelected {
			textColor = mr.theme.SelectedColor
		}

		// Render ROM item
		romName := rom.Name
		if len(romName) > 40 { // Truncate long names
			romName = romName[:37] + "..."
		}

		// Add metadata if available and requested
		if showMetadata && !rom.IsDirectory {
			if metadata, err := mr.getMetadataForROM(rom.Path); err == nil && metadata != nil {
				if metadata.Author != "" {
					romName += fmt.Sprintf(" - %s", metadata.Author)
				}
			}
		}

		err := mr.baseRenderer.RenderText(romName, layout.MarginLeft, y, textColor)
		if err != nil {
			return err
		}

		// Show additional info for selected item
		if isSelected && showMetadata && !rom.IsDirectory {
			mr.renderSelectedROMDetails(rom, y+mr.layout.CharHeight+2)
		}

		y += mr.layout.CharHeight + layout.ItemSpacing
	}

	return nil
}

// RenderSearchBar renders the search input bar
func (mr *MenuRenderer) RenderSearchBar(query string, isActive bool) error {
	if mr.baseRenderer == nil {
		return fmt.Errorf("renderer not initialized")
	}

	layout := mr.layout.GetDefaultConfig()

	// Search bar background
	searchBarY := layout.MarginTop
	searchBarText := "Search: " + query
	if isActive {
		searchBarText += "_" // Show cursor
	}

	color := mr.theme.TextColor
	if isActive {
		color = mr.theme.SelectedColor
	}

	return mr.baseRenderer.RenderText(searchBarText, layout.MarginLeft, searchBarY, color)
}

// RenderSortIndicator renders the current sort criteria and direction
func (mr *MenuRenderer) RenderSortIndicator(criteria string, ascending bool) error {
	if mr.baseRenderer == nil {
		return fmt.Errorf("renderer not initialized")
	}

	layout := mr.layout.GetDefaultConfig()

	direction := "↑"
	if !ascending {
		direction = "↓"
	}

	sortText := fmt.Sprintf("Sort: %s %s", criteria, direction)

	// Position in top-right corner
	x := mr.layout.WindowWidth - len(sortText)*mr.layout.CharWidth - layout.MarginRight
	y := layout.MarginTop

	return mr.baseRenderer.RenderText(sortText, x, y, mr.theme.TextColor)
}

// RenderStatusBar renders a status bar at the bottom
func (mr *MenuRenderer) RenderStatusBar(status string) error {
	if mr.baseRenderer == nil {
		return fmt.Errorf("renderer not initialized")
	}

	layout := mr.layout.GetDefaultConfig()

	// Position at bottom
	y := mr.layout.WindowHeight - layout.MarginBottom - mr.layout.CharHeight

	return mr.baseRenderer.RenderText(status, layout.MarginLeft, y, mr.theme.TextColor)
}

// RenderScreenTitle renders the title for the current screen
func (mr *MenuRenderer) RenderScreenTitle(title string) error {
	if mr.baseRenderer == nil {
		return fmt.Errorf("renderer not initialized")
	}

	layout := mr.layout.GetDefaultConfig()

	// Center the title
	titleWidth := len(title) * mr.layout.CharWidth
	x := (mr.layout.WindowWidth - titleWidth) / 2
	y := layout.MarginTop

	return mr.baseRenderer.RenderText(title, x, y, mr.theme.SelectedColor)
}

// renderSelectedROMDetails renders detailed info for selected ROM
func (mr *MenuRenderer) renderSelectedROMDetails(rom ROMInfo, y int) error {
	metadata, err := mr.getMetadataForROM(rom.Path)
	if err != nil || metadata == nil {
		return nil // No metadata available
	}

	layout := mr.layout.GetDefaultConfig()
	x := layout.MarginLeft + 20 // Indent details

	// Show description if available
	if len(metadata.Description) > 0 && metadata.Description[0] != "" {
		desc := metadata.Description[0]
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		mr.baseRenderer.RenderText("  "+desc, x, y, mr.theme.TextColor)
		y += mr.layout.CharHeight + 2
	}

	// Show controls if available
	if len(metadata.Controls) > 0 && metadata.Controls[0] != "" {
		controls := "Controls: " + metadata.Controls[0]
		if len(controls) > 60 {
			controls = controls[:57] + "..."
		}
		mr.baseRenderer.RenderText("  "+controls, x, y, mr.theme.TextColor)
	}

	return nil
}

// getMetadataForROM retrieves metadata for a ROM (placeholder implementation)
func (mr *MenuRenderer) getMetadataForROM(romPath string) (*ROMMetadata, error) {
	// This would be implemented to interface with the browser's metadata system
	// For now, return nil to indicate no metadata
	return nil, fmt.Errorf("metadata not available")
}

// RenderFavoriteIndicator renders a favorite star indicator
func (mr *MenuRenderer) RenderFavoriteIndicator(x, y int, isFavorite bool) error {
	if !isFavorite {
		return nil
	}

	return mr.baseRenderer.RenderText("⭐", x, y, mr.theme.SelectedColor)
}

// RenderItemCount renders the current item count and total
func (mr *MenuRenderer) RenderItemCount(current, total int) error {
	if mr.baseRenderer == nil {
		return fmt.Errorf("renderer not initialized")
	}

	layout := mr.layout.GetDefaultConfig()

	countText := fmt.Sprintf("%d/%d", current+1, total)

	// Position in bottom-right corner
	x := mr.layout.WindowWidth - len(countText)*mr.layout.CharWidth - layout.MarginRight
	y := mr.layout.WindowHeight - layout.MarginBottom - mr.layout.CharHeight

	return mr.baseRenderer.RenderText(countText, x, y, mr.theme.TextColor)
}
