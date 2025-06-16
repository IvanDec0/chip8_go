package config

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// ValidationResult represents the result of configuration validation
type ValidationResult struct {
	IsValid  bool
	Errors   []ValidationError
	Warnings []ValidationWarning
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

// ValidationWarning represents a configuration validation warning
type ValidationWarning struct {
	Field   string
	Value   interface{}
	Message string
}

// ConfigValidator provides configuration validation functionality
type ConfigValidator struct {
	strictMode bool
}

// NewConfigValidator creates a new configuration validator
func NewConfigValidator(strictMode bool) *ConfigValidator {
	return &ConfigValidator{
		strictMode: strictMode,
	}
}

// ValidateUserConfig validates a complete user configuration
func (cv *ConfigValidator) ValidateUserConfig(config *UserConfig) *ValidationResult {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   make([]ValidationError, 0),
		Warnings: make([]ValidationWarning, 0),
	}

	// Validate menu configuration
	cv.validateMenuConfig(&config.MenuConfig, result)

	// Validate theme configuration
	cv.validateThemeConfig(&config.ThemeConfig, result)

	// Validate input configuration
	cv.validateInputConfig(&config.InputConfig, result)

	// Validate emulator configuration
	cv.validateEmulatorConfig(&config.EmulatorConfig, result)

	// Validate cache configuration
	cv.validateCacheConfig(&config.CacheConfig, result)

	// Validate overall configuration consistency
	cv.validateConfigConsistency(config, result)

	// Set overall validity
	result.IsValid = len(result.Errors) == 0

	return result
}

// validateMenuConfig validates menu configuration
func (cv *ConfigValidator) validateMenuConfig(config *MenuConfig, result *ValidationResult) {
	// Validate default ROM directory
	if config.DefaultROMDirectory != "" {
		if !cv.isValidDirectory(config.DefaultROMDirectory) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "MenuConfig.DefaultROMDirectory",
				Value:   config.DefaultROMDirectory,
				Message: "Directory does not exist or is not accessible",
			})
		}
	}

	// Validate sort criteria
	validSortOptions := []string{"name", "date", "size", "recent"}
	if !cv.isValidEnum(config.SortBy, validSortOptions) {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "MenuConfig.SortBy",
			Value:   config.SortBy,
			Message: fmt.Sprintf("Invalid sort option. Valid options: %v", validSortOptions),
		})
	}

	// Validate MaxRecentROMs
	if config.MaxRecentROMs < 1 || config.MaxRecentROMs > 1000 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "MenuConfig.MaxRecentROMs",
			Value:   config.MaxRecentROMs,
			Message: "MaxRecentROMs must be between 1 and 1000",
		})
	}

	// Warning for high MaxRecentROMs
	if config.MaxRecentROMs > 100 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "MenuConfig.MaxRecentROMs",
			Value:   config.MaxRecentROMs,
			Message: "High MaxRecentROMs value may impact performance",
		})
	}
}

// validateThemeConfig validates theme configuration
func (cv *ConfigValidator) validateThemeConfig(config *ThemeConfig, result *ValidationResult) {
	// Validate theme name
	if config.Name == "" {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "ThemeConfig.Name",
			Value:   config.Name,
			Message: "Theme name cannot be empty",
		})
	}

	// Validate colors
	colors := map[string]Color{
		"BackgroundColor": config.BackgroundColor,
		"TextColor":       config.TextColor,
		"SelectedColor":   config.SelectedColor,
		"BorderColor":     config.BorderColor,
		"AccentColor":     config.AccentColor,
	}

	for colorName, color := range colors {
		if !cv.isValidColor(color) {
			result.Errors = append(result.Errors, ValidationError{
				Field:   fmt.Sprintf("ThemeConfig.%s", colorName),
				Value:   color,
				Message: "Invalid color values",
			})
		}
	}

	// Validate font size
	if config.FontSize < 8 || config.FontSize > 72 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "ThemeConfig.FontSize",
			Value:   config.FontSize,
			Message: "Font size must be between 8 and 72",
		})
	}

	// Warning for accessibility
	if cv.areColorsTooSimilar(config.BackgroundColor, config.TextColor) {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "ThemeConfig",
			Value:   "BackgroundColor/TextColor",
			Message: "Background and text colors may have poor contrast",
		})
	}
}

