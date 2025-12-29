package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppSettings represents application settings stored on disk
type AppSettings struct {
	DeveloperMode bool   `json:"developerMode"`
	Theme         string `json:"theme"`
	ItemsPerPage  int    `json:"itemsPerPage"`
}

// DefaultSettings returns default application settings
func DefaultSettings() *AppSettings {
	return &AppSettings{
		DeveloperMode: false,
		Theme:         "dark",
		ItemsPerPage:  25,
	}
}

// getSettingsPath returns the path to the settings file
func getSettingsPath() (string, error) {
	// Use the same path as data: %APPDATA%\AutoManagement
	dataPath, err := getDataPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(dataPath, "settings.json"), nil
}

// LoadSettings loads settings from disk
func LoadSettings() (*AppSettings, error) {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return DefaultSettings(), nil
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		// File doesn't exist, return defaults
		return DefaultSettings(), nil
	}

	var settings AppSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return DefaultSettings(), nil
	}

	return &settings, nil
}

// SaveSettings saves settings to disk
func SaveSettings(settings *AppSettings) error {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
}

// UpdateDeveloperMode updates only the developer mode setting
func UpdateDeveloperMode(enabled bool) error {
	settings, _ := LoadSettings()
	settings.DeveloperMode = enabled
	return SaveSettings(settings)
}

// IsDeveloperMode checks if developer mode is enabled
func IsDeveloperMode() bool {
	settings, _ := LoadSettings()
	return settings.DeveloperMode
}
