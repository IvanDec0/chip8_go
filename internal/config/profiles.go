package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ProfileManager manages configuration profiles
type ProfileManager struct {
	profilesDir    string
	currentProfile string
	profiles       map[string]*UserConfig
	configManager  *ConfigManager
}

// Profile represents a named configuration profile
type Profile struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Created     time.Time  `json:"created"`
	Modified    time.Time  `json:"modified"`
	Config      UserConfig `json:"config"`
	Tags        []string   `json:"tags,omitempty"`
	IsDefault   bool       `json:"isDefault,omitempty"`
}

// NewProfileManager creates a new profile manager
func NewProfileManager(configManager *ConfigManager) *ProfileManager {
	profilesDir := filepath.Join(filepath.Dir(configManager.configPath), "profiles")

	pm := &ProfileManager{
		profilesDir:   profilesDir,
		profiles:      make(map[string]*UserConfig),
		configManager: configManager,
	}

	// Create profiles directory if it doesn't exist
	os.MkdirAll(profilesDir, 0755)

	// Load existing profiles
	pm.loadProfiles()

	return pm
}

// CreateProfile creates a new configuration profile
func (pm *ProfileManager) CreateProfile(name, description string, config *UserConfig) error {
	if name == "" {
		return fmt.Errorf("profile name cannot be empty")
	}

	// Validate profile name
	if !isValidProfileName(name) {
		return fmt.Errorf("invalid profile name: %s", name)
	}

	// Check if profile already exists
	if _, exists := pm.profiles[name]; exists {
		return fmt.Errorf("profile '%s' already exists", name)
	}

	// Create profile
	profile := &Profile{
		Name:        name,
		Description: description,
		Created:     time.Now(),
		Modified:    time.Now(),
		Config:      *config.Clone(),
		Tags:        []string{},
		IsDefault:   false,
	}

	// Save profile to file
	err := pm.saveProfile(profile)
	if err != nil {
		return fmt.Errorf("failed to save profile: %w", err)
	}

	// Add to in-memory cache
	pm.profiles[name] = &profile.Config

	return nil
}

// LoadProfile loads a configuration profile
func (pm *ProfileManager) LoadProfile(name string) error {
	config, exists := pm.profiles[name]
	if !exists {
		return fmt.Errorf("profile '%s' not found", name)
	}

	// Apply the profile configuration
	err := pm.configManager.SetUserConfig(config)
	if err != nil {
		return fmt.Errorf("failed to apply profile: %w", err)
	}

	pm.currentProfile = name
	return nil
}

// DeleteProfile deletes a configuration profile
func (pm *ProfileManager) DeleteProfile(name string) error {
	if name == "default" {
		return fmt.Errorf("cannot delete default profile")
	}

	if pm.currentProfile == name {
		return fmt.Errorf("cannot delete currently active profile")
	}

	// Remove from memory
	delete(pm.profiles, name)

	// Remove file
	profilePath := filepath.Join(pm.profilesDir, name+".json")
	err := os.Remove(profilePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete profile file: %w", err)
	}

	return nil
}

// ListProfiles returns all available profiles
func (pm *ProfileManager) ListProfiles() []string {
	profiles := make([]string, 0, len(pm.profiles))
	for name := range pm.profiles {
		profiles = append(profiles, name)
	}
	sort.Strings(profiles)
	return profiles
}

// GetProfile returns a specific profile configuration
func (pm *ProfileManager) GetProfile(name string) (*UserConfig, error) {
	config, exists := pm.profiles[name]
	if !exists {
		return nil, fmt.Errorf("profile '%s' not found", name)
	}
	return config.Clone(), nil
}

// UpdateProfile updates an existing profile
func (pm *ProfileManager) UpdateProfile(name string, config *UserConfig) error {
	if _, exists := pm.profiles[name]; !exists {
		return fmt.Errorf("profile '%s' not found", name)
	}

	// Load existing profile to preserve metadata
	profile, err := pm.loadProfileFromFile(name)
	if err != nil {
		return fmt.Errorf("failed to load profile: %w", err)
	}

	// Update configuration and timestamp
	profile.Config = *config.Clone()
	profile.Modified = time.Now()

	// Save updated profile
	err = pm.saveProfile(profile)
	if err != nil {
		return fmt.Errorf("failed to save updated profile: %w", err)
	}

	// Update in-memory cache
	pm.profiles[name] = &profile.Config

	return nil
}

// ExportProfile exports a profile to a file
func (pm *ProfileManager) ExportProfile(name, exportPath string) error {
	config, exists := pm.profiles[name]
	if !exists {
		return fmt.Errorf("profile '%s' not found", name)
	}

	// Create export profile structure
	profile := &Profile{
		Name:        name,
		Description: fmt.Sprintf("Exported profile: %s", name),
		Created:     time.Now(),
		Modified:    time.Now(),
		Config:      *config.Clone(),
	}

	// Write to export file
	file, err := os.Create(exportPath)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(profile)
	if err != nil {
		return fmt.Errorf("failed to encode profile: %w", err)
	}

	return nil
}