// validateInputConfig validates input configuration
func (cv *ConfigValidator) validateInputConfig(config *InputConfig, result *ValidationResult) {
	// Validate key repeat settings
	if config.KeyRepeatDelay > 1000 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "InputConfig.KeyRepeatDelay",
			Value:   config.KeyRepeatDelay,
			Message: "High key repeat delay may feel unresponsive",
		})
	}

	if config.KeyRepeatInterval < 10 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "InputConfig.KeyRepeatInterval",
			Value:   config.KeyRepeatInterval,
			Message: "Very low key repeat interval may cause input flooding",
		})
	}

	// Validate custom key bindings
	for key, binding := range config.CustomKeyBindings {
		if key == "" || binding == "" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "InputConfig.CustomKeyBindings",
				Value:   fmt.Sprintf("%s -> %s", key, binding),
				Message: "Key bindings cannot be empty",
			})
		}
	}
}

// validateEmulatorConfig validates emulator configuration
func (cv *ConfigValidator) validateEmulatorConfig(config *EmulatorConfig, result *ValidationResult) {
	// Validate scale
	if config.DefaultScale < 1 || config.DefaultScale > 50 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "EmulatorConfig.DefaultScale",
			Value:   config.DefaultScale,
			Message: "Default scale must be between 1 and 50",
		})
	}

	// Warning for extreme scales
	if config.DefaultScale > 20 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "EmulatorConfig.DefaultScale",
			Value:   config.DefaultScale,
			Message: "Very high scale may cause display issues on smaller screens",
		})
	}

	// Validate speed
	if config.DefaultSpeed < 60 || config.DefaultSpeed > 5000 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "EmulatorConfig.DefaultSpeed",
			Value:   config.DefaultSpeed,
			Message: "Default speed must be between 60 and 5000 Hz",
		})
	}

	// Warning for extreme speeds
	if config.DefaultSpeed > 2000 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "EmulatorConfig.DefaultSpeed",
			Value:   config.DefaultSpeed,
			Message: "Very high speed may break compatibility with some ROMs",
		})
	}

	// Validate audio settings
	if config.AudioVolume < 0.0 || config.AudioVolume > 1.0 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "EmulatorConfig.AudioVolume",
			Value:   config.AudioVolume,
			Message: "Audio volume must be between 0.0 and 1.0",
		})
	}

	validFrequencies := []int{22050, 44100, 48000}
	if !cv.isValidIntEnum(config.AudioFrequency, validFrequencies) {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "EmulatorConfig.AudioFrequency",
			Value:   config.AudioFrequency,
			Message: fmt.Sprintf("Unusual audio frequency. Common values: %v", validFrequencies),
		})
	}

	if config.BeepFrequency < 100.0 || config.BeepFrequency > 2000.0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "EmulatorConfig.BeepFrequency",
			Value:   config.BeepFrequency,
			Message: "Beep frequency outside typical range (100-2000 Hz)",
		})
	}
}

// validateCacheConfig validates cache configuration
func (cv *ConfigValidator) validateCacheConfig(config *CacheConfig, result *ValidationResult) {
	// Validate cache TTL values
	if config.MetadataCacheTTL < time.Second {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "CacheConfig.MetadataCacheTTL",
			Value:   config.MetadataCacheTTL,
			Message: "Very short cache TTL may impact performance",
		})
	}

	if config.DirectoryCacheTTL < time.Second {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "CacheConfig.DirectoryCacheTTL",
			Value:   config.DirectoryCacheTTL,
			Message: "Very short cache TTL may impact performance",
		})
	}

	// Validate cache size
	if config.MaxCacheSize < 10 {
		result.Errors = append(result.Errors, ValidationError{
			Field:   "CacheConfig.MaxCacheSize",
			Value:   config.MaxCacheSize,
			Message: "Cache size too small (minimum 10)",
		})
	}

	if config.MaxCacheSize > 100000 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "CacheConfig.MaxCacheSize",
			Value:   config.MaxCacheSize,
			Message: "Very large cache size may consume excessive memory",
		})
	}

	// Validate cleanup interval
	if config.CleanupInterval < time.Minute {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "CacheConfig.CleanupInterval",
			Value:   config.CleanupInterval,
			Message: "Frequent cache cleanup may impact performance",
		})
	}
}

// validateConfigConsistency validates overall configuration consistency
func (cv *ConfigValidator) validateConfigConsistency(config *UserConfig, result *ValidationResult) {
	// Check if caching is disabled but metadata is enabled
	if config.MenuConfig.ShowMetadata && !config.CacheConfig.MetadataCacheEnabled {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "Configuration",
			Value:   "Metadata/Cache",
			Message: "Metadata enabled but metadata caching disabled - may impact performance",
		})
	}

	// Check for conflicting settings
	if !config.MenuConfig.EnableSearch && config.MenuConfig.CaseSensitiveSearch {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "MenuConfig",
			Value:   "Search settings",
			Message: "Case-sensitive search enabled but search is disabled",
		})
	}

	// Check audio consistency
	if !config.EmulatorConfig.AudioEnabled && config.EmulatorConfig.AudioVolume > 0 {
		result.Warnings = append(result.Warnings, ValidationWarning{
			Field:   "EmulatorConfig",
			Value:   "Audio settings",
			Message: "Audio volume set but audio is disabled",
		})
	}
}

