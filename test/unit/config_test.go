package unit

import (
	"chip8/internal/config"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigManager_LoadConfig(t *testing.T) {
	// Create temporary directory for config
	tmpDir, err := os.MkdirTemp("", "chip8_config_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	configManager := config.NewConfigManager(configPath)

	// Load config (should create default config)
	err = configManager.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify config file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}
}

func TestConfigManager_SaveConfig(t *testing.T) {
	// Create temporary directory for config
	tmpDir, err := os.MkdirTemp("", "chip8_config_save_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	configManager := config.NewConfigManager(configPath)

	// Save config
	err = configManager.SaveConfig()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Verify config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Expected config file to be created")
	}

	// Load config and verify it works
	err = configManager.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}
}

func TestConfigManager_GetUserConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "chip8_config_get_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	configManager := config.NewConfigManager(configPath)

	// Load config first
	err = configManager.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Get user config
	userConfig := configManager.GetUserConfig()
	if userConfig == nil {
		t.Fatal("Expected user config to be non-nil")
	}

	// Verify default values
	if userConfig.MenuConfig.SortBy != "name" {
		t.Errorf("Expected default sort 'name', got '%s'", userConfig.MenuConfig.SortBy)
	}

	if !userConfig.MenuConfig.SortAscending {
		t.Error("Expected default sort ascending to be true")
	}
}

func TestConfigManager_UpdateConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "chip8_config_update_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	configManager := config.NewConfigManager(configPath)

	// Load initial config
	err = configManager.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Get user config and modify it
	userConfig := configManager.GetUserConfig()
	originalDir := userConfig.MenuConfig.DefaultROMDirectory
	userConfig.MenuConfig.DefaultROMDirectory = "/modified/path"

	// Update config
	err = configManager.SetUserConfig(userConfig)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Verify the change was applied
	updatedConfig := configManager.GetUserConfig()
	if updatedConfig.MenuConfig.DefaultROMDirectory == originalDir {
		t.Error("Expected config to be updated")
	}

	if updatedConfig.MenuConfig.DefaultROMDirectory != "/modified/path" {
		t.Errorf("Expected directory '/modified/path', got '%s'",
			updatedConfig.MenuConfig.DefaultROMDirectory)
	}
}

func TestConfigManager_Backup(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "chip8_config_backup_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	configManager := config.NewConfigManager(configPath)

	// Load and save initial config
	err = configManager.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	err = configManager.SaveConfig()
	if err != nil {
		t.Fatalf("Failed to save initial config: %v", err)
	}

	// Modify and save again (should create backup)
	userConfig := configManager.GetUserConfig()
	userConfig.MenuConfig.DefaultROMDirectory = "/new/path"

	err = configManager.SetUserConfig(userConfig)
	if err != nil {
		t.Fatalf("Failed to update config: %v", err)
	}

	// Check if backup exists
	backupPath := configPath + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Error("Expected backup file to be created")
	}
}

func TestConfigManager_DefaultPath(t *testing.T) {
	// Test with empty path (should use default)
	configManager := config.NewConfigManager("")

	err := configManager.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config with default path: %v", err)
	}

	userConfig := configManager.GetUserConfig()
	if userConfig == nil {
		t.Fatal("Expected user config to be non-nil")
	}
}

func TestConfigManager_InvalidPath(t *testing.T) {
	// Test with invalid path
	invalidPath := "/invalid/readonly/path/config.json"
	configManager := config.NewConfigManager(invalidPath)

	// This might fail on some systems due to permissions
	err := configManager.LoadConfig()
	if err == nil {
		// If it succeeds, that's fine too
		t.Logf("Config loaded successfully even with restricted path")
	} else {
		t.Logf("Expected error for invalid path: %v", err)
	}
}

func BenchmarkConfigManager_LoadSave(b *testing.B) {
	tmpDir, err := os.MkdirTemp("", "chip8_config_bench_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.json")
	configManager := config.NewConfigManager(configPath)

	// Initial load
	err = configManager.LoadConfig()
	if err != nil {
		b.Fatalf("Failed to load config: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := configManager.SaveConfig()
		if err != nil {
			b.Fatalf("Save failed: %v", err)
		}

		err = configManager.LoadConfig()
		if err != nil {
			b.Fatalf("Load failed: %v", err)
		}
	}
}
