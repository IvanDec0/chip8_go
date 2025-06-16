package ui

import (
	"chip8/internal/menu"
	"time"
)

// MessageType represents the type of message to display
type MessageType int

const (
	MessageError MessageType = iota
	MessageInfo
	MessageWarning
	MessageSuccess
)

// Message represents a timed message to display
type Message struct {
	Text      string
	Type      MessageType
	ExpiresAt time.Time
}

// MessageRenderer handles error messages and notifications
type MessageRenderer struct {
	baseRenderer BaseRenderer
	theme        MessageTheme
	currentMsg   *Message
}

// MessageTheme holds color configuration for different message types
type MessageTheme struct {
	ErrorColor   menu.Color
	InfoColor    menu.Color
	WarningColor menu.Color
	SuccessColor menu.Color
	Background   menu.Color
}

// BaseRenderer interface to avoid import cycles
type BaseRenderer interface {
	RenderText(text string, x, y int, color menu.Color) error
	ClearWithColor(color menu.Color)
	Present()
}

// DefaultMessageTheme returns a default message theme
func DefaultMessageTheme() MessageTheme {
	return MessageTheme{
		ErrorColor:   menu.Color{R: 255, G: 100, B: 100, A: 255}, // Light red
		InfoColor:    menu.Color{R: 100, G: 150, B: 255, A: 255}, // Light blue
		WarningColor: menu.Color{R: 255, G: 200, B: 100, A: 255}, // Light orange
		SuccessColor: menu.Color{R: 100, G: 255, B: 100, A: 255}, // Light green
		Background:   menu.Color{R: 0, G: 0, B: 0, A: 200},       // Semi-transparent black
	}
}

// NewMessageRenderer creates a new message renderer
func NewMessageRenderer(baseRenderer BaseRenderer) *MessageRenderer {
	return &MessageRenderer{
		baseRenderer: baseRenderer,
		theme:        DefaultMessageTheme(),
		currentMsg:   nil,
	}
}

// ShowError displays an error message for the specified duration
func (mr *MessageRenderer) ShowError(message string, duration time.Duration) {
	mr.showMessage(message, MessageError, duration)
}

// ShowInfo displays an info message for the specified duration
func (mr *MessageRenderer) ShowInfo(message string, duration time.Duration) {
	mr.showMessage(message, MessageInfo, duration)
}

// ShowWarning displays a warning message for the specified duration
func (mr *MessageRenderer) ShowWarning(message string, duration time.Duration) {
	mr.showMessage(message, MessageWarning, duration)
}

// ShowSuccess displays a success message for the specified duration
func (mr *MessageRenderer) ShowSuccess(message string, duration time.Duration) {
	mr.showMessage(message, MessageSuccess, duration)
}

// showMessage shows a message of the specified type
func (mr *MessageRenderer) showMessage(text string, msgType MessageType, duration time.Duration) {
	mr.currentMsg = &Message{
		Text:      text,
		Type:      msgType,
		ExpiresAt: time.Now().Add(duration),
	}
}

// Update checks if the current message has expired and clears it
func (mr *MessageRenderer) Update() {
	if mr.currentMsg != nil && time.Now().After(mr.currentMsg.ExpiresAt) {
		mr.currentMsg = nil
	}
}

// HasActiveMessage returns true if there's an active message to display
func (mr *MessageRenderer) HasActiveMessage() bool {
	mr.Update() // Update expiration first
	return mr.currentMsg != nil
}

// RenderMessage renders the current active message if any
func (mr *MessageRenderer) RenderMessage(windowWidth, windowHeight int) error {
	if !mr.HasActiveMessage() {
		return nil
	}

	msg := mr.currentMsg

	// Calculate message position (centered)
	charWidth := 8
	charHeight := 8
	padding := 8

	textWidth := len(msg.Text) * charWidth
	textHeight := charHeight

	msgWidth := textWidth + padding*2
	msgHeight := textHeight + padding*2

	x := (windowWidth - msgWidth) / 2
	y := (windowHeight - msgHeight) / 2

	// Choose color based on message type
	var textColor menu.Color
	switch msg.Type {
	case MessageError:
		textColor = mr.theme.ErrorColor
	case MessageInfo:
		textColor = mr.theme.InfoColor
	case MessageWarning:
		textColor = mr.theme.WarningColor
	case MessageSuccess:
		textColor = mr.theme.SuccessColor
	default:
		textColor = mr.theme.InfoColor
	}

	// Render message text
	textX := x + padding
	textY := y + padding

	return mr.baseRenderer.RenderText(msg.Text, textX, textY, textColor)
}

// ClearMessage clears any active message
func (mr *MessageRenderer) ClearMessage() {
	mr.currentMsg = nil
}

// SetTheme updates the message theme
func (mr *MessageRenderer) SetTheme(theme MessageTheme) {
	mr.theme = theme
}