// ImportProfile imports a profile from a file
func (pm *ProfileManager) ImportProfile(importPath string) error {
	file, err := os.Open(importPath)
	if err != nil {
		return fmt.Errorf("failed to open import file: %w", err)
	}
	defer file.Close()

	var profile Profile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&profile)
	if err != nil {
		return fmt.Errorf("failed to decode profile: %w", err)
	}

	// Validate imported profile
	profile.Config.Validate()

	// Check for name conflicts
	originalName := profile.Name
	counter := 1
	for _, exists := pm.profiles[profile.Name]; exists; _, exists = pm.profiles[profile.Name] {
		profile.Name = fmt.Sprintf("%s_%d", originalName, counter)
		counter++
	}

	// Save imported profile
	profile.Modified = time.Now()
	err = pm.saveProfile(&profile)
	if err != nil {
		return fmt.Errorf("failed to save imported profile: %w", err)
	}

	// Add to memory
	pm.profiles[profile.Name] = &profile.Config

	return nil
}

// GetCurrentProfile returns the name of the currently active profile
func (pm *ProfileManager) GetCurrentProfile() string {
	return pm.currentProfile
}

// CreatePresetProfiles creates commonly used configuration presets
func (pm *ProfileManager) CreatePresetProfiles() error {
	presets := map[string]*Profile{
		"gaming": {
			Name:        "gaming",
			Description: "Optimized for gaming with fast performance and enhanced features",
			Created:     time.Now(),
			Modified:    time.Now(),
			Config: UserConfig{
				MenuConfig: MenuConfig{
					Enabled:             true,
					DefaultROMDirectory: "",
					SortBy:              "name",
					SortAscending:       true,
					ShowMetadata:        true,
					ShowRecentROMs:      true,
					MaxRecentROMs:       50,
					EnableSearch:        true,
					CaseSensitiveSearch: false,
				},
				ThemeConfig: ThemeConfig{
					Name:            "Dark",
					BackgroundColor: Color{32, 32, 32, 255},
					TextColor:       Color{255, 255, 255, 255},
					SelectedColor:   Color{0, 128, 255, 255},
					BorderColor:     Color{64, 64, 64, 255},
					AccentColor:     Color{0, 255, 128, 255},
					FontSize:        14,
				},
				EmulatorConfig: EmulatorConfig{
					DefaultScale:       12,
					DefaultSpeed:       800,
					DefaultWindowTitle: "CHIP-8 Gaming",
					AudioEnabled:       true,
					AudioVolume:        0.8,
					AudioFrequency:     44100,
					BeepFrequency:      440.0,
				},
				CacheConfig: CacheConfig{
					MetadataCacheEnabled:  true,
					MetadataCacheTTL:      time.Minute * 5,
					DirectoryCacheEnabled: true,
					DirectoryCacheTTL:     time.Minute * 2,
					MaxCacheSize:          2000,
					CleanupInterval:       time.Minute * 10,
				},
				LastUpdated: time.Now(),
				Version:     "1.0.0",
			},
			Tags: []string{"gaming", "performance", "preset"},
		},

		"development": {
			Name:        "development",
			Description: "Development and testing configuration with debug features",
			Created:     time.Now(),
			Modified:    time.Now(),
			Config: UserConfig{
				MenuConfig: MenuConfig{
					Enabled:             true,
					DefaultROMDirectory: "",
					SortBy:              "date",
					SortAscending:       false,
					ShowMetadata:        true,
					ShowRecentROMs:      true,
					MaxRecentROMs:       100,
					EnableSearch:        true,
					CaseSensitiveSearch: false,
				},
				ThemeConfig: ThemeConfig{
					Name:            "Light",
					BackgroundColor: Color{240, 240, 240, 255},
					TextColor:       Color{0, 0, 0, 255},
					SelectedColor:   Color{0, 100, 200, 255},
					BorderColor:     Color{180, 180, 180, 255},
					AccentColor:     Color{255, 140, 0, 255},
					FontSize:        12,
				},
				EmulatorConfig: EmulatorConfig{
					DefaultScale:       8,
					DefaultSpeed:       500,
					DefaultWindowTitle: "CHIP-8 Development",
					AudioEnabled:       false,
					AudioVolume:        0.5,
					AudioFrequency:     22050,
					BeepFrequency:      440.0,
				},
				CacheConfig: CacheConfig{
					MetadataCacheEnabled:  false,
					MetadataCacheTTL:      time.Minute * 1,
					DirectoryCacheEnabled: false,
					DirectoryCacheTTL:     time.Second * 30,
					MaxCacheSize:          500,
					CleanupInterval:       time.Minute * 5,
				},
				LastUpdated: time.Now(),
				Version:     "1.0.0",
			},
			Tags: []string{"development", "debug", "preset"},
		},

		"performance": {
			Name:        "performance",
			Description: "Optimized for maximum performance and minimal resource usage",
			Created:     time.Now(),
			Modified:    time.Now(),
			Config: UserConfig{
				MenuConfig: MenuConfig{
					Enabled:             true,
					DefaultROMDirectory: "",
					SortBy:              "name",
					SortAscending:       true,
					ShowMetadata:        false,
					ShowRecentROMs:      true,
					MaxRecentROMs:       10,
					EnableSearch:        true,
					CaseSensitiveSearch: false,
				},
				ThemeConfig: ThemeConfig{
					Name:            "Classic",
					BackgroundColor: Color{0, 0, 0, 255},
					TextColor:       Color{0, 255, 0, 255},
					SelectedColor:   Color{0, 255, 255, 255},
					BorderColor:     Color{128, 128, 128, 255},
					AccentColor:     Color{255, 255, 0, 255},
					FontSize:        12,
				},
				EmulatorConfig: EmulatorConfig{
					DefaultScale:       8,
					DefaultSpeed:       600,
					DefaultWindowTitle: "CHIP-8 Performance",
					AudioEnabled:       true,
					AudioVolume:        0.6,
					AudioFrequency:     22050,
					BeepFrequency:      440.0,
				},
				CacheConfig: CacheConfig{
					MetadataCacheEnabled:  true,
					MetadataCacheTTL:      time.Minute * 10,
					DirectoryCacheEnabled: true,
					DirectoryCacheTTL:     time.Minute * 5,
					MaxCacheSize:          500,
					CleanupInterval:       time.Minute * 15,
				},
				LastUpdated: time.Now(),
				Version:     "1.0.0",
			},
			Tags: []string{"performance", "minimal", "preset"},
		},
	}

	for name, profile := range presets {
		// Check if preset already exists
		if _, exists := pm.profiles[name]; exists {
			continue // Skip existing presets
		}

		err := pm.saveProfile(profile)
		if err != nil {
			return fmt.Errorf("failed to create preset profile '%s': %w", name, err)
		}

		pm.profiles[name] = &profile.Config
	}

	return nil
}

