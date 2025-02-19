package usecase

import (
	"encoding/json"
	"os"
	"path/filepath"

	"main/internal/domain"
)

type SettingsManager struct {
	settings     domain.Settings
	settingsPath string
}

func NewSettingsManager() *SettingsManager {
	homeDir, _ := os.UserHomeDir()
	return &SettingsManager{
		settings: domain.Settings{
			ThresholdSeconds: 15,
		},
		settingsPath: filepath.Join(homeDir, ".workpulse", "settings.json"),
	}
}

func (m *SettingsManager) LoadSettings() error {
	data, err := os.ReadFile(m.settingsPath)
	if os.IsNotExist(err) {
		return m.SaveSettings()
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &m.settings)
}

func (m *SettingsManager) SaveSettings() error {
	dir := filepath.Dir(m.settingsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(m.settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.settingsPath, data, 0644)
}

func (m *SettingsManager) GetSettings() domain.Settings {
	return m.settings
}

func (m *SettingsManager) UpdateSettings(settings domain.Settings) error {
	m.settings = settings
	return m.SaveSettings()
} 