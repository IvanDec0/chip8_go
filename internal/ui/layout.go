package ui

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