// Helper validation functions

// isValidDirectory checks if a directory path is valid and accessible
func (cv *ConfigValidator) isValidDirectory(path string) bool {
	if path == "" {
		return true // Empty is valid (will use default)
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// isValidEnum checks if a value is in a list of valid options
func (cv *ConfigValidator) isValidEnum(value string, validOptions []string) bool {
	for _, option := range validOptions {
		if value == option {
			return true
		}
	}
	return false
}

// isValidIntEnum checks if an integer value is in a list of valid options
func (cv *ConfigValidator) isValidIntEnum(value int, validOptions []int) bool {
	for _, option := range validOptions {
		if value == option {
			return true
		}
	}
	return false
}

// isValidColor validates color values
func (cv *ConfigValidator) isValidColor(color Color) bool {
	// All color values are valid since they're uint8
	// But we can check for extreme values or patterns
	return true
}

// areColorsTooSimilar checks if two colors have poor contrast
func (cv *ConfigValidator) areColorsTooSimilar(color1, color2 Color) bool {
	// Simple contrast check - calculate difference in brightness
	brightness1 := (int(color1.R)*299 + int(color1.G)*587 + int(color1.B)*114) / 1000
	brightness2 := (int(color2.R)*299 + int(color2.G)*587 + int(color2.B)*114) / 1000

	diff := brightness1 - brightness2
	if diff < 0 {
		diff = -diff
	}

	// If brightness difference is less than 125, consider it poor contrast
	return diff < 125
}

// SanitizeUserConfig sanitizes and fixes common configuration issues
func (cv *ConfigValidator) SanitizeUserConfig(config *UserConfig) *UserConfig {
	sanitized := config.Clone()

	// Sanitize menu config
	if sanitized.MenuConfig.MaxRecentROMs < 1 {
		sanitized.MenuConfig.MaxRecentROMs = 10
	}
	if sanitized.MenuConfig.MaxRecentROMs > 1000 {
		sanitized.MenuConfig.MaxRecentROMs = 100
	}

	// Sanitize theme config
	if sanitized.ThemeConfig.FontSize < 8 {
		sanitized.ThemeConfig.FontSize = 12
	}
	if sanitized.ThemeConfig.FontSize > 72 {
		sanitized.ThemeConfig.FontSize = 16
	}

	// Sanitize emulator config
	if sanitized.EmulatorConfig.DefaultScale < 1 {
		sanitized.EmulatorConfig.DefaultScale = 10
	}
	if sanitized.EmulatorConfig.DefaultScale > 50 {
		sanitized.EmulatorConfig.DefaultScale = 20
	}

	if sanitized.EmulatorConfig.DefaultSpeed < 60 {
		sanitized.EmulatorConfig.DefaultSpeed = 700
	}
	if sanitized.EmulatorConfig.DefaultSpeed > 5000 {
		sanitized.EmulatorConfig.DefaultSpeed = 1000
	}

	if sanitized.EmulatorConfig.AudioVolume < 0.0 {
		sanitized.EmulatorConfig.AudioVolume = 0.0
	}
	if sanitized.EmulatorConfig.AudioVolume > 1.0 {
		sanitized.EmulatorConfig.AudioVolume = 1.0
	}

	// Sanitize cache config
	if sanitized.CacheConfig.MaxCacheSize < 10 {
		sanitized.CacheConfig.MaxCacheSize = 100
	}

	if sanitized.CacheConfig.MetadataCacheTTL < time.Second {
		sanitized.CacheConfig.MetadataCacheTTL = time.Minute * 5
	}

	if sanitized.CacheConfig.DirectoryCacheTTL < time.Second {
		sanitized.CacheConfig.DirectoryCacheTTL = time.Minute * 2
	}

	if sanitized.CacheConfig.CleanupInterval < time.Minute {
		sanitized.CacheConfig.CleanupInterval = time.Minute * 10
	}

	return sanitized
}

// GenerateConfigReport generates a human-readable configuration report
func (cv *ConfigValidator) GenerateConfigReport(config *UserConfig) string {
	var report strings.Builder

	report.WriteString("=== CONFIGURATION REPORT ===\n\n")

	// Menu configuration
	report.WriteString("MENU CONFIGURATION:\n")
	report.WriteString(fmt.Sprintf("  ROM Directory: %s\n", config.MenuConfig.DefaultROMDirectory))
	report.WriteString(fmt.Sprintf("  Sort By: %s (%s)\n", config.MenuConfig.SortBy,
		map[bool]string{true: "ascending", false: "descending"}[config.MenuConfig.SortAscending]))
	report.WriteString(fmt.Sprintf("  Show Metadata: %t\n", config.MenuConfig.ShowMetadata))
	report.WriteString(fmt.Sprintf("  Max Recent ROMs: %d\n", config.MenuConfig.MaxRecentROMs))
	report.WriteString(fmt.Sprintf("  Search Enabled: %t\n", config.MenuConfig.EnableSearch))
	report.WriteString("\n")

	// Theme configuration
	report.WriteString("THEME CONFIGURATION:\n")
	report.WriteString(fmt.Sprintf("  Theme Name: %s\n", config.ThemeConfig.Name))
	report.WriteString(fmt.Sprintf("  Font Size: %d\n", config.ThemeConfig.FontSize))
	report.WriteString(fmt.Sprintf("  Background: RGB(%d,%d,%d)\n",
		config.ThemeConfig.BackgroundColor.R, config.ThemeConfig.BackgroundColor.G, config.ThemeConfig.BackgroundColor.B))
	report.WriteString(fmt.Sprintf("  Text Color: RGB(%d,%d,%d)\n",
		config.ThemeConfig.TextColor.R, config.ThemeConfig.TextColor.G, config.ThemeConfig.TextColor.B))
	report.WriteString("\n")

	// Emulator configuration
	report.WriteString("EMULATOR CONFIGURATION:\n")
	report.WriteString(fmt.Sprintf("  Default Scale: %d\n", config.EmulatorConfig.DefaultScale))
	report.WriteString(fmt.Sprintf("  Default Speed: %d Hz\n", config.EmulatorConfig.DefaultSpeed))
	report.WriteString(fmt.Sprintf("  Audio Enabled: %t\n", config.EmulatorConfig.AudioEnabled))
	if config.EmulatorConfig.AudioEnabled {
		report.WriteString(fmt.Sprintf("  Audio Volume: %.1f%%\n", config.EmulatorConfig.AudioVolume*100))
		report.WriteString(fmt.Sprintf("  Audio Frequency: %d Hz\n", config.EmulatorConfig.AudioFrequency))
		report.WriteString(fmt.Sprintf("  Beep Frequency: %.1f Hz\n", config.EmulatorConfig.BeepFrequency))
	}
	report.WriteString("\n")

	// Cache configuration
	report.WriteString("CACHE CONFIGURATION:\n")
	report.WriteString(fmt.Sprintf("  Metadata Cache: %t (TTL: %v)\n",
		config.CacheConfig.MetadataCacheEnabled, config.CacheConfig.MetadataCacheTTL))
	report.WriteString(fmt.Sprintf("  Directory Cache: %t (TTL: %v)\n",
		config.CacheConfig.DirectoryCacheEnabled, config.CacheConfig.DirectoryCacheTTL))
	report.WriteString(fmt.Sprintf("  Max Cache Size: %d entries\n", config.CacheConfig.MaxCacheSize))
	report.WriteString(fmt.Sprintf("  Cleanup Interval: %v\n", config.CacheConfig.CleanupInterval))
	report.WriteString("\n")

	// Validation
	validation := cv.ValidateUserConfig(config)
	report.WriteString("VALIDATION RESULTS:\n")
	report.WriteString(fmt.Sprintf("  Status: %s\n", map[bool]string{true: "VALID", false: "INVALID"}[validation.IsValid]))

	if len(validation.Errors) > 0 {
		report.WriteString("  ERRORS:\n")
		for _, err := range validation.Errors {
			report.WriteString(fmt.Sprintf("    - %s: %s\n", err.Field, err.Message))
		}
	}

	if len(validation.Warnings) > 0 {
		report.WriteString("  WARNINGS:\n")
		for _, warning := range validation.Warnings {
			report.WriteString(fmt.Sprintf("    - %s: %s\n", warning.Field, warning.Message))
		}
	}

	if validation.IsValid && len(validation.Warnings) == 0 {
		report.WriteString("  No issues found.\n")
	}

	return report.String()
}

// ValidateConfigFile validates a configuration file at the given path
func (cv *ConfigValidator) ValidateConfigFile(configPath string) (*ValidationResult, error) {
	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %s", configPath)
	}

	// Try to load the configuration
	cm := NewConfigManager(configPath)
	err := cm.LoadConfig()
	if err != nil {
		return &ValidationResult{
			IsValid: false,
			Errors: []ValidationError{{
				Field:   "ConfigFile",
				Value:   configPath,
				Message: fmt.Sprintf("Failed to load configuration: %v", err),
			}},
		}, nil
	}

	// Validate the loaded configuration
	return cv.ValidateUserConfig(cm.GetUserConfig()), nil
}