// loadProfiles loads all profiles from the profiles directory
func (pm *ProfileManager) loadProfiles() {
	entries, err := os.ReadDir(pm.profilesDir)
	if err != nil {
		return // Profiles directory doesn't exist or can't be read
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			profileName := strings.TrimSuffix(entry.Name(), ".json")
			profile, err := pm.loadProfileFromFile(profileName)
			if err == nil {
				pm.profiles[profileName] = &profile.Config
			}
		}
	}
}

// loadProfileFromFile loads a single profile from file
func (pm *ProfileManager) loadProfileFromFile(name string) (*Profile, error) {
	profilePath := filepath.Join(pm.profilesDir, name+".json")

	file, err := os.Open(profilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var profile Profile
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&profile)
	if err != nil {
		return nil, err
	}

	// Validate loaded configuration
	profile.Config.Validate()

	return &profile, nil
}

// saveProfile saves a profile to file
func (pm *ProfileManager) saveProfile(profile *Profile) error {
	profilePath := filepath.Join(pm.profilesDir, profile.Name+".json")

	// Create temporary file for atomic write
	tempPath := profilePath + ".tmp"
	file, err := os.Create(tempPath)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(profile)
	file.Close()

	if err != nil {
		os.Remove(tempPath)
		return err
	}

	// Atomic rename
	return os.Rename(tempPath, profilePath)
}

// isValidProfileName validates a profile name
func isValidProfileName(name string) bool {
	if len(name) == 0 || len(name) > 50 {
		return false
	}

	// Check for valid characters (alphanumeric, underscore, hyphen)
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '_' || char == '-') {
			return false
		}
	}

	// Check for reserved names
	reserved := []string{"default", "backup", "temp", "config"}
	for _, reservedName := range reserved {
		if strings.ToLower(name) == reservedName {
			return false
		}
	}

	return true
}

// GetProfileInfo returns detailed information about a profile
func (pm *ProfileManager) GetProfileInfo(name string) (*Profile, error) {
	return pm.loadProfileFromFile(name)
}

// RenameProfile renames an existing profile
func (pm *ProfileManager) RenameProfile(oldName, newName string) error {
	if !isValidProfileName(newName) {
		return fmt.Errorf("invalid new profile name: %s", newName)
	}

	if _, exists := pm.profiles[newName]; exists {
		return fmt.Errorf("profile '%s' already exists", newName)
	}

	// Load existing profile
	profile, err := pm.loadProfileFromFile(oldName)
	if err != nil {
		return fmt.Errorf("failed to load profile: %w", err)
	}

	// Update profile name and save with new name
	profile.Name = newName
	profile.Modified = time.Now()

	err = pm.saveProfile(profile)
	if err != nil {
		return fmt.Errorf("failed to save renamed profile: %w", err)
	}

	// Remove old profile file
	oldPath := filepath.Join(pm.profilesDir, oldName+".json")
	os.Remove(oldPath)

	// Update in-memory cache
	pm.profiles[newName] = pm.profiles[oldName]
	delete(pm.profiles, oldName)

	// Update current profile if it was renamed
	if pm.currentProfile == oldName {
		pm.currentProfile = newName
	}

	return nil
}

// GetProfilesDirectory returns the profiles directory path
func (pm *ProfileManager) GetProfilesDirectory() string {
	return pm.profilesDir
}
