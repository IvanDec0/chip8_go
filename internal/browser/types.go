package browser

import (
	"chip8/internal/menu"
)

// Browser handles ROM file discovery and navigation
type Browser interface {
	ScanDirectory(path string) ([]menu.ROMInfo, error)
	GetROMs() []menu.ROMInfo
	SetCurrentDirectory(path string) error
	GetCurrentDirectory() string
	NavigateUp() error
	Refresh() error
}

// BrowserError represents browser-specific errors
type BrowserError struct {
	Type    BrowserErrorType
	Message string
	Cause   error
}

// Error implements the error interface
func (e *BrowserError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// BrowserErrorType represents different browser error types
type BrowserErrorType int

const (
	ErrorFileAccess BrowserErrorType = iota
	ErrorInvalidDirectory
	ErrorPermissionDenied
	ErrorNotFound
)

// String returns a string representation of the browser error type
func (t BrowserErrorType) String() string {
	switch t {
	case ErrorFileAccess:
		return "FileAccess"
	case ErrorInvalidDirectory:
		return "InvalidDirectory"
	case ErrorPermissionDenied:
		return "PermissionDenied"
	case ErrorNotFound:
		return "NotFound"
	default:
		return "Unknown"
	}
}
