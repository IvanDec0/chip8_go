package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// ConfigManager handles configuration persistence
type ConfigManager struct {
	configPath    string
	backupPath    string
	userConfig    *UserConfig
	defaultConfig *UserConfig
}

// NewConfigManager creates a new configuration manager
func NewConfigManager(configPath string) *ConfigManager {
	if configPath == "" {
		configPath = getDefaultConfigPath()
	}

	backupPath := configPath + ".backup"

	return &ConfigManager{
		configPath:    configPath,
		backupPath:    backupPath,
		userConfig:    DefaultUserConfig(),
		defaultConfig: DefaultUserConfig(),
	}
}

// LoadConfig loads configuration from file
func (cm *ConfigManager) LoadConfig() error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Config file doesn't exist, create default one
		return cm.SaveConfig()
	}

	// Load config from file
	file, err := os.Open(cm.configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var loadedConfig UserConfig
	if err := decoder.Decode(&loadedConfig); err != nil {
		// Config file is corrupted, try to restore from backup
		if backupErr := cm.RestoreFromBackup(); backupErr != nil {
			// Backup also failed, use default config
			cm.userConfig = DefaultUserConfig()
			return fmt.Errorf("config file corrupted and backup failed, using defaults: %w", err)
		}
		return nil
	}

	// Validate loaded config
	loadedConfig.Validate()
	cm.userConfig = &loadedConfig

	return nil
}

// SaveConfig saves configuration to file
func (cm *ConfigManager) SaveConfig() error {
	// Create backup before saving
	if err := cm.BackupConfig(); err != nil {
		// Log warning but don't fail the save
		fmt.Printf("Warning: failed to create config backup: %v\n", err)
	}

	// Update timestamp
	cm.userConfig.LastUpdated = time.Now()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(cm.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write to temporary file first
	tempPath := cm.configPath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp config file: %w", err)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(cm.userConfig); err != nil {
		file.Close()
		os.Remove(tempPath)
		return fmt.Errorf("failed to encode config: %w", err)
	}

	file.Close()

	// Atomically replace the config file
	if err := os.Rename(tempPath, cm.configPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("failed to save config file: %w", err)
	}

	return nil
}

// GetUserConfig returns the current user configuration
func (cm *ConfigManager) GetUserConfig() *UserConfig {
	return cm.userConfig
}

// SetUserConfig sets a new user configuration
func (cm *ConfigManager) SetUserConfig(config *UserConfig) error {
	if config == nil {
		return fmt.Errorf("config cannot be nil")
	}

	// Validate config before setting
	config.Validate()
	cm.userConfig = config

	return cm.SaveConfig()
}

// UpdateUserConfig updates specific parts of the configuration
func (cm *ConfigManager) UpdateUserConfig(updater func(*UserConfig)) error {
	if updater == nil {
		return fmt.Errorf("updater function cannot be nil")
	}

	// Create a clone to avoid modifying the original if validation fails
	clone := cm.userConfig.Clone()
	updater(clone)
	clone.Validate()

	cm.userConfig = clone
	return cm.SaveConfig()
}

// ResetToDefaults resets configuration to default values
func (cm *ConfigManager) ResetToDefaults() error {
	cm.userConfig = DefaultUserConfig()
	return cm.SaveConfig()
}

// BackupConfig creates a backup of the current configuration
func (cm *ConfigManager) BackupConfig() error {
	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		return nil // No config to backup
	}

	// Copy config file to backup location
	return copyFile(cm.configPath, cm.backupPath)
}

// RestoreFromBackup restores configuration from backup
func (cm *ConfigManager) RestoreFromBackup() error {
	// Check if backup exists
	if _, err := os.Stat(cm.backupPath); os.IsNotExist(err) {
		return fmt.Errorf("no backup file found")
	}

	// Copy backup to config location
	if err := copyFile(cm.backupPath, cm.configPath); err != nil {
		return fmt.Errorf("failed to restore from backup: %w", err)
	}

	// Load the restored config
	return cm.LoadConfig()
}

// RestoreConfig restores configuration from a specific backup file
func (cm *ConfigManager) RestoreConfig(backupPath string) error {
	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backupPath)
	}

	// Validate backup file by trying to load it
	file, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var testConfig UserConfig
	if err := decoder.Decode(&testConfig); err != nil {
		return fmt.Errorf("backup file is corrupted: %w", err)
	}

	// Copy backup to config location
	if err := copyFile(backupPath, cm.configPath); err != nil {
		return fmt.Errorf("failed to restore from backup: %w", err)
	}

	// Load the restored config
	return cm.LoadConfig()
}

// GetConfigPath returns the path to the configuration file
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configPath
}

// GetBackupPath returns the path to the backup configuration file
func (cm *ConfigManager) GetBackupPath() string {
	return cm.backupPath
}

// HasBackup checks if a backup file exists
func (cm *ConfigManager) HasBackup() bool {
	_, err := os.Stat(cm.backupPath)
	return err == nil
}

// MigrateConfig migrates configuration from an older version
func (cm *ConfigManager) MigrateConfig() error {
	// This can be extended in the future for version migrations
	if cm.userConfig.Version == "" {
		cm.userConfig.Version = "1.0"
		return cm.SaveConfig()
	}
	return nil
}

// ValidateConfig validates the current configuration
func (cm *ConfigManager) ValidateConfig() []string {
	var warnings []string

	// Check if ROM directory exists
	if cm.userConfig.MenuConfig.DefaultROMDirectory != "" {
		if _, err := os.Stat(cm.userConfig.MenuConfig.DefaultROMDirectory); os.IsNotExist(err) {
			warnings = append(warnings, fmt.Sprintf("Default ROM directory does not exist: %s", cm.userConfig.MenuConfig.DefaultROMDirectory))
		}
	}

	// Check key bindings
	for action, key := range cm.userConfig.InputConfig.CustomKeyBindings {
		if _, valid := cm.userConfig.InputConfig.GetKeyCode(key); !valid {
			warnings = append(warnings, fmt.Sprintf("Invalid key binding for %s: %s", action, key))
		}
	}

	return warnings
}

// getDefaultConfigPath returns the default configuration file path
func getDefaultConfigPath() string {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			configDir = os.Getenv("USERPROFILE")
		}
	case "darwin":
		configDir = os.Getenv("HOME")
		if configDir != "" {
			configDir = filepath.Join(configDir, "Library", "Application Support")
		}
	default: // Linux and other Unix-like systems
		//configDir = os.Getenv("XDG_CONFIG_HOME")
		configDir = (".")
		if configDir == "" {
			home := os.Getenv("HOME")
			if home != "" {
				configDir = filepath.Join(home, ".config")
			}
		}
	}

	// Fallback to current directory if all else fails
	if configDir == "" {
		configDir = "."
	}

	return filepath.Join(configDir, "chip8_save", "config.json")
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	buffer := make([]byte, 32*1024) // 32KB buffer
	for {
		n, err := sourceFile.Read(buffer)
		if n > 0 {
			if _, writeErr := destFile.Write(buffer[:n]); writeErr != nil {
				return writeErr
			}
		}
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return err
		}
	}

	return nil
}
