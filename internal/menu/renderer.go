package menu

import (
	"fmt"
	"path/filepath"
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
